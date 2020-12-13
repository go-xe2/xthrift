package pdl

import (
	"encoding/json"
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"os"
	"testing"
)

func TestNewFileProject(t *testing.T) {
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
	data, err := json.MarshalIndent(project, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	f, err := xfile.OpenWithFlag("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto/demo.pdl", os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = project.SaveProject(f)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("保存协议项目成功")

	fmt.Println("测试打开协议项目文件")
	proj, err := NewFileProject("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto1")
	if err != nil {
		t.Fatal(err)
	}
	f1, err := xfile.OpenWithFlag("/Users/mx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto/demo.pdl", os.O_RDONLY)
	if err != nil {
		t.Fatal(err)
	}
	defer f1.Close()
	err = proj.LoadProject(f1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("打开协议项目成功")
	if err = proj.Check(); err != nil {
		t.Fatal(err)
	}
	fmt.Println("项目内容:")
	data, err = json.MarshalIndent(proj, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
