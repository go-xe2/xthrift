package rpcRouter

import (
	"bytes"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/netstream"
)

var _ netstream.ClientStreamHandler = (*TRouterClient)(nil)

func (p *TRouterClient) OnRecv(conn netstream.StreamConn, data []byte) {
	buf := bytes.NewBuffer(data)
	proto := NewRouterBinaryProto(buf)
	pktType, pktId, err := proto.ReadPacketBegin()
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("routerClient recv data.")
	defer proto.ReadPacketEnd()
	switch pktType {
	case REG_RES_PACKET:
		p.recvRegResult(pktId, proto)
		break
	case ERR_RES_PACKET:
		p.recvError(proto, pktId)
		break
	case CALL_PACKET:
		namespace, method, seqId, err := proto.ReadCallBegin()
		if err != nil {
			p.sendError(pktId, err.Error(), -1)
			xlog.Error(err)
			return
		}
		callData, err := proto.ReadData()
		if err != nil {
			p.sendError(pktId, err.Error(), -1)
			xlog.Error(err)
			return
		}
		if err := proto.ReadCallEnd(); err != nil {
			p.sendError(pktId, err.Error(), -1)
			xlog.Error(err)
			return
		}
		p.recvCall(pktId, namespace, method, seqId, callData)
		break
	case REPLY_PACKET:
		break
	}
}

func (p *TRouterClient) OnCall(conn netstream.StreamConn, data []byte) (result []byte, err error) {
	return nil, nil
}

func (p *TRouterClient) OnConnect(conn netstream.StreamConn) {

}

// 断线重连接成功
func (p *TRouterClient) OnReconnect(conn netstream.StreamConn) {
	// 断线重新连接后，需要向服务端重新注册协议，以便服务端绑定路由客户端与服务命名空间之间的关系
	xlog.Debug("rpcClient onReconnect.")
	allProjects := p.pdlStore.AllProject()
	for _, pn := range allProjects {
		if err := p.sendRegProject(pn.PDL, pn.MD5); err != nil {
			xlog.Error(err)
		}
	}
}

func (p *TRouterClient) OnDisconnect(conn netstream.StreamConn) {

}

func (p *TRouterClient) OnRequest(reqId string, namespace string, data []byte) {

}

func (p *TRouterClient) OnResponse(reqId string, data []byte) {

}
