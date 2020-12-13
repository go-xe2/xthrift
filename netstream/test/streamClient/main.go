/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-04 08:59
* Description:
*****************************************************************/

package main

import (
	"fmt"
	"github.com/go-xe2/xthrift/netstream"
	"sync"
	"time"
)

type tclientHandler struct {
}

var _ netstream.ClientStreamHandler = (*tclientHandler)(nil)

func (p *tclientHandler) OnRecv(conn netstream.StreamConn, data []byte) {
	fmt.Println("recv data:", string(data))
}

func (p *tclientHandler) OnCall(conn netstream.StreamConn, data []byte) (result []byte, err error) {
	fmt.Println("客户端被调用,输入数据:", string(data))
	return []byte("hello, this is client result data."), nil
}

func (p *tclientHandler) OnConnect(conn netstream.StreamConn) {
	fmt.Println("client onConnect.")
}

func (p *tclientHandler) OnReconnect(conn netstream.StreamConn) {
	fmt.Println("client reconnect.")
}

func (p *tclientHandler) OnDisconnect(conn netstream.StreamConn) {
	fmt.Println("client:", conn.Id(), ", onDisconnect")
}

func (p *tclientHandler) OnRequest(reqId string, namespace string, data []byte) {
	fmt.Println("收到请求 reqId:", reqId, ",namespace:", namespace, ", data:", string(data))
}

func (p *tclientHandler) OnResponse(reqId string, data []byte) {
	fmt.Println("收到请求回复 reqId:", reqId, ", data:", string(data))
}

func startClient() string {
	options := netstream.NewStmClientOptions()
	options.SetWriteTimeout(3 * time.Minute)
	options.SetReadTimeout(3 * time.Minute)
	options.SetHeartbeatSpeed(1 * time.Minute)
	options.SetAllowMaxLoss(3)

	client, err := netstream.NewStreamClient("127.0.0.1:8000", options)
	if err != nil {
		panic(err)
	}
	var handler = &tclientHandler{}
	client.SetHandler(handler)
	if e := client.Open(); e != nil {
		fmt.Println("open netClient error:", e)
	}
	fmt.Println("2秒后发送:hello, netstream.")
	time.Sleep(2 * time.Second)
	fmt.Println("开始发送")
	client.Send([]byte("hello, netstream."))
	fmt.Println("数据发送完成.")
	fmt.Println("秒后调用call")
	time.Sleep(2 * time.Second)
	result, e := client.Call([]byte("客户端呼叫数据"))
	fmt.Println("客户端调用返回:", string(result), ", err:", e)
	fmt.Println("准备发送请求")
	if e := client.Request(client.MakeRequestId(), "mnyun.com", []byte("我是客户端调用数据")); e != nil {
		fmt.Println("客户端调用出错:")
	}
	if e := client.Serve(nil); e != nil {
		fmt.Println("客户端已结束，出错:", e)
	}
	return client.Id()
}

func startClientWithCount(wg *sync.WaitGroup, n int) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := startClient()
			fmt.Println("=======>> ^^^^^ 客户端:", id, "已关闭")
		}()
	}
}

func main() {
	// 创建1000个客户端测试
	var wg sync.WaitGroup
	var count = 3
	startClientWithCount(&wg, count)
	wg.Wait()
	fmt.Println(count, "个客户端都已经关闭")
}
