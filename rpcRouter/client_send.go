package rpcRouter

import (
	"context"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
	"time"
)

func (p *TRouterClient) send(data []byte) error {
	if _, err := p.client.Send(data); err != nil {
		return err
	}
	return nil
}

func (p *TRouterClient) sendError(pktId int64, msg string, code int32) {
	data, err := makeErrorData(pktId, msg, code)
	if err != nil {
		xlog.Error(err)
	}
	if err := p.send(data); err != nil {
		xlog.Error(err)
	}
}

func (p *TRouterClient) sendRegProject(proj *pdl.FileProject, md5 string) error {
	pktId := time.Now().UnixNano()

	data := makeRegData(context.Background(), p.clientId, pktId, proj, md5)
	return p.send(data)
	//
	//buf := bytes.NewBuffer([]byte{})
	//proto := NewRouterBinaryProto(buf)
	//pktId := time.Now().UnixNano()
	//projDataBuf := bytes.NewBuffer([]byte{})
	//if err := proj.SaveProject(projDataBuf); err != nil {
	//	return err
	//}
	//if err := proto.WritePacketBegin(REG_PACKET, pktId); err != nil {
	//	return err
	//}
	//if err := proto.WriteRegBegin(p.clientId, proj.GetProjectName(), md5); err != nil {
	//	return err
	//}
	//if err := proto.WriteData(projDataBuf.Bytes()); err != nil {
	//	return err
	//}
	//if err := proto.WriteRegEnd(); err != nil {
	//	return err
	//}
	//if err := proto.WritePacketEnd(); err != nil {
	//	return err
	//}
	//if err := proto.Flush(context.Background()); err != nil {
	//	return nil
	//}
	//return p.send(buf.Bytes())
}
