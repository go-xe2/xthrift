/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 16:35
* Description:
*****************************************************************/

package rpcRouter

import (
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/sync/xsafeMap"
	"github.com/go-xe2/xthrift/gateway"
	"github.com/go-xe2/xthrift/netstream"
	"github.com/go-xe2/xthrift/pdlQrySvc"
	"github.com/go-xe2/xthrift/regcenter"
	"github.com/go-xe2/xthrift/xhttpServer"
	"net"
)

type TRouterServer struct {
	options    *TSvrOptions
	stmSvr     *netstream.TStreamServer
	center     *regcenter.TRegCenter
	stmSvrOpts *netstream.TStmServerOptions
	httpSvr    *xhttpServer.THttpServer
	httpQry    *pdlQrySvc.TService
	// 客户id与连接id的映射关系
	clientConnIds map[string]string
	// 连接id与客户端id的映射关系
	connClientIds map[string]string
	callTimeouts  *xsafeMap.TStrAnyMap
	pdlProjects   map[string]*tProjectInfo
	gateway       *gateway.TSvrHttpHandler
}

func NewServer(options *TSvrOptions, httpLst, svcLst net.Listener) *TRouterServer {
	opts := defaultSvrOptions
	if options != nil {
		opts = options
	}
	stmSvrOpts := netstream.NewStmServerOptions()
	stmSvrOpts.SetAllowMaxLoss(opts.AllowMaxLoss)
	stmSvrOpts.SetConnectTimeout(opts.ConnectTimeout)
	stmSvrOpts.SetHeartbeatSpeed(opts.Heartbeat)
	stmSvrOpts.SetReadTimeout(opts.ReadTimeout)
	stmSvrOpts.SetWriteTimeout(opts.ReadTimeout)
	stmSvrOpts.SetRecvBufSize(opts.RecvBufSize)
	stmSvrOpts.SetSendBufSize(opts.SendBufSize)
	inst := &TRouterServer{
		center:        regcenter.NewRegCenter(opts.HostPath, opts.HostExt, opts.PDLPath, opts.PDLExt, opts.WatchHostChanged, opts.WatchPDLChanged),
		stmSvrOpts:    stmSvrOpts,
		options:       opts,
		clientConnIds: make(map[string]string),
		connClientIds: make(map[string]string),
		pdlProjects:   make(map[string]*tProjectInfo),
		callTimeouts:  xsafeMap.NewStrAnyMap(),
	}
	inst.stmSvr = netstream.NewStreamServerByListener(svcLst, stmSvrOpts)
	inst.httpSvr = xhttpServer.NewHttpServerFromLst(opts.BaseRouter, httpLst)
	inst.httpQry = pdlQrySvc.NewService(inst.httpSvr, inst.center)
	inst.httpQry.RegisterHandlers(false)
	inst.stmSvr.SetHandler(inst)
	inst.center.PdlStore().SetHandler(inst)
	inst.gateway = gateway.NewSvrHttpHandler(inst.center.PDLQuery(), inst, opts.BaseRouter, opts.BaseNamespace)
	inst.httpSvr.PathPrefixHandle("", opts.BaseRouter, inst.gateway, "调用服务", "调用内部服务")
	return inst
}

func (p *TRouterServer) Serve() error {
	if err := p.center.Load(); err != nil {
		return err
	}
	go func() {
		if err := p.httpSvr.Serve(); err != nil {
			xlog.Error(err)
		}
	}()
	return p.stmSvr.Serve()
}

func (p *TRouterServer) Stop() error {
	if p.httpSvr.IsStart() {
		if err := p.httpSvr.Stop(); err != nil {
			xlog.Error(err)
		}
	}
	return p.stmSvr.Stop()
}
