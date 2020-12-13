/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 17:16
* Description:
*****************************************************************/

package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"io"
	"os"
)

type TFileIOFileManager struct {
	files map[string]*os.File
}

var _ FileIOManager = (*TFileIOFileManager)(nil)

func NewFileIOFileManager() FileIOManager {
	return &TFileIOFileManager{
		files: make(map[string]*os.File),
	}
}

func (p *TFileIOFileManager) Create(ns string, fileName string) (io.Writer, error) {
	f, err := xfile.OpenWithFlag(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		return nil, err
	}
	p.files[fileName] = f
	return f, nil
}

func (p *TFileIOFileManager) Close(ns string, fileName string) {
	if f, ok := p.files[fileName]; ok {
		if err := f.Close(); err != nil {
			fmt.Println("close file:", fileName, ", error:", err)
		}
		delete(p.files, fileName)
	}
}
