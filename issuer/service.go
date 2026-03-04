package issuer

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/deckarep/golang-set"

	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"zkrevoke/config"
	"zkrevoke/crypto2"
	"zkrevoke/model"
	"zkrevoke/utils"
)

func (issuer *Issuer) RequestVC() (*model.VerifiableCredential, *utils.PublicParams, int, time.Time, int) {
	issuer.vcCounter = issuer.vcCounter + 1
	vcID := fmt.Sprintf("vcID#%s", strconv.Itoa(issuer.vcCounter))
	vc := issuer.verfiableCredentials[vcID]
	pki := &utils.PublicParams{}
	pki.Ccs = issuer.ccs
	pki.EddsaPublicKey = issuer.eddsaPublicKey
	pki.ZkpProvingKey = issuer.zkpProvingKey
	pki.ZkpVerifyingKey = issuer.zkpVerifyingKey
	return &vc, pki, issuer.Duration, issuer.InitialTimeStamp, issuer.NumberOfTokensInCircuit
}
func (issuer *Issuer) GenerateDummyVC() *model.VerifiableCredential {
	seed := rand.Text()
	issuer.Lock()
	issuer.vcCounter = issuer.vcCounter + 1
	issuer.Unlock()
	vcID := fmt.Sprintf("vcID#%s", strconv.Itoa(issuer.vcCounter))
	newVC := &model.VerifiableCredential{}
	newVC.Id = vcID
	newVC.Seed = seed
	issuer.Lock()
	issuer.verfiableCredentials[vcID] = *newVC
	issuer.seedStore[vcID] = seed
	issuer.Unlock()
	return newVC
}

func (issuer *Issuer) GenerateVC(Holder_PublicKey []byte) (*model.VerifiableCredential, *utils.PublicParams, int, time.Time, int) {
	seed := rand.Text()
	issuer.Lock()
	issuer.vcCounter = issuer.vcCounter + 1
	issuer.Unlock()
	vcID := fmt.Sprintf("vcID#%s", strconv.Itoa(issuer.vcCounter))
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(issuer.TotalNumberOfEpocs)))
	validFrom := time.Now()
	vadlidFromStr := strconv.Itoa(int(validFrom.Unix()))
	validUntilDuration := int(random.Int64()) % issuer.TotalNumberOfEpocs
	if validUntilDuration < issuer.MinimumExpiryDurationOfVC {
		validUntilDuration = issuer.MinimumExpiryDurationOfVC
	}
	validUntil := validFrom.Add(time.Duration(validUntilDuration) * 86400 * time.Second)
	//zap.S().Infoln("Number of blocks vc is valid: ", utils.GetNumberOfBlocksVCisValid(int(validFrom.Unix()), int(validUntil.Unix())))
	validUntilStr := strconv.Itoa(int(validUntil.Unix()))
	start := time.Now()
	newVC, eddsaSigningTime, _ := model.CreateEmploymentProofVC(vcID, seed, Holder_PublicKey, issuer.eddsaPrivateKey, vadlidFromStr, validUntilStr, issuer.SelectiveDisclosureExtension)
	end := time.Now()
	issuer.Lock()
	issuer.verfiableCredentials[vcID] = *newVC
	issuer.seedStore[vcID] = seed
	issuer.Unlock()

	pki := &utils.PublicParams{}
	pki.Ccs = issuer.ccs
	pki.EddsaPublicKey = issuer.eddsaPublicKey
	pki.ZkpProvingKey = issuer.zkpProvingKey
	pki.ZkpVerifyingKey = issuer.zkpVerifyingKey

	issuer.Result.AddVCGenerationTime(end.Sub(start))
	issuer.Result.SetEDDSASignTime(eddsaSigningTime)

	vcJson := newVC.Json()
	issuer.Result.AddVCSize(len(vcJson))
	zap.S().Infoln("*****ISSUER****: issued new VC: ", vcID)
	return newVC, pki, issuer.Duration, issuer.InitialTimeStamp, issuer.NumberOfTokensInCircuit
}

func (issuer *Issuer) BulkIssueVCs(numberOfVCs int) ([]*model.VerifiableCredential, *utils.PublicParams, int, time.Time) {
	var vcs []*model.VerifiableCredential

	for i := 0; i < numberOfVCs; i++ {
		newVC, _, _, _, _ := issuer.GenerateVC(nil)
		vcs = append(vcs, newVC)

	}

	pki := &utils.PublicParams{}
	pki.Ccs = issuer.ccs
	pki.EddsaPublicKey = issuer.eddsaPublicKey
	pki.ZkpProvingKey = issuer.zkpProvingKey
	pki.ZkpVerifyingKey = issuer.zkpVerifyingKey

	return vcs, pki, issuer.Duration, issuer.InitialTimeStamp
}

/*
Returns
1) (int) revocation time: Time to revoke VCs
*/
func (issuer *Issuer) RevokeVCsRandomly(numberOfRevokedVCs int) int {

	revokedVCIds := mapset.NewSet()
	var vcIDs []string
	for _, vc := range issuer.verfiableCredentials {
		vcIDs = append(vcIDs, vc.Id)
	}

	revocationTime := 0
	totalVCs := len(issuer.verfiableCredentials)
	for i := 0; i < numberOfRevokedVCs; {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(totalVCs)))
		vcID := vcIDs[index.Int64()]
		if revokedVCIds.Contains(vcID) {
			continue
		} else {
			revokedVCIds.Add(vcID)
			start := time.Now()
			issuer.RevokeVC(vcID)
			end := time.Since(start)
			revocationTime = revocationTime + int(end.Nanoseconds())
			i++
		}
	}
	if numberOfRevokedVCs < 5 {
		zap.S().Infoln("****ISSUER****: revoked vc IDs: ", revokedVCIds.String())
	}
	return revocationTime

}

func (issuer *Issuer) RevokeVC(vcID string) {
	issuer.Lock()
	issuer.RevokedVCIDs = append(issuer.RevokedVCIDs, vcID)
	issuer.Unlock()
}

func (issuer *Issuer) GetIssuedVCIDs() []string {
	var vcIDs []string
	for _, vc := range issuer.verfiableCredentials {
		vcIDs = append(vcIDs, vc.Id)
	}
	return vcIDs
}

func (issuer *Issuer) GetRevokedTokens() [][]byte {
	var tokens [][]byte
	epoch := issuer.GetCurrentEpoch()
	for i := 0; i < len(issuer.RevokedVCIDs); i++ {
		seed := issuer.seedStore[issuer.RevokedVCIDs[i]]
		token_gen_start := time.Now()
		token := utils.ComputeToken(epoch, seed)
		token_gen_end := time.Now()
		issuer.Result.AddTokenGenerationTime(token_gen_end.Sub(token_gen_start))
		tokens = append(tokens, token)
	}
	return tokens
}

func (issuer *Issuer) CalculateTimeToComputeTokensGivenEpoch(epoch int) time.Duration {
	var tokens [][]byte
	token_gen_start := time.Now()
	for i := 0; i < len(issuer.RevokedVCIDs); i++ {
		seed := issuer.seedStore[issuer.RevokedVCIDs[i]]
		token := utils.ComputeToken(epoch, seed)
		tokens = append(tokens, token)
	}
	token_gen_end := time.Since(token_gen_start)
	return token_gen_end
}

func (issuer Issuer) GetCurrentEpoch() int {
	currentTime := time.Now()
	epoch := (int(currentTime.Sub(issuer.InitialTimeStamp).Seconds())) / issuer.Duration
	return epoch
}

func (issuer Issuer) GetTokenSize() int {
	seed := rand.Text()
	epoch := issuer.GetCurrentEpoch()
	token := utils.ComputeToken(epoch, seed)
	return len(token)
}

func (issuer *Issuer) PublishRevokedTokens() (uint64, string, int, int, int) {
	issuer.RemoveExpiredVCs()
	if len(issuer.RevokedVCIDs) > 0 {
		return issuer.PublishRevokedTokensToBlockchain()
	} else {
		return 0, "", 0, 0, 0
	}

}

func (issuer *Issuer) RemoveExpiredVCs() {
	var expiredVCIDs []string
	var UpdatedRevokedVCIDs []string
	for i := 0; i < len(issuer.RevokedVCIDs); i++ {
		vc := issuer.verfiableCredentials[issuer.RevokedVCIDs[i]]
		valid_from, _ := strconv.Atoi(vc.Metadata.ValidFrom)
		valid_until, _ := strconv.Atoi(vc.Metadata.ValidUntil)
		//numberOfEpochs := (valid_until - valid_from) / issuer.Duration
		numberOfEpochs := utils.GetNumberOfBlocksVCisValid(valid_from, valid_until)
		if numberOfEpochs < issuer.GetCurrentEpoch() {
			expiredVCIDs = append(expiredVCIDs, issuer.RevokedVCIDs[i])
		} else {
			UpdatedRevokedVCIDs = append(UpdatedRevokedVCIDs, issuer.RevokedVCIDs[i])
		}
	}
	if len(expiredVCIDs) > 0 {
		issuer.RevokedVCIDs = UpdatedRevokedVCIDs
		zap.S().Info("***ISSUER*****: Cleaned up expired VCs: ", expiredVCIDs)
	}
}

func (issuer *Issuer) GenerateAndStoreVCs(conf config.Config) {
	var vcs []*model.VerifiableCredential

	for i := 0; i < int(conf.Params.TotalVCs); i++ {
		newVC, _, _, _, _ := issuer.GenerateVC(nil)
		vcs = append(vcs, newVC)
	}
	var signature []byte
	for i := 0; i < len(vcs[0].Proofs); i++ {
		if vcs[0].Proofs[i].Type == string(model.ProofTypeEDDSA) {
			signature = vcs[0].Proofs[i].ProofValue
		}
	}
	zap.S().Infoln("****ISSUER*****: Generating VC: Id:", vcs[0].Id,
		" seed: ", vcs[0].Seed,
		"valid until: ", vcs[0].Metadata.ValidUntil,
		"signature: ", hex.EncodeToString(signature))
	zap.S().Infoln("***ISSUER***: Public key: ", hex.EncodeToString(issuer.eddsaPublicKey.Bytes()))
	filename := fmt.Sprintf("issuer/vc_data%d.json", int(conf.Params.TotalVCs))
	_, err := os.Open(filename)
	if err != nil {
		_, err2 := os.Create(filename)
		if err2 != nil {
			zap.S().Errorln("ERROR - vc_data.json file creation error")
		}
	}

	jsonRes, err3 := json.MarshalIndent(vcs, "", "")
	if err3 != nil {
		zap.S().Errorln("ERROR - marshalling the results")
	}

	err = ioutil.WriteFile(filename, jsonRes, 0644)
	if err != nil {
		zap.S().Errorln("unable to write VCs to file")
	}
}

/*
LoadStoredVCs loads pre-generated VCs to issuer's storage
*/
func (issuer *Issuer) LoadStoredVCs(conf config.Config) {
	start := time.Now()
	var VCs []*model.VerifiableCredential
	filename := fmt.Sprintf("issuer/vc_data%d.json", int(conf.Params.TotalVCs))
	jsonFile, err := os.Open(filename)
	if err != nil {
		zap.S().Errorln("Error opening vcs file:", err)
	} else {
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &VCs)
	}

	var newVCs []*model.VerifiableCredential
	for _, vc := range VCs {
		jsonObj, _ := json.MarshalIndent(vc, "", "  ")
		issuer.Result.AddVCSize(len(jsonObj))
		vc1 := model.JsonToEmploymentProofVC(jsonObj)
		newVCs = append(newVCs, vc1)

		f := bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(vc1.Id))
		_, _ = f.Write([]byte(vc1.Seed))

		msg := f.Sum(nil)
		var signature []byte
		for i := 0; i < len(vc1.Proofs); i++ {
			if vc1.Proofs[i].Type == string(model.ProofTypeEDDSA) {
				signature = vc1.Proofs[i].ProofValue
			}
		}
		if conf.DEBUG == true {
			//zap.S().Infoln("****ISSUER*****: Stored VC: Id:", vc1.Id,
			//	" seed: ", vc1.Seed,
			//	"valid until: ", vc1.Metadata.ValidUntil,
			//	"signature: ", hex.EncodeToString(signature))
			//zap.S().Infoln("***ISSUER***: Public key: ", hex.EncodeToString(issuer.eddsaPublicKey.Bytes()))
			status := crypto2.Verify_EDDSA(issuer.eddsaPublicKey, msg, signature)
			if status == false {
				zap.S().Fatalln("***ISSUER***Loading pre-generated VCs: Unable to verify signature")
			}
		}

	}
	//zap.S().Infoln("***ISSUER*****: Stored VCs: ", newVCs[0])
	//timestampofFirstvc, _ := strconv.Atoi(VCs[0].Metadata.ValidFrom)
	//issuer.InitialTimeStamp = time.Unix(int64(timestampofFirstvc), 0)
	for _, vc := range newVCs {
		issuer.seedStore[vc.Id] = vc.Seed
		issuer.verfiableCredentials[vc.Id] = *vc
	}
	end := time.Now()
	zap.S().Infoln("***ISSUER*****: Loaded pre-generated VCs: ", len(newVCs), "\t took: ", end.Sub(start).Minutes(), "minutes")
}

func (issuer *Issuer) GenerateAndStoreVCsWithoutSignatures(conf config.Config) {
	var vcs []*model.VerifiableCredential

	var wg sync.WaitGroup

	for i := 0; i < int(conf.Params.TotalVCs); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newVC := issuer.GenerateDummyVC()
			vcs = append(vcs, newVC)

			zap.S().Infoln("****ISSUER*****: Generated new VC: Id:", newVC.Id,
				" seed: ", newVC.Seed)

		}()
	}
	wg.Wait()
	if conf.DEBUG {
		zap.S().Infoln("***ISSUER***: Public key: ", hex.EncodeToString(issuer.eddsaPublicKey.Bytes()))
	}
	filename := fmt.Sprintf("issuer/vc_data%d.json", int(conf.Params.TotalVCs))
	_, err := os.Open(filename)
	if err != nil {
		_, err2 := os.Create(filename)
		if err2 != nil {
			zap.S().Errorln("ERROR - vc_data.json file creation error")
		}
	}

	jsonRes, err3 := json.MarshalIndent(vcs, "", "")
	if err3 != nil {
		zap.S().Errorln("ERROR - marshalling the results")
	}

	err = ioutil.WriteFile(filename, jsonRes, 0644)
	if err != nil {
		zap.S().Errorln("unable to write VCs to file")
	}
	zap.S().Infoln("******ISSUER*****: Generated ", conf.Params.TotalVCs, "VCs")
}

func (issuer *Issuer) LoadStoredVCsWithoutSignatures(conf config.Config) {
	start := time.Now()
	var VCs []*model.VerifiableCredential
	filename := fmt.Sprintf("issuer/vc_data%d.json", conf.Params.TotalVCs)
	jsonFile, err := os.Open(filename)
	if err != nil {
		zap.S().Errorln("Error opening vcs file:", err)
	} else {
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &VCs)
	}

	var newVCs []*model.VerifiableCredential
	for _, vc := range VCs {
		jsonObj, _ := json.MarshalIndent(vc, "", "  ")
		issuer.Result.AddVCSize(len(jsonObj))
		vc1 := model.JsonToEmploymentProofVC(jsonObj)
		newVCs = append(newVCs, vc1)
	}
	//zap.S().Infoln("***ISSUER*****: Stored VCs: ", newVCs[0])
	//timestampofFirstvc, _ := strconv.Atoi(VCs[0].Metadata.ValidFrom)
	//issuer.InitialTimeStamp = time.Unix(int64(timestampofFirstvc), 0)
	for _, vc := range newVCs {
		issuer.seedStore[vc.Id] = vc.Seed
		issuer.verfiableCredentials[vc.Id] = *vc
	}
	end := time.Now()
	zap.S().Infoln("***ISSUER*****: Loaded pre-generated VCs: ", len(newVCs), "\t took: ", end.Sub(start).Minutes(), "minutes")
}
