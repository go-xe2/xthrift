package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
	"math"
	"sort"
)

type FileStruct struct {
	// 所在的空间
	Type    *FileDataType    `json:"type"`
	Summary string           `json:"summary"`
	Fields  []*FileDataField `json:"fields"`
	// 引用次数
	RefCount int `json:"refCount"`
}

func NewFileStructType(namespace, name string, summary string) *FileStruct {
	typ := NewBaseFileDataType(SPD_STRUCT)
	typ.Namespace = namespace
	typ.TypName = name
	return &FileStruct{
		Type:     typ,
		Summary:  summary,
		Fields:   make([]*FileDataField, 0),
		RefCount: 0,
	}
}

func NewFileExceptType(namespace, name string, summary string) *FileStruct {
	typ := NewBaseFileDataType(SPD_EXCEPTION)
	typ.Namespace = namespace
	typ.TypName = name
	return &FileStruct{
		Type:    typ,
		Summary: summary,
		Fields:  make([]*FileDataField, 0),
	}
}

func (p *FileStruct) SetSummary(s string) {
	p.Summary = s
}

func (p *FileStruct) AddField(field *FileDataField) {
	if p.Fields == nil {
		p.Fields = make([]*FileDataField, 0)
	}
	p.Fields = append(p.Fields, field)
}

func (p *FileStruct) AddFieldByName(id int16, name string, typ *FileDataType, summary string) *FileDataField {
	result := NewFileDataField(id, name, typ, summary)
	p.AddField(result)
	return result
}

func (p *FileStruct) IncRef() {
	p.RefCount++
}

func (p *FileStruct) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_STRUCT_NODE)); err != nil {
		return err
	}
	if err = p.Type.Save(writer); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Summary); err != nil {
		return err
	}
	size := len(p.Fields)
	if size > math.MaxInt32 {
		size = math.MaxInt32
	}
	if err = writer.WriteInt32(int32(size)); err != nil {
		return err
	}
	// 按字段id排序
	sort.Slice(p.Fields, func(i, j int) bool {
		return p.Fields[i].Id-p.Fields[j].Id < 0
	})

	for i := 0; i < size; i++ {
		if i >= math.MaxInt32 {
			break
		}
		if err = p.Fields[i].Save(writer); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileStruct) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_STRUCT_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_STRUCT_NODE, nt)
		}
	}
	if p.Type == nil {
		p.Type = NewBaseFileDataType(SPD_STRUCT)
	}
	if err = p.Type.Load(reader); err != nil {
		return err
	}
	if p.Summary, err = reader.ReadStr(); err != nil {
		return err
	}
	n, e := reader.ReadInt32()
	if e != nil {
		return e
	}
	size := int(n)
	for i := 0; i < size; i++ {
		fd := NewFileDataField(0, "", nil, "")
		if err = fd.Load(reader); err != nil {
			return err
		}
		p.AddField(fd)
	}
	return nil
}
