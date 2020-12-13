/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-26 17:16
* Description:
*****************************************************************/

package regcenter

import (
	"github.com/go-xe2/x/os/xfileNotify"
	"github.com/go-xe2/x/sync/xsafeMap"
	"github.com/go-xe2/xthrift/pdl"
	"sync/atomic"
)

type TPDLStoreChangedEventFun func(fileName, projName string, namespaces map[string]*pdl.TPDLNamespace)

type PDLStoreHandler interface {
	Install(proj *pdl.FileProject, md5 string)
	UnInstall(proj *pdl.FileProject)
}

type PDLStore interface {
	SetOnChanged(fun TPDLStoreChangedEventFun)
	AddProjectFromContent(content []byte) (proj *pdl.FileProject, err error)
	AddProjectFromBase64(base64 []byte) (proj *pdl.FileProject, err error)
	AddProject(proj *pdl.FileProject) (err error)
	RemoveProject(projName string) error
	Load() error
	EnableFileWatch() error
	DisableFileWatch()
	IsFileWatch() bool
	SetHandler(handler PDLStoreHandler)
	GetProjectByName(projectName string) *TPDLProjectInfo
	AllProject() []*TPDLProjectInfo
	pdl.PDLQuery
}

type TPDLProjectInfo struct {
	PDL *pdl.FileProject
	MD5 string
}

type TPDLStore struct {
	savePath string
	// 协议文件后缀
	fileExt     string
	nsContainer *pdl.TPDLContainer
	// 协议文件md5值
	filesMd5 map[string]string
	// 项目名称与文件路径映射
	projectFiles  map[string]string
	watcher       *xfileNotify.TWatcher
	watchCallback *xfileNotify.Callback
	enableWatch   bool
	onChanged     TPDLStoreChangedEventFun
	handler       PDLStoreHandler
	allProjects   *xsafeMap.TStrAnyMap
	nDisableWatch int32
}

func NewPDLStore(savePath string, fileExt string, enableWatch bool) *TPDLStore {
	inst := &TPDLStore{
		savePath:      savePath,
		fileExt:       fileExt,
		filesMd5:      make(map[string]string),
		nsContainer:   pdl.NewPDLContainer(),
		projectFiles:  make(map[string]string),
		enableWatch:   enableWatch,
		watcher:       nil,
		watchCallback: nil,
		allProjects:   xsafeMap.NewStrAnyMap(),
	}
	atomic.StoreInt32(&inst.nDisableWatch, 1)
	return inst
}

func (p *TPDLStore) SetOnChanged(fun TPDLStoreChangedEventFun) {
	p.onChanged = fun
}

func (p *TPDLStore) GetSavePath() string {
	return p.savePath
}

func (p *TPDLStore) GetEnableFileWatch() bool {
	return p.enableWatch
}

func (p *TPDLStore) SetHandler(handler PDLStoreHandler) {
	p.handler = handler
}
