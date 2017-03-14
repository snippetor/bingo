// Code generated by protoc-gen-gogo.
// source: rpc_msg.proto
// DO NOT EDIT!

/*
	Package rpc is a generated protocol buffer package.

	It is generated from these files:
		rpc_msg.proto

	It has these top-level messages:
		MethodCall
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
	RPC_MSGID_CALL RPC_MSGID = -256
)

var RPC_MSGID_name = map[int32]string{
	-256: "CALL",
}
var RPC_MSGID_value = map[string]int32{
	"CALL": -256,
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

type MethodCall struct {
	MethodName string            `protobuf:"bytes,1,req,name=methodName" json:"methodName"`
	Args       map[string]string `protobuf:"bytes,2,rep,name=args" json:"args,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Version    string            `protobuf:"bytes,3,opt,name=version" json:"version"`
}

func (m *MethodCall) Reset()                    { *m = MethodCall{} }
func (m *MethodCall) String() string            { return proto.CompactTextString(m) }
func (*MethodCall) ProtoMessage()               {}
func (*MethodCall) Descriptor() ([]byte, []int) { return fileDescriptorRpcMsg, []int{0} }

func (m *MethodCall) GetMethodName() string {
	if m != nil {
		return m.MethodName
	}
	return ""
}

func (m *MethodCall) GetArgs() map[string]string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *MethodCall) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*MethodCall)(nil), "rpc.MethodCall")
	proto.RegisterEnum("rpc.RPC_MSGID", RPC_MSGID_name, RPC_MSGID_value)
}
func (m *MethodCall) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MethodCall) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(len(m.MethodName)))
	i += copy(dAtA[i:], m.MethodName)
	if len(m.Args) > 0 {
		for k, _ := range m.Args {
			dAtA[i] = 0x12
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
	dAtA[i] = 0x1a
	i++
	i = encodeVarintRpcMsg(dAtA, i, uint64(len(m.Version)))
	i += copy(dAtA[i:], m.Version)
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
func (m *MethodCall) Size() (n int) {
	var l int
	_ = l
	l = len(m.MethodName)
	n += 1 + l + sovRpcMsg(uint64(l))
	if len(m.Args) > 0 {
		for k, v := range m.Args {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovRpcMsg(uint64(len(k))) + 1 + len(v) + sovRpcMsg(uint64(len(v)))
			n += mapEntrySize + 1 + sovRpcMsg(uint64(mapEntrySize))
		}
	}
	l = len(m.Version)
	n += 1 + l + sovRpcMsg(uint64(l))
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
func (m *MethodCall) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MethodCall: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MethodCall: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MethodName", wireType)
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
			m.MethodName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000001)
		case 2:
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
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
			m.Version = string(dAtA[iNdEx:postIndex])
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
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("methodName")
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
	// 225 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x2a, 0x48, 0x8e,
	0xcf, 0x2d, 0x4e, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0x3a,
	0xc6, 0xc8, 0xc5, 0xe5, 0x9b, 0x5a, 0x92, 0x91, 0x9f, 0xe2, 0x9c, 0x98, 0x93, 0x23, 0xa4, 0xc2,
	0xc5, 0x95, 0x0b, 0xe6, 0xf9, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x69, 0x70, 0x3a, 0xb1,
	0x9c, 0xb8, 0x27, 0xcf, 0x10, 0x84, 0x24, 0x2e, 0xa4, 0xcb, 0xc5, 0x92, 0x58, 0x94, 0x5e, 0x2c,
	0xc1, 0xa4, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xa9, 0x57, 0x54, 0x90, 0xac, 0x87, 0x30, 0x44, 0xcf,
	0xb1, 0x28, 0xbd, 0xd8, 0x35, 0xaf, 0xa4, 0xa8, 0x32, 0x08, 0xac, 0x4c, 0x48, 0x8e, 0x8b, 0xbd,
	0x2c, 0xb5, 0xa8, 0x38, 0x33, 0x3f, 0x4f, 0x82, 0x59, 0x81, 0x11, 0x6e, 0x22, 0x4c, 0x50, 0xca,
	0x99, 0x8b, 0x13, 0xae, 0x45, 0x48, 0x8c, 0x8b, 0x39, 0x3b, 0xb5, 0x52, 0x82, 0x11, 0x49, 0x21,
	0x48, 0x40, 0x48, 0x8a, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x34, 0x55, 0x82, 0x09, 0x49, 0x06, 0x22,
	0x64, 0xc5, 0x64, 0xc1, 0xa8, 0x25, 0xc7, 0xc5, 0x19, 0x14, 0xe0, 0x1c, 0xef, 0x1b, 0xec, 0xee,
	0xe9, 0x22, 0x24, 0xc8, 0xc5, 0xe2, 0xec, 0xe8, 0xe3, 0x23, 0xd0, 0xf0, 0xef, 0x3f, 0x04, 0x30,
	0x3a, 0x09, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13,
	0x1e, 0xcb, 0x31, 0x00, 0x02, 0x00, 0x00, 0xff, 0xff, 0x11, 0xc7, 0xa9, 0xe3, 0x0f, 0x01, 0x00,
	0x00,
}
