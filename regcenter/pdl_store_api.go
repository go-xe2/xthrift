/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-26 17:16
* Description:
*****************************************************************/

package regcenter

import (
	"bytes"
	"fmt"
	"github.com/go-xe2/x/crypto/xmd5"
	"github.com/go-xe2/x/encoding/xbase64"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TPDLStore) doInstall(pdl *pdl.FileProject, md5 string) {
	if pdl == nil {
		return
	}
	p.allProjects.Set(pdl.GetProjectName(), &TPDLProjectInfo{PDL: pdl, MD5: md5})
	if p.handler != nil {
		p.handler.Install(pdl, md5)
	}
}

func (p *TPDLStore) doUninstall(pdl *pdl.FileProject) {
	if pdl == nil {
		return
	}
	p.allProjects.Remove(pdl.GetProjectName())
	if p.handler != nil {
		p.handler.UnInstall(pdl)
	}
}

func (p *TPDLStore) GetProjectByName(projectName string) *TPDLProjectInfo {
	if v := p.allProjects.Get(projectName); v != nil {
		return v.(*TPDLProjectInfo)
	}
	return nil
}

func (p *TPDLStore) AllProject() []*TPDLProjectInfo {
	result := make([]*TPDLProjectInfo, 0)
	keys := p.allProjects.Keys()
	for _, k := range keys {
		v := p.allProjects.Get(k)
		if v != nil {
			result = append(result, v.(*TPDLProjectInfo))
		}
	}
	return result
}

func (p *TPDLStore) loadFile(fileName string, oldMd5 string) (project *pdl.FileProject, md5 string, err error) {
	if !xfile.Exists(fileName) {
		return nil, "", fmt.Errorf("协议文件%s不存在", fileName)
	}
	md5, err = xmd5.EncryptFile(fileName)
	if oldMd5 == md5 {
		// 文件未变动，不需要安装
		return nil, md5, nil
	}
	project = pdl.NewEmptyFileProject()
	if err = project.LoadFromFile(fileName); err != nil {
		return nil, md5, err
	}
	if p.nsContainer.IsInstall(project.GetProjectName()) {
		if err := p.nsContainer.Uninstall(project.GetProjectName()); err != nil {
			return nil, md5, err
		}
		p.doUninstall(project)
	}
	// 安装协议
	xlog.Debug("加载协议:", project.GetProjectName(), ", md5:", md5)
	if err := p.nsContainer.Install(project); err != nil {
		return nil, md5, err
	}
	p.doInstall(project, md5)
	return project, md5, nil
}

func (p *TPDLStore) Load() error {
	if p.savePath == "" {
		return nil
	}
	if !xfile.Exists(p.savePath) {
		return nil
	}
	p.DisableFileWatch()
	defer p.EnableFileWatch()
	files, err := xfile.ScanDir(p.savePath, fmt.Sprintf("*%s", p.fileExt), true)
	if err != nil {
		return err
	}
	p.projectFiles = make(map[string]string)
	p.filesMd5 = make(map[string]string)
	for _, fileName := range files {
		realPath := xfile.RealPath(fileName)
		proj, md5, err := p.loadFile(realPath, "")
		if err != nil {
			xlog.Error(err)
			continue
		}
		p.filesMd5[realPath] = md5
		p.projectFiles[proj.GetProjectName()] = realPath
		//// 加载时不比较md5
		//p.doInstall(proj, md5)
	}
	return nil
}

func (p *TPDLStore) AddProjectFromContent(content []byte) (proj *pdl.FileProject, err error) {
	reader := bytes.NewBuffer(content)
	project := pdl.NewEmptyFileProject()
	if err := project.LoadProject(reader); err != nil {
		return nil, err
	}
	return project, p.AddProject(project)
}

func (p *TPDLStore) AddProjectFromBase64(base64 []byte) (proj *pdl.FileProject, err error) {
	data, err := xbase64.Decode(base64)
	if err != nil {
		return nil, err
	}
	return p.AddProjectFromContent(data)
}

func (p *TPDLStore) AddProject(proj *pdl.FileProject) (err error) {
	p.DisableFileWatch()
	defer p.EnableFileWatch()
	if proj == nil {
		return nil
	}
	// 获取项目md5值
	projFile := ""
	if s, ok := p.projectFiles[proj.GetProjectName()]; ok {
		projFile = s
	}
	md5 := ""
	if projFile != "" {
		if s, ok := p.filesMd5[projFile]; ok {
			md5 = s
		}
	}
	var projMd5 = ""
	if md5 != "" {
		// 获取当前项目的md5值
		buf := bytes.NewBuffer([]byte{})
		if err := proj.SaveProject(buf); err != nil {
			return err
		}
		projMd5, err = xmd5.Encrypt(buf.Bytes())
		if err != nil {
			return err
		}
		xlog.Debug("注册协议项目:"+proj.GetProjectName(), ", 原md5:", md5, ", 当前md5:", projMd5)
		if projMd5 == md5 {
			// 项目已经存在且没有变动
			return nil
		}
		// 删除原协议文件
		projFileName := xfile.Join(p.savePath, fmt.Sprintf("%s%s", proj.GetProjectName(), p.fileExt))
		if xfile.Exists(projFileName) {
			// 删除之前协议文件
			if err := xfile.Remove(projFileName); err != nil {
				return err
			}
		}
	}
	// 文件不存在或内容变动
	projFileName := xfile.Join(p.savePath, fmt.Sprintf("%s%s", proj.GetProjectName(), p.fileExt))
	if !xfile.Exists(p.savePath) {
		if err := xfile.Mkdir(p.savePath); err != nil {
			return err
		}
	}
	if err := proj.SaveToFile(projFileName); err != nil {
		return err
	}
	// 保存协议文件
	p.projectFiles[proj.GetProjectName()] = projFileName
	if projMd5 == "" {
		projMd5, err = xmd5.EncryptFile(projFileName)
		if err != nil {
			return err
		}
	}
	p.filesMd5[projFileName] = projMd5
	// 卸载之前协议
	if p.nsContainer.IsInstall(proj.GetProjectName()) {
		if err := p.nsContainer.Uninstall(proj.GetProjectName()); err != nil {
			return err
		}
		p.doUninstall(proj)
	}
	// 安装协议文件
	err = p.nsContainer.Install(proj)
	if err != nil {
		return err
	}
	p.doInstall(proj, projMd5)
	return nil
}

func (p *TPDLStore) RemoveProject(projName string) error {
	fileName, ok := p.projectFiles[projName]
	if !ok {
		// 未安装协议
		return nil
	}
	if err := p.nsContainer.Uninstall(projName); err != nil {
		return err
	}
	proj := p.GetProjectByName(projName)
	if proj != nil {
		p.doUninstall(proj.PDL)
	}
	projFileName := xfile.Join(p.savePath, fmt.Sprintf("%s%s", projName, p.fileExt))
	if xfile.Exists(projFileName) {
		// 删除之前协议文件
		if err := xfile.Remove(projFileName); err != nil {
			xlog.Debug(err)
		}
	}
	delete(p.projectFiles, projName)
	delete(p.filesMd5, fileName)
	return nil
}

func (p *TPDLStore) GetProjectNamespaces(projName string) map[string]*pdl.TPDLNamespace {
	if projName == "" {
		return nil
	}
	return p.nsContainer.GetProjectNamespaces(projName)
}
