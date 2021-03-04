package server

import (
	"fmt"
	"log"
	"net"

	"github.com/winjeg/toy/commands"
	"github.com/winjeg/toy/conn"
)

func Serve(store commands.RedisCommands, password string, port int) {
	conn.SetPassword(password)
	server, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	log.Println("server started.")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer server.Close()
	for {
		netConn, err := server.Accept()
		if err != nil {
			continue
		}
		// key point of choosing which store you want
		c := conn.NewConn(netConn, store)
		go c.HandleCommands()
	}
}
