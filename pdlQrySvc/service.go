/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 16:59
* Description:
*****************************************************************/

package pdlQrySvc

import (
	"github.com/go-xe2/xthrift/regcenter"
	"github.com/go-xe2/xthrift/xhttpServer"
)

type TService struct {
	server    *xhttpServer.THttpServer
	regCenter *regcenter.TRegCenter
}

func NewService(svr *xhttpServer.THttpServer, regCt *regcenter.TRegCenter) *TService {
	return &TService{
		server:    svr,
		regCenter: regCt,
	}
}
