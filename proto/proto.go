package proto

//import (
//	"errors"
//	"strconv"
//)
//
//// 根据默认协议将消息结构体转换为[]byte
//// 注意：此方法要传入指针
//func (c *Codec) Marshal(v interface{}) (net.MessageBody, error) {
//	return codec.marshal(v)
//}
//
//// 尝试使用注册的消息ID和结构体对来解析消息体，如果解析无误将返回对应的消息结构体
//// 注意：此方法返回的是消息结构体的指针，而非值
//func (c *Codec) UnmarshalTo(msgId net.MessageId, data net.MessageBody, collection IProtoCollection) (interface{}, error) {
//	if v, ok := collection.Get(msgId, protoVersion); ok {
//		if err := codec.unmarshal(data, v); err != nil {
//			return nil, err
//		}
//		return v, nil
//	} else if v, ok := collection.GetDefault(msgId); ok {
//		if err := codec.unmarshal(data, v); err != nil {
//			return nil, err
//		}
//		return v, nil
//	}
//	return nil, errors.New("no found message struct for msgid #" + strconv.Itoa(int(msgId)))
//}
//
//// 直接解析成结构体
//func (c *Codec) Unmarshal(data net.MessageBody, v interface{}) error {
//	if err := codec.unmarshal(data, v); err != nil {
//		return err
//	}
//	return nil
//}