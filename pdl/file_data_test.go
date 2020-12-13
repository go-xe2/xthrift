package pdl

import (
	"encoding/json"
	"testing"
)

func TestNewFileData(t *testing.T) {
	file := NewFileData(nil)
	err := file.OpenFile("/Users/ytx/mx/mny3/github.com/go-xe2/xthrift/pdl/proto/com/mnyun/reg/user/regUserSvc.yaml")
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := json.MarshalIndent(file, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))
}
