package demo

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type HelloResult struct {
	*pdl.TDynamicStructBase

	Status        int32  `thrift:"status,1,required" json:"status"`
	Msg           string `thrift:"msg,2,required" json:"msg"`
	Content       string `thrift:"content,3,required" json:"content"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*HelloResult)(nil)
var _ thrift.TStruct = (*HelloResult)(nil)

func NewHelloResult() *HelloResult {
	inst := &HelloResult{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *HelloResult) init() *HelloResult {
	p.fieldNameMaps["Status"] = "Status"
	p.fieldNameMaps["status"] = "Status"

	p.fields["Status"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloResult)
		n32 := t.Int32(val)
		thisObj.Status = n32

		return true
	})

	p.fieldNameMaps["Msg"] = "Msg"
	p.fieldNameMaps["msg"] = "Msg"

	p.fields["Msg"] = pdl.NewStructFieldInfo(2, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloResult)
		s := t.String(val)
		thisObj.Msg = s

		return true
	})

	p.fieldNameMaps["Content"] = "Content"
	p.fieldNameMaps["content"] = "Content"

	p.fields["Content"] = pdl.NewStructFieldInfo(3, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*HelloResult)
		s := t.String(val)
		thisObj.Content = s

		return true
	})

	return p
}

func (p *HelloResult) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "status") {
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
			p.Status = n
		}
		if (fdId > 0 && fdId == 2) || (fdId <= 0 && fdName == "msg") {
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
			p.Msg = s
		}
		if (fdId > 0 && fdId == 3) || (fdId <= 0 && fdName == "content") {
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
			p.Content = s
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

func (p *HelloResult) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("hello_result"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("status", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.Status); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("msg", thrift.STRING, 2); err != nil {
		return err
	}
	if err := out.WriteString(p.Msg); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("content", thrift.STRING, 3); err != nil {
		return err
	}
	if err := out.WriteString(p.Content); err != nil {
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

func (p *HelloResult) NewInstance() pdl.DynamicStruct {
	return NewHelloResult()
}

func (p *HelloResult) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *HelloResult) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
