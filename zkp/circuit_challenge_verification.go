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
	"github.com/consensys/gnark/std/hash/mimc"
	"go.uber.org/zap"
)

type CircuitChallengeVerification struct {
	CurveID tedwards.ID `gnark:",public"`

	Challenge       frontend.Variable `gnark:",public"` // challenge is provided by verifier when requesting VP from a holder
	HolderRadomness frontend.Variable // random value used by holder
	Hash1           frontend.Variable `gnark:",public"` // H(Challenge || HolderRadomness)
}

func NewCircuitForChallengeVerification() constraint.ConstraintSystem {

	//logger.Disable()
	var circuit CircuitChallengeVerification
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
	**** H(challenge||holder_randomness)) == Hash1

-------------- verifies that the hash1 is computed from the challenge issued by a verifier and a secret randomness provided by the holder

*/

func (circuit *CircuitChallengeVerification) Define(api frontend.API) error {

	//**** H(challenge||holder_randomness)) == Hash1
	mimc3, _ := mimc.NewMiMC(api)
	mimc3.Write(circuit.Challenge, circuit.HolderRadomness)
	computedHash2 := mimc3.Sum()
	api.AssertIsEqual(circuit.Hash1, computedHash2)

	return nil
}

/*
PrivateWitnessGeneration generates the private witness
*/
func PrivateWitnessGenerationForChallengeVerification(pubParams PublicWitnessParameters, privParams PrivateWitnessParameters) witness.Witness {
	assignment := &CircuitChallengeVerification{
		HolderRadomness: frontend.Variable(privParams.HolderRadomness),
		CurveID:         tedwards.BN254,
		Challenge:       frontend.Variable(pubParams.Challenge),
		Hash1:           frontend.Variable(pubParams.Hash1),
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		zap.S().Errorf("error generating witness: %v", err)
	}

	return witness
}

/*
PublicWitnessGeneration generates the public witness given public inputs to the circuit
*/
func PublicWitnessGenerationForChallengeVerification(pubParams PublicWitnessParameters) witness.Witness {
	assignment := &CircuitChallengeVerification{
		CurveID:   tedwards.BN254,
		Challenge: frontend.Variable(pubParams.Challenge),
		Hash1:     frontend.Variable(pubParams.Hash1),
	}

	publicWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	//fmt.Println("public witness: ", publicWitness)
	if err != nil {
		fmt.Errorf("error generating public witness: %v", err)
	}

	return publicWitness
}
