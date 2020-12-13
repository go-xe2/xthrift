/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 15:34
* Description:
*****************************************************************/

package pdl

import (
	"encoding/json"
	"testing"
)

func TestNewPDLContainer(t *testing.T) {
	container := NewPDLContainer()
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

	err = container.Install(project)
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(container, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
