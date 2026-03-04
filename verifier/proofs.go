package verifier

import (
	"zkrevoke/zkp"
)

func (verifier Verifier) VerifyZKPProofInVP(epochs [][]byte, tokens [][]byte, proof []byte, challenge []byte, hash1 []byte, validUntil []byte, claims_hash []byte) bool {

	zkpVerifyingKey := verifier.zkpVerifyingKey
	eddsaPublicKey := verifier.eddsaPublicKey
	pubParams := zkp.PublicWitnessParameters{
		PublicKey:  eddsaPublicKey,
		Hash1:      hash1,
		Challenge:  challenge,
		Epochs:     epochs,
		ValidUntil: validUntil,
		Tokens:     tokens,
		ClaimsHash: claims_hash,
	}
	publicWitness := zkp.PublicWitnessGeneration(pubParams)
	zkpProof, _ := zkp.BytesToGrothProof(proof)
	//zap.S().Infoln("****VERIFIER****: zkp proof: ", hex.EncodeToString(proof))
	status := zkp.VerifyGroth(zkpProof, zkpVerifyingKey, publicWitness)

	return status
}
