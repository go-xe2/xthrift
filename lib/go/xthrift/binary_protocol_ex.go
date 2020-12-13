/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 11:54
* Description:
*****************************************************************/

package xthrift

import (
	. "github.com/apache/thrift/lib/go/thrift"
)

// 继承thrift.TBinaryProtocol, 修正为传输数据结构名称及字段名
type TBinaryProtocolEx struct {
	*TBinaryProtocol
	protocolType TProtocolType
}

var _ DynamicProtocol = (*TBinaryProtocolEx)(nil)

func NewBinaryProtocolEx(trans TTransport) TProtocol {
	return &TBinaryProtocolEx{
		protocolType:    UnknownProtocolType,
		TBinaryProtocol: NewTBinaryProtocolTransport(trans),
	}
}

func (p *TBinaryProtocolEx) GetProtocolType() TProtocolType {
	return p.protocolType
}

func (p *TBinaryProtocolEx) WriteMessageBegin(name string, msgType TMessageType, seqId int32) error {
	if err := p.WriteByte(int8(BinaryProtocolType)); err != nil {
		return err
	}
	return p.TBinaryProtocol.WriteMessageBegin(name, msgType, seqId)
}

func (p *TBinaryProtocolEx) ReadMessageBegin() (name string, msgType TMessageType, seqId int32, err error) {
	var n int8 = 0
	if n, err = p.ReadByte(); err != nil {
		return
	} else {
		p.protocolType = TProtocolType(n)
	}
	return p.TBinaryProtocol.ReadMessageBegin()
}

func (p *TBinaryProtocolEx) WriteStructBegin(name string) error {
	return p.WriteString(name)
}

func (p *TBinaryProtocolEx) ReadStructBegin() (name string, err error) {
	s, err := p.ReadString()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (p *TBinaryProtocolEx) WriteFieldBegin(name string, fdType TType, fdId int16) error {
	if err := p.WriteByte(int8(fdType)); err != nil {
		return err
	}
	if err := p.WriteString(name); err != nil {
		return err
	}
	if err := p.WriteI16(fdId); err != nil {
		return err
	}
	return nil
}

func (p *TBinaryProtocolEx) ReadFieldBegin() (name string, fdType TType, fdId int16, err error) {
	var e error
	name = ""
	fdId = 0
	n, e := p.ReadByte()
	if e != nil {
		err = NewTProtocolException(err)
		return
	}
	fdType = TType(n)
	if n != STOP {
		if name, err = p.ReadString(); err != nil {
			return
		}
		if fdId, err = p.ReadI16(); err != nil {
			return
		}
	}
	return
}
