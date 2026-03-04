package zkp

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	twistededwards_zkp "github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	eddsa_zkp "github.com/consensys/gnark/std/signature/eddsa"
	"math/rand"
	"time"
)

type eddsaCircuit struct {
	curveID   tedwards.ID
	PublicKey eddsa_zkp.PublicKey `gnark:",public"`
	Message   frontend.Variable   //
	Signature eddsa_zkp.Signature `gnark:",public"`
}

func (circuit *eddsaCircuit) Define(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	curve, err := twistededwards_zkp.NewEdCurve(api, circuit.curveID)
	if err != nil {
		return err
	}

	return eddsa_zkp.Verify(curve, circuit.Signature, circuit.Message, circuit.PublicKey, &mimc)

}

func test_eddsa_circuit() {
	msg := []byte("test")
	fmt.Println("msg: ", msg)

	randomness := rand.New(rand.NewSource(time.Now().Unix()))
	privateKey, err := eddsa.New(tedwards.BN254, randomness)
	publicKey := privateKey.Public()

	if err != nil {
		fmt.Println("error in generating private key: ", err)
	}
	hFunc := hash.MIMC_BN254.New()

	signature, err := privateKey.Sign(msg, hFunc)
	if err != nil {
		fmt.Println("signature generation failed: ", err)
	}
	fmt.Println("public key: ", publicKey)
	fmt.Println("message: ", msg)
	fmt.Println("signature: ", signature)

	var circuit eddsaCircuit
	circuit.curveID = tedwards.ID(ecc.BN254)
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	pk, vk := SetupGroth(ccs)

	assignment := &eddsaCircuit{
		curveID: tedwards.BN254,
	}

	// public key bytes
	_publicKey := publicKey.Bytes()

	// assign public key values
	assignment.PublicKey.Assign(tedwards.BN254, _publicKey[:32])

	fmt.Println("signature: ", signature)
	assignment.Signature.Assign(tedwards.BN254, signature)

	assignment.Message = msg

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		fmt.Errorf("error generating witness: %v", err)
	} else {
		fmt.Println("witness: ", witness)
	}
	publicWitness, _ := witness.Public()
	fmt.Println("public witness: ", publicWitness)

	proof := ProveGroth(ccs, pk, witness)
	status := VerifyGroth(proof, vk, publicWitness)

	fmt.Println("Verification status: ", status)
}
