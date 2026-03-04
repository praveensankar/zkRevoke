package zkp

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/profile"
	twistededwards_zkp "github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	eddsa_zkp "github.com/consensys/gnark/std/signature/eddsa"
	"go.uber.org/zap"
)

type CircuitSignatureVerifcation struct {
	CurveID tedwards.ID `gnark:",public"`

	Seed       frontend.Variable // seed is the private input known only to the holder and the issuer
	ClaimsHash frontend.Variable `gnark:",public"` // hash of claims
	ValidUntil frontend.Variable `gnark:",public"` // ValidUntil is the expiration date of the VC.

	PublicKey eddsa_zkp.PublicKey `gnark:",public"` // public key of the issuer
	Signature eddsa_zkp.Signature // Signature ← Sign(PublicKey,  H(seed || h(claims) || validUntil))
}

func NewCircuitSignatureVerification() constraint.ConstraintSystem {

	//logger.Disable()
	var circuit CircuitSignatureVerifcation
	circuit.CurveID = tedwards.ID(ecc.BN254)

	p := profile.Start()
	// Issuer creates a circuit that defines our ZKP constraints
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		zap.S().Infoln("error compiling circuit: ", err)
	}
	p.Stop()
	return ccs
}

/*
Define defines our circuit. The circuit consists of the following conditions:

	**** Verify(PK, H(seed || H(claims) || validUntil)) == sig


Verify(PK, H(seed || claims || validUntil)) == sig:
-------------- verifies that the issuer's signature is present on vc_id, seed and expiry date of VC

*/

func (circuit *CircuitSignatureVerifcation) Define(api frontend.API) error {
	curve, err := twistededwards_zkp.NewEdCurve(api, circuit.CurveID)
	if err != nil {
		return err
	}

	//**** Verify(PK, H(seed || claims || validUntil)) == sig
	mimc1, _ := mimc.NewMiMC(api)
	mimc1.Write(circuit.Seed)
	mimc1.Write(circuit.ClaimsHash)
	mimc1.Write(circuit.ValidUntil)
	computedHash := mimc1.Sum()
	mimc2, _ := mimc.NewMiMC(api)
	err = eddsa_zkp.Verify(curve, circuit.Signature, computedHash, circuit.PublicKey, &mimc2)

	return nil
}

/*
PrivateWitnessGeneration generates the private witness
*/
func PrivateWitnessGenerationForCircuitSignatureVerification(pubParams PublicWitnessParameters, privParams PrivateWitnessParameters) witness.Witness {
	assignment := &CircuitSignatureVerifcation{
		Seed:       frontend.Variable(privParams.Seed),
		CurveID:    tedwards.BN254,
		ClaimsHash: frontend.Variable(pubParams.ClaimsHash),
		ValidUntil: frontend.Variable(pubParams.ValidUntil),
	}

	// public key bytes
	_publicKey := pubParams.PublicKey.Bytes()

	// assign public key values
	assignment.PublicKey.Assign(tedwards.BN254, _publicKey[:32])

	assignment.Signature.Assign(tedwards.BN254, privParams.Signature)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		zap.S().Errorf("error generating witness: %v", err)
	}
	return witness
}

/*
PublicWitnessGeneration generates the public witness given public inputs to the circuit
*/
func PublicWitnessGenerationForCircuitSignatureVerification(pubParams PublicWitnessParameters) witness.Witness {
	assignment := &CircuitSignatureVerifcation{
		CurveID:    tedwards.BN254,
		ValidUntil: frontend.Variable(pubParams.ValidUntil),
		ClaimsHash: frontend.Variable(pubParams.ClaimsHash),
	}

	// public key bytes
	_publicKey := pubParams.PublicKey.Bytes()

	// assign public key values
	assignment.PublicKey.Assign(tedwards.BN254, _publicKey[:32])

	publicWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	//fmt.Println("public witness: ", publicWitness)
	if err != nil {
		fmt.Errorf("error generating public witness: %v", err)
	}

	return publicWitness
}
