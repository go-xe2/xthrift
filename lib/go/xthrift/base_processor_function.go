/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-24 15:05
* Description:
*****************************************************************/

package xthrift

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type BaseProcessorFunction interface {
	Process(cxt context.Context, method string, seqId int32, in thrift.TProtocol, out thrift.TProtocol) (bool, thrift.TException)
	GetInputArgsInstance() thrift.TStruct
	Invoke(args thrift.TStruct) (thrift.TStruct, error)
}

type TBaseProcessorFunction struct {
	this BaseProcessorFunction
}

var _ BaseProcessorFunction = (*TBaseProcessorFunction)(nil)

func NewBaseProcessorFunction(inherited ...BaseProcessorFunction) *TBaseProcessorFunction {
	inst := &TBaseProcessorFunction{
		this: nil,
	}
	inst.this = inst
	if len(inherited) > 0 && inherited[0] != nil {
		if s, ok := inherited[0].(BaseProcessorFunction); ok {
			inst.this = s
		}
	}
	return inst
}

func (p *TBaseProcessorFunction) GetInputArgsInstance() thrift.TStruct {
	if p.this != p {
		// 调用继承类的方法
		return p.this.GetInputArgsInstance()
	}
	return nil
}

func (p *TBaseProcessorFunction) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	return nil, nil
}

func (p *TBaseProcessorFunction) Process(cxt context.Context, method string, seqId int32, in thrift.TProtocol, out thrift.TProtocol) (bool, thrift.TException) {
	args := p.GetInputArgsInstance()
	if args == nil {
		if e := in.Skip(thrift.STRUCT); e != nil {
			return false, e
		}
	} else {
		if e := args.Read(in); e != nil {
			return false, e
		}
	}
	result, err := p.this.Invoke(args)
	if err != nil {
		appEx := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, err.Error())
		if e := out.WriteMessageBegin(method, thrift.EXCEPTION, seqId); e != nil {
			return false, e
		}
		if e := appEx.Write(out); e != nil {
			return false, e
		}
		if e := out.WriteMessageEnd(); e != nil {
			return false, e
		}
		if e := out.Flush(cxt); e != nil {
			return false, e
		}
		return false, appEx
	}
	if e := out.WriteMessageBegin(method, thrift.REPLY, seqId); e != nil {
		return false, e
	}
	if e := result.Write(out); e != nil {
		return false, e
	}
	if e := out.WriteMessageEnd(); e != nil {
		return false, e
	}
	if e := out.Flush(cxt); e != nil {
		return false, e
	}
	return true, nil
}
