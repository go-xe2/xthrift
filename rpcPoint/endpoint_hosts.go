/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-18 16:49
* Description:
*****************************************************************/

package rpcPoint

import (
	"errors"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/xhttpServer"
)

// handler /pdl/nshost
func (p *TEndPointServer) ServiceHosts(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	svc := t.String(params["svc"])
	if svc == "" {
		p.returnError(req, 1, "请输入要查询的服务")
	}
	items := p.regCenter.HostStore().GetSvcHosts(svc)
	p.returnSuccess(req, svc, items)
}

func (p *TEndPointServer) GetServiceHostClient(serviceName string) (*tInnerClient, error) {
	if pool := p.hostClients.Get(serviceName); pool != nil {
		c := pool.(*tHostClientPool).Get()
		if c == nil {
			p.hostClients.Remove(serviceName)
			return nil, errors.New("没有可用的服务资源")
		}
		return c, nil
	}
	hosts := p.regCenter.HostStore().GetSvcHosts(serviceName)
	size := len(hosts)

	if size == 0 {
		return nil, errors.New("没有可用的服务资源")
	}
	pool := newHostClientPool(
		p,
		hosts,
		p.options.ReadTimeout,
		p.options.WriteTimeout,
		p.options.ConnectTimeout,
		p.options.ConnectFailRetry,
		p.options.ClientPoolKeepAlive,
		p.options.Heartbeat,
		p.options.HeartbeatLoss)

	p.hostClients.Set(serviceName, pool)
	c := pool.Get()

	if c == nil {
		return nil, errors.New("没有可用的服务资源")
	}
	return c, nil
}
