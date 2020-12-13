package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xstream"
	"sort"
	"strings"
)

type FileData struct {
	project   *FileProject
	namespace *FileNamespace
	filePath  string
	fileName  string
	// namespace
	Namespace string
	// imports节点
	Imports []string
	// defTypes 节点
	Typedefs map[string]*FileTypeDef
	// type节点
	Types    map[string]*FileStruct
	Services map[string]*FileService
}

func NewFileData(project *FileProject) *FileData {
	return &FileData{
		project:  project,
		filePath: "",
		fileName: "",
		Imports:  make([]string, 0),
		Typedefs: make(map[string]*FileTypeDef, 0),
		Types:    make(map[string]*FileStruct, 0),
		Services: make(map[string]*FileService, 0),
	}
}

func (p *FileData) bindNamespace(ns *FileNamespace) {
	p.namespace = ns
}

func (p *FileData) SetNamespace(namespace string) {
	p.Namespace = namespace
}

func (p *FileData) AddImport(namespace string) {
	p.Imports = append(p.Imports, namespace)
}

func (p *FileData) AddTypedef(typeDef *FileTypeDef) {
	p.Typedefs[typeDef.Name] = typeDef
}

func (p *FileData) AddType(typ *FileStruct) {
	p.Types[typ.Type.TypName] = typ
}

func (p *FileData) AddService(svc *FileService) {
	p.Services[svc.Name] = svc
}

func (p *FileData) QryTypedef(defName string) *FileTypeDef {
	if v, ok := p.Typedefs[defName]; ok {
		return v
	}
	return nil
}

func (p *FileData) QryType(typeName string) *FileStruct {
	if v, ok := p.Types[typeName]; ok {
		return v
	}
	return nil
}

func (p *FileData) QryService(svcName string) *FileService {
	if svc, ok := p.Services[svcName]; ok {
		return svc
	}
	return nil
}

func (p *FileData) checkDataType(dt *FileDataType, imports map[string]int) (TProtoBaseType, error) {
	switch dt.Type {
	case SPD_LIST:
		elemT, err := p.checkDataType(dt.ElemType, imports)
		if err != nil {
			return SPD_VOID, fmt.Errorf("list列表项定义错误:%s", err.Error())
		}
		dt.ElemType.Type = elemT
		break
	case SPD_SET:
		elemT, err := p.checkDataType(dt.ElemType, imports)
		if err != nil {
			return SPD_VOID, fmt.Errorf("set列表项定义错误:%s", err.Error())
		}
		dt.ElemType.Type = elemT
		break
	case SPD_MAP:
		kt, err := p.checkDataType(dt.KeyType, imports)
		if err != nil {
			return SPD_VOID, fmt.Errorf("map的key定义错误:%s", err.Error())
		}
		dt.KeyType.Type = kt
		vt, err := p.checkDataType(dt.ValType, imports)
		if err != nil {
			return SPD_VOID, fmt.Errorf("map的val定义错误:%s", err.Error())
		}
		dt.ValType.Type = vt
		break
	case SPD_UNKNOWN:
		if strings.Index(dt.TypName, ".") > 0 {
			s1, s2 := NamespaceLastName(dt.TypName)
			if _, ok := imports[s1]; !ok {
				return SPD_VOID, fmt.Errorf("未导入类型%s的命名空间", dt.TypName)
			}
			_, df := p.project.QryTypeDefByNS(s1, s2)
			if df != nil {
				// 增加引用计数
				imports[s1]++
				dt.Namespace = s1
				dt.TypName = s2
				return SPD_TYPEDEF, nil
			}
			_, t := p.project.QryTypeByNS(s1, s2)
			if t == nil {
				return SPD_VOID, fmt.Errorf("数据类型%s未定义", dt.TypName)
			}
			// 增加引用计数
			imports[s1]++
			dt.Namespace = s1
			dt.TypName = s2
			return t.Type.Type, nil
		}
		// 在本命名空间中查找
		if df := p.namespace.QryTypedef(dt.TypName); df != nil {
			dt.Type = SPD_TYPEDEF
			dt.Namespace = p.Namespace
			return SPD_TYPEDEF, nil
		}
		if t := p.namespace.QryType(dt.TypName); t != nil {
			dt.Type = t.Type.Type
			dt.Namespace = p.Namespace
			return t.Type.Type, nil
		}
		return SPD_VOID, fmt.Errorf("类型%s未定义", dt.TypName)
	}
	return dt.Type, nil
}

func (p *FileData) checkTypedef(imports map[string]int) error {
	for _, dt := range p.Typedefs {
		t, err := p.checkDataType(dt.OrgType, imports)
		if err != nil {
			return fmt.Errorf("类型别名%s定义错误，类型%s定义错误:%s", dt.Name, dt.OrgType.TypName, err)
		}
		dt.OrgType.Type = t
	}
	return nil
}

func (p *FileData) checkTypes(imports map[string]int) error {
	for _, dt := range p.Types {
		ids := make(map[int16]bool)
		for _, fd := range dt.Fields {
			t, err := p.checkDataType(fd.FieldType, imports)
			if err != nil {
				return fmt.Errorf("字段%s.%s类型定义错误:%s", dt.Type.TypName, fd.Name, err)
			}
			if _, ok := ids[fd.Id]; ok {
				return fmt.Errorf("字段%s.%sID:%d重复", dt.Type.TypName, fd.Name, fd.Id)
			}
			ids[fd.Id] = true
			fd.FieldType.Type = t
		}
	}
	return nil
}

func (p *FileData) checkServices(imports map[string]int) error {
	for _, svc := range p.Services {
		if err := p.checkService(svc, imports); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileData) checkService(svc *FileService, imports map[string]int) error {
	for _, m := range svc.Methods {
		t, err := p.checkDataType(m.Result, imports)
		if err != nil {
			return fmt.Errorf("接口%s.%s返回类型定义错误:%s", svc.Name, m.Name, err)
		}
		m.Result.Type = t
		t, err = p.checkDataType(m.Exception, imports)
		if err != nil {
			return fmt.Errorf("接口%s.%sthrow类型定义错误:%s", svc.Name, m.Name, err)
		}
		ids := make(map[int16]bool)
		for _, arg := range m.Args {
			t, err := p.checkDataType(arg.FieldType, imports)
			if err != nil {
				return fmt.Errorf("接口%s.%s输入参数%s类型定义错误:%s", svc.Name, m.Name, arg.Name, err)
			}
			if _, ok := ids[arg.Id]; ok {
				return fmt.Errorf("接口%s.%s输入参数%sID:%d重复", svc.Name, m.Name, arg.Name, arg.Id)
			}
			ids[arg.Id] = true
			arg.FieldType.Type = t
		}
	}
	return nil
}

// 检查数据定义
func (p *FileData) Check() error {
	importMap := make(map[string]int)
	for _, s := range p.Imports {
		importMap[s] = 0
	}
	// 检查typedef节点数据定义
	if err := p.checkTypedef(importMap); err != nil {
		return err
	}
	if err := p.checkTypes(importMap); err != nil {
		return err
	}
	if err := p.checkServices(importMap); err != nil {
		return err
	}
	return nil
}

func (p *FileData) Margin(other *FileData) {
	if other == nil {
		return
	}
	imports := make(map[string]int)
	for _, s := range p.Imports {
		imports[s] = 1
	}
	// 合并imports
	for _, s := range other.Imports {
		if _, ok := imports[s]; !ok {
			p.Imports = append(p.Imports, s)
		}
	}
	// 合并typedefs
	for k, def := range other.Typedefs {
		if _, ok := p.Typedefs[k]; !ok {
			p.Typedefs[k] = def
		}
	}
	// 合并types
	for k, t := range other.Types {
		if _, ok := p.Types[k]; !ok {
			p.Types[k] = t
		}
	}
	// 合并services
	for k, svc := range other.Services {
		self, ok := p.Services[k]
		if !ok {
			p.Services[k] = svc
		} else {
			// 合并service的接口
			self.Margin(svc)
		}
	}
}

func (p *FileData) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_FILE_NODE)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.fileName); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Namespace); err != nil {
		return err
	}
	size := len(p.Typedefs)
	if err = writer.WriteInt64(int64(size)); err != nil {
		return err
	}
	defKeys := make([]string, size)
	i := 0
	for k := range p.Typedefs {
		defKeys[i] = k
		i++
	}
	sort.Slice(defKeys, func(i, j int) bool {
		return strings.Compare(defKeys[i], defKeys[j]) < 0
	})

	for _, k := range defKeys {
		def := p.Typedefs[k]
		if err = def.Save(writer); err != nil {
			return err
		}
	}
	size = len(p.Types)
	if err = writer.WriteInt64(int64(size)); err != nil {
		return err
	}

	typKeys := make([]string, size)
	i = 0
	for k := range p.Types {
		typKeys[i] = k
		i++
	}
	sort.Slice(typKeys, func(i, j int) bool {
		return strings.Compare(typKeys[i], typKeys[j]) < 0
	})

	for _, k := range typKeys {
		t := p.Types[k]
		if err = t.Save(writer); err != nil {
			return err
		}
	}

	size = len(p.Services)
	if err = writer.WriteInt64(int64(size)); err != nil {
		return err
	}

	svcKeys := make([]string, size)
	i = 0
	for k := range p.Services {
		svcKeys[i] = k
		i++
	}
	sort.Slice(svcKeys, func(i, j int) bool {
		return strings.Compare(typKeys[i], typKeys[j]) < 0
	})

	for _, k := range svcKeys {
		svc := p.Services[k]
		if err = svc.Save(writer); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileData) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_FILE_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_FILE_NODE, nt)
		}
	}
	if p.fileName, err = reader.ReadStr(); err != nil {
		return err
	}
	if p.Namespace, err = reader.ReadStr(); err != nil {
		return err
	}
	n, e := reader.ReadInt64()
	if e != nil {
		return e
	}
	size := int(n)
	for i := 0; i < size; i++ {
		def := NewFileTypeDef("", nil)
		if err = def.Load(reader); err != nil {
			return err
		}
		p.Typedefs[def.Name] = def
	}

	n, e = reader.ReadInt64()
	if e != nil {
		return e
	}
	size = int(n)
	for i := 0; i < size; i++ {
		t := NewFileStructType("", "", "")
		if err = t.Load(reader); err != nil {
			return err
		}
		p.Types[t.Type.TypName] = t
	}

	n, e = reader.ReadInt64()
	if e != nil {
		return e
	}
	size = int(n)
	for i := 0; i < size; i++ {
		svc := NewFileService(p.Namespace, "", "")
		if err = svc.Load(reader); err != nil {
			return err
		}
		p.Services[svc.Name] = svc
	}
	return nil
}

func (p *FileData) Export(expt FileExport) error {
	ext := xfile.Ext(p.fileName)
	name := p.fileName[:len(p.fileName)-len(ext)]
	w, cxt, err := expt.BeginFileWrite(p.namespace, name)
	if err != nil {
		return err
	}
	defer expt.EndFileWrite(w, p.namespace, name)
	if err := expt.WriteNamespace(w, cxt, p.Namespace); err != nil {
		return err
	}
	if err := expt.WriteImports(w, cxt, p.Imports); err != nil {
		return err
	}
	if err := expt.WriteTypedefs(w, cxt, p.Typedefs); err != nil {
		return err
	}
	if err := expt.WriteTypes(w, cxt, p.Types); err != nil {
		return err
	}
	if err := expt.WriteServices(w, cxt, p.Services); err != nil {
		return err
	}
	if err := expt.Flush(w, cxt); err != nil {
		return err
	}
	return nil
}
