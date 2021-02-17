package parser

import (
	"log"
	"strconv"

	"github.com/winjeg/toy/resp"
)

type RedisCommand struct {
	Cmd  string
	Args []string
}

// parse raw input being sent by redis clients to an string slice
// it shall be empty if there is something wrong
// redis benchmark 会一块传很多命令过来， 因此要把所有命令解析出来且正确返回才行，否则会卡主
// 也就是必须要实现pipeline
func Parse(raw []byte) []RedisCommand {
	arrLen := getFirstArrLen(raw)
	if arrLen == 0 {
		return nil
	}

	var startIdx int
	var bulkLen int64
	result := make([]string, 0, arrLen)
	cmds := make([]RedisCommand, 0, 4)

	for i := 0; i < len(raw); i++ {
		// 获取bulk的长度
		if raw[i] == resp.BulkStrVal && bulkLen == 0 {
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
				cmd := RedisCommand{
					Cmd:  result[0],
					Args: result[1:],
				}
				cmds = append(cmds, cmd)
				i += int(bulkLen) + 2
				bulkLen = 0
				arrLen = getFirstArrLen(raw[i:])
				if arrLen == 0 {
					break
				}
				result = make([]string, 0, arrLen)
				continue
			}
			i += int(bulkLen)
			bulkLen = 0
		}
	}
	return cmds
}

func getFirstArrLen(raw []byte) int64 {
	if len(raw) <= 0 {
		return 0
	}
	if raw[0] != resp.ArrayVal {
		return 0
	}

	// PING 这个command是错的
	// 网络连接应该可以边读边写双工
	// 获取第一个长度

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
		return 0
	}
	return arrLen
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
