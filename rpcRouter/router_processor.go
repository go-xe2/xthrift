package rpcRouter

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/comm/protoTrans"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"strings"
)

type TRouterProcessor struct {
	caller RouterCaller
}

var _ thrift.TProcessor = (*TRouterProcessor)(nil)

type TRouterProcessorFactory struct {
	caller RouterCaller
}

var _ thrift.TProcessorFactory = (*TRouterProcessorFactory)(nil)

func (p *TRouterProcessorFactory) GetProcessor(trans thrift.TTransport) thrift.TProcessor {
	return NewRouterProcessor(p.caller)
}

func NewRouterProcessorFactory(caller RouterCaller) *TRouterProcessorFactory {
	return &TRouterProcessorFactory{caller: caller}
}

func NewRouterProcessor(caller RouterCaller) thrift.TProcessor {
	return &TRouterProcessor{caller: caller}
}

func (p *TRouterProcessor) returnError(cxt context.Context, msg string, seqId int32, out thrift.TProtocol, err error) thrift.TApplicationException {
	appErr := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, err.Error())
	if err := out.WriteMessageBegin(msg, thrift.EXCEPTION, seqId); err != nil {
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

// 只支持binaryProtocol及TFrameTransport封包
func (p *TRouterProcessor) Process(ctx context.Context, in, out thrift.TProtocol) (bool, thrift.TException) {
	buf := thrift.NewTMemoryBuffer()
	frame := thrift.NewTFramedTransport(buf)
	inBufProto := xthrift.NewBinaryProtocolEx(frame)
	name, typeId, seqid, err := protoTrans.ProtocolTransform(in, inBufProto)

	if err != nil {
		return false, err
	}
	if typeId == thrift.ONEWAY+1 {
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
	if typeId != thrift.CALL && typeId != thrift.ONEWAY {
		return true, nil
	}

	v := strings.Split(name, xthrift.NAMESPACE_SEPARATOR)
	if len(v) < 2 {
		_ = p.returnError(ctx, name, seqid, out, fmt.Errorf("未找到服务: %s. 客户端请使用TNamespaceProtocol协议调用", name))
		return true, nil
	}
	namespace := strings.Join(v[:len(v)-1], xthrift.NAMESPACE_SEPARATOR)
	method := v[len(v)-1]
	// 转发服务
	if err := frame.Flush(ctx); err != nil {
		_ = p.returnError(ctx, method, seqid, out, err)
		return true, nil
	}
	res, err := p.caller.RouterCall(ctx, namespace, method, seqid, buf.Bytes())
	if err != nil {
		xlog.Debug("rpc call error:", err)
		_ = p.returnError(ctx, method, seqid, out, err)
		return true, nil
	}
	if _, err := out.Transport().Write(res[4:]); err != nil {
		xlog.Error(err)
		return false, err
	}
	if err := out.Flush(ctx); err != nil {
		xlog.Error(err)
		return false, err
	}
	return true, nil
}
