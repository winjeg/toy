package server

import (
	"fmt"
	"github.com/winjeg/toy/resp"
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
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return
		}
		for  n >= 8 {
			cml = append(cml, buf...)
			buf = make([]byte, 8)
			n, err = conn.Read(buf)
			if err != nil {
				conn.Close()
				return
			}
			if n < 8 {
				cml = append(cml, buf...)
			}
		}
		fmt.Println(resp.Parse(cml))
		conn.Write([]byte("+ok\r\n"))
	}
}
