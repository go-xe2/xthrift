/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-26 17:30
* Description:
*****************************************************************/

package regcenter

import (
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xfileNotify"
	"github.com/go-xe2/x/os/xlog"
	"sync/atomic"
)

func (p *THostStore) initFileWatch() (err error) {
	p.fileWatcher, err = xfileNotify.New()
	if err != nil {
		return err
	}
	// 初始化的
	p.watchLayout = -1
	return nil
}

func (p *THostStore) EnableFileWatch() (err error) {
	if !p.enableFileWatch {
		return nil
	}
	if p.fileWatcher == nil {
		if err := p.initFileWatch(); err != nil {
			return err
		}
	}
	path := p.HostFilePath()
	if !xfile.Exists(path) {
		return nil
	}
	p.watchLayout++
	if n := atomic.LoadInt32(&p.nDisableWatch); n == 0 {
		return nil
	}
	if p.watchCallback != nil {
		return nil
	}
	p.watchCallback, err = p.fileWatcher.Add(path, p.fileWatchCallback)
	if err != nil {
		p.watchCallback = nil
		return err
	}
	atomic.StoreInt32(&p.nDisableWatch, 0)
	return nil
}

func (p *THostStore) IsFileWatch() bool {
	return p.watchCallback != nil
}

func (p *THostStore) DisableFileWatch() {
	if !p.enableFileWatch {
		return
	}
	p.watchLayout--
	if p.watchLayout > 0 || p.watchCallback == nil {
		return
	}
	if n := atomic.LoadInt32(&p.nDisableWatch); n > 0 {
		return
	}
	p.fileWatcher.RemoveCallback(p.watchCallback.Id)
	p.watchCallback = nil
	atomic.StoreInt32(&p.nDisableWatch, 1)
}

func (p *THostStore) fileWatchCallback(event *xfileNotify.TEvent) {
	if n := atomic.LoadInt32(&p.nDisableWatch); n > 0 {
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
	baseName := xfile.Basename(filePath)
	md5 := ""
	if s := p.filesMd5.Get(baseName); s != "" {
		md5 = s
	}
	if event.IsCreate() || event.IsWrite() {
		// 文件变动
		fileId := 0
		if n, ok := p.fileIds[baseName]; ok {
			fileId = n
		} else {
			p.maxFileId++
			fileId = p.maxFileId
		}
		if s, err := p.loadFile(filePath, md5, fileId); err != nil {
			xlog.Error(err)
		} else {
			if s != md5 {
				p.filesMd5.Set(baseName, s)
				if p.OnChanged != nil {
					p.OnChanged(p, fileId)
				}
			}
		}
	} else if event.IsRemove() {
		// 删除文件
		p.RemoveFile(filePath)
		if p.OnChanged != nil {
			p.OnChanged(p, -1)
		}
	}
}
