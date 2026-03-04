package zkp

import (
	"bytes"
	"encoding/binary"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zkrevoke/crypto2"
)

func GenerateToken(seed []byte, epoch []byte) []byte {
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write(seed)
	_, _ = f.Write(epoch)
	return f.Sum(nil)
}

func GenerateDigestOfClaim(claim []byte, challenge []byte, holder_randomness []byte) []byte {
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write(claim)
	_, _ = f.Write(challenge)
	_, _ = f.Write(holder_randomness)
	return f.Sum(nil)
}

/*
Generates the following inputs:  seed, epochs, tokens, challenge, holder_randomness, claims, hash_digets, valid_until
*/
func Generate_Inputs_V4(n int, l int) ([]byte, [][]byte, [][]byte, []byte, []byte, [][]byte, [][]byte, []byte) {

	seed := []byte("seed#fdsd")
	challenge := []byte("challenge#123")
	holder_randomness := []byte("holder#123")
	valid_until := []byte("expiration_period#123")
	var epochs [][]byte
	var tokens [][]byte
	var claims [][]byte
	var hash_digets [][]byte
	for i := 0; i < n; i++ {
		epoch := []byte(strconv.Itoa(i))
		epochs = append(epochs, epoch)
		token := GenerateToken(seed, epoch)
		tokens = append(tokens, token)
	}

	for i := 0; i < l; i++ {
		claim := []byte("claim: " + strconv.Itoa(i+1))
		claims = append(claims, claim)
		digest := GenerateDigestOfClaim(claim, challenge, holder_randomness)
		hash_digets = append(hash_digets, digest)
	}
	zap.S().Infoln("\t seed: ", seed, "\t challenge: ", challenge, "\t epochs: ", epochs)
	zap.S().Infoln("tokens: H(seed || epoch)): ", tokens)
	zap.S().Infoln("claims: ", claims)

	//var msg []byte
	//msg = append(msg, vc_id...)
	//msg = append(msg, seed...)

	return seed, epochs, tokens, challenge, holder_randomness, claims, hash_digets, valid_until
}

/*
Contains all the workflow and interactions.
*/
func Test_Circuit_V4() {
	n := 1
	l := 3
	// The following elements are encoded in a VC
	seed, epochs, tokens, challenge, holder_randomness, claims, hash_digests, valid_until := Generate_Inputs_V4(n, l)

	// private key - used by the issuer to sign VCs
	// public key - used by the verifier to verify VCs
	privateKey, publicKey := crypto2.Generate_EDDSA_Keypairs()
	//privateKey_Holder, publicKey_Holder := crypto.Generate_EDDSA_Keypairs()
	zap.S().Infoln("private key size: ", binary.Size(privateKey))
	zap.S().Infoln("public key size: ", binary.Size(publicKey))

	//x_point, y_point := crypto.Parse_PublicKey(publicKey_Holder.Bytes())
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write([]byte(seed))
	for i := 0; i < l; i++ {
		_, _ = f.Write([]byte(claims[i]))
	}
	_, _ = f.Write([]byte(valid_until))
	msg := f.Sum(nil)

	signature, _ := crypto2.Sign_EDDSA(privateKey, msg)

	//signature_Holder, _ := crypto.Sign_EDDSA(privateKey_Holder, msg2)
	ccs := NewCircuitv4(n, l)

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)
	zap.S().Infoln("proving key size: ", binary.Size(buf.Bytes()), "\t len: ", buf.Len())
	buf = new(bytes.Buffer)
	vk.WriteTo(buf)
	zap.S().Infoln("verification key size: ", binary.Size(buf.Bytes()), "len: ", buf.Len())
	zap.S().Infoln("****Issuer****: setup success")

	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := PrivateWitnessParametersV4{
		Seed:            seed,
		HolderRadomness: holder_randomness,
		Signature:       signature,
		Claims:          claims,
	}
	pubParams := PublicWitnessParametersV4{
		Epochs:      epochs,
		Tokens:      tokens,
		PublicKey:   publicKey,
		Challenge:   challenge,
		HashDigests: hash_digests,
		ValidUntil:  valid_until,
	}
	witness := PrivateWitnessGenerationV4(pubParams, privParams)
	zap.S().Infoln("witness: ", witness)
	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	zap.S().Infoln("****Holder****: witness:  size: ", binary.Size(buf.Bytes()), "\t time (in micro seconds): ", end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	zap.S().Infoln("****Holder****: proof:  size: ", binary.Size(buf.Bytes()), "\t time (in micro seconds): ", end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := PublicWitnessGenerationV4(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness size: ", binary.Size(buf.Bytes()), "\t time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t time (in micro seconds): ", end.Microseconds())
}
