/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 14:44
* Description:
*****************************************************************/

package xthrift

import "github.com/apache/thrift/lib/go/thrift"

type TBinaryProtocolExFactory struct {
}

func (pf *TBinaryProtocolExFactory) GetProtocol(trans thrift.TTransport) thrift.TProtocol {
	return NewBinaryProtocolEx(trans)
}

func NewBinaryProtocolExFactory() thrift.TProtocolFactory {
	return &TBinaryProtocolExFactory{}
}
