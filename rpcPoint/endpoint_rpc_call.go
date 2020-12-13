/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 16:22
* Description:
*****************************************************************/

package rpcPoint

import (
	"context"
	"fmt"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
)

// 检查服务接口是否存在
func (p *TEndPointServer) ApiExists(fullSvcName string, method string) (exists bool, service *pdl.FileService, m *pdl.FileServiceMethod) {
	svc := p.regCenter.PDLQuery().GetServiceByFullName(fullSvcName)
	if svc == nil {
		return false, nil, nil
	}
	if m := svc.QryMethod(method); m == nil {
		return false, svc, nil
	} else {
		return true, svc, m
	}
}

// 调用服务接接口, 实现RpcService接口方法
// 输入数据封包格式应与提供的服务一致
func (p *TEndPointServer) RpcCall(ctx context.Context, fullSvcName string, method string, seqId int32, input []byte) (result []byte, err error) {
	svc := p.GetPdlQuery().GetServiceByFullName(fullSvcName)
	if svc == nil {
		return nil, fmt.Errorf("服务%s不存在", fullSvcName)
	}
	if m := svc.QryMethod(method); m == nil {
		return nil, fmt.Errorf("服务%s不存在接口%s", fullSvcName, method)
	}
	client, e := p.GetServiceHostClient(fullSvcName)
	if e != nil {
		return nil, e
	}
	// 未使用thrift.TFrame封包
	//outBufTrans := thrift.NewTMemoryBuffer()

	xlog.Debug("准备调用服务:", fullSvcName, method, ", 服务地址:", client.host, ", 端口:", client.port)
	return client.Call(ctx, method, seqId, input)
}
