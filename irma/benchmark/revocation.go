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

type ResultRevocation struct {
	TotalRevokedVCsPerEpoch  int `json:"total_revoked_vcs_per_epoch"`
	TimeToRevokeVCs          int `json:"time_to_revoke_vcs"`
	WitnessUpdateSize        int `json:"witness_update_size"`
	TotalValidVCs            int `json:"total_valid_vcs"`
	RevocationRate           int `json:"revocation_rate"`
	IssuerBandwidth          int `json:"issuer_bandwidth"`
	TotalHoldersWithValidVCs int `json:"total_holders_with_valid"`
	TotalEpochs              int `json:"total_epochs"`
	CurrentEpoch             int `json:"current_epoch"`
}

type ResultRevocationList struct {
	Results []ResultRevocation
	mutex   sync.Mutex
}

func (list *ResultRevocationList) Add(result ResultRevocation, witnessUpdateMessageWithoutRepetition bool) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.Results = append(list.Results, result)
	WriteResultRevocationToFile(result, false, witnessUpdateMessageWithoutRepetition)
}

func (list *ResultRevocationList) Get(index int) ResultRevocation {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.Results[index]
}

func ComputeAverageResultRevocation(results []ResultRevocation, witnessUpdateMessageWithoutRepetition bool) []*ResultRevocation {

	if len(results) == 0 {
		filename := fmt.Sprintf("irma/benchmark/results/result_revocation.json")
		if witnessUpdateMessageWithoutRepetition == true {
			filename = fmt.Sprintf("irma/benchmark/results/result_revocation_witness_update_without_repetition.json")
		}

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: TotalvalidVCs#TotalRevokedVCsPerEpoch#CurrentEpoch#TotalEpochs
	AvgResults := make(map[string]*ResultRevocation)

	//key: TotalVCs#TotalRevokedVCsPerEpoch
	isExists := make(map[string]bool)

	for _, res := range results {

		id := fmt.Sprintf(strconv.Itoa(res.TotalValidVCs), "#", strconv.Itoa(res.TotalRevokedVCsPerEpoch), "#", strconv.Itoa(res.CurrentEpoch), "#", strconv.Itoa(res.TotalEpochs))
		_, okay := isExists[id]
		if okay {
			AvgResults[id].TimeToRevokeVCs = (AvgResults[id].TimeToRevokeVCs + res.TimeToRevokeVCs) / 2
			AvgResults[id].WitnessUpdateSize = (AvgResults[id].WitnessUpdateSize + res.WitnessUpdateSize) / 2
			AvgResults[id].IssuerBandwidth = (AvgResults[id].IssuerBandwidth + res.IssuerBandwidth) / 2
		} else {
			result := &ResultRevocation{}
			result.CurrentEpoch = res.CurrentEpoch
			result.TotalEpochs = res.TotalEpochs
			result.TotalValidVCs = res.TotalValidVCs
			result.RevocationRate = res.RevocationRate
			result.TotalRevokedVCsPerEpoch = res.TotalRevokedVCsPerEpoch
			result.TimeToRevokeVCs = res.TimeToRevokeVCs
			result.WitnessUpdateSize = res.WitnessUpdateSize
			result.IssuerBandwidth = res.IssuerBandwidth
			result.TotalHoldersWithValidVCs = res.TotalHoldersWithValidVCs
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

func WriteResultRevocationToFile(result ResultRevocation, isavg bool, witnessUpdateMessageWithoutRepetition bool) {

	var results []ResultRevocation
	filename := fmt.Sprintf("irma/benchmark/results/result_revocation.json")
	if isavg == true && witnessUpdateMessageWithoutRepetition == false {
		filename = fmt.Sprintf("irma/benchmark/results/result_revocation_avg.json")
	}
	if isavg == false && witnessUpdateMessageWithoutRepetition == true {
		filename = fmt.Sprintf("irma/benchmark/results/result_revocation_witness_update_without_repetition.json")
	}
	if isavg == true && witnessUpdateMessageWithoutRepetition == true {
		filename = fmt.Sprintf("irma/benchmark/results/result_revocation_avg_witness_update_without_repetition.json")
	}

	jsonFile, err := os.Open(filename)
	if err != nil {
		jsonFile2, err2 := os.Create(filename)
		if err2 != nil {
			zap.S().Errorln("ERROR: " + filename + " creation error")
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
}

func ResetRevocationFiles(witnessUpdateMessageWithoutRepetition bool) {

	if witnessUpdateMessageWithoutRepetition == false {
		filename1 := fmt.Sprintf("irma/benchmark/results/result_revocation.json")
		filename2 := fmt.Sprintf("irma/benchmark/results/result_revocation_avg.json")
		os.Remove(filename1)
		os.Remove(filename2)

	}

	if witnessUpdateMessageWithoutRepetition == true {
		filename3 := fmt.Sprintf("irma/benchmark/results/result_revocation_avg_witness_update_without_repetition.json")
		filename4 := fmt.Sprintf("irma/benchmark/results/result_revocation_witness_update_without_repetition.json")
		os.Remove(filename3)
		os.Remove(filename4)
	}
}
