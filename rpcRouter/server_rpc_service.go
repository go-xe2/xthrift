/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-08 12:00
* Description:
*****************************************************************/

package rpcRouter

import (
	"github.com/go-xe2/xthrift/gateway"
	"github.com/go-xe2/xthrift/pdl"
	"golang.org/x/net/context"
)

var _ gateway.RpcService = (*TRouterServer)(nil)

func (p *TRouterServer) ApiExists(fullSvcName string, method string) (exists bool, service *pdl.FileService, m *pdl.FileServiceMethod) {
	svc := p.center.PDLQuery().GetServiceByFullName(fullSvcName)
	if svc == nil {
		return false, nil, nil
	}
	if m := svc.QryMethod(method); m == nil {
		return false, svc, nil
	} else {
		return true, svc, m
	}
}

func (p *TRouterServer) RpcCall(ctx context.Context, fullSvcName string, method string, seqId int32, input []byte) (result []byte, err error) {
	return p.RouterCall(ctx, fullSvcName, method, seqId, input)
}
