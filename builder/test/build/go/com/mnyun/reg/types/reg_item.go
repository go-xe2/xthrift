package types

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
)

type RegItem struct {
	*pdl.TDynamicStructBase

	Id            int32      `thrift:"id,1,required" json:"id"`
	ParentId      int32      `thrift:"parent_id,2,required" json:"parent_id"`
	Name          string     `thrift:"name,3,required" json:"name"`
	Level         int8       `thrift:"level,4,required" json:"level"`
	ChildCount    int32      `thrift:"child_count,5,required" json:"child_count"`
	ParIds        string     `thrift:"par_ids,6,required" json:"par_ids"`
	Path          string     `thrift:"path,7,required" json:"path"`
	Time          int64      `thrift:"time,8,required" json:"time"`
	Data          []*RegItem `thrift:"data,9,required" json:"data"`
	fieldNameMaps map[string]string
	fields        map[string]*pdl.TStructFieldInfo
}

var _ pdl.DynamicStruct = (*RegItem)(nil)
var _ thrift.TStruct = (*RegItem)(nil)

func NewRegItem() *RegItem {
	inst := &RegItem{
		fieldNameMaps: make(map[string]string),
		fields:        make(map[string]*pdl.TStructFieldInfo),
	}
	inst.TDynamicStructBase = pdl.NewBasicStruct(inst)
	return inst.init()
}

func (p *RegItem) init() *RegItem {
	p.fieldNameMaps["Id"] = "Id"
	p.fieldNameMaps["id"] = "Id"

	p.fields["Id"] = pdl.NewStructFieldInfo(1, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		n32 := t.Int32(val)
		thisObj.Id = n32

		return true
	})

	p.fieldNameMaps["ParentId"] = "ParentId"
	p.fieldNameMaps["parent_id"] = "ParentId"

	p.fields["ParentId"] = pdl.NewStructFieldInfo(2, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		n32 := t.Int32(val)
		thisObj.ParentId = n32

		return true
	})

	p.fieldNameMaps["Name"] = "Name"
	p.fieldNameMaps["name"] = "Name"

	p.fields["Name"] = pdl.NewStructFieldInfo(3, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		s := t.String(val)
		thisObj.Name = s

		return true
	})

	p.fieldNameMaps["Level"] = "Level"
	p.fieldNameMaps["level"] = "Level"

	p.fields["Level"] = pdl.NewStructFieldInfo(4, thrift.I08, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		n8 := t.Int8(val)
		thisObj.Level = n8

		return true
	})

	p.fieldNameMaps["ChildCount"] = "ChildCount"
	p.fieldNameMaps["child_count"] = "ChildCount"

	p.fields["ChildCount"] = pdl.NewStructFieldInfo(5, thrift.I32, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		n32 := t.Int32(val)
		thisObj.ChildCount = n32

		return true
	})

	p.fieldNameMaps["ParIds"] = "ParIds"
	p.fieldNameMaps["par_ids"] = "ParIds"

	p.fields["ParIds"] = pdl.NewStructFieldInfo(6, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		s := t.String(val)
		thisObj.ParIds = s

		return true
	})

	p.fieldNameMaps["Path"] = "Path"
	p.fieldNameMaps["path"] = "Path"

	p.fields["Path"] = pdl.NewStructFieldInfo(7, thrift.STRING, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		s := t.String(val)
		thisObj.Path = s

		return true
	})

	p.fieldNameMaps["Time"] = "Time"
	p.fieldNameMaps["time"] = "Time"

	p.fields["Time"] = pdl.NewStructFieldInfo(8, thrift.I64, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		n64 := t.Int64(val)
		thisObj.Time = n64

		return true
	})

	p.fieldNameMaps["Data"] = "Data"
	p.fieldNameMaps["data"] = "Data"

	p.fields["Data"] = pdl.NewStructFieldInfo(9, thrift.LIST, func(obj pdl.DynamicStruct, val interface{}) bool {
		thisObj := obj.(*RegItem)
		if lst, ok := val.([]*RegItem); ok {
			thisObj.Data = lst
			return true
		}
		return false
	})

	return p
}

func (p *RegItem) Read(in thrift.TProtocol) error {
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
		if (fdId > 0 && fdId == 1) || (fdId <= 0 && fdName == "id") {
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
			p.Id = n
		}
		if (fdId > 0 && fdId == 2) || (fdId <= 0 && fdName == "parent_id") {
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
			p.ParentId = n
		}
		if (fdId > 0 && fdId == 3) || (fdId <= 0 && fdName == "name") {
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
		if (fdId > 0 && fdId == 4) || (fdId <= 0 && fdName == "level") {
			if fdId > 0 && fdType != thrift.I08 {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			n, err := in.ReadByte()
			if err != nil {
				return err
			}
			p.Level = n
		}
		if (fdId > 0 && fdId == 5) || (fdId <= 0 && fdName == "child_count") {
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
			p.ChildCount = n
		}
		if (fdId > 0 && fdId == 6) || (fdId <= 0 && fdName == "par_ids") {
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
			p.ParIds = s
		}
		if (fdId > 0 && fdId == 7) || (fdId <= 0 && fdName == "path") {
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
			p.Path = s
		}
		if (fdId > 0 && fdId == 8) || (fdId <= 0 && fdName == "time") {
			if fdId > 0 && fdType != thrift.I64 {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			n, err := in.ReadI64()
			if err != nil {
				return err
			}
			p.Time = n
		}
		if (fdId > 0 && fdId == 9) || (fdId <= 0 && fdName == "data") {
			if fdId > 0 && fdType != thrift.LIST {
				if err := in.Skip(fdType); err != nil {
					return err
				}
				if err := in.ReadFieldEnd(); err != nil {
					return err
				}
				continue
			}
			isMatch = true
			elemType, size, err := in.ReadListBegin()
			if err != nil {
				return err
			}
			if elemType != thrift.STRUCT {
				return thrift.NewTApplicationException(thrift.INVALID_PROTOCOL, "协议数据类型不匹配")
			}
			lst := make([]*RegItem, size)
			for j := 0; j < size; j++ {
				rec := NewRegItem()
				if err := rec.Read(in); err != nil {
					return err
				}
				lst[j] = rec
			}
			if err := in.ReadListEnd(); err != nil {
				return err
			}
			p.Data = lst
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

func (p *RegItem) Write(out thrift.TProtocol) error {
	if err := out.WriteStructBegin("reg_item"); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("id", thrift.I32, 1); err != nil {
		return err
	}
	if err := out.WriteI32(p.Id); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("parent_id", thrift.I32, 2); err != nil {
		return err
	}
	if err := out.WriteI32(p.ParentId); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("name", thrift.STRING, 3); err != nil {
		return err
	}
	if err := out.WriteString(p.Name); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("level", thrift.I08, 4); err != nil {
		return err
	}
	if err := out.WriteByte(p.Level); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("child_count", thrift.I32, 5); err != nil {
		return err
	}
	if err := out.WriteI32(p.ChildCount); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("par_ids", thrift.STRING, 6); err != nil {
		return err
	}
	if err := out.WriteString(p.ParIds); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("path", thrift.STRING, 7); err != nil {
		return err
	}
	if err := out.WriteString(p.Path); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if err := out.WriteFieldBegin("time", thrift.I64, 8); err != nil {
		return err
	}
	if err := out.WriteI64(p.Time); err != nil {
		return err
	}
	if err := out.WriteFieldEnd(); err != nil {
		return err
	}
	if p.Data != nil {
		if err := out.WriteFieldBegin("data", thrift.LIST, 9); err != nil {
			return err
		}
		lstSize := len(p.Data)
		if err := out.WriteListBegin(thrift.STRUCT, lstSize); err != nil {
			return err
		}
		for i := 0; i < lstSize; i++ {
			if err := p.Data[i].Write(out); err != nil {
				return err
			}
		}
		if err := out.WriteListEnd(); err != nil {
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

// 字段Data读取方法,未设置时返回默认值
func (p *RegItem) GetData() []*RegItem {
	if p.Data == nil {
		p.Data = make([]*RegItem, 0)

	}
	return p.Data
}

func (p *RegItem) NewInstance() pdl.DynamicStruct {
	return NewRegItem()
}

func (p *RegItem) AllFields() map[string]*pdl.TStructFieldInfo {
	return p.fields
}

func (p *RegItem) FieldNameMaps() map[string]string {
	return p.fieldNameMaps
}
