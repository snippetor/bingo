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
	"sync"
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
	callTasks        *sync.Map //[utils.Identity]*callTask
	isMonitorRunning bool
}

func (w *callSyncWorker) waitingResult(t *callTask) {
	if w.callTasks == nil {
		w.callTasks = &sync.Map{}
	}
	w.callTasks.Store(t.seq, t)
	if !w.isMonitorRunning {
		w.isMonitorRunning = true
		// RPC调用超时检测
		go func() {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case now := <-ticker.C:
					w.callTasks.Range(func(k, v interface{}) bool {
						if v == nil {
							w.callTasks.Delete(k)
						} else {
							task := v.(*callTask)
							if task == nil || task.conn == nil || task.conn.GetState() == net.STATE_CLOSED || task.c == nil { //无效任务，移除
								w.callTasks.Delete(k)
							} else {
								if (now.UnixNano() - task.t) > int64(2*time.Second) { // 两秒超时
									task.c(nil) // 返回nil
									w.callTasks.Delete(k)
									fwlogger.E("-- RPC method callback timeout, %s(%d) %d %s --", task.method, k, task.conn.Identity(), task.conn.Address())
								}
							}
						}
						return true
					})
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
	if t, ok := w.callTasks.Load(callSeq); ok {
		if t != nil {
			t.(*callTask).c(result)
			w.callTasks.Delete(callSeq)
			return
		}
	}
	fwlogger.W("-- RPC method callback ignored, no found call task for %d --", callSeq)
}
