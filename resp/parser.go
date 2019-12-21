package resp

import (
	"log"
	"strconv"
)

// parse raw input being sent by redis clients to an string slice
// it shall be empty if there is something wrong
func Parse(raw []byte) []string {
	if len(raw) <= 0 {
		return nil
	}
	if raw[0] != ArrayVal {
		return nil
	}
	var arrLenData []byte
	for i := range raw {
		if raw[i] == '\r' {
			arrLenData = raw[1:i]
			break
		}
	}
	arrLen, err := strconv.ParseInt(string(arrLenData), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	var result = make([]string, 0, arrLen)
	var startIdx int
	var bulkLen int64
	for i := 0; i < len(raw); i++ {
		if raw[i] == BulkStrVal && bulkLen == 0 {
			startIdx = i
		}
		if startIdx > 0 && raw[i] == '\r' {
			bulkLen, _ = strconv.ParseInt(string(raw[startIdx+1:i]), 10, 64)
			startIdx = 0
			i += 2
		}
		if bulkLen > 0 {
			// i start
			result = append(result, string(raw[i:i+int(bulkLen)]))
			if len(result) >= int(arrLen) {
				break
			}
			i += int(bulkLen)
			bulkLen = 0
		}
	}
	return result
}

/**
all commands are send via  array
pipeline 就是多条命令一起发， 其实没什么的
multi先发, exec 也是多次命令一起发
*/

///

type Client struct {
	Args [][]byte
}

func (c *Client) FromArgs(rawArray [][]byte) {

}
