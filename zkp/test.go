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

func MimcHash(seed []byte, epoch []byte) []byte {
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write(seed)
	_, _ = f.Write(epoch)
	return f.Sum(nil)
}

func ComputeHashOnClaims(claims [][]byte) []byte {
	f := bn254_mimc.NewMiMC()
	for i := 0; i < len(claims); i++ {
		_, _ = f.Write(claims[i])
	}
	return f.Sum(nil)
}

/*
Generates the following inputs: vc_id, seed, epoch, token, msg
*/
func Generate_Inputs(n int) ([]byte, []byte, [][]byte, [][]byte, []byte, []byte) {
	vc_id := []byte("id:2")
	seed := []byte("seed#fdsd")
	challenge := []byte("challenge#123")
	holder_randomness := []byte("holder#123")

	var epochs [][]byte
	var tokens [][]byte

	for i := 0; i < n; i++ {
		epoch := []byte(strconv.Itoa(i))
		epochs = append(epochs, epoch)
		token := MimcHash(seed, epoch)
		tokens = append(tokens, token)
	}

	//zap.S().Infoln("vc id: ", vc_id, "\t seed: ", seed, "\t challenge: ", challenge, "\t epochs: ", epochs)
	//zap.S().Infoln("tokens: H(seed || epoch)): ", tokens)

	//var msg []byte
	//msg = append(msg, vc_id...)
	//msg = append(msg, seed...)

	return vc_id, seed, epochs, tokens, challenge, holder_randomness
}

/*
Contains all the workflow and interactions.
*/
func Test_Circuit() {

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
	zap.S().Infoln("private key size: ", binary.Size(privateKey))
	zap.S().Infoln("public key size: ", binary.Size(publicKey))

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
	zap.S().Infoln("proving key size: ", binary.Size(buf.Bytes()), "\t len: ", buf.Len())
	buf = new(bytes.Buffer)
	vk.WriteTo(buf)
	zap.S().Infoln("verification key size: ", binary.Size(buf.Bytes()), "len: ", buf.Len())
	zap.S().Infoln("****Issuer****: setup success")

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
	publicWitness := PublicWitnessGeneration(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	zap.S().Infoln("****Verifier****: public witness size: ", binary.Size(buf.Bytes()), "\t time (in micro seconds): ", end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	zap.S().Infoln("****Verifier****: Verification status: ", status, "\t time (in micro seconds): ", end.Microseconds())
	zap.S().Infoln("************************************************")
}
