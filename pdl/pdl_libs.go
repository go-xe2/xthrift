/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-12 11:51
* Description:
*****************************************************************/

package pdl

import (
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"regexp"
	"strings"
)

var protoSuffix = []string{".json", ".yaml"}
var protoSuffixCount = len(protoSuffix)

func NamespaceToPath(namespace string) (path string, fileName string) {
	if namespace == "" {
		return "", ""
	}
	items := strings.Split(namespace, ".")
	size := len(items)
	if size == 1 {
		return "", items[0]
	}
	return strings.Join(items[:size-1], xfile.Separator), items[size-1]
}

func makeProtoFileName(workPath, path string, fileName string, ext string) string {
	return xfile.Join(workPath, path, fileName) + ext
}

func GetNamespaceFile(workPath string, namespace string) (string, error) {
	path, file := NamespaceToPath(namespace)
	var fileName = ""
	for i := 0; i < protoSuffixCount; i++ {
		s := makeProtoFileName(workPath, path, file, protoSuffix[i])
		if xfile.Exists(s) {
			fileName = s
			break
		}
	}
	if fileName == "" {
		return "", fmt.Errorf("在工作目录%s中找不到命名空间%s的文件", workPath, namespace)
	}
	return fileName, nil
}

func NamespaceLastName(namespace string) (ownerName string, baseName string) {
	if namespace == "" {
		return "", ""
	}
	items := strings.Split(namespace, ".")
	size := len(items)
	if size == 0 {
		return "", ""
	} else if size == 1 {
		return "", items[0]
	}
	return strings.Join(items[:size-1], "."), items[size-1]
}

var protoListReg = regexp.MustCompile(`^list<(.+)>$`)
var protoSetReg = regexp.MustCompile(`^set<(.+)>$`)
var protoMapReg = regexp.MustCompile(`^map<(.+),(.+)>$`)

func MatchProtoListType(typeName string) (isMatch bool, elemType string) {
	items := protoListReg.FindStringSubmatch(typeName)
	isMatch = len(items) == 2
	if isMatch {
		elemType = items[1]
	}
	return
}

func MatchProtoSetType(typeName string) (isMatch bool, elemType string) {
	items := protoSetReg.FindStringSubmatch(typeName)
	isMatch = len(items) == 2
	if isMatch {
		elemType = items[1]
	}
	return
}

func MatchProtoMapType(typeName string) (isMatch bool, keyType, valType string) {
	items := protoMapReg.FindStringSubmatch(typeName)
	isMatch = len(items) == 3
	if isMatch {
		keyType = items[1]
		valType = items[2]
	}
	return
}
