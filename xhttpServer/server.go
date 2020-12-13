/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-17 16:09
* Description:
*****************************************************************/

package xhttpServer

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"strings"
	"sync"
)

type THttpServer struct {
	mu         *mux.Router
	svr        *http.Server
	listAddr   string
	lst        net.Listener
	wg         sync.WaitGroup
	baseRouter string
	handlers   map[string]*THttpRouterInfo
	subRouters map[string]*TSubRouter
}

func NewHttpServer(baseRouter, addr string) *THttpServer {
	inst := &THttpServer{
		baseRouter: baseRouter,
		listAddr:   addr,
		mu:         mux.NewRouter(),
		handlers:   make(map[string]*THttpRouterInfo),
		subRouters: make(map[string]*TSubRouter),
	}
	return inst
}

func NewHttpServerFromLst(baseRouter string, lst net.Listener) *THttpServer {
	inst := &THttpServer{
		baseRouter: baseRouter,
		listAddr:   lst.Addr().String(),
		lst:        lst,
		mu:         mux.NewRouter(),
		handlers:   make(map[string]*THttpRouterInfo),
		subRouters: make(map[string]*TSubRouter),
	}
	return inst
}

func (p *THttpServer) SetListener(lst net.Listener) {
	p.lst = lst
}

func (p *THttpServer) Handle(method, pattern string, handler HttpServerHandler, title string, summary string) {
	if _, ok := p.handlers[pattern]; !ok {
		h := newServerHandler(p.baseRouter, handler)
		h.SetTitle(title)
		h.SetSummary(summary)
		p.handlers[pattern] = NewHttpRouterInfo(h, method, pattern)
		r := p.mu.Handle(pattern, h)
		if method != "" {
			items := strings.Split(method, ",")
			r.Methods(items...)
		}
	}
}

func (p *THttpServer) HandleFun(method, pattern string, fun HttpServerHandleFun, title, summary string) {
	key := fmt.Sprintf("%s-%s", method, pattern)
	if _, ok := p.handlers[key]; !ok {
		h := newServerHandleByFun(p.baseRouter, fun)
		h.SetTitle(title)
		h.SetSummary(summary)
		p.handlers[key] = NewHttpRouterInfo(h, method, pattern)
		r := p.mu.Handle(pattern, h)
		if method != "" {
			items := strings.Split(method, ",")
			r.Methods(items...)
		}
	}
}

func (p *THttpServer) PathPrefixHandle(method string, path string, handler HttpServerHandler, title, summary string) {
	key := fmt.Sprintf("%s/*", path)
	if _, ok := p.handlers[key]; !ok {
		h := newServerHandler(p.baseRouter, handler)
		h.SetTitle(title)
		h.SetSummary(summary)
		p.handlers[key] = NewHttpRouterInfo(h, method, key)
		r := p.mu.PathPrefix(path).Handler(h)
		if method != "" {
			items := strings.Split(method, ",")
			r.Methods(items...)
		}
	}
}

func (p *THttpServer) PathPrefixHandleFun(method string, path string, fun HttpServerHandleFun, title, summary string) {
	key := fmt.Sprintf("%s/*", path)
	if _, ok := p.handlers[key]; !ok {
		h := newServerHandleByFun(p.baseRouter, fun)
		h.SetTitle(title)
		h.SetSummary(summary)
		p.handlers[key] = NewHttpRouterInfo(h, method, key)
		r := p.mu.PathPrefix(path).Handler(h)
		if method != "" {
			items := strings.Split(method, ",")
			r.Methods(items...)
		}
	}
}

func (p *THttpServer) SubRouter(path string) *TSubRouter {
	if v, ok := p.subRouters[path]; ok {
		return v
	}
	result := newSubRouter(p.mu.PathPrefix(path).Subrouter())
	p.subRouters[path] = result
	return result
}

func (p *THttpServer) Routers() map[string]*THttpRouterInfo {
	return p.handlers
}

func (p *THttpServer) Serve() error {
	if p.svr != nil {
		return nil
	}
	var err error
	if p.lst == nil {
		p.lst, err = net.Listen("tcp", p.listAddr)
		if err != nil {
			return err
		}
	}
	p.svr = &http.Server{Addr: p.listAddr, Handler: p.mu}
	if e := p.svr.Serve(p.lst); e != nil {
		if e == http.ErrServerClosed {
			return nil
		}
		return e
	}
	return nil
}

func (p *THttpServer) Stop() error {
	if p.svr == nil {
		return errors.New("服务未启动")
	}
	if err := p.svr.Shutdown(context.Background()); err != nil {
		return err
	}
	// 等待处理结束
	p.wg.Wait()
	p.svr = nil
	return nil
}

func (p *THttpServer) IsStart() bool {
	return p.svr != nil
}

func (p *THttpServer) ListenAddr() string {
	return p.listAddr
}
