package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcRemoveResultArgs struct {
	*pdl.TDynamicStructBase

	RegId         int32 `thrift:"reg_id,1,required" json:"reg_id"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcRemoveResultArgs)(nil)
var _ thrift.TStruct = (*RegSvcRemoveResultArgs)(nil)

func NewRegSvcRemoveResultArgs() *RegSvcRemoveResultArgs {
	inst := &RegSvcRemoveResultArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcRemoveResultArgs) init() *RegSvcRemoveResultArgs {
	p.fieldNameMaps["RegId"] = "RegId"
	p.fieldNameMaps["reg_id"] = "RegId"

	p.fields["RegId"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcRemoveResultArgs)
		n32 := t.Int32(val)
		thisObj.RegId = n32

		return true
	})

	return p
}

func (p *RegSvcRemoveResultArgs) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "reg_id") {
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
			p.RegId = n
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

func (p *RegSvcRemoveResultArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_remove_result_args"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("reg_id", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.RegId); err != nil {
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

func (p *RegSvcRemoveResultArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcRemoveResultArgs()
}

func (p *RegSvcRemoveResultArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcRemoveResultArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
