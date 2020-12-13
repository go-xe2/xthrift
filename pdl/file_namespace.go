package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
	"sort"
	"strings"
)

type FileNamespace struct {
	project   *FileProject
	Namespace string               `json:"namespace"`
	Files     map[string]*FileData `json:"files"`
}

func NewFileNamespace(project *FileProject, namespace string) *FileNamespace {
	return &FileNamespace{
		project:   project,
		Namespace: namespace,
		Files:     make(map[string]*FileData),
	}
}

func (p *FileNamespace) AddFile(file *FileData) {
	p.Files[file.fileName] = file
}

func (p *FileNamespace) QryTypedef(defName string) *FileTypeDef {
	for _, f := range p.Files {
		if v := f.QryTypedef(defName); v != nil {
			return v
		}
	}
	return nil
}

func (p *FileNamespace) QryType(typeName string) *FileStruct {
	for _, f := range p.Files {
		if v := f.QryType(typeName); v != nil {
			return v
		}
	}
	return nil
}

func (p *FileNamespace) QryService(svcName string) *FileService {
	for _, f := range p.Files {
		if v := f.QryService(svcName); v != nil {
			return v
		}
	}
	return nil
}

func (p *FileNamespace) QryMethod(svcName, methodName string) *FileServiceMethod {
	svc := p.QryService(svcName)
	if svc == nil {
		return nil
	}
	if m := svc.QryMethod(methodName); m != nil {
		return m
	}
	return nil
}

// 合并空间
func (p *FileNamespace) Margin(other *FileNamespace) {
	if other == nil {
		return
	}
	if other.Namespace != p.Namespace {
		return
	}
	// 合并文件
	for szName, f := range other.Files {
		self, ok := p.Files[szName]
		if !ok {
			p.Files[szName] = f
		} else {
			self.Margin(f)
		}
	}
}

func (p *FileNamespace) Check() error {
	for _, f := range p.Files {
		if err := f.Check(); err != nil {
			return fmt.Errorf("文件%s错误:%s", f.fileName, err)
		}
	}
	return nil
}

func (p *FileNamespace) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_NAMESPACE_NODE)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Namespace); err != nil {
		return err
	}
	size := len(p.Files)
	if err = writer.WriteInt64(int64(size)); err != nil {
		return err
	}
	// 对文件名进行排序，以防止项目每次保存时文件md5值变动
	keys := make([]string, size)
	i := 0
	for k := range p.Files {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})

	for _, k := range keys {
		f := p.Files[k]
		if err = f.Save(writer); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileNamespace) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_NAMESPACE_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_NAMESPACE_NODE, nt)
		}
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
		f := NewFileData(p.project)
		if err = f.Load(reader); err != nil {
			return err
		}
		f.bindNamespace(p)
		p.Files[f.fileName] = f
	}
	return nil
}

func (p *FileNamespace) Export(expt FileExport) error {
	if err := expt.BeginNamespace(p.Namespace); err != nil {
		return err
	}
	defer expt.EndNamespace(p.Namespace)
	for _, f := range p.Files {
		if err := f.Export(expt); err != nil {
			return err
		}
	}
	return nil
}
