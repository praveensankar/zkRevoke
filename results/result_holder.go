package results

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
	"zkrevoke/config"
)

type VPRelatedMetrics struct {
	NumberOfTokenBlocks    int `json:"number_of_token_blocks"`
	NumberOfTokensInaBlock int `json:"number_of_tokens_ina_block"`
	SizeOfTokenBlock       int `json:"size_of_token_block"`
	VPSize                 int `json:"vp_size"`
	VPGenerationTime       int `json:"vPGenerationTime"`
	AllZKPProofsGenTime    int `json:"all_zkp_proofs_gen_time"`
	SingleZKPProofGenTime  int `json:"single_zkp_proofs_gen_time"`
}
type ResultHolder struct {
	EpochDuration      int `json:"epochDuration"`
	NumberOfTokensInVP int `json:"numberOfTokensInVP"`

	GrothProofGenerationTime    []int `json:"grothproofGenerationTime"`
	GrothProofGenerationTimeAvg int   `json:"grothproofGenerationTimeAvg"`
	GrothProofSize              []int `json:"grothproofSize"`
	GrothProofSizeAvg           int   `json:"grothproofSizeAvg"`

	VPSizeMetrics          []VPRelatedMetrics `json:"vPSizeMetrics"`
	TokenGenerationTime    []int              `json:"tokenGenerationTime"`
	TokenGenerationTimeAvg int                `json:"tokenGenerationTimeAvg"`
}

func NewResultHolder(conf config.Config, numberOfTokensInVP int) *ResultHolder {
	res := ResultHolder{}
	res.SetEpochDuration(int(conf.Params.EpochDuration))
	res.SetNumberOfTokensInVP(numberOfTokensInVP)
	return &res
}

func (r *ResultHolder) SetEpochDuration(epochDuration int) {
	r.EpochDuration = epochDuration
}

func (r *ResultHolder) SetNumberOfTokensInVP(numberOfTokensInaVP int) {
	r.NumberOfTokensInVP = numberOfTokensInaVP
}

func (r *ResultHolder) AddTokenGenerationTime(time time.Duration) {
	r.TokenGenerationTime = append(r.TokenGenerationTime, int(time.Microseconds()))
}

func (r *ResultHolder) ComputeAvgTokenGenerationTime() {
	res := 0
	for i := 0; i < len(r.TokenGenerationTime); i++ {
		res += r.TokenGenerationTime[i]
	}
	res = res / len(r.TokenGenerationTime)
	r.TokenGenerationTimeAvg = res
}

func (r *ResultHolder) AddVPRelatedMetrics(metrics VPRelatedMetrics) {
	metric := VPRelatedMetrics{}
	metric.NumberOfTokensInaBlock = metrics.NumberOfTokensInaBlock
	metric.SizeOfTokenBlock = metrics.SizeOfTokenBlock
	metric.NumberOfTokenBlocks = metrics.NumberOfTokenBlocks
	metric.VPSize = metrics.VPSize
	metric.VPGenerationTime = metrics.VPGenerationTime
	metric.AllZKPProofsGenTime = metrics.AllZKPProofsGenTime
	metric.SingleZKPProofGenTime = metrics.SingleZKPProofGenTime
	r.VPSizeMetrics = append(r.VPSizeMetrics, metric)
}

func (r *ResultHolder) AddGrothProofSize(proofSize int) {
	r.GrothProofSize = append(r.GrothProofSize, proofSize)
}

func (r *ResultHolder) ComputeAvgGrothProofSize() {
	res := 0
	for i := 0; i < len(r.GrothProofSize); i++ {
		res += r.GrothProofSize[i]
	}
	res = res / len(r.GrothProofSize)
	r.GrothProofSizeAvg = res
}

func (r *ResultHolder) AddGrothProofGenerationTime(grothProofGenerationTime time.Duration) {
	r.GrothProofGenerationTime = append(r.GrothProofGenerationTime, int(grothProofGenerationTime.Microseconds()))
}

func (r *ResultHolder) ComputeAvgGrothProofGenerationTime() {
	res := 0
	for i := 0; i < len(r.GrothProofGenerationTime); i++ {
		res += r.GrothProofGenerationTime[i]
	}
	res = res / len(r.GrothProofGenerationTime)
	r.GrothProofGenerationTimeAvg = res
}

func WriteResultHolderToFile(result ResultHolder) {

	var results []ResultHolder
	//filename := fmt.Sprintf("results/results_computed.json")
	filename := fmt.Sprintf("results/result_holder.json")
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

func (r *ResultHolder) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultHolder(jsonObj []byte) *ResultHolder {
	res := ResultHolder{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func (r ResultHolder) String() string {
	var response string

	response = response + "Epoch Duration (in seconds): " + fmt.Sprintf("%d", r.EpochDuration) + "\n"
	response = response + "Number of tokens in a VP (in micro seconds): " + fmt.Sprintf("%d", r.NumberOfTokensInVP) + "\n"

	response = response + "Time to generate tokens (in micro seconds):" + fmt.Sprintf("%v", r.TokenGenerationTime) + "\n"
	response = response + "Average time to generate tokens (in micro seconds):" + fmt.Sprintf("%v", r.TokenGenerationTimeAvg) + "\n"

	response = response + "Groth16 proof size (in bytes): " + fmt.Sprintf("%v", r.GrothProofSize) + "\n"
	response = response + "Average time to generate groth16 proof (in micro seconds):" + fmt.Sprintf("%v", r.GrothProofGenerationTimeAvg) + "\n"

	return response
}
