/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-21 09:59
* Description:
*****************************************************************/

package pdl

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"sort"
)

type tYamlWriter struct {
	indent int
	w      *bytes.Buffer
}

func newYamlWriter() *tYamlWriter {
	return &tYamlWriter{
		w: bytes.NewBuffer([]byte{}),
	}
}

func (p *tYamlWriter) IncIndent() {
	p.indent += 1
}

func (p *tYamlWriter) DecIndent() {
	p.indent -= 1
}

func (p *tYamlWriter) write(str string) {
	szIndent := xstring.MakeStrAndFill(p.indent*2, " ", "")
	str = szIndent + str
	if _, err := p.w.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func (p *tYamlWriter) WriteBegin() error {
	p.write("# 服务协议定义文件\n")
	return nil
}

func (p *tYamlWriter) WriteNamespace(namespace string) error {
	p.write("\n")
	p.write("# 协议命名空间\n")
	p.write("namespace: ")
	p.write(namespace)
	p.write("\n")
	return nil
}

func (p *tYamlWriter) WriteImports(imports []string) error {
	p.write("# 引用协议文件\n")
	p.write("imports:\n")
	p.IncIndent()
	for _, k := range imports {
		p.write("- " + k + "\n")
	}
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteBasicBegin() error {
	p.write("\n")
	p.write("# 服务协议基础数据类型\n")
	return nil
}

func (p *tYamlWriter) WriteBasic(basicTypes []string) error {
	p.write("basic:\n")
	p.IncIndent()
	defer p.DecIndent()
	for _, s := range basicTypes {
		p.write("- " + s)
		p.write("\n")
	}
	return nil
}

func (p *tYamlWriter) WriteBasicEnd() error {
	return nil
}

func (p *tYamlWriter) WriteTypeDefBegin() error {
	p.write("\n")
	p.write("# 数据类型别名定义节点\n")
	p.write("typeDefs:\n")
	p.IncIndent()
	return nil
}

func (p *tYamlWriter) WriteTypeDefs(defs map[string]*FileTypeDef) error {
	for k, v := range defs {
		p.write(k + ": " + v.OrgType.Name() + "\n")
	}
	return nil
}

func (p *tYamlWriter) WriteTypeDefEnd() error {
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteTypesBegin() error {
	p.write("\n")
	p.write("# 数据类型定义节点\n")
	p.write("types:\n")
	p.IncIndent()
	return nil
}

func (p *tYamlWriter) WriteTypes(typs map[string]*FileStruct) error {
	for k, v := range typs {
		if e := p.WriteStructBegin("", k); e != nil {
			return e
		}
		if e := p.WriteStruct("", v); e != nil {
			return e
		}
		if e := p.WriteStructEnd(""); e != nil {
			return e
		}
	}
	return nil
}

func (p *tYamlWriter) WriteTypesEnd() error {
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteInterfacesBegin() error {
	p.write("\n")
	p.write("# 服务接口定义节点\n")
	p.write("interfaces:\n")
	p.IncIndent()
	return nil
}

func (p *tYamlWriter) WriteInterfaces(interfaces map[string]*FileService) error {
	for k, v := range interfaces {
		if e := p.WriteServiceBegin("", k); e != nil {
			return e
		}
		if e := p.WriteService("", v); e != nil {
			return e
		}
		if e := p.WriteServiceEnd(""); e != nil {
			return e
		}
	}
	return nil
}

func (p *tYamlWriter) WriteInterfacesEnd() error {
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteServiceBegin(indent string, svcName string) error {
	p.write("# 服务" + svcName + "定义开始\n")
	p.write(svcName + ":\n")
	p.IncIndent()
	return nil
}

func (p *tYamlWriter) WriteService(indent string, service *FileService) error {
	methods := service.Methods
	for k, m := range methods {
		if e := p.WriteServiceMethodBegin("", k); e != nil {
			return e
		}
		if e := p.WriteServiceMethod("", m); e != nil {
			return e
		}
		if e := p.WriteServiceMethodEnd(""); e != nil {
			return e
		}
	}
	return nil
}

func (p *tYamlWriter) WriteServiceEnd(indent string) error {
	p.DecIndent()
	p.write("# 服务定义节点结束\n")
	return nil
}

func (p *tYamlWriter) WriteServiceMethodBegin(indent string, methodName string) error {
	p.write("# 定义接口方法" + methodName + "开始\n")
	p.write(methodName + ":\n")
	p.IncIndent()
	return nil
}

func (p tYamlWriter) WriteServiceMethod(indent string, method *FileServiceMethod) error {
	p.write("summary: " + method.Summary)
	p.write("\n")
	p.write("# 接口输入参数\n")
	p.write("args:\n")
	p.IncIndent()
	args := method.Args
	sort.Slice(args, func(i, j int) bool {
		return args[i].Id-args[j].Id < 0
	})
	for _, v := range args {
		if e := p.WriteFieldBegin("", v.Name); e != nil {
			return e
		}
		if e := p.WriteField("", v); e != nil {
			return e
		}
		if e := p.WriteFieldEnd(""); e != nil {
			return e
		}
	}
	p.DecIndent()
	p.write("\n")
	p.write("# 接口返回数据类型\n")
	p.write("results: " + method.Result.Name() + "\n")
	if method.Exception != nil && method.Exception.Type != SPD_VOID {
		p.write("# 接口可能抛出异常类型\n")
		p.write("throw: " + method.Exception.Name() + "\n")
	}
	return nil
}

func (p *tYamlWriter) WriteServiceMethodEnd(indent string) error {
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteFieldBegin(indent string, fieldName string) error {
	p.write("# 定义字段" + fieldName + "\n")
	p.write(fieldName + ":\n")
	p.IncIndent()
	return nil
}

func (p *tYamlWriter) WriteField(indent string, field *FileDataField) error {
	p.write("id: " + fmt.Sprintf("%d\n", field.Id))
	p.write("type: " + field.FieldType.Name() + "\n")
	if field.Summary != "" {
		p.write("summary: " + field.Summary + "\n")
	}
	if field.Limit != SPDLimitRequired {
		p.write("limit: " + field.Limit.String() + "\n")
	}
	if field.Rule != "" {
		p.write("valid: " + field.Rule + "\n")
	}
	return nil
}

func (p *tYamlWriter) WriteFieldEnd(indent string) error {
	p.DecIndent()
	p.write("\n")
	return nil
}

func (p *tYamlWriter) WriteStructBegin(indent string, struName string) error {
	p.write("# 定义数据结构" + struName + "\n")
	p.write(struName + ": ")
	p.write("\n")
	p.IncIndent()
	return nil
}

func (p *tYamlWriter) WriteStruct(indent string, stru *FileStruct) error {
	p.write("type: struct\n")
	p.write("summary: " + stru.Summary + "\n")
	p.write("fields:\n")
	p.IncIndent()
	fields := make([]*FileDataField, len(stru.Fields))
	i := 0
	for _, v := range stru.Fields {
		fields[i] = v
		i++
	}
	// 对字段id排序
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Id-fields[j].Id < 0
	})

	for _, v := range fields {
		if e := p.WriteFieldBegin("", v.Name); e != nil {
			return e
		}
		if e := p.WriteField("", v); e != nil {
			return e
		}
		if e := p.WriteFieldEnd(""); e != nil {
			return e
		}
	}
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteStructEnd(indent string) error {
	p.DecIndent()
	return nil
}

func (p *tYamlWriter) WriteEnd() error {
	return nil
}

func (p *tYamlWriter) Data() []byte {
	return p.w.Bytes()
}
