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
	"github.com/go-xe2/x/crypto/xmd5"
	"github.com/go-xe2/x/encoding/xparser"
	"github.com/go-xe2/x/encoding/xyaml"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
	"io/ioutil"
	"os"
)

func (p *THostStore) loadFile(fileName string, md5 string, fileId int) (fileMd5 string, err error) {
	xlog.Debug("load host file:", fileName)
	if !xfile.Exists(fileName) {
		return "", nil
	}
	if curMd5, err := xmd5.EncryptFile(fileName); err != nil {
		return "", err
	} else {
		if md5 != curMd5 {
			fileMd5 = curMd5
		} else {
			// 文件没有变动，不处理
			return md5, nil
		}
	}
	file, err := xfile.OpenWithFlag(fileName, os.O_RDONLY)
	if err != nil {
		return fileMd5, err
	}
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return fileMd5, err
	}
	defer file.Close()

	parser, err := xparser.LoadContent(fileData)
	if err != nil {
		return fileMd5, err
	}
	if fileId >= 0 {
		p.removeNodeByFileId(fileId)
	}
	mp := parser.ToMap()
	for k, v := range mp {
		items, ok := v.([]interface{})
		if !ok {
			continue
		}
		var hostItems map[string]*THostStoreToken
		if tmp, ok := p.items[k]; ok {
			hostItems = tmp
		} else {
			hostItems = make(map[string]*THostStoreToken)
		}
		for _, item := range items {
			host := ""
			port := 0
			project := ""
			itemMp, itemOk := item.(map[string]interface{})
			if !itemOk {
				continue
			}
			if s, ok := itemMp["host"].(string); ok {
				host = s
			}
			if n, ok := itemMp["port"]; ok {
				port = t.Int(n)
			}
			if s, ok := itemMp["project"].(string); ok {
				project = s
			}
			if host == "" || port == 0 {
				continue
			}
			hostItems[fmt.Sprintf("%s:%d", host, port)] = &THostStoreToken{Project: project, Host: host, Port: port, fileId: fileId, Ext: 0}
		}
		if len(hostItems) > 0 {
			p.items[k] = hostItems
		} else {
			delete(p.items, k)
		}
	}
	return fileMd5, nil
}

func (p *THostStore) loadHostFileMd5() error {
	md5File := xfile.Join(p.savePath, "host_file.md5")
	if !xfile.Exists(md5File) {
		return nil
	}
	fileData := xfile.GetBinContents(md5File)
	p.filesMd5.Clear()
	var mp map[string]interface{}
	err := xyaml.DecodeTo(fileData, &mp)
	if err != nil {
		xlog.Debug("load file md5 error:", err)
		return err
	}
	for k, v := range mp {
		if s, ok := v.(string); ok {
			p.filesMd5.Set(k, s)
		}
	}
	return nil
}

func (p *THostStore) saveHostFileMd5() error {
	md5Bytes, err := xyaml.Encode(p.filesMd5)
	if err != nil {
		return err
	}
	md5File := xfile.Join(p.savePath, "host_file.md5")
	file, err := xfile.OpenWithFlag(md5File, os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		xlog.Debug("open md5 file error:", err)
		return err
	}
	defer file.Close()
	if _, err := file.Write(md5Bytes); err != nil {
		xlog.Debug("save file md5 error:", err)
		return err
	}
	return nil
}

func (p *THostStore) Load() error {
	if !xfile.Exists(p.HostFilePath()) {
		if err := xfile.Mkdir(p.HostFilePath()); err != nil {
			return err
		}
	}
	files, err := xfile.ScanDir(p.HostFilePath(), fmt.Sprintf("*%s", p.fileExt), true)
	if err != nil {
		xlog.Debug("scanDir error:", err)
		return err
	}

	// 清空原数据
	p.items = make(map[string]map[string]*THostStoreToken)
	p.fileIds = make(map[string]int)
	p.maxFileId = 0
	// 加载md5文件
	//if err := p.loadHostFileMd5(); err != nil {
	//	xlog.Debug("loadFileMd5 error:", err)
	//	return err
	//}
	p.filesMd5.Clear()

	for _, fileName := range files {
		// 生成文件id

		baseName := xfile.Basename(fileName)
		fileId := p.maxFileId
		p.maxFileId++
		p.fileIds[baseName] = fileId
		var md5 = ""
		if p.filesMd5.Contains(baseName) {
			md5 = p.filesMd5.Get(baseName)
		}
		if curMd5, err := p.loadFile(fileName, md5, fileId); err != nil {
			xlog.Error(err)
		} else {
			if curMd5 != md5 {
				p.filesMd5.Set(baseName, curMd5)
			}
		}
	}
	return p.saveHostFileMd5()
}

func (p *THostStore) AddHostWithProject(proj *pdl.FileProject, host string, port int, ext ...int) {
	if proj == nil {
		return
	}
	services := proj.AllServices()
	for k := range services {
		p.AddHost(proj.GetProjectName(), k, host, port, ext...)
	}
}

func (p *THostStore) AddHost(project string, svcFullName string, host string, port int, ext ...int) {
	k := fmt.Sprintf("%s:%d", host, port)
	node, ok := p.items[svcFullName]
	fileName := fmt.Sprintf("%s%s", project, p.fileExt)
	fileId := 0
	if n, ok := p.fileIds[fileName]; ok {
		fileId = n
	} else {
		p.maxFileId++
		fileId = p.maxFileId
		p.fileIds[fileName] = fileId

	}
	nExt := 0
	if len(ext) > 0 {
		nExt = ext[0]
	}
	if ok {
		if _, ok := node[k]; !ok {
			node[fmt.Sprintf("%s:%d", host, port)] = &THostStoreToken{Project: project, Host: host, Port: port, isSaved: false, fileId: fileId, Ext: nExt}
		}
	} else {
		node = make(map[string]*THostStoreToken)
		node[k] = &THostStoreToken{Project: project, Host: host, Port: port, isSaved: false, fileId: fileId, Ext: nExt}
	}
	p.items[svcFullName] = node
}

func (p *THostStore) HasProject(project string) bool {
	isRef := false
	for _, host := range p.items {
		if isRef {
			break
		}
		for _, node := range host {
			if node.Project == project {
				isRef = true
				break
			}
		}
	}
	return isRef
}

// 注销提供服务地址
func (p *THostStore) RemoveHost(host string, port int) error {
	k := fmt.Sprintf("%s:%d", host, port)
	isModify := false
	project := ""
	for _, host := range p.items {
		for curK, node := range host {
			if curK == k {
				if project == "" {
					project = node.Project
				}
				delete(host, curK)
				isModify = true
			}
		}
	}
	if isModify {
		if err := p.Save(); err != nil {
			return err
		}
	}
	return nil
}

func (p *THostStore) RemoveProject(project string) error {
	isModify := false
	for _, host := range p.items {
		for k, node := range host {
			if node.Project == project {
				delete(host, k)
				isModify = true
			}
		}
	}
	if isModify {
		if err := p.Save(); err != nil {
			return err
		}
		if p.savePath == "" {
			return nil
		}
	}
	return nil
}

func (p *THostStore) removeNodeByFileId(fileId int) {
	for _, host := range p.items {
		for k, node := range host {
			if node.fileId == fileId {
				delete(host, k)
			}
		}
	}
}

func (p *THostStore) RemoveFile(fileName string) {
	p.DisableFileWatch()
	defer func() {
		if err := p.EnableFileWatch(); err != nil {
			xlog.Error(err)
		}
	}()
	baseName := xfile.Basename(fileName)
	fileId, ok := p.fileIds[baseName]
	if !ok {
		fileId = 0
	}
	isModify := false
	for _, host := range p.items {
		for k, node := range host {
			if node.fileId == fileId {
				delete(host, k)
				isModify = true
			}
		}
	}
	if isModify {
		p.filesMd5.Remove(baseName)
		delete(p.fileIds, baseName)
		pathName := xfile.Join(p.savePath, fileName)
		if xfile.Exists(pathName) {
			if err := xfile.Remove(pathName); err != nil {
				xlog.Error(err)
			}
		}
	}
}

func (p *THostStore) saveFile(fileName string, items map[string][]*THostStoreToken) error {
	p.DisableFileWatch()
	defer func() {
		if err := p.EnableFileWatch(); err != nil {
			xlog.Error(err)
		}
	}()
	if !xfile.Exists(p.savePath) {
		if err := xfile.Mkdir(p.savePath); err != nil {
			return err
		}
	}
	pathName := xfile.Join(p.savePath, fileName)
	bytes, err := xyaml.Encode(items)
	if err != nil {
		return err
	}
	file, err := xfile.OpenWithFlag(pathName, os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(bytes); err != nil {
		return err
	}
	md5, err := xmd5.Encrypt(bytes)
	if err != nil {
		return err
	}
	p.filesMd5.Set(fileName, md5)
	return nil
}

func (p *THostStore) Save() error {
	// host按文件保存
	files := make(map[int]map[string][]*THostStoreToken)

	for svcName, v := range p.items {
		for _, node := range v {
			file, fileOk := files[node.fileId]
			if !fileOk {
				file = make(map[string][]*THostStoreToken)
				files[node.fileId] = file
			}
			fileItems, svcOk := file[svcName]
			if !svcOk {
				fileItems = make([]*THostStoreToken, 0)
			}
			fileItems = append(fileItems, node)
			file[svcName] = fileItems
		}
	}
	fileNames := make(map[int]string)
	for k, v := range p.fileIds {
		fileNames[v] = k
	}

	var fileName = ""
	p.DisableFileWatch()
	defer func() {
		if err := p.EnableFileWatch(); err != nil {
			xlog.Error(err)
		}
	}()
	// 保存文件
	for fileId, items := range files {
		fileName = ""
		if s, ok := fileNames[fileId]; !ok {
			continue
		} else {
			fileName = s
		}
		if err := p.saveFile(fileName, items); err != nil {
			xlog.Error(err)
		}
	}
	return nil
}

func (p *THostStore) AllHosts() map[string][]*THostStoreToken {
	result := make(map[string][]*THostStoreToken)
	for svc, items := range p.items {
		arr := make([]*THostStoreToken, len(items))
		i := 0
		for _, node := range items {
			arr[i] = node
			i++
		}
		result[svc] = arr
	}
	return result
}

func (p *THostStore) FileHosts(fileId int) map[string][]*THostStoreToken {
	result := make(map[string][]*THostStoreToken)
	for svc, items := range p.items {
		for _, node := range items {
			if node.fileId == fileId {
				arr, ok := result[svc]
				if !ok {
					arr = make([]*THostStoreToken, 0)
				}
				arr = append(arr, node)
				result[svc] = arr
			}
		}
	}
	return result
}

func (p *THostStore) FileHostsByName(fileName string) map[string][]*THostStoreToken {
	if n, ok := p.fileIds[fileName]; ok {
		return p.FileHosts(n)
	}
	return nil
}

func (p *THostStore) AllFileID() map[string]int {
	return p.fileIds
}

func (p *THostStore) AllFileMd5() map[string]string {
	return p.filesMd5.Map()
}

func (p *THostStore) GetSvcHosts(fullSvcName string) []*THostStoreToken {
	hosts := make([]*THostStoreToken, 0)
	hostNode, ok := p.items[fullSvcName]
	if !ok {
		return hosts
	}
	for _, n := range hostNode {
		hosts = append(hosts, n)
	}
	return hosts
}
