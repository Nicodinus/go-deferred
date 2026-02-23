package tests

import (
	"crypto/rand"
	"math/big"
)

func RandRange(min, max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		panic(err)
	}
	return nBig.Int64() + min
}
