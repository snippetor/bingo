// Code generated by protoc-gen-gogo.
// source: rpc_msg.proto
// DO NOT EDIT!

/*
	Package rpc is a generated protocol buffer package.

	It is generated from these files:
		rpc_msg.proto

	It has these top-level messages:
		RPCHandShake
		RPCMethodCall
		RPCMethodReturn
*/
package rpc

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"
import github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type RPC_MSGID int32

const (
	RPC_MSGID_HANDSHAKE     RPC_MSGID = -256
	RPC_MSGID_CALL          RPC_MSGID = -512
	RPC_MSGID_CALL_NORETURN RPC_MSGID = -513
	RPC_MSGID_RETURN        RPC_MSGID = -514
)

var RPC_MSGID_name = map[int32]string{
	-256: "HANDSHAKE",
	-512: "CALL",
	-513: "CALL_NORETURN",
	-514: "RETURN",
}
var RPC_MSGID_value = map[string]int32{
	"HANDSHAKE":     -256,
	"CALL":          -512,
	"CALL_NORETURN": -513,
	"RETURN":        -514,
}

func (x RPC_MSGID) Enum() *RPC_MSGID {
	p := new(RPC_MSGID)
	*p = x
	return p
}
func (x RPC_MSGID) String() string {
	return proto.EnumName(RPC_MSGID_name, int32(x))
}
func (x *RPC_MSGID) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(RPC_MSGID_value, data, "RPC_MSGID")
	if err != nil {
		return err
	}
	*x = RPC_MSGID(value)
	return nil
}
func (RPC_MSGID) EnumDescriptor() ([]byte, []int) { return fileDescriptorRpcMsg, []int{0} }

type RPCHandShake struct {
	EndName string `protobuf:"bytes,1,req,name=endName" json:"endName"`
}

func (m *RPCHandShake) Reset()                    { *m = RPCHandShake{} }
func (m *RPCHandShake) String() string            { return proto.CompactTextString(m) }
func (*RPCHandShake) ProtoMessage()               {}
func (*RPCHandShake) Descriptor() ([]byte, []int) { return fileDescriptorRpcMsg, []int{0} }

func (m *RPCHandShake) GetEndName() string {
	if m != nil {
		return m.EndName
	}
	return ""
}

type RPCMethodCall struct {
	CallSeq int32             `protobuf:"varint,1,req,name=call_seq,json=callSeq" json:"call_seq"`
	Method  string            `protobuf:"bytes,2,req,name=method" json:"method"`
	Args    map[string]string `protobuf:"bytes,3,rep,name=args" json:"args,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Caller  string            `protobuf:"bytes,4,req,name=caller" json:"caller"`
}

func (m *RPCMethodCall) Reset()                    { *m = RPCMethodCall{} }
func (m *RPCMethodCall) String() string            { return proto.CompactTextString(m) }
func (*RPCMethodCall) ProtoMessage()               {}
func (*RPCMethodCall) Descriptor() ([]byte, []int) { return fileDescriptorRpcMsg, []int{1} }

func (m *RPCMethodCall) GetCallSeq() int32 {
	if m != nil {
		return m.CallSeq
	}
	return 0
}

func (m *RPCMethodCall) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *RPCMethodCall) GetArgs() map[string]string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *RPCMethodCall) GetCaller() string {
	if m != nil {
		return m.Caller
	}
	return ""
}

type RPCMethodReturn struct {
	CallSeq int32             `protobuf:"varint,1,req,name=call_seq,json=callSeq" json:"call_seq"`
	Method  string            `protobuf:"bytes,2,req,name=method" json:"method"`
	Returns map[string]string `protobuf:"bytes,3,rep,name=returns" json:"returns,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *RPCMethodReturn) Reset()                    { *m = RPCMethodReturn{} }
func (m *RPCMethodReturn) String() string            { return proto.CompactTextString(m) }
func (*RPCMethodReturn) ProtoMessage()               {}
func (*RPCMethodReturn) Descriptor() ([]byte, []int) { return fileDescriptorRpcMsg, []int{2} }

func (m *RPCMethodReturn) GetCallSeq() int32 {
	if m != nil {
		return m.CallSeq
	}
	return 0
}

func (m *RPCMethodReturn) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *RPCMethodReturn) GetReturns() map[string]string {
	if m != nil {
		return m.Returns
	}
	return nil
}

func init() {
	proto.RegisterType((*RPCHandShake)(nil), "rpc.RPCHandShake")
	proto.RegisterType((*RPCMethodCall)(nil), "rpc.RPCMethodCall")
	proto.RegisterType((*RPCMethodReturn)(nil), "rpc.RPCMethodReturn")
	proto.RegisterEnum("rpc.RPC_MSGID", RPC_MSGID_name, RPC_MSGID_value)
}
func (m *RPCHandShake) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RPCHandShake) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(len(m.EndName)))
	i += copy(dAtA[i:], m.EndName)
	return i, nil
}

func (m *RPCMethodCall) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RPCMethodCall) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(m.CallSeq))
	dAtA[i] = 0x12
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(len(m.Method)))
	i += copy(dAtA[i:], m.Method)
	if len(m.Args) > 0 {
		for k, _ := range m.Args {
			dAtA[i] = 0x1a
			i++
			v := m.Args[k]
			mapSize := 1 + len(k) + sovRpcMsg(uint64(len(k))) + 1 + len(v) + sovRpcMsg(uint64(len(v)))
			i = encodeVarintRpcMsg(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintRpcMsg(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintRpcMsg(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	dAtA[i] = 0x22
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(len(m.Caller)))
	i += copy(dAtA[i:], m.Caller)
	return i, nil
}

func (m *RPCMethodReturn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RPCMethodReturn) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(m.CallSeq))
	dAtA[i] = 0x12
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(len(m.Method)))
	i += copy(dAtA[i:], m.Method)
	if len(m.Returns) > 0 {
		for k, _ := range m.Returns {
			dAtA[i] = 0x1a
			i++
			v := m.Returns[k]
			mapSize := 1 + len(k) + sovRpcMsg(uint64(len(k))) + 1 + len(v) + sovRpcMsg(uint64(len(v)))
			i = encodeVarintRpcMsg(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintRpcMsg(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintRpcMsg(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	return i, nil
}

func encodeFixed64RpcMsg(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32RpcMsg(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintRpcMsg(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RPCHandShake) Size() (n int) {
	var l int
	_ = l
	l = len(m.EndName)
	n += 1 + l + sovRpcMsg(uint64(l))
	return n
}

func (m *RPCMethodCall) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovRpcMsg(uint64(m.CallSeq))
	l = len(m.Method)
	n += 1 + l + sovRpcMsg(uint64(l))
	if len(m.Args) > 0 {
		for k, v := range m.Args {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovRpcMsg(uint64(len(k))) + 1 + len(v) + sovRpcMsg(uint64(len(v)))
			n += mapEntrySize + 1 + sovRpcMsg(uint64(mapEntrySize))
		}
	}
	l = len(m.Caller)
	n += 1 + l + sovRpcMsg(uint64(l))
	return n
}

func (m *RPCMethodReturn) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovRpcMsg(uint64(m.CallSeq))
	l = len(m.Method)
	n += 1 + l + sovRpcMsg(uint64(l))
	if len(m.Returns) > 0 {
		for k, v := range m.Returns {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovRpcMsg(uint64(len(k))) + 1 + len(v) + sovRpcMsg(uint64(len(v)))
			n += mapEntrySize + 1 + sovRpcMsg(uint64(mapEntrySize))
		}
	}
	return n
}

func sovRpcMsg(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRpcMsg(x uint64) (n int) {
	return sovRpcMsg(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RPCHandShake) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpcMsg
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RPCHandShake: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RPCHandShake: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EndName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000001)
		default:
			iNdEx = preIndex
			skippy, err := skipRpcMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpcMsg
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("endName")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RPCMethodCall) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpcMsg
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RPCMethodCall: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RPCMethodCall: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallSeq", wireType)
			}
			m.CallSeq = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CallSeq |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Method", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Method = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000002)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var keykey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				keykey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			var stringLenmapkey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLenmapkey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLenmapkey := int(stringLenmapkey)
			if intStringLenmapkey < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postStringIndexmapkey := iNdEx + intStringLenmapkey
			if postStringIndexmapkey > l {
				return io.ErrUnexpectedEOF
			}
			mapkey := string(dAtA[iNdEx:postStringIndexmapkey])
			iNdEx = postStringIndexmapkey
			if m.Args == nil {
				m.Args = make(map[string]string)
			}
			if iNdEx < postIndex {
				var valuekey uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRpcMsg
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					valuekey |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				var stringLenmapvalue uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRpcMsg
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLenmapvalue |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLenmapvalue := int(stringLenmapvalue)
				if intStringLenmapvalue < 0 {
					return ErrInvalidLengthRpcMsg
				}
				postStringIndexmapvalue := iNdEx + intStringLenmapvalue
				if postStringIndexmapvalue > l {
					return io.ErrUnexpectedEOF
				}
				mapvalue := string(dAtA[iNdEx:postStringIndexmapvalue])
				iNdEx = postStringIndexmapvalue
				m.Args[mapkey] = mapvalue
			} else {
				var mapvalue string
				m.Args[mapkey] = mapvalue
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Caller", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Caller = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000004)
		default:
			iNdEx = preIndex
			skippy, err := skipRpcMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpcMsg
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("call_seq")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("method")
	}
	if hasFields[0]&uint64(0x00000004) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("caller")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RPCMethodReturn) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpcMsg
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RPCMethodReturn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RPCMethodReturn: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallSeq", wireType)
			}
			m.CallSeq = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CallSeq |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Method", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Method = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000002)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Returns", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var keykey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				keykey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			var stringLenmapkey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLenmapkey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLenmapkey := int(stringLenmapkey)
			if intStringLenmapkey < 0 {
				return ErrInvalidLengthRpcMsg
			}
			postStringIndexmapkey := iNdEx + intStringLenmapkey
			if postStringIndexmapkey > l {
				return io.ErrUnexpectedEOF
			}
			mapkey := string(dAtA[iNdEx:postStringIndexmapkey])
			iNdEx = postStringIndexmapkey
			if m.Returns == nil {
				m.Returns = make(map[string]string)
			}
			if iNdEx < postIndex {
				var valuekey uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRpcMsg
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					valuekey |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				var stringLenmapvalue uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRpcMsg
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLenmapvalue |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLenmapvalue := int(stringLenmapvalue)
				if intStringLenmapvalue < 0 {
					return ErrInvalidLengthRpcMsg
				}
				postStringIndexmapvalue := iNdEx + intStringLenmapvalue
				if postStringIndexmapvalue > l {
					return io.ErrUnexpectedEOF
				}
				mapvalue := string(dAtA[iNdEx:postStringIndexmapvalue])
				iNdEx = postStringIndexmapvalue
				m.Returns[mapkey] = mapvalue
			} else {
				var mapvalue string
				m.Returns[mapkey] = mapvalue
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpcMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpcMsg
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("call_seq")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("method")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRpcMsg(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRpcMsg
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRpcMsg
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthRpcMsg
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRpcMsg
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipRpcMsg(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthRpcMsg = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRpcMsg   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("rpc_msg.proto", fileDescriptorRpcMsg) }

var fileDescriptorRpcMsg = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x8f, 0x4f, 0x8f, 0x9a, 0x40,
	0x18, 0xc6, 0x1d, 0x40, 0x2d, 0x6f, 0x35, 0xa5, 0xd3, 0xc4, 0x10, 0x62, 0xd0, 0x7a, 0x32, 0x3d,
	0x90, 0xa6, 0xa7, 0xa6, 0x3d, 0x21, 0xda, 0xda, 0x54, 0xa9, 0x19, 0xda, 0x33, 0x21, 0x38, 0xc1,
	0x44, 0x44, 0x1c, 0xb0, 0x89, 0x37, 0x3f, 0x42, 0x3f, 0x96, 0xc7, 0x5e, 0x7a, 0x6d, 0x36, 0xfa,
	0x31, 0x76, 0xa3, 0x1b, 0x40, 0x08, 0x5e, 0x77, 0xb9, 0x90, 0xfc, 0x9e, 0x3f, 0xef, 0x3c, 0xd0,
	0x64, 0xa1, 0x6b, 0xaf, 0x22, 0x4f, 0x0b, 0xd9, 0x3a, 0x5e, 0x63, 0x9e, 0x85, 0x6e, 0x4f, 0x83,
	0x06, 0x99, 0x19, 0x63, 0x27, 0x98, 0x5b, 0x0b, 0x67, 0x49, 0xb1, 0x0a, 0x75, 0x1a, 0xcc, 0x4d,
	0x67, 0x45, 0x65, 0xd4, 0xe5, 0xfa, 0xe2, 0x40, 0x38, 0xfc, 0xef, 0x54, 0x48, 0x0e, 0x7b, 0x27,
	0x04, 0x4d, 0x32, 0x33, 0xa6, 0x34, 0x5e, 0xac, 0xe7, 0x86, 0xe3, 0xfb, 0xb8, 0x03, 0x2f, 0x5c,
	0xc7, 0xf7, 0xed, 0x88, 0x6e, 0xd2, 0x48, 0x35, 0x8f, 0x24, 0xd4, 0xa2, 0x1b, 0xdc, 0x86, 0xda,
	0x2a, 0xb5, 0xcb, 0x5c, 0xa9, 0xf1, 0xca, 0xf0, 0x7b, 0x10, 0x1c, 0xe6, 0x45, 0x32, 0xdf, 0xe5,
	0xfb, 0x2f, 0x3f, 0xb4, 0x35, 0x16, 0xba, 0xda, 0xcd, 0x01, 0x4d, 0x67, 0x5e, 0x34, 0x0a, 0x62,
	0xb6, 0x23, 0xa9, 0x33, 0xe9, 0x4b, 0xaa, 0x29, 0x93, 0x85, 0x72, 0x5f, 0xc6, 0x14, 0x03, 0xc4,
	0x22, 0x80, 0x5b, 0xc0, 0x2f, 0xe9, 0x4e, 0x46, 0x5d, 0x54, 0xf8, 0x12, 0x80, 0x15, 0xa8, 0xfe,
	0x76, 0xfc, 0x2d, 0x95, 0xb9, 0x92, 0x92, 0xa1, 0x4f, 0xdc, 0x47, 0xd4, 0xfb, 0x87, 0xe0, 0x55,
	0xf1, 0x08, 0x42, 0xe3, 0x2d, 0x0b, 0x9e, 0xbb, 0xf3, 0x33, 0xd4, 0x59, 0x5a, 0x94, 0x4f, 0x7d,
	0x7b, 0x3b, 0x35, 0xbb, 0xa2, 0x65, 0xbf, 0xeb, 0xde, 0x3c, 0xa1, 0x7c, 0x81, 0x46, 0x59, 0x78,
	0xea, 0xae, 0x77, 0x1e, 0x88, 0x64, 0x66, 0xd8, 0x53, 0xeb, 0xeb, 0xb7, 0x21, 0x6e, 0x81, 0x38,
	0xd6, 0xcd, 0xa1, 0x35, 0xd6, 0xbf, 0x8f, 0xa4, 0xfd, 0xf9, 0x92, 0x7d, 0x08, 0xbf, 0x06, 0xc1,
	0xd0, 0x27, 0x13, 0x69, 0xff, 0x50, 0x20, 0x05, 0x9a, 0x09, 0xb2, 0xcd, 0x1f, 0x64, 0xf4, 0xf3,
	0x17, 0x31, 0xa5, 0xcb, 0x7d, 0xa1, 0xbd, 0x81, 0xda, 0x15, 0x9e, 0x0b, 0x38, 0x90, 0x0e, 0x47,
	0x15, 0xfd, 0x3d, 0xaa, 0xe8, 0xee, 0xa8, 0xa2, 0x3f, 0x27, 0xb5, 0xf2, 0x18, 0x00, 0x00, 0xff,
	0xff, 0x12, 0x79, 0xd9, 0x78, 0x7d, 0x02, 0x00, 0x00,
}
