package user

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/reg/types"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcGetCityListResult struct {
	*pdl.TDynamicStructBase

	Success       []*types.RegItem `thrift:"success,1,optional" json:"success"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcGetCityListResult)(nil)
var _ thrift.TStruct = (*RegSvcGetCityListResult)(nil)

func NewRegSvcGetCityListResult() *RegSvcGetCityListResult {
	inst := &RegSvcGetCityListResult{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcGetCityListResult) init() *RegSvcGetCityListResult {
	p.fieldNameMaps["Success"] = "Success"
	p.fieldNameMaps["success"] = "Success"

	p.fields["Success"] = pdl.NewStructFieldInfo(1, thrift.LIST, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcGetCityListResult)
		if lst, ok := val.([]*types.RegItem); ok {
			thisObj.Success = lst
			return true
		}
		return false
	})

	return p
}

func (p *RegSvcGetCityListResult) Read(in thrift.TProtocol) error {
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
			if fdId > 0 && fdType != thrift.LIST {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			elemType, size, err := in.ReadListBegin()
			if err != nil {
				return err
			}
			if elemType != thrift.STRUCT {
				return thrift.NewTApplicationException(thrift.INVALID_PROTOCOL, "协议数据类型不匹配")
			}
			lst := make([]*types.RegItem, size)
			for j := 0; j < size; j++ {
				rec := types.NewRegItem()
				if err := rec.Read(in); err != nil {
					return err
				}
				lst[j] = rec
			}
			if err := in.ReadListEnd(); err != nil {
				return err
			}
			p.Success = lst
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

func (p *RegSvcGetCityListResult) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_get_city_list_result"); err != nil {
		return err
	}
	if p.Success != nil {
		if err := out.WriteFieldBegin("success", thrift.LIST, 1); err != nil {
			return err
		}
		lstSize := len(p.Success)
		if err := out.WriteListBegin(thrift.STRUCT, lstSize); err != nil {
			return err
		}
		for i := 0; i < lstSize; i++ {
			if err := p.Success[i].Write(out); err != nil {
				return err
			}
		}
		if err := out.WriteListEnd(); err != nil {
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
func (p *RegSvcGetCityListResult) GetSuccess() []*types.RegItem {
	if p.Success == nil {
		p.Success = make([]*types.RegItem, 0)

	}
	return p.Success
}

func (p *RegSvcGetCityListResult) NewInstance() pdl.DynamicStruct {
	return NewRegSvcGetCityListResult()
}

func (p *RegSvcGetCityListResult) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcGetCityListResult) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
