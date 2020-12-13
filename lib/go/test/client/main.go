/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-27 14:58
* Description:
*****************************************************************/

package main

import (
	"bufio"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/lib/go/test/comm/service/com/mnyun/demo"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"golang.org/x/net/context"
	"os"
	"strings"
)

func ClientSayHello(c *demo.HelloServiceClient) {
	result, err := c.SayHello(context.Background(), "yy")
	if err != nil {
		fmt.Println("call SayHello error:", err)
	} else {
		fmt.Println("call SayHello success", "status:", result.Status, ", msg:", result.Msg, ", data:", result.Content)
	}
}

const hostAddr = "127.0.0.1:3001"

func main() {
	socket, err := thrift.NewTSocket(hostAddr)
	if err != nil {
		xlog.Error(err)
	}
	trans := thrift.NewTFramedTransport(socket)
	if err := trans.Open(); err != nil {
		xlog.Error(err)
	}
	fmt.Println("连接服务端成功，服务端:", hostAddr)

	fac := xthrift.NewNamespaceProtocolFactory("com.mnyun.demo.helloService", xthrift.NewBinaryProtocolEx(trans))
	c := demo.NewHelloServiceClient(trans, fac, fac)

	fmt.Println("连接服务成功,输入r重新调用一次，输入c退出!")
	inputReader := bufio.NewReader(os.Stdin) //创建一个读取器，并将其与标准输入绑定。
	for {
		input, err := inputReader.ReadString('\n') //读取器对象提供一个方法 ReadString(delim byte) ，该方法从输入中读取内容，直到碰到 delim 指定的字符，然后将读取到的内容连同 delim 字符一起放到缓冲区。
		if err != nil {
			fmt.Printf("读取输入出错: %s", err)
		} else {
			input = strings.Trim(input, " ")
			input = strings.Replace(input, "\n", "", -1)
			fmt.Println("命令:", input)
			if input == "r" {
				ClientSayHello(c)
			} else if input == "c" {
				break
			}
		}
	}
}
