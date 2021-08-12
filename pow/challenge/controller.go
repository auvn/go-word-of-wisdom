package challenge

import (
	"crypto/sha1"
)

var Simple = NewSimpleChallenge()

type Controller struct {
	ChooseVerifier
	Solver
}

func NewSimpleChallenge() *Controller {
	hFunc := sha1.New
	hSumBits := uint(sha1.Size << 3)

	return &Controller{
		ChooseVerifier: NewLeadingZeros(hFunc, hSumBits, 14),
		Solver:         NewLeadingZerosSolver(hFunc, hSumBits),
	}
}
