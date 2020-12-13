/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 16:33
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	. "github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenMapWriteCode(ns *FileNamespace, protoVar string, mpName string, mp *FileDataType) string {
	if mp.Type != SPD_MAP {
		return ""
	}

	result := bytes.NewBufferString("")

	result.WriteString(fmt.Sprintf("mpSize := len(%s)\n", mpName))
	result.WriteString(fmt.Sprintf("if err := %s.WriteMapBegin(%s, %s, mpSize); err != nil {\n", protoVar, comm.ProtoTypeToThriftDefine(mp.KeyType.Type),
		comm.ProtoTypeToThriftDefine(mp.ValType.Type)))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	// for
	result.WriteString(fmt.Sprintf("for k, v := range %s {\n", mpName))

	// 写入键名
	switch mp.KeyType.Type {
	case SPD_STR:
		s := p.GenStrWriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_BOOL:
		s := p.GenBoolWriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I08:
		s := p.GenI08WriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I16:
		s := p.GenI16WriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I32:
		s := p.GenI32WriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I64:
		s := p.GenI64WriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_DOUBLE:
		s := p.GenDoubleWriteCode(protoVar, "k", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	}

	// 写入键值
	switch mp.ValType.Type {
	case SPD_STR:
		s := p.GenStrWriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_BOOL:
		s := p.GenBoolWriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I08:
		s := p.GenI08WriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I16:
		s := p.GenI16WriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I32:
		s := p.GenI32WriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_I64:
		s := p.GenI64WriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_DOUBLE:
		s := p.GenDoubleWriteCode(protoVar, "v", "err")
		result.WriteString(comm.AppendIndent(1, s))
		break
	case SPD_LIST:
		subStr := p.GenListWriteCode(ns, protoVar, "v", mp.ValType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_SET:
		subStr := p.GenListWriteCode(ns, protoVar, "v", mp.ValType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_MAP:
		subStr := p.GenMapWriteCode(ns, protoVar, "v", mp.ValType)
		result.WriteString(comm.AppendIndent(1, subStr))
		result.WriteString("\n")
		break
	case SPD_STRUCT:
		result.WriteString(fmt.Sprintf("if err := v.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\treturn err\n")
		result.WriteString("}\n")
		break
	case SPD_EXCEPTION:
		result.WriteString(fmt.Sprintf("if err := v.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\treturn err\n")
		result.WriteString("}\n")
		break
	}

	// end for
	result.WriteString("}\n")

	// WriteListEnd
	result.WriteString(fmt.Sprintf("if err := %s.WriteMapEnd(); err != nil {\n", protoVar))
	result.WriteString("\treturn err\n")
	result.WriteString("}\n")

	return result.String()
}
