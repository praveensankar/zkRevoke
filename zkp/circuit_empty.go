package zkp

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/logger"
	"go.uber.org/zap"
)

/*
Contains an empty circuit to benchmark the overhead of groth16
*/
type CircuitEmpty struct {
	CurveID tedwards.ID `gnark:",public"`
}

func NewCircuitEmpty() constraint.ConstraintSystem {

	logger.Disable()
	var circuit CircuitEmpty
	circuit.CurveID = tedwards.ID(ecc.BN254)
	// Issuer creates a circuit that defines our ZKP constraints
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		zap.S().Infoln("error compiling circuit: ", err)
	}

	return ccs
}

/*
Define defines our circuit. The circuit consists of the following conditions:

	**** Verify(PK, H(seed || claims)) == sig
	**** H(Seed || Epoch[i]) == token[i]: i = 1...numberOfTokens
	**** H(claims[i] || challenge || holder_randomness)) == HashDigests[i]: i = 1...numberOfClaims


Verify(PK, H(seed || claims)) == sig:
-------------- verifies that the issuer's signature is present on seed and claims
H(Seed || Epoch) == token:
-------------- verifies that the token is generated from the secret seed and an epoch value
H(claims[i] || challenge || holder_randomness)) == HashDigests[i]
-------------- verifies that the hash digests of claims are computed using claims encoded in the VC. Furthermore, the hash
digests are randomized using a secret randomness provided by the holder and a challenge provided by the verifier

*/

func (circuit *CircuitEmpty) Define(api frontend.API) error {

	return nil
}

/*
PrivateWitnessGeneration generates the private witness
*/
func PrivateWitnessGenerationForEmptyCircuit() witness.Witness {
	assignment := &CircuitEmpty{
		CurveID: tedwards.BN254,
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
func PublicWitnessGenerationForEmptyCircuit() witness.Witness {
	assignment := &CircuitEmpty{
		CurveID: tedwards.BN254,
	}

	publicWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	//fmt.Println("public witness: ", publicWitness)
	if err != nil {
		fmt.Errorf("error generating public witness: %v", err)
	}

	return publicWitness
}
