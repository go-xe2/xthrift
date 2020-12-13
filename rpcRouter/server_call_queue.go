package rpcRouter

import (
	"errors"
	"fmt"
	"github.com/go-xe2/x/os/xlog"

	"time"
)

type tCallInfo struct {
	id     string
	svc    *TRouterServer
	pktId  int64
	sender ServerSender
	timer  *time.Timer
}

func newCallInfo(id string, svc *TRouterServer, pktId int64, sender ServerSender) *tCallInfo {
	return &tCallInfo{
		id:     id,
		svc:    svc,
		sender: sender,
		pktId:  pktId,
	}
}

func (p *tCallInfo) start(timeout time.Duration) *tCallInfo {
	xlog.Debug("tCallInfo start timeout:", timeout)
	p.timer = time.AfterFunc(timeout, func() {
		p.svc.removeCall(p.id)
		p.sender.SendErr(p.pktId, errors.New("访问超时"), -1)
	})
	return p
}

func (p *TRouterServer) removeCall(callId string) {
	p.callTimeouts.Remove(callId)
}

// 设置调用超时
func (p *TRouterServer) callTimeout(targetId string, pktId int64, sender ServerSender) {
	id := fmt.Sprintf("%s:%d", targetId, pktId)
	info := newCallInfo(id, p, pktId, sender)
	p.callTimeouts.Set(id, info)
	xlog.Debug("callTimeout call id:", id, ", timeout:", p.options.ReadTimeout, ", pktId:", pktId)
	info.start(p.options.ReadTimeout)
}

// 检查调用是否已经超时
func (p *TRouterServer) callReply(fromId string, pktId int64, pktData []byte) error {
	id := fmt.Sprintf("%s:%d", fromId, pktId)
	v := p.callTimeouts.Get(id)
	xlog.Debug("rpcRouter callReply id:", id, ", v:", v)
	if v == nil {
		xlog.Debug("rpcRouter callReply id:", id, "，超时, pktId:", pktId)
		// 超时
		return nil
	}
	p.removeCall(id)
	info := v.(*tCallInfo)
	info.timer.Stop()
	info.sender.SendPacket(pktData)
	return nil
}
