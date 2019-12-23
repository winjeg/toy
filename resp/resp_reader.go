package resp

import (
	"net"
)

const (
	readSize = 32
)

type Reader struct {
	conn net.Conn
}

func (r *Reader) Read() ([]byte, error) {
	cml := make([]byte, 0, readSize)
	buf := make([]byte, readSize)
	n, err := r.conn.Read(buf)
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
		n, err = r.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		if n < readSize {
			cml = append(cml, buf...)
		}
	}
	return cml, nil
}
