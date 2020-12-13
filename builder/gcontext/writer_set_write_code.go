/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 16:29
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	. "github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenSetWriteCode(ns *FileNamespace, protoVar string, setName string, set *FileDataType) string {
	if set.Type != SPD_SET {
		return ""
	}
	result := bytes.NewBufferString("")

	result.WriteString(fmt.Sprintf("lstSize := len(%s)\n", setName))
	result.WriteString(fmt.Sprintf("if err := %s.WriteSetBegin(%s, lstSize); err != nil {\n", protoVar, comm.ProtoTypeToThriftDefine(set.ElemType.Type)))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	// for
	result.WriteString("for i := 0; i < lstSize; i++ {\n")

	switch set.ElemType.Type {
	case SPD_STR:
		s := p.GenStrWriteCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_BOOL:
		s := p.GenBoolReadCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I08:
		s := p.GenI08WriteCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I16:
		s := p.GenI16WriteCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I32:
		s := p.GenI32WriteCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I64:
		s := p.GenI64WriteCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_DOUBLE:
		s := p.GenDoubleWriteCode(protoVar, fmt.Sprintf("%s[i]", setName), "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_LIST:
		subStr := p.GenListWriteCode(ns, protoVar, fmt.Sprintf("%s[i]", setName), set.ElemType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_SET:
		subStr := p.GenSetWriteCode(ns, protoVar, fmt.Sprintf("%s[i]", setName), set.ElemType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_MAP:
		subStr := p.GenMapWriteCode(ns, protoVar, fmt.Sprintf("%s[i]", setName), set.ElemType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_STRUCT:
		result.WriteString(fmt.Sprintf("if err := %s[i].Write(%s); err != nil {\n", setName, protoVar))
		result.WriteString("\treturn err\n")
		result.WriteString("}\n")
		break
	case SPD_EXCEPTION:
		result.WriteString(fmt.Sprintf("if err := %s[i].Write(%s); err != nil {\n", setName, protoVar))
		result.WriteString("\treturn err\n")
		result.WriteString("}\n")
		break
	}

	// end for
	result.WriteString("}\n")

	// WriteListEnd
	result.WriteString(fmt.Sprintf("if err := %s.WriteSetEnd(); err != nil {\n", protoVar))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	return result.String()
}
