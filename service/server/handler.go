package server

import (
	"net"
)

type ConnectionHandler interface {
	ServeConn(c net.Conn)
}

type ConnectionHandlerFunc func(c net.Conn)

func (fn ConnectionHandlerFunc) ServeConn(c net.Conn) {
	fn(c)
}
