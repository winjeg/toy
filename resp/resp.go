package resp

import (
	"fmt"
	"strings"
)

type responseWrapper struct {
}

func (rw *responseWrapper) SimpleString(str string) string {
	return SimpleStrPrefix + str + EndStr
}

func (rw *responseWrapper) ErrorString(str string) string {
	return ErrorStrPrefix + str + EndStr
}

func (rw *responseWrapper) IntegerStr(val int64) string {
	return IntegerPrefix + fmt.Sprintf("%d", val) + EndStr
}

func (rw *responseWrapper) BulkStr(str string) string {
	lenStr := fmt.Sprintf("%d", len(str))
	return BulkStrPrefix + lenStr + EndStr + str + EndStr
}

func (rw *responseWrapper) ArrayStr(elements []interface{}) string {
	if len(elements) == 0 {
		return EmptyArray
	}
	var builder strings.Builder
	builder.WriteString(ArrayPrefix)
	builder.WriteString(fmt.Sprintf("%d", len(elements)))
	builder.WriteString(EndStr)
	for i := range elements {
		switch elements[i].(type) {
		case int64:
			builder.WriteString(rw.IntegerStr(elements[i].(int64)))
			break
		default:
			builder.WriteString(rw.BulkStr(elements[i].(string)))
		}
	}
	return builder.String()
}
