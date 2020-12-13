package pdl

import (
	"errors"
	"fmt"
	"github.com/go-xe2/x/crypto/xmd5"
	"github.com/go-xe2/x/encoding/xyaml"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/os/xstream"
	"github.com/go-xe2/x/type/xstring"
	"io"
	"os"
	"sort"
	"strings"
)

const ProjectFileFlag int32 = 0x0A35AF
const ProjectVer int16 = 0x00f1

type FileProject struct {
	projectId string
	// 项目名称
	ProjectName string                    `json:"projectName"`
	Path        string                    `json:"path"`
	Namespaces  map[string]*FileNamespace `json:"namespaces"`
	errors      []error
	services    map[string]*FileService
}

func NewEmptyFileProject() *FileProject {
	return &FileProject{
		Path:       "",
		Namespaces: make(map[string]*FileNamespace),
	}
}

func NewFileProject(path string) (*FileProject, error) {
	if !xfile.Exists(path) {
		if err := xfile.Mkdir(path); err != nil {
			return nil, fmt.Errorf("路径%s不存在, 创建失败:%s", path, err)
		}
	}
	return &FileProject{
		Path:       path,
		Namespaces: make(map[string]*FileNamespace),
	}, nil
}

func (p *FileProject) GetServices() map[string]*FileService {
	if p.services == nil {
		p.services = make(map[string]*FileService)
		p.findAllServices()
	}
	return p.services
}

func (p *FileProject) findAllServices() {
	for _, ns := range p.Namespaces {
		for _, file := range ns.Files {
			for _, svc := range file.Services {
				keyName := fmt.Sprintf("%s.%s", ns.Namespace, xstring.LcFirst(svc.Name))
				p.services[keyName] = svc
			}
		}
	}
}

func (p *FileProject) GetNamespace(namespace string) *FileNamespace {
	if v, ok := p.Namespaces[namespace]; ok {
		return v
	}
	return nil
}

func (p *FileProject) GetOrSetNamespace(namespace string) *FileNamespace {
	if v, ok := p.Namespaces[namespace]; ok {
		return v
	} else {
		v = NewFileNamespace(p, namespace)
		p.Namespaces[namespace] = v
		return v
	}
}

// 加载项目
func (p *FileProject) Load() error {
	items, err := xfile.ScanDir(p.Path, "*.yaml, *.json", true)
	if err != nil {
		return err
	}
	p.errors = make([]error, 0)
	for _, f := range items {
		file := NewFileData(p)
		if err := file.OpenFile(f); err != nil {
			p.errors = append(p.errors, fmt.Errorf("加载文件%s错误:%s", f, err))
			xlog.Warning("加载文件:", f, ", 出错:", err)
			continue
		}
		ns := p.GetOrSetNamespace(file.Namespace)
		file.bindNamespace(ns)
		ns.AddFile(file)
	}
	projFile := xfile.Join(p.Path, "pdl.proj")
	if xfile.Exists(projFile) {
		fileData := xfile.GetBinContents(projFile)
		if fileData != nil {
			var mp = map[string]interface{}{}
			err := xyaml.DecodeTo(fileData, &mp)
			if err != nil {
				return err
			}
			if s, ok := mp["name"].(string); ok {
				p.ProjectName = s
			}
		}
	}
	if err := p.loadProjectInfo(); err != nil {
		return err
	}
	return p.Check()
}

func (p *FileProject) loadProjectInfo() error {
	projFile := xfile.Join(p.Path, "pdl.proj")
	if xfile.Exists(projFile) {
		fileData := xfile.GetBinContents(projFile)
		if fileData != nil {
			var mp = map[string]interface{}{}
			err := xyaml.DecodeTo(fileData, &mp)
			if err != nil {
				return err
			}
			if s, ok := mp["name"].(string); ok {
				p.ProjectName = s
			}
		}
	}
	return nil
}

func (p *FileProject) saveProjectInfo() error {
	projFile := xfile.Join(p.Path, "pdl.proj")
	info := map[string]interface{}{
		"name": p.ProjectName,
	}
	data, err := xyaml.Encode(info)
	if err != nil {
		return err
	}
	return xfile.PutBinContents(projFile, data)
}

func (p *FileProject) Errors() []error {
	return p.errors
}

func (p *FileProject) QryNamespace(name string) *FileNamespace {
	if ns, ok := p.Namespaces[name]; ok {
		return ns
	}
	return nil
}

func (p *FileProject) AllNamespaces() []string {
	result := make([]string, len(p.Namespaces))
	i := 0
	for k := range p.Namespaces {
		result[i] = k
		i++
	}
	return result
}

func (p *FileProject) SetProjectName(name string) {
	p.ProjectName = name
	if err := p.saveProjectInfo(); err != nil {
		xlog.Error("save projectInfo error:", err)
	}
}

func (p *FileProject) GetProjectName() string {
	return p.ProjectName
}

func (p *FileProject) GetProjectId() string {
	if p.projectId == "" {
		svcs := p.AllServices()
		svcNames := make([]string, len(svcs))
		i := 0
		for k := range svcs {
			svcNames[i] = k
		}
		sort.Slice(svcNames, func(i, j int) bool {
			return strings.Compare(svcNames[i], svcNames[j]) < 0
		})
		hashCode, _ := xmd5.Encrypt(strings.Join(svcNames, ""))
		p.projectId = hashCode
	}
	return p.projectId
}

func (p *FileProject) AllServices() map[string]*FileService {
	result := make(map[string]*FileService)
	for _, ns := range p.Namespaces {
		for _, f := range ns.Files {
			for _, svc := range f.Services {
				result[fmt.Sprintf("%s.%s", ns.Namespace, xstring.LcFirst(svc.Name))] = svc
			}
		}
	}
	return result
}

func (p *FileProject) QryTypedef(defName string) (*FileNamespace, *FileTypeDef) {
	if strings.Index(defName, ".") > 0 {
		ns, dtName := NamespaceLastName(defName)
		return p.QryTypeDefByNS(ns, dtName)
	}
	for _, ns := range p.Namespaces {
		if dt := ns.QryTypedef(defName); dt != nil {
			return ns, dt
		}
	}
	return nil, nil
}

func (p *FileProject) QryTypeDefByNS(namespace string, defName string) (*FileNamespace, *FileTypeDef) {
	if ns, ok := p.Namespaces[namespace]; ok {
		if v := ns.QryTypedef(defName); v != nil {
			return ns, v
		}
		return ns, nil
	}
	return nil, nil
}

func (p *FileProject) QryType(typeName string) (*FileNamespace, *FileStruct) {
	if strings.Index(typeName, ".") > 0 {
		ns, dt := NamespaceLastName(typeName)
		return p.QryTypeByNS(ns, dt)
	}
	for _, ns := range p.Namespaces {
		if dt := ns.QryType(typeName); dt != nil {
			return ns, dt
		}
	}
	return nil, nil
}

func (p *FileProject) QryTypeByNS(namespace string, typeName string) (*FileNamespace, *FileStruct) {
	if ns, ok := p.Namespaces[namespace]; ok {
		if dt := ns.QryType(typeName); dt != nil {
			return ns, dt
		}
		return ns, nil
	}
	return nil, nil
}

func (p *FileProject) QryService(svcName string) (*FileNamespace, *FileService) {
	if strings.Index(svcName, ".") > 0 {
		ns, dtn := NamespaceLastName(svcName)
		return p.QryServiceByNS(ns, dtn)
	}
	for _, ns := range p.Namespaces {
		if dt := ns.QryService(svcName); dt != nil {
			return ns, dt
		}
	}
	return nil, nil
}

func (p *FileProject) QryServiceByNS(namespace string, svcName string) (*FileNamespace, *FileService) {
	if ns, ok := p.Namespaces[namespace]; ok {
		if dt := ns.QryService(svcName); dt != nil {
			return ns, dt
		}
		return ns, nil
	}
	return nil, nil
}

func (p *FileProject) QryMethod(svcName string, methodName string) (*FileNamespace, *FileService, *FileServiceMethod) {
	ns, svc := p.QryService(svcName)
	if svc != nil {
		if dt := svc.QryMethod(methodName); dt != nil {
			return ns, svc, dt
		}
		return ns, svc, nil
	}
	return ns, nil, nil
}

func (p *FileProject) QryMethodByNS(namespace string, svcName string, methodName string) (*FileNamespace, *FileService, *FileServiceMethod) {
	ns, svc := p.QryServiceByNS(namespace, svcName)
	if svc != nil {
		if dt := svc.QryMethod(methodName); dt != nil {
			return ns, svc, dt
		}
		return ns, svc, nil
	}
	return ns, nil, nil
}

// 合并项目包
func (p *FileProject) Margin(other *FileProject) {
	if other == nil {
		return
	}
	for k, ns := range other.Namespaces {
		self, ok := p.Namespaces[k]
		if !ok {
			p.Namespaces[k] = ns
		} else {
			self.Margin(ns)
		}
	}
}

// 保存项目
func (p *FileProject) SaveProject(w io.Writer) (err error) {
	writer := xstream.NewLeStreamWriter(w)
	if err = writer.WriteInt32(ProjectFileFlag); err != nil {
		return err
	}
	if err = writer.WriteInt16(ProjectVer); err != nil {
		return err
	}
	if err = writer.WriteInt8(int8(FNT_PROJECT_NODE)); err != nil {
		return err
	}
	if err = writer.WriteStr(p.ProjectName); err != nil {
		return err
	}
	size := len(p.Namespaces)
	if err = writer.WriteInt64(int64(size)); err != nil {
		return err
	}
	keys := make([]string, size)
	i := 0
	for k := range p.Namespaces {
		keys[i] = k
		i++
	}
	// 排序
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	for _, k := range keys {
		ns := p.Namespaces[k]
		if err = ns.Save(writer); err != nil {
			return err
		}
	}
	return nil
}

func (p *FileProject) LoadFromFile(fileName string) error {
	if !xfile.Exists(fileName) {
		return fmt.Errorf("协议文件%s不存在", fileName)
	}
	file, err := xfile.OpenWithFlag(fileName, os.O_RDONLY)
	if err != nil {
		return err
	}
	defer file.Close()
	return p.LoadProject(file)
}

func (p *FileProject) SaveToFile(fileName string) error {
	if xfile.Exists(fileName) {
		return errors.New("协议项目已经存在")
	}
	file, err := xfile.OpenWithFlag(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := p.SaveProject(file); err != nil {
		return err
	}
	return nil
}

func (p *FileProject) LoadProject(r io.Reader) (err error) {
	reader := xstream.NewLeStreamReader(r)
	if fg, err := reader.ReadInt32(); err != nil {
		return err
	} else {
		if fg != ProjectFileFlag {
			return errors.New("不是有效的协议项目文件")
		}
	}
	if v, err := reader.ReadInt16(); err != nil {
		return err
	} else {
		if v != ProjectVer {
			return errors.New("协议版本不正确")
		}
	}

	if n, err := reader.ReadInt8(); err != nil {
		return err
	} else {
		if nt := TFileNodeType(n); nt != FNT_PROJECT_NODE {
			return fmt.Errorf("数据损坏名不是有效的协议项目文件,期望%s类型，实际为%s类型", FNT_PROJECT_NODE, nt)
		}
	}
	// 读取项目名称
	if p.ProjectName, err = reader.ReadStr(); err != nil {
		return err
	}
	n, e := reader.ReadInt64()
	if e != nil {
		return e
	}
	if p.Namespaces == nil {
		p.Namespaces = make(map[string]*FileNamespace)
	}
	size := int(n)
	for i := 0; i < size; i++ {
		ns := NewFileNamespace(p, "")
		if err = ns.Load(reader); err != nil {
			return err
		}
		p.Namespaces[ns.Namespace] = ns
	}
	return nil
}

func (p *FileProject) Check() error {
	for _, ns := range p.Namespaces {
		if err := ns.Check(); err != nil {
			return fmt.Errorf("检查空间:%s错误:%s", ns.Namespace, err)
		}
	}
	return nil
}

func (p *FileProject) Export(exp FileExport) error {
	if err := exp.BeginProjectWrite(); err != nil {
		return err
	}
	defer exp.EndProjectWrite()
	for _, ns := range p.Namespaces {
		if err := ns.Export(exp); err != nil {
			return err
		}
	}
	return nil
}
