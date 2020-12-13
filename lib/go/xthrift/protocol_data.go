/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-30 16:58
* Description:
*****************************************************************/

package xthrift

import (
	"github.com/apache/thrift/lib/go/thrift"
)

const TProtocolMapKeyTypes = thrift.STRING | thrift.I08 | thrift.I16 | thrift.I32 | thrift.I64 | thrift.DOUBLE

type TProtocolDataType thrift.TType

var MESSAGE TProtocolDataType = thrift.UTF16 + 1

func (pt TProtocolDataType) String() string {
	if pt == MESSAGE {
		return "message"
	}
	return thrift.TType(pt).String()
}

type ProtocolDataToken interface {
	GetType() TProtocolDataType
	Write(out thrift.TTransport) error
	Read(in thrift.TTransport) error
}
