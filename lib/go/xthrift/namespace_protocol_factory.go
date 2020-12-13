/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-23 11:03
* Description:
*****************************************************************/

package xthrift

import "github.com/apache/thrift/lib/go/thrift"

type NamespaceProtocolFactory interface {
	Namespace() string
	ConvertProtocol(toType TProtocolType) thrift.TProtocol
}

type TNamespaceProtocolFactory struct {
	protocol  thrift.TProtocol
	namespace string
}

var _ thrift.TProtocolFactory = (*TNamespaceProtocolFactory)(nil)

func NewNamespaceProtocolFactory(namespace string, protocol thrift.TProtocol) *TNamespaceProtocolFactory {
	return &TNamespaceProtocolFactory{
		namespace: namespace,
		protocol:  protocol,
	}
}

func (p *TNamespaceProtocolFactory) GetProtocol(trans thrift.TTransport) thrift.TProtocol {
	namespaceProto := NewNamespaceProtocol(p.protocol, p.namespace)
	return namespaceProto
}

func (p *TNamespaceProtocolFactory) GetNamespace() string {
	return p.namespace
}

func (p *TNamespaceProtocolFactory) ConvertProtocol(toType TProtocolType) thrift.TProtocol {
	switch toType {
	case BinaryProtocolType:
		if tmp, ok := p.protocol.(NamespaceProtocol); ok {
			if _, ok := tmp.Protocol().(*TBinaryProtocolEx); ok {
				return tmp
			} else {
				return NewNamespaceProtocol(NewBinaryProtocolEx(p.protocol.Transport()), p.namespace)
			}
		} else {
			if tmp, ok := p.protocol.(*TBinaryProtocolEx); ok {
				return tmp
			} else {
				return NewBinaryProtocolEx(p.protocol.Transport())
			}
		}
	}
	return p.protocol
}
