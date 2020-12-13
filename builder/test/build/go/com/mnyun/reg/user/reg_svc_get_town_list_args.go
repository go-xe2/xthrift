package user

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcGetTownListArgs struct {
	*pdl.TDynamicStructBase

	CountyId      int32 `thrift:"county_id,1,required" json:"county_id"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcGetTownListArgs)(nil)
var _ thrift.TStruct = (*RegSvcGetTownListArgs)(nil)

func NewRegSvcGetTownListArgs() *RegSvcGetTownListArgs {
	inst := &RegSvcGetTownListArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcGetTownListArgs) init() *RegSvcGetTownListArgs {
	p.fieldNameMaps["CountyId"] = "CountyId"
	p.fieldNameMaps["county_id"] = "CountyId"

	p.fields["CountyId"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcGetTownListArgs)
		n32 := t.Int32(val)
		thisObj.CountyId = n32

		return true
	})

	return p
}

func (p *RegSvcGetTownListArgs) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "county_id") {
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
			p.CountyId = n
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

func (p *RegSvcGetTownListArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_get_town_list_args"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("county_id", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.CountyId); err != nil {
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

func (p *RegSvcGetTownListArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcGetTownListArgs()
}

func (p *RegSvcGetTownListArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcGetTownListArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
