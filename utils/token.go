package utils

import (
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"strconv"
)

func MimcHash(seed []byte, epoch []byte) []byte {
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write(seed)
	_, _ = f.Write(epoch)
	return f.Sum(nil)
}

/*
ComputeToken computes time-based token
*/
func ComputeToken(epoch int, seed string) []byte {
	token := MimcHash([]byte(seed), []byte(strconv.Itoa(epoch)))
	return token
}
