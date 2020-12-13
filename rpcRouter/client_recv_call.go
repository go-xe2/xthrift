package rpcRouter

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

func (p *TRouterClient) recvCall(pktId int64, namespace string, method string, seqId int32, data []byte) {
	cxt := context.Background()
	xlog.Debug("recvCall pktId:", pktId)
	result, err := p.rpcCall.RpcCall(cxt, namespace, method, seqId, data)
	if err != nil {
		errData := makeProtoErrorData(cxt, method, seqId, err)
		resData := makeCallReplyData(cxt, pktId, namespace, method, seqId, errData)
		if e := p.send(resData); e != nil {
			xlog.Error(e)
		}
	}
	resData := makeCallReplyData(cxt, pktId, namespace, method, seqId, result)
	xlog.Debug("recvCall reply pktId:", pktId)
	if e := p.send(resData); e != nil {
		xlog.Error(e)
	}
}

func makeProtoErrorData(ctx context.Context, method string, seqId int32, err error) []byte {
	buf := thrift.NewTMemoryBuffer()
	frame := thrift.NewTFramedTransport(buf)
	proto := xthrift.NewBinaryProtocolEx(frame)
	appErr := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, err.Error())
	if err := proto.WriteMessageBegin(method, thrift.EXCEPTION, seqId); err != nil {
		xlog.Error(err)
	}
	if err := appErr.Write(proto); err != nil {
		xlog.Error(err)
	}
	if err := proto.WriteMessageEnd(); err != nil {
		xlog.Error(err)
	}
	if err := proto.Flush(ctx); err != nil {
		xlog.Error(err)
	}
	return buf.Bytes()
}
