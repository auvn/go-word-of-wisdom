package challenge

type ChooseVerifier interface {
	Choose(salt []byte) *Puzzle
	Verify(salt []byte, s *Solution) bool
}

type Solver interface {
	Solve(puzzle *Puzzle) *Solution
}

type Rand interface {
	Read([]byte) (int, error)
	Uint64() uint64
}
