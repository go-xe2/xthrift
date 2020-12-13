package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
	"math"
	"sort"
)

type FileServiceMethod struct {
	Name      string           `json:"name"`
	Summary   string           `json:"summary"`
	Args      []*FileDataField `json:"args"`
	Result    *FileDataType    `json:"resultType"`
	Exception *FileDataType    `json:"exception"`
}

func NewFileServiceMethod(name string, summary string) *FileServiceMethod {
	return &FileServiceMethod{
		Name:      name,
		Summary:   summary,
		Args:      make([]*FileDataField, 0),
		Result:    nil,
		Exception: nil,
	}
}

func (p *FileServiceMethod) SetResult(result *FileDataType) {
	p.Result = result
}

func (p *FileServiceMethod) SetException(exception *FileDataType) {
	p.Exception = exception
}

func (p *FileServiceMethod) AddArg(arg *FileDataField) {
	p.Args = append(p.Args, arg)
}

func (p *FileServiceMethod) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_METHOD_NODE)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Name); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Summary); err != nil {
		return err
	}
	if p.Result == nil {
		p.Result = NewBaseFileDataType(SPD_VOID)
	}
	if err = p.Result.Save(writer); err != nil {
		return err
	}
	if p.Exception == nil {
		p.Exception = NewBaseFileDataType(SPD_VOID)
	}
	if err = p.Exception.Save(writer); err != nil {
		return err
	}
	size := len(p.Args)
	if size > math.MaxInt16 {
		size = math.MaxInt16
	}
	if err = writer.WriteInt16(int16(size)); err != nil {
		return err
	}
	sort.Slice(p.Args, func(i, j int) bool {
		return p.Args[i].Id-p.Args[j].Id < 0
	})

	for i := 0; i < size; i++ {
		if err = p.Args[i].Save(writer); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileServiceMethod) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_METHOD_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_METHOD_NODE, nt)
		}
	}
	if p.Name, err = reader.ReadStr(); err != nil {
		return err
	}
	if p.Summary, err = reader.ReadStr(); err != nil {
		return err
	}
	if p.Result == nil {
		p.Result = NewBaseFileDataType(SPD_UNKNOWN)
	}
	if p.Exception == nil {
		p.Exception = NewBaseFileDataType(SPD_UNKNOWN)
	}
	if err = p.Result.Load(reader); err != nil {
		return err
	}
	if err = p.Exception.Load(reader); err != nil {
		return err
	}
	n, e := reader.ReadInt16()
	if e != nil {
		return err
	}
	size := int(n)
	for i := 0; i < size; i++ {
		fd := NewFileDataField(0, "", nil, "")
		if err = fd.Load(reader); err != nil {
			return err
		}
		p.AddArg(fd)
	}
	return nil
}
