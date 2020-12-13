/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 12:43
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/xthrift/builder"
	"os"
)

type TWriter struct {
	file         *os.File
	fileFullName string
	fileName     string
	buf          *bytes.Buffer
	cxt          *TContext
	// 使用到的命名空间
	imports   map[string]bool
	namespace string
}

var _ builder.CodeWriter = (*TWriter)(nil)

func NewWriter(cxt *TContext, fileName string) (w *TWriter, err error) {
	dir := xfile.Dir(fileName)
	if !xfile.Exists(dir) {
		if err = xfile.Mkdir(dir); err != nil {
			return
		}
	}
	inst := &TWriter{
		fileFullName: fileName,
		cxt:          cxt,
		buf:          bytes.NewBufferString(""),
		imports:      make(map[string]bool),
	}
	if inst.file, err = xfile.OpenWithFlag(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR); err != nil {
		return
	}
	inst.fileName = xfile.Basename(fileName)
	return inst, nil
}
