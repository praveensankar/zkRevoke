package crypto2

import (
	"encoding/binary"
	"encoding/hex"
	"go.uber.org/zap"
)

func Test_EDDSA() {

	vc_id := []byte("id:2")
	seed := []byte("seed#fdsd")
	var msg []byte
	msg = append(msg, vc_id...)
	msg = append(msg, seed...)
	privateKey, publicKey := Generate_EDDSA_Keypairs()

	privKeyBytes := privateKey.Bytes()
	pubKeyBytes := publicKey.Bytes()
	zap.S().Infoln("eddsa private key hex:", hex.EncodeToString(privKeyBytes))
	zap.S().Infoln("eddsa public key hex:", hex.EncodeToString(pubKeyBytes))
	zap.S().Infoln("private key size: ", binary.Size(privateKey))
	zap.S().Infoln("public key size: ", binary.Size(publicKey))
	signature, _ := Sign_EDDSA(privateKey, msg)
	zap.S().Infoln("****CRYPTO EDDSA**** \t signature:", signature)
}

func Test_BBS() {
	bbsPrivateKey, bbsPublicKey := Generate_BBS_KeyPair()
	bbsPrivateKeyBytes, _ := bbsPrivateKey.Marshal()
	bbsPublicKeyBytes, _ := bbsPublicKey.Marshal()

	zap.S().Infoln("bbs private key hex:", hex.EncodeToString(bbsPrivateKeyBytes))
	zap.S().Infoln("bbs public key hex:", hex.EncodeToString(bbsPublicKeyBytes))
	zap.S().Infoln("bbs public key size (in bytes): ", binary.Size(bbsPublicKeyBytes))
}

func Test_LoadKeys() {

	eddsaPrivateKey, _ := hex.DecodeString("e9370d0522d74b1dc8915b65886028240f4e70edfdd5600af3e5caed28e49c8279a5703413d4dcdd99b5f8b1759c19f7a86e474eb3d0eb6c72cb45b272c41eb0b093d80d38356e4bd1fd955d5a0c7b7ba86f00d7bb7112d14c6cf0a9be3009da")
	bbsPrivateKey, _ := hex.DecodeString("4328f1116c1282298667d79ec0a2a2ebabc9ef3c008bbbb336eee0c966b3e6cc")

	eddsaSK, eddsaPK := BytesToEDDSAKeys(eddsaPrivateKey)
	bbsSK, bbsPK := BytesToBBSPrivateKey(bbsPrivateKey)

	eddsaSKBytes := eddsaSK.Bytes()
	eddsaPKBytes := eddsaPK.Bytes()
	zap.S().Infoln("eddsa private key hex:", hex.EncodeToString(eddsaSKBytes))
	zap.S().Infoln("eddsa public key hex:", hex.EncodeToString(eddsaPKBytes))

	bbsPrivateKeyBytes, _ := bbsSK.Marshal()
	bbsPublicKeyBytes, _ := bbsPK.Marshal()

	zap.S().Infoln("bbs private key hex:", hex.EncodeToString(bbsPrivateKeyBytes))
	zap.S().Infoln("bbs public key hex:", hex.EncodeToString(bbsPublicKeyBytes))

	vc_id := []byte("id:2")
	seed := []byte("seed#fdsd")
	var msg []byte
	msg = append(msg, vc_id...)
	msg = append(msg, seed...)

}
