/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-24 11:46
* Description:
*****************************************************************/

package rpcPoint

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/sync/xsafeMap"
	"io"
	"time"
)

type tClientCallInfo struct {
	err    error
	data   []byte
	status chan int
}

func newClientCallInfo() *tClientCallInfo {
	return &tClientCallInfo{
		err:    nil,
		status: make(chan int, 1),
	}
}

func (p *tClientCallInfo) Wait(timeout time.Duration) (data []byte, err error) {
	time.AfterFunc(timeout, func() {
		select {
		case <-p.status:
			return
		default:
			p.err = errors.New("请求超时1")
			p.status <- 1
		}
	})
	select {
	case n := <-p.status:
		if n > 0 {
			return nil, p.err
		}
		return p.data, nil
	}
}

func (p *tClientCallInfo) CallResult(data []byte, err error) {
	select {
	case <-p.status:
		return
	default:
	}
	if err != nil {
		p.err = err
		p.status <- 1
	} else {
		p.data = data
		p.err = nil
		p.status <- 0
	}
}

type tClientCallData struct {
	method string
	seqId  int32
	data   []byte
}

func newClientCallData(method string, seqId int32, data []byte) *tClientCallData {
	return &tClientCallData{
		method: method,
		seqId:  seqId,
		data:   data,
	}
}

type tInnerClient struct {
	sckt           *thrift.TSocket
	pool           *tHostClientPool
	host           string
	port           int
	trans          thrift.TTransport
	inProto        thrift.TProtocol
	outProto       thrift.TProtocol
	writeTimeout   time.Duration
	readTimeout    time.Duration
	connectTimeout time.Duration
	expire         time.Time
	lastAlive      time.Time
	heartFailCount int
	sendBuf        chan *tClientCallData
	close          chan byte
	callResults    *xsafeMap.TStrAnyMap
	protoFactory   thrift.TProtocolFactory
	heartbeatLoss  int
	sendLayout     int
}

func newInnerClient(pool *tHostClientPool, protoFac thrift.TProtocolFactory, host string, port int, writeTimeout, readTimeout, connectTimeout time.Duration, heartbeatLoss int) (*tInnerClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	sckt, err := thrift.NewTSocket(addr)
	if err != nil {
		return nil, err
	}
	frameTrans := thrift.NewTFramedTransport(sckt)
	if err := sckt.SetTimeout(connectTimeout); err != nil {
		return nil, err
	}
	if err := sckt.Open(); err != nil {
		return nil, err
	}
	inst := &tInnerClient{
		sckt:           sckt,
		pool:           pool,
		trans:          frameTrans,
		host:           host,
		port:           port,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		connectTimeout: connectTimeout,
		protoFactory:   protoFac,
		inProto:        protoFac.GetProtocol(frameTrans),
		outProto:       protoFac.GetProtocol(sckt),
		sendBuf:        make(chan *tClientCallData, 256),
		close:          make(chan byte, 1),
		callResults:    xsafeMap.NewStrAnyMap(),
		heartbeatLoss:  heartbeatLoss,
		sendLayout:     0,
	}
	go func() {
		if err := inst.sendLoop(); err != nil {
			xlog.Error(err)
		}
	}()
	go func() {
		if err := inst.readLoop(); err != nil {
			xlog.Error(err)
		}
	}()
	return inst, nil
}

func (p *tInnerClient) IsOpen() bool {
	select {
	case <-p.close:
		return false
	default:
		return p.trans != nil
	}
}

func (p *tInnerClient) Close() error {
	select {
	case <-p.close:
		return nil
	default:
	}
	close(p.close)
	keys := p.callResults.Keys()
	for _, k := range keys {
		c := p.callResults.Get(k).(*tClientCallInfo)
		c.CallResult(nil, errors.New("断开连接"))
		p.callResults.Remove(k)
	}

	if p.trans != nil && p.trans.IsOpen() {
		err := p.trans.Close()
		p.trans = nil
		if err != nil {
			return err
		}
	}
	fmt.Println("断开连接")
	return nil
}

func (p *tInnerClient) sendLoop() error {
	ctx := context.Background()
	for {
		select {
		case data, ok := <-p.sendBuf:
			if !ok {
				return nil
			}
			if err := p.doSend(ctx, data.data); err != nil {
				key := fmt.Sprintf("%s_%d", data.method, data.seqId)
				if v := p.callResults.Get(key); v != nil {
					call := v.(*tClientCallInfo)
					call.CallResult(nil, err)
					p.callResults.Remove(key)
				}
			}
			break
		case <-p.close:
			return nil
		}
	}
}

func (p *tInnerClient) sendData(cxt context.Context, method string, seqId int32, data []byte) ([]byte, error) {
	call := newClientCallInfo()
	key := fmt.Sprintf("%s_%d", method, seqId)
	xlog.Debug("发送数据 key:", key)
	p.callResults.Set(key, call)

	sendData := newClientCallData(method, seqId, data)
	select {
	case <-p.close:
		return nil, errors.New("已经断开连接")
	default:
		p.sendBuf <- sendData
	}
	return call.Wait(p.readTimeout)
}

func (p *tInnerClient) readLoop() error {
	ctx := context.Background()
	buf := thrift.NewTMemoryBuffer()
	frameTrans := thrift.NewTFramedTransport(buf)
	out := p.protoFactory.GetProtocol(frameTrans)

	for {
		select {
		case <-p.close:
			return nil
		default:
		}

		buf.Truncate(0)
		if err := p.sckt.SetTimeout(p.readTimeout); err != nil {
			xlog.Debug("set read timeout error:", err)
		}
		msgName, msgType, seqId, err := ProtocolTransform(p.inProto, out)
		if err != nil {
			if err.Error() == io.EOF.Error() {
				xlog.Debug("服务端断开连接:", err)
				return p.Close()
			}
			if p.sendLayout > 0 {
				if p.heartFailCount > p.heartbeatLoss {
					return p.Close()
				}
				p.heartFailCount++
			}
			continue
		}
		if err := out.Flush(ctx); err != nil {
			xlog.Error(err)
		}
		if msgType == thrift.ONEWAY+1 {
			p.sendLayout--
			p.heartFailCount = 0
			p.lastAlive = time.Now()
			// heartbeat包
			xlog.Debug("====> 收到收跳包")
			continue
		}
		if msgType == thrift.REPLY || msgType == thrift.EXCEPTION {
			p.lastAlive = time.Now()
			p.heartFailCount = 0
			p.sendLayout--
			// 更新连接空闲到期时间
			p.expire = time.Now().Add(p.pool.keepAlive)
			key := fmt.Sprintf("%s_%d", msgName, seqId)
			xlog.Debug("返回数据 key:", key)
			if v := p.callResults.Get(key); v != nil {
				call := v.(*tClientCallInfo)
				call.CallResult(buf.Bytes(), err)
				p.callResults.Remove(key)
			}
		}
	}
}

func (p *tInnerClient) sendHeartbeat() {
	ctx := context.Background()
	buf := thrift.NewTMemoryBuffer()
	frameBuf := thrift.NewTFramedTransport(buf)
	out := p.protoFactory.GetProtocol(frameBuf)
	if err := out.WriteMessageBegin("heartbeat", thrift.ONEWAY+1, -1); err != nil {
		xlog.Error(err)
	}
	if err := out.WriteMessageEnd(); err != nil {
		xlog.Error(err)
	}
	if err := out.Flush(ctx); err != nil {
		xlog.Error(err)
	}
	if err := p.doSend(ctx, buf.Bytes()); err != nil {
		xlog.Error(err)
	}
	xlog.Debug("发送心跳包")
}

func (p *tInnerClient) doSend(cxt context.Context, data []byte) error {
	if !p.trans.IsOpen() {
		return errors.New("连接超时或已断开")
	}
	if err := p.sckt.SetTimeout(p.writeTimeout); err != nil {
		return err
	}
	if _, err := p.outProto.Transport().Write(data); err != nil {
		return err
	}
	if err := p.outProto.Transport().Flush(cxt); err != nil {
		return err
	}
	p.sendLayout++
	xlog.Debug("数据发送成功.")
	return nil
}

func (p *tInnerClient) Call(cxt context.Context, method string, seqId int32, data []byte) (result []byte, err error) {
	result, err = p.sendData(cxt, method, seqId, data)
	if err == nil {
		p.pool.Put(p)
	}
	return
}
