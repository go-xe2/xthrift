/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 17:20
* Description:
*****************************************************************/

package pdl

import (
	"fmt"
	"testing"
)

func TestFileProject_Export(t *testing.T) {
	project, err := NewFileProject("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto")
	if err != nil {
		t.Fatal(err)
	}
	if err := project.Load(); err != nil {
		t.Fatal(err)
	}
	errs := project.Errors()
	if len(errs) > 0 {
		t.Log(errs)
	}
	ioMgr := NewFileIOFileManager()
	jsonExport := NewFileExportJson("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/jsonProto", ioMgr)
	if err := project.Export(jsonExport); err != nil {
		t.Fatal(err)
	}
	fmt.Println("导出json格式成功")

	yamlExport := NewFileExportYaml("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/yamlProto", ioMgr)
	if err := project.Export(yamlExport); err != nil {
		t.Fatal(err)
	}
	fmt.Println("导出yaml格式成功")

	thriftExport := NewFileExportThrift("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/thriftProto", ioMgr, "go,java,php")
	if err := project.Export(thriftExport); err != nil {
		t.Fatal(err)
	}
	fmt.Println("导出thrift格式成功")
}
