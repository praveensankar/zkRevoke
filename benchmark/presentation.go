package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type ResultPresentation struct {
	VPValidityPeriod           int `json:"vp_validity_period"`
	NumberOfTokensInCircuit    int `json:"number_of_tokens_in_circuit"`
	NumberOfZKPProofs          int `json:"number_of_zkp_proofs"`
	ZKPProofSize               int `json:"zkp_proof_size"`
	TotalZKPProofSize          int `json:"total_zkp_proof_size"`
	TimeToGenerateOneZKPProof  int `json:"time_to_generate_one_zkp_proof"`
	TimeToGenerateAllZKPProofs int `json:"time_to_generate_all_zkp_proofs"`
}

func ComputeAverageResultPresentation(results []ResultPresentation) []*ResultPresentation {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_presentation.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: <m>#<k>
	AvgResults := make(map[string]*ResultPresentation)

	//key: <m>#<k>
	isExists := make(map[string]bool)
	for _, res := range results {
		m := res.VPValidityPeriod
		k := res.NumberOfTokensInCircuit

		id := fmt.Sprintf("%d", m, "#%d", k)
		_, okay := isExists[id]
		if okay {
			AvgResults[id].NumberOfZKPProofs = (AvgResults[id].NumberOfZKPProofs + res.NumberOfZKPProofs) / 2
			AvgResults[id].ZKPProofSize = (AvgResults[id].ZKPProofSize + res.ZKPProofSize) / 2
			AvgResults[id].TotalZKPProofSize = (AvgResults[id].TotalZKPProofSize + res.TotalZKPProofSize) / 2

			AvgResults[id].TimeToGenerateOneZKPProof = (AvgResults[id].TimeToGenerateOneZKPProof + res.TimeToGenerateOneZKPProof) / 2
			AvgResults[id].TimeToGenerateAllZKPProofs = (AvgResults[id].TimeToGenerateAllZKPProofs + res.TimeToGenerateAllZKPProofs) / 2

		} else {
			result := &ResultPresentation{}
			result.VPValidityPeriod = m
			result.NumberOfTokensInCircuit = k
			result.NumberOfZKPProofs = res.NumberOfZKPProofs
			result.ZKPProofSize = res.ZKPProofSize
			result.TotalZKPProofSize = res.TotalZKPProofSize
			result.TimeToGenerateOneZKPProof = res.TimeToGenerateOneZKPProof
			result.TimeToGenerateAllZKPProofs = res.TimeToGenerateAllZKPProofs

			AvgResults[id] = result
			isExists[id] = true
		}

	}

	var finalResults []*ResultPresentation
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func WriteResultPresentationToFile(result ResultPresentation, isavg bool) {

	var results []ResultPresentation
	filename := fmt.Sprintf("benchmark/results/result_presentation.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_presentation_avg.json")
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

func ResetPresentationFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_presentation.json")
	filename2 := fmt.Sprintf("benchmark/results/result_presentation_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
