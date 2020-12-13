/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 18:35
* Description:
*****************************************************************/

package xhttpServer

import (
	"github.com/gorilla/mux"
	"strings"
)

type TSubRouter struct {
	r          *mux.Router
	baseRouter string
}

func newSubRouter(r *mux.Router) *TSubRouter {
	return &TSubRouter{
		r: r,
	}
}

func (p *TSubRouter) Handle(method string, pattern string, handler HttpServerHandler) {
	h := newServerHandler(p.baseRouter, handler)
	r := p.r.Handle(pattern, h)
	if method != "" {
		items := strings.Split(method, ",")
		r.Methods(items...)
	}
}

func (p *TSubRouter) HandleFun(method string, pattern string, fun HttpServerHandleFun) {
	h := newServerHandleByFun(p.baseRouter, fun)
	r := p.r.Handle(pattern, h)
	if method != "" {
		items := strings.Split(method, ",")
		r.Methods(items...)
	}
}
