package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcGetCityListArgs struct {
	*pdl.TDynamicStructBase

	ProvinceId    int32 `thrift:"province_id,1,required" json:"province_id"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcGetCityListArgs)(nil)
var _ thrift.TStruct = (*RegSvcGetCityListArgs)(nil)

func NewRegSvcGetCityListArgs() *RegSvcGetCityListArgs {
	inst := &RegSvcGetCityListArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcGetCityListArgs) init() *RegSvcGetCityListArgs {
	p.fieldNameMaps["ProvinceId"] = "ProvinceId"
	p.fieldNameMaps["province_id"] = "ProvinceId"

	p.fields["ProvinceId"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcGetCityListArgs)
		n32 := t.Int32(val)
		thisObj.ProvinceId = n32

		return true
	})

	return p
}

func (p *RegSvcGetCityListArgs) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "province_id") {
			if fdId > 0 && fdType != thrift.I32 {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			n, err := in.ReadI32()
			if err != nil {
				return err
			}
			p.ProvinceId = n
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

func (p *RegSvcGetCityListArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_get_city_list_args"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("province_id", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.ProvinceId); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldStop(); err != nil {
		return err
	}
	if err := out.WriteStructEnd(); err != nil {
		return err
	}
	return nil

}

func (p *RegSvcGetCityListArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcGetCityListArgs()
}

func (p *RegSvcGetCityListArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcGetCityListArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
