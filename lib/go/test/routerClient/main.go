/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-29 16:42
* Description:
*****************************************************************/

package main

import (
	"fmt"
	"github.com/go-xe2/xthrift/lib/go/test/comm"
	"github.com/go-xe2/xthrift/lib/go/xthrift/router"
	"time"
)

func main() {
	helloHandler := comm.NewSayHelloService()
	processor1 := comm.NewSayHelloServiceProcessor(helloHandler)
	client, err := router.NewRouterClient("127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	err = client.RegisterLocalService("demo.hello", processor1)
	if err != nil {
		panic(err)
	}
	if err := client.Open(); err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("10秒后获取已在远程注册的服务")
		time.Sleep(10 * time.Second)
		fmt.Println("开始获取注册的服务列表:")
		items, err := client.GetNamespaces()
		fmt.Println("获取结果：", items, ", 错误:", err)
	}()
	if err := client.Serve(); err != nil {
		fmt.Println("err:", err)
	}
}
