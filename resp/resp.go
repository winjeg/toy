package resp

import (
	"strconv"
	"strings"
)

type responseWrapper struct {
}

func (rw *responseWrapper) SimpleString(str string) []byte {
	return []byte(SimpleStrPrefix + str + EndStr)
}

func (rw *responseWrapper) ErrorString(str string) []byte {
	return []byte(ErrorStrPrefix + str + EndStr)
}

func (rw *responseWrapper) IntegerStr(val int64) []byte {
	return []byte(IntegerPrefix + strconv.FormatInt(val, 10) + EndStr)
}

func (rw *responseWrapper) BulkStr(raw []byte) []byte {
	if len(raw) == 0 {
		return []byte(NullBulkStr)
	}
	lenStr := strconv.Itoa(len(raw))
	return []byte(BulkStrPrefix + lenStr + EndStr + string(raw) + EndStr)
}

func (rw *responseWrapper) ArrayStr(elements []interface{}) []byte {
	if len(elements) == 0 {
		return []byte(EmptyArray)
	}
	var builder strings.Builder
	builder.WriteString(ArrayPrefix)
	builder.WriteString(strconv.FormatInt(int64(len(elements)), 10))
	builder.WriteString(EndStr)
	for i := range elements {
		switch elements[i].(type) {
		case int64:
			builder.Write(rw.IntegerStr(elements[i].(int64)))
			break
		default:
			builder.Write(rw.BulkStr([]byte(elements[i].(string))))
		}
	}
	return []byte(builder.String())
}
