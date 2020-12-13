/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-24 14:34
* Description:
*****************************************************************/

package xthrift

import "github.com/apache/thrift/lib/go/thrift"

type DynamicProtocol interface {
	thrift.TProtocol
	GetProtocolType() TProtocolType
}
