package rpcRouter

import "github.com/go-xe2/x/os/xlog"

func (p *TRouterClient) recvError(proto RouterProto, pktId int64) {
	msg, code, err := proto.ReadError()
	if err != nil {
		xlog.Error(err)
	}
	xlog.Debug("注册出错: pktId:", pktId, "msg:", msg, ", code:", code)
}
