package rpcRouter

import (
	"bytes"
	"context"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
)

// 发送错误消息
func makeErrorData(pktId int64, msg string, code int32) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	proto := NewRouterBinaryProto(buf)
	if err := proto.WritePacketBegin(ERR_RES_PACKET, pktId); err != nil {
		return nil, err
	}
	if err := proto.WriteError(msg, code); err != nil {
		return nil, err
	}
	if err := proto.WritePacketEnd(); err != nil {
		return nil, err
	}
	if err := proto.Flush(context.Background()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func makeCallReplyData(ctx context.Context, pktId int64, namespace string, method string, seqId int32, data []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	proto := NewRouterBinaryProto(buf)
	if err := proto.WritePacketBegin(REPLY_PACKET, pktId); err != nil {
		xlog.Error(err)
	}
	if err := proto.WriteCallBegin(namespace, method, seqId); err != nil {
		xlog.Error(err)
	}
	if err := proto.WriteData(data); err != nil {
		xlog.Error(err)
	}
	if err := proto.WriteCallEnd(); err != nil {
		xlog.Error(err)
	}
	if err := proto.WritePacketEnd(); err != nil {
		xlog.Error(err)
	}
	if err := proto.Flush(ctx); err != nil {
		xlog.Error(err)
	}
	return buf.Bytes()
}

func makeRegResData(ctx context.Context, pktId int64, projectName string, md5 string) []byte {
	outBuf := bytes.NewBuffer([]byte{})
	outProto := NewRouterBinaryProto(outBuf)
	if err := outProto.WritePacketBegin(REG_RES_PACKET, pktId); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteData([]byte(projectName + "," + md5)); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WritePacketEnd(); err != nil {
		xlog.Error(err)
	}
	if err := outProto.Flush(ctx); err != nil {
		xlog.Error(err)
	}
	return outBuf.Bytes()
}

func makeRegData(ctx context.Context, clientId string, pktId int64, proj *pdl.FileProject, md5 string) []byte {
	outBuf := bytes.NewBuffer([]byte{})
	outProto := NewRouterBinaryProto(outBuf)
	projBuf := bytes.NewBuffer([]byte{})
	if err := proj.SaveProject(projBuf); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WritePacketBegin(REG_PACKET, pktId); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteRegBegin(clientId, proj.GetProjectName(), md5); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteData(projBuf.Bytes()); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteRegEnd(); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WritePacketEnd(); err != nil {
		xlog.Error(err)
	}
	if err := outProto.Flush(ctx); err != nil {
		xlog.Error(err)
	}
	return outBuf.Bytes()
}

func makeCallData(ctx context.Context, pktId int64, namespace string, method string, seqId int32, rpcData []byte) []byte {
	outBuf := bytes.NewBuffer([]byte{})
	outProto := NewRouterBinaryProto(outBuf)

	if err := outProto.WritePacketBegin(CALL_PACKET, pktId); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteCallBegin(namespace, method, seqId); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteData(rpcData); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WriteCallEnd(); err != nil {
		xlog.Error(err)
	}
	if err := outProto.WritePacketEnd(); err != nil {
		xlog.Error(err)
	}
	if err := outProto.Flush(ctx); err != nil {
		xlog.Error(err)
	}
	return outBuf.Bytes()
}
