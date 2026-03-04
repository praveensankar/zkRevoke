package zkp

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"go.uber.org/zap"
)

/*
Setup takes a constraint system as a input and generates the proving and verification key

Output:
1) proving key- plonk.ProvingKey
2) verification key- plonk.VerifyingKey
*/
func SetupGroth(ccs constraint.ConstraintSystem) (groth16.ProvingKey, groth16.VerifyingKey) {
	
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		fmt.Println("error set up: ", err)
	}
	return pk, vk
}

func ProveGroth(ccs constraint.ConstraintSystem, pk groth16.ProvingKey, witness witness.Witness) groth16.Proof {
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		zap.S().Errorf("***ZKP***:error generating proof: %v", err)
		return nil
	} else {

		return proof
	}

}

func VerifyGroth(proof groth16.Proof, vk groth16.VerifyingKey, publicWitness witness.Witness) bool {
	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Errorf("error verifiying the proof: ", err)
		return false
	}
	return true
}

func GrothProofToBytes(proof groth16.Proof) ([]byte, error) {
	var buf bytes.Buffer
	_, err := proof.WriteTo(&buf)
	return buf.Bytes(), err
}

func GrothProvingKeyToBytes(publicKey groth16.ProvingKey) ([]byte, error) {
	var buf bytes.Buffer
	_, err := publicKey.WriteTo(&buf)
	return buf.Bytes(), err
}

func BytesToGrothProvingKey(data []byte) (groth16.ProvingKey, error) {
	pk := groth16.NewProvingKey(ecc.BN254)
	_, err := pk.ReadFrom(bytes.NewReader(data))
	return pk, err
}

func GrothVerifyingKeyToBytes(verifyingKey groth16.VerifyingKey) ([]byte, error) {
	var buf bytes.Buffer
	_, err := verifyingKey.WriteTo(&buf)
	return buf.Bytes(), err
}

func BytesToGrothVerifyingKey(data []byte) (groth16.VerifyingKey, error) {
	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err := vk.ReadFrom(bytes.NewReader(data))
	return vk, err
}

func BytesToGrothProof(data []byte) (groth16.Proof, error) {
	proof := groth16.NewProof(ecc.BN254) // Specify curve explicitly
	_, err := proof.ReadFrom(bytes.NewReader(data))
	return proof, err
}
