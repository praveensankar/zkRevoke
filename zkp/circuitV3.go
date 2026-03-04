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

type Circuit struct {
	CurveID tedwards.ID `gnark:",public"`

	Seed   frontend.Variable   // seed is the private input known only to the holder and the issuer
	Epochs []frontend.Variable `gnark:",public"` // Epoch indicates the current epoch.  The number of durations elapsed from the issuance of VC.
	Tokens []frontend.Variable `gnark:",public"` // token = H(Seed || Epoch). token is inserted into a revocation list.

	ClaimsHash frontend.Variable   `gnark:",public"` // hash of claims
	ValidUntil frontend.Variable   `gnark:",public"` // ValidUntil is the expiration date of the VC.
	PublicKey  eddsa_zkp.PublicKey `gnark:",public"` // public key of the issuer
	Signature  eddsa_zkp.Signature // Signature ← Sign(PublicKey,  H(seed || h(claims) || validUntil))

	Challenge       frontend.Variable `gnark:",public"` // challenge is provided by verifier when requesting VP from a holder
	HolderRadomness frontend.Variable // random value used by holder
	Hash1           frontend.Variable `gnark:",public"` // H(Challenge || HolderRadomness)
}

func NewCircuit(n int) constraint.ConstraintSystem {

	//logger.Disable()
	var circuit Circuit
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

	**** Verify(PK, H(seed || H(claims) || validUntil)) == sig
	**** H(challenge||holder_randomness)) == Hash1
	**** H(Seed || Epoch[i]) == token[i]: i = 1...numberOfTokens

Verify(PK, H(seed || claims || validUntil)) == sig:
-------------- verifies that the issuer's signature is present on vc_id, seed and expiry date of VC
H(challenge||holder_randomness)) == Hash1:
-------------- verifies that the hash1 is computed from the challenge issued by a verifier and a secret randomness provided by the holder
H(Seed || Epoch) == token:
-------------- verifies that the token is generated from the secret seed and an epoch value
*/

func (circuit *Circuit) Define(api frontend.API) error {
	curve, err := twistededwards_zkp.NewEdCurve(api, circuit.CurveID)
	if err != nil {
		return err
	}

	numberOfTokens := len(circuit.Tokens)

	//**** Verify(PK, H(seed || claims || validUntil)) == sig
	mimc1, _ := mimc.NewMiMC(api)
	mimc1.Write(circuit.Seed)
	mimc1.Write(circuit.ClaimsHash)
	mimc1.Write(circuit.ValidUntil)
	computedHash := mimc1.Sum()
	mimc2, _ := mimc.NewMiMC(api)
	err = eddsa_zkp.Verify(curve, circuit.Signature, computedHash, circuit.PublicKey, &mimc2)

	//**** H(challenge||holder_randomness)) == Hash1
	mimc3, err := mimc.NewMiMC(api)
	mimc3.Write(circuit.Challenge, circuit.HolderRadomness)
	computedHash2 := mimc3.Sum()
	api.AssertIsEqual(circuit.Hash1, computedHash2)

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
func PrivateWitnessGeneration(pubParams PublicWitnessParameters, privParams PrivateWitnessParameters) witness.Witness {
	assignment := &Circuit{
		Seed:            frontend.Variable(privParams.Seed),
		HolderRadomness: frontend.Variable(privParams.HolderRadomness),
		CurveID:         tedwards.BN254,
		Epochs:          make([]frontend.Variable, len(pubParams.Epochs)),
		Tokens:          make([]frontend.Variable, len(pubParams.Tokens)),
		Challenge:       frontend.Variable(pubParams.Challenge),
		Hash1:           frontend.Variable(pubParams.Hash1),
		ClaimsHash:      frontend.Variable(pubParams.ClaimsHash),
		ValidUntil:      frontend.Variable(pubParams.ValidUntil),
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Epochs[i] = frontend.Variable(pubParams.Epochs[i])
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Tokens[i] = frontend.Variable(pubParams.Tokens[i])
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
func PublicWitnessGeneration(pubParams PublicWitnessParameters) witness.Witness {
	assignment := &Circuit{
		CurveID:    tedwards.BN254,
		Epochs:     make([]frontend.Variable, len(pubParams.Epochs)),
		Tokens:     make([]frontend.Variable, len(pubParams.Tokens)),
		Challenge:  frontend.Variable(pubParams.Challenge),
		Hash1:      frontend.Variable(pubParams.Hash1),
		ValidUntil: frontend.Variable(pubParams.ValidUntil),
		ClaimsHash: frontend.Variable(pubParams.ClaimsHash),
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Epochs[i] = frontend.Variable(pubParams.Epochs[i])
	}

	for i := 0; i < len(pubParams.Epochs); i++ {
		assignment.Tokens[i] = frontend.Variable(pubParams.Tokens[i])
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
