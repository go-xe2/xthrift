/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-17 15:12
* Description:
*****************************************************************/

package rpcPoint

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/sync/xsafeMap"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"github.com/go-xe2/xthrift/regcenter"
	"github.com/go-xe2/xthrift/rpcRouter"
	"github.com/go-xe2/xthrift/xhttpServer"
	"github.com/modood/table"
	"net"
	"time"
)

type TEndPointServer struct {
	server     *xhttpServer.THttpServer
	baseRouter string

	isRun       bool
	options     *TOptions
	hostClients *xsafeMap.TStrAnyMap

	thriftSvc               *xthrift.TXServer
	transFac                thrift.TTransportFactory
	inProtoFac, outProtoFac thrift.TProtocolFactory
	regCenter               *regcenter.TRegCenter
	httpLst                 net.Listener
	thriftLst               net.Listener
	routerCli               *rpcRouter.TRouterClient
}

func NewEndPointServer(options *TOptions) *TEndPointServer {
	opts := options
	if opts == nil {
		opts = defaultOptions
	}

	inst := &TEndPointServer{
		options:     opts,
		server:      xhttpServer.NewHttpServer(options.BaseRouter, options.HttpAddr),
		hostClients: xsafeMap.NewStrAnyMap(),
		regCenter:   regcenter.NewRegCenter(opts.HostPath, opts.HostExt, opts.PDLPath, opts.PDLExt, opts.WatchHostChanged, opts.WatchPDLChanged),
	}
	//inst.regCenter.PdlStore().SetHandler(inst)
	if opts.EnableRouter {
		rc, err := rpcRouter.NewClient(opts.RouterId, inst, inst.regCenter.PdlStore(), opts.RouterSvr, opts.Router)
		if err != nil {
			xlog.Debug("创建router客户端失败:", err)
		} else {
			inst.routerCli = rc
		}
	}
	return inst.initDefaultHandle()
}

func (p *TEndPointServer) initDefaultHandle() *TEndPointServer {
	p.setHandler("get", "/pdl/help", p.help, "服务接口帮助", "参数：无")
	p.setHandler("get", "/pdl/namespaces", p.QryNamespaces, "获取所有命名空间列表", "参数: 无")
	p.setHandler("get", "/pdl/namespace", p.QryNamespace, "获取命名空间定义", "参数: ns: 空间名")
	p.setHandler("get", "/pdl/svcHost", p.ServiceHosts, "获取提供服务的服务器地址", "参数: svc: 服务名称，包含服务所在的空间名")
	p.setHandler("get", "/pdl/services", p.QryServices, "获取命名空间下的所有服务列表", "参数: ns: 空间名")
	p.setHandler("get", "/pdl/service", p.QryMethods, "获取服务下的所有接口", "参数: ns: 空间名, svc: 服务名称")
	p.setHandler("get", "/pdl/method", p.QryServiceMethod, "获取接口定义", "参数: ns: 空间名, svc: 服务名， mn: 要查询的接口名称")
	p.setHandler("get", "/pdl/type", p.QryType, "查询空间中的数据类型定义", "参数: ns:空间名, dt:数据类型名，不包含命名空间")
	p.setHandler("get", "/pdl/typedef", p.QryTypedef, "查询空间中的数据类型别名", "参数: ns:空间名, dt:数据类型别名，不包含命名空间")
	p.setHandler("post", "/pdl/reg", p.RegisterService, "注册命名空间服务", "参数: host:服务地址，port:服务端口，data:传入的协议文件数据内容，格式为base64编码字符串")
	p.setHandler("post", "/pdl/unRegHost", p.UnRegisterHost, "注销提供服务的服务器", "参数: host:服务地址，port:服务端口")
	p.setHandler("post", "/pdl/unRegProj", p.UnRegProject, "注销服务协议", "参数: proj:服务协议名称")

	return p
}

func (p *TEndPointServer) SetHttpListener(lst net.Listener) {
	p.httpLst = lst
	p.server.SetListener(lst)
}

func (p *TEndPointServer) SetThriftListener(lst net.Listener) {
	p.thriftLst = lst
}

func (p *TEndPointServer) setHandler(method string, pattern string, fn xhttpServer.HttpServerHandleFun, title, summary string) {
	p.server.HandleFun(method, pattern, fn, title, summary)
}

func (p *TEndPointServer) Handle(method, pattern string, handler xhttpServer.HttpServerHandler, title, summary string) {
	p.server.Handle(method, pattern, handler, title, summary)
}

func (p *TEndPointServer) HandleFun(method, pattern string, fun xhttpServer.HttpServerHandleFun, title, summary string) {
	p.server.HandleFun(method, pattern, fun, title, summary)
}

func (p *TEndPointServer) PathPrefixHandle(method string, pattern string, handler xhttpServer.HttpServerHandler, title, summary string) {
	p.server.PathPrefixHandle(method, pattern, handler, title, summary)
}

func (p *TEndPointServer) PathPrefixHandleFun(method string, pattern string, fun xhttpServer.HttpServerHandleFun, title, summary string) {
	p.server.PathPrefixHandleFun(method, pattern, fun, title, summary)
}

func (p *TEndPointServer) Serve() error {
	if p.isRun {
		return nil
	}
	var err error
	if p.httpLst == nil {
		p.httpLst, err = net.Listen("tcp", p.options.HttpAddr)
		if err != nil {
			return err
		}
		p.server.SetListener(p.httpLst)
	}

	if p.thriftLst == nil {
		p.thriftLst, err = net.Listen("tcp", p.options.ThriftAddr)
		if err != nil {
			return err
		}
	}

	if p.routerCli != nil {
		// 连接路由
		if err := p.routerCli.Open(); err != nil {
			xlog.Debug("连接路由出错:", err)
			p.routerCli = nil
		}
	}

	if err := p.regCenter.Load(); err != nil {
		return err
	}

	if p.options.WatchPDLChanged {
		if err := p.regCenter.WatchPDLChanged(); err != nil {
			return err
		}
	}
	if p.options.WatchHostChanged {
		if err := p.regCenter.WatchHostChanged(); err != nil {
			return err
		}
	}

	if err := p.InitInnerServer(); err != nil {
		return err
	}

	var errChan = make(chan error, 1)
	go func() {
		if e := p.server.Serve(); e != nil {
			select {
			case <-errChan:
				// 已经关闭
				xlog.Error("服务出错:", e)
				return
			default:
			}
			errChan <- e
			xlog.Error("服务出错:", e)
		}
	}()
	go func() {
		if e := p.InnerServerStart(); e != nil {
			select {
			case <-errChan:
				xlog.Error("内部服务启动出错:", e)
				return
			default:
			}
			errChan <- e
			xlog.Error("内部服务启动出错:", e)
		}
	}()
	time.AfterFunc(1*time.Second, func() {
		close(errChan)
	})
	select {
	case e, ok := <-errChan:
		if ok {
			_ = p.Stop()
			return e
		}
		break
	}
	xlog.Info("http端口号:", p.server.ListenAddr())
	xlog.Info("thrift服务端口号:", p.options.ThriftAddr)
	p.PrintAllHandle()

	return nil
}

func (p *TEndPointServer) Stop() error {
	if p.thriftSvc != nil {
		if err := p.InnerServerStop(); err != nil {
			return err
		}
	}
	if p.server != nil {
		return p.server.Stop()
	}
	return nil
}

func (p *TEndPointServer) returnError(req *xhttpServer.THttpRequest, status int, msg string) {
	req.ReturnJson(NewResResult(status, msg, nil))
}

func (p *TEndPointServer) returnErrorf(req *xhttpServer.THttpRequest, status int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	p.returnError(req, status, s)
}

func (p *TEndPointServer) returnSuccess(req *xhttpServer.THttpRequest, msg string, data interface{}) {
	req.ReturnJson(NewResResult(0, msg, data))
}

func (p *TEndPointServer) PrintAllHandle() {
	//xlog.Info("所有服务路由:")
	//xlog.Info(strings.Repeat("-", 30))
	//xlog.Info("\t请求方法\t\t\t|\trouter\t\t\t|\t说明\t\t\t")
	type infoLine struct {
		Method  string
		Router  string
		Summary string
	}
	lines := make([]infoLine, 1)
	lines[0] = infoLine{"请求方法", "路由", "说明"}
	routers := p.server.Routers()
	for _, node := range routers {
		lines = append(lines, infoLine{node.GetMethod(), node.GetPattern(), node.GetTitle() + " " + node.GetSummary()})
		//xlog.Info(fmt.Sprintf("\t%s\t\t\t|\t%s\t\t\t|\t%s", items[0], items[1], node.summary))
	}
	s := table.Table(lines)
	fmt.Print(s)
	fmt.Print("\n")
	//xlog.Info(strings.Repeat("-", 30))
}
