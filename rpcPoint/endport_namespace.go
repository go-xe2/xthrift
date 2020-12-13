/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-18 14:57
* Description:
*****************************************************************/

package rpcPoint

import (
	"fmt"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/xhttpServer"
)

// handler /pdl/unRegProj
func (p *TEndPointServer) UnRegProject(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	proj := t.String(params["proj"])
	if proj == "" {
		req.ReturnJson(NewResResult(1, "参数proj不能为空", nil))
	}

	if err := p.regCenter.HostStore().RemoveProject(proj); err != nil {
		req.ReturnJson(NewResResult(1, err.Error(), nil))
	}
	if err := p.regCenter.PdlStore().RemoveProject(proj); err != nil {
		req.ReturnJson(NewResResult(1, err.Error(), nil))
	}
	p.returnSuccess(req, "注册服务协议成功", nil)
}

// handler, /pdl/namespaces
func (p *TEndPointServer) QryNamespaces(req *xhttpServer.THttpRequest) {
	items := p.regCenter.PDLQuery().AllNamespaces()
	p.returnSuccess(req, "", items)
}

// handler /pdl/namespace?ns=命名空间
func (p *TEndPointServer) QryNamespace(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	ns := t.String(params["ns"])
	if ns == "" {
		req.ReturnJson(NewResResult(1, "参数ns不能为空", nil))
	}

	n := p.regCenter.PDLQuery().QryNamespace(ns)
	if n == nil {
		p.returnErrorf(req, 1, "命名空间%s不存在", ns)
	}
	p.returnSuccess(req, ns, n)
}

// handler /pdl/services
func (p *TEndPointServer) QryServices(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	ns := t.String(params["ns"])
	if ns == "" {
		p.returnError(req, 1, "参数ns不能为空")
	}
	_, items := p.regCenter.PDLQuery().QryServices(ns)
	if items == nil {
		p.returnErrorf(req, 1, "命名空间%s没有定义服务", ns)
	}
	p.returnSuccess(req, ns, items)
}

// handler /pdl/service
func (p *TEndPointServer) QryMethods(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	nsName := t.String(params["ns"])
	svcName := t.String(params["svc"])
	if nsName == "" {
		p.returnError(req, 1, "参数ns不能为空")
	}
	if svcName == "" {
		p.returnError(req, 1, "参数svc不参为空")
	}
	n, svc := p.regCenter.PDLQuery().QryServiceByNS(nsName, svcName)
	if svc == nil {
		if n != nil {
			p.returnErrorf(req, 1, "命名空间%s不存在服务%s", nsName, svcName)
		}
		p.returnErrorf(req, 1, "不存在命名空间%s", nsName)
	}
	p.returnSuccess(req, fmt.Sprintf("%s.%s", nsName, svcName), svc)
}

// handler /pdl/method
func (p *TEndPointServer) QryServiceMethod(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	nsName := t.String(params["ns"])
	svcName := t.String(params["svc"])
	mName := t.String(params["mn"])

	if nsName == "" {
		p.returnError(req, 1, "参数ns不能为空")
	}
	if svcName == "" {
		p.returnError(req, 1, "参数svc不参为空")
	}
	if mName == "" {
		p.returnError(req, 1, "参数mn不能为空")
	}

	ns, svc, m := p.regCenter.PDLQuery().QryMethodByNS(nsName, svcName, mName)
	if ns == nil {
		p.returnErrorf(req, 1, "命名空间%s不存在", nsName)
	}
	if svc == nil {
		p.returnErrorf(req, 1, "命名空间%s不存在服务%s", nsName, svcName)
	}
	if m == nil {
		p.returnErrorf(req, 1, "服务%s.%s不存在接口%s", nsName, svcName, mName)
	}
	p.returnSuccess(req, fmt.Sprintf("%s.%s.%s", nsName, svcName, mName), m)
}

// handler /pdl/type
func (p *TEndPointServer) QryType(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	ns := t.String(params["ns"])
	dt := t.String(params["dt"])

	if ns == "" {
		p.returnError(req, 1, "参数ns不能为空")
	}
	if dt == "" {
		p.returnError(req, 1, "参数dt不能为空")
	}

	_, dType := p.regCenter.PDLQuery().QryTypeByNS(ns, dt)
	if dType == nil {
		p.returnError(req, 1, fmt.Sprintf("服务%s中不存在类型%s", ns, dt))
	}
	p.returnSuccess(req, fmt.Sprintf("%s.%s", ns, dt), dType)
}

// handler /pdl/typedef
func (p *TEndPointServer) QryTypedef(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	ns := t.String(params["ns"])
	dt := t.String(params["dt"])

	if ns == "" {
		p.returnError(req, 1, "参数ns不能为空")
	}
	if dt == "" {
		p.returnError(req, 1, "参数dt不能为空")
	}

	_, dType := p.regCenter.PDLQuery().QryTypeDefByNS(ns, dt)
	if dType == nil {
		p.returnError(req, 1, fmt.Sprintf("服务%s中不存在类型%s", ns, dt))
	}
	p.returnSuccess(req, fmt.Sprintf("%s.%s", ns, dt), dType)
}
