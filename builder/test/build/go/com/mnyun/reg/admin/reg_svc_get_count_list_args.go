package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcGetCountListArgs struct {
	*pdl.TDynamicStructBase

	CityId        int32 `thrift:"city_id,1,required" json:"city_id"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcGetCountListArgs)(nil)
var _ thrift.TStruct = (*RegSvcGetCountListArgs)(nil)

func NewRegSvcGetCountListArgs() *RegSvcGetCountListArgs {
	inst := &RegSvcGetCountListArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcGetCountListArgs) init() *RegSvcGetCountListArgs {
	p.fieldNameMaps["CityId"] = "CityId"
	p.fieldNameMaps["city_id"] = "CityId"

	p.fields["CityId"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcGetCountListArgs)
		n32 := t.Int32(val)
		thisObj.CityId = n32

		return true
	})

	return p
}

func (p *RegSvcGetCountListArgs) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "city_id") {
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
			p.CityId = n
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

func (p *RegSvcGetCountListArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_get_count_list_args"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("city_id", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.CityId); err != nil {
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

func (p *RegSvcGetCountListArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcGetCountListArgs()
}

func (p *RegSvcGetCountListArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcGetCountListArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
