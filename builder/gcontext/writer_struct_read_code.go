/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 14:12
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	. "github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenStructReadCode(ns *FileNamespace, protoVar string, stName string, stru *FileStruct, fields []*FileDataField) string {
	result := bytes.NewBufferString("")

	size := len(fields)

	if size == 0 {
		result.WriteString(fmt.Sprintf("if err := %s.Skip(thrift.STRUCT); err != nil {\n", protoVar))
		result.WriteString("\treturn err\n")
		result.WriteString("}\n")
		result.WriteString("return nil\n")
		return result.String()
	}

	result.WriteString(fmt.Sprintf("_, err := %s.ReadStructBegin()\n", protoVar))
	result.WriteString("if err != nil {\n")
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	result.WriteString("var nMaxLoop = 512\n")
	result.WriteString("nLoop := 0\n")
	result.WriteString("var isMatch bool\n")
	result.WriteString("for {\n")
	result.WriteString("	\t// 防止协议数据错误，无thrift.STOP时无限循环\n")
	result.WriteString("\tnLoop++\n")
	result.WriteString("\tif nLoop >= nMaxLoop {\n")
	result.WriteString(fmt.Sprintf("\t\t_ = %s.Skip(thrift.STRUCT)\n", protoVar))
	result.WriteString("\t\treturn nil\n")
	result.WriteString("\t}\n")

	result.WriteString("\tisMatch = false\n")

	result.WriteString(fmt.Sprintf("\tfdName, fdType, fdId, err := %s.ReadFieldBegin()\n", protoVar))
	result.WriteString("\tif err != nil {\n")
	result.WriteString("\t\treturn err\n")
	result.WriteString("\t}\n")

	result.WriteString("\tif fdType == thrift.STOP {\n")
	result.WriteString("\t\tbreak\n")
	result.WriteString("\t}\n")

	result.WriteString("\tif fdType == thrift.VOID {\n")

	result.WriteString(fmt.Sprintf("\t\tif err := %s.ReadFieldEnd(); err != nil {\n", protoVar))
	result.WriteString("\t\t\treturn err\n")
	result.WriteString("\t\t}\n")

	result.WriteString("\t\tcontinue\n")
	result.WriteString("\t}\n")

	// 字段读取数据
	for i := 0; i < size; i++ {
		field := fields[i]
		fdName := p.GenFieldNameCode(field)
		inOutName := p.GetIdentTransportName(fdName)

		// if (fdId > 0 && fdId == %d) || (fdId <= 0 && fdName ==field.name) {
		result.WriteString(fmt.Sprintf("\tif (fdId > 0 && fdId == %d) || (fdId <= 0 && fdName ==\"%s\") {\n", field.Id, inOutName))
		//	if fdId > 0 && fdType != field.Type {
		// 字段与期望字段类型不一致时，忽略字段内容，读取下一个字段值
		result.WriteString(fmt.Sprintf("\t\tif fdId > 0 && fdType != %s {\n", comm.ProtoTypeToThriftDefine(field.FieldType.Type)))
		result.WriteString(fmt.Sprintf("\t\t\tif err := %s.Skip(fdType); err !=nil {\n", protoVar))
		result.WriteString("\t\t\t\treturn err\n")
		result.WriteString("\t\t\t}\n")

		result.WriteString(fmt.Sprintf("\t\t\tif err := %s.ReadFieldEnd(); err != nil {\n", protoVar))
		result.WriteString("\t\t\t\treturn err\n")
		result.WriteString("\t\t\t}\n")
		// 	continue
		result.WriteString("\t\t\tcontinue\n")
		//	}
		result.WriteString("\t\t}\n")

		result.WriteString("\t\tisMatch = true\n")

		switch field.FieldType.Type {
		case SPD_STR:
			s := p.GenStrReadCode(protoVar, "s", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "s")))
			result.WriteString("\n")
			break
		case SPD_BOOL:
			s := p.GenBoolReadCode(protoVar, "b", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "b")))
			result.WriteString("\n")
			break
		case SPD_I08:
			s := p.GenI08ReadCode(protoVar, "n", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "n")))
			result.WriteString("\n")
			break
		case SPD_I16:
			s := p.GenI16ReadCode(protoVar, "n", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "n")))
			result.WriteString("\n")
			break
		case SPD_I32:
			s := p.GenI32ReadCode(protoVar, "n", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "n")))
			result.WriteString("\n")
			break
		case SPD_I64:
			s := p.GenI64ReadCode(protoVar, "n", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "n")))
			result.WriteString("\n")
			break
		case SPD_DOUBLE:
			s := p.GenDoubleReadCode(protoVar, "d", "err")
			s = comm.AppendIndent(2, s)
			result.WriteString(s)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "d")))
			result.WriteString("\n")
			break
		case SPD_LIST:
			subStr := p.GenListReadCode(ns, protoVar, "lst", field.FieldType)
			subStr = comm.AppendIndent(2, subStr)
			result.WriteString(subStr)
			result.WriteString("\n")
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "lst")))
			result.WriteString("\n")
			break
		case SPD_SET:
			subStr := p.GenSetReadCode(ns, protoVar, "st", field.FieldType)
			subStr = comm.AppendIndent(2, subStr)
			result.WriteString(subStr)
			result.WriteString(comm.AppendIndent(2, p.GenFieldAssignValueCode(ns, "p", field, "st")))
			result.WriteString("\n")
			break
		case SPD_MAP:
			subStr := p.GenMapReadCode(ns, protoVar, "mp", field.FieldType)
			subStr = comm.AppendIndent(2, subStr)
			result.WriteString(subStr)
			result.WriteString("\n")
			result.WriteString(fmt.Sprintf("\t\t%s.%s = mp\n", stName, fdName))
			break
		case SPD_STRUCT:
			s := p.GenCreateStructTypeCode(ns, fmt.Sprintf("%s.%s ", stName, fdName), "=", field.FieldType)
			result.WriteString(comm.AppendIndent(2, s))
			result.WriteString("\n")
			result.WriteString(fmt.Sprintf("\t\tif err := %s.%s.Read(%s); err != nil {\n", stName, fdName, protoVar))
			result.WriteString("\t\t\treturn err\n")
			result.WriteString("\t\t}\n")
			break
		case SPD_EXCEPTION:
			s := p.GenCreateExceptionTypeCode(ns, fmt.Sprintf("%s.%s ", stName, fdName), "=", field.FieldType)
			result.WriteString(comm.AppendIndent(2, s))
			result.WriteString("\n")
			result.WriteString(fmt.Sprintf("\t\tif err := %s.%s.Read(%s); err != nil {\n", stName, fdName, protoVar))
			result.WriteString("\t\t\treturn err\n")
			result.WriteString("\t\t}\n")
			break
		}

		// end if
		result.WriteString("\t}\n")
	}

	result.WriteString("\tif !isMatch {\n")
	result.WriteString(fmt.Sprintf("\t\tif err := %s.Skip(fdType); err != nil {\n", protoVar))
	result.WriteString("\t\t\treturn err\n")
	result.WriteString("\t\t}\n")
	result.WriteString("\t}\n")

	result.WriteString(fmt.Sprintf("\tif err := %s.ReadFieldEnd(); err != nil {\n", protoVar))
	result.WriteString("\t\treturn err\n")
	result.WriteString("\t}\n")

	// end for
	result.WriteString("}\n") // end for

	result.WriteString(fmt.Sprintf("if err := %s.ReadStructEnd(); err != nil {\n", protoVar))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")
	result.WriteString("return nil\n")
	return result.String()
}
