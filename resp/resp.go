package resp

import (
	"fmt"
	"strings"
)

type responseWrapper struct {
}

func (rw *responseWrapper) SimpleString(str string) []byte {
	return []byte(SimpleStrPrefix + str + EndStr)
}

func (rw *responseWrapper) ErrorString(str string) []byte {
	return []byte( ErrorStrPrefix + str + EndStr)
}

func (rw *responseWrapper) IntegerStr(val int64) []byte {
	return []byte(IntegerPrefix + fmt.Sprintf("%d", val) + EndStr)
}

func (rw *responseWrapper) BulkStr(str string) []byte {
	if len(str) == 0 {
		return []byte(NullBulkStr)
	}
	lenStr := fmt.Sprintf("%d", len(str))
	return []byte(BulkStrPrefix + lenStr + EndStr + str + EndStr)
}

func (rw *responseWrapper) ArrayStr(elements []interface{}) []byte {
	if len(elements) == 0 {
		return []byte(EmptyArray)
	}
	var builder strings.Builder
	builder.WriteString(ArrayPrefix)
	builder.WriteString(fmt.Sprintf("%d", len(elements)))
	builder.WriteString(EndStr)
	for i := range elements {
		switch elements[i].(type) {
		case int64:
			builder.Write(rw.IntegerStr(elements[i].(int64)))
			break
		default:
			builder.Write(rw.BulkStr(elements[i].(string)))
		}
	}
	return []byte(builder.String())
}
