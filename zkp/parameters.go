package zkp

import (
	"encoding/hex"
	"github.com/consensys/gnark-crypto/signature"
	"zkrevoke/utils"
)

type PrivateWitnessParameters struct {
	Signature       []byte
	Seed            []byte
	HolderRadomness []byte
}

type PublicWitnessParameters struct {
	PublicKey  signature.PublicKey
	Hash1      []byte
	Challenge  []byte
	Epochs     [][]byte
	ValidUntil []byte
	Tokens     [][]byte
	ClaimsHash []byte
}

type PrivateWitnessParametersV4 struct {
	Seed            []byte
	HolderRadomness []byte
	Signature       []byte
	Claims          [][]byte
}

type PublicWitnessParametersV4 struct {
	Epochs      [][]byte
	Tokens      [][]byte
	PublicKey   signature.PublicKey
	Challenge   []byte
	ValidUntil  []byte
	HashDigests [][]byte
}

func (params PrivateWitnessParameters) String() string {
	var response string
	response = response + " \n private witness parameters: "
	response = response + "seed: " + utils.GetShortString(hex.EncodeToString(params.Seed)) + ", "
	response = response + "holder randomness: " + utils.GetShortString(hex.EncodeToString(params.HolderRadomness)) + ", "
	response = response + "issuer signature: " + utils.GetShortString(hex.EncodeToString(params.Signature)) + ", "
	return response
}

func (params PublicWitnessParameters) String() string {
	var response string
	response = response + " \n public witness parameters: "
	response = response + "Epochs: ("
	for i := 0; i < len(params.Epochs); i++ {
		response = response + hex.EncodeToString(params.Epochs[i]) + ", "
	}
	response = response + "), "
	response = response + "Tokens: ("
	for i := 0; i < len(params.Tokens); i++ {
		response = response + utils.GetShortString(hex.EncodeToString(params.Tokens[i])) + ", "
	}
	response = response + "), "
	response = response + "issuer's public key: " + utils.GetShortString(hex.EncodeToString(params.PublicKey.Bytes())) + ", "
	response = response + "verifier's challenge: " + utils.GetShortString(hex.EncodeToString(params.Challenge)) + ", "
	response = response + "hash of (challenge and secret randomness): " + utils.GetShortString(hex.EncodeToString(params.Hash1)) + ", "
	response = response + "validity of VC: " + utils.GetShortString(hex.EncodeToString(params.ValidUntil))
	return response
}
