package challenge

import (
	"fmt"

	"github.com/auvn/go-word-of-wisdom/pow"
)

type Nonce = [8]byte

func NewNonce(nonce uint64) Nonce {
	var bb Nonce
	pow.BytesOrder.PutUint64(bb[:], nonce)
	return bb
}

type ZerosTarget = byte
type SolutionBytes = [8]byte
type BytesSize [2]byte

func (s BytesSize) Num() uint16 {
	value := pow.BytesOrder.Uint16(s[:])

	if value > maxBytesSize {
		value = maxBytesSize
	}

	return value
}

type Hash = []byte

type Puzzle struct {
	Nonce      Nonce
	Zeros      ZerosTarget
	PuzzleSize BytesSize
	Puzzle     Hash
}

func (p Puzzle) String() string {
	return fmt.Sprintf(
		"%x:%x:%x:%x",
		p.Nonce, p.Zeros, p.PuzzleSize, p.Puzzle)
}

type Solution struct {
	Nonce    Nonce
	Solution SolutionBytes
}

func (s Solution) String() string {
	return fmt.Sprintf(
		"%x:%x",
		s.Nonce, s.Solution)
}

const (
	maxBytesSize uint16 = 64
)
