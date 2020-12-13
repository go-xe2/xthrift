/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 16:16
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	. "github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenListWriteCode(ns *FileNamespace, protoVar string, lstName string, lst *FileDataType) string {
	if lst.Type != SPD_LIST {
		return ""
	}

	result := bytes.NewBufferString("")

	result.WriteString(fmt.Sprintf("lstSize := len(%s)\n", lstName))
	result.WriteString(fmt.Sprintf("if err := %s.WriteListBegin(%s, lstSize); err != nil {\n", protoVar, comm.ProtoTypeToThriftDefine(lst.ElemType.Type)))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	// for
	result.WriteString("for i := 0; i < lstSize; i++ {\n")

	switch lst.ElemType.Type {
	case SPD_STR:
		s := p.GenStrWriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_BOOL:
		s := p.GenBoolWriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I08:
		s := p.GenI08WriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I16:
		s := p.GenI16WriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I32:
		s := p.GenI32WriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I64:
		s := p.GenI64WriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_DOUBLE:
		s := p.GenDoubleWriteCode(protoVar, fmt.Sprintf("%s[i]", lstName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_LIST:
		subStr := p.GenListWriteCode(ns, protoVar, fmt.Sprintf("%s[i]", lstName), lst.ElemType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_SET:
		subStr := p.GenSetWriteCode(ns, protoVar, fmt.Sprintf("%s[i]", lstName), lst.ElemType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_MAP:
		subStr := p.GenMapWriteCode(ns, protoVar, fmt.Sprintf("%s[i]", lstName), lst.ElemType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_STRUCT:
		result.WriteString(fmt.Sprintf("\tif err := %s[i].Write(%s); err != nil {\n", lstName, protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		break
	case SPD_EXCEPTION:
		result.WriteString(fmt.Sprintf("\tif err := %s[i].Write(%s); err != nil {\n", lstName, protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		break
	}

	// end for
	result.WriteString("}\n")

	// WriteListEnd
	result.WriteString(fmt.Sprintf("if err := %s.WriteListEnd(); err != nil {\n", protoVar))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	return result.String()
}
