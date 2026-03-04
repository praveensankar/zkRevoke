package crypto2

import (
	"crypto/rand"
	"github.com/suutaku/go-bbs/pkg/bbs"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

type BBS struct {
	PublicKey  *bbs.PublicKey
	PrivateKey *bbs.PrivateKey
}

func Generate_BBS_KeyPair() (*bbs.PrivateKey, *bbs.PublicKey) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		zap.S().Infoln("BBS - error while generating random string: %s", err)
	}
	//zap.S().Infoln("BBS - key pair generation: seed", seed)
	publicKey, privateKey, err := bbs.GenerateKeyPair(sha3.New512, seed)
	if err != nil {
		zap.S().Infoln("BBS - error creating new key pair: ", err)
	}
	return privateKey, publicKey

}

func Sign_BBS(privateKey *bbs.PrivateKey, messages [][]byte) []byte {
	bbsInstance := bbs.NewBbs()
	signature, err := bbsInstance.SignWithKey(messages, privateKey)
	if err != nil {
		zap.S().Infoln("BBS - error signing: ", err)
	}
	return signature
}

func Verify_BBS(publicKey []byte, signature []byte, messages [][]byte) bool {

	bbsInstance := &bbs.Bbs{}
	err := bbsInstance.Verify(messages, signature, publicKey)
	if err != nil {
		zap.S().Infoln("BBS - verification failed: ", err)
		return false
	}

	//zap.S().Infoln("BBS - digital signature verification successful")
	return true
}

/*
Generate_BBS_Proof function generates proof for selective disclosure

Input:

	publicKey: public key
	signature : digital signature of the complete messages
	messages: all the messages
	revealedIndexes: list of indexes that need to be revealed

Output:

	(proof)[]byte - bbs proof
	(nonce)[]byte - nonce used for unlinkability
*/
func Generate_BBS_Proof(publicKey []byte, signature []byte, messages [][]byte, revealedIndexes []int) ([]byte, []byte) {

	nonce := make([]byte, 32)
	_, err := rand.Read(nonce)
	if err != nil {
		zap.S().Infoln("BBS - error while generating random string: %s", err)
	}
	bbsInstance := bbs.NewBbs()
	//pk , _ := bbs.UnmarshalPublicKey(publicKey)
	//zap.S().Infoln("BBS - Selective disclosure - public key: ", pk)
	proof, err := bbsInstance.DeriveProof(messages, signature, nonce, publicKey, revealedIndexes)
	if err != nil {
		zap.S().Infoln("BBS - error creating proof for selective disclosure: ", err)
	}
	return proof, nonce
}

func Verify_BBS_Proof(publicKey []byte, proof []byte, selectiveMessages [][]byte, nonce []byte) bool {

	bbsInstance := bbs.NewBbs()

	err := bbsInstance.VerifyProof(selectiveMessages, proof, nonce, publicKey)
	if err != nil {
		zap.S().Infoln("BBS - selective disclosure verification failed: ", err)
		return false
	}

	//zap.S().Infoln("BBS - selective disclosure verification successful")
	return true
}

/*
PublicKeyToString function converts public key to byte

Input:

	publicKey

Returns:

	byte (public key)
*/
func BBSPublicKey_To_Bytes(publicKey *bbs.PublicKey) []byte {
	res, err := publicKey.Marshal()
	if err != nil {
		zap.S().Infoln("BBS - error marshing public key")
	}

	//zap.S().Infoln("BBS - public key byte: ",res)
	return res
}

func BBSPrivateKeyToBytes(privateKey *bbs.PrivateKey) []byte {
	res, err := privateKey.Marshal()
	if err != nil {
		zap.S().Infoln("BBS - error marshalling private key")
	}
	return res
}

func BytesToBBSPrivateKey(data []byte) (*bbs.PrivateKey, *bbs.PublicKey) {
	sk, _ := bbs.UnmarshalPrivateKey(data)
	pk := sk.PublicKey()
	return sk, pk
}
