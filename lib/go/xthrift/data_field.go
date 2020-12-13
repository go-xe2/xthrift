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

type TPField struct {
	Name  string
	Id    int16
	Type  TProtocolDataType
	Value ProtocolDataToken
}

var _ ProtocolDataToken = (*TPField)(nil)

func NewPField(name string, id int16, typ TProtocolDataType) *TPField {
	return &TPField{
		Name:  name,
		Id:    id,
		Type:  typ,
		Value: nil,
	}
}

func (p *TPField) Write(out thrift.TTransport) error {
	return nil
}

func (p *TPField) Read(in thrift.TTransport) error {
	return nil
}

func (p *TPField) GetType() TProtocolDataType {
	return p.Type
}

func (p *TPField) SetValue(val ProtocolDataToken) error {
	if val == nil {
		return nil
	}
	if val.GetType() != p.Type {
		return fmt.Errorf("字段类型为%s,不能接受%s类型的值", p.Type, val.GetType())
	}
	p.Value = val
	return nil
}
