/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 10:30
* Description:
*****************************************************************/

package netstream

import (
	"errors"
	"github.com/go-xe2/x/core/logger"
	"github.com/go-xe2/x/os/xlog"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type StreamConn interface {
	Id() string
	// 发送数据
	Send(data []byte) (int, error)
	// 发送并等待数据返回
	Call(data []byte) ([]byte, error)
	SendTo(targetId string, data []byte) (int, error)
	CallTo(targetId string, data []byte) ([]byte, error)
	// 发送请求
	Request(reqId string, namespace string, body []byte) error
	// 请求回复
	Response(reqId string, body []byte) error

	SendHeartbeat() error
	Disconnect()
	UpdateHeartbeat(isLoss bool)
	HeartbeatLossCount() int
	Heartbeat() time.Time
	Close() error
}

type StreamConnHandler interface {
	OnRecv(conn StreamConn, data []byte)
	OnCall(conn StreamConn, data []byte) ([]byte, error)
	// 收到消息并回复
	OnSendTo(conn StreamConn, toConn string, data []byte)
	OnCallTo(conn StreamConn, toConn string, data []byte) ([]byte, error)
	// 收到请求
	OnRequest(reqConn StreamConn, reqId string, namespace string, body []byte)
	// 收到请求回复
	OnResponse(resConn StreamConn, reqId string, body []byte)
	OnHeartbeat(conn StreamConn)
	OnDisconnect(conn StreamConn)
	OnConnect(conn StreamConn)
	OnReady(conn StreamConn)
}

type streamConnBuf struct {
	buf     []byte
	timeout time.Duration
}

type streamReadBuf struct {
	cmd  StreamCmdType
	data []byte
}

type tStreamConn struct {
	wg   sync.WaitGroup
	id   string
	conn net.Conn
	addr net.Addr
	// 待发送队列
	sendBuf chan *streamConnBuf
	// 待读取队列
	recvBuf chan *streamReadBuf
	// 接收数据缓存区
	handler       StreamConnHandler
	close         chan byte
	heartbeatLoss int32
	heartbeat     int64
	options       *TConnOptions
	dataCache     []byte
}

var _ StreamConn = (*tStreamConn)(nil)

func newStreamConn(handler StreamConnHandler, addr net.Addr, options *TConnOptions) *tStreamConn {
	if options == nil {
		options = DefaultConnOptions
	}
	return &tStreamConn{
		addr:          addr,
		id:            MakeConnId(),
		handler:       handler,
		heartbeatLoss: 0,
		heartbeat:     0,
		options:       options,
	}
}

func newStreamConnByConn(handler StreamConnHandler, conn net.Conn, options *TConnOptions) *tStreamConn {
	if options == nil {
		options = DefaultConnOptions
	}
	inst := &tStreamConn{
		options:       options,
		conn:          conn,
		addr:          conn.RemoteAddr(),
		id:            MakeConnId(),
		handler:       handler,
		heartbeatLoss: 0,
		heartbeat:     0,
	}
	inst.initChannels()
	return inst
}

func (p *tStreamConn) log(level logger.LogLevel, args ...interface{}) {
	Log("tStreamConn", p.options.GetLogger(), level, args...)
}

func (p *tStreamConn) SetHandler(handler StreamConnHandler) {
	p.handler = handler
}

func (p *tStreamConn) Id() string {
	return p.id
}

// 发送数据
func (p *tStreamConn) Send(data []byte) (int, error) {
	return p.send(StreamCmdSend, data)
}

func (p *tStreamConn) Call(data []byte) ([]byte, error) {
	seqId := MakeSendSeqId()
	callData := packeCallData(seqId, data)
	_, err := p.send(StreamCmdCall, callData)
	if err != nil {
		return nil, err
	}
	return sendReplies.GetResult(seqId, p.options.GetReadTimeout())
}

func (p *tStreamConn) SendTo(targetId string, data []byte) (int, error) {
	sendData := packeSendToData(targetId, data)
	return p.send(StreamCmdSendTo, sendData)
}

func (p *tStreamConn) CallTo(targetId string, data []byte) ([]byte, error) {
	seqId := MakeSendSeqId()
	callData := packeCallToData(seqId, targetId, data)
	if _, err := p.send(StreamCmdCallTo, callData); err != nil {
		return nil, err
	}
	return sendReplies.GetResult(seqId, p.options.GetReadTimeout())
}

func (p *tStreamConn) Request(reqId string, namespace string, body []byte) error {
	data := packeRequestData(reqId, namespace, body)
	if _, err := p.send(StreamCmdRequest, data); err != nil {
		return err
	}
	return nil
}

func (p *tStreamConn) Response(reqId string, body []byte) error {
	data := packeResponseData(reqId, body)
	if _, err := p.send(StreamCmdResponse, data); err != nil {
		return err
	}
	return nil
}

func (p *tStreamConn) UpdateHeartbeat(isLoss bool) {
	if isLoss {
		atomic.AddInt32(&p.heartbeatLoss, 1)
		atomic.StoreInt64(&p.heartbeat, time.Now().Unix())
	} else {
		atomic.StoreInt32(&p.heartbeatLoss, 0)
		atomic.StoreInt64(&p.heartbeat, time.Now().Unix())
	}
}

func (p *tStreamConn) HeartbeatLossCount() int {
	n := int(atomic.LoadInt32(&p.heartbeatLoss))
	return n
}

func (p *tStreamConn) Heartbeat() time.Time {
	t := atomic.LoadInt64(&p.heartbeat)
	return time.Unix(t, 0)
}

func (p *tStreamConn) SendHeartbeat() error {
	xlog.Debug("发送收跳====>>")
	if _, err := p.send(StreamCmdHeartbeat, nil); err != nil {
		return err
	}
	return nil
}

func (p *tStreamConn) SendClientId() error {
	data := packeClientId(p.id)
	if _, err := p.send(StreamCmdClientId, data); err != nil {
		return err
	}
	p.log(logger.LEVEL_DEBU, "发送客户端ID:", p.id, "成功")
	return nil
}

func (p *tStreamConn) pushDeadline(read, write bool) {
	rTimeout := p.options.GetReadTimeout()
	wTimeout := p.options.GetWriteTimeout()
	if read && write {
		if rTimeout <= 0 && wTimeout <= 0 {
			return
		}
		if rTimeout <= 0 {
			rTimeout = wTimeout
		}
		if wTimeout <= 0 {
			wTimeout = rTimeout
		}
		rt := time.Now().Add(time.Duration(rTimeout))
		wt := time.Now().Add(time.Duration(wTimeout))
		_ = p.conn.SetReadDeadline(rt)
		_ = p.conn.SetWriteDeadline(wt)
	} else if read {
		if rTimeout <= 0 {
			return
		}
		t := time.Now().Add(rTimeout)
		_ = p.conn.SetReadDeadline(t)
	} else if write {
		if wTimeout <= 0 {
			return
		}
		t := time.Now().Add(wTimeout)
		_ = p.conn.SetWriteDeadline(t)
	}
}

func (p *tStreamConn) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if e := p.sendProcessLoop(); e != nil {
			p.log(logger.LEVEL_DEBU, "发送数据协程处理结束，出错:", e)
		}
		p.closeLoop()
		p.log(logger.LEVEL_DEBU, "1.sendProcessLoop结束")
	}()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if e := p.recvProcessLoop(); e != nil {
			p.log(logger.LEVEL_DEBU, "接收数据协程处理结束，出错:", e)
		}
		p.closeLoop()
		p.log(logger.LEVEL_DEBU, "2.recvProcessLoop结束")
	}()
	p.wg.Add(1)
	go func() {
		p.wg.Done()
		if e := p.processReadLoop(); e != nil {
			p.log(logger.LEVEL_DEBU, "数据读取处理协程结束，出错:", e)
		}
		p.closeLoop()
		p.log(logger.LEVEL_DEBU, "3.processReadLoop结束")
	}()
	go func() {
		p.wg.Wait()
		// 触发断线事件
		p.doOnDisconnect()
	}()
	xlog.Debug("服务处理协程启动完成.")
}

func (p *tStreamConn) WaitStop() {
	p.wg.Wait()
}

func (p *tStreamConn) initChannels() {
	p.sendBuf = make(chan *streamConnBuf, p.options.GetSendBufSize())
	p.recvBuf = make(chan *streamReadBuf, p.options.GetRecvBufSize())
	p.close = make(chan byte, 1)
	p.dataCache = make([]byte, 0)
}

func (p *tStreamConn) Open() error {
	if p.IsOpen() {
		return errors.New("已经打开连接")
	}
	if p.addr == nil {
		return errors.New("未设置连接地址")
	}
	if len(p.addr.Network()) == 0 {
		return errors.New("未设置网络协议名称")
	}
	if len(p.addr.String()) == 0 {
		return errors.New("网络地址无效")
	}
	var err error
	if p.conn, err = net.DialTimeout(p.addr.Network(), p.addr.String(), p.options.GetConnectTimeout()); err != nil {
		return errors.New(err.Error())
	}
	p.initChannels()
	p.Start()
	p.doOnConnect()
	return nil
}

func (p *tStreamConn) ReOpen() error {
	p.log(logger.LEVEL_DEBU, "准备重试连接")
	defer func() {
		p.log(logger.LEVEL_DEBU, "重试连接完成")
	}()
	p.wg.Wait()
	if p.IsOpen() {
		p.innerClose()
	}
	var err error
	if p.conn, err = net.DialTimeout(p.addr.Network(), p.addr.String(), p.options.GetConnectTimeout()); err != nil {
		return errors.New(err.Error())
	}
	//p.initChannels()
	p.close = make(chan byte, 1)
	p.Start()
	p.doOnConnect()
	return nil
}

func (p *tStreamConn) send(cmd StreamCmdType, data []byte) (int, error) {
	// 已经关闭连接，不再发送数据
	select {
	case <-p.close:
		// 已经关闭服务，则关闭发送队列
		close(p.sendBuf)
		return 0, errors.New("已关闭连接")
	default:
	}
	packet := packeData(data, cmd)
	buf := &streamConnBuf{
		buf:     packet,
		timeout: p.options.GetWriteTimeout(),
	}
	select {
	case p.sendBuf <- buf:
		// 发送队列区如果已满，等待发送直到有空位或关闭服务
		return len(data), nil
	case <-p.close:
		// 服务已关闭，关闭发送消息队列
		close(p.sendBuf)
		return 0, errors.New("连接已关闭")
	}
}

func (p *tStreamConn) sendProcessLoop() error {
	p.pushDeadline(false, true)
	for {
		var buf *streamConnBuf
		var ok bool
		select {
		case <-p.close:
			return nil
		case buf, ok = <-p.sendBuf:
		}
		if !ok {
			// 发送队列已关闭
			return nil
		}
		if buf == nil {
			continue
		}
		if _, e := p.write(buf.buf); e != nil {
			p.log(logger.LEVEL_WARN, "发送数据出错:", e)
			continue
		}
	}
}

// 数据接收处理过程,recvBuf的生产者
func (p *tStreamConn) recvProcessLoop() error {
	buf := make([]byte, 4096)
	p.dataCache = make([]byte, 0)
	p.pushDeadline(true, false)
	for {
		select {
		case <-p.close:
			// 已关闭服务，关闭接收数据队列
			close(p.recvBuf)
			return nil
		default:
			// 进入默认接收数据处理过程
		}
		n, err := p.read(buf)
		if err != nil {
			if err == io.EOF {
				// 连接已经断开
				return io.EOF
			} else if netErr, ok := err.(net.Error); ok {
				xlog.Debug("recvProcessLoop read timeout")
				if netErr.Timeout() {
					continue
				}
				// 连接已经断开
				return io.EOF
			}
			p.log(logger.LEVEL_WARN, "接收数据协程读取数据出错:", err)
			continue
		}
		p.UpdateHeartbeat(false)
		if n > 0 {
			p.dataCache = append(p.dataCache, buf[:n]...)
		}
		xlog.Debug("recv process recv len:", n)
		// 解码所有包直到剩余数据不足一个包头为止
		for len(p.dataCache) >= PacketHeadSize {
			ready, data, cmd, remainder := unpackeData(p.dataCache)
			xlog.Debug("recv process ready:", ready, ", remainder len:", len(remainder))
			if !ready {
				// 数据未准备好，跳出当前过程
				break
			}
			p.dataCache = remainder
			readBuf := &streamReadBuf{
				cmd:  cmd,
				data: data,
			}
			select {
			case p.recvBuf <- readBuf:
				// 如果接通收缓存区数据已满，则等待有空位或服务关闭
				xlog.Debug("消息数据放入队列 cmd:", cmd, ", data size:", len(data))
				continue
			case <-p.close:
				// 服务已经关闭，则关闭接收列队列
				close(p.recvBuf)
				return nil
			}
		}
	}
}

// 数据读取处理过程，recvBuf的消费者
func (p *tStreamConn) processReadLoop() error {
	for {
		var buf *streamReadBuf
		var ok bool
		select {
		case <-p.close:
			return nil
		case buf, ok = <-p.recvBuf:
		}
		if !ok {
			// 队列已经关闭
			return nil
		}
		if buf == nil {
			continue
		}
		go func() {
			if e := p.doRecv(buf.cmd, buf.data); e != nil {
				p.log(logger.LEVEL_WARN, "读取数据处理出错:", e)
			}
		}()
	}
}

func (p *tStreamConn) doRecv(cmd StreamCmdType, data []byte) error {
	switch cmd {
	case StreamCmdHeartbeat:
		p.doOnHeartbeat()
		break
	case StreamCmdConnect:
		p.doOnConnect()
		break
	case StreamCmdDisconnect:
		// 收到当前当前连接通知消息,准备断开当前连接
		select {
		case <-p.close:
			// 已经关闭服务
			return nil
		default:
			// 关闭服务
			go p.innerClose()
			return nil
		}
	case StreamCmdSend:
		return p.doRecvSend(data)
	case StreamCmdCall:
		return p.doRecvCall(data)
	case StreamCmdCallReply:
		return p.doRecvCallReply(data)
	case StreamCmdSendTo:
		return p.doRecvSendTo(data)
	case StreamCmdCallTo:
		return p.doRecvCallTo(data)
	case StreamCmdClientId:
		return p.doRecvId(data)
	case StreamCmdClientReady:
		return p.DoRecvClientReady(data)
	case StreamCmdResponse:
		return p.doRecvResponse(data)
	case StreamCmdRequest:
		return p.doRecvRequest(data)
	}
	return nil
}

// 发送客户端已准备就绪消息
func (p *tStreamConn) SendReady() error {
	if _, err := p.send(StreamCmdClientReady, nil); err != nil {
		return err
	}
	p.log(logger.LEVEL_DEBU, "发送就绪消息成功")
	return nil
}

// 收到客户端准备就绪消息， 发送客户id给客户端
func (p *tStreamConn) DoRecvClientReady(data []byte) error {
	// 服务端方法
	xlog.Debug("收到客户端准备就绪消息")
	if p.handler != nil {
		p.handler.OnReady(p)
	}
	return p.SendClientId()
}

func (p *tStreamConn) doRecvId(data []byte) error {
	szId := unpackeClientId(data)
	if szId != "" {
		p.id = szId
		if p.handler != nil {
			// 准备就续
			p.handler.OnReady(p)
		}
		p.log(logger.LEVEL_DEBU, "收到客户端ID:", szId)
	}
	return nil
}

func (p *tStreamConn) doRecvRequest(data []byte) error {
	if p.handler != nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					p.log(logger.LEVEL_WARN, "处理请求时出错:", e)
				}
				reqId, namespace, body := unpackRequestData(data)
				if reqId == "" {
					return
				}
				p.handler.OnRequest(p, reqId, namespace, body)
			}()
		}()
	}
	return nil
}

func (p *tStreamConn) doRecvResponse(data []byte) error {
	if p.handler != nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					p.log(logger.LEVEL_WARN, "处理回复出错:", e)
				}
				reqId, body := unpackResponseData(data)
				if reqId == "" {
					return
				}
				p.handler.OnResponse(p, reqId, body)
			}()
		}()
	}
	return nil
}

func (p *tStreamConn) doRecvSendTo(data []byte) error {
	targetId, buf := unpackeSendToData(data)
	return p.doSendTo(targetId, buf)
}

func (p *tStreamConn) doRecvCall(data []byte) error {
	seqId, buf := unpackeCallData(data)
	if seqId == 0 {
		return nil
	}
	if p.handler != nil {
		ret, err := p.handler.OnCall(p, buf)
		replyData := packeCallReplyData(seqId, ret, err)
		if _, err := p.send(StreamCmdCallReply, replyData); err != nil {
			return err
		}
		return nil
	}
	// 默认返回空数据
	replyData := packeCallReplyData(seqId, nil, nil)
	if _, err := p.send(StreamCmdCallReply, replyData); err != nil {
		return err
	}
	return nil
}

func (p *tStreamConn) doRecvCallReply(data []byte) error {
	seqId, rpData, err := unpackCallReplyData(data)
	if seqId == 0 {
		return nil
	}
	if p.handler != nil {
		sendReplies.Reply(seqId, rpData, err)
		return nil
	}
	return nil
}

func (p *tStreamConn) doRecvCallTo(data []byte) error {
	seqId, targetId, buf := unpackeCallToData(data)
	if seqId == 0 {
		return nil
	}
	if p.handler != nil {
		ret, err := p.handler.OnCallTo(p, targetId, buf)
		replyData := packeCallReplyData(seqId, ret, err)
		if _, err := p.send(StreamCmdCallReply, replyData); err != nil {
			return err
		}
		return nil
	}
	// 默认处理方式
	replyData := packeCallReplyData(seqId, nil, nil)
	if _, err := p.send(StreamCmdCallReply, replyData); err != nil {
		return err
	}
	return nil
}

func (p *tStreamConn) Conn() net.Conn {
	return p.conn
}

func (p *tStreamConn) isStop() bool {
	select {
	case <-p.close:
		return true
	default:
	}
	return false
}

// Returns true if the connection is open
func (p *tStreamConn) IsOpen() bool {
	if p.conn == nil {
		return false
	}
	return !p.isStop()
}

func (p *tStreamConn) doRecvSend(data []byte) error {
	if p.handler != nil {
		p.handler.OnRecv(p, data)
	}
	return nil
}

func (p *tStreamConn) doSendTo(targetId string, data []byte) error {
	if p.handler != nil {
		p.handler.OnSendTo(p, targetId, data)
	}
	return nil
}

func (p *tStreamConn) doOnDisconnect() {
	if p.handler != nil {
		p.handler.OnDisconnect(p)
	}
}

func (p *tStreamConn) doOnConnect() {
	if p.handler != nil {
		p.handler.OnConnect(p)
	}
}

func (p *tStreamConn) doOnHeartbeat() {
	if p.handler != nil {
		p.handler.OnHeartbeat(p)
	}
}

func (p *tStreamConn) Disconnect() {
	select {
	case <-p.close:
		// 已关闭服务，不发送数据
		return
	default:
	}
	p.pushDeadline(false, true)
	data := packeData(nil, StreamCmdDisconnect)
	if _, err := p.write(data); err != nil {
		p.log(logger.LEVEL_WARN, "发送断线消息出错:", err)
	}
}

func (p *tStreamConn) closeLoop() {
	select {
	case <-p.close:
		// 已经关闭服务
		return
	default:
		close(p.close)
	}
}

func (p *tStreamConn) doClose(closer func() error) (err error) {
	// 这里需要改成直接发送disconnect消息
	p.closeLoop()
	p.WaitStop()
	if closer != nil {
		if e := closer(); e != nil {
			return e
		}
	}
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			p.log(logger.LEVEL_WARN, "关闭连接出错:", err)
		}
		p.conn = nil
	}
	xlog.Debug("doClose before doOnDisconnect.")
	p.doOnDisconnect()
	xlog.Debug("doClose after doOnDisconnect.")
	return nil
}

// 内部使用的关闭方法，不发送disconnect消息,不调用waitStop
func (p *tStreamConn) innerClose() {
	p.closeLoop()
	p.WaitStop()
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			p.log(logger.LEVEL_WARN, "关闭连接出错:", err)
		}
		p.conn = nil
	}
}

// 当前连接，在连接关闭之前先发送disconnect消息
func (p *tStreamConn) Close() error {
	return p.doClose(func() error {
		p.Disconnect()
		return nil
	})
}

func (p *tStreamConn) Addr() net.Addr {
	return p.addr
}

func (p *tStreamConn) RemoteAddr() net.Addr {
	if !p.IsOpen() {
		return nil
	}
	return p.conn.RemoteAddr()
}

func (p *tStreamConn) read(buf []byte) (int, error) {
	if !p.IsOpen() {
		return 0, errors.New("未打开连接")
	}
	n, err := p.conn.Read(buf)
	return n, err
}

func (p *tStreamConn) write(buf []byte) (int, error) {
	if !p.IsOpen() {
		return 0, errors.New("未打开连接")
	}
	return p.conn.Write(buf)
}
