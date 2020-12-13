/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 09:06
* Description: 使用命名空间的协议处理器
*****************************************************************/

package xthrift

import (
	"context"
	"fmt"
	. "github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/core/exception"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/sync/xsafeMap"
	"strings"
)

var namespaceProcessInstances = xsafeMap.NewStrAnyMap()

// 命名空间分割符
const NAMESPACE_SEPARATOR = "."

type TNamespaceProcessor struct {
	processors       map[string]TProcessor
	defaultProcessor TProcessor
}

var _ TProcessor = (*TNamespaceProcessor)(nil)

func NamespaceProcessor(mgrName ...string) *TNamespaceProcessor {
	name := "defaultProcessor"
	if len(mgrName) > 0 {
		name = mgrName[0]
	}
	if v := namespaceProcessInstances.Get(name); v != nil {
		return v.(*TNamespaceProcessor)
	}
	inst := &TNamespaceProcessor{
		processors:       make(map[string]TProcessor),
		defaultProcessor: DefaultServiceProcessor,
	}
	namespaceProcessInstances.Set(name, inst)
	return inst
}

func (p *TNamespaceProcessor) Namespaces() []string {
	result := make([]string, 0)
	if p.defaultProcessor != nil {
		result = append(result, "[default]")
	}
	for k := range p.processors {
		result = append(result, k)
	}
	return result
}

func (p *TNamespaceProcessor) RegisterNamespace(serviceName string, processor TProcessor) error {
	if processor == nil {
		return nil
	}
	if _, ok := p.processors[serviceName]; ok {
		return exception.NewText("命名空间%s已经被注册过")
	}
	p.processors[serviceName] = processor
	return nil
}

func (p *TNamespaceProcessor) RegisterDefault(processor TProcessor) {
	p.defaultProcessor = processor
}

func (p *TNamespaceProcessor) returnError(cxt context.Context, msg string, seqId int32, out TProtocol, err error) TApplicationException {
	appErr := NewTApplicationException(UNKNOWN_APPLICATION_EXCEPTION, err.Error())
	if err := out.WriteMessageBegin(msg, EXCEPTION, seqId); err != nil {
		xlog.Error(err)
	}
	if err := appErr.Write(out); err != nil {
		xlog.Error(err)
	}
	if err := out.WriteMessageEnd(); err != nil {
		xlog.Error(err)
	}
	if err := out.Flush(cxt); err != nil {
		xlog.Error(err)
	}
	return appErr
}

func (p *TNamespaceProcessor) Process(ctx context.Context, in, out TProtocol) (bool, TException) {
	name, typeId, seqid, err := in.ReadMessageBegin()
	fmt.Println("process name:", name, ", typeId:", typeId, ", seqId:", seqid, ", err:", err)
	if err != nil {
		return false, err
	}
	if typeId == ONEWAY+1 {
		xlog.Debug("心跳包.")
		if err := out.WriteMessageBegin(name, typeId, seqid); err != nil {
			return false, err
		}
		if err := out.WriteMessageEnd(); err != nil {
			return false, err
		}
		if err := out.Flush(ctx); err != nil {
			return false, err
		}
		return true, nil
	}
	if typeId != CALL && typeId != ONEWAY {
		return true, nil
	}
	inputProto := in
	if proto, ok := in.(DynamicProtocol); ok {
		protoType := proto.GetProtocolType()
		switch protoType {
		case BinaryProtocolType:
			inputProto = NewBinaryProtocolEx(in.Transport())
			break
		}
	}
	var callProto TProtocol
	v := strings.Split(name, NAMESPACE_SEPARATOR)
	if len(v) < 2 {
		if p.defaultProcessor != nil {
			callProto = NewProtocolStore(inputProto, name, typeId, seqid)
			return p.defaultProcessor.Process(ctx, callProto, out)
		}
		_ = p.returnError(ctx, name, seqid, out, fmt.Errorf("未找到服务: %s. 客户端请使用TNamespaceProtocol协议调用", name))
		return true, nil
	}
	namespace := strings.Join(v[:len(v)-1], NAMESPACE_SEPARATOR)
	msgName := v[len(v)-1]
	processor, ok := p.processors[namespace]
	if !ok {
		_ = p.returnError(ctx, msgName, seqid, out, fmt.Errorf("服务%s不存在，是否已经注册", namespace))
		return true, nil
	}
	callProto = NewProtocolStore(inputProto, msgName, typeId, seqid)
	return processor.Process(ctx, callProto, out)
}
