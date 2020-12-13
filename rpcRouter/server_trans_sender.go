package rpcRouter

import "github.com/go-xe2/x/os/xlog"

type tServerTransSender struct {
	svc          *TRouterServer
	fromClientId string
}

var _ ServerSender = (*tServerTransSender)(nil)

func newServerTransSender(svc *TRouterServer, fromClientId string) ServerSender {
	return &tServerTransSender{
		svc:          svc,
		fromClientId: fromClientId,
	}
}

func (p *tServerTransSender) SendPacket(pktData []byte) {
	p.svc.send(p.fromClientId, pktData)
}

func (p *tServerTransSender) SendErr(pktId int64, err error, code int32) {
	if err := p.svc.sendError(p.fromClientId, pktId, err.Error(), code); err != nil {
		xlog.Error(err)
	}
}
