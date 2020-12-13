/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 14:54
* Description:
*****************************************************************/

package netstream

import (
	"errors"
	"fmt"
	"github.com/go-xe2/x/sync/xsafeMap"
	"sync/atomic"
	"time"
)

type tSendReply struct {
	seqId   int64
	buf     []byte
	ready   chan byte
	err     error
	timeout time.Duration
	reqTime time.Time
	mgr     *tSendReplies
}

func newSendReply(mgr *tSendReplies, seqId int64, timeout time.Duration) *tSendReply {
	return &tSendReply{
		mgr:     mgr,
		seqId:   seqId,
		ready:   make(chan byte, 1),
		buf:     nil,
		reqTime: time.Now(),
		timeout: timeout,
	}
}

func (p *tSendReply) GetResult() (data []byte, err error) {
	go func() {
		// 设置超时时间
		time.Sleep(p.timeout)
		select {
		case <-p.ready:
			// 已经处理好，不再处理
			return
		default:
			p.buf = nil
			p.err = errors.New("访问超时")
			p.mgr.removeWait(p.seqId)
			close(p.ready)
		}
	}()
	select {
	case <-p.ready:
		return p.buf, p.err
	}
}

func (p *tSendReply) Reply(data []byte, err error) {
	p.err = err
	p.buf = data
	select {
	case <-p.ready:
		// 数据已经返回，不需要处理
		return
	default:
		p.buf = data
		p.err = err
		p.mgr.removeWait(p.seqId)
		close(p.ready)
	}
}

type tSendReplies struct {
	isCheckRun int32
	waits      *xsafeMap.TIntAnyMap
	closed     int32
}

var sendReplies = &tSendReplies{
	isCheckRun: 0,
	waits:      xsafeMap.NewIntAnyMap(),
	closed:     0,
}

func (p *tSendReplies) Open() {
	atomic.StoreInt32(&p.closed, 0)
}

func (p *tSendReplies) GetResult(seqId int64, timeout time.Duration) (data []byte, err error) {
	if atomic.LoadInt32(&p.closed) != 0 {
		return nil, errors.New("已关闭连接")
	}
	if p.waits.Contains(int(seqId)) {
		return nil, errors.New(fmt.Sprintf("请求序列%d重复", int(seqId)))
	}
	reply := newSendReply(p, seqId, timeout)
	p.waits.Set(int(seqId), reply)
	return reply.GetResult()
}

func (p *tSendReplies) Reply(seqId int64, data []byte, err error) {
	if tmp := p.waits.Get(int(seqId)); tmp != nil {
		r := tmp.(*tSendReply)
		r.Reply(data, err)
	}
}

func (p *tSendReplies) Close() {
	atomic.StoreInt32(&p.closed, 1)
}

func (p *tSendReplies) removeWait(seqId int64) {
	p.waits.Remove(int(seqId))
}
