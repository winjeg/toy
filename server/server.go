package server

import (
	"fmt"
	"log"
	"net"

	"github.com/winjeg/toy/resp"
)

const defaultSize = 32

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
		cml, err := readFromConn(conn)
		if err != nil {
			return
		}

		// here to route commands
		fmt.Println(resp.Parse(cml))

		writeErr := write2Conn(conn)
		if writeErr != nil {
			return
		}
	}
}

func readFromConn(conn net.Conn) ([]byte, error) {
	cml := make([]byte, 0, defaultSize)
	buf := make([]byte, defaultSize)
	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if n <= defaultSize {
		cml = append(cml, buf...)
		return cml, nil
	}

	for n >= defaultSize {
		cml = append(cml, buf...)
		buf = make([]byte, defaultSize)
		n, err = conn.Read(buf)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if n < defaultSize {
			cml = append(cml, buf...)
		}
	}
	return cml, nil
}

func write2Conn(conn net.Conn) error {
	_, err := conn.Write([]byte("+ok\r\n"))
	return err
}
