/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 17:14
* Description:
*****************************************************************/

package pdlQrySvc

import (
	"github.com/go-xe2/xthrift/xhttpServer"
)

func (p *TService) RegisterHandlers(canRegPdl bool) {
	p.setHandler("get", "/pdl/help", p.help, "服务接口帮助", "参数：无")
	p.setHandler("get", "/pdl/namespaces", p.QryNamespaces, "获取所有命名空间列表", "参数: 无")
	p.setHandler("get", "/pdl/namespace", p.QryNamespace, "获取命名空间定义", "参数: ns: 空间名")
	p.setHandler("get", "/pdl/svcHost", p.ServiceHosts, "获取提供服务的服务器地址", "参数: svc: 服务名称，包含服务所在的空间名")
	p.setHandler("get", "/pdl/services", p.QryServices, "获取命名空间下的所有服务列表", "参数: ns: 空间名")
	p.setHandler("get", "/pdl/service", p.QryMethods, "获取服务下的所有接口", "参数: ns: 空间名, svc: 服务名称")
	p.setHandler("get", "/pdl/method", p.QryServiceMethod, "获取接口定义", "参数: ns: 空间名, svc: 服务名， mn: 要查询的接口名称")
	p.setHandler("get", "/pdl/type", p.QryType, "查询空间中的数据类型定义", "参数: ns:空间名, dt:数据类型名，不包含命名空间")
	p.setHandler("get", "/pdl/typedef", p.QryTypedef, "查询空间中的数据类型别名", "参数: ns:空间名, dt:数据类型别名，不包含命名空间")
	if canRegPdl {
		p.setHandler("post", "/pdl/reg", p.RegisterService, "注册命名空间服务", "参数: host:服务地址，port:服务端口，data:传入的协议文件数据内容，格式为base64编码字符串")
		p.setHandler("post", "/pdl/unRegHost", p.UnRegisterHost, "注销提供服务的服务器", "参数: host:服务地址，port:服务端口")
		p.setHandler("post", "/pdl/unRegProj", p.UnRegProject, "注销服务协议", "参数: proj:服务协议名称")
	}
}

func (p *TService) setHandler(method string, path string, handlerFun xhttpServer.HttpServerHandleFun, title, comment string) {
	p.server.HandleFun(method, path, handlerFun, title, comment)
}
