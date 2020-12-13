/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 15:47
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	. "github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenStructWriteCode(ns *FileNamespace, protoVar string, stName string, stru *FileStruct, fields []*FileDataField) string {
	size := len(fields)

	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if err := %s.WriteStructBegin(\"%s\"); err != nil {\n", protoVar, p.GetIdentTransportName(stru.Type.TypName)))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	for i := 0; i < size; i++ {
		field := fields[i]
		allowNil := field.Limit == SPDLimitOptional || field.FieldType.Type == SPD_STRUCT || field.FieldType.Type == SPD_EXCEPTION || field.FieldType.Type == SPD_LIST ||
			field.FieldType.Type == SPD_SET

		indent := 0
		if allowNil {
			indent = 1

			result.WriteString(fmt.Sprintf("if %s.%s != nil {\n", stName, p.GenFieldNameCode(field)))

			result.WriteString(fmt.Sprintf("\tif err := %s.WriteFieldBegin(\"%s\", %s, %d); err != nil {\n", protoVar, p.GetIdentTransportName(field.Name), comm.ProtoTypeToThriftDefine(field.FieldType.Type), field.Id))
			result.WriteString("\t\treturn err\n")
			result.WriteString("\t}\n")
		} else {
			result.WriteString(fmt.Sprintf("if err := %s.WriteFieldBegin(\"%s\", %s, %d); err != nil {\n", protoVar, p.GetIdentTransportName(field.Name), comm.ProtoTypeToThriftDefine(field.FieldType.Type), field.Id))
			result.WriteString("\treturn err\n")
			result.WriteString("}\n")
		}

		switch field.FieldType.Type {
		case SPD_STR:
			s := p.GenStrWriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_BOOL:
			s := p.GenBoolWriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_I08:
			s := p.GenI08WriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_I16:
			s := p.GenI16WriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_I32:
			s := p.GenI32WriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_I64:
			s := p.GenI64WriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_DOUBLE:
			s := p.GenDoubleWriteCode(protoVar, p.GenFieldAssignToCode(ns, stName, field), "err")
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")
			break
		case SPD_LIST:

			s := p.GenListWriteCode(ns, protoVar, p.GenFieldAssignToCode(ns, stName, field), field.FieldType)
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")

			break
		case SPD_SET:

			s := p.GenSetWriteCode(ns, protoVar, p.GenFieldAssignToCode(ns, stName, field), field.FieldType)
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")

			break
		case SPD_MAP:

			s := p.GenMapWriteCode(ns, protoVar, fmt.Sprintf("%s.%s", stName, p.GenFieldNameCode(field)), field.FieldType)
			result.WriteString(comm.AppendIndent(indent, s))
			result.WriteString("\n")

			break
		case SPD_STRUCT:

			result.WriteString(comm.AppendIndent(indent, fmt.Sprintf("if err := %s.%s.Write(%s); err != nil {\n", stName, p.GenFieldNameCode(field), protoVar)))
			result.WriteString(comm.AppendIndent(indent, "\treturn err\n"))
			result.WriteString(comm.AppendIndent(indent, "}\n"))

			break
		case SPD_EXCEPTION:

			result.WriteString(comm.AppendIndent(indent, fmt.Sprintf("if err := %s.%s.Write(%s); err != nil {\n", stName, p.GenFieldNameCode(field), protoVar)))
			result.WriteString(comm.AppendIndent(indent, "\treturn err\n"))
			result.WriteString(comm.AppendIndent(indent, "}\n"))

			break
		}

		if allowNil {
			result.WriteString(fmt.Sprintf("\tif err := %s.WriteFieldEnd(); err != nil {\n", protoVar))
			result.WriteString("\t\treturn err\n")
			result.WriteString("\t}\n")

			result.WriteString("}\n")

		} else {
			result.WriteString(fmt.Sprintf("if err := %s.WriteFieldEnd(); err != nil {\n", protoVar))
			result.WriteString("\treturn err\n")
			result.WriteString("}\n")
		}

	}

	result.WriteString(fmt.Sprintf("if err := %s.WriteFieldStop(); err != nil {\n", protoVar))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	result.WriteString(fmt.Sprintf("if err := %s.WriteStructEnd(); err != nil {\n", protoVar))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	result.WriteString("return nil\n")

	return result.String()
}
