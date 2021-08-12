package powproto

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/auvn/go-word-of-wisdom/pow/challenge"
	"github.com/auvn/go-word-of-wisdom/service/server"
)

// TODO: memory bound/hard algos
// TODO: connection could be put in a waiting pool until solution is not received - epoll unix kernel feature
// TODO: connection could be remade by client in case of recoverable salt for the puzzle

type ChallengeHandler struct {
	parent server.ConnectionHandler
	proto  Proto
}

func NewChallengeHandler(parent server.ConnectionHandler, proto Proto) *ChallengeHandler {
	return &ChallengeHandler{
		parent: parent,
		proto:  proto,
	}
}

func (h *ChallengeHandler) ServeConn(conn net.Conn) {
	ok, err := h.proto.Run(conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		fmt.Println(fmt.Errorf("run protocol: %w", err))
		h.closeConn(conn)
		return
	}

	// failed pow verification
	if !ok {
		h.closeConn(conn)
		return
	}

	// run is successful
	h.parent.ServeConn(conn)
}

func (h *ChallengeHandler) closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		fmt.Println("conn.Close:", err)
	}
}

type Proto interface {
	Run(conn net.Conn) (bool, error)
}

type KeepAliveLink struct {
	challenge challenge.ChooseVerifier
}

func NewKeepAliveLink(cv challenge.ChooseVerifier) *KeepAliveLink {
	return &KeepAliveLink{
		challenge: cv,
	}
}

func (d *KeepAliveLink) Run(conn net.Conn) (ok bool, err error) {
	// TODO: set deadline according to difficulty
	deadline := time.Now().Add(2 * time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		return false, fmt.Errorf("set deadline: %w", err)
	}

	defer func() {
		if dErr := conn.SetDeadline(time.Time{}); dErr != nil {
			if err != nil {
				err = fmt.Errorf("%w: reset deadline: %s", err, dErr)
				return
			}

			ok = false
			err = fmt.Errorf("reset deadline: %w", dErr)
		}
	}()

	proto := challenge.Protocol{ReadWriter: conn}

	puzzle := d.challenge.Choose(nil)
	if err := proto.WritePuzzle(puzzle); err != nil {
		return false, fmt.Errorf("write puzzle: %w", err)
	}

	solution, err := proto.ReadSolution()
	if err != nil {
		return false, fmt.Errorf("read solution: %w", err)
	}

	// checking nonce number is not compromised
	if puzzle.Nonce != solution.Nonce {
		return false, nil
	}

	return d.challenge.Verify(nil, solution), nil
}
