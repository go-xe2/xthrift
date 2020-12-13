package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/pdl"
)

type RegSvcGetProvincesListArgs struct {
	*pdl.TDynamicStructBase

	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegSvcGetProvincesListArgs)(nil)
var _ thrift.TStruct = (*RegSvcGetProvincesListArgs)(nil)

func NewRegSvcGetProvincesListArgs() *RegSvcGetProvincesListArgs {
	inst := &RegSvcGetProvincesListArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegSvcGetProvincesListArgs) init() *RegSvcGetProvincesListArgs {
	return p
}

func (p *RegSvcGetProvincesListArgs) Read(in thrift.TProtocol) error {
	if err := in.Skip(thrift.STRUCT); err != nil {
		return err
	}
	return nil

}

func (p *RegSvcGetProvincesListArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_svc_get_provinces_list_args"); err != nil {
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

func (p *RegSvcGetProvincesListArgs) NewInstance() pdl.DynamicStruct {
	return NewRegSvcGetProvincesListArgs()
}

func (p *RegSvcGetProvincesListArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegSvcGetProvincesListArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
