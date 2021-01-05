/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 16:14
* Description:
*****************************************************************/

package gateway

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"github.com/go-xe2/xthrift/pdl"
	"github.com/go-xe2/xthrift/xhttpServer"
	"strings"
	"sync"
	"sync/atomic"
)

type TSvrHttpHandler struct {
	rootNamespace string
	baseRouter    string
	qry           pdl.PDLQuery
	rpc           RpcService
	seqId         int32
	rw            sync.RWMutex
}

var _ xhttpServer.HttpServerHandler = (*TSvrHttpHandler)(nil)

func NewSvrHttpHandler(qry pdl.PDLQuery, rpc RpcService, baseRouter string, rootNamespace string) *TSvrHttpHandler {
	inst := &TSvrHttpHandler{
		rpc:           rpc,
		qry:           qry,
		rootNamespace: rootNamespace,
		baseRouter:    baseRouter,
		seqId:         0,
	}
	atomic.StoreInt32(&inst.seqId, 0)
	return inst
}

func (p *TSvrHttpHandler) Handle(req *xhttpServer.THttpRequest) {
	router := req.GetRouterName()
	if p.baseRouter != "" {
		if xstring.StartWith(router, p.baseRouter) {
			router = router[len(p.baseRouter):]
		}
	}
	items := make([]string, 0)
	tmp := strings.Split(router, "/")
	for _, s := range tmp {
		if s != "" {
			items = append(items, s)
		}
	}
	routeCount := len(items)
	if routeCount == 0 {
		req.ReturnJson(map[string]interface{}{
			"status": 1,
			"msg":    "请求地址无效",
		})
	}

	var svcName, methodName = "", ""
	if routeCount > 1 {
		svcName = p.rootNamespace + strings.Join(items[:len(items)-1], ".")
		methodName = items[len(items)-1]
	} else {
		methodName = items[0]
		svcName = p.rootNamespace
	}
	fullName := fmt.Sprintf("%s.%s", svcName, methodName)
	exists, svc, method := p.rpc.ApiExists(svcName, methodName)
	if !exists {
		req.ReturnJson(map[string]interface{}{
			"status": 1,
			"msg":    fmt.Sprintf("服务接口%s.%s不存在", svcName, methodName),
		})
	}
	ctx := context.Background()

	protoFac := xthrift.NewBinaryProtocolExFactory()
	p.rw.Lock()
	defer p.rw.Unlock()
	seqId := atomic.AddInt32(&p.seqId, 1)
	inArgs, err := MakeCallMessageFromMap(ctx, p.qry, fullName, seqId, svc, method, protoFac, req.GetParams())
	if err != nil {
		req.ReturnJson(map[string]interface{}{
			"status": 1,
			"msg":    fmt.Sprintf("%s", err.Error()),
		})
	}
	output, err := p.rpc.RpcCall(ctx, svcName, methodName, p.seqId, inArgs)
	if err != nil {
		req.ReturnJson(map[string]interface{}{
			"status": 1,
			"msg":    fmt.Sprintf("%s", err.Error()),
		})
	}
	xlog.Debug("rpc result data size:", len(output))
	outTrans := &thrift.TMemoryBuffer{Buffer: bytes.NewBuffer(output)}
	outFrameTrans := thrift.NewTFramedTransport(outTrans)
	outProto := xthrift.NewBinaryProtocolEx(outFrameTrans)

	//if _, err := outTrans.Write(output); err != nil {
	//	xlog.Error(err)
	//	req.ReturnJson(map[string]interface{}{
	//		"status": 1,
	//		"msg":    "内部服务错误",
	//	})
	//}

	resultData, err := MapFromCallReply(p.qry, fullName, svc, seqId, method, outProto)
	if err != nil {
		xlog.Error(err)
		req.ReturnJson(map[string]interface{}{
			"status": 1,
			"msg":    fmt.Sprintf("读取返回值出错:%s", err),
		})
	}
	if mp, ok := resultData.(map[string]interface{}); ok {
		if _, ok1 := mp["status"]; ok1 {
			req.ReturnJson(resultData)
			return
		}
	}
	req.ReturnJson(map[string]interface{}{
		"status":  0,
		"msg":     "",
		"content": resultData,
	})
}
