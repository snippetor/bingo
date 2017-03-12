package protocol

import "testing"

func BenchmarkGoGoProtobufMarshal(b *testing.B) {
	b.StopTimer()
	p := protocolProtoBuf{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.marshal(&Person{1, "carl"})
	}
}

func BenchmarkGoGoProtobufUnmarshal(b *testing.B) {
	b.StopTimer()
	p := protocolProtoBuf{}
	bytes, _ := p.marshal(&Person{1, "carl"})
	p1 := Person{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.unmarshal(bytes, &p1)
	}
}

