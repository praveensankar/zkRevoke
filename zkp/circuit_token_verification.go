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

type CircuitTokenVerification struct {
	CurveID tedwards.ID `gnark:",public"`

	Seed   frontend.Variable   // seed is the private input known only to the holder and the issuer
	Epochs []frontend.Variable `gnark:",public"` // Epoch indicates the current epoch.  The number of durations elapsed from the issuance of VC.
	Tokens []frontend.Variable `gnark:",public"` // token = H(Seed || Epoch). token is inserted into a revocation list.

}

func NewCircuitForTokenVerification(n int) constraint.ConstraintSystem {

	//logger.Disable()
	var circuit CircuitTokenVerification
	circuit.CurveID = tedwards.ID(ecc.BN254)
	circuit.Tokens = make([]frontend.Variable, n)
	circuit.Epochs = make([]frontend.Variable, n)
	for i := 0; i < n; i++ {
		circuit.Tokens[i] = frontend.Variable(i)
		circuit.Epochs[i] = frontend.Variable(i)
	}
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


	**** H(Seed || Epoch[i]) == token[i]: i = 1...numberOfTokens

H(Seed || Epoch) == token:
-------------- verifies that the token is generated from the secret seed and an epoch value
*/

func (circuit *CircuitTokenVerification) Define(api frontend.API) error {

	numberOfTokens := len(circuit.Tokens)

	//**** H(Seed || Epoch[i]) == token[i]: i = 1...numberOfTokens
	for i := 0; i < numberOfTokens; i++ {
		mimc5, _ := mimc.NewMiMC(api)
		mimc5.Write(circuit.Seed, circuit.Epochs[i])
		computedHash3 := mimc5.Sum()

		// 	**** H(VC_ID || seed) == token
		api.AssertIsEqual(circuit.Tokens[i], computedHash3)
	}

	return nil
}

/*
PrivateWitnessGeneration generates the private witness
*/
func PrivateWitnessGenerationForTokenVerificationCircuit(pubParams PublicWitnessParameters, privParams PrivateWitnessParameters) witness.Witness {
	assignment := &CircuitTokenVerification{
		Seed:    frontend.Variable(privParams.Seed),
		CurveID: tedwards.BN254,
		Epochs:  make([]frontend.Variable, len(pubParams.Epochs)),
		Tokens:  make([]frontend.Variable, len(pubParams.Tokens)),
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Epochs[i] = frontend.Variable(pubParams.Epochs[i])
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Tokens[i] = frontend.Variable(pubParams.Tokens[i])
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
func PublicWitnessGenerationForTokenVerificationCircuit(pubParams PublicWitnessParameters) witness.Witness {
	assignment := &CircuitTokenVerification{
		CurveID: tedwards.BN254,
		Epochs:  make([]frontend.Variable, len(pubParams.Epochs)),
		Tokens:  make([]frontend.Variable, len(pubParams.Tokens)),
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Epochs[i] = frontend.Variable(pubParams.Epochs[i])
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Tokens[i] = frontend.Variable(pubParams.Tokens[i])
	}

	publicWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	//fmt.Println("public witness: ", publicWitness)
	if err != nil {
		fmt.Errorf("error generating public witness: %v", err)
	}

	return publicWitness
}
