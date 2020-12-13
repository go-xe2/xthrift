/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 17:34
* Description:
*****************************************************************/

package xhttpServer

import (
	"github.com/go-xe2/x/os/xlog"
	"net/http"
)

type HttpServerHandleFun func(req *THttpRequest)

type HttpServerHandler interface {
	Handle(req *THttpRequest)
}

type tServerHandlerStore struct {
	baseRouter string
	handler    HttpServerHandler
	fun        HttpServerHandleFun
	// 服务说明
	summary string
	// 服务名称
	title string
}

type THttpRouterInfo struct {
	handler *tServerHandlerStore
	method  string
	pattern string
}

func NewHttpRouterInfo(handler *tServerHandlerStore, method string, pattern string) *THttpRouterInfo {
	return &THttpRouterInfo{
		handler: handler,
		method:  method,
		pattern: pattern,
	}
}

func (p *THttpRouterInfo) GetMethod() string {
	return p.method
}

func (p *THttpRouterInfo) GetPattern() string {
	return p.pattern
}

func (p *THttpRouterInfo) GetSummary() string {
	return p.handler.summary
}

func (p *THttpRouterInfo) GetTitle() string {
	return p.handler.title
}

var _ http.Handler = (*tServerHandlerStore)(nil)

func newServerHandler(baseRouter string, handler HttpServerHandler) *tServerHandlerStore {
	return &tServerHandlerStore{
		handler:    handler,
		baseRouter: baseRouter,
		fun:        nil,
	}
}

func newServerHandleByFun(baseRouter string, fun HttpServerHandleFun) *tServerHandlerStore {
	return &tServerHandlerStore{
		handler:    nil,
		baseRouter: baseRouter,
		fun:        fun,
	}
}

func (p *tServerHandlerStore) SetTitle(title string) {
	p.title = title
}

func (p *tServerHandlerStore) SetSummary(summary string) {
	p.summary = summary
}

func (p *tServerHandlerStore) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if p.handler == nil && p.fun == nil {
		_, _ = res.Write([]byte("server is run."))
		res.WriteHeader(http.StatusOK)
		return
	}
	httpReq := NewHttpRequest(p.baseRouter, req, res)
	p.doRequest(httpReq)
}

func (p *tServerHandlerStore) doRequest(req *THttpRequest) {
	defer func() {
		if e := recover(); e != nil {
			if st, ok := e.(HttpProcessStatus); ok {
				switch st.Code() {
				case HPS_END_WRITE:
					data := req.Buffer().Bytes()
					if len(data) > 0 {
						_, _ = req.GetResponse().Write(data)
					}
					break
				case HPS_UN_SUPPORT:
					req.WriteHeader(http.StatusHTTPVersionNotSupported)
					break
				case HPS_INNER_ERROR:
					req.WriteHeader(http.StatusInternalServerError)
					if st.Error() != nil {
						xlog.Error(st.Error())
					}
					break
				default:
					if st.Error() != nil {
						xlog.Error(st.Error())
						req.WriteHeader(http.StatusInternalServerError)
					}
				}
			} else if err, ok := e.(error); ok {
				xlog.Error(err)
				req.Clear()
				req.WriteHeader(http.StatusInternalServerError)
			}
		}
	}()
	if p.handler != nil {
		p.handler.Handle(req)
	} else if p.fun != nil {
		p.fun(req)
	}
}
