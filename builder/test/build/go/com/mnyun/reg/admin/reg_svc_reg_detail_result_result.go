package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/reg/types"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcRegDetailResultResult struct {
	*pdl.TDynamicStructBase

	Success       *types.RegItem `thrift:"success,1,optional" json:"success"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcRegDetailResultResult)(nil)
var _ thrift.TStruct = (*RegSvcRegDetailResultResult)(nil)

func NewRegSvcRegDetailResultResult() *RegSvcRegDetailResultResult {
	inst := &RegSvcRegDetailResultResult{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcRegDetailResultResult) init() *RegSvcRegDetailResultResult {
	p.fieldNameMaps["Success"] = "Success"
	p.fieldNameMaps["success"] = "Success"

	p.fields["Success"] = pdl.NewStructFieldInfo(1, thrift.STRUCT, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcRegDetailResultResult)
		if stru, ok := val.(*types.RegItem); ok {
			thisObj.Success = stru
			return true
		}
		return false
	})

	return p
}

func (p *RegSvcRegDetailResultResult) Read(in thrift.TProtocol) error {
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
			p.Success = types.NewRegItem()
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

func (p *RegSvcRegDetailResultResult) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_reg_detail_result_result"); err != nil {
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
func (p *RegSvcRegDetailResultResult) GetSuccess() *types.RegItem {
	if p.Success == nil {
		p.Success = types.NewRegItem()

	}
	return p.Success
}

func (p *RegSvcRegDetailResultResult) NewInstance() pdl.DynamicStruct {
	return NewRegSvcRegDetailResultResult()
}

func (p *RegSvcRegDetailResultResult) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcRegDetailResultResult) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
