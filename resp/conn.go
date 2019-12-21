// this is only the page of

package resp

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Conn struct {
	c net.Conn
	w *Writer
	r *Reader
}

func NewConn(c net.Conn) *Conn {
	return &Conn{
		c: c,
		w: &Writer{c},
		r: &Reader{conn: c},
	}
}

func (c Conn) Do() {
	log.Printf("clients connected from %s\n", c.c.RemoteAddr().String())
	defer c.c.Close()
	for {
		cml, err := c.r.Read()
		if err != nil {
			errMsg := fmt.Sprintf("-error read conn: %s", err.Error())
			c.w.Write([]byte(errMsg))
			return
		}
		// here to route commands
		args := Parse(cml)
		result := handleArgs(args)
		writeErr := c.w.Write(result)
		if writeErr != nil {
			return
		}
	}
}

func handleArgs(args []string) []byte {
	// must return some kind of right type of the response or the connection will be close by the client
	// if the command is not certain you have to return error or ok, otherwise the client will close the connection.
	switch strings.ToLower(args[0]) {
	case "auth", "set":
		return []byte("+ok\r\n")
	case "get":
		return []byte("*1\r\n$3\r\nabc\r\n")
	default:
		return []byte(fmt.Sprintf("-unknown command %s\r\n", args[0]))
	}
}
