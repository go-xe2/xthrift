/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-30 17:12
* Description:
*****************************************************************/

package xthrift

import "github.com/apache/thrift/lib/go/thrift"

type TPMap struct {
	KeyType thrift.TType
	ValType thrift.TType
	Size    int32
	Body    []*TPField
}
