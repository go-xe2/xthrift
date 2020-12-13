/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 15:45
* Description:
*****************************************************************/

package pdl

import "io"

type FileExport interface {
	BeginProjectWrite() error
	EndProjectWrite()
	BeginNamespace(name string) error
	EndNamespace(name string)
	BeginFileWrite(ns *FileNamespace, fileName string) (w io.Writer, cxt interface{}, err error)

	WriteNamespace(w io.Writer, cxt interface{}, namespace string) error
	WriteImports(w io.Writer, cxt interface{}, im []string) error
	WriteTypedefs(w io.Writer, cxt interface{}, defs map[string]*FileTypeDef) error
	WriteTypes(w io.Writer, cxt interface{}, types map[string]*FileStruct) error
	WriteServices(w io.Writer, cxt interface{}, ss map[string]*FileService) error
	Flush(w io.Writer, cxt interface{}) error
	EndFileWrite(w io.Writer, ns *FileNamespace, fileName string)
}
