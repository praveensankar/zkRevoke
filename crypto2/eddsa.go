package crypto2

import (
	"crypto/rand"
	"fmt"
	edwardsbn254 "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark-crypto/signature/eddsa"
	"go.uber.org/zap"
)

type EdDSA struct {
	PrivateKey signature.Signer
	publicKey  signature.PublicKey
}

/*
Generates new private key and public key pair for ZKP
*/
func Generate_EDDSA_Keypairs() (signature.Signer, signature.PublicKey) {
	randomness := rand.Reader
	privateKey, err := eddsa.New(tedwards.BN254, randomness)
	publicKey := privateKey.Public()
	if err != nil {
		fmt.Println("error in generating private key: ", err)
	}
	return privateKey, publicKey
}

/*
Sign_EDDSA signs messages using EDDSA signature algorithm
Inputs:
1) (privateKey) signature.Signer
2) (msg) []byte

Returns:
1) (signature) []byte
2) error
*/
func Sign_EDDSA(privateKey signature.Signer, msg []byte) ([]byte, error) {
	hFunc := hash.MIMC_BN254.New()

	// Signature is encoded in a VC
	signature, err := privateKey.Sign(msg, hFunc)
	if err != nil {
		zap.S().Infoln("signature generation failed: ", err)
	}

	//zap.S().Infoln("****EDDSA****:  signature: ", signature)
	publicKey := privateKey.Public()
	// checks whether the generated is correct or not
	isValid, err := publicKey.Verify(signature, msg, hFunc)
	if isValid == false {
		zap.S().Infoln("signature verification failed")
	}
	return signature, nil
}

/*
Verify_EDDSA verifies the digial signature using eddsa algorithm
Returns
True - if signature is valid
False - otherwise
*/
func Verify_EDDSA(publicKey signature.PublicKey, msg []byte, signature []byte) bool {
	// checks whether the generated is correct or not
	hFunc := hash.MIMC_BN254.New()
	isValid, err := publicKey.Verify(signature, msg, hFunc)
	if err != nil {
		zap.S().Infoln("signature verification failed: ", err)
	}
	if isValid == true {
		//zap.S().Infoln("****EDDSA****: signature verification successful")
	} else {
		zap.S().Infoln("****EDDSA****: signature verification failed")
	}
	return isValid
}

// Parse_PublicKey parses a compressed binary point into uncompressed P.X and P.Y
func Parse_PublicKey(publickey []byte) ([]byte, []byte) {

	var pointbn254 edwardsbn254.PointAffine
	if _, err := pointbn254.SetBytes(publickey[:32]); err != nil {
		return nil, nil
	}
	a := pointbn254.X.Bytes()
	b := pointbn254.Y.Bytes()
	return a[:], b[:]
}

func BytesToEDDSAKeys(secret_key_bytes []byte) (signature.Signer, signature.PublicKey) {
	randomness := rand.Reader
	privateKey, _ := eddsa.New(tedwards.BN254, randomness)
	privateKey.SetBytes(secret_key_bytes)
	return privateKey, privateKey.Public()
}
