package rpcRouter

import (
	"bytes"
	"errors"
	"fmt"
)

type tServerRPCSender struct {
	status chan int
	data   []byte
	err    error
}

var _ ServerSender = (*tServerRPCSender)(nil)

func newServerRPCSender() *tServerRPCSender {
	return &tServerRPCSender{
		status: make(chan int, 1),
		data:   nil,
		err:    nil,
	}
}

func (p *tServerRPCSender) SendPacket(pktData []byte) {
	buf := bytes.NewBuffer(pktData)
	proto := NewRouterBinaryProto(buf)
	pktType, _, err := proto.ReadPacketBegin()
	if err != nil {
		p.setResult(nil, err)
		return
	}
	defer proto.ReadPacketEnd()
	switch pktType {
	case REPLY_PACKET:
		_, _, _, err := proto.ReadCallBegin()
		if err != nil {
			p.setResult(nil, err)
			return
		}
		data, err := proto.ReadData()
		if err != nil {
			p.setResult(nil, err)
			return
		}
		if err := proto.ReadCallEnd(); err != nil {
			p.setResult(nil, err)
			return
		}
		p.setResult(data, nil)
		break
	case ERR_RES_PACKET:
		s, code, err := proto.ReadError()
		if err != nil {
			p.setResult(nil, err)
		} else {
			p.setResult(nil, fmt.Errorf("%s(%d)", s, code))
		}
		break
	default:
		p.setResult(nil, errors.New("未知数据包类型"))
	}
}

func (p *tServerRPCSender) setResult(data []byte, err error) {
	select {
	case _, ok := <-p.status:
		if !ok {
			return
		}
	default:
	}
	p.err = err
	p.data = data
	if err != nil {
		p.status <- 1
	} else {
		p.status <- 0
	}
}

func (p *tServerRPCSender) SendErr(pktId int64, err error, code int32) {
	p.setResult(nil, err)
}

func (p *tServerRPCSender) Wait() ([]byte, error) {
	select {
	case _, ok := <-p.status:
		if !ok {
			return nil, errors.New("请求超时2")
		}
		return p.data, p.err
	}
}
