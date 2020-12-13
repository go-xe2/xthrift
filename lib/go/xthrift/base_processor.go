/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 12:22
* Description:
*****************************************************************/

package xthrift

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type BaseProcessor interface {
	thrift.TProcessor
	GetFunction(fnName string) BaseProcessorFunction
}

type TBaseProcessor struct {
	this   BaseProcessor
	fnMaps map[string]BaseProcessorFunction
}

var _ thrift.TProcessor = (*TBaseProcessor)(nil)

var DefaultServiceProcessor = NewBaseProcessor()

func NewBaseProcessor(inherited ...thrift.TProcessor) *TBaseProcessor {
	inst := &TBaseProcessor{
		fnMaps: make(map[string]BaseProcessorFunction),
	}
	inst.this = inst
	if len(inherited) > 0 && inherited[0] != nil {
		if s, ok := inherited[0].(BaseProcessor); ok {
			inst.this = s
		} else {
			inst.this = inst
		}
	}
	return inst
}

func (p *TBaseProcessor) RegisterFunction(name string, fn BaseProcessorFunction) {
	p.fnMaps[name] = fn
}

func (p *TBaseProcessor) GetFunction(fnName string) BaseProcessorFunction {
	if fn, ok := p.fnMaps[fnName]; ok {
		return fn
	}
	return nil
}

func (p *TBaseProcessor) Process(cxt context.Context, in, out thrift.TProtocol) (bool, thrift.TException) {
	msg, _, seqId, err := in.ReadMessageBegin()
	if err != nil {
		if e := in.Skip(thrift.STRUCT); err != nil {
			return false, e
		}
		if e := in.ReadMessageEnd(); err != nil {
			return false, e
		}
		return false, err
	}
	fn := p.this.GetFunction(msg)
	// 检查接口是否存在，不存在返回错误信息
	if fn == nil {
		if e := in.ReadMessageEnd(); e != nil {
			return false, e
		}
		appEx := thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, fmt.Sprintf("服务不存在接口%s,或是否已经注册接口?", msg))
		if e := out.WriteMessageBegin(msg, thrift.EXCEPTION, seqId); e != nil {
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
		fmt.Println("写入接口不存在错误信息成功======>>")
		return false, appEx
	}
	// 创建输入参数
	b, ex := fn.Process(cxt, msg, seqId, in, out)
	if e := in.ReadMessageEnd(); e != nil {
		return false, e
	}
	return b, ex
}
