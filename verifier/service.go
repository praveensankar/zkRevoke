package verifier

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"go.uber.org/zap"
	"slices"
	"strconv"
	"time"
	"zkrevoke/model"
	"zkrevoke/results"
	utils2 "zkrevoke/utils"
)

func (verifier *Verifier) RequestVP() string {
	verifier.Lock()
	challenge := rand.Text()
	verifier.VPCounter = verifier.VPCounter + 1
	verifier.Challenge[challenge] = verifier.VPCounter
	verifier.Unlock()
	return challenge
}
func (verifier *Verifier) ReceiveVP(challenge string, vp model.VerifiablePresentation) {
	verifier.Lock()
	verifier.verifiablePresentations = append(verifier.verifiablePresentations, vp)
	verifier.Unlock()
	verifier.VerifyVP(challenge, vp)
}

func (verifier *Verifier) VerifyVP(challenge string, vp model.VerifiablePresentation) bool {

	var res string
	res = res + "\n ****Verifier***** Received a new VP:"
	vp_ver_metrics := results.VPVerificationMetrics{}
	vp_ver_start := time.Now()

	var zkp_ver_times []int

	claims := vp.Messages
	var claims_hash []byte
	if verifier.SelectiveDisclosureExtension == false {
		revealedClaims := vp.Messages.(model.EmploymentClaims)
		claims_hash = revealedClaims.GenerateHashDigest()

		res = res + "\t { (claims- " + fmt.Sprintf("%v", claims) + " ),"
	} else {
		claims_hash = verifier.VerifySelectiveDisclosure(vp.Messages, vp.SelectiveDisclosureElements)
		claimsInVP := vp.SelectiveDisclosureElements.SelectClaims.(model.SampleEmploymentProofPresentation)
		employeeDesignation := claimsInVP.EmployeeDesignation
		salary := claimsInVP.Salary
		res = res + "\t { (claims- " + fmt.Sprintf("%v", employeeDesignation) + fmt.Sprintf("%v", salary) + " ),"
	}

	currentEpoch := verifier.GetCurrentEpoch()
	// check whether the first epoch included in the VP matches with the current epoch
	if int(vp.TokenPresentations[0].Epochs[0]) != currentEpoch {
		zap.S().Infoln("****VERIFIER****: wrong epoch in the VP: \t current epoch: ", currentEpoch,
			"\t starting epoch in VP: ", int(vp.TokenPresentations[0].Epochs[0]))
		return false
	}

	vpCounter := verifier.Challenge[challenge]
	numberOfTokenPresentations := len(vp.TokenPresentations)
	vp_ver_metrics.NumberOfTokenBlocksInVP = numberOfTokenPresentations
	number_of_tokens := 0

	res = res + "(number of token blocks: " + strconv.Itoa(numberOfTokenPresentations) + "),"
	for i := 0; i < numberOfTokenPresentations; i++ {

		epochSet := make(map[int]bool)
		var epochs [][]byte
		var tokensBytes [][]byte
		var tokensStr []string
		tokenStorages := make([]*TokenStorage, 0)
		for j := 0; j < len(vp.TokenPresentations[i].Tokens); j++ {
			number_of_tokens = number_of_tokens + 1
			epoch := vp.TokenPresentations[i].Epochs[j]

			if epochSet[int(epoch)] == false {
				if int(epoch) != currentEpoch {
					zap.S().Infoln("****VERIFIER****: wrong epoch in the VP: \t current epoch: ", currentEpoch,
						"\t starting epoch in VP: ", epoch)
					return false
				}
				currentEpoch++
			}

			token := vp.TokenPresentations[i].Tokens[j]

			vpExpiryDate, _ := strconv.Atoi(vp.ValidUntil)
			validFrom, _ := strconv.Atoi(vp.ValidFrom)
			//numberOfEpochsVCisValid := (vpExpiryDate - validFrom) / verifier.Duration
			numberOfEpochsVCisValid := utils2.GetNumberOfBlocksVCisValid(validFrom, vpExpiryDate)
			//zap.S().Infoln("Verifier: number of epochs vc is valid: ", numberOfEpochsVCisValid)

			tokenStorage := TokenStorage{
				Epoch:                   int(epoch),
				Token:                   token,
				VCExpiryDate:            int64(vpExpiryDate),
				VCValidFrom:             int64(validFrom),
				NumberOfEpochsVCisValid: numberOfEpochsVCisValid,
			}

			if epochSet[int(epoch)] == false {
				tokenStorages = append(tokenStorages, &tokenStorage)
				tokensStr = append(tokensStr, utils2.GetShortString(hex.EncodeToString(token)))
			}

			epochs = append(epochs, []byte(strconv.Itoa(int(epoch))))
			tokensBytes = append(tokensBytes, token)

			epochSet[int(epoch)] = true
		}
		zkpProof := vp.TokenPresentations[i].ZKPProof.ProofValue

		zkp_start := time.Now()
		zkpStatus := verifier.VerifyZKPProofInVP(epochs, tokensBytes, zkpProof, []byte(challenge), vp.Hash1, []byte(vp.ValidUntil), claims_hash)
		zkp_end := time.Since(zkp_start)
		zkp_ver_times = append(zkp_ver_times, int(zkp_end.Microseconds()))
		if zkpStatus == true {
			for _, tokenStorage := range tokenStorages {
				verifier.Tokens[vpCounter] = append(verifier.Tokens[vpCounter], tokenStorage)
			}
		} else {
			return false
		}
		//res = res + "\n (tokens: " + fmt.Sprintf("%v", tokensStr) + "), (number of tokens: " + strconv.Itoa(len(vp.TokenPresentations[i].Tokens))
		//res = res + "), (zkp status- " + fmt.Sprintf("%v", zkpStatus) + ")"
	}

	vpExpiryDate, _ := strconv.Atoi(vp.ValidUntil)
	validFrom, _ := strconv.Atoi(vp.ValidFrom)
	numberOfEpochsVCisValid := utils2.GetNumberOfBlocksVCisValid(validFrom, vpExpiryDate)
	res = res + " \t (expiry date in number of epochs: " + strconv.Itoa(numberOfEpochsVCisValid) + ",) \n "

	zap.S().Infoln(res)

	vp_ver_end := time.Now()

	avg_zkp_verification_time := 0
	for i := 0; i < len(zkp_ver_times); i++ {
		avg_zkp_verification_time += int(zkp_ver_times[i])
	}
	vp_ver_metrics.GrothProofVerificationTimeTotal = avg_zkp_verification_time
	avg_zkp_verification_time = avg_zkp_verification_time / len(zkp_ver_times)
	vp_ver_metrics.NumberOfTokensInVP = number_of_tokens
	vp_ver_metrics.VPVerificationTime = int(vp_ver_end.Sub(vp_ver_start).Microseconds())
	vp_ver_metrics.GrothProofVerificationTime = avg_zkp_verification_time
	verifier.Result.AddVPRelatedMetrics(vp_ver_metrics)

	return true
}

/*
VerifyTokens verify revocation status of tokens received by the verifier.
The tokens corresponding to the current epoch are verified.
*/
func (verifier *Verifier) VerifyTokens() {

	expiredVPs := make([]int, 0)
	_, Tokens := verifier.GetRevokedTokens()
	currentEpoch := verifier.GetCurrentEpoch()
	for vpIndex, tokenStorage := range verifier.Tokens {
		for i := 0; i < len(tokenStorage); i++ {
			var revocationStatus string

			if int(tokenStorage[i].Epoch) == currentEpoch {
				token := tokenStorage[i].Token
				if tokenStorage[i].NumberOfEpochsVCisValid >= currentEpoch {

					if slices.Contains(Tokens, hex.EncodeToString(token)) == true {
						revocationStatus = "revoked"
					} else {
						revocationStatus = "valid"
					}
					zap.S().Infoln("****Verifier*****: Verify Token: ", utils2.GetShortString(hex.EncodeToString(token)),
						"\t from VP: ", vpIndex,
						"\t epoch in VP: ", int(tokenStorage[i].Epoch),
						"\t current epoch: ", currentEpoch,
						"\t status: ", revocationStatus)
				} else {
					expiredVPs = append(expiredVPs, vpIndex)
					revocationStatus = "expired"
				}
			}
		}
	}
	if len(expiredVPs) > 0 {
		verifier.RemoveExpiredVPs(expiredVPs)
	}
}

func (verifier *Verifier) RemoveExpiredVPs(expiredVPs []int) {
	for i := 0; i < len(expiredVPs); i++ {
		delete(verifier.Tokens, expiredVPs[i])
	}
	zap.S().Info("***Verifier*****: Cleaned up tokens corresponding to expired VCs: ", expiredVPs)

}

func (verifier *Verifier) VerifySelectiveDisclosure(messages interface{}, elements model.SelectiveDisclosureElements) []byte {
	indexes := elements.Indexes
	claims := elements.SelectClaims.(model.SampleEmploymentProofPresentation)
	employeeDesignation := claims.EmployeeDesignation
	salary := claims.Salary

	var claims_hash []byte
	if individual_hashes, ok := messages.([][]byte); ok {

		g := bn254_mimc.NewMiMC()
		_, _ = g.Write([]byte(employeeDesignation))
		h1 := g.Sum(nil)

		if !bytes.Equal(h1, individual_hashes[indexes[0]]) {
			zap.S().Errorln("Wrong claims")
			return nil
		}

		g2 := bn254_mimc.NewMiMC()
		_, _ = g2.Write([]byte(strconv.Itoa(salary)))
		h2 := g2.Sum(nil)

		if !bytes.Equal(h2, individual_hashes[indexes[1]]) {
			zap.S().Errorln("Wrong claims")
			return nil
		}

		f := bn254_mimc.NewMiMC()
		for i := 0; i < len(individual_hashes); i++ {
			_, _ = f.Write(individual_hashes[i])
		}
		claims_hash = f.Sum(nil)

		zap.S().Infoln("****Verifier*****: Received a VP: claims- designation: ", employeeDesignation, "\t salary: ", salary)

	}
	return claims_hash
}
