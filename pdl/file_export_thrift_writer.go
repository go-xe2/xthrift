/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-21 10:19
* Description:
*****************************************************************/

package pdl

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"sort"
	"strings"
)

type tThriftWriter struct {
	namespace string
	indent    int
	w         *bytes.Buffer
	lang      string
}

func newThriftWriter(namespace string, lang string) *tThriftWriter {
	return &tThriftWriter{
		namespace: namespace,
		indent:    0,
		lang:      lang,
		w:         bytes.NewBuffer([]byte{}),
	}
}

func (p *tThriftWriter) IncIndent() {
	p.indent += 1
}

func (p *tThriftWriter) DecIndent() {
	p.indent -= 1
}

func (p *tThriftWriter) write(str string) {
	szIndent := xstring.MakeStrAndFill(p.indent, "\t", "")
	str = szIndent + str
	if _, err := p.w.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func (p *tThriftWriter) WriteBegin() error {
	return nil
}

func (p *tThriftWriter) WriteNamespace(namespace string) error {
	szLang := p.lang
	if szLang == "" {
		szLang = "go"
	}
	langs := strings.Split(szLang, ",")
	for _, lang := range langs {
		p.write(fmt.Sprintf("namespace "+lang+" %s\n", namespace))
	}
	return nil
}

func (p *tThriftWriter) WriteImports(imports []string) error {
	for _, k := range imports {
		p.write(fmt.Sprintf("include '%s'\n", k+".thrift"))
	}
	return nil
}

func (p *tThriftWriter) WriteBasicBegin() error {
	return nil
}

func (p *tThriftWriter) WriteBasic(basicTypes []string) error {
	return nil
}

func (p *tThriftWriter) WriteBasicEnd() error {
	return nil
}

func (p *tThriftWriter) WriteTypeDefBegin() error {
	return nil
}

func (p *tThriftWriter) writeBasicTypeDef(thriftType string, typ TProtoBaseType) {
	if thriftType != typ.String() {
		p.write("typedef " + thriftType + " " + typ.String() + "\n")
	}
}

func (p *tThriftWriter) WriteTypeDefs(defs map[string]*FileTypeDef) error {
	// 写入pro的定义
	p.writeBasicTypeDef("string", SPD_STR)
	p.writeBasicTypeDef("bool", SPD_BOOL)
	p.writeBasicTypeDef("i8", SPD_I08)
	p.writeBasicTypeDef("i16", SPD_I16)
	p.writeBasicTypeDef("i32", SPD_I32)
	p.writeBasicTypeDef("i64", SPD_I64)
	p.writeBasicTypeDef("double", SPD_DOUBLE)

	for k, v := range defs {
		p.write(fmt.Sprintf("typedef %s %s\n", v.OrgType.FullName(p.namespace), k))
	}
	return nil
}

func (p *tThriftWriter) WriteTypeDefEnd() error {
	return nil
}

func (p *tThriftWriter) WriteTypesBegin() error {
	return nil
}

func (p *tThriftWriter) WriteTypes(types map[string]*FileStruct) error {
	for k, v := range types {
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

func (p *tThriftWriter) WriteTypesEnd() error {
	return nil
}

func (p *tThriftWriter) WriteInterfacesBegin() error {
	return nil
}

func (p *tThriftWriter) WriteInterfaces(interfaces map[string]*FileService) error {
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

func (p *tThriftWriter) WriteInterfacesEnd() error {
	return nil
}

func (p *tThriftWriter) WriteServiceBegin(indent string, svcName string) error {
	p.write(fmt.Sprintf("service %s {\n", svcName))
	p.IncIndent()
	return nil
}

func (p *tThriftWriter) WriteService(indent string, service *FileService) error {
	for k, method := range service.Methods {
		if e := p.WriteServiceMethodBegin("", k); e != nil {
			return e
		}
		if e := p.WriteServiceMethod("", method); e != nil {
			return e
		}
		if e := p.WriteServiceMethodEnd(""); e != nil {
			return e
		}
	}
	return nil
}

func (p *tThriftWriter) WriteServiceEnd(indent string) error {
	p.DecIndent()
	p.write("}\n")
	return nil
}

func (p *tThriftWriter) WriteServiceMethodBegin(indent string, methodName string) error {
	return nil
}

func (p *tThriftWriter) WriteServiceMethod(indent string, method *FileServiceMethod) error {
	if method.Summary != "" {
		p.write("// " + method.Summary + "\n")
	}
	s := fmt.Sprintf("%s %s(", method.Result.FullName(p.namespace), method.Name)
	i := 0
	size := len(method.Args)
	for _, arg := range method.Args {
		s += fmt.Sprintf("%d:", arg.Id)

		if arg.Limit == SPDLimitOptional {
			s += "optional "
		}
		typeName := arg.FieldType.FullName(p.namespace)
		s += fmt.Sprintf("%s %s", typeName, arg.Name)
		i++
		if i < size {
			s += ","
		}
	}
	s += ")"
	if method.Exception != nil && method.Exception.Type != SPD_VOID {
		s += fmt.Sprintf("(1: %s err)", method.Exception.FullName(p.namespace))
	}
	s += ";\n"
	p.write(s)
	return nil
}

func (p *tThriftWriter) WriteServiceMethodEnd(indent string) error {
	return nil
}

func (p *tThriftWriter) WriteFieldBegin(indent string, fieldName string) error {
	return nil
}

func (p *tThriftWriter) WriteField(indent string, field *FileDataField) error {
	if field.Summary != "" {
		p.write("// " + field.Summary + "\n")
	}
	s := fmt.Sprintf("%d:", field.Id)
	if field.Limit == SPDLimitOptional {
		s += "optional "
	}
	typeName := field.FieldType.FullName(p.namespace)
	s += fmt.Sprintf("%s %s", typeName, field.Name)
	s += ";\n"
	p.write(s)
	return nil
}

func (p *tThriftWriter) WriteFieldEnd(indent string) error {
	return nil
}

func (p *tThriftWriter) WriteStructBegin(indent string, struName string) error {

	return nil
}

func (p *tThriftWriter) WriteStruct(indent string, stru *FileStruct) error {
	if stru.Summary != "" {
		p.write("// " + stru.Summary + "\n")
	}
	p.write(fmt.Sprintf("struct %s {\n", stru.Type.TypName))
	p.IncIndent()

	fields := make([]*FileDataField, len(stru.Fields))
	i := 0
	for _, f := range stru.Fields {
		fields[i] = f
		i++
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Id-fields[j].Id < 0
	})
	for _, f := range fields {
		if e := p.WriteFieldBegin("", f.Name); e != nil {
			return e
		}
		if e := p.WriteField("", f); e != nil {
			return e
		}
		if e := p.WriteFieldEnd(""); e != nil {
			return e
		}
	}
	p.DecIndent()
	p.write("}\n")
	return nil
}

func (p *tThriftWriter) WriteStructEnd(indent string) error {
	return nil
}

func (p *tThriftWriter) WriteEnd() error {
	return nil
}

func (p *tThriftWriter) WriteComment(s string) {
	p.write(fmt.Sprintf("// %s\n", s))
}

func (p *tThriftWriter) Data() []byte {
	return p.w.Bytes()
}
