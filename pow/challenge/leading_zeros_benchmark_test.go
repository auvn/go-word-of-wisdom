package challenge

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
)

func BenchmarkLeadingZerosChoose(b *testing.B) {
	bench := func(
		b *testing.B,
		hFunc func() hash.Hash,
		sumBits uint,
		zeros ZerosTarget,
	) {

		salt := []byte("my secret salt")
		algo := NewLeadingZeros(hFunc, sumBits, zeros)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			puzzle := algo.Choose(salt)
			if len(puzzle.Puzzle) == 0 {
				b.Error("empty puzzle")
			}
		}

	}

	b.Run("SHA1", func(b *testing.B) {
		bench(b, sha1.New, sha1.Size<<3, 10)
	})

	b.Run("SHA256", func(b *testing.B) {
		bench(b, sha256.New, sha256.Size<<3, 10)
	})

	b.Run("SHA512", func(b *testing.B) {
		bench(b, sha512.New, sha512.Size<<3, 10)
	})
}

func BenchmarkLeadingZerosVerify(b *testing.B) {
	bench := func(b *testing.B, hFunc func() hash.Hash, sumBits uint, zeros ZerosTarget) {
		salt := []byte("my super salt")
		algo := NewLeadingZeros(hFunc, sumBits, zeros)
		algoSolver := NewLeadingZerosSolver(hFunc, sumBits)

		solution := algoSolver.Solve(algo.Choose(salt))

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ok := algo.Verify(salt, solution)
			if !ok {
				b.Error("verification failed")
			}
		}
	}

	b.Run("SHA1", func(b *testing.B) {
		hFunc := sha1.New
		hBits := uint(sha1.Size << 3)

		b.Run("10 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 10)
		})

		b.Run("15 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 15)
		})
	})

	b.Run("SHA256", func(b *testing.B) {
		hFunc := sha256.New
		hBits := uint(sha256.Size << 3)

		b.Run("10 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 10)
		})

		b.Run("15 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 15)
		})
	})
}

func BenchmarkLeadingZerosSolve(b *testing.B) {
	bench := func(b *testing.B, hFunc func() hash.Hash, sumBits uint, zeros ZerosTarget) {
		salt := []byte("my super salt")
		algo := NewLeadingZeros(hFunc, sumBits, zeros)
		algoSolver := NewLeadingZerosSolver(hFunc, sumBits)
		puzzle := algo.Choose(salt)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			algoSolver.Solve(puzzle)
		}
	}

	b.Run("SHA1", func(b *testing.B) {
		hFunc := sha1.New
		hBits := uint(sha1.Size << 3)

		b.Run("10 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 10)
		})

		b.Run("13 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 13)
		})

		b.Run("15 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 15)
		})

		b.Run("20 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 20)
		})
	})

	b.Run("SHA256", func(b *testing.B) {
		hFunc := sha256.New
		hBits := uint(sha256.Size << 3)

		b.Run("10 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 10)
		})

		b.Run("13 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 13)
		})

		b.Run("15 Zeros", func(b *testing.B) {
			bench(b, hFunc, hBits, 15)
		})
	})
}
