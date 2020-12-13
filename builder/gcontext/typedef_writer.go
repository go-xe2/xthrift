/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-13 14:30
* Description:
*****************************************************************/

package gcontext

import (
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
)

type TTypedefCodeWriter struct {
	*TWriter
}

var _ builder.TypedefCodeWriter = (*TTypedefCodeWriter)(nil)

func NewTypedefCodeWriter(cxt *TContext, fileName string) (w *TTypedefCodeWriter, err error) {
	inst := &TTypedefCodeWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return
	}
	return inst, nil
}

func (p *TTypedefCodeWriter) WriteDef(ns *pdl.FileNamespace, def *pdl.FileTypeDef) (ident string, err error) {
	ident = def.Name
	p.Write("type ")
	p.Write(def.Name + "\t")
	p.Write(p.GenDataTypeCode(ns, def.OrgType))
	p.Write("\n")
	return ident, nil
}
