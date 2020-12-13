/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 17:36
* Description:
*****************************************************************/

package rpcRouter

import "context"

type TPacketType int8

const (
	// reg -> reg_res; call -> reply
	// 协议注册包
	REG_PACKET TPacketType = iota
	// 注册返回
	REG_RES_PACKET
	// 调用接口数据包
	CALL_PACKET
	// 接口调用回复数据包
	REPLY_PACKET
	// 异常返回
	ERR_RES_PACKET
	// 未知类型数据包
	UNKOWN_PACKET
)

func (p TPacketType) String() string {
	switch p {
	case REG_PACKET:
		return "REGISTER"
	case REG_RES_PACKET:
		return "REG_RES"
	case CALL_PACKET:
		return "CALL"
	case REPLY_PACKET:
		return "REPLY"
	case ERR_RES_PACKET:
		return "ERR_RES"
	case UNKOWN_PACKET:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

type RouterProto interface {
	WritePacketBegin(ptType TPacketType, pktId int64) error
	WritePacketEnd() error

	WriteCallBegin(namesapce string, method string, seqId int32) error
	WriteCallEnd() error

	WriteRegBegin(clientId string, project string, md5 string) error
	WriteRegEnd() error

	WriteData(data []byte) error

	WriteError(msg string, code int32) error

	ReadPacketBegin() (ptType TPacketType, pktId int64, err error)
	ReadPacketEnd() error
	ReadCallBegin() (namespace string, method string, seqId int32, err error)
	ReadCallEnd() error

	ReadRegBegin() (clientId string, project string, md5 string, err error)
	ReadRegEnd() error

	ReadData() (data []byte, err error)

	ReadError() (msg string, code int32, err error)
	Flush(ctx context.Context) error
}
