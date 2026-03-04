/*
Consists of following functions:
1) GenerateVP(numberOfEpochs int, vc model.VerifiableCredential) *model.VerifiablePresentation
*/
package holder

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"strconv"
	"time"
	"zkrevoke/model"
	"zkrevoke/results"
	"zkrevoke/utils"
)

/*
GenerateVP function generates a new VP
vp, zkp_proof_generation_times, zkp_proof_sizes

Inputs:
(numberOfEpochs) int - holder can compute tokens for numberOfEpochs starting from the current epoch
(challenge) []byte - challenge provided by verifier

Returns:
1) (*model.VerifiablePresentation): newly generated vp
*/
func (holder *Holder) GenerateVP(numberOfEpochs int, challenge []byte) *model.VerifiablePresentation {

	epoch := holder.GetCurrentEpoch()

	vpGenStart := time.Now()
	vc := holder.verfiableCredentials[0]
	claimSet := vc.CredentialSubject
	claims := claimSet[0].(model.EmploymentClaims)

	var vp model.VerifiablePresentation

	if holder.SelectiveDisclosureExtension == false {
		vp.Messages = claims
	} else {
		var selectClaims model.SampleEmploymentProofPresentation
		selectClaims.EmployeeDesignation = claims.EmployeeDesignation
		selectClaims.Salary = claims.Salary

		// index starts with 0
		sde := model.SelectiveDisclosureElements{
			SelectClaims: selectClaims,
			Indexes:      []int{3, 4},
		}
		vp.SelectiveDisclosureElements = sde
	}

	vp.ValidFrom = vc.Metadata.ValidFrom
	vp.ValidUntil = vc.Metadata.ValidUntil

	holder_randomness := rand.Text()
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write(challenge)
	_, _ = f.Write([]byte(holder_randomness))

	msg := f.Sum(nil)
	vp.Hash1 = msg

	var claims_hash []byte

	if holder.SelectiveDisclosureExtension == true {

		var individual_hashes [][]byte

		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployeeID))
		hash1 := f.Sum(nil)
		individual_hashes = append(individual_hashes, hash1)

		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployeeName))
		hash2 := f.Sum(nil)
		individual_hashes = append(individual_hashes, hash2)

		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployerName))
		hash3 := f.Sum(nil)
		individual_hashes = append(individual_hashes, hash3)

		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployeeDesignation))
		hash4 := f.Sum(nil)
		individual_hashes = append(individual_hashes, hash4)

		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(strconv.Itoa(claims.Salary)))
		hash5 := f.Sum(nil)
		individual_hashes = append(individual_hashes, hash5)

		f = bn254_mimc.NewMiMC()
		_, _ = f.Write(hash1)
		_, _ = f.Write(hash2)
		_, _ = f.Write(hash3)
		_, _ = f.Write(hash4)
		_, _ = f.Write(hash5)
		claims_hash = f.Sum(nil)
		vp.ClaimsHash = claims_hash
		vp.Messages = individual_hashes

	} else {
		g := bn254_mimc.NewMiMC()
		_, _ = g.Write([]byte(claims.EmployeeID))
		_, _ = g.Write([]byte(claims.EmployeeName))
		_, _ = g.Write([]byte(claims.EmployerName))
		_, _ = g.Write([]byte(claims.EmployeeDesignation))
		_, _ = g.Write([]byte(strconv.Itoa(claims.Salary)))
		claims_hash = g.Sum(nil)
		vp.ClaimsHash = claims_hash
	}
	//eddsa_signature, err := crypto.Sign_EDDSA(holder.holder_PrivateKey, []byte(msg))
	//vp.Holder_randomness = []byte(holder_randomness)
	//vp.Holder_signature = eddsa_signature
	//if err != nil {
	//	zap.S().Fatalln("***HOLDER***: failed to sign verifier's challenge and holder's randomness")
	//}

	numberOfTokens := 0
	done := false

	// compute tokens and zkp proofs for future epochs
	singleZKPProofGenTimeTotal := 0
	for i := 0; i < numberOfEpochs; i++ {
		tokenPresentation := model.TokenPresentation{}
		if done == true {
			break
		}
		var epochs []uint
		var tokens [][]byte
		var tokensStr []string
		for j := 0; j < holder.NumberOfTokensInCircuit; j++ {
			token_gen_start := time.Now()
			token := utils.ComputeToken(epoch, vc.Seed)
			token_gen_end := time.Now()
			holder.Result.AddTokenGenerationTime(token_gen_end.Sub(token_gen_start))
			epochs = append(epochs, uint(epoch))
			tokens = append(tokens, token)
			tokensStr = append(tokensStr, utils.GetShortString(hex.EncodeToString(token)))
			//zap.S().Infoln("*****HOLDER******: epoch: ", epoch, "\t token: ", utils.GetShortString(hex.EncodeToString(token)))
			numberOfTokens++
			if numberOfTokens >= numberOfEpochs {
				done = true
			} else {
				epoch++
			}
		}

		startZKPProofGenTime := time.Now()
		zkp_proof := holder.Generate_ZKP_Proof(epochs, tokens, challenge, []byte(holder_randomness), msg, vc, claims_hash)
		endZKPProofGenTime := time.Since(startZKPProofGenTime)
		singleZKPProofGenTimeTotal += int(endZKPProofGenTime.Microseconds())
		for x := 0; x < len(epochs); x++ {
			tokenPresentation.Epochs = append(tokenPresentation.Epochs, uint(epochs[x]))
			tokenPresentation.Tokens = append(tokenPresentation.Tokens, tokens[x])
		}
		tokenPresentation.ZKPProof = *zkp_proof
		vp.TokenPresentations = append(vp.TokenPresentations, tokenPresentation)
	}

	singleZKPProofGenTime := singleZKPProofGenTimeTotal / len(vp.TokenPresentations)
	vpGenEnd := time.Now()
	vpGenTime := vpGenEnd.Sub(vpGenStart)

	vpJson := vp.Json()
	sizeOfTokenBlock, _ := json.Marshal(vp.TokenPresentations[0])
	vpRelatedMetrics := results.VPRelatedMetrics{
		NumberOfTokenBlocks:    len(vp.TokenPresentations),
		NumberOfTokensInaBlock: holder.NumberOfTokensInCircuit,
		SizeOfTokenBlock:       len(sizeOfTokenBlock),
		VPSize:                 len(vpJson),
		VPGenerationTime:       int(vpGenTime.Microseconds()),
		SingleZKPProofGenTime:  singleZKPProofGenTime,
		AllZKPProofsGenTime:    singleZKPProofGenTimeTotal,
	}

	holder.Result.AddVPRelatedMetrics(vpRelatedMetrics)

	return &vp
}
