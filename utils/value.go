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

type Value interface {
	Set(i interface{})
	Get() interface{}
	GetInt() int
	GetInt8() int8
	GetInt16() int16
	GetInt32() int32
	GetInt64() int64
	GetUint() uint
	GetUint8() uint8
	GetUint16() uint16
	GetUint32() uint32
	GetUint64() uint64
	GetFloat32() float32
	GetFloat64() float64
	GetString() string
	GetBool() bool
	GetIntArray() []int
	GetInt8Array() []int8
	GetInt16Array() []int16
	GetInt32Array() []int32
	GetInt64Array() []int64
	GetUintArray() []uint
	GetUint8Array() []uint8
	GetUint16Array() []uint16
	GetUint32Array() []uint32
	GetUint64Array() []uint64
	GetFloat32Array() []float32
	GetFloat64Array() []float64
	GetStringArray() []string
	GetBoolArray() []bool
	GetByteArray() []byte
}

type value struct {
	inner interface{}
}

func NewValue(v interface{}) Value {
	return &value{v}
}

func (v *value) Set(i interface{}) {
	v.inner = i
}

func (v *value) Get() interface{} {
	return v.inner
}

func (v *value) GetInt() int {
	return v.inner.(int)
}

func (v *value) GetInt8() int8 {
	return v.inner.(int8)
}

func (v *value) GetInt16() int16 {
	return v.inner.(int16)
}

func (v *value) GetInt32() int32 {
	return v.inner.(int32)
}

func (v *value) GetInt64() int64 {
	return v.inner.(int64)
}

func (v *value) GetUint() uint {
	return v.inner.(uint)
}

func (v *value) GetUint8() uint8 {
	return v.inner.(uint8)
}

func (v *value) GetUint16() uint16 {
	return v.inner.(uint16)
}

func (v *value) GetUint32() uint32 {
	return v.inner.(uint32)
}

func (v *value) GetUint64() uint64 {
	return v.inner.(uint64)
}

func (v *value) GetFloat32() float32 {
	return v.inner.(float32)
}

func (v *value) GetFloat64() float64 {
	return v.inner.(float64)
}

func (v *value) GetString() string {
	return v.inner.(string)
}

func (v *value) GetBool() bool {
	return v.inner.(bool)
}

func (v *value) GetIntArray() []int {
	return v.inner.([]int)
}

func (v *value) GetInt8Array() []int8 {
	return v.inner.([]int8)
}

func (v *value) GetInt16Array() []int16 {
	return v.inner.([]int16)
}

func (v *value) GetInt32Array() []int32 {
	return v.inner.([]int32)
}

func (v *value) GetInt64Array() []int64 {
	return v.inner.([]int64)
}

func (v *value) GetUintArray() []uint {
	return v.inner.([]uint)
}

func (v *value) GetUint8Array() []uint8 {
	return v.inner.([]uint8)
}

func (v *value) GetUint16Array() []uint16 {
	return v.inner.([]uint16)
}

func (v *value) GetUint32Array() []uint32 {
	return v.inner.([]uint32)
}

func (v *value) GetUint64Array() []uint64 {
	return v.inner.([]uint64)
}

func (v *value) GetFloat32Array() []float32 {
	return v.inner.([]float32)
}

func (v *value) GetFloat64Array() []float64 {
	return v.inner.([]float64)
}

func (v *value) GetStringArray() []string {
	return v.inner.([]string)
}

func (v *value) GetBoolArray() []bool {
	return v.inner.([]bool)
}

func (v *value) GetByteArray() []byte {
	return v.inner.([]byte)
}
