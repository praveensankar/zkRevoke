package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
)

type ResultIssuance struct {
	TotalVCs                 int `json:"total_vcs"`
	TimeToGenerateOneVC      int `json:"time_to_generate_one_vc"`
	TimeToGenerateAllVCs     int `json:"time_to_generate_all_vcs"`
	TimeToGenerateOneWitness int `json:"time_to_generate_one_witness"`
	TimeToGenerateAllWitness int `json:"time_to_generate_all_witness"`
}

func ComputeAverageResultIssuance(results []ResultIssuance) []*ResultIssuance {

	if len(results) == 0 {
		filename := fmt.Sprintf("irma/results/result_issuance.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: TotalVCs
	AvgResults := make(map[string]*ResultIssuance)

	//key: TotalVCs
	isExists := make(map[string]bool)

	for _, res := range results {

		id := fmt.Sprintf(strconv.Itoa(res.TotalVCs))
		_, okay := isExists[id]
		if okay {
			AvgResults[id].TimeToGenerateOneVC = (AvgResults[id].TimeToGenerateOneVC + res.TimeToGenerateOneVC) / 2
			AvgResults[id].TimeToGenerateAllVCs = (AvgResults[id].TimeToGenerateAllVCs + res.TimeToGenerateAllVCs) / 2
			AvgResults[id].TimeToGenerateOneWitness = (AvgResults[id].TimeToGenerateOneWitness + res.TimeToGenerateOneWitness) / 2
			AvgResults[id].TimeToGenerateAllWitness = (AvgResults[id].TimeToGenerateAllWitness + res.TimeToGenerateAllWitness) / 2
		} else {
			result := &ResultIssuance{}
			result.TotalVCs = res.TotalVCs
			result.TimeToGenerateOneVC = res.TimeToGenerateOneVC
			result.TimeToGenerateAllVCs = res.TimeToGenerateAllVCs
			result.TimeToGenerateOneWitness = res.TimeToGenerateOneWitness
			result.TimeToGenerateAllWitness = res.TimeToGenerateAllWitness
			AvgResults[id] = result
			isExists[id] = true
		}

	}

	var finalResults []*ResultIssuance
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func WriteResultIssuanceToFile(result ResultIssuance, isavg bool) {

	var results []ResultIssuance
	filename := fmt.Sprintf("irma/benchmark/results/result_issuance.json")
	if isavg {
		filename = fmt.Sprintf("irma/benchmark/results/result_issuance_avg.json")
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

func ResetIssuanceFiles() {
	filename1 := fmt.Sprintf("irma/benchmark/results/result_issuance.json")
	filename2 := fmt.Sprintf("irma/benchmark/results/result_issuance_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
