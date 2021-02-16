package resp

import (
	"net"
)

const (
	readSize = 4096
)

type Reader struct {
	Conn net.Conn
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
		n, err = r.Conn.Read(buf)
		if err != nil {
			return nil, err
		}
		if n < readSize {
			cml = append(cml, buf...)
		}
	}
	return cml, nil
}
