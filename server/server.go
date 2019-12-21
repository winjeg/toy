package server

import (
	"fmt"
	"log"
	"net"

	"github.com/winjeg/toy/resp"
)

const defaultSize = 32

func Run(port int) {
	server, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		c := resp.NewConn(conn)
		go c.Do()
	}
}
