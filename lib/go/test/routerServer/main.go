/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-29 16:14
* Description:
*****************************************************************/

package main

import (
	"fmt"
	"github.com/go-xe2/xthrift/lib/go/xthrift/router"
)

func main() {
	processorFactory := router.NewServiceProcessorFactory()
	svc, err := router.NewServer(":8001", ":8000", processorFactory)
	if err != nil {
		panic(err)
	}
	fmt.Println("准备启动服务路由服务器，注册服务地址:", ":8001", "接口服务地址:", ":8000")
	if err := svc.Serve(); err != nil {
		fmt.Println("routerServer error:", err)
	}
}
