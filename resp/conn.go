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

func (c *Conn) handleArgs(args []string) []byte {
	if len(args) == 0 {
		return []byte{0}
	}
	// must return some kind of right type of the response or the connection will be close by the client
	// if the command is not certain you have to return error or ok, otherwise the client will close the connection.
	if strings.EqualFold(strings.ToLower(args[0]), "auth") {
		if len(args) != 2 || !strings.EqualFold(args[1], serverPassword) {
			return wrapper.ErrorString("auth failed!")
		}
		c.lock.Lock()
		c.authOk = true
		c.lock.Unlock()
		return okResp
	}
/*	if !c.authOk {
		return wrapper.ErrorString("NOAUTH Authentication required")
	}*/
	switch strings.ToLower(args[0]) {
	case "auth":
		if len(args) != 2 || !strings.EqualFold(args[1], "foobar") {
			return wrapper.ErrorString("auth failed!")
		}
		c.lock.Lock()
		c.authOk = true
		c.lock.Unlock()
		return wrapper.SimpleString("ok")
	case "get":
		return wrapper.BulkStr(c.store.Get(args[1]))
	case "ping":
		return pongResp
	case "set":
		if len(args) != 3 {
			return wrapper.ErrorString("wrong arguments")
		}
		c.store.Set(args[1], []byte(args[2]))
		return wrapper.SimpleString("ok")
	default:
		return wrapper.ErrorString("unknown command!")
	}
}
