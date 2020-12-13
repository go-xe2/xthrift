package types

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type HelloData struct {
	*pdl.TDynamicStructBase

	Name          *string `thrift:"name,1,optional" json:"name"`
	Sex           bool    `thrift:"sex,2,required" json:"sex"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*HelloData)(nil)
var _ thrift.TStruct = (*HelloData)(nil)

func NewHelloData() *HelloData {
	inst := &HelloData{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *HelloData) init() *HelloData {
	p.fieldNameMaps["Name"] = "Name"
	p.fieldNameMaps["name"] = "Name"

	p.fields["Name"] = pdl.NewStructFieldInfo(1, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloData)
		s := t.String(val)
		thisObj.Name = &s

		return true
	})

	p.fieldNameMaps["Sex"] = "Sex"
	p.fieldNameMaps["sex"] = "Sex"

	p.fields["Sex"] = pdl.NewStructFieldInfo(2, thrift.BOOL, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloData)
		b := t.Bool(val)
		thisObj.Sex = b

		return true
	})

	return p
}

func (p *HelloData) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "name") {
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
			p.Name = &s
		}
		if (fdId > 0 && fdId == 2) || (fdId <= 0 && fdName == "sex") {
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
			p.Sex = b
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

func (p *HelloData) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("hello_data"); err != nil {
		return err
	}
	if p.Name != nil {
		if err := out.WriteFieldBegin("name", thrift.STRING, 1); err != nil {
			return err
		}
		if err := out.WriteString(*p.Name); err != nil {
			return err
		}
		if err := out.WriteFieldEnd(); err != nil {
			return err
		}
	}
	if err := out.WriteFieldBegin("sex", thrift.BOOL, 2); err != nil {
		return err
	}
	if err := out.WriteBool(p.Sex); err != nil {
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

// 字段Name读取方法,未设置时返回默认值
func (p *HelloData) GetName() string {
	if p.Name == nil {
		s := ""
		p.Name = &s

	}
	return *p.Name
}

func (p *HelloData) NewInstance() pdl.DynamicStruct {
	return NewHelloData()
}

func (p *HelloData) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *HelloData) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
