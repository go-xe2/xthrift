/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-19 12:00
* Description:
*****************************************************************/

package registerCli

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"testing"
)

func TestNewRegisterClient(t *testing.T) {
	client := NewRegisterClient("http://127.0.0.1:8000/pdl/reg")
	client.SetWorkPath("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto")
	fileName1 := "/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto/com/mnyun/reg/user/regUserSvc.yaml"
	//fileName1 := "/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto/com/mnyun/reg/user/tmp.yaml"
	if !xfile.Exists(fileName1) {
		fmt.Println("========> file not exists.")
	}
	err := client.Register("192.168.1.111", 3003, fileName1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("注册成功")
}
