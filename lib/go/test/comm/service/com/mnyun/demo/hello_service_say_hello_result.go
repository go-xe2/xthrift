package demo

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/pdl"
)

type HelloServiceSayHelloResult struct {
	*pdl.TDynamicStructBase

	Success       *HelloResult `thrift:"success,1,optional" json:"success"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*HelloServiceSayHelloResult)(nil)
var _ thrift.TStruct = (*HelloServiceSayHelloResult)(nil)

func NewHelloServiceSayHelloResult() *HelloServiceSayHelloResult {
	inst := &HelloServiceSayHelloResult{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *HelloServiceSayHelloResult) init() *HelloServiceSayHelloResult {
	p.fieldNameMaps["Success"] = "Success"
	p.fieldNameMaps["success"] = "Success"

	p.fields["Success"] = pdl.NewStructFieldInfo(1, thrift.STRUCT, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloServiceSayHelloResult)
		if stru, ok := val.(*HelloResult); ok {
			thisObj.Success = stru
			return true
		}
		return false
	})

	return p
}

func (p *HelloServiceSayHelloResult) Read(in thrift.TProtocol) error {
	_, err := in.ReadStructBegin()
	if err != nil {
		return err
	}
	var nMaxLoop = 512
	nLoop := 0
	var isMatch bool
	for {
		// 防止协议数据错误，无thrift.STOP时无限循环
		nLoop++
		if nLoop >= nMaxLoop {
			_ = in.Skip(thrift.STRUCT)
			return nil
		}
		isMatch = false
		fdName, fdType, fdId, err := in.ReadFieldBegin()
		if err != nil {
			return err
		}
		if fdType == thrift.STOP {
			break
		}
		if fdType == thrift.VOID {
			if err := in.ReadFieldEnd(); err != nil {
				return err
			}
			continue
		}
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "success") {
			if fdId > 0 && fdType != thrift.STRUCT {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			p.Success = NewHelloResult()
			if err := p.Success.Read(in); err != nil {
				return err
			}
		}
		if !isMatch {
			if err := in.Skip(fdType); err != nil {
				return err
			}
		}
		if err := in.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := in.ReadStructEnd(); err != nil {
		return err
	}
	return nil

}

func (p *HelloServiceSayHelloResult) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("hello_service_say_hello_result"); err != nil {
		return err
	}
	if p.Success != nil {
		if err := out.WriteFieldBegin("success", thrift.STRUCT, 1); err != nil {
			return err
		}
		if err := p.Success.Write(out); err != nil {
			return err
		}
		if err := out.WriteFieldEnd(); err != nil {
			return err
		}
	}
	if err := out.WriteFieldStop(); err != nil {
		return err
	}
	if err := out.WriteStructEnd(); err != nil {
		return err
	}
	return nil

}

// 字段Success读取方法,未设置时返回默认值
func (p *HelloServiceSayHelloResult) GetSuccess() *HelloResult {
	if p.Success == nil {
		p.Success = NewHelloResult()

	}
	return p.Success
}

func (p *HelloServiceSayHelloResult) NewInstance() pdl.DynamicStruct {
	return NewHelloServiceSayHelloResult()
}

func (p *HelloServiceSayHelloResult) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *HelloServiceSayHelloResult) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
