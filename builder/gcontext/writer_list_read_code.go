/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 14:32
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenListReadCode(ns *pdl.FileNamespace, protoVar string, lstName string, lst *pdl.FileDataType) string {
	if lst.Type != pdl.SPD_LIST {
		return ""
	}
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("elemType, size, err := %s.ReadListBegin()\n", protoVar))
	result.WriteString("if err != nil {\n")
	result.WriteString("	\treturn err\n")
	result.WriteString("}\n")
	result.WriteString(fmt.Sprintf("if elemType != %s {\n", comm.ProtoTypeToThriftDefine(lst.ElemType.Type)))
	result.WriteString("\treturn thrift.NewTApplicationException(thrift.INVALID_PROTOCOL, \"协议数据类型不匹配\")\n")
	result.WriteString("}\n")
	result.WriteString(fmt.Sprintf("%s := make(%s, size)\n", lstName, p.GenDataTypeCode(ns, lst)))
	result.WriteString("for j := 0; j < size; j++ {\n")
	switch lst.ElemType.Type {
	case pdl.SPD_STR:
		s := p.GenStrReadCode(protoVar, "s", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = s\n", lstName))
		break
	case pdl.SPD_BOOL:
		s := p.GenBoolReadCode(protoVar, "b", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = b\n", lstName))
		break
	case pdl.SPD_I08:
		s := p.GenI08ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n\n", lstName))
		break
	case pdl.SPD_I16:
		s := p.GenI16ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n[i]\n", lstName))
		break
	case pdl.SPD_I32:
		s := p.GenI32ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n\n", lstName))
		break
	case pdl.SPD_I64:
		s := p.GenI64ReadCode(protoVar, "n", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = n\n", lstName))
		break
	case pdl.SPD_DOUBLE:
		s := p.GenDoubleReadCode(protoVar, "d", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		result.WriteString(fmt.Sprintf("\t%s[j] = d\n", lstName))
		break
	case pdl.SPD_LIST:
		subStr := p.GenListReadCode(ns, protoVar, "subItem", lst.ElemType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		result.WriteString(fmt.Sprintf("\t%s[j] = subItem\n", lstName))
		break
	case pdl.SPD_SET:
		subStr := p.GenSetReadCode(ns, protoVar, "subSet", lst.ElemType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		result.WriteString(fmt.Sprintf("\t%s[j] = subSet\n", lstName))
		break
	case pdl.SPD_MAP:
		subStr := p.GenMapReadCode(ns, protoVar, "subMp", lst.ElemType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		result.WriteString(fmt.Sprintf("\t%s[j] = subMp\n", lstName))
		break
	case pdl.SPD_STRUCT:
		result.WriteString(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "rec", ":=", lst.ElemType)))
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\tif err := rec.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		result.WriteString(fmt.Sprintf("\t%s[j] = rec\n", lstName))
		break
	case pdl.SPD_EXCEPTION:
		result.WriteString(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "rec", ":=", lst.ElemType)))
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\tif err := rec.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		result.WriteString(fmt.Sprintf("\t%s[j] = rec\n", lstName))
		break
	}
	result.WriteString("}\n")
	result.WriteString(fmt.Sprintf("if err := %s.ReadListEnd(); err != nil {\n", protoVar))
	result.WriteString("		return err\n")
	result.WriteString("}\n")
	return result.String()
}
