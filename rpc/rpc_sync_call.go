package rpc

import (
	"time"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo"
)

type callTask struct {
	seq        utils.Identity
	method     string
	conn       net.IConn
	msg        net.MessageBody
	c          RPCCallback
	t          int64
	retryTimes int
}

type callSyncWorker struct {
	callTasks        map[utils.Identity]*callTask
	isMonitorRunning bool
}

func (w *callSyncWorker) waitingResult(t *callTask) {
	if w.callTasks == nil {
		w.callTasks = make(map[utils.Identity]*callTask)
	}
	w.callTasks[t.seq] = t
	if !w.isMonitorRunning {
		w.isMonitorRunning = true
		// RPC调用超时检测
		go func() {
			t := time.NewTicker(time.Second)
			for {
				select {
				case now := <-t.C:
					for seq, t := range w.callTasks {
						if t == nil || t.conn == nil || t.conn.GetState() == net.STATE_CLOSED || t.c == nil { //无效任务，移除
							delete(w.callTasks, seq)
						} else {
							if (now.UnixNano() - t.t) > int64(2*time.Second) { // 两秒超时
								bingo.W("-- RPC method callback timeout, retry it %s(%d) %d %s --", t.method, seq, t.conn.Identity(), t.conn.Address())
								if t.retryTimes == 5 { // 重试5次，放弃
									delete(w.callTasks, seq)
									bingo.E("-- RPC method callback recall end, retry time is more than 5 %s(%d) %d %s --", t.method, seq, t.conn.Identity(), t.conn.Address())
								} else {
									t.t = 0
									// recall
									t.conn.Send(net.MessageId(RPC_MSGID_CALL), t.msg)
									t.retryTimes++
								}
							}
						}
					}
				}
			}
		}()
	}
}

func (w *callSyncWorker) receiveResult(callSeq utils.Identity, result *Result) {
	if w.callTasks == nil {
		bingo.E("-- RPC callSyncWorker.callTask not init --")
		return
	}
	if t, ok := w.callTasks[callSeq]; ok {
		t.c(result)
		delete(w.callTasks, callSeq)
	} else {
		bingo.W("-- RPC method callback ignored, no found call task for %d --", callSeq)
	}
}
