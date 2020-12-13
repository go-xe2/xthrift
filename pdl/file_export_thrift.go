/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 16:03
* Description:
*****************************************************************/

package pdl

import (
	"errors"
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"io"
	"strings"
)

type tThriftWriteCxt struct {
	namespace string
	w         io.Writer
	buf       *tThriftWriter
	saveFile  string
	srcFile   string
	status    int
	imports   map[string]bool
}

func newThriftWriteCxt(namespace string, w io.Writer, lang string, saveFile string) *tThriftWriteCxt {
	return &tThriftWriteCxt{
		namespace: namespace,
		w:         w,
		buf:       newThriftWriter(namespace, lang),
		imports:   make(map[string]bool),
		status:    0,
		saveFile:  saveFile,
	}
}

type TFileExportThrift struct {
	root  string
	mgr   FileIOManager
	files map[string]*tThriftWriteCxt
	lang  string
}

var _ FileExport = (*TFileExportThrift)(nil)

func NewFileExportThrift(root string, mgr FileIOManager, lang string) FileExport {
	return &TFileExportThrift{
		root:  root,
		mgr:   mgr,
		lang:  lang,
		files: make(map[string]*tThriftWriteCxt),
	}
}

func (p *TFileExportThrift) BeginProjectWrite() error {
	return nil
}

func (p *TFileExportThrift) EndProjectWrite() {
}

func (p *TFileExportThrift) BeginNamespace(ns string) error {
	if !xfile.Exists(p.root) {
		if err := xfile.Mkdir(p.root); err != nil {
			return err
		}
	}
	file := xfile.Join(p.root, ns+".thrift")
	w, err := p.mgr.Create(ns, file)
	if err != nil {
		return err
	}
	p.files[ns] = newThriftWriteCxt(ns, w, p.lang, file)
	return nil
}

func (p *TFileExportThrift) EndNamespace(ns string) {
	cxt := p.files[ns]
	if cxt == nil {
		return
	}
	// 写入包名
	szLang := p.lang
	if szLang == "" {
		szLang = "go"
	}
	langs := strings.Split(szLang, ",")
	for _, lang := range langs {
		_, err := cxt.w.Write([]byte(fmt.Sprintf("namespace %s %s\n", lang, ns)))
		if err != nil {
			panic(err)
		}
	}

	// 写入引用
	for k := range cxt.imports {
		if k == ns {
			continue
		}
		_, err := cxt.w.Write([]byte(fmt.Sprintf("include '%s'\n", k+".thrift")))
		if err != nil {
			panic(err)
		}
	}
	// 写缓存文件
	_, err := cxt.w.Write(cxt.buf.Data())
	if err != nil {
		panic(err)
	}
	p.mgr.Close(ns, cxt.saveFile)
	delete(p.files, ns)
}

func (p *TFileExportThrift) BeginFileWrite(ns *FileNamespace, fileName string) (w io.Writer, cxt interface{}, err error) {
	writerCxt := p.files[ns.Namespace]
	if writerCxt == nil {
		return nil, nil, errors.New("未调用BeginNamespace方法初始化")
	}
	writerCxt.buf.WriteComment(fmt.Sprintf("文件:%s开始", fileName))
	return writerCxt.w, writerCxt, nil
}

func (p TFileExportThrift) WriteNamespace(w io.Writer, cxt interface{}, namespace string) error {
	return nil
}

func (p *TFileExportThrift) WriteImports(w io.Writer, cxt interface{}, im []string) error {
	writeCxt := cxt.(*tThriftWriteCxt)
	// 收集空间中引用到其他空间的文件
	for _, s := range im {
		if _, ok := writeCxt.imports[s]; !ok {
			writeCxt.imports[s] = true
		}
	}
	return nil
}

func (p *TFileExportThrift) WriteTypedefs(w io.Writer, cxt interface{}, defs map[string]*FileTypeDef) error {
	wCxt := cxt.(*tThriftWriteCxt)
	iw := wCxt.buf
	if err := iw.WriteTypeDefBegin(); err != nil {
		return err
	}
	if err := iw.WriteTypeDefs(defs); err != nil {
		return err
	}
	return iw.WriteTypeDefEnd()
}

func (p *TFileExportThrift) WriteTypes(w io.Writer, cxt interface{}, types map[string]*FileStruct) error {
	wCxt := cxt.(*tThriftWriteCxt)
	iw := wCxt.buf
	if err := iw.WriteTypesBegin(); err != nil {
		return err
	}
	if err := iw.WriteTypes(types); err != nil {
		return err
	}
	return iw.WriteTypesEnd()
}

func (p *TFileExportThrift) WriteServices(w io.Writer, cxt interface{}, ss map[string]*FileService) error {
	wCxt := cxt.(*tThriftWriteCxt)
	iw := wCxt.buf
	if err := iw.WriteInterfacesBegin(); err != nil {
		return err
	}
	if err := iw.WriteInterfaces(ss); err != nil {
		return err
	}
	return iw.WriteInterfacesEnd()
}

func (p *TFileExportThrift) Flush(w io.Writer, cxt interface{}) error {
	wCxt := cxt.(*tThriftWriteCxt)
	iw := wCxt.buf
	return iw.WriteEnd()
}

func (p *TFileExportThrift) EndFileWrite(w io.Writer, ns *FileNamespace, fileName string) {
	wCxt := p.files[ns.Namespace]
	iw := wCxt.buf
	iw.WriteComment(fmt.Sprintf("文件%s结束", fileName))
}
