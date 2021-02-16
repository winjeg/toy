package parser

import (
	"testing"
)

const (
	str = "*2\r\n$4\r\nAUTH\r\n$6\r\n123456\r\n*3\r\n$3\r\nSET\r\n$16\r\nkey:000059531779\r\n$3\r\nxxx\r\n"
)

func TestParse(t *testing.T) {
	data := Parse([]byte(str))
	if len(data) != 2 {
		t.FailNow()
	}
}

func BenchmarkParse(b *testing.B) {
	data := []byte(str)
	for i := 0; i < b.N; i++ {
		Parse(data)
	}
}
