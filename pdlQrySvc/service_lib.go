/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 17:05
* Description:
*****************************************************************/

package pdlQrySvc

import (
	"fmt"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/xhttpServer"
)

func (p *TService) returnSuccess(req *xhttpServer.THttpRequest, msg string, data interface{}) {
	req.ReturnJson(MakeResData(0, msg, data))
}

func (p *TService) returnError(req *xhttpServer.THttpRequest, status int, msg string) {
	req.ReturnJson(MakeResData(status, msg, nil))
}

func (p *TService) returnErrorf(req *xhttpServer.THttpRequest, status int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	p.returnError(req, status, s)
}

// handler /pdl/unRegProj
func (p *TService) UnRegProject(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	proj := t.String(params["proj"])
	if proj == "" {
		req.ReturnJson(MakeResData(1, "参数proj不能为空", nil))
	}

	if err := p.regCenter.HostStore().RemoveProject(proj); err != nil {
		req.ReturnJson(MakeResData(1, err.Error(), nil))
	}
	if err := p.regCenter.PdlStore().RemoveProject(proj); err != nil {
		req.ReturnJson(MakeResData(1, err.Error(), nil))
	}
	p.returnSuccess(req, "注册服务协议成功", nil)
}

// handler, /pdl/namespaces
func (p *TService) QryNamespaces(req *xhttpServer.THttpRequest) {
	items := p.regCenter.PDLQuery().AllNamespaces()
	p.returnSuccess(req, "", items)
}

// handler /pdl/namespace?ns=命名空间
func (p *TService) QryNamespace(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	ns := t.String(params["ns"])
	if ns == "" {
		req.ReturnJson(MakeResData(1, "参数ns不能为空", nil))
	}

	n := p.regCenter.PDLQuery().QryNamespace(ns)
	if n == nil {
		p.returnErrorf(req, 1, "命名空间%s不存在", ns)
	}
	p.returnSuccess(req, ns, n)
}

// handler /pdl/services
func (p *TService) QryServices(req *xhttpServer.THttpRequest) {
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
func (p *TService) QryMethods(req *xhttpServer.THttpRequest) {
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
func (p *TService) QryServiceMethod(req *xhttpServer.THttpRequest) {
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
func (p *TService) QryType(req *xhttpServer.THttpRequest) {
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
func (p *TService) QryTypedef(req *xhttpServer.THttpRequest) {
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

// handler /pdl/nshost
func (p *TService) ServiceHosts(req *xhttpServer.THttpRequest) {
	params := req.GetParams()
	svc := t.String(params["svc"])
	if svc == "" {
		p.returnError(req, 1, "请输入要查询的服务")
	}
	items := p.regCenter.HostStore().GetSvcHosts(svc)
	p.returnSuccess(req, svc, items)
}
