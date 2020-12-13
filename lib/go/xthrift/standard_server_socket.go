/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 15:04
* Description:
*****************************************************************/

package xthrift

import (
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/core/exception"
	"net"
	"sync"
	"time"
)

type TStandardServerSocket struct {
	listener      net.Listener
	addr          net.Addr
	clientTimeout time.Duration

	// Protects the interrupted value to make it thread safe.
	mu          sync.RWMutex
	interrupted bool
}

func NewStandardServerSocket(addr string) (svc *TStandardServerSocket, err error) {
	return NewStandardSocketTimeout(addr, 0)
}

func NewStandardSocketTimeout(listenAddr string, clientTimeout time.Duration) (*TStandardServerSocket, error) {
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &TStandardServerSocket{addr: addr, clientTimeout: clientTimeout}, nil
}

func NewStandardSocketFromListener(listener net.Listener) (svc *TStandardServerSocket, err error) {
	if listener == nil {
		return nil, errors.New("listener is nil")
	}
	return &TStandardServerSocket{
		listener:      listener,
		addr:          listener.Addr(),
		clientTimeout: 0,
	}, nil
}

// Creates a TServerSocket from a net.Addr
func NewStandardSocketFromAddrTimeout(addr net.Addr, clientTimeout time.Duration) *TStandardServerSocket {
	return &TStandardServerSocket{addr: addr, clientTimeout: clientTimeout}
}

func (p *TStandardServerSocket) Listen() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.IsListening() {
		return nil
	}
	l, err := net.Listen(p.addr.Network(), p.addr.String())
	if err != nil {
		return err
	}
	p.listener = l
	return nil
}

func (p *TStandardServerSocket) Accept() (thrift.TTransport, error) {
	p.mu.RLock()
	interrupted := p.interrupted
	p.mu.RUnlock()

	if interrupted {
		return nil, exception.NewText("Transport Interrupted")
	}
	p.mu.Lock()
	listener := p.listener
	p.mu.Unlock()
	if listener == nil {
		return nil, thrift.NewTTransportException(thrift.NOT_OPEN, "No underlying server socket")
	}

	conn, err := listener.Accept()
	if err != nil {
		return nil, thrift.NewTTransportExceptionFromError(err)
	}
	return thrift.NewTSocketFromConnTimeout(conn, p.clientTimeout), nil
}

// Checks whether the socket is listening.
func (p *TStandardServerSocket) IsListening() bool {
	return p.listener != nil
}

// Connects the socket, creating a new socket object if necessary.
func (p *TStandardServerSocket) Open() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.IsListening() {
		return thrift.NewTTransportException(thrift.ALREADY_OPEN, "Server socket already open")
	}
	if l, err := net.Listen(p.addr.Network(), p.addr.String()); err != nil {
		return err
	} else {
		p.listener = l
	}
	return nil
}

func (p *TStandardServerSocket) Addr() net.Addr {
	if p.listener != nil {
		return p.listener.Addr()
	}
	return p.addr
}

func (p *TStandardServerSocket) Close() error {
	var err error
	p.mu.Lock()
	if p.IsListening() {
		err = p.listener.Close()
		p.listener = nil
	}
	p.mu.Unlock()
	return err
}

func (p *TStandardServerSocket) Interrupt() error {
	p.mu.Lock()
	p.interrupted = true
	p.mu.Unlock()
	p.Close()

	return nil
}
