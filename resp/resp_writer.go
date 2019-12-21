package resp

import "net"

type Writer struct {
	Conn net.Conn
}

func (w *Writer) Write(data []byte) error {
	_, err := w.Conn.Write(data)
	return err
}
