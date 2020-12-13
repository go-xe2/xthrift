/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-17 15:52
* Description:
*****************************************************************/

package xhttpServer

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/xtime"
	"github.com/go-xe2/x/utils/xrand"
	"mime/multipart"
	"os"
)

type File interface {
	Header() *multipart.FileHeader
	Size() int64
	Ext() string
	FileName() string
	Save(fileName string) error
	MakeFileName() string
}

type THttpFile struct {
	header *multipart.FileHeader
}

var _ File = (*THttpFile)(nil)

func NewHttpFile(header *multipart.FileHeader) *THttpFile {
	return &THttpFile{
		header: header,
	}
}

func (p *THttpFile) MakeFileName() string {
	return fmt.Sprintf("%s_%s.%s", xtime.Now().String(), xrand.Str(5), p.Ext())
}

func (p *THttpFile) Header() *multipart.FileHeader {
	return p.header
}

func (p *THttpFile) Size() int64 {
	return p.header.Size
}

func (p *THttpFile) Ext() string {
	return xfile.Ext(p.FileName())
}

func (p *THttpFile) FileName() string {
	return p.header.Filename
}

func (p *THttpFile) Save(fileName string) error {
	dir := xfile.Dir(fileName)
	if !xfile.Exists(dir) {
		if err := xfile.Mkdir(dir); err != nil {
			return err
		}
	}
	f, err := p.header.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建输出文件
	out, err := xfile.OpenWithFlag(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 1024)
	total := p.header.Size
	for {
		if total <= 0 {
			break
		}
		n, err := f.Read(buf)
		if err != nil {
			return err
		}
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
			total = total - int64(n)
		}
	}
	return nil
}
