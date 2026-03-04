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
	twistededwards_zkp "github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	eddsa_zkp "github.com/consensys/gnark/std/signature/eddsa"
	"go.uber.org/zap"
)

type CircuitV4 struct {
	CurveID tedwards.ID `gnark:",public"`

	Seed       frontend.Variable   // seed is the private input known only to the holder and the issuer
	ValidUntil frontend.Variable   `gnark:",public"` // vc expiry time
	Claims     []frontend.Variable // claims encoded in the VC
	Epochs     []frontend.Variable `gnark:",public"` // Epoch indicates the current epoch.  The number of durations elapsed from the issuance of VC.
	Tokens     []frontend.Variable `gnark:",public"` // token = H(Seed || Epoch). token is inserted into a revocation list.

	PublicKey eddsa_zkp.PublicKey `gnark:",public"` // public key of the issuer
	Signature eddsa_zkp.Signature // Signature ← Sign(PublicKey,  H(seed || claims || validUntil))

	Challenge         frontend.Variable   `gnark:",public"` // challenge is provided by verifier when requesting VP from a holder
	Holder_randomness frontend.Variable   // random value used by holder
	HashDigests       []frontend.Variable `gnark:",public"` // HashDigests_i = H(claims_i || Challenge || Holder_randomness)
}

func NewCircuitv4(n int, l int) constraint.ConstraintSystem {

	logger.Disable()
	var circuit CircuitV4
	circuit.CurveID = tedwards.ID(ecc.BN254)
	circuit.Tokens = make([]frontend.Variable, n)
	circuit.Epochs = make([]frontend.Variable, n)
	for i := 0; i < n; i++ {
		circuit.Tokens[i] = frontend.Variable(i)
		circuit.Epochs[i] = frontend.Variable(i)
	}

	circuit.Claims = make([]frontend.Variable, l)
	circuit.HashDigests = make([]frontend.Variable, l)
	for i := 0; i < l; i++ {
		circuit.Claims[i] = frontend.Variable(i)
		circuit.HashDigests[i] = frontend.Variable(i)
	}

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

func (circuit *CircuitV4) Define(api frontend.API) error {
	curve, err := twistededwards_zkp.NewEdCurve(api, circuit.CurveID)
	if err != nil {
		return err
	}

	numberOfClaims := len(circuit.Claims)
	numberOfTokens := len(circuit.Tokens)

	//**** Verify(PK, H(seed || claims || validUntil)) == sig
	mimc1, _ := mimc.NewMiMC(api)
	mimc1.Write(circuit.Seed)

	for i := 0; i < numberOfClaims; i++ {
		mimc1.Write(circuit.Claims[i])
	}
	mimc1.Write(circuit.ValidUntil)
	computedHash := mimc1.Sum()
	mimc2, _ := mimc.NewMiMC(api)
	err = eddsa_zkp.Verify(curve, circuit.Signature, computedHash, circuit.PublicKey, &mimc2)

	//**** H(Seed || Epoch[i]) == token[i]: i = 1...numberOfTokens

	for i := 0; i < numberOfTokens; i++ {
		mimc5, _ := mimc.NewMiMC(api)
		mimc5.Write(circuit.Seed, circuit.Epochs[i])
		computedHash3 := mimc5.Sum()

		// 	**** H(VC_ID || seed) == token
		api.AssertIsEqual(circuit.Tokens[i], computedHash3)
	}

	//**** H(claims[i] || challenge || holder_randomness)) == HashDigests[i]
	for i := 0; i < numberOfClaims; i++ {
		mimc3, _ := mimc.NewMiMC(api)
		mimc3.Write(circuit.Claims[i], circuit.Challenge, circuit.Holder_randomness)
		computedHash2 := mimc3.Sum()
		api.AssertIsEqual(circuit.HashDigests[i], computedHash2)

	}

	return nil
}

/*
PrivateWitnessGeneration generates the private witness
*/
func PrivateWitnessGenerationV4(pubParams PublicWitnessParametersV4, privParams PrivateWitnessParametersV4) witness.Witness {
	assignment := &CircuitV4{
		Seed:              frontend.Variable(privParams.Seed),
		CurveID:           tedwards.BN254,
		Claims:            make([]frontend.Variable, len(privParams.Claims)),
		HashDigests:       make([]frontend.Variable, len(pubParams.HashDigests)),
		Epochs:            make([]frontend.Variable, len(pubParams.Epochs)),
		Tokens:            make([]frontend.Variable, len(pubParams.Tokens)),
		Challenge:         frontend.Variable(pubParams.Challenge),
		Holder_randomness: frontend.Variable(privParams.HolderRadomness),
		ValidUntil:        frontend.Variable(pubParams.ValidUntil),
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Epochs[i] = frontend.Variable(pubParams.Epochs[i])
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Tokens[i] = frontend.Variable(pubParams.Tokens[i])
	}

	for i := 0; i < len(privParams.Claims); i++ {
		assignment.Claims[i] = frontend.Variable(privParams.Claims[i])
	}

	for i := 0; i < len(pubParams.HashDigests); i++ {
		assignment.HashDigests[i] = frontend.Variable(pubParams.HashDigests[i])
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
func PublicWitnessGenerationV4(pubParams PublicWitnessParametersV4) witness.Witness {
	assignment := &CircuitV4{
		CurveID:     tedwards.BN254,
		Epochs:      make([]frontend.Variable, len(pubParams.Epochs)),
		Tokens:      make([]frontend.Variable, len(pubParams.Tokens)),
		Challenge:   frontend.Variable(pubParams.Challenge),
		HashDigests: make([]frontend.Variable, len(pubParams.HashDigests)),
		ValidUntil:  frontend.Variable(pubParams.ValidUntil),
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Epochs[i] = frontend.Variable(pubParams.Epochs[i])
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Tokens[i] = frontend.Variable(pubParams.Tokens[i])
	}

	for i := 0; i < len(pubParams.HashDigests); i++ {
		assignment.HashDigests[i] = frontend.Variable(pubParams.HashDigests[i])
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
