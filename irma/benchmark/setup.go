package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type ResultSetup struct {
	PrivateKeySize     int `json:"private_key_size"`
	PublicKeySize      int `json:"public_key_size"`
	KeyGenTime         int `json:"key_gen_time"`
	AccumulatorGenTime int `json:"accumulator_gen_time"`
	AccumulatorSize    int `json:"accumulator_size"`
}

func ComputeAverageResultSetup(results []ResultSetup) ResultSetup {

	if len(results) == 0 {
		filename := fmt.Sprintf("irma/results/result_setup.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	firstResult := true

	var result ResultSetup
	for _, res := range results {
		result.PrivateKeySize += res.PrivateKeySize
		result.PublicKeySize += res.PublicKeySize
		result.KeyGenTime += res.KeyGenTime
		result.AccumulatorGenTime += res.AccumulatorGenTime
		result.AccumulatorSize += res.AccumulatorSize

		if firstResult == false {
			result.PrivateKeySize = result.PrivateKeySize / 2
			result.PublicKeySize = result.PublicKeySize / 2
			result.KeyGenTime = result.KeyGenTime / 2
			result.AccumulatorGenTime = result.AccumulatorGenTime / 2
			result.AccumulatorSize = result.AccumulatorSize / 2
		}
		firstResult = false
	}
	return result
}

func WriteResultSetupToFile(result ResultSetup, isavg bool) {

	var results []ResultSetup
	filename := fmt.Sprintf("irma/benchmark/results/result_setup.json")
	if isavg {
		filename = fmt.Sprintf("irma/benchmark/results/result_setup_avg.json")
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

func ResetSetupFiles() {
	filename1 := fmt.Sprintf("irma/benchmark/results/result_setup.json")
	filename2 := fmt.Sprintf("irma/benchmark/results/result_setup_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
