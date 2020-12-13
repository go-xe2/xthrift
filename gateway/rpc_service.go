/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 16:26
* Description:
*****************************************************************/

package gateway

import (
	"github.com/go-xe2/xthrift/pdl"
	"golang.org/x/net/context"
)

// 远程调用服务接口
type RpcService interface {
	ApiExists(fullSvcName string, method string) (exists bool, service *pdl.FileService, m *pdl.FileServiceMethod)
	RpcCall(ctx context.Context, fullSvcName string, method string, seqId int32, input []byte) (result []byte, err error)
}
