/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-27 16:22
* Description:
*****************************************************************/

package regcenter

import (
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xfileNotify"
	"github.com/go-xe2/x/os/xlog"
	"sync/atomic"
)

func (p *TPDLStore) initFileWatch() error {
	if p.watcher != nil {
		return nil
	}
	if w, err := xfileNotify.New(); err != nil {
		return err
	} else {
		p.watcher = w
	}
	return nil
}

func (p *TPDLStore) EnableFileWatch() error {
	if !p.enableWatch {
		return nil
	}
	if p.watcher == nil {
		if err := p.initFileWatch(); err != nil {
			return err
		}
	}
	if n := atomic.LoadInt32(&p.nDisableWatch); n == 0 {
		// 已经启动
		return nil
	}
	if p.watchCallback != nil {
		return nil
	}
	if !xfile.Exists(p.savePath) {
		if err := xfile.Mkdir(p.savePath); err != nil {
			return err
		}
	}
	if c, err := p.watcher.Add(p.savePath, p.fileWatchCallback); err != nil {
		return err
	} else {
		p.watchCallback = c
		atomic.StoreInt32(&p.nDisableWatch, 0)
	}
	return nil
}

func (p *TPDLStore) IsFileWatch() bool {
	return p.watchCallback != nil
}

func (p *TPDLStore) DisableFileWatch() {
	if !p.enableWatch {
		return
	}
	if n := atomic.LoadInt32(&p.nDisableWatch); n > 0 {
		return
	}
	if p.watchCallback == nil {
		return
	}
	p.watcher.RemoveCallback(p.watchCallback.Id)
	p.watchCallback = nil
	atomic.StoreInt32(&p.nDisableWatch, 1)
}

func (p *TPDLStore) fileWatchCallback(event *xfileNotify.TEvent) {
	defer func() {
		if e := recover(); e != nil {
			xlog.Debug(e)
		}
	}()
	if n := atomic.LoadInt32(&p.nDisableWatch); n != 0 {
		return
	}
	filePath := event.Path
	if xfile.IsDir(filePath) {
		return
	}
	fileExt := xfile.Ext(filePath)
	if fileExt != p.fileExt {
		return
	}
	if event.IsWrite() || event.IsCreate() {
		md5, ok := p.filesMd5[filePath]
		if !ok {
			md5 = ""
		}
		if !xfile.IsReadable(filePath) {
			return
		}
		if proj, curMd5, err := p.loadFile(filePath, md5); err != nil {
			xlog.Error(err)
		} else {
			if proj == nil {
				return
			}
			p.projectFiles[proj.GetProjectName()] = filePath
			p.filesMd5[filePath] = curMd5
			if p.onChanged != nil {
				nsMap := p.GetProjectNamespaces(proj.GetProjectName())
				if nsMap != nil {
					p.onChanged(filePath, proj.GetProjectName(), nsMap)
				}
			}
		}
	} else if event.IsRemove() {
		// 删除文件，卸载协议
		projName := ""
		for proj, fileName := range p.projectFiles {
			if fileName == filePath {
				projName = proj
				break
			}
		}
		if projName == "" {
			return
		}
		if err := p.RemoveProject(projName); err != nil {
			xlog.Errorf("remove project:%s error:%s", projName, err.Error())
		} else {
			if p.onChanged != nil {
				p.onChanged(filePath, projName, nil)
			}
		}
	}
}
