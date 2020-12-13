/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 14:35
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/xthrift/builder/comm"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenMapReadCode(ns *pdl.FileNamespace, protoVar string, mpName string, mp *pdl.FileDataType) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("keyType, valType, size, err := %s.ReadMapBegin()\n", protoVar))
	result.WriteString("if err != nil {\n")
	result.WriteString("		return err\n")
	result.WriteString("}\n")
	result.WriteString(fmt.Sprintf("if keyType != %s || valType != %s {\n", comm.ProtoTypeToThriftDefine(mp.KeyType.Type), comm.ProtoTypeToThriftDefine(mp.ValType.Type)))
	result.WriteString("\treturn thrift.NewTApplicationException(thrift.INVALID_PROTOCOL, \"协议数据类型不匹配\")\n")
	result.WriteString("}\n")

	result.WriteString(fmt.Sprintf("%s := make(%s)\n", mpName, p.GenDataTypeCode(ns, mp)))
	result.WriteString("for i := 0; i < size; i++ {\n")

	// 读取键名
	switch mp.KeyType.Type {
	case pdl.SPD_STR:
		s := p.GenStrReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_BOOL:
		s := p.GenBoolReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I08:
		s := p.GenI08ReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I16:
		s := p.GenI16ReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I32:
		s := p.GenI32ReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I64:
		s := p.GenI64ReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_DOUBLE:
		s := p.GenDoubleReadCode(protoVar, "key", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	default:
		result.WriteString("		break")
	}

	// 读取值
	switch mp.ValType.Type {
	case pdl.SPD_STR:
		s := p.GenStrReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_BOOL:
		s := p.GenBoolReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I08:
		s := p.GenI08ReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I16:
		s := p.GenI16ReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I32:
		s := p.GenI32ReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_I64:
		s := p.GenI64ReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_DOUBLE:
		s := p.GenDoubleReadCode(protoVar, "val", "err")
		s = comm.AppendIndent(1, s)
		result.WriteString(s)
		break
	case pdl.SPD_LIST:
		subStr := p.GenListReadCode(ns, protoVar, "val", mp.ValType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		break
	case pdl.SPD_SET:
		subStr := p.GenSetReadCode(ns, protoVar, "val", mp.ValType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		break
	case pdl.SPD_MAP:
		subStr := p.GenMapReadCode(ns, protoVar, "val", mp.ValType)
		subStr = comm.AppendIndent(1, subStr)
		result.WriteString(subStr)
		break
	case pdl.SPD_STRUCT:
		result.WriteString(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "val", ":=", mp.ValType)))
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\tif err := val.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		break
	case pdl.SPD_EXCEPTION:
		result.WriteString(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "val", ":=", mp.ValType)))
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\tif err := val.Read(%s); err != nil {\n", protoVar))
		result.WriteString("\t\treturn err\n")
		result.WriteString("\t}\n")
		break
	}
	// 加入到map
	result.WriteString(fmt.Sprintf("\t%s[key] = val\n", mpName))

	result.WriteString("}\n")

	result.WriteString(fmt.Sprintf("if err := %s.ReadMapEnd(); err != nil {\n", protoVar))
	result.WriteString("		return err\n")
	result.WriteString("}\n")
	return result.String()
}
