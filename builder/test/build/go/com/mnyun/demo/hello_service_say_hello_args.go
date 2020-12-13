package demo

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type HelloServiceSayHelloArgs struct {
	*pdl.TDynamicStructBase

	Name          *string `thrift:"name,1,optional" json:"name"`
	Age           int32   `thrift:"age,2,required" json:"age"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*HelloServiceSayHelloArgs)(nil)
var _ thrift.TStruct = (*HelloServiceSayHelloArgs)(nil)

func NewHelloServiceSayHelloArgs() *HelloServiceSayHelloArgs {
	inst := &HelloServiceSayHelloArgs{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *HelloServiceSayHelloArgs) init() *HelloServiceSayHelloArgs {
	p.fieldNameMaps["Name"] = "Name"
	p.fieldNameMaps["name"] = "Name"

	p.fields["Name"] = pdl.NewStructFieldInfo(1, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloServiceSayHelloArgs)
		s := t.String(val)
		thisObj.Name = &s

		return true
	})

	p.fieldNameMaps["Age"] = "Age"
	p.fieldNameMaps["age"] = "Age"

	p.fields["Age"] = pdl.NewStructFieldInfo(2, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloServiceSayHelloArgs)
		n32 := t.Int32(val)
		thisObj.Age = n32

		return true
	})

	return p
}

func (p *HelloServiceSayHelloArgs) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 2) || (fdId <= 0 && fdName == "age") {
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
			p.Age = n
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

func (p *HelloServiceSayHelloArgs) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("hello_service_say_hello_args"); err != nil {
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
	if err := out.WriteFieldBegin("age", thrift.I32, 2); err != nil {
		return err
	}
	if err := out.WriteI32(p.Age); err != nil {
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
func (p *HelloServiceSayHelloArgs) GetName() string {
	if p.Name == nil {
		s := ""
		p.Name = &s

	}
	return *p.Name
}

func (p *HelloServiceSayHelloArgs) NewInstance() pdl.DynamicStruct {
	return NewHelloServiceSayHelloArgs()
}

func (p *HelloServiceSayHelloArgs) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *HelloServiceSayHelloArgs) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
