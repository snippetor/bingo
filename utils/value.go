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

package utils

type Value struct {
	inner interface{}
}

func (v *Value) Set(i interface{}) {
	v.inner = i
}

func (v *Value) Get() interface{} {
	return v.inner
}

func (v *Value) GetInt() int {
	return v.inner.(int)
}

func (v *Value) GetInt8() int8 {
	return v.inner.(int8)
}

func (v *Value) GetInt16() int16 {
	return v.inner.(int16)
}

func (v *Value) GetInt32() int32 {
	return v.inner.(int32)
}

func (v *Value) GetInt64() int64 {
	return v.inner.(int64)
}

func (v *Value) GetUint() uint {
	return v.inner.(uint)
}

func (v *Value) GetUint8() uint8 {
	return v.inner.(uint8)
}

func (v *Value) GetUint16() uint16 {
	return v.inner.(uint16)
}

func (v *Value) GetUint32() uint32 {
	return v.inner.(uint32)
}

func (v *Value) GetUint64() uint64 {
	return v.inner.(uint64)
}

func (v *Value) GetFloat32() float32 {
	return v.inner.(float32)
}

func (v *Value) GetFloat64() float64 {
	return v.inner.(float64)
}

func (v *Value) GetString() string {
	return v.inner.(string)
}

func (v *Value) GetBool() bool {
	return v.inner.(bool)
}

func (v *Value) GetIntArray() []int {
	return v.inner.([]int)
}

func (v *Value) GetInt8Array() []int8 {
	return v.inner.([]int8)
}

func (v *Value) GetInt16Array() []int16 {
	return v.inner.([]int16)
}

func (v *Value) GetInt32Array() []int32 {
	return v.inner.([]int32)
}

func (v *Value) GetInt64Array() []int64 {
	return v.inner.([]int64)
}

func (v *Value) GetUintArray() []uint {
	return v.inner.([]uint)
}

func (v *Value) GetUint8Array() []uint8 {
	return v.inner.([]uint8)
}

func (v *Value) GetUint16Array() []uint16 {
	return v.inner.([]uint16)
}

func (v *Value) GetUint32Array() []uint32 {
	return v.inner.([]uint32)
}

func (v *Value) GetUint64Array() []uint64 {
	return v.inner.([]uint64)
}

func (v *Value) GetFloat32Array() []float32 {
	return v.inner.([]float32)
}

func (v *Value) GetFloat64Array() []float64 {
	return v.inner.([]float64)
}

func (v *Value) GetStringArray() []string {
	return v.inner.([]string)
}

func (v *Value) GetBoolArray() []bool {
	return v.inner.([]bool)
}

func (v *Value) GetByteArray() []byte {
	return v.inner.([]byte)
}
