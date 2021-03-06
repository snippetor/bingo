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

package codec

import "github.com/snippetor/bingo/net"

type CodecType byte

const (
	Json     CodecType = iota
	Protobuf
)

// 数据传输协议，即包体格式
type Codec interface {
	Marshal(interface{}) (net.MessageBody, error)
	Unmarshal(net.MessageBody, interface{}) error
	Type() int
}

func NewCodec(t CodecType) Codec {
	switch t {
	case Json:
		return &json{}
	case Protobuf:
		return &protoBuf{}
	}
	return nil
}

var (
	JsonCodec     Codec
	ProtobufCodec Codec
)

func init() {
	JsonCodec = NewCodec(Json)
	ProtobufCodec = NewCodec(Protobuf)
}
