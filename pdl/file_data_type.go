package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
)

type FileDataType struct {
	// 类型
	Type TProtoBaseType `json:"type"`
	// 类型名称,struct,exception类型
	TypName   string `json:"typeName,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	// Type为list,set时的列表项类型
	ElemType *FileDataType `json:"elemType,omitempty"`
	// Type为map时的key及val类型
	KeyType *FileDataType `json:"keyType,omitempty"`
	ValType *FileDataType `json:"valType,omitempty"`
}

var BasicTypeNames = map[string]TProtoBaseType{
	"void": SPD_VOID,
	"str":  SPD_STR,
	"bl":   SPD_BOOL,
	"i8":   SPD_I08,
	"i16":  SPD_I16,
	"i32":  SPD_I32,
	"i64":  SPD_I64,
	"dl":   SPD_DOUBLE,
	//"list":      SPD_LIST,
	//"set":       SPD_SET,
	//"struct":    SPD_STRUCT,
	//"exception": SPD_EXCEPTION,
}

// 创建基础类型,str,bl,i8,i16,i32,i64,dl
func NewBaseFileDataType(typ TProtoBaseType) *FileDataType {
	return &FileDataType{
		Type:     typ,
		TypName:  "",
		ElemType: nil,
		KeyType:  nil,
		ValType:  nil,
	}
}

func NewUnknownFileDataType(typName string) *FileDataType {
	return &FileDataType{
		Type:     SPD_UNKNOWN,
		TypName:  typName,
		ElemType: nil,
		KeyType:  nil,
		ValType:  nil,
	}
}

func NewListFileDataType(elemType *FileDataType) *FileDataType {
	return &FileDataType{
		Type:     SPD_LIST,
		TypName:  "",
		ElemType: elemType,
		KeyType:  nil,
		ValType:  nil,
	}
}

func NewSetFileDataType(elemType *FileDataType) *FileDataType {
	return &FileDataType{
		Type:     SPD_SET,
		TypName:  "",
		ElemType: elemType,
		KeyType:  nil,
		ValType:  nil,
	}
}

func NewMapFileDataType(keyType, valType *FileDataType) *FileDataType {
	return &FileDataType{
		Type:     SPD_MAP,
		TypName:  "",
		ElemType: nil,
		KeyType:  keyType,
		ValType:  valType,
	}
}

func NewStructFileDataType(name string) *FileDataType {
	return &FileDataType{
		Type:     SPD_STRUCT,
		TypName:  name,
		ElemType: nil,
		KeyType:  nil,
		ValType:  nil,
	}
}

func NewExceptFileDataType(name string) *FileDataType {
	return &FileDataType{
		Type:     SPD_EXCEPTION,
		TypName:  name,
		ElemType: nil,
		KeyType:  nil,
		ValType:  nil,
	}
}

func NewFileDataTypeFromStr(str string) *FileDataType {
	if str == "" {
		return NewBaseFileDataType(SPD_VOID)
	}
	if n, ok := BasicTypeNames[str]; ok {
		return NewBaseFileDataType(n)
	}
	if b, elem := MatchProtoListType(str); b {
		et := NewFileDataTypeFromStr(elem)
		return NewListFileDataType(et)
	} else if b, elem := MatchProtoSetType(str); b {
		et := NewFileDataTypeFromStr(elem)
		return NewSetFileDataType(et)
	} else if b, key, val := MatchProtoMapType(str); b {
		kt := NewFileDataTypeFromStr(key)
		vt := NewFileDataTypeFromStr(val)
		return NewMapFileDataType(kt, vt)
	}
	return NewUnknownFileDataType(str)
}

func (p *FileDataType) Save(writer xstream.StreamWriter) (err error) {
	if err := writer.WriteInt8(int8(FNT_DATATYPE_NODE)); err != nil {
		return err
	}
	if err = writer.WriteInt8(int8(p.Type)); err != nil {
		return err
	}
	switch p.Type {
	case SPD_LIST, SPD_SET:
		if err := p.ElemType.Save(writer); err != nil {
			return err
		}
		break
	case SPD_MAP:
		if err := p.KeyType.Save(writer); err != nil {
			return err
		}
		if err := p.ValType.Save(writer); err != nil {
			return err
		}
		break
	case SPD_STRUCT, SPD_EXCEPTION, SPD_TYPEDEF:
		if err := writer.WriteStr(p.TypName); err != nil {
			return err
		}
		if err := writer.WriteStr(p.Namespace); err != nil {
			return err
		}
		break
	}
	return nil
}

func (p *FileDataType) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_DATATYPE_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_DATATYPE_NODE, nt)
		}
	}
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		p.Type = TProtoBaseType(n)
	}
	switch p.Type {
	case SPD_LIST, SPD_SET:
		if p.ElemType == nil {
			p.ElemType = NewBaseFileDataType(SPD_UNKNOWN)
		}
		if err = p.ElemType.Load(reader); err != nil {
			return err
		}
		break
	case SPD_MAP:
		if p.KeyType == nil {
			p.KeyType = NewBaseFileDataType(SPD_UNKNOWN)
		}
		if p.ValType == nil {
			p.ValType = NewBaseFileDataType(SPD_UNKNOWN)
		}
		if err = p.KeyType.Load(reader); err != nil {
			return err
		}
		if err = p.ValType.Load(reader); err != nil {
			return err
		}
		break
	case SPD_STRUCT, SPD_EXCEPTION, SPD_TYPEDEF:
		if p.TypName, err = reader.ReadStr(); err != nil {
			return err
		}
		if p.Namespace, err = reader.ReadStr(); err != nil {
			return err
		}
		break
	}
	return nil
}

func (p *FileDataType) FullName(curNamespace string) string {
	switch p.Type {
	case SPD_STR:
		return "str"
	case SPD_BOOL:
		return "bl"
	case SPD_I08:
		return "i8"
	case SPD_I16:
		return "i16"
	case SPD_I32:
		return "i32"
	case SPD_I64:
		return "i64"
	case SPD_DOUBLE:
		return "dl"
	case SPD_LIST:
		s := ""
		if p.ElemType != nil {
			s = p.ElemType.FullName(curNamespace)
		}
		return "list<" + s + ">"
	case SPD_SET:
		s := ""
		if p.ElemType != nil {
			s = p.ElemType.FullName(curNamespace)
		}
		return "set<" + s + ">"
	case SPD_MAP:
		sk, sv := "", ""
		if p.KeyType != nil {
			sk = p.KeyType.FullName(curNamespace)
		}
		if p.ValType != nil {
			sv = p.ValType.FullName(curNamespace)
		}
		return "map<" + sk + "," + sv + ">"
	case SPD_STRUCT:
		if p.Namespace != "" && p.Namespace != curNamespace {
			return fmt.Sprintf("%s.%s", p.Namespace, p.TypName)
		}
		return p.TypName
	case SPD_EXCEPTION:
		if p.Namespace != "" && p.Namespace != curNamespace {
			return fmt.Sprintf("%s.%s", p.Namespace, p.TypName)
		}
		return p.TypName
	case SPD_TYPEDEF:
		if p.Namespace != "" && p.Namespace != curNamespace {
			return fmt.Sprintf("%s.%s", p.Namespace, p.TypName)
		}
		return p.TypName
	}
	return ""
}

func (p *FileDataType) Name() string {
	switch p.Type {
	case SPD_STR:
		return "str"
	case SPD_BOOL:
		return "bl"
	case SPD_I08:
		return "i8"
	case SPD_I16:
		return "i16"
	case SPD_I32:
		return "i32"
	case SPD_I64:
		return "i64"
	case SPD_DOUBLE:
		return "dl"
	case SPD_LIST:
		s := ""
		if p.ElemType != nil {
			s = p.ElemType.Name()
		}
		return "list<" + s + ">"
	case SPD_SET:
		s := ""
		if p.ElemType != nil {
			s = p.ElemType.Name()
		}
		return "set<" + s + ">"
	case SPD_MAP:
		sk, sv := "", ""
		if p.KeyType != nil {
			sk = p.KeyType.Name()
		}
		if p.ValType != nil {
			sv = p.ValType.Name()
		}
		return "map<" + sk + "," + sv + ">"
	case SPD_STRUCT:
		if p.Namespace != "" {
			return fmt.Sprintf("%s.%s", p.Namespace, p.TypName)
		}
		return p.TypName
	case SPD_EXCEPTION:
		if p.Namespace != "" {
			return fmt.Sprintf("%s.%s", p.Namespace, p.TypName)
		}
		return p.TypName
	case SPD_TYPEDEF:
		return p.TypName
	}
	return ""
}

func (p *FileDataType) String() string {
	return p.Name()
}
