/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 17:51
* Description:
*****************************************************************/

package rpcRouter

import (
	"context"
	"github.com/go-xe2/x/type/xbinary"
	"io"
)

type TRouterBinaryProto struct {
	trans io.ReadWriter
}

var _ RouterProto = (*TRouterBinaryProto)(nil)

func NewRouterBinaryProto(trans io.ReadWriter) RouterProto {
	return &TRouterBinaryProto{
		trans: trans,
	}
}

func (p *TRouterBinaryProto) WriteError(msg string, code int32) error {
	if err := p.writeStr(msg); err != nil {
		return err
	}
	if err := p.writeInt32(code); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) ReadError() (msg string, code int32, err error) {
	if msg, err = p.readStr(); err != nil {
		return
	}
	if code, err = p.readInt32(); err != nil {
		return
	}
	return
}

func (p *TRouterBinaryProto) write(data []byte) error {
	if _, err := p.trans.Write(data); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) read(data []byte) error {
	if _, err := p.trans.Read(data); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) writeInt8(v int8) error {
	return p.write(xbinary.BeEncodeInt8(v))
}

func (p *TRouterBinaryProto) writeInt32(v int32) error {
	return p.write(xbinary.BeEncodeInt32(v))
}

func (p *TRouterBinaryProto) writeInt64(v int64) error {
	return p.write(xbinary.BeEncodeInt64(v))
}

func (p *TRouterBinaryProto) readInt8() (int8, error) {
	buf := make([]byte, 1)
	if err := p.read(buf); err != nil {
		return 0, err
	}
	return xbinary.BeDecodeToInt8(buf), nil
}

func (p *TRouterBinaryProto) readInt32() (int32, error) {
	buf := make([]byte, 4)
	if err := p.read(buf); err != nil {
		return 0, err
	}
	return xbinary.BeDecodeToInt32(buf), nil
}

func (p *TRouterBinaryProto) readInt64() (int64, error) {
	buf := make([]byte, 8)
	if err := p.read(buf); err != nil {
		return 0, err
	}
	return xbinary.BeDecodeToInt64(buf), nil
}

func (p *TRouterBinaryProto) readStr() (string, error) {
	size, err := p.readInt32()
	if err != nil {
		return "", err
	}
	buf := make([]byte, size)
	if err := p.read(buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (p *TRouterBinaryProto) writeStr(str string) error {
	size := len(str)
	if err := p.writeInt32(int32(size)); err != nil {
		return err
	}
	return p.write([]byte(str))
}

func (p *TRouterBinaryProto) WritePacketBegin(ptType TPacketType, pktId int64) error {
	if err := p.writeInt8(int8(ptType)); err != nil {
		return err
	}
	if err := p.writeInt64(pktId); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) WritePacketEnd() error {
	return nil
}

func (p *TRouterBinaryProto) WriteCallBegin(namesapce string, method string, seqId int32) error {
	if err := p.writeStr(namesapce); err != nil {
		return err
	}
	if err := p.writeStr(method); err != nil {
		return err
	}
	if err := p.writeInt32(seqId); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) WriteCallEnd() error {
	return nil
}

func (p *TRouterBinaryProto) WriteRegBegin(clientId string, project string, md5 string) error {
	if err := p.writeStr(clientId); err != nil {
		return err
	}
	if err := p.writeStr(project); err != nil {
		return err
	}
	if err := p.writeStr(md5); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) WriteRegEnd() error {
	return nil
}

func (p *TRouterBinaryProto) WriteData(data []byte) error {
	if err := p.writeInt64(int64(len(data))); err != nil {
		return err
	}
	if err := p.write(data); err != nil {
		return err
	}
	return nil
}

func (p *TRouterBinaryProto) ReadPacketBegin() (ptType TPacketType, pktId int64, err error) {
	var n int8
	n, err = p.readInt8()
	if err != nil {
		return
	}
	ptType = TPacketType(n)
	if ptType >= UNKOWN_PACKET {
		return
	}
	if pktId, err = p.readInt64(); err != nil {
		return
	}
	return
}

func (p *TRouterBinaryProto) ReadPacketEnd() error {
	return nil
}

func (p *TRouterBinaryProto) ReadCallBegin() (namespace string, method string, seqId int32, err error) {
	if namespace, err = p.readStr(); err != nil {
		return
	}
	if method, err = p.readStr(); err != nil {
		return
	}
	if seqId, err = p.readInt32(); err != nil {
		return
	}
	return
}

func (p *TRouterBinaryProto) ReadCallEnd() error {
	return nil
}

func (p *TRouterBinaryProto) ReadRegBegin() (clientId string, project string, md5 string, err error) {
	if clientId, err = p.readStr(); err != nil {
		return
	}
	if project, err = p.readStr(); err != nil {
		return
	}
	if md5, err = p.readStr(); err != nil {
		return
	}
	return
}

func (p *TRouterBinaryProto) ReadRegEnd() error {
	return nil
}

func (p *TRouterBinaryProto) ReadData() (data []byte, err error) {
	var size int64 = 0
	if size, err = p.readInt64(); err != nil {
		return nil, err
	}
	data = make([]byte, size)
	if err = p.read(data); err != nil {
		return
	}
	return
}

func (p *TRouterBinaryProto) Flush(ctx context.Context) error {
	if f, ok := p.trans.(Flusher); ok {
		return f.Flush(ctx)
	}
	return nil
}
