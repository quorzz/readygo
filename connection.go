package readygo

import (
	"bufio"
	"github.com/quorzz/redis-protocol"
	"net"
)

type connection struct {
	br *bufio.Reader
	bw *bufio.Writer
}

func NewConnection(netConn net.Conn) *connection {
	return &connection{
		br: bufio.NewReader(netConn),
		bw: bufio.NewWriter(netConn),
	}
}

func (conn *connection) Send(args ...interface{}) error {

	rawCmds, err := protocol.PackCommand(args...)
	if err != nil {
		return err
	}

	if _, err := conn.bw.Write(rawCmds); err != nil {
		return err
	}

	return nil
}

func (conn *connection) Flush() error {
	if err := conn.bw.Flush(); err != nil {
		return err
	}
	return nil
}

func (conn *connection) Receive() (*protocol.Message, error) {

	if message, err := protocol.UnpackFromReader(conn.br); err != nil {
		return nil, err
	} else {
		return message, nil
	}
}

func (conn *connection) Execute(args ...interface{}) (*protocol.Message, error) {

	if err := conn.Send(args...); err != nil {
		return nil, err
	}

	if err := conn.Flush(); err != nil {
		return nil, err
	}

	if message, err := conn.Receive(); err != nil {
		return nil, err
	} else {
		return message, nil
	}
}
