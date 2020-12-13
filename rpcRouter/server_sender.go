package rpcRouter

type ServerSender interface {
	SendPacket(pktData []byte)
	SendErr(pktId int64, err error, code int32)
}
