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

type ResultListCommitment struct {
	TotalRevokedVCsPerEpoch        int `json:"total_revoked_vcs_per_epoch"`
	TotalEpochs                    int `json:"total_epochs"`
	CurrentEpoch                   int `json:"current_epoch"`
	RevocationRate                 int `json:"revocation_rate"`
	TotalValidVCs                  int `json:"total_valid_vcs"`
	SizeOfTheListAtTheCurrentEpoch int `json:"size_of_the_list_at_the_current_epoch"`
	TimeToCreateCommitment         int `json:"time_to_create_commitment"`
	TimeToVerifyCommitment         int `json:"time_to_verify_commitment"`
}

type ResultListCommitmentList struct {
	Results []ResultListCommitment
	mutex   sync.Mutex
}

func (list *ResultListCommitmentList) Add(result ResultListCommitment) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.Results = append(list.Results, result)
	WriteResultListCommitmentToFile(result, false)
}

func (list *ResultListCommitmentList) Get(index int) ResultListCommitment {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.Results[index]
}

func ComputeAverageResultListCommitment(results []ResultListCommitment) []*ResultListCommitment {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_list_commitment.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: TotalvalidVCs#TotalRevokedVCsPerEpoch#CurrentEpoch#TotalEpochs
	AvgResults := make(map[string]*ResultListCommitment)
	isExists := make(map[string]bool)
	for _, res := range results {
		id := fmt.Sprintf(strconv.Itoa(res.TotalValidVCs), "#", strconv.Itoa(res.TotalRevokedVCsPerEpoch), "#", strconv.Itoa(res.CurrentEpoch), "#", strconv.Itoa(res.TotalEpochs))
		_, okay := isExists[id]
		if okay {
			AvgResults[id].TimeToCreateCommitment = (AvgResults[id].TimeToCreateCommitment + res.TimeToCreateCommitment) / 2
			AvgResults[id].TimeToVerifyCommitment = (AvgResults[id].TimeToVerifyCommitment + res.TimeToVerifyCommitment) / 2

		} else {
			result := &ResultListCommitment{}
			result.TotalRevokedVCsPerEpoch = res.TotalRevokedVCsPerEpoch
			result.TotalValidVCs = res.TotalValidVCs
			result.RevocationRate = res.RevocationRate
			result.TotalEpochs = res.TotalEpochs
			result.CurrentEpoch = res.CurrentEpoch

			result.TimeToCreateCommitment = res.TimeToCreateCommitment
			result.TimeToVerifyCommitment = res.TimeToVerifyCommitment
			result.SizeOfTheListAtTheCurrentEpoch = res.SizeOfTheListAtTheCurrentEpoch

			AvgResults[id] = result
			isExists[id] = true
		}
	}

	var finalResults []*ResultListCommitment
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func (r *ResultListCommitment) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultResultListCommitment(jsonObj []byte) *ResultListCommitment {
	res := ResultListCommitment{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func WriteResultListCommitmentToFile(result ResultListCommitment, isavg bool) {

	var results []ResultListCommitment
	filename := fmt.Sprintf("benchmark/results/result_list_commitment.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_list_commitment_avg.json")
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

func ResetCommitmentFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_list_commitment.json")
	filename2 := fmt.Sprintf("benchmark/results/result_list_commitment_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
