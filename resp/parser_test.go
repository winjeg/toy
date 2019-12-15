package resp

import (
	"fmt"
	"testing"
)

const (
	str = "*2\r\n$7\r\nCOMMAND\r\n$3\r\n12313131\r\n"
)

func TestParse(t *testing.T) {
	data := Parse([]byte(str))
	fmt.Println(data)
}
