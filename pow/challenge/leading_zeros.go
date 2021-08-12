package challenge

import (
	"hash"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/auvn/go-word-of-wisdom/pow"
)

type LeadingZeros struct {
	rnd         *rand.Rand
	newHash     func() hash.Hash
	hashSumBits uint
	target      *big.Int
	zerosTarget ZerosTarget
}

func NewLeadingZeros(newHash func() hash.Hash, sumBits uint, zeros ZerosTarget) *LeadingZeros {
	targetBits := sumBits - uint(zeros)
	target := big.NewInt(1)
	target.Lsh(target, targetBits)

	return &LeadingZeros{
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
		newHash:     newHash,
		hashSumBits: targetBits,
		target:      target,
		zerosTarget: zeros,
	}
}

func (c *LeadingZeros) Choose(salt []byte) *Puzzle {
	var nonce Nonce
	pow.BytesOrder.PutUint64(nonce[:], c.rnd.Uint64())

	h := c.newHash()
	h.Write(salt)
	h.Write(nonce[:])
	sum := h.Sum(nil)

	var puzzleSize BytesSize
	pow.BytesOrder.PutUint16(puzzleSize[:], uint16(len(sum)))

	return &Puzzle{
		Nonce:      nonce,
		Zeros:      c.zerosTarget,
		PuzzleSize: puzzleSize,
		Puzzle:     h.Sum(nil),
	}
}

func (c *LeadingZeros) Verify(salt []byte, s *Solution) bool {
	h := c.newHash()
	h.Write(salt)
	h.Write(s.Nonce[:])
	puzzle := h.Sum(nil)

	return verifyLeadingZerosSolution(c.newHash(), puzzle, s.Solution[:], c.target)
}

type LeadingZerosSolver struct {
	newHash     func() hash.Hash
	hashSumBits uint
}

func NewLeadingZerosSolver(newHash func() hash.Hash, sumBits uint) *LeadingZerosSolver {
	return &LeadingZerosSolver{
		newHash:     newHash,
		hashSumBits: sumBits,
	}
}

func (s *LeadingZerosSolver) Solve(puzzle *Puzzle) *Solution {
	h := s.newHash()

	target := big.NewInt(1)
	target.Lsh(target, s.hashSumBits-uint(puzzle.Zeros))

	var solution [8]byte
	for i := uint64(0); i < math.MaxUint64; i++ {
		attempt := solution[:]
		pow.BytesOrder.PutUint64(attempt, i)

		h.Reset()

		if verifyLeadingZerosSolution(h, puzzle.Puzzle, attempt, target) {
			break
		}
	}

	return &Solution{
		Nonce:    puzzle.Nonce,
		Solution: solution,
	}
}

func verifyLeadingZerosSolution(h hash.Hash, puzzle, solution []byte, target *big.Int) bool {
	h.Write(puzzle)
	h.Write(solution)
	sum := h.Sum(nil)

	var result big.Int
	result.SetBytes(sum)

	return result.Cmp(target) != 1
}
