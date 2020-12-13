package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
)

type FileDataField struct {
	Id        int16           `json:"id"`
	Name      string          `json:"name"`
	FieldType *FileDataType   `json:"fdType"`
	Summary   string          `json:"summary"`
	Limit     ProtoFieldLimit `json:"limit"`
	Rule      string          `json:"rule"`
}

func NewFileDataField(id int16, name string, typ *FileDataType, summary string) *FileDataField {
	return &FileDataField{
		Id:        id,
		Name:      name,
		FieldType: typ,
		Summary:   summary,
	}
}

func (p *FileDataField) SetLimit(l ProtoFieldLimit) {
	p.Limit = l
}

func (p *FileDataField) SetRule(r string) {
	p.Rule = r
}

func (p *FileDataField) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_FIELD_NODE)); err != nil {
		return err
	}
	if err = writer.WriteInt16(p.Id); err != nil {
		return
	}
	if err = writer.WriteStr(p.Name); err != nil {
		return
	}
	if err = p.FieldType.Save(writer); err != nil {
		return
	}
	if err = writer.WriteStr(p.Summary); err != nil {
		return
	}
	if err = writer.WriteInt8(int8(p.Limit)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Rule); err != nil {
		return err
	}
	return nil
}

func (p *FileDataField) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_FIELD_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_FIELD_NODE, nt)
		}
	}
	if p.Id, err = reader.ReadInt16(); err != nil {
		return err
	}
	if p.Name, err = reader.ReadStr(); err != nil {
		return err
	}
	if p.FieldType == nil {
		p.FieldType = NewBaseFileDataType(SPD_UNKNOWN)
	}
	if err = p.FieldType.Load(reader); err != nil {
		return err
	}
	if p.Summary, err = reader.ReadStr(); err != nil {
		return err
	}
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		p.Limit = ProtoFieldLimit(n)
	}
	if p.Rule, err = reader.ReadStr(); err != nil {
		return err
	}
	return nil
}
