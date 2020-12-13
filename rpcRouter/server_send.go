package rpcRouter

import (
	"github.com/go-xe2/xthrift/netstream"
)

func (p *TRouterServer) sendErrorByConn(conn netstream.StreamConn, pktId int64, msg string, code int32) error {
	data, err := makeErrorData(pktId, msg, code)
	if err != nil {
		return err
	}
	return p.sendByConn(conn, data)
}

func (p *TRouterServer) sendByConn(conn netstream.StreamConn, data []byte) error {
	if _, err := conn.Send(data); err != nil {
		return err
	}
	return nil
}

// 向客户端发送数据
func (p *TRouterServer) send(clientId string, data []byte) error {
	conn := p.getClientConn(clientId)
	if conn == nil {
		return nil
	}
	return p.sendByConn(conn, data)
}

func (p *TRouterServer) sendError(clientId string, pktId int64, msg string, code int32) error {
	data, err := makeErrorData(pktId, msg, code)
	if err != nil {
		return err
	}
	return p.send(clientId, data)
}
