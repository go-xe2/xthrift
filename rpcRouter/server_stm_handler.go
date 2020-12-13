/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 17:28
* Description:
*****************************************************************/

package rpcRouter

import (
	"bytes"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/netstream"
)

var _ netstream.ServerStreamHandler = (*TRouterServer)(nil)

func (p *TRouterServer) OnRecv(conn netstream.StreamConn, data []byte) {
	buf := bytes.NewBuffer(data)
	proto := NewRouterBinaryProto(buf)
	packetType, pktId, err := proto.ReadPacketBegin()
	if err != nil {
		xlog.Error(err)
		_ = p.sendErrorByConn(conn, -1, err.Error(), -1)
		return
	}
	defer proto.ReadPacketEnd()
	err = p.processRecv(conn, pktId, proto, data, packetType)
	if err != nil {
		_ = p.sendErrorByConn(conn, pktId, err.Error(), -1)
	}
}

func (p *TRouterServer) OnCall(conn netstream.StreamConn, data []byte) (result []byte, err error) {
	return nil, nil
}

func (p *TRouterServer) OnConnect(conn netstream.StreamConn) {
}

func (p *TRouterServer) OnReconnect(conn netstream.StreamConn) {

}

func (p *TRouterServer) OnDisconnect(conn netstream.StreamConn) {

}

func (p *TRouterServer) OnSendTo(conn netstream.StreamConn, target netstream.StreamConn, data []byte) {

}

func (p *TRouterServer) OnCallTo(conn netstream.StreamConn, target netstream.StreamConn, data []byte) (result []byte, err error) {
	return nil, nil
}

func (p *TRouterServer) OnRequest(reqConn netstream.StreamConn, reqId string, namespace string, data []byte) {
}
