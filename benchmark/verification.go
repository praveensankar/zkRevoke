package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type ResultVerification struct {
	VPValidityPeriod        int `json:"vp_validity_period"`
	ZKPProofVerTime         int `json:"zkp_proof_ver_time"`
	NumberOfTokensInCircuit int `json:"number_of_tokens_in_circuit"`
}

func ComputeAverageResultVerification(results []ResultVerification) []*ResultVerification {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_verification.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: <m>#<k>
	AvgResults := make(map[string]*ResultVerification)

	//key: <m>#<k>
	isExists := make(map[string]bool)
	for _, res := range results {
		m := res.VPValidityPeriod
		k := res.NumberOfTokensInCircuit
		id := fmt.Sprintf("%d#", m, "%d", k)
		_, okay := isExists[id]
		if okay {

			AvgResults[id].ZKPProofVerTime = (AvgResults[id].ZKPProofVerTime + res.ZKPProofVerTime) / 2

		} else {
			result := &ResultVerification{}
			result.VPValidityPeriod = m
			result.NumberOfTokensInCircuit = k
			result.ZKPProofVerTime = res.ZKPProofVerTime

			AvgResults[id] = result
			isExists[id] = true
		}

	}

	var finalResults []*ResultVerification
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func WriteResultVerificationToFile(result ResultVerification, isavg bool) {

	var results []ResultVerification
	filename := fmt.Sprintf("benchmark/results/result_verification.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_verification_avg.json")
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

func ResetVerificationFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_verification.json")
	filename2 := fmt.Sprintf("benchmark/results/result_verification_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
