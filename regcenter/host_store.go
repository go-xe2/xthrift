/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-26 17:00
* Description:
*****************************************************************/

package regcenter

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xfileNotify"
	"github.com/go-xe2/x/sync/xsafeMap"
	"github.com/go-xe2/xthrift/pdl"
	"sync/atomic"
)

type THostStoreChangedEventFun func(store *THostStore, fileId int)

type HostStore interface {
	SetOnChanged(fun THostStoreChangedEventFun)
	GetSavePath() string
	GetEnableFileWatch() bool
	Load() error
	AddHostWithProject(proj *pdl.FileProject, host string, port int, ext ...int)
	AddHost(project string, svcFullNae string, host string, port int, ext ...int)
	HasProject(project string) bool
	RemoveHost(host string, port int) error
	RemoveProject(project string) error
	RemoveFile(fileName string)
	Save() error
	AllHosts() map[string][]*THostStoreToken
	FileHosts(fileId int) map[string][]*THostStoreToken
	FileHostsByName(fileName string) map[string][]*THostStoreToken
	AllFileID() map[string]int
	AllFileMd5() map[string]string
	GetSvcHosts(fullSvcName string) []*THostStoreToken
	EnableFileWatch() (err error)
	DisableFileWatch()
	IsFileWatch() bool
}

type THostStore struct {
	// 当前最大文件id
	maxFileId int
	// 文件路径名称与id映射关系
	fileIds         map[string]int
	filesMd5        *xsafeMap.TStrStrMap
	fileWatcher     *xfileNotify.TWatcher
	watchCallback   *xfileNotify.Callback
	nDisableWatch   int32
	watchLayout     int
	savePath        string
	fileExt         string
	items           map[string]map[string]*THostStoreToken
	enableFileWatch bool
	OnChanged       THostStoreChangedEventFun
}

var _ HostStore = (*THostStore)(nil)

func NewHostStore(savePath string, fileExt string, enableFileWatch bool) *THostStore {
	inst := &THostStore{
		savePath:        savePath,
		fileExt:         fileExt,
		watchLayout:     0,
		watchCallback:   nil,
		fileWatcher:     nil,
		enableFileWatch: enableFileWatch,
		filesMd5:        xsafeMap.NewStrStrMap(),
		fileIds:         make(map[string]int),
		items:           make(map[string]map[string]*THostStoreToken),
	}
	atomic.StoreInt32(&inst.nDisableWatch, 1)
	return inst
}

func (p *THostStore) SetOnChanged(fun THostStoreChangedEventFun) {
	p.OnChanged = fun
}

func (p *THostStore) HostFilePath() string {
	fileName := xfile.Join(p.savePath)
	return fileName
}

func (p *THostStore) MakeHostFileName(projectId string) string {
	return xfile.Join(p.savePath, fmt.Sprintf("%s.%s", projectId, p.fileExt))
}

func (p *THostStore) GetSavePath() string {
	return p.savePath
}

func (p *THostStore) GetEnableFileWatch() bool {
	return p.enableFileWatch
}
