/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-13 14:34
* Description:
*****************************************************************/

package gcontext

import (
	"errors"
	"fmt"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/builder/comm"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
)

type TStructCodeWriter struct {
	*TWriter
	fields []*pdl.FileDataField
	// 指定引用的字段
	ptrFields []*pdl.FileDataField
}

var _ builder.StructCodeWriter = (*TStructCodeWriter)(nil)

func NewStructCodeWriter(cxt *TContext, fileName string) (w *TStructCodeWriter, err error) {
	inst := &TStructCodeWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

//func (p *TStructCodeWriter) WriteInclude(namespaces []string, files map[string]string) error {
//	if err := p.TWriter.WriteInclude(namespaces, files); err != nil {
//		return err
//	}
//	p.Write("import \"github.com/apache/thrift/lib/go/thrift\"\n")
//	//p.Write("import \"github.com/go-xe2/xthrift/lib/go/xthrift\"\n")
//	p.Write("import \"github.com/go-xe2/x/type/t\"\n")
//	p.Write("import \"github.com/go-xe2/xthrift/pdl\"\n")
//	return nil
//}

func (p *TStructCodeWriter) WriteStructBegin(ns *pdl.FileNamespace, stru *pdl.FileStruct) (ident string, err error) {
	// 处理字段，对字段排序处理及检出指针字段
	ident = p.GenStructNameCode(stru)

	size := len(stru.Fields)
	p.fields = make([]*pdl.FileDataField, size)
	p.ptrFields = make([]*pdl.FileDataField, 0)
	i := 0

	for _, f := range stru.Fields {
		p.fields[i] = f
		i++
		if (f.Limit == pdl.SPDLimitOptional) || f.FieldType.Type == pdl.SPD_STRUCT ||
			f.FieldType.Type == pdl.SPD_EXCEPTION || f.FieldType.Type == pdl.SPD_MAP ||
			f.FieldType.Type == pdl.SPD_LIST || f.FieldType.Type == pdl.SPD_SET {
			p.ptrFields = append(p.ptrFields, f)
		}
	}

	// 排序处理
	sort.Slice(p.fields, func(i, j int) bool {
		return p.fields[i].Id-p.fields[j].Id < 0
	})

	p.Write("\n")
	p.Write(fmt.Sprintf("type %s struct {\n", ident))
	return ident, nil
}

func (p *TStructCodeWriter) writerStructField(ns *pdl.FileNamespace, stru *pdl.FileStruct, field *pdl.FileDataField) error {
	userNamespaces := p.getDataTypeNamespace(ns, field.FieldType)
	for _, s := range userNamespaces {
		p.Import(s, true)
	}

	transName := p.GetIdentTransportName(field.Name)
	p.Write("\t")
	p.Write(p.GenFieldNameCode(field))
	p.Write(" ")
	p.Write(p.GenFieldDefineTypeCode(ns, field))
	p.Write("`")
	p.Write(fmt.Sprintf("thrift:\"%s,%d,%s\"", transName, field.Id, field.Limit))
	p.Write(fmt.Sprintf(" json:\"%s\"", transName))
	p.Write("`\n")
	return nil
}

func (p *TStructCodeWriter) WriteStruct(ns *pdl.FileNamespace, stru *pdl.FileStruct) error {
	if p.fields == nil {
		return errors.New("未正常调用WriteStructBegin初始化或已经调用WriteStructEnd结束创建")
	}
	p.Write("\t*pdl.TDynamicStructBase\n")
	p.Write("\n")
	size := len(p.fields)
	for i := 0; i < size; i++ {
		field := p.fields[i]
		if err := p.writerStructField(ns, stru, field); err != nil {
			return err
		}
	}
	// 输出字段设置函数字段
	p.Write("\tfieldNameMaps map[string]string\n")
	p.Write("\tfields map[string]*pdl.TStructFieldInfo\n")

	return nil
}

func (p *TStructCodeWriter) writeStructConstructorMethod(ns *pdl.FileNamespace, stru *pdl.FileStruct) error {
	// 引用github.com/apache/thrift/lib/go/thrift
	p.Import("github.com/apache/thrift/lib/go/thrift", false)
	p.Import("github.com/go-xe2/xthrift/pdl", false)

	struName := p.GenStructNameCode(stru)
	p.Write("\n")
	p.Write(fmt.Sprintf("var _ pdl.DynamicStruct = (*%s)(nil)\n", struName))
	p.Write(fmt.Sprintf("var _ thrift.TStruct = (*%s)(nil)", struName))
	p.Write("\n")
	p.Write("\n")
	p.Write(fmt.Sprintf("func New%s() *%s {\n", struName, struName))
	p.Write(fmt.Sprintf("\tinst := &%s{\n", struName))
	p.Write("\t\tfieldNameMaps: make(map[string]string),\n")
	p.Write("\t\tfields: make(map[string]*pdl.TStructFieldInfo),\n")
	p.Write("\t}\n")

	p.Write("\tinst.TDynamicStructBase = pdl.NewBasicStruct(inst)\n")
	p.Write("\treturn inst.init()\n")
	p.Write("}\n")

	return nil
}

func (p *TStructCodeWriter) writeStructMethodDefine(stVar string, struName string, define string) {
	p.Write(fmt.Sprintf("func (%s *%s) %s {\n", stVar, struName, define))
}

func (p *TStructCodeWriter) writeStructMethodDefineEnd() {
	p.Write("}\n")
}

func (p *TStructCodeWriter) writeStructInitMethod(ns *pdl.FileNamespace, stVar string, struName string, stru *pdl.FileStruct) error {
	p.Write("\n")
	p.writeStructMethodDefine(stVar, struName, fmt.Sprintf("init() *%s", struName))

	size := len(p.fields)
	for i := 0; i < size; i++ {
		field := p.fields[i]
		fdName := p.GenFieldNameCode(field)

		p.Write(fmt.Sprintf("\t%s.fieldNameMaps[\"%s\"] = \"%s\"\n", stVar, fdName, fdName))
		p.Write(fmt.Sprintf("\t%s.fieldNameMaps[\"%s\"] = \"%s\"\n", stVar, p.GetIdentTransportName(fdName), fdName))

		p.Write("\n")

		switch field.FieldType.Type {
		case pdl.SPD_STR:
			p.Import("github.com/go-xe2/x/type/t", false)

			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\ts := t.String(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "s")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_BOOL:
			p.Import("github.com/go-xe2/x/type/t", false)
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\tb := t.Bool(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "b")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_I08:
			p.Import("github.com/go-xe2/x/type/t", false)
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\tn8 := t.Int8(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "n8")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_I16:
			p.Import("github.com/go-xe2/x/type/t", false)
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\tn16 := t.Int16(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "n16")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_I32:
			p.Import("github.com/go-xe2/x/type/t", false)
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\tn32 := t.Int32(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "n32")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_I64:
			p.Import("github.com/go-xe2/x/type/t", false)
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\tn64 := t.Int64(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "n64")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_DOUBLE:
			p.Import("github.com/go-xe2/x/type/t", false)
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write("\t\tf64 := t.Float64(val)\n")
			p.Write(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "thisObj", field, "f64")))
			p.Write("\n")
			p.Write("\t\treturn true\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_LIST:
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write(fmt.Sprintf("\t\tif lst, ok := val.(%s); ok {\n", p.GenDataTypeCode(ns, field.FieldType)))
			p.Write(fmt.Sprintf("\t\t\tthisObj.%s = lst\n", fdName))
			p.Write("\t\t\treturn true\n")
			p.Write("\t\t}\n")
			p.Write("\t\treturn false\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_SET:
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write(fmt.Sprintf("\t\tif set, ok := val.(%s); ok {\n", p.GenDataTypeCode(ns, field.FieldType)))
			p.Write(fmt.Sprintf("\t\t\tthisObj.%s = set\n", fdName))
			p.Write("\t\t\treturn true\n")
			p.Write("\t\t}\n")
			p.Write("\t\treturn false\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_MAP:
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write(fmt.Sprintf("\t\tif mp, ok := val.(%s); ok {\n", p.GenDataTypeCode(ns, field.FieldType)))
			p.Write(fmt.Sprintf("\t\t\tthisObj.%s = mp\n", fdName))
			p.Write("\t\t\treturn true\n")
			p.Write("\t\t}\n")
			p.Write("\t\treturn false\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_STRUCT:
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write(fmt.Sprintf("\t\tif stru, ok := val.(%s); ok {\n", p.GenDataTypeCode(ns, field.FieldType)))
			p.Write(fmt.Sprintf("\t\t\tthisObj.%s = stru\n", fdName))
			p.Write("\t\t\treturn true\n")
			p.Write("\t\t}\n")
			p.Write("\t\treturn false\n")
			p.Write("\t})\n")
			break
		case pdl.SPD_EXCEPTION:
			p.Write(fmt.Sprintf("\t%s.fields[\"%s\"] = pdl.NewStructFieldInfo(%d, %s, func(obj pdl.DynamicStruct, val interface{}) bool {\n", stVar, fdName, field.Id, comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
			p.Write(fmt.Sprintf("\t\tthisObj := obj.(*%s)\n", struName))
			p.Write(fmt.Sprintf("\t\tif stru, ok := val.(%s); ok {\n", p.GenDataTypeCode(ns, field.FieldType)))
			p.Write(fmt.Sprintf("\t\t\tthisObj.%s = stru\n", fdName))
			p.Write("\t\t\treturn true\n")
			p.Write("\t\t}\n")
			p.Write("\t\treturn false\n")
			p.Write("\t})\n")
			break
		}
		p.Write("\n")
	}
	p.Write(fmt.Sprintf("\treturn %s\n", stVar))
	p.writeStructMethodDefineEnd()
	return nil
}

func (p *TStructCodeWriter) WriteStructEnd(ns *pdl.FileNamespace, stru *pdl.FileStruct) error {
	p.Write("}\n")
	p.Write("\n")

	// 写构造方法
	if err := p.writeStructConstructorMethod(ns, stru); err != nil {
		return err
	}

	if err := p.writeStructInitMethod(ns, "p", p.GenStructNameCode(stru), stru); err != nil {
		return err
	}
	if err := p.writeStructReadMethod(ns, stru); err != nil {
		return err
	}
	if err := p.writeStructWriteMethod(ns, stru); err != nil {
		return err
	}
	if err := p.writerStructFieldGetterMethods(ns, stru); err != nil {
		return err
	}
	if err := p.writeDynamicStructImplement(ns, "p", stru); err != nil {
		return err
	}

	p.fields = nil
	p.ptrFields = nil
	return nil
}

func (p *TStructCodeWriter) writeStructReadMethod(ns *pdl.FileNamespace, stru *pdl.FileStruct) error {
	// 对字段进行排序
	p.Import("github.com/apache/thrift/lib/go/thrift", false)
	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) Read(in thrift.TProtocol) error {\n", p.GenStructNameCode(stru)))
	str := p.GenStructReadCode(ns, "in", "p", stru, p.fields)
	p.Write(comm.AppendIndent(1, str))
	p.Write("\n")
	p.Write("}\n")
	return nil
}

func (p *TStructCodeWriter) writeStructWriteMethod(ns *pdl.FileNamespace, stru *pdl.FileStruct) error {
	p.Import("github.com/apache/thrift/lib/go/thrift", false)
	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) Write(out thrift.TProtocol) error {\n", p.GenStructNameCode(stru)))
	str := p.GenStructWriteCode(ns, "out", "p", stru, p.fields)
	p.Write(comm.AppendIndent(1, str))
	p.Write("\n")
	p.Write("}\n")
	return nil
}

// 创建struct对象指针字段读取方法
func (p *TStructCodeWriter) writerStructFieldGetterMethods(ns *pdl.FileNamespace, stru *pdl.FileStruct) error {
	size := len(p.ptrFields)
	for i := 0; i < size; i++ {
		if err := p.writerStructFieldGetterMethod(ns, "p", stru, p.ptrFields[i]); err != nil {
			return err
		}
	}
	return nil
}

func (p *TStructCodeWriter) writerStructFieldGetterMethod(ns *pdl.FileNamespace, stVar string, stru *pdl.FileStruct, field *pdl.FileDataField) error {
	fdName := p.GenFieldNameCode(field)
	p.Write("\n")
	p.Write(fmt.Sprintf("// 字段%s读取方法,未设置时返回默认值\n", fdName))

	resultType := ""
	defReturn := ""
	resultField := ""
	switch field.FieldType.Type {
	case pdl.SPD_STR:
		resultType = "string"
		defReturn = fmt.Sprintf("s :=\"\"\n%s.%s = &s", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_BOOL:
		resultType = "bool"
		defReturn = fmt.Sprintf("b :=false\n%s.%s = &b", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_I08:
		resultType = "int8"
		defReturn = fmt.Sprintf("var n8 int8 = 0\n%s.%s = &n8", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_I16:
		resultType = "int16"
		defReturn = fmt.Sprintf("var n16 int16 = 0\n%s.%s = &n16", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_I32:
		resultType = "int32"
		defReturn = fmt.Sprintf("var n32 int32 = 0\n%s.%s = &n32", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_I64:
		resultType = "int64"
		defReturn = fmt.Sprintf("var n64 int64 = 0\n%s.%s = &n64", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_DOUBLE:
		resultType = "float64"
		defReturn = fmt.Sprintf("var f64 = 0\n%s.%s = &f64", stVar, fdName)
		resultField = fmt.Sprintf("*%s.%s", stVar, fdName)
		break
	case pdl.SPD_LIST:
		resultType = p.GenDataTypeCode(ns, field.FieldType)
		defReturn = fmt.Sprintf("%s.%s = make("+resultType+",0)", stVar, fdName)
		resultField = fmt.Sprintf("%s.%s", stVar, fdName)
		break
	case pdl.SPD_MAP:
		resultType = p.GenDataTypeCode(ns, field.FieldType)
		defReturn = fmt.Sprintf("%s.%s = make("+resultType+")", stVar, fdName)
		resultField = fmt.Sprintf("%s.%s", stVar, fdName)
		break
	case pdl.SPD_SET:
		resultType = p.GenDataTypeCode(ns, field.FieldType)
		defReturn = fmt.Sprintf("%s.%s = make("+resultType+",0)", stVar, fdName)
		resultField = fmt.Sprintf("%s.%s", stVar, fdName)
		break
	case pdl.SPD_STRUCT, pdl.SPD_EXCEPTION:
		resultType = p.GenDataTypeCode(ns, field.FieldType)
		defReturn = p.GenCreateStructTypeCode(ns, fmt.Sprintf("%s.%s", stVar, fdName), "=", field.FieldType)
		resultField = fmt.Sprintf("%s.%s", stVar, fdName)
		break
		//case pdl.SPD_EXCEPTION:
		//	resultType = p.GenDataTypeCode(ns, field.FieldType)
		//	defReturn = p.GenCreateExceptionTypeCode(ns, fmt.Sprintf("%s.%s", stVar, fdName), "", field.FieldType)
		//	resultField = fmt.Sprintf("%s.%s", stVar, fdName)
		//	break
	}

	p.Write(fmt.Sprintf("func (%s *%s) Get%s() %s {\n", stVar, p.GenStructNameCode(stru), fdName, resultType))

	// if p.field != nil { return p.field }
	p.Write(fmt.Sprintf("\tif %s.%s == nil {\n", stVar, fdName))
	p.Write(fmt.Sprintf("%s\n", comm.AppendIndent(2, defReturn)))
	p.Write("\t}\n")

	// return default value
	p.Write(fmt.Sprintf("\treturn %s\n", resultField))

	// end func
	p.Write("}\n")
	return nil
}

// 写实现DynamicStruct接口方法
func (p *TStructCodeWriter) writeDynamicStructImplement(ns *pdl.FileNamespace, stVar string, stru *pdl.FileStruct) error {
	p.Import("github.com/go-xe2/xthrift/pdl", false)

	struName := p.GenStructNameCode(stru)
	// 写NewInstance() DynamicStruct方法
	p.Write("\n")
	p.writeStructMethodDefine("p", struName, "NewInstance() pdl.DynamicStruct")
	p.Write(fmt.Sprintf("\treturn New%s()\n", struName))
	p.writeStructMethodDefineEnd()

	// 写AllFields() map[string]*TStructFieldInfo方法
	p.Write("\n")
	p.writeStructMethodDefine("p", struName, "AllFields() map[string]*pdl.TStructFieldInfo")
	p.Write(fmt.Sprintf("\treturn %s.fields\n", stVar))
	p.writeStructMethodDefineEnd()

	// 写FieldNameMaps() map[string]string
	p.Write("\n")
	p.writeStructMethodDefine("p", struName, "FieldNameMaps() map[string]string")
	p.Write(fmt.Sprintf("\treturn %s.fieldNameMaps\n", stVar))
	p.writeStructMethodDefineEnd()
	return nil
}
