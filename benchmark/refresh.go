package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"os"
)

type ResultRefresh struct {
	NumberOfRevokedVCs   int      `json:"numberOfRevokedVCs"`
	RevocationRate       int      `json:"revocation_rate"`
	TokenSize            int      `json:"tokenSize"`
	Gas                  uint64   `json:"gas"`
	Cost                 *big.Int `json:"cost"`
	NumberOfTransactions int      `json:"numberOfTransactions"`
	TotalEpochs          int      `json:"total_epochs"`
	CurrentEpoch         int      `json:"current_epoch"`
	Time                 int      `json:"time"`
}

func (res *ResultRefresh) SetNumberOfRevokedVCs(numRevokedVCs int) {
	res.NumberOfRevokedVCs = numRevokedVCs
}

func (res *ResultRefresh) SetTokenSize(tokenSize int) {
	res.TokenSize = tokenSize
}

func (res *ResultRefresh) SetCost(cost *big.Int) {
	res.Cost = cost
}

func (res *ResultRefresh) SetGas(gas uint64) {
	res.Gas = gas
}

func (res *ResultRefresh) SetNumberOfTransactions(numTransactions int) {
	res.NumberOfTransactions = numTransactions
}

func (res *ResultRefresh) SetTimeToRefreshTheList(t int) {
	res.Time = t
}

func ComputeAverageResultRefresh(results []ResultRefresh) []*ResultRefresh {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_refresh.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	// key: NumberOfRevokedVCs
	AvgResults := make(map[int]*ResultRefresh)

	for _, res := range results {
		_, exists := AvgResults[res.NumberOfRevokedVCs]
		if exists {
			AvgResults[res.NumberOfRevokedVCs].TokenSize = (AvgResults[res.NumberOfRevokedVCs].TokenSize + res.TokenSize) / 2
			AvgResults[res.NumberOfRevokedVCs].Gas = (AvgResults[res.NumberOfRevokedVCs].Gas + res.Gas) / 2
			AvgResults[res.NumberOfRevokedVCs].Cost.Add(AvgResults[res.NumberOfRevokedVCs].Cost, res.Cost)
			AvgResults[res.NumberOfRevokedVCs].Cost.Div(AvgResults[res.NumberOfRevokedVCs].Cost, big.NewInt(2))
			AvgResults[res.NumberOfRevokedVCs].NumberOfTransactions = (AvgResults[res.NumberOfRevokedVCs].NumberOfTransactions + res.NumberOfTransactions) / 2
			AvgResults[res.NumberOfRevokedVCs].Time = (AvgResults[res.NumberOfRevokedVCs].Time + res.Time) / 2

		} else {
			result := &ResultRefresh{}
			result.NumberOfRevokedVCs = res.NumberOfRevokedVCs
			result.TokenSize = res.TokenSize
			result.Gas = res.Gas
			result.Cost = res.Cost
			result.NumberOfTransactions = res.NumberOfTransactions
			result.Time = res.Time
			result.RevocationRate = res.RevocationRate
			result.TotalEpochs = res.TotalEpochs
			result.CurrentEpoch = res.CurrentEpoch
			AvgResults[res.NumberOfRevokedVCs] = result
		}
	}

	var finalResults []*ResultRefresh
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func (r *ResultRefresh) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultRefresh(jsonObj []byte) *ResultRefresh {
	res := ResultRefresh{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func WriteResultRefreshToFile(result ResultRefresh, isavg bool) {

	var results []ResultRefresh
	filename := fmt.Sprintf("benchmark/results/result_refresh.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_refresh_avg.json")
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

func (r ResultRefresh) String() string {
	var response string
	response = response + "Revoked VCs: " + fmt.Sprintf("%d", r.NumberOfRevokedVCs) + "\t"
	response = response + "Token size (in bytes): " + fmt.Sprintf("%d", r.TokenSize) + "\t"
	response = response + "Cost to store token (in wei):" + fmt.Sprintf("%s", r.Cost) + "\t"
	response = response + "Gas used:" + fmt.Sprintf("%d", r.Gas) + "\t"
	response = response + "Number of transactions (each transaction requires a separate block):" + fmt.Sprintf("%d", r.NumberOfTransactions) + "\n"
	return response
}

func ResetRefreshFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_refresh.json")
	filename2 := fmt.Sprintf("benchmark/results/result_refresh_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}

func WriteTokenStorageResultToFile(result ResultRefresh, isavg bool) {

	var results []ResultRefresh
	filename := fmt.Sprintf("benchmark/results/result_token_storage.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_token_storage_avg.json")
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

func ResetTokenStorageResultFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_token_storage.json")
	filename2 := fmt.Sprintf("benchmark/results/result_token_storage_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
