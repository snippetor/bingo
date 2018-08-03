package utils

import (
	"testing"
	"fmt"
)

func TestGenIdentity(t *testing.T) {
	i := NewIdentifier(8)

	//for k := 0; k < 100000; k++ {
	//	fmt.Println(i.GenIdentity())
	//}

	for k := 0; k < 10000; k++ {
		go func() {
			fmt.Println(i.GenIdentity())
		}()
	}
}

func BenchmarkGenIdentity(b *testing.B) {
	b.StopTimer()
	id := NewIdentifier(8)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		id.GenIdentity()
	}
}
