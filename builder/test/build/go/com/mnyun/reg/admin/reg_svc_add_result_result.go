package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcAddResultResult struct {
	*pdl.TDynamicStructBase

	Success       *bool `thrift:"success,1,optional" json:"success"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcAddResultResult)(nil)
var _ thrift.TStruct = (*RegSvcAddResultResult)(nil)

func NewRegSvcAddResultResult() *RegSvcAddResultResult {
	inst := &RegSvcAddResultResult{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcAddResultResult) init() *RegSvcAddResultResult {
	p.fieldNameMaps["Success"] = "Success"
	p.fieldNameMaps["success"] = "Success"

	p.fields["Success"] = pdl.NewStructFieldInfo(1, thrift.BOOL, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcAddResultResult)
		b := t.Bool(val)
		thisObj.Success = &b

		return true
	})

	return p
}

func (p *RegSvcAddResultResult) Read(in thrift.TProtocol) error {
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
			if fdId > 0 && fdType != thrift.BOOL {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			b, err := in.ReadBool()
			if err != nil {
				return err
			}
			p.Success = &b
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

func (p *RegSvcAddResultResult) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_add_result_result"); err != nil {
		return err
	}
	if p.Success != nil {
		if err := out.WriteFieldBegin("success", thrift.BOOL, 1); err != nil {
			return err
		}
		if err := out.WriteBool(*p.Success); err != nil {
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
func (p *RegSvcAddResultResult) GetSuccess() bool {
	if p.Success == nil {
		b := false
		p.Success = &b

	}
	return *p.Success
}

func (p *RegSvcAddResultResult) NewInstance() pdl.DynamicStruct {
	return NewRegSvcAddResultResult()
}

func (p *RegSvcAddResultResult) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcAddResultResult) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
