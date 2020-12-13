package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xstream"
	"math"
	"sort"
	"strings"
)

type FileService struct {
	Namespace string                        `json:"namespace"`
	Name      string                        `json:"name"`
	Summary   string                        `json:"summary"`
	Methods   map[string]*FileServiceMethod `json:"methods"`
}

func NewFileService(namespace string, name string, summary string) *FileService {
	return &FileService{
		Namespace: namespace,
		Name:      name,
		Summary:   summary,
		Methods:   make(map[string]*FileServiceMethod),
	}
}

func (p *FileService) AddMethod(m *FileServiceMethod) {
	p.Methods[m.Name] = m
}

func (p *FileService) SetSummary(s string) {
	p.Summary = s
}

func (p *FileService) QryMethod(methodName string) *FileServiceMethod {
	if m, ok := p.Methods[methodName]; ok {
		return m
	}
	return nil
}

func (p *FileService) Margin(other *FileService) {
	if other == nil {
		return
	}
	for k, m := range other.Methods {
		if _, ok := p.Methods[k]; !ok {
			p.Methods[k] = m
		}
	}
}

func (p *FileService) Save(writer xstream.StreamWriter) (err error) {
	if err = writer.WriteInt8(int8(FNT_SERVICE_NODE)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Name); err != nil {
		return err
	}
	if err = writer.WriteStr(p.Summary); err != nil {
		return err
	}
	size := len(p.Methods)
	if size > math.MaxInt32 {
		size = math.MaxInt32
	}
	if err = writer.WriteInt32(int32(size)); err != nil {
		return err
	}
	methodKeys := make([]string, size)
	i := 0
	for k := range p.Methods {
		methodKeys[i] = k
		i++
	}
	sort.Slice(methodKeys, func(i, j int) bool {
		return strings.Compare(methodKeys[i], methodKeys[j]) < 0
	})

	for _, k := range methodKeys {
		m := p.Methods[k]
		if err = m.Save(writer); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileService) Load(reader xstream.StreamReader) (err error) {
	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_SERVICE_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_SERVICE_NODE, nt)
		}
	}
	if p.Name, err = reader.ReadStr(); err != nil {
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
		m := NewFileServiceMethod("", "")
		if err = m.Load(reader); err != nil {
			return err
		}
		p.Methods[m.Name] = m
	}
	return nil
}
