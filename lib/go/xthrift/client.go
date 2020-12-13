/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-23 14:47
* Description:
*****************************************************************/

package xthrift

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type TXClient struct {
	seqId         int32
	trans         thrift.TTransport
	inputFactory  thrift.TProtocolFactory
	outputFactory thrift.TProtocolFactory
}

var _ thrift.TClient = (*TXClient)(nil)

func NewClient(trans thrift.TTransport, inFac, outFac thrift.TProtocolFactory) *TXClient {
	return &TXClient{
		trans:         trans,
		inputFactory:  inFac,
		outputFactory: outFac,
	}
}

func (p *TXClient) Send(ctx context.Context, out thrift.TProtocol, method string, seqId int32, args thrift.TStruct) error {
	// Set headers from context object on THeaderProtocol
	if headerProt, ok := out.(*thrift.THeaderProtocol); ok {
		headerProt.ClearWriteHeaders()
		for _, key := range thrift.GetWriteHeaderList(ctx) {
			if value, ok := thrift.GetHeader(ctx, key); ok {
				headerProt.SetWriteHeader(key, value)
			}
		}
	}
	if err := out.WriteMessageBegin(method, thrift.CALL, seqId); err != nil {
		return err
	}
	if err := args.Write(out); err != nil {
		return err
	}
	if err := out.WriteMessageEnd(); err != nil {
		return err
	}
	return out.Flush(ctx)
}

func (p *TXClient) Recv(in thrift.TProtocol, method string, seqId int32, result thrift.TStruct) error {
	rMethod, rTypeId, rSeqId, err := in.ReadMessageBegin()
	var readPro = in

	if tmpPro, ok := in.(DynamicProtocol); ok {
		protoType := tmpPro.GetProtocolType()
		if fac, ok := p.inputFactory.(NamespaceProtocolFactory); ok {
			readPro = fac.ConvertProtocol(protoType)
		} else {
			switch protoType {
			case BinaryProtocolType:
				if np, ok := in.(NamespaceProtocol); ok {
					if _, ok := np.Protocol().(*TBinaryProtocolEx); !ok {
						readPro = NewNamespaceProtocol(NewBinaryProtocolEx(in.Transport()), np.Namespace())
					}
				} else {
					if _, ok := in.(*TBinaryProtocolEx); !ok {
						readPro = NewBinaryProtocolEx(in.Transport())
					}
				}
			}
		}
	}
	if err != nil {
		return err
	}
	if method != rMethod {
		fmt.Println("recv method:", rMethod, ", need method:", method)
		return thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, fmt.Sprintf("wrong method name:%s", rMethod))
	} else if seqId != rSeqId {
		return thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, fmt.Sprintf("out of order sequence response:%s", rMethod))
	} else if rTypeId == thrift.EXCEPTION {
		var exception = thrift.NewTApplicationException(0, "")
		if err := exception.Read(readPro); err != nil {
			return err
		}
		if err := in.ReadMessageEnd(); err != nil {
			return err
		}
		return exception
	} else if rTypeId != thrift.REPLY {
		return thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, fmt.Sprintf("%s: invalid message type", method))
	}
	if err := result.Read(readPro); err != nil {
		return err
	}
	return in.ReadMessageEnd()
}

func (p *TXClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	p.seqId++
	seqId := p.seqId
	if err := p.Send(ctx, p.outputFactory.GetProtocol(p.trans), method, seqId, args); err != nil {
		return err
	}
	// method is oneway
	if result == nil {
		return nil
	}
	return p.Recv(p.inputFactory.GetProtocol(p.trans), method, seqId, result)
}
