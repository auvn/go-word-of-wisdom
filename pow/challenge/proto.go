package challenge

import (
	"fmt"
	"io"
)

type Protocol struct {
	ReadWriter io.ReadWriter
}

func (rw *Protocol) ReadSolution() (*Solution, error) {
	nonce, err := rw.readNonce()
	if err != nil {
		return nil, err
	}

	var solution SolutionBytes
	if err := rw.readExact(solution[:]); err != nil {
		return nil, fmt.Errorf("read solution: %w", err)
	}

	return &Solution{
		Nonce:    nonce,
		Solution: solution,
	}, nil
}

func (rw *Protocol) WriteSolution(s *Solution) error {
	if err := rw.writeNonce(s.Nonce[:]); err != nil {
		return err
	}

	if err := rw.writeExact(s.Solution[:]); err != nil {
		return fmt.Errorf("write solution: %w", err)
	}
	return nil
}

func (rw *Protocol) WritePuzzle(p *Puzzle) error {
	if err := rw.writeNonce(p.Nonce[:]); err != nil {
		return err
	}

	if err := rw.writeExact([]byte{p.Zeros}); err != nil {
		return fmt.Errorf("write zeros: %w", err)
	}

	if err := rw.writeExact(p.PuzzleSize[:]); err != nil {
		return fmt.Errorf("write puzzle size: %w", err)
	}

	if err := rw.writeExact(p.Puzzle); err != nil {
		return fmt.Errorf("write puzzle: %w", err)
	}

	return nil
}

func (rw *Protocol) ReadPuzzle() (*Puzzle, error) {
	nonce, err := rw.readNonce()
	if err != nil {
		return nil, err
	}

	var zeros [1]byte
	if err := rw.readExact(zeros[:]); err != nil {
		return nil, fmt.Errorf("read zeros: %w", err)
	}

	var puzzleSize BytesSize
	if err := rw.readExact(puzzleSize[:]); err != nil {
		return nil, fmt.Errorf("read puzzle size: %w", err)
	}

	puzzle := make([]byte, puzzleSize.Num())
	if err := rw.readExact(puzzle); err != nil {
		return nil, fmt.Errorf("read puzzle: %w", err)
	}

	return &Puzzle{
		Nonce:      nonce,
		Zeros:      zeros[0],
		PuzzleSize: puzzleSize,
		Puzzle:     puzzle,
	}, nil
}

func (rw *Protocol) readNonce() (Nonce, error) {
	var nonce Nonce
	if err := rw.readExact(nonce[:]); err != nil {
		return nonce, fmt.Errorf("read nonce")
	}

	return nonce, nil
}

func (rw *Protocol) writeNonce(nonce []byte) error {
	if err := rw.writeExact(nonce); err != nil {
		return fmt.Errorf("write nonce: %w", err)
	}

	return nil
}

func (rw *Protocol) readExact(bb []byte) error {
	n, err := rw.ReadWriter.Read(bb)
	if err != nil {
		return err
	}

	if n != len(bb) {
		return fmt.Errorf("incomplete read: want: %d, got: %d", len(bb), n)
	}

	return nil
}

func (rw *Protocol) writeExact(bb []byte) error {
	n, err := rw.ReadWriter.Write(bb)
	if err != nil {
		return err
	}

	if n != len(bb) {
		return fmt.Errorf("incomplete write: want: %d, got: %d", len(bb), n)
	}

	return nil
}
