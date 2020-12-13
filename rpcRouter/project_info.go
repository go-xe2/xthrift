package rpcRouter

import (
	"github.com/go-xe2/xthrift/pdl"
	"time"
)

type tProjectInfo struct {
	md5     string
	upgrade int64
	project *pdl.FileProject
}

func newProjectInfo(md5 string, proj *pdl.FileProject) *tProjectInfo {
	return &tProjectInfo{
		md5:     md5,
		project: proj,
		upgrade: time.Now().Unix(),
	}
}
