/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-04 09:05
* Description:
*****************************************************************/

package main

import (
	"errors"
	"fmt"
	"github.com/go-xe2/xthrift/netstream"
)

type tserverHandler struct {
}

func (p *tserverHandler) OnReconnect(conn netstream.StreamConn) {
	panic("implement me")
}

var _ netstream.ServerStreamHandler = (*tserverHandler)(nil)

func (p *tserverHandler) OnRecv(conn netstream.StreamConn, data []byte) {
	fmt.Println("recv data:", string(data))
	fmt.Println("开始返回原数据")
	conn.Send(data)
}

func (p *tserverHandler) OnCall(conn netstream.StreamConn, data []byte) (result []byte, err error) {
	fmt.Println("服务端被调用,输入数据:", string(data))
	return []byte("回复[" + conn.Id() + "]客户端: hello, this is server result data."), errors.New("测试返回错误信息")
}

func (p *tserverHandler) OnConnect(conn netstream.StreamConn) {
	fmt.Println("客户Id:", conn.Id(), "已连接.")
	//fmt.Println("测试调用客户端")
	//result, err := conn.Call([]byte("这是来自服务端的调用"), 2*time.Second)
	//fmt.Println("服务端调用返回:", string(result), ", err:", err)
}

func (p *tserverHandler) OnDisconnect(conn netstream.StreamConn) {
	fmt.Println("客户Id:", conn.Id(), "断开连接.")
}

func (p *tserverHandler) OnSendTo(conn netstream.StreamConn, target netstream.StreamConn, data []byte) {
	fmt.Println("serverHandler onSendTo, fromId:", conn.Id(), ", targetId:", target.Id())
	target.Send(data)
}

func (p *tserverHandler) OnCallTo(conn netstream.StreamConn, target netstream.StreamConn, data []byte) (result []byte, err error) {
	fmt.Println("serverHandler OnCallTo, fromId:", conn.Id(), ", targetId:", target.Id())
	return target.Call(data)
}

func (p *tserverHandler) OnRequest(reqConn netstream.StreamConn, reqId string, namespace string, data []byte) {
	fmt.Println("收到请求 reqId", reqId, ",namespace:", namespace, ", data:", string(data))
	if e := reqConn.Response(reqId, []byte("这是服务端直接返回的数据")); e != nil {
		fmt.Println("服务端返回数据出错:", e)
	}
}

func main() {

	var handler = &tserverHandler{}
	var options = netstream.NewStmServerOptions()
	server, err := netstream.NewStreamServer(":8000", options)
	if err != nil {
		panic(err)
	}
	server.SetHandler(handler)
	fmt.Println("开始监听服务端口:8000")
	if e := server.Serve(); e != nil {
		fmt.Println("server.serve error:", e)
	}
}
