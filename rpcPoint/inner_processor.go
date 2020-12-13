/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-24 11:24
* Description:
*****************************************************************/

package rpcPoint

import (
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"strings"
)
import "golang.org/x/net/context"

type tInnerProcessor struct {
	server *TEndPointServer
	trans  thrift.TTransport
}

var _ thrift.TProcessor = (*tInnerProcessor)(nil)

func NewInnerProcessor(server *TEndPointServer) *tInnerProcessor {
	return &tInnerProcessor{
		server: server,
	}
}

type tInnerProcessorFactory struct {
	server *TEndPointServer
}

var _ thrift.TProcessorFactory = (*tInnerProcessorFactory)(nil)

func newInnerProcessorFactory(server *TEndPointServer) *tInnerProcessorFactory {
	return &tInnerProcessorFactory{
		server: server,
	}
}

func (p *tInnerProcessorFactory) GetProcessor(trans thrift.TTransport) thrift.TProcessor {
	return &tInnerProcessor{
		server: p.server,
		trans:  trans,
	}
}

func (p *tInnerProcessor) returnError(ctx context.Context, out thrift.TProtocol, msgName string, seqId int32, err error) error {
	appErr := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, err.Error())
	if err := out.WriteMessageBegin(msgName, thrift.EXCEPTION, seqId); err != nil {
		return err
	}
	if err := appErr.Write(out); err != nil {
		return err
	}
	if err := out.WriteMessageEnd(); err != nil {
		return err
	}
	if err := out.Flush(ctx); err != nil {
		return err
	}
	return nil
}

func (p *tInnerProcessor) returnFrameError(ctx context.Context, out thrift.TProtocol, msgName string, seqId int32, err error) error {
	outBuf := thrift.NewTMemoryBuffer()
	outFrameTrans := thrift.NewTFramedTransport(outBuf)
	outFrameProto := p.server.outProtoFac.GetProtocol(outFrameTrans)
	e1 := p.returnError(ctx, outFrameProto, msgName, seqId, err)
	if e1 != nil {
		return e1
	}
	if _, e1 := out.Transport().Write(outBuf.Bytes()); e1 != nil {
		return e1
	}
	if e1 := out.Transport().Flush(ctx); e1 != nil {
		return e1
	}
	return nil
}

func (p *tInnerProcessor) Process(ctx context.Context, in, out thrift.TProtocol) (bool, thrift.TException) {
	inBufTrans := thrift.NewTMemoryBuffer()
	frameInTrans := thrift.NewTFramedTransport(inBufTrans)
	inBuf := p.server.inProtoFac.GetProtocol(frameInTrans)

	msgName, msgType, seqId, err := ProtocolTransform(in, inBuf)

	items := strings.Split(msgName, ".")
	if len(items) < 2 {
		_ = p.returnFrameError(ctx, out, msgName, seqId, errors.New("请使用命命空间限定名称方式调用"))
		return true, nil
	}
	fullSvcName := strings.Join(items[:len(items)-1], ".")
	methodName := items[len(items)-1]

	xlog.Debug("recv message:", msgName, ", seqId:", seqId, ", msgType:", msgType, ", err:", err)
	if err != nil {
		return false, err
	}
	if err := frameInTrans.Flush(ctx); err != nil {
		return false, p.returnFrameError(ctx, out, methodName, seqId, err)
	}

	if msgType != thrift.CALL && msgType != thrift.ONEWAY {
		return false, nil
	}

	result, e := p.server.RpcCall(ctx, fullSvcName, methodName, seqId, inBufTrans.Bytes())

	if e != nil {
		xlog.Error("call ", msgName, " error:", e)
		if err := p.returnFrameError(ctx, out, methodName, seqId, e); err != nil {
			return false, err
		}
		return true, nil
		//appErr := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, e.Error())
		//if err := out.WriteMessageBegin(msgName, thrift.EXCEPTION, seqId); err != nil {
		//	return false, err
		//}
		//if err := appErr.Write(out); err != nil {
		//	return false, err
		//}
		//if err := out.WriteMessageEnd(); err != nil {
		//	return false, err
		//}
		//if err := out.Flush(ctx); err != nil {
		//	return false, err
		//}
		//return true, nil
	}
	if _, err := out.Transport().Write(result); err != nil {
		return false, err
	}
	if err := out.Transport().Flush(ctx); err != nil {
		return false, err
	}
	xlog.Debug("返回数据成功")
	return true, nil
}
