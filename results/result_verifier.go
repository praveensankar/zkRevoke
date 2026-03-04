package results

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"zkrevoke/config"
)

type VPVerificationMetrics struct {
	NumberOfTokenBlocksInVP         int `json:"number_of_token_blocks"`
	NumberOfTokensInVP              int `json:"number_of_tokens_in_vp"`
	VPVerificationTime              int `json:"vp_verification_time"`
	GrothProofVerificationTime      int `json:"groth_proof_verification_time"`
	GrothProofVerificationTimeTotal int `json:"groth_proof_verification_time_total"`
}

type ResultVerifier struct {
	TotalVCs           int `json:"totalVCs"`
	NumberOfRevokedVCs int `json:"numberOfRevokedVCs"`
	TotalEpochs        int `json:"totalEpochs"`
	EpochDuration      int `json:"epochDuration"`

	VPMetrics []VPVerificationMetrics `json:"vp_metrics"`

	SizeOfRevocationList int `json:"sizeOfRevocationList"`
}

func NewResultVerifier(conf config.Config) *ResultVerifier {
	res := ResultVerifier{}
	res.SetTotalEpochs(int(conf.Params.ExpirationPeriod))
	res.SetTotalVCs(int(conf.Params.TotalVCs))
	//res.SetRevokedVCs(int(conf.ExpectedNumberofRevokedVCs))
	res.SetEpochDuration(int(conf.Params.EpochDuration))

	return &res
}

func (r *ResultVerifier) SetTotalVCs(totalVCs int) {
	r.TotalVCs = totalVCs
}

func (r *ResultVerifier) SetRevokedVCs(revokedVCs int) {
	r.NumberOfRevokedVCs = revokedVCs
}

func (r *ResultVerifier) SetTotalEpochs(totalEpochs int) {
	r.TotalEpochs = totalEpochs
}

func (r *ResultVerifier) SetEpochDuration(epochDuration int) {
	r.EpochDuration = epochDuration
}

func (r *ResultVerifier) AddVPRelatedMetrics(metrics VPVerificationMetrics) {
	metric := VPVerificationMetrics{}
	metric.NumberOfTokenBlocksInVP = metrics.NumberOfTokenBlocksInVP
	metric.NumberOfTokensInVP = metrics.NumberOfTokensInVP
	metric.GrothProofVerificationTime = metrics.GrothProofVerificationTime
	metric.VPVerificationTime = metrics.VPVerificationTime
	metric.GrothProofVerificationTimeTotal = metrics.GrothProofVerificationTimeTotal
	r.VPMetrics = append(r.VPMetrics, metric)
}

func (r *ResultVerifier) SetRevocationListSize(revocationListSize int) {
	r.SizeOfRevocationList = revocationListSize
}

func WriteResultVerifierToFile(result ResultVerifier) {

	var results []ResultVerifier
	//filename := fmt.Sprintf("results/results_computed.json")
	filename := fmt.Sprintf("results/result_%d.json", result.TotalVCs)
	jsonFile, err := os.Open(filename)
	if err != nil {
		jsonFile2, err2 := os.Create(filename)
		if err2 != nil {
			zap.S().Errorln("ERROR - results.json file creation error")
		}
		resJson, _ := ioutil.ReadAll(jsonFile2)
		json.Unmarshal(resJson, &results)
		results = append(results, result)
	} else {
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)
		results = append(results, result)
	}

	jsonRes, err3 := json.MarshalIndent(results, "", "")
	if err3 != nil {
		zap.S().Errorln("ERROR - marshalling the results")
	}

	err = ioutil.WriteFile(filename, jsonRes, 0644)
	if err != nil {
		zap.S().Errorln("unable to write results to file")
	}
	//zap.S().Errorln("RESULTS - successfully written to the file")
	zap.S().Info(result.String())
}

func (r *ResultVerifier) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultVerifier(jsonObj []byte) *ResultVerifier {
	res := ResultVerifier{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func (r ResultVerifier) String() string {
	var response string
	response = response + "Total VCs: " + fmt.Sprintf("%d", r.TotalVCs) + "\n"
	response = response + "Revoked VCs: " + fmt.Sprintf("%d", r.NumberOfRevokedVCs) + "\n"
	response = response + "Total Epochs: " + fmt.Sprintf("%d", r.TotalEpochs) + "\n"
	response = response + "Epoch Duration (in seconds): " + fmt.Sprintf("%d", r.EpochDuration) + "\n"

	return response
}
