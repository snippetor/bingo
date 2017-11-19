// Copyright 2017 bingo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"time"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/log/fwlogger"
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
								if t, ok := w.callTasks[seq]; ok {
									t.c(nil) // 返回nil
									delete(w.callTasks, seq)
								} else {
									fwlogger.W("-- RPC method callback ignored, no found call task for %d --", seq)
								}
								fwlogger.E("-- RPC method callback timeout, %s(%d) %d %s --", t.method, seq, t.conn.Identity(), t.conn.Address())
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
		fwlogger.E("-- RPC callSyncWorker.callTask not init --")
		return
	}
	if t, ok := w.callTasks[callSeq]; ok {
		t.c(result)
		delete(w.callTasks, callSeq)
	} else {
		fwlogger.W("-- RPC method callback ignored, no found call task for %d --", callSeq)
	}
}
