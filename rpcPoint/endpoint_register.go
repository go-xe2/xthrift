/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-19 11:21
* Description:
*****************************************************************/

package rpcPoint

import (
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/xhttpServer"
)

// handler /pdl/reg
func (p *TEndPointServer) RegisterService(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	host := t.String(params["host"])
	port := t.Int(params["port"])
	pdlData := t.String(params["pdl"])
	if host == "" || port <= 0 {
		p.returnError(req, 1, "服务地址或端口不能为空")
	}
	if pdlData == "" {
		p.returnError(req, 1, "未上协议数据")
	}
	if err := p.saveHostNamespace(host, port, pdlData); err != nil {
		p.returnErrorf(req, 1, err.Error())
	}
	p.returnSuccess(req, "注册成功", nil)
}

// handler /pdl/unRegHost
func (p *TEndPointServer) UnRegisterHost(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	host := t.String(params["host"])
	port := t.Int(params["port"])
	if host == "" || port <= 0 {
		p.returnError(req, 1, "服务地址或端口不能为空")
	}
	if err := p.regCenter.HostStore().RemoveHost(host, port); err != nil {
		p.returnError(req, 1, err.Error())
	}
	if err := p.regCenter.HostStore().Save(); err != nil {
		p.returnErrorf(req, 1, err.Error())
	}
	p.returnSuccess(req, "注册成功", nil)
}

func (p *TEndPointServer) saveHostNamespace(host string, port int, data string) error {
	proj, err := p.regCenter.PdlStore().AddProjectFromBase64([]byte(data))
	if err != nil {
		return err
	}
	p.regCenter.HostStore().AddHostWithProject(proj, host, port)
	if err := p.regCenter.HostStore().Save(); err != nil {
		return err
	}
	return nil
}
