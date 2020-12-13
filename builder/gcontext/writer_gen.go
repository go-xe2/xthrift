/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 15:04
* Description:
*****************************************************************/

package gcontext

import (
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TWriter) GenDataTypeCode(ns *pdl.FileNamespace, dataType *pdl.FileDataType) string {
	// 前掇
	prefix := ""
	// 类型名称
	szName := ""
	switch dataType.Type {
	case pdl.SPD_STR:
		szName = "string"
		break
	case pdl.SPD_BOOL:
		szName = "bool"
		break
	case pdl.SPD_I08:
		szName = "int8"
		break
	case pdl.SPD_I16:
		szName = "int16"
		break
	case pdl.SPD_I32:
		szName = "int32"
		break
	case pdl.SPD_I64:
		szName = "int64"
		break
	case pdl.SPD_DOUBLE:
		szName = "float64"
		break
	case pdl.SPD_LIST:
		elemType := p.GenDataTypeCode(ns, dataType.ElemType)
		szName = "[]" + elemType
		break
	case pdl.SPD_SET:
		elemType := p.GenDataTypeCode(ns, dataType.ElemType)
		szName = "[]" + elemType
		break
	case pdl.SPD_MAP:
		keyType := p.GenDataTypeCode(ns, dataType.KeyType)
		valType := p.GenDataTypeCode(ns, dataType.ValType)
		szName = "map[" + keyType + "]" + valType
		break
	case pdl.SPD_STRUCT, pdl.SPD_EXCEPTION:
		if dataType.Namespace != "" && dataType.Namespace != ns.Namespace {
			_, prefix = pdl.NamespaceLastName(dataType.Namespace)
			szName = fmt.Sprintf("*%s.%s", prefix, dataType.TypName)
		} else {
			szName = "*" + dataType.TypName
		}
		break
	default:
		szName = dataType.TypName
	}
	return szName
}

func (p *TWriter) GenFieldDefineTypeCode(ns *pdl.FileNamespace, field *pdl.FileDataField) string {
	szName := p.GenDataTypeCode(ns, field.FieldType)
	if field.Limit == pdl.SPDLimitOptional && !(field.FieldType.Type == pdl.SPD_MAP || field.FieldType.Type == pdl.SPD_LIST ||
		field.FieldType.Type == pdl.SPD_SET ||
		field.FieldType.Type == pdl.SPD_STRUCT || field.FieldType.Type == pdl.SPD_EXCEPTION) {
		return "*" + szName
	}
	return szName
}

// 生成创建struct类型实例代码
func (p *TWriter) GenCreateStructTypeCode(ns *pdl.FileNamespace, varName string, assign string, stru *pdl.FileDataType) string {
	prefix := ""
	szName := xstring.UcFirst(stru.TypName)
	if ns == nil || (stru.Namespace != ns.Namespace) {
		_, prefix = pdl.NamespaceLastName(stru.Namespace)
	}
	if prefix != "" {
		szName = prefix + ".New" + szName
	} else {
		szName = "New" + szName
	}
	return fmt.Sprintf("%s %s %s()", varName, assign, szName)
}

// 生成创建exception类型实例代码
func (p *TWriter) GenCreateExceptionTypeCode(ns *pdl.FileNamespace, varName string, assign string, except *pdl.FileDataType) string {
	prefix := ""
	szName := xstring.UcFirst(except.TypName)
	if ns == nil || (except.Namespace != ns.Namespace) {
		_, prefix = pdl.NamespaceLastName(except.Namespace)
	}
	if prefix != "" {
		szName = prefix + ".New" + szName
	} else {
		szName = "New" + szName
	}
	return fmt.Sprintf("%s %s %s()", varName, assign, szName)
}

// 生成字段赋值代码
func (p *TWriter) GenFieldAssignValueCode(ns *pdl.FileNamespace, struName string, field *pdl.FileDataField, value string) string {
	fdName := struName + "." + xstring.UcFirst(field.Name)
	if field.Limit == pdl.SPDLimitOptional && !(field.FieldType.Type == pdl.SPD_STRUCT || field.FieldType.Type == pdl.SPD_EXCEPTION ||
		field.FieldType.Type == pdl.SPD_MAP || field.FieldType.Type == pdl.SPD_LIST || field.FieldType.Type == pdl.SPD_SET) {
		value = "&" + value
	}
	return fmt.Sprintf("%s = %s", fdName, value)
}

func (p *TWriter) GenFieldAssignToCode(ns *pdl.FileNamespace, struName string, field *pdl.FileDataField) string {
	fdName := struName + "." + xstring.UcFirst(field.Name)
	if field.Limit == pdl.SPDLimitOptional && !(field.FieldType.Type == pdl.SPD_EXCEPTION || field.FieldType.Type == pdl.SPD_STRUCT ||
		field.FieldType.Type == pdl.SPD_MAP || field.FieldType.Type == pdl.SPD_LIST || field.FieldType.Type == pdl.SPD_SET) {
		fdName = "*" + fdName
	}
	return fdName
}

func (p *TWriter) GenFieldNameCode(field *pdl.FileDataField) string {
	return xstring.UcFirst(field.Name)
}

func (p *TWriter) GenStructNameCode(stru *pdl.FileStruct) string {
	return xstring.UcFirst(stru.Type.TypName)
}

func (p *TWriter) GenExceptNameCode(expt *pdl.FileStruct) string {
	return xstring.UcFirst(expt.Type.TypName)
}

func (p *TWriter) GenServiceNameCode(svc *pdl.FileService) string {
	return xstring.UcFirst(svc.Name)
}

func (p *TWriter) GenServiceMethodNameCode(method *pdl.FileServiceMethod) string {
	return xstring.UcFirst(method.Name)
}
