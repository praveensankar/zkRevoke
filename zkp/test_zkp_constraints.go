package zkp

import (
	"bytes"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"go.uber.org/zap"
	"time"
	"zkrevoke/crypto2"
)

func TestEmptyCircuit() {

	zap.S().Infoln("************************ Empty Circuit ************************")
	ccs := NewCircuitEmpty()

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	zap.S().Infoln("number of constraints in the circuit: ", ccs.GetNbConstraints())
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()

	witness := PrivateWitnessGenerationForEmptyCircuit()

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	zap.S().Infoln("****Holder****: private witness generation time (in micro seconds): ", end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	zap.S().Infoln("****Holder****: proof generation time (in micro seconds): ", end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := PublicWitnessGenerationForEmptyCircuit()
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness generation time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t verification time (in micro seconds): ", end.Microseconds())
	zap.S().Infoln("************************************************")
}

func TestSignatureVerificationCircuit() {

	zap.S().Infoln("************************ Signature Verification Circuit ************************")
	// The following elements are encoded in a VC
	validUntil, seed, _, _, _, _ := Generate_Inputs(1)

	var claims [][]byte
	claims = append(claims, []byte("name:bob"))
	claims = append(claims, []byte("age:25"))
	claims = append(claims, []byte("salary:10000"))
	claims = append(claims, []byte("role:CTO"))
	claims = append(claims, []byte("employer:E1"))
	// private key - used by the issuer to sign VCs
	// public key - used by the verifier to verify VCs
	privateKey, publicKey := crypto2.Generate_EDDSA_Keypairs()
	//privateKey_Holder, publicKey_Holder := crypto.Generate_EDDSA_Keypairs()

	//x_point, y_point := crypto.Parse_PublicKey(publicKey_Holder.Bytes())
	f := bn254_mimc.NewMiMC()

	_, _ = f.Write([]byte(seed))
	_, _ = f.Write(ComputeHashOnClaims(claims))
	_, _ = f.Write([]byte(validUntil))
	msg := f.Sum(nil)

	signature, _ := crypto2.Sign_EDDSA(privateKey, msg)

	//signature_Holder, _ := crypto.Sign_EDDSA(privateKey_Holder, msg2)
	ccs := NewCircuitSignatureVerification()

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	zap.S().Infoln("number of constraints in the circuit: ", ccs.GetNbConstraints())
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := PrivateWitnessParameters{
		Seed:      seed,
		Signature: signature,
	}
	pubParams := PublicWitnessParameters{
		PublicKey:  publicKey,
		ValidUntil: validUntil,
		ClaimsHash: ComputeHashOnClaims(claims),
	}
	witness := PrivateWitnessGenerationForCircuitSignatureVerification(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	zap.S().Infoln("****Holder****: private witness generation time time (in micro seconds): ", end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	zap.S().Infoln("****Holder****: proof generation time (in micro seconds): ", end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := PublicWitnessGenerationForCircuitSignatureVerification(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness generation time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t time (in micro seconds): ", end.Microseconds())
	zap.S().Infoln("************************************************")
}

func TestTokenVerificationCircuit() {

	zap.S().Infoln("************************ Token Verification Circuit ************************")
	n := 1
	// The following elements are encoded in a VC
	_, seed, epochs, tokens, _, _ := Generate_Inputs(1)

	var claims [][]byte
	claims = append(claims, []byte("name:bob"))
	claims = append(claims, []byte("age:25"))
	claims = append(claims, []byte("salary:10000"))
	claims = append(claims, []byte("role:CTO"))
	claims = append(claims, []byte("employer:E1"))
	// private key - used by the issuer to sign VCs
	// public key - used by the verifier to verify VCs

	//signature_Holder, _ := crypto.Sign_EDDSA(privateKey_Holder, msg2)
	ccs := NewCircuitForTokenVerification(n)

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)
	zap.S().Infoln("number of constraints in the circuit: ", ccs.GetNbConstraints())
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := PrivateWitnessParameters{
		Seed: seed,
	}
	pubParams := PublicWitnessParameters{
		Epochs:     epochs,
		Tokens:     tokens,
		ClaimsHash: ComputeHashOnClaims(claims),
	}
	witness := PrivateWitnessGenerationForTokenVerificationCircuit(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	zap.S().Infoln("****Holder****: private witness generation time time (in micro seconds): ", end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	zap.S().Infoln("****Holder****: proof generation time (in micro seconds): ", end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := PublicWitnessGenerationForTokenVerificationCircuit(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness generation time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t time (in micro seconds): ", end.Microseconds())
	zap.S().Infoln("************************************************")
}

func TestChallengeVerificationCircuit() {

	zap.S().Infoln("************************ Challenge Verification Circuit ************************")

	// The following elements are encoded in a VC
	_, _, _, _, challenge, holder_randomness := Generate_Inputs(1)

	var claims [][]byte
	claims = append(claims, []byte("name:bob"))
	claims = append(claims, []byte("age:25"))
	claims = append(claims, []byte("salary:10000"))
	claims = append(claims, []byte("role:CTO"))
	claims = append(claims, []byte("employer:E1"))

	f2 := bn254_mimc.NewMiMC()
	_, _ = f2.Write(challenge)
	_, _ = f2.Write(holder_randomness)
	msg2 := f2.Sum(nil)

	ccs := NewCircuitForChallengeVerification()

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	zap.S().Infoln("number of constraints in the circuit: ", ccs.GetNbConstraints())

	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := PrivateWitnessParameters{
		HolderRadomness: holder_randomness,
	}
	pubParams := PublicWitnessParameters{
		Hash1:      msg2,
		Challenge:  challenge,
		ClaimsHash: ComputeHashOnClaims(claims),
	}
	witness := PrivateWitnessGenerationForChallengeVerification(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	zap.S().Infoln("****Holder****: private witness generation time time (in micro seconds): ", end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	zap.S().Infoln("****Holder****: proof generation time (in micro seconds): ", end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := PublicWitnessGenerationForChallengeVerification(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness generation time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t time (in micro seconds): ", end.Microseconds())
	zap.S().Infoln("************************************************")
}

func TestCompleteCircuit() {

	zap.S().Infoln("************************ Complete Circuit ************************")
	n := 1
	// The following elements are encoded in a VC
	validUntil, seed, epochs, tokens, challenge, holder_randomness := Generate_Inputs(1)

	var claims [][]byte
	claims = append(claims, []byte("name:bob"))
	claims = append(claims, []byte("age:25"))
	claims = append(claims, []byte("salary:10000"))
	claims = append(claims, []byte("role:CTO"))
	claims = append(claims, []byte("employer:E1"))
	// private key - used by the issuer to sign VCs
	// public key - used by the verifier to verify VCs
	privateKey, publicKey := crypto2.Generate_EDDSA_Keypairs()
	//privateKey_Holder, publicKey_Holder := crypto.Generate_EDDSA_Keypairs()

	//x_point, y_point := crypto.Parse_PublicKey(publicKey_Holder.Bytes())
	f := bn254_mimc.NewMiMC()

	_, _ = f.Write([]byte(seed))
	_, _ = f.Write(ComputeHashOnClaims(claims))
	_, _ = f.Write([]byte(validUntil))
	msg := f.Sum(nil)

	signature, _ := crypto2.Sign_EDDSA(privateKey, msg)

	f2 := bn254_mimc.NewMiMC()
	_, _ = f2.Write(challenge)
	_, _ = f2.Write(holder_randomness)
	msg2 := f2.Sum(nil)

	//signature_Holder, _ := crypto.Sign_EDDSA(privateKey_Holder, msg2)
	ccs := NewCircuit(n)

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	zap.S().Infoln("number of constraints in the circuit: ", ccs.GetNbConstraints())
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := PrivateWitnessParameters{
		Seed:            seed,
		HolderRadomness: holder_randomness,
		Signature:       signature,
	}
	pubParams := PublicWitnessParameters{
		PublicKey:  publicKey,
		Hash1:      msg2,
		Challenge:  challenge,
		Epochs:     epochs,
		ValidUntil: validUntil,
		Tokens:     tokens,
		ClaimsHash: ComputeHashOnClaims(claims),
	}
	witness := PrivateWitnessGeneration(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	zap.S().Infoln("****Holder****: private witness generation time time (in micro seconds): ", end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	zap.S().Infoln("****Holder****: proof generation time (in micro seconds): ", end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := PublicWitnessGeneration(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness generation time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t time (in micro seconds): ", end.Microseconds())
	zap.S().Infoln("************************************************")
}

func TestConstraintsInCircuit() {
	for i := 1; i <= 5; i++ {
		//TestEmptyCircuit()
		//TestTokenVerificationCircuit()
		//TestChallengeVerificationCircuit()
		//TestSignatureVerificationCircuit()
		TestCompleteCircuit()
	}
}
