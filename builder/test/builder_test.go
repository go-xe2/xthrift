/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 09:20
* Description:
*****************************************************************/

package test

import (
	"fmt"
	pbuilder "github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/builder/gcontext"
	"github.com/go-xe2/xthrift/pdl"
	"testing"
)

func TestBuildGo(t *testing.T) {
	project, err := pdl.NewFileProject("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto")
	if err != nil {
		t.Fatal(err)
	}
	if err := project.Load(); err != nil {
		t.Fatal(err)
	}

	cxt := gcontext.NewContext("./build/go", "github.com/go-xe2/xthrift/builder/test/build/go")
	builder := pbuilder.NewProtoBuilder(cxt, project)
	_, err = builder.Build()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("生成golang成功")
}
