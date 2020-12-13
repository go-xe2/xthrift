/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-30 17:10
* Description:
*****************************************************************/

package xthrift

import "github.com/apache/thrift/lib/go/thrift"

type TPMessage struct {
	Name  string
	Type  thrift.TMessageType
	SeqId int32
}

var _ ProtocolDataToken = (*TPMessage)(nil)

func NewPMessage(name string, typ thrift.TMessageType, seqId int32) *TPMessage {
	return &TPMessage{
		Name:  name,
		Type:  typ,
		SeqId: seqId,
	}
}

func (p *TPMessage) Write(out thrift.TTransport) error {
	return nil
}

func (p *TPMessage) Read(in thrift.TTransport) error {
	return nil
}

func (p *TPMessage) GetType() TProtocolDataType {
	return MESSAGE
}
