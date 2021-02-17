package resp

import (
	"fmt"
	"net"
)

const (
	readSize = 3
)

type Reader struct {
	Conn net.Conn
	Data []byte // data buffer.
}

func (r *Reader) Read() ([]byte, error) {
	cml := make([]byte, 0, readSize)
	buf := make([]byte, readSize)
	n, err := r.Conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < readSize {
		cml = append(cml, buf...)
		return cml, nil
	}

	for n >= readSize {
		cml = append(cml, buf...)
		buf = make([]byte, readSize)
		fmt.Println(string(cml))
		n, err = r.Conn.Read(buf)
		if err != nil {
			return nil, err
		}
		if n < readSize {
			cml = append(cml, buf...)
			break
		}
	}
	return cml, nil
}
