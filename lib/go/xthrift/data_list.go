/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-30 17:11
* Description:
*****************************************************************/

package xthrift

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type TPList struct {
	ElemType TProtocolDataType
	Body     []ProtocolDataToken
}

var _ ProtocolDataToken = (*TPList)(nil)

func NewPList(elemType TProtocolDataType) *TPList {
	return &TPList{
		ElemType: elemType,
		Body:     make([]ProtocolDataToken, 0),
	}
}

func (p *TPList) Write(out thrift.TTransport) error {
	return nil
}

func (p *TPList) Read(in thrift.TTransport) error {
	return nil
}

func (p *TPList) GetType() TProtocolDataType {
	return thrift.LIST
}

func (p *TPList) GetElemType() TProtocolDataType {
	return p.ElemType
}

func (p *TPList) AddItem(item ProtocolDataToken) error {
	if item == nil {
		return nil
	}
	if item.GetType() != p.ElemType {
		return fmt.Errorf("列表项期望类型%s，不接受类型:%s", p.ElemType, item.GetType())
	}
	p.Body = append(p.Body, item)
	return nil
}
