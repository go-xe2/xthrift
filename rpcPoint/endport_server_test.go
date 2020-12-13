/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-18 11:12
* Description:
*****************************************************************/

package rpcPoint

import (
	"testing"
)

func TestNewEndPortServer(t *testing.T) {
	options := NewOptions(":8000")
	options.SetHostPath("./hosts")
	options.SetPdlPath("./pdls")
	options.SetBaseRouter("")
	server := NewEndPointServer(options)
	if err := server.Serve(); err != nil {
		t.Fatal(err)
	}
	select {}
}
