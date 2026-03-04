package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

/*
RevocationTime - it is measured in nano seconds (Please note this since all other measurements in micro seconds).
*/
type ResultRevocation struct {
	TotalRevokedVCsPerEpoch int `json:"total_revoked_vcs_per_epoch"`
	TimeToRevokeVCs         int `json:"time_to_revoke_vcs"`
	TotalValidVCs           int `json:"total_valid_vcs"`
	RevocationRate          int `json:"revocation_rate"`
	IssuerBandwidth         int `json:"issuer_bandwidth"`
	TotalEpochs             int `json:"total_epochs"`
	CurrentEpoch            int `json:"current_epoch"`
	TokenSize               int `json:"tokenSize"`
}

type ResultRevocationList struct {
	Results []ResultRevocation
	mutex   sync.Mutex
}

func (list *ResultRevocationList) Add(result ResultRevocation) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.Results = append(list.Results, result)
	WriteResultRevocationToFile(result, false)
}

func (list *ResultRevocationList) Get(index int) ResultRevocation {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.Results[index]
}

func ComputeAverageResultRevocation(results []ResultRevocation) []*ResultRevocation {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_revocation.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: TotalvalidVCs#TotalRevokedVCsPerEpoch#CurrentEpoch#TotalEpochs
	AvgResults := make(map[string]*ResultRevocation)
	isExists := make(map[string]bool)
	for _, res := range results {
		id := fmt.Sprintf(strconv.Itoa(res.TotalValidVCs), "#", strconv.Itoa(res.TotalRevokedVCsPerEpoch), "#", strconv.Itoa(res.CurrentEpoch), "#", strconv.Itoa(res.TotalEpochs))
		_, okay := isExists[id]
		if okay {
			AvgResults[id].TimeToRevokeVCs = (AvgResults[id].TimeToRevokeVCs + res.TimeToRevokeVCs) / 2
			AvgResults[id].IssuerBandwidth = (AvgResults[id].IssuerBandwidth + res.IssuerBandwidth) / 2
			AvgResults[id].TokenSize = (AvgResults[id].TokenSize + res.TokenSize) / 2
		} else {
			result := &ResultRevocation{}
			result.TotalRevokedVCsPerEpoch = res.TotalRevokedVCsPerEpoch
			result.TimeToRevokeVCs = res.TimeToRevokeVCs
			result.TotalValidVCs = res.TotalValidVCs
			result.RevocationRate = res.RevocationRate
			result.IssuerBandwidth = res.IssuerBandwidth
			result.TotalEpochs = res.TotalEpochs
			result.CurrentEpoch = res.CurrentEpoch
			result.TokenSize = res.TokenSize
			AvgResults[id] = result
			isExists[id] = true
		}
	}

	var finalResults []*ResultRevocation
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func (r *ResultRevocation) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultRevocation(jsonObj []byte) *ResultRevocation {
	res := ResultRevocation{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func WriteResultRevocationToFile(result ResultRevocation, isavg bool) {

	var results []ResultRevocation
	filename := fmt.Sprintf("benchmark/results/result_revocation.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_revocation_avg.json")
	}
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
}

func ResetRevocationFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_revocation.json")
	filename2 := fmt.Sprintf("benchmark/results/result_revocation_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
