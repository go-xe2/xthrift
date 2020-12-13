/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 11:28
* Description:
*****************************************************************/

package netstream

import (
	"errors"
	"github.com/go-xe2/x/core/logger"
	"github.com/go-xe2/x/sync/xsafeMap"
	"net"
	"sync"
	"sync/atomic"
)

type TStreamServer struct {
	wg            sync.WaitGroup
	listener      net.Listener
	addr          net.Addr
	heartbeatRun  int32
	clients       *xsafeMap.TStrAnyMap
	handler       ServerStreamHandler
	isListener    int32
	isInterrupted int32
	closed        chan byte
	options       *TStmServerOptions
	requests      *xsafeMap.TStrAnyMap
}

var _ StreamConnHandler = (*TStreamServer)(nil)

func NewStreamServer(listenAddr string, options *TStmServerOptions) (*TStreamServer, error) {
	if options == nil {
		options = DefaultStmServerOptions
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &TStreamServer{
		options:      options,
		addr:         addr,
		clients:      xsafeMap.NewStrAnyMap(),
		heartbeatRun: 0,
		closed:       make(chan byte, 1),
		requests:     xsafeMap.NewStrAnyMap(),
	}, nil
}

func NewStreamServerByListener(listener net.Listener, options *TStmServerOptions) *TStreamServer {
	if options == nil {
		options = DefaultStmServerOptions
	}
	return &TStreamServer{
		options:      options,
		addr:         listener.Addr(),
		listener:     listener,
		clients:      xsafeMap.NewStrAnyMap(),
		heartbeatRun: 0,
		closed:       make(chan byte, 1),
		requests:     xsafeMap.NewStrAnyMap(),
	}
}

func (p *TStreamServer) Log(level logger.LogLevel, args ...interface{}) {
	Log("TStreamServer", p.options.GetLogger(), level, args...)
}

func (p *TStreamServer) SetHandler(handler ServerStreamHandler) {
	p.handler = handler
}

func (p *TStreamServer) listen() error {
	if p.isListening() {
		return nil
	}
	if p.listener == nil {
		l, err := net.Listen(p.addr.Network(), p.addr.String())
		if err != nil {
			return err
		}
		p.listener = l
	}
	atomic.StoreInt32(&p.isListener, 1)
	return nil
}

func (p *TStreamServer) getInterrupted() bool {
	n := atomic.LoadInt32(&p.isInterrupted)
	return n != 0
}

func (p *TStreamServer) getListener() net.Listener {
	if n := atomic.LoadInt32(&p.isListener); n != 0 {
		return p.listener
	}
	return nil
}

func (p *TStreamServer) accept() (*tStreamConn, error) {
	interrupted := p.getInterrupted()
	if interrupted {
		return nil, errors.New("已中断服务")
	}
	listener := p.getListener()
	if listener == nil {
		return nil, errors.New("未设置服务监听地址")
	}
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	return newStreamConnByConn(p, conn, p.options.TConnOptions), nil
}

func (p *TStreamServer) isListening() bool {
	n := atomic.LoadInt32(&p.isListener)
	return n != 0 && p.listener != nil
}

func (p *TStreamServer) Open() error {
	if p.isListening() {
		return errors.New("套接字服务已经打开")
	}
	if l, err := net.Listen(p.addr.Network(), p.addr.String()); err != nil {
		return err
	} else {
		p.listener = l
	}
	return nil
}

func (p *TStreamServer) Addr() net.Addr {
	if p.listener != nil {
		return p.listener.Addr()
	}
	return p.addr
}

func (p *TStreamServer) Close() error {
	var err error
	select {
	case <-p.closed:
		return nil
	default:
	}
	close(p.closed)
	if p.isListening() {
		atomic.StoreInt32(&p.isListener, 0)
		err = p.listener.Close()
		p.listener = nil
	}
	return err
}

func (p *TStreamServer) interrupt() error {
	atomic.StoreInt32(&p.isInterrupted, 1)
	if e := p.Close(); e != nil {
		p.Log(logger.LEVEL_WARN, "中断服务出错:", e)
	}
	return nil
}

func (p *TStreamServer) innerAccept() error {
	client, err := p.accept()
	if err != nil {
		// 服务关闭退出
		return nil
	}
	select {
	case <-p.closed:
		// 已经关闭服务
		return nil
	default:
	}
	if client != nil {
		p.wg.Add(1)
		id := client.Id()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					p.Log(logger.LEVEL_WARN, "结束客户端[", id, "]连接出错:", e)
				}
			}()
			defer p.wg.Done()
			client.Start()
			p.heartbeatProcessLoop()
			client.doOnConnect()
			client.WaitStop()
		}()
	}
	return nil
}

func (p *TStreamServer) acceptLoop() error {
	for {
		select {
		case <-p.closed:
			// 已经关闭服务
			return nil
		default:
		}
		err := p.innerAccept()
		if err != nil {
			p.Log(logger.LEVEL_ERRO, "出错，服务已退出:", err)
			return err
		}
	}
}

func (p *TStreamServer) Serve() error {
	err := p.listen()
	if err != nil {
		return err
	}
	return p.acceptLoop()
}

func (p *TStreamServer) Stop() error {
	select {
	case <-p.closed:
		return nil
	default:
	}
	// 不再接受新的连接
	// 关闭已连接的客户端
	p.clients.Foreach(func(k string, v interface{}) bool {
		cli := v.(*tStreamConn)
		if e := cli.Close(); e != nil {
			p.Log(logger.LEVEL_WARN, "关闭客户端出错:", e)
		}
		return true
	})
	// 关闭服务
	_ = p.interrupt()
	p.wg.Wait()
	return nil
}
