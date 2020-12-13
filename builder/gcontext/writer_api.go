/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 12:48
* Description:
*****************************************************************/

package gcontext

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
	"strings"
)

func (p *TWriter) Context() builder.Context {
	return p.cxt
}

func (p *TWriter) ModuleName() string {
	return p.cxt.GetModuleName()
}

func (p *TWriter) FileName() string {
	return p.fileName
}

func (p *TWriter) Write(str string) {
	if _, err := p.buf.WriteString(str); err != nil {
		panic(err)
	}
}

func (p *TWriter) Flush() error {
	_, s := pdl.NamespaceLastName(p.namespace)
	if _, err := p.file.WriteString(fmt.Sprintf("\npackage %s\n", s)); err != nil {
		return err
	}
	importNs := make([]string, 0)

	for k, b := range p.imports {
		name := ""
		if b {
			s, f := pdl.NamespaceToPath(k)
			name = xfile.Join(s, f)
			if p.ModuleName() != "" {
				name = xfile.Join(p.ModuleName(), name)
			}
		} else {
			name = k
		}
		importNs = append(importNs, name)
	}
	sort.Slice(importNs, func(i, j int) bool {
		return strings.Compare(importNs[i], importNs[j]) < 0
	})

	if len(importNs) > 0 {
		if _, err := p.file.WriteString("\n\nimport (\n"); err != nil {
			return err
		}
		for _, k := range importNs {
			if _, err := p.file.WriteString(fmt.Sprintf("\t\"%s\"\n", k)); err != nil {
				return err
			}
		}
		if _, err := p.file.WriteString(")\n"); err != nil {
			return err
		}
	}
	if _, err := p.file.WriteString(p.buf.String()); err != nil {
		return err
	}
	return nil
}

func (p *TWriter) Import(namespace string, inner bool) {
	if _, ok := p.imports[namespace]; !ok {
		p.imports[namespace] = inner
	}
}

func (p *TWriter) WriteNamespace(namespace string) error {
	p.namespace = namespace
	return nil
}

func (p *TWriter) Close() error {
	if p.file != nil {
		return p.file.Close()
	}
	return nil
}

// 获取标识符，在数据交互层中的名称
// 返回驼峰写法转下划线连接的字符串
func (p *TWriter) GetIdentTransportName(name string) string {
	return xstring.Camel2UnderScore(name, "_")
}

func (p *TWriter) getDataTypeNamespace(currentNs *pdl.FileNamespace, dataType *pdl.FileDataType) []string {
	switch dataType.Type {
	case pdl.SPD_LIST:
		return p.getDataTypeNamespace(currentNs, dataType.ElemType)
	case pdl.SPD_SET:
		return p.getDataTypeNamespace(currentNs, dataType.ElemType)
	case pdl.SPD_MAP:
		a1 := p.getDataTypeNamespace(currentNs, dataType.KeyType)
		a2 := p.getDataTypeNamespace(currentNs, dataType.ValType)
		result := make([]string, 0)
		for _, s := range a1 {
			result = append(result, s)
		}
		for _, s := range a2 {
			result = append(result, s)
		}
		return result
	default:
		if dataType.Namespace == "" {
			return nil
		}
		if dataType.Namespace != currentNs.Namespace {
			return []string{dataType.Namespace}
		}
	}
	return nil
}
