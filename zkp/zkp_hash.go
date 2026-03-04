package zkp

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
)

type hashCircuit struct {
	Input  frontend.Variable
	Digest frontend.Variable `gnark:",public"`
}

func (circuit *hashCircuit) Define(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(circuit.Input)
	computedHash := mimc.Sum()
	fmt.Println("computed hash: ", computedHash, "\n input hash: ", circuit.Digest)
	api.AssertIsEqual(circuit.Digest, computedHash)

	return nil

}

func test_hash_circuit() {

	input := []byte("test")

	f := bn254_mimc.NewMiMC()
	_, _ = f.Write(input)
	digest := f.Sum(nil)

	fmt.Println("input: ", input)
	fmt.Println("digest: ", digest)

	var circuit hashCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	pk, vk := SetupGroth(ccs)

	assignment := &hashCircuit{
		Input:  frontend.Variable(input),
		Digest: frontend.Variable(digest),
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		fmt.Errorf("error generating witness: %v", err)
	} else {
		fmt.Println("witness: ", witness)
	}
	publicWitness, err := witness.Public()

	if err != nil {
		fmt.Println("error in creating the public witness")
	}
	fmt.Println("public witness: ", publicWitness)

	proof := ProveGroth(ccs, pk, witness)
	status := VerifyGroth(proof, vk, publicWitness)

	fmt.Println("Verification status: ", status)
}
