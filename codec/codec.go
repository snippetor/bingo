package codec

type CodecType byte

const (
	Json     CodecType = iota
	Protobuf
)

func NewCodec(t CodecType) ICodec {
	switch t {
	case Json:
		return ICodec(&json{})
	case Protobuf:
		return ICodec(&protoBuf{})
	}
	return nil
}
