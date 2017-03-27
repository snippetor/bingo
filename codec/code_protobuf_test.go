package codec

import "testing"

func BenchmarkGoGoProtobufMarshal(b *testing.B) {
	b.StopTimer()
	p := protoBuf{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.Marshal(&Person{1, "carl"})
	}
}

func BenchmarkGoGoProtobufUnmarshal(b *testing.B) {
	b.StopTimer()
	p := protoBuf{}
	bytes, _ := p.Marshal(&Person{1, "carl"})
	p1 := Person{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.Unmarshal(bytes, &p1)
	}
}

