package utils

import (
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
)

/*
PublicParams stores public keys used by issuer.
This is temporary.
Once smart contracts are ready, store the keys in the smart contract.
Maybe, this can be used for local testing.
*/
type PublicParams struct {
	Ccs             constraint.ConstraintSystem
	EddsaPublicKey  signature.PublicKey
	ZkpProvingKey   groth16.ProvingKey
	ZkpVerifyingKey groth16.VerifyingKey
}
