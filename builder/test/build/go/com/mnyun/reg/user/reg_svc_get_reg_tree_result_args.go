package user

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcGetRegTreeResultArgs struct {
	*pdl.TDynamicStructBase

	ParId         *int32 `thrift:"par_id,1,optional" json:"par_id"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcGetRegTreeResultArgs)(nil)
var _ thrift.TStruct = (*RegSvcGetRegTreeResultArgs)(nil)

func NewRegSvcGetRegTreeResultArgs() *RegSvcGetRegTreeResultArgs {
	inst := &RegSvcGetRegTreeResultArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcGetRegTreeResultArgs) init() *RegSvcGetRegTreeResultArgs {
	p.fieldNameMaps["ParId"] = "ParId"
	p.fieldNameMaps["par_id"] = "ParId"

	p.fields["ParId"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcGetRegTreeResultArgs)
		n32 := t.Int32(val)
		thisObj.ParId = &n32

		return true
	})

	return p
}

func (p *RegSvcGetRegTreeResultArgs) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "par_id") {
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
			p.ParId = &n
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

func (p *RegSvcGetRegTreeResultArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_get_reg_tree_result_args"); err != nil {
		return err
	}
	if p.ParId != nil {
		if err := out.WriteFieldBegin("par_id", thrift.I32, 1); err != nil {
			return err
		}
		if err := out.WriteI32(*p.ParId); err != nil {
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

// 字段ParId读取方法,未设置时返回默认值
func (p *RegSvcGetRegTreeResultArgs) GetParId() int32 {
	if p.ParId == nil {
		var n32 int32 = 0
		p.ParId = &n32

	}
	return *p.ParId
}

func (p *RegSvcGetRegTreeResultArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcGetRegTreeResultArgs()
}

func (p *RegSvcGetRegTreeResultArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcGetRegTreeResultArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
