/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-27 16:12
* Description:
*****************************************************************/

package main

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/lib/go/test/comm"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"golang.org/x/net/context"
	"sync"
	"sync/atomic"
	"time"
)

type TServerVisitor struct {
	closed int32
	mu     sync.RWMutex
	wg     sync.WaitGroup
	trans  thrift.TServerTransport
}

func NewServerVisitor(addr string) (*TServerVisitor, error) {
	sck, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return nil, err
	}
	inst := &TServerVisitor{
		trans: sck,
	}
	return inst, nil
}

func (p *TServerVisitor) innerAccept() (int32, error) {
	client, err := p.trans.Accept()
	fmt.Println("innerAccept client:", client, ", err:", err)
	p.mu.Lock()
	defer p.mu.Unlock()
	closed := atomic.LoadInt32(&p.closed)
	if closed != 0 {
		return closed, nil
	}
	if err != nil {
		return 0, err
	}
	if client != nil {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			if err := p.processRequests(client); err != nil {
				xlog.Error("error processing request:", err)
			}
		}()
	}
	return 0, nil
}

func (p *TServerVisitor) AcceptLoop() error {
	for {
		closed, err := p.innerAccept()
		if err != nil {
			return err
		}
		if closed != 0 {
			return nil
		}
	}
}

func (p *TServerVisitor) Serve() error {
	if err := p.trans.Listen(); err != nil {
		return err
	}
	if err := p.AcceptLoop(); err != nil {
		return err
	}
	return nil
}

func (p *TServerVisitor) processRequests(client thrift.TTransport) error {
	// 连接成功时，2称后，调用客户端方法测试
	var n = 0
	xlog.Info("每隔5秒调用客户端方法,5次后退出")
	for {
		if n >= 5 {
			break
		}
		fmt.Printf("第%d次调用", n)
		time.Sleep(5 * time.Second)
		xlog.Info("开始调用客户端方法")
		n++
		var cxt = context.Background()
		protocol := xthrift.NewNamesapceProtocolFactory("demo.hello", xthrift.NewBinaryProtocolEx(client))
		clientSvc := comm.NewSayHelloClient(client, protocol, protocol)
		result, err := clientSvc.SayHello(cxt, "我是服务端，我调用客户端方法")
		if err != nil {
			xlog.Error("调用出错:", err)
		} else {
			xlog.Info("服务端调用客户端方法成功，返回:", result)
		}
	}
	fmt.Println("断开客户端.")
	client.Close()
	return nil
}

func main() {
	svcSocket, err := NewServerVisitor(":8001")
	if err != nil {
		panic(err)
	}
	xlog.Info("中心服务启动成功，监听端口8001")
	if e := svcSocket.Serve(); e != nil {
		xlog.Error(e)
	}
}
