// this is only the page of

package resp

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/winjeg/toy/commands"
)

type Conn struct {
	c      net.Conn
	w      *Writer
	r      *Reader
	authOk bool
	lock   sync.Mutex
	store  commands.RedisCommands
}

func NewConn(c net.Conn, store commands.RedisCommands) *Conn {
	return &Conn{
		c:     c,
		w:     &Writer{c},
		r:     &Reader{conn: c},
		store: store, // TODO change this to use whatever store, or implementation you want
	}
}

func (c *Conn) Do() {
	log.Printf("clients connected from %s\n", c.c.RemoteAddr().String())
	defer func() {
		fmt.Println("closing connection with defer ")
		c.c.Close()
	}()
	for {
		cml, err := c.r.Read()
		if err != nil {
			return
		}
		// here to route commands
		args := Parse(cml)
		result := c.handleArgs(args)
		writeErr := c.w.Write(result)
		if writeErr != nil {
			return
		}
	}
}

var (
	wrapper        = new(responseWrapper)
	serverPassword = "foobar"
	passLock       sync.Mutex
	okResp         = []byte("+OK\r\n")
	pongResp       = []byte("+PONG\r\n")
)

func SetPassword(password string) {
	passLock.Lock()
	serverPassword = password
	passLock.Unlock()
}

func (c *Conn) handleArgs(cmds []RedisCommand) []byte {
	if len(cmds) == 0 {
		return []byte{0}
	}
	// must return some kind of right type of the response or the connection will be close by the client
	// if the command is not certain you have to return error or ok, otherwise the client will close the connection.
	result := make([]byte, 0, 64)
	for i := range cmds {
		cmd := strings.ToLower(cmds[i].Cmd)
		switch cmd {
		case "auth":
			authArgs := cmds[i].Args
			if len(authArgs) <= 0 {
				result = append(result, wrapper.ErrorString("wrong arguments")...)
				continue
			}
			if strings.EqualFold(serverPassword, authArgs[0]) {
				c.lock.Lock()
				c.authOk = true
				c.lock.Unlock()
				result = append(result, okResp...)
				continue
			} else {
				result = append(result, wrapper.ErrorString("auth failed!")...)
				continue
			}
		case "get":
			if !c.authOk {
				result = append(result, wrapper.ErrorString("Authentication required")...)
				continue
			}
			if len(cmds[i].Args) < 1 {
				result = append(result, wrapper.ErrorString("wrong arguments")...)
				continue
			}
			result = append(result, wrapper.BulkStr(c.store.Get(cmds[i].Args[0]))...)
		case "ping":
			if !c.authOk {
				result = append(result, wrapper.ErrorString("Authentication required")...)
				continue
			}
			result = append(result, pongResp...)
		case "set":
			if !c.authOk {
				result = append(result, wrapper.ErrorString("Authentication required")...)
				continue
			}
			if len(cmds[i].Args) < 2 {
				result = append(result, wrapper.ErrorString("wrong arguments")...)
				continue
			}
			c.store.Set(cmds[i].Args[0], []byte(cmds[i].Args[1]))
			result = append(result, wrapper.SimpleString("ok")...)
		case "config":
			if !c.authOk {
				result = append(result, wrapper.ErrorString("Authentication required")...)
				continue
			}
			if len(cmds[i].Args) < 2 {
				result = append(result, wrapper.ErrorString("wrong arguments")...)
				continue
			}
			if strings.EqualFold(string(cmds[i].Args[0]), "get") {
				result = append(result, okResp...)
				continue
			}
			result = append(result, wrapper.ErrorString(fmt.Sprintf("unknown command config %s", string(cmds[i].Args[1])))...)
		default:
			result = append(result, wrapper.ErrorString("unknown command!")...)
		}
	}
	return result
}
