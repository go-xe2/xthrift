/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-07 15:30
* Description:
*****************************************************************/

package pdl

import (
	"encoding/json"
	"github.com/apache/thrift/lib/go/thrift"
)

type TProtoBaseType int8

var ProtoBasicTypes = []string{
	"void",
	"str",
	"bl",
	"i8",
	"i16",
	"i32",
	"i64",
	"idl",
	"list",
	"set",
	"map",
	"struct",
	"exception",
}

const (
	SPD_VOID TProtoBaseType = iota
	SPD_STR
	SPD_BOOL
	SPD_I08
	SPD_I16
	SPD_I32
	SPD_I64
	SPD_DOUBLE
	SPD_LIST
	SPD_SET
	SPD_MAP
	SPD_STRUCT
	SPD_EXCEPTION
	SPD_TYPEDEF
	SPD_UNKNOWN
)

func (spd TProtoBaseType) String() string {
	switch spd {
	case SPD_VOID:
		return "void"
	case SPD_STR:
		return "str"
	case SPD_BOOL:
		return "bl"
	case SPD_I08:
		return "i8"
	case SPD_I16:
		return "i16"
	case SPD_I32:
		return "i32"
	case SPD_I64:
		return "i64"
	case SPD_DOUBLE:
		return "dl"
	case SPD_LIST:
		return "list"
	case SPD_SET:
		return "set"
	case SPD_MAP:
		return "map"
	case SPD_STRUCT:
		return "struct"
	case SPD_EXCEPTION:
		return "exception"
	case SPD_TYPEDEF:
		return "typedef"
	case SPD_UNKNOWN:
		return "unknown"
	}
	return "unknown"
}

func (spd *TProtoBaseType) MarshalJSON() (data []byte, err error) {
	s := spd.String()
	return json.Marshal(s)
}

func (spd TProtoBaseType) ThriftType() thrift.TType {
	switch spd {
	case SPD_STR:
		return thrift.STRING
	case SPD_BOOL:
		return thrift.BOOL
	case SPD_I08:
		return thrift.BYTE
	case SPD_I16:
		return thrift.I16
	case SPD_I32:
		return thrift.I32
	case SPD_I64:
		return thrift.I64
	case SPD_DOUBLE:
		return thrift.DOUBLE
	case SPD_LIST:
		return thrift.LIST
	case SPD_SET:
		return thrift.SET
	case SPD_MAP:
		return thrift.MAP
	case SPD_STRUCT:
		return thrift.STRUCT
	}
	return thrift.VOID
}

type ProtoFieldLimit int8

const (
	SPDLimitRequired ProtoFieldLimit = iota
	SPDLimitOptional
)

func (spdf ProtoFieldLimit) String() string {
	switch spdf {
	case SPDLimitRequired:
		return "required"
	case SPDLimitOptional:
		return "optional"
	default:
		return "required"
	}
}

func (spdf *ProtoFieldLimit) MarshalJSON() (data []byte, err error) {
	s := spdf.String()
	return json.Marshal(s)
}
