/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-27 17:02
* Description:
*****************************************************************/

package test

import (
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
	"github.com/go-xe2/xthrift/regcenter"
	"testing"
)

func TestPDLStore(t *testing.T) {
	store := regcenter.NewPDLStore("./protocol", ".pdl", true)
	err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	store.SetOnChanged(func(fileName, projName string, namespaces map[string]*pdl.TPDLNamespace) {
		if namespaces == nil {
			xlog.Info("删除文件:", fileName, ", projName:", projName)
		} else {
			xlog.Info("changed file:", fileName, ", projName:", projName)
			xlog.Info("changed namespaces:", namespaces)
		}
		xlog.Info("所有空间:", store.AllNamespaces())
		xlog.Info("所有协议项目文件:", store.AllFiles())
	})

	if err := store.EnableFileWatch(); err != nil {
		xlog.Error(err)
	}
	xlog.Info("所有空间:", store.AllNamespaces())
	xlog.Info("所有协议项目文件:", store.AllFiles())
	select {}
}
