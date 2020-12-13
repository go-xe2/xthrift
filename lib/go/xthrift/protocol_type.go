/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-24 14:28
* Description:
*****************************************************************/

package xthrift

type TProtocolType byte

const (
	UnknownProtocolType TProtocolType = iota
	BinaryProtocolType
)

func (pt TProtocolType) String() string {
	switch pt {
	case BinaryProtocolType:
		return "BinaryProtocolType"
	}
	return "UnknownProtocolType"
}
