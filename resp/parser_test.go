package resp

import (
	"testing"
)

const (
	str = "*2\r\n$7\r\nCOMMAND\r\n$3\r\n12313131\r\n"
)

func TestParse(t *testing.T) {
	data := Parse([]byte(str))
	if len(data) <= 0 {
		t.FailNow()
	}
}

func BenchmarkParse(b *testing.B) {
	data := []byte(str)
	for i := 0; i < b.N; i++ {
		Parse(data)
	}
}
