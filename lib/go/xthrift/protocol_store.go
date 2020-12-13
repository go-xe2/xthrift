/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-24 15:38
* Description:
*****************************************************************/

package xthrift

import (
	"github.com/apache/thrift/lib/go/thrift"
)

type TProtocolStore struct {
	thrift.TProtocol
	msgType thrift.TMessageType
	seqId   int32
	msg     string
}

var _ thrift.TProtocol = (*TProtocolStore)(nil)

func NewProtocolStore(proto thrift.TProtocol, msg string, msgType thrift.TMessageType, seqId int32) thrift.TProtocol {
	return &TProtocolStore{
		TProtocol: proto,
		msgType:   msgType,
		msg:       msg,
		seqId:     seqId,
	}
}

func (p *TProtocolStore) ReadMessageBegin() (name string, msgType thrift.TMessageType, seqId int32, err error) {
	return p.msg, p.msgType, p.seqId, nil
}
