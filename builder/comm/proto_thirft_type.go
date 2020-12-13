/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-13 15:35
* Description:
*****************************************************************/

package comm

import (
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/pdl"
	"strings"
)

var ProtoThriftTypes = map[pdl.TProtoBaseType]string{
	pdl.SPD_STR:       "thrift.STRING",
	pdl.SPD_BOOL:      "thrift.BOOL",
	pdl.SPD_I08:       "thrift.I08",
	pdl.SPD_I16:       "thrift.I16",
	pdl.SPD_I32:       "thrift.I32",
	pdl.SPD_I64:       "thrift.I64",
	pdl.SPD_DOUBLE:    "thrift.DOUBLE",
	pdl.SPD_LIST:      "thrift.LIST",
	pdl.SPD_SET:       "thrift.SET",
	pdl.SPD_MAP:       "thrift.MAP",
	pdl.SPD_STRUCT:    "thrift.STRUCT",
	pdl.SPD_EXCEPTION: "thrift.STRUCT",
}

func ProtoTypeToThriftDefine(typ pdl.TProtoBaseType) string {
	if s, ok := ProtoThriftTypes[typ]; ok {
		return s
	}
	return "thrift.VOID"
}

func AppendIndent(indent int, str string) string {
	if indent <= 0 {
		return str
	}
	szIndent := ""
	if indent > 0 {
		szIndent = xstring.MakeStrAndFill(indent, "\t", "")
	}
	items := strings.Split(str, "\n")
	result := make([]string, 0)
	for i := 0; i < len(items); i++ {
		if items[i] != "" {
			result = append(result, szIndent+items[i]+"\n")
		}
	}
	return strings.Join(result, "")
}
