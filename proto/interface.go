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

package proto

import "github.com/snippetor/bingo/net"

// 协议集合，用于存储消息ID和结构体的对应关系集合
type IProtoCollection interface {
	PutDefault(id net.MessageId, v interface{})
	Put(id net.MessageId, v interface{}, protoVersion string)
	GetDefault(id net.MessageId) (interface{}, bool)
	Get(id net.MessageId, protoVersion string) (interface{}, bool)
	RemoveDefault(id net.MessageId)
	Remove(id net.MessageId, protoVersion string)
	Clear()
	Size() int
	//Dump()
}
