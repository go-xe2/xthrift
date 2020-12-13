package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
)

type FileTypeDef struct {
	Name    string        `json:"name"`
	OrgType *FileDataType `json:"orgType"`
}

func NewFileTypeDef(name string, orgType *FileDataType) *FileTypeDef {
	return &FileTypeDef{
		Name:    name,
		OrgType: orgType,
	}
}

func (p *FileTypeDef) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_TYPEDEF_NODE)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Name); err != nil {
		return err
	}
	if err = p.OrgType.Save(writer); err != nil {
		return err
	}
	return nil
}

func (p *FileTypeDef) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_TYPEDEF_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_TYPEDEF_NODE, nt)
		}
	}
	if p.Name, err = reader.ReadStr(); err != nil {
		return err
	}
	if p.OrgType == nil {
		p.OrgType = NewBaseFileDataType(SPD_UNKNOWN)
	}
	if err = p.OrgType.Load(reader); err != nil {
		return err
	}
	return nil
}
