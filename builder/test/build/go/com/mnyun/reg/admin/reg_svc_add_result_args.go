package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcAddResultArgs struct {
	*pdl.TDynamicStructBase

	ParId         int32  `thrift:"par_id,1,required" json:"par_id"`
	Name          string `thrift:"name,2,required" json:"name"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcAddResultArgs)(nil)
var _ thrift.TStruct = (*RegSvcAddResultArgs)(nil)

func NewRegSvcAddResultArgs() *RegSvcAddResultArgs {
	inst := &RegSvcAddResultArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcAddResultArgs) init() *RegSvcAddResultArgs {
	p.fieldNameMaps["ParId"] = "ParId"
	p.fieldNameMaps["par_id"] = "ParId"

	p.fields["ParId"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcAddResultArgs)
		n32 := t.Int32(val)
		thisObj.ParId = n32

		return true
	})

	p.fieldNameMaps["Name"] = "Name"
	p.fieldNameMaps["name"] = "Name"

	p.fields["Name"] = pdl.NewStructFieldInfo(2, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegSvcAddResultArgs)
		s := t.String(val)
		thisObj.Name = s

		return true
	})

	return p
}

func (p *RegSvcAddResultArgs) Read(in thrift.TProtocol) error {
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
			p.ParId = n
		}
		if (fdId > 0 && fdId == 2) || (fdId <= 0 && fdName == "name") {
			if fdId > 0 && fdType != thrift.STRING {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			s, err := in.ReadString()
			if err != nil {
				return err
			}
			p.Name = s
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

func (p *RegSvcAddResultArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_add_result_args"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("par_id", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.ParId); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("name", thrift.STRING, 2); err != nil {
		return err
	}
	if err := out.WriteString(p.Name); err != nil {
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

func (p *RegSvcAddResultArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcAddResultArgs()
}

func (p *RegSvcAddResultArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcAddResultArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
