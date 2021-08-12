package challenge

import (
	"fmt"
	"io"
)

type Protocol struct {
	ReadWriter io.ReadWriter
}

func (p *Protocol) ReadSolution() (*Solution, error) {
	nonce, err := p.readNonce()
	if err != nil {
		return nil, err
	}

	var solution SolutionBytes
	if err := p.readExact(solution[:]); err != nil {
		return nil, fmt.Errorf("read solution: %w", err)
	}

	return &Solution{
		Nonce:    nonce,
		Solution: solution,
	}, nil
}

func (p *Protocol) WriteSolution(s *Solution) error {
	if err := p.writeNonce(s.Nonce[:]); err != nil {
		return err
	}

	if err := p.writeExact(s.Solution[:]); err != nil {
		return fmt.Errorf("write solution: %w", err)
	}
	return nil
}

func (p *Protocol) WritePuzzle(puzzle *Puzzle) error {
	if err := p.writeNonce(puzzle.Nonce[:]); err != nil {
		return err
	}

	if err := p.writeExact([]byte{puzzle.Zeros}); err != nil {
		return fmt.Errorf("write zeros: %w", err)
	}

	if err := p.writeExact(puzzle.PuzzleSize[:]); err != nil {
		return fmt.Errorf("write puzzle size: %w", err)
	}

	if err := p.writeExact(puzzle.Puzzle); err != nil {
		return fmt.Errorf("write puzzle: %w", err)
	}

	return nil
}

func (p *Protocol) ReadPuzzle() (*Puzzle, error) {
	nonce, err := p.readNonce()
	if err != nil {
		return nil, err
	}

	var zeros [1]byte
	if err := p.readExact(zeros[:]); err != nil {
		return nil, fmt.Errorf("read zeros: %w", err)
	}

	var puzzleSize BytesSize
	if err := p.readExact(puzzleSize[:]); err != nil {
		return nil, fmt.Errorf("read puzzle size: %w", err)
	}

	puzzle := make([]byte, puzzleSize.Num())
	if err := p.readExact(puzzle); err != nil {
		return nil, fmt.Errorf("read puzzle: %w", err)
	}

	return &Puzzle{
		Nonce:      nonce,
		Zeros:      zeros[0],
		PuzzleSize: puzzleSize,
		Puzzle:     puzzle,
	}, nil
}

func (p *Protocol) readNonce() (Nonce, error) {
	var nonce Nonce
	if err := p.readExact(nonce[:]); err != nil {
		return nonce, fmt.Errorf("read nonce")
	}

	return nonce, nil
}

func (p *Protocol) writeNonce(nonce []byte) error {
	if err := p.writeExact(nonce); err != nil {
		return fmt.Errorf("write nonce: %w", err)
	}

	return nil
}

func (p *Protocol) readExact(bb []byte) error {
	n, err := p.ReadWriter.Read(bb)
	if err != nil {
		return err
	}

	if n != len(bb) {
		return fmt.Errorf("incomplete read: want: %d, got: %d", len(bb), n)
	}

	return nil
}

func (p *Protocol) writeExact(bb []byte) error {
	n, err := p.ReadWriter.Write(bb)
	if err != nil {
		return err
	}

	if n != len(bb) {
		return fmt.Errorf("incomplete write: want: %d, got: %d", len(bb), n)
	}

	return nil
}
