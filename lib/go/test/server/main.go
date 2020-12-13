/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-27 14:56
* Description:
*****************************************************************/

package main

import (
	"fmt"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/lib/go/test/comm/helloHandler"
	"github.com/go-xe2/xthrift/lib/go/test/comm/service/pdlSvr"
)

func main() {
	svc, err := pdlSvr.NewPdlSvrServer(":3000")
	if err != nil {
		xlog.Error(err)
	}
	handler := &helloHandler.THelloServiceHandler{}
	svc.RegisterDemoHelloService(handler)

	// 注册多个命名空间
	defer svc.Stop()
	fmt.Println("服务启动成功,端口号:", 3000)
	svc.Serve()
}
