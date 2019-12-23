// this is only the page of

package resp

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Conn struct {
	c      net.Conn
	w      *Writer
	r      *Reader
	authOk bool
	lock   sync.Mutex
}

func NewConn(c net.Conn) *Conn {
	return &Conn{
		c: c,
		w: &Writer{c},
		r: &Reader{conn: c},
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

func (c *Conn) handleArgs(args []string) []byte {
	// must return some kind of right type of the response or the connection will be close by the client
	// if the command is not certain you have to return error or ok, otherwise the client will close the connection.
	wrapper := new(responseWrapper)
	switch strings.ToLower(args[0]) {
	case "auth":
		// TODO choose this password from any other client
		if len(args) != 2 || !strings.EqualFold(args[1], "foobar") {
			return wrapper.ErrorString("auth failed!")
		}
		c.lock.Lock()
		c.authOk = true
		c.lock.Unlock()
		return wrapper.SimpleString("ok")
	default:
		if c.authOk {
			return wrapper.SimpleString("ok!")
		} else {
			return wrapper.ErrorString("need auth first!")
		}
	}
}
