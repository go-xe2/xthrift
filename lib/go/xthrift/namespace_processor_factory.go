/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 10:41
* Description: 命名空间协议处理器工厂
*****************************************************************/

package xthrift

import . "github.com/apache/thrift/lib/go/thrift"

type TNamespaceProcessorFactory struct {
	processor TProcessor
}

var _ TProcessorFactory = (*TNamespaceProcessorFactory)(nil)

func (pf *TNamespaceProcessorFactory) GetProcessor(trans TTransport) TProcessor {
	if pf.processor == nil {
		return NamespaceProcessor()
	}
	return pf.processor
}

func NewNamespaceProcessorFactory(processor TProcessor) TProcessorFactory {
	return &TNamespaceProcessorFactory{
		processor: processor,
	}
}

func NewNamespaceProcessorFactoryDefault() TProcessorFactory {
	return &TNamespaceProcessorFactory{
		processor: nil,
	}
}
