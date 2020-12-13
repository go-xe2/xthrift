package rpcRouter

import (
	"context"
	"errors"
	"github.com/go-xe2/x/utils/xrand"
	"github.com/go-xe2/xthrift/regcenter"
	"time"
)

func (p *TRouterServer) RouterCall(ctx context.Context, namespace string, method string, seqId int32, rpcData []byte) ([]byte, error) {
	hosts := p.center.HostStore().GetSvcHosts(namespace)
	onLines := make([]*regcenter.THostStoreToken, 0)
	for _, h := range hosts {
		if h.Ext > 0 {
			onLines = append(onLines, h)
		}
	}
	if len(onLines) == 0 {
		return nil, errors.New("没有可用的服务资源")
	}
	// 随机获取一条链接
	idx := 0
	onLineCount := len(hosts)
	if onLineCount > 0 {
		idx = xrand.N(0, onLineCount-1)
	}
	host := hosts[idx]
	clientId := host.Host
	targetConn := p.getClientConn(clientId)
	if targetConn == nil {
		return nil, errors.New("服务已离线或没有可用的服务资源")
	}
	//p.callTimeout(clientId, fromClientId, pktId)
	pktId := time.Now().UnixNano()
	pktData := makeCallData(ctx, pktId, namespace, method, seqId, rpcData)

	rpcSender := newServerRPCSender()
	p.callTimeout(clientId, pktId, rpcSender)
	if _, err := targetConn.Send(pktData); err != nil {
		return nil, err
	}
	return rpcSender.Wait()
}
