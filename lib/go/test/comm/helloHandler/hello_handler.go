/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-24 14:07
* Description:
*****************************************************************/

package helloHandler

import (
	"fmt"
	"github.com/go-xe2/xthrift/lib/go/test/comm/service/com/mnyun/demo"
)

type THelloServiceHandler struct {
}

func (THelloServiceHandler) SayHello(name string) *demo.HelloResult {
	result := demo.NewHelloResult()
	result.Status = 0
	result.Msg = "调用成功"
	result.Content = fmt.Sprintf("你好,%s!!", name)
	return result
}

var _ demo.HelloService = (*THelloServiceHandler)(nil)
