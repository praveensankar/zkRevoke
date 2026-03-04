package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type ResultIssuance struct {
	EDDSASignatureTime int `json:"eddsa_signature_time"`
	EDDSASignatureSize int `json:"eddsa_signature_size"`
	SeedSize           int `json:"seed_size"`
}

func ComputeAverageResultIssuance(results []ResultIssuance) ResultIssuance {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_issuance.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	firstResult := true

	var result ResultIssuance
	for _, res := range results {
		result.EDDSASignatureSize += res.EDDSASignatureSize
		result.EDDSASignatureTime += res.EDDSASignatureTime

		result.SeedSize += res.SeedSize
		if firstResult == false {
			result.EDDSASignatureTime = result.EDDSASignatureTime / 2
			result.EDDSASignatureSize = result.EDDSASignatureSize / 2
			result.SeedSize = result.SeedSize / 2
		}
		firstResult = false
	}
	return result
}

func WriteResultIssuanceToFile(result ResultIssuance, isavg bool) {

	var results []ResultIssuance
	filename := fmt.Sprintf("benchmark/results/result_issuance.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_issuance_avg.json")
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
	filename1 := fmt.Sprintf("benchmark/results/result_issuance.json")
	filename2 := fmt.Sprintf("benchmark/results/result_issuance_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
