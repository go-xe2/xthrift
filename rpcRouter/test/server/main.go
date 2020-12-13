package main

import (
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"github.com/go-xe2/xthrift/rpcRouter"
	"net"
)

func main() {
	opts := rpcRouter.NewSvcOptions()
	opts.PDLPath = "./pdl"
	opts.PDLExt = ".pdl"
	opts.HostPath = "hosts"
	opts.HostExt = ".host"
	opts.BaseRouter = "/v1"
	opts.BaseNamespace = "com.mnyun."

	httpLst, err := net.Listen("tcp", opts.HttpAddr)
	if err != nil {
		panic(err)
	}
	svcLst, err := net.Listen("tcp", opts.RouterAddr)
	if err != nil {
		panic(err)
	}
	svc := rpcRouter.NewServer(opts, httpLst, svcLst)
	xlog.Info("http服务:", opts.HttpAddr)
	xlog.Info("router服务:", opts.RouterAddr)

	processorFac := rpcRouter.NewRouterProcessorFactory(svc)

	rpcTrans, err := xthrift.NewServer(":3001")
	if err != nil {
		panic(err)
	}
	rpcTrans.SetProcessorFac(processorFac)
	go func() {
		xlog.Info("rpc服务端口:3001")
		rpcTrans.Serve()
	}()
	svc.Serve()
}
