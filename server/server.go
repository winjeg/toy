package server

import (
	"fmt"
	"log"
	"net"
)

func ServeRedis() {
	server, err := net.Listen("tcp", "127.0.0.1:6379")
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
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	fmt.Printf("clients connected from %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	for {
		cml := make([]byte, 0, 16)
		buf := make([]byte, 8)
		for n, err := conn.Read(buf); n >= 8; {
			if err != nil {
				conn.Close()
				return
			}
			cml = append(cml, buf...)
			buf = make([]byte, 8)
			n, err = conn.Read(buf)
			if n < 8 {
				cml = append(cml, buf...)
			}
		}
		fmt.Println(string(cml))
		conn.Write([]byte("+ok\r\n"))
	}
}
