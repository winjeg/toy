package server

import (
	"fmt"
	"log"
	"net"

	"github.com/winjeg/toy/commands"
	"github.com/winjeg/toy/resp"
)

func Run(store commands.RedisCommands, password string, port int) {
	resp.SetPassword(password)
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
		// key point of choosing which store you want
		c := resp.NewConn(conn, store)
		go c.Do()
	}
}
