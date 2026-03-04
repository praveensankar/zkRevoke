package holder

import (
	"encoding/json"
	"strconv"
	"time"
	"zkrevoke/model"
	"zkrevoke/zkp"
)

func (holder *Holder) Generate_ZKP_Proof(epochs []uint, tokens [][]byte, challenge []byte, holder_randomness []byte, hash1 []byte, vc model.VerifiableCredential, claims_hash []byte) *model.Proof {

	start := time.Now()
	var eddsaSignature []byte

	if vc.Proofs[0].Type == string(model.ProofTypeEDDSA) {
		eddsaSignature = vc.Proofs[0].ProofValue
	}

	var epochBytes [][]byte
	for i := 0; i < len(epochs); i++ {
		epochBytes = append(epochBytes, []byte(strconv.Itoa(int(epochs[i]))))
	}

	var claims [][]byte

	jsonData, _ := json.Marshal(vc.CredentialSubject[0])
	var employment_claims model.EmploymentClaims
	json.Unmarshal(jsonData, &employment_claims)
	claims = append(claims, []byte(employment_claims.EmployeeID))
	claims = append(claims, []byte(employment_claims.EmployeeName))
	claims = append(claims, []byte(employment_claims.EmployerName))
	claims = append(claims, []byte(employment_claims.EmployeeDesignation))
	claims = append(claims, []byte(strconv.Itoa(employment_claims.Salary)))

	//zap.S().Infoln("***HOLDER****: eddsa signature: ", hex.EncodeToString(eddsaSignature))
	privParams := zkp.PrivateWitnessParameters{
		Signature:       eddsaSignature,
		Seed:            []byte(vc.Seed),
		HolderRadomness: holder_randomness,
	}
	pubParams := zkp.PublicWitnessParameters{
		PublicKey:  holder.eddsaPublicKey,
		Hash1:      hash1,
		Challenge:  challenge,
		Epochs:     epochBytes,
		ValidUntil: []byte(vc.Metadata.ValidUntil),
		Tokens:     tokens,
		ClaimsHash: claims_hash,
	}

	witness := zkp.PrivateWitnessGeneration(pubParams, privParams)
	proof := zkp.ProveGroth(holder.ccs, holder.zkpProvingKey, witness)
	//zap.S().Infoln("***HOLDER****: zkp proof: ", proof)
	proof_in_bytes, _ := zkp.GrothProofToBytes(proof)
	proofSize := len(proof_in_bytes)
	vp_proof_zkp := &model.Proof{
		Type:         string(model.ProofTypeZKP),
		ProofPurpose: "ZKP Proof",
		Cryptosuite:  "zkp",
		Created:      "",
		Expires:      "",
		Nonce:        "",
		ProofValue:   proof_in_bytes,
	}
	end := time.Since(start)
	holder.Result.AddGrothProofGenerationTime(end)
	holder.Result.AddGrothProofSize(proofSize)
	return vp_proof_zkp
}
