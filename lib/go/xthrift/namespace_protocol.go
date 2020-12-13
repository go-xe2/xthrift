/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 09:13
* Description: 使用命名空间调用的客户端协议代理器
*****************************************************************/

package xthrift

import (
	"fmt"
	. "github.com/apache/thrift/lib/go/thrift"
)

// 使用命名空间的协议调用,该协议只使用在客户端，服务器端使用TNamespaceProcessor接收处理数据
// 示例
//socket := thrift.NewTSocketFromAddrTimeout(addr, TIMEOUT)
//transport := thrift.NewTFramedTransport(socket)
//protocol := thrift.NewTBinaryProtocolTransport(transport)
//
//mp := thrift.NewNamespaceProtocol(protocol, "Calculator")
//service := Calculator.NewCalculatorClient(mp)
//
//mp2 := thrift.NewNamespaceProtocol(protocol, "WeatherReport")
//service2 := WeatherReport.NewWeatherReportClient(mp2)
//
//err := transport.Open()
//if err != nil {
//t.Fatal("Unable to open client socket", err)
//}
//
//fmt.Println(service.Add(2,2))
//fmt.Println(service2.GetTemperature())

type NamespaceProtocol interface {
	TProtocol
	Namespace() string
	Protocol() TProtocol
}

type TNamespaceProtocol struct {
	TProtocol
	namespace string
}

var _ DynamicProtocol = (*TNamespaceProtocol)(nil)
var _ NamespaceProtocol = (*TNamespaceProtocol)(nil)

func NewNamespaceProtocol(protocol TProtocol, namespace string) *TNamespaceProtocol {
	return &TNamespaceProtocol{
		TProtocol: protocol,
		namespace: namespace,
	}
}

func (p *TNamespaceProtocol) Protocol() TProtocol {
	return p.TProtocol
}

func (p *TNamespaceProtocol) Namespace() string {
	return p.namespace
}

func (p *TNamespaceProtocol) GetProtocolType() TProtocolType {
	if tmp, ok := p.TProtocol.(DynamicProtocol); ok {
		return tmp.GetProtocolType()
	}
	return UnknownProtocolType
}

func (p *TNamespaceProtocol) WriteMessageBegin(name string, msgType TMessageType, seqId int32) error {
	if msgType == CALL || msgType == ONEWAY {
		return p.TProtocol.WriteMessageBegin(fmt.Sprintf("%s%s%s", p.namespace, NAMESPACE_SEPARATOR, name), msgType, seqId)
	} else {
		return p.TProtocol.WriteMessageBegin(name, msgType, seqId)
	}
}
