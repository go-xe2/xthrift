/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-27 10:54
* Description:
*****************************************************************/

package test

import (
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/regcenter"
	"testing"
)

func TestHostStore(t *testing.T) {
	store := regcenter.NewHostStore("./hostStore", ".host", true)
	if err := store.Load(); err != nil {
		xlog.Error(err)
	}
	//store.AddHost("123456", "com.mnyun.demo", "192.168.0.11", 3000)
	//store.AddHost("123456", "com.mnyun.demo", "192.168.0.13", 3003)
	//store.AddHost("222222", "com.mnyun.pro", "192.168.0.12", 3000)
	if err := store.SaveHosts(); err != nil {
		t.Error(err)
	}
	xlog.Info(store.AllHosts())
	xlog.Info("com.mnyun.demo:", store.GetSvcHosts("com.mnyun.demo"))
	xlog.Info("com.mnyun.pro", store.GetSvcHosts("com.mnyun.pro"))
	xlog.Info("project 123456 exists:", store.HasProject("123456"))
	xlog.Info("project 333333 exists:", store.HasProject("333333"))

	store.SetOnChanged(func(store *regcenter.THostStore, fileId int) {
		if fileId == -1 {
			xlog.Info("changed: delete host")
		} else {
			xlog.Info("changed:", store.FileHosts(fileId))
		}
	})
	select {}
}
