package rpcRouter

import "github.com/go-xe2/x/os/xlog"

func (p *TRouterClient) recvRegResult(pktId int64, proto RouterProto) {
	s, err := proto.ReadData()
	if err != nil {
		xlog.Error(err)
	}
	xlog.Info("pktId:", pktId, "注册协议成功:", s)
}
