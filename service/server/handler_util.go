package server

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"
)

type ConnectionMonitor struct {
	served uint64
}

func NewConnectionMonitor() *ConnectionMonitor {
	m := ConnectionMonitor{}
	go m.monitor()
	return &m
}

func (m *ConnectionMonitor) ServeConn(c net.Conn) {
	atomic.AddUint64(&m.served, 1)
}

func (m *ConnectionMonitor) monitor() {
	var buf bytes.Buffer

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for range t.C {
		_, _ = fmt.Fprintln(&buf, "Stats: ")
		_, _ = fmt.Fprintf(&buf, "Served: %d\n", atomic.LoadUint64(&m.served))
		_, _ = fmt.Fprintln(&buf)
		_, _ = buf.WriteTo(os.Stdout)
		buf.Reset()
	}
}

type MultiConnectionHandler struct {
	hh []ConnectionHandler
}

func NewMultiConnectionHandler(hh ...ConnectionHandler) *MultiConnectionHandler {
	if len(hh) == 0 {
		panic("empty handlers")
	}

	return &MultiConnectionHandler{
		hh: hh,
	}
}

func (h *MultiConnectionHandler) ServeConn(c net.Conn) {
	for _, handler := range h.hh {
		handler.ServeConn(c)
	}
}
