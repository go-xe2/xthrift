/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 10:33
* Description:
*****************************************************************/

package netstream

import (
	"errors"
	"github.com/go-xe2/x/core/logger"
	"net"
)

// 网络流客户端
// 功能特点
// 1.与网络流服务端使用心跳方式保持长连接
// 2.具有心跳检查及断线重连功能
// 3.具有双向通讯功能
// 4.具有点对点发送功能
// 5.具有同步发送数据和异步发送数据功能

type StreamClient interface {
	IsOpen() bool
	Addr() net.Addr
	RemoteAddr() net.Addr
	Send(data []byte) (int, error)
	Call(data []byte) (result []byte, err error)
	SendTo(targetId string, data []byte) (int, error)
	CallTo(targetId string, data []byte) (result []byte, err error)
	Request(reqId string, namespace string, body []byte) error
	Response(reqId string, body []byte) error
	MakeRequestId() string
	Id() string
	Open() error
	Close() error
	OnReady(conn StreamConn)
	SetHandler(handler ClientStreamHandler)
}

type StreamPoint interface {
	ClientStreamHandler
	Start(client StreamClient)
	Stop()
}

type TStreamClient struct {
	// 是否已经关闭
	heartbeatIsRun int32
	conn           *tStreamConn
	handler        ClientStreamHandler
	// 当前重试连接的次数
	tryCount    int
	closed      chan byte
	options     *TStmClientOptions
	retryStatus int32
}

var _ StreamClient = (*TStreamClient)(nil)

var _ StreamConnHandler = (*TStreamClient)(nil)

func NewStreamClient(hostPort string, options *TStmClientOptions) (*TStreamClient, error) {
	if options == nil {
		options = DefaultStmClientOptions
	}
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		return nil, err
	}
	inst := &TStreamClient{
		// 当前重试连接的次数
		options:  options,
		tryCount: 0,
	}
	inst.conn = newStreamConn(inst, addr, options.TConnOptions)
	return inst, nil
}

func (p *TStreamClient) SetHandler(handler ClientStreamHandler) {
	p.handler = handler
}

func (p *TStreamClient) Id() string {
	return p.conn.Id()
}

// 发送数据
func (p *TStreamClient) Send(data []byte) (int, error) {
	if p.IsOpen() {
		return p.conn.Send(data)
	}
	return 0, errors.New("未打开连接或已关闭连接")
}

func (p *TStreamClient) Call(data []byte) ([]byte, error) {
	if !p.IsOpen() {
		return nil, errors.New("连接未打开或已断开")
	}
	return p.conn.Call(data)
}

func (p *TStreamClient) SendTo(targetId string, data []byte) (int, error) {
	if p.IsOpen() {
		return p.conn.SendTo(targetId, data)
	}
	return 0, errors.New("未打开连接或已关闭连接")
}

func (p *TStreamClient) CallTo(targetId string, data []byte) ([]byte, error) {
	if !p.IsOpen() {
		return nil, errors.New("未打开连接或已关闭连接")
	}
	return p.conn.CallTo(targetId, data)
}

func (p *TStreamClient) Open() error {
	if e := p.conn.Open(); e != nil {
		return e
	}
	p.closed = make(chan byte, 1)
	return nil
}

func (p *TStreamClient) Serve(point StreamPoint) error {
	if point != nil {
		p.handler = point
	}
	if !p.IsOpen() {
		return errors.New("连接未打开，请先打开连接")
	}
	if point != nil {
		point.Start(p)
	}
	select {
	case <-p.closed:
	}
	if point != nil {
		point.Stop()
	}
	return nil
}

// Retrieve the underlying net.Conn
func (p *TStreamClient) Conn() StreamConn {
	return p.conn
}

// Returns true if the connection is open
func (p *TStreamClient) IsOpen() bool {
	return p.conn.IsOpen()
}

func (p *TStreamClient) innerClose() {
	if p.conn.IsOpen() {
		p.conn.innerClose()
	}
}

// Closes the socket.
func (p *TStreamClient) Close() error {
	select {
	case <-p.closed:
		return nil
	default:
	}
	if !p.IsOpen() {
		return nil
	}
	close(p.closed)
	return p.conn.Close()
}

func (p *TStreamClient) Addr() net.Addr {
	return p.conn.addr
}

func (p *TStreamClient) RemoteAddr() net.Addr {
	if p.conn == nil {
		return nil
	}
	return p.conn.RemoteAddr()
}

func (p *TStreamClient) Log(level logger.LogLevel, args ...interface{}) {
	Log("TStreamClient", p.options.GetLogger(), level, args...)
}

func (p *TStreamClient) Request(reqId string, namespace string, body []byte) error {
	if p.IsOpen() {
		return p.conn.Request(reqId, namespace, body)
	}
	return errors.New("客户端未连接")
}

func (p *TStreamClient) Response(reqId string, body []byte) error {
	if p.IsOpen() {
		return p.conn.Response(reqId, body)
	}
	return errors.New("客户端未连接")
}

func (p *TStreamClient) MakeRequestId() string {
	return MakeRequestId()
}
