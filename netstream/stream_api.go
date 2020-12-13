/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-05 10:20
* Description:
*****************************************************************/

package netstream

import (
	"fmt"
	"github.com/go-xe2/x/core/logger"
	"github.com/go-xe2/x/encoding/xhashGen"
	"github.com/go-xe2/x/os/xlog"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

var globalMU sync.Mutex

func MakeConnId() string {
	//globalMU.Lock()
	//defer globalMU.Unlock()
	//n := xrand.N(1000, 9999)
	//nano := time.Now().UnixNano()
	//s := fmt.Sprintf("%d%d", nano, n)
	//if len(s) > 8 {
	//	return s[8:]
	//}
	s := bson.NewObjectId().Hex()
	return s
}

func MakeRequestId() string {
	return bson.NewObjectId().Hex()
}

func MakeSendSeqId() int64 {
	s := bson.NewObjectId()
	seqId := xhashGen.BKDRHash64([]byte(s))
	return int64(seqId)
}

func Log(tag string, lg logger.ILogger, level logger.LogLevel, args ...interface{}) {
	if lg == nil {
		return
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(tag, ", logger api error[", time.Now().String(), "]:", e)
			}
		}()
		switch level {
		case logger.LEVEL_DEBU:
			lg.Debug(tag, args...)
			return
		case logger.LEVEL_INFO:
			lg.Info(tag, args...)
			return
		case logger.LEVEL_WARN:
			lg.Warning(tag, args...)
			return
		case logger.LEVEL_ERRO:
			lg.Error(tag, args...)
			return
		case logger.LEVEL_NOTI:
			lg.Notice(tag, args...)
			return
		case logger.LEVEL_CRIT:
			lg.Critical(tag, args...)
			return
		case logger.LEVEL_DEV:
			lg.Print(tag, uint8(level), args...)
			return
		case xlog.LEVEL_PROD:
			lg.Print(tag, uint8(level), args...)
			return
		}
	}()
}
