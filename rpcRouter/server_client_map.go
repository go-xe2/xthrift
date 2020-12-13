package rpcRouter

import "github.com/go-xe2/xthrift/netstream"

func (p *TRouterServer) saveClient(clientId string, connId string) {
	p.clientConnIds[clientId] = connId
	p.connClientIds[connId] = clientId
}

func (p *TRouterServer) getConnIdByClientId(clientId string) string {
	if s, ok := p.clientConnIds[clientId]; ok {
		return s
	}
	return ""
}

func (p *TRouterServer) getClientIdByConnId(connId string) string {
	if s, ok := p.connClientIds[connId]; ok {
		return s
	}
	return ""
}

func (p *TRouterServer) RemoveClientId(clientId string) {
	connId := p.getConnIdByClientId(clientId)
	if connId != "" {
		delete(p.connClientIds, connId)
	}
	delete(p.clientConnIds, clientId)
}

func (p *TRouterServer) getClientConn(clientId string) netstream.StreamConn {
	connId := p.getConnIdByClientId(clientId)
	if connId == "" {
		return nil
	}
	return p.stmSvr.GetClient(connId)
}
