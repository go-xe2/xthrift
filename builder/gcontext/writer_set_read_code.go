/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 14:33
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenSetReadCode(ns *pdl.FileNamespace, protoVar string, setName string, set *pdl.FileDataType) string {
	if set.Type != pdl.SPD_SET {
		return ""
	}
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("elemType, size, err := %s.ReadSetBegin()\n", protoVar))
	result.WriteString("if err != nil {\n")
	result.WriteString("	\treturn err\n")
	result.WriteString("}\n")
	result.WriteString(fmt.Sprintf("if elemType != %s {\n", comm.ProtoTypeToThriftDefine(set.ElemType.Type)))
	result.WriteString("\treturn thrift.NewTApplicationException(thrift.INVALID_PROTOCOL, \"协议数据类型不匹配\")\n")
	result.WriteString("}\n")
	result.WriteString(fmt.Sprintf("%s := make(%s, size)\n", setName, p.GenDataTypeCode(ns, set)))
	result.WriteString("for j := 0; j < size; j++ {\n")
	switch set.ElemType.Type {
	case pdl.SPD_STR:
		s := p.GenStrReadCode(protoVar, "s", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = s\n", setName))
		break
	case pdl.SPD_BOOL:
		s := p.GenBoolReadCode(protoVar, "b", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = b\n", setName))
		break
	case pdl.SPD_I08:
		s := p.GenI08ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n\n", setName))
		break
	case pdl.SPD_I16:
		s := p.GenI16ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n[i]\n", setName))
		break
	case pdl.SPD_I32:
		s := p.GenI32ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n\n", setName))
		break
	case pdl.SPD_I64:
		s := p.GenI64ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n\n", setName))
		break
	case pdl.SPD_DOUBLE:
		s := p.GenDoubleReadCode(protoVar, "d", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = d\n", setName))
		break
	case pdl.SPD_LIST:
		subStr := p.GenListReadCode(ns, protoVar, "subItem", set.ElemType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		result.WriteString(fmt.Sprintf("\t%s[j] = subItem\n", setName))
		break
	case pdl.SPD_SET:
		subStr := p.GenSetReadCode(ns, protoVar, "subSet", set.ElemType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		result.WriteString(fmt.Sprintf("\t%s[j] = subSet\n", setName))
		break
	case pdl.SPD_MAP:
		subStr := p.GenMapReadCode(ns, protoVar, "subMp", set.ElemType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		result.WriteString(fmt.Sprintf("\t%s[j] = subMp\n", setName))
		break
	case pdl.SPD_STRUCT:
		result.WriteString(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "rec", ":=", set.ElemType)))
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\tif err := rec.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		result.WriteString(fmt.Sprintf("\t%s[j] = rec\n", setName))
		break
	case pdl.SPD_EXCEPTION:
		result.WriteString(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "rec", ":=", set.ElemType)))
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\tif err := rec.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		result.WriteString(fmt.Sprintf("\t%s[j] = rec\n", setName))
		break
	}
	result.WriteString("}\n")
	result.WriteString("if err := %s.ReadSetEnd(); err != nil {\n")
	result.WriteString("		return err\n")
	result.WriteString("}\n")
	return result.String()
}
