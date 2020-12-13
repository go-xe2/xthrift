/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 16:02
* Description:
*****************************************************************/

package pdl

import (
	"github.com/go-xe2/x/os/xfile"
	"io"
	"strings"
)

type TFileExportYaml struct {
	root    string
	mgr     FileIOManager
	fileMap map[*FileNamespace]string
}

var _ FileExport = (*TFileExportYaml)(nil)

func NewFileExportYaml(root string, mgr FileIOManager) FileExport {
	return &TFileExportYaml{
		root:    root,
		mgr:     mgr,
		fileMap: make(map[*FileNamespace]string),
	}
}

func (p *TFileExportYaml) BeginProjectWrite() error {
	return nil
}

func (p *TFileExportYaml) EndProjectWrite() {
}

func (p *TFileExportYaml) BeginNamespace(ns string) error {
	return nil
}

func (p *TFileExportYaml) EndNamespace(ns string) {
}

func (p *TFileExportYaml) BeginFileWrite(ns *FileNamespace, fileName string) (w io.Writer, cxt interface{}, err error) {
	path := strings.Replace(ns.Namespace, ".", xfile.Separator, -1)
	realPath := xfile.Join(p.root, path)
	if !xfile.Exists(realPath) {
		if err := xfile.Mkdir(realPath); err != nil {
			return nil, nil, err
		}
	}
	file := xfile.Join(realPath, fileName+".yaml")
	w, err = p.mgr.Create(ns.Namespace, file)
	if err != nil {
		return nil, nil, err
	}
	p.fileMap[ns] = file
	iw := newYamlWriter()
	if err := iw.WriteBegin(); err != nil {
		return nil, nil, err
	}
	return w, iw, nil
}

func (TFileExportYaml) WriteNamespace(w io.Writer, cxt interface{}, namespace string) error {
	iw := cxt.(*tYamlWriter)
	if err := iw.WriteNamespace(namespace); err != nil {
		return err
	}
	if err := iw.WriteBasicBegin(); err != nil {
		return err
	}
	if err := iw.WriteBasic(ProtoBasicTypes); err != nil {
		return err
	}
	return iw.WriteBasicEnd()
}

func (TFileExportYaml) WriteImports(w io.Writer, cxt interface{}, im []string) error {
	iw := cxt.(*tYamlWriter)
	return iw.WriteImports(im)
}

func (TFileExportYaml) WriteTypedefs(w io.Writer, cxt interface{}, defs map[string]*FileTypeDef) error {
	iw := cxt.(*tYamlWriter)
	if err := iw.WriteTypeDefBegin(); err != nil {
		return err
	}
	if err := iw.WriteTypeDefs(defs); err != nil {
		return err
	}
	return iw.WriteTypeDefEnd()
}

func (TFileExportYaml) WriteTypes(w io.Writer, cxt interface{}, types map[string]*FileStruct) error {
	iw := cxt.(*tYamlWriter)
	if err := iw.WriteTypesBegin(); err != nil {
		return err
	}
	if err := iw.WriteTypes(types); err != nil {
		return err
	}
	return iw.WriteTypesEnd()
}

func (TFileExportYaml) WriteServices(w io.Writer, cxt interface{}, ss map[string]*FileService) error {
	iw := cxt.(*tYamlWriter)
	if err := iw.WriteInterfacesBegin(); err != nil {
		return err
	}
	if err := iw.WriteInterfaces(ss); err != nil {
		return err
	}
	return iw.WriteInterfacesEnd()
}

func (TFileExportYaml) Flush(w io.Writer, cxt interface{}) error {
	iw := cxt.(*tYamlWriter)
	if err := iw.WriteEnd(); err != nil {
		return err
	}
	if _, err := w.Write(iw.Data()); err != nil {
		return err
	}
	return nil
}

func (p *TFileExportYaml) EndFileWrite(w io.Writer, ns *FileNamespace, fileName string) {
	if f, ok := p.fileMap[ns]; ok {
		p.mgr.Close(ns.Namespace, f)
	}
}
