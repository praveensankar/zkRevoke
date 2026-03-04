package benchmark

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"os"
	"zkrevoke/issuer"
)

type ResultCircuit struct {
	NumberOfTokensInCircuit     int      `json:"numberOfTokensInCircuit"`
	GrothCCSSize                int      `json:"grothCcsSize"`
	GrothCCSCost                *big.Int `json:"GrothccsCost"`
	GrothCCSGas                 uint64   `json:"grothCcsGas"`
	GrothCCSNumberOfConstraints int64    `json:"grothccsNumberOfConstraints"`
	GrothCCSTime                int      `json:"grothccsTime"`
}
type ResultSetup struct {
	ZKPCircuitResults        []ResultCircuit `json:"zkpcircuitResult"`
	GrothProvingKeySize      int             `json:"grothProvingKeySize"`
	GrothVerificationKeySize int             `json:"grothVerificationKeySize"`
	GrothVerificationKeyCost *big.Int        `json:"grothVerificationKeyCost"`
	GrothVerificationKeyGas  uint64          `json:"grothVerificationKeyGas"`

	EDDSAPrivateKeySize int      `json:"eddsaPrivateKeySize"`
	EDDSAPublicKeySize  int      `json:"eddsaPublicKeySize"`
	EDDSAPublicKeyCost  *big.Int `json:"eddsaPublicKeyCost"`
	EDDSAPublicKeyGas   uint64   `json:"eddsaPublicKeyGas"`
	EDDSAKeyGenTime     int      `json:"eddsaKeyGenTime"`

	ContractDeploymentCost *big.Int `json:"contractDeploymentCost"`
	ContractDeploymentGas  uint64   `json:"contractDeploymentGas"`
	AvgGasPriceMarch2025   int      `json:"avgGasPriceMarch2025"`
}

func (r *ResultCircuit) SetNumberOfTokensInCircuit(numberOfTokensInCircuit int) {
	r.NumberOfTokensInCircuit = numberOfTokensInCircuit
}

func (r *ResultCircuit) SetGrothCCSSize(grothCCSSize int) {
	r.GrothCCSSize = grothCCSSize
}

func (r *ResultCircuit) SetGrothCCSGas(gas uint64) {
	r.GrothCCSGas = gas
}

func (r *ResultCircuit) SetGrothCCSTime(grothCCSTime int) {
	r.GrothCCSTime = grothCCSTime
}

func (r *ResultCircuit) SetGrothCCSCost(grothCCSSCost *big.Int) {
	r.GrothCCSCost = grothCCSSCost
}

func (r *ResultCircuit) SetGrothCCSNumberOfConstraints(grothCCSNumberOfConstraints int64) {
	r.GrothCCSNumberOfConstraints = grothCCSNumberOfConstraints
}

func (r *ResultSetup) SetEDDSAPrivateKeySize(eddsAPrivateKeySize int) {
	r.EDDSAPrivateKeySize = eddsAPrivateKeySize
}

func (r *ResultSetup) SetEDDSAPublicKeySize(eddsAPublicKeySize int) {
	r.EDDSAPublicKeySize = eddsAPublicKeySize
}

func (r *ResultSetup) SetEDDSAPublicKeyCost(eddsaPublicKeyCost *big.Int) {
	r.EDDSAPublicKeyCost = eddsaPublicKeyCost
}

func (r *ResultSetup) SetEDDSAPublicKeyGas(gas uint64) {
	r.EDDSAPublicKeyGas = gas
}

func (r *ResultSetup) SetEDDSAKeyGenTime(eddsaKeyGentime int) {
	r.EDDSAKeyGenTime = eddsaKeyGentime
}

func (r *ResultSetup) SetGrothProvingKeySize(grothProvingKeySize int) {
	r.GrothProvingKeySize = grothProvingKeySize
}

func (r *ResultSetup) SetGrothVerificationKeySize(grothVerificationKeySize int) {
	r.GrothVerificationKeySize = grothVerificationKeySize
}

func (r *ResultSetup) SetGrothVerificationKeyCost(grothVerificationKeyCost *big.Int) {
	r.GrothVerificationKeyCost = grothVerificationKeyCost
}

func (r *ResultSetup) SetGrothVerificationKeyGas(gas uint64) {
	r.GrothVerificationKeyGas = gas
}

func (r *ResultSetup) SetContractDeploymentCost(contractDeploymentCost *big.Int) {
	r.ContractDeploymentCost = contractDeploymentCost
}
func (r *ResultSetup) SetContractDeploymentGas(gasUsed uint64) {
	r.ContractDeploymentGas = uint64(gasUsed)
}
func (r *ResultSetup) SetAvgGasPriceMarch2025(avgGasPriceMarch2025 int) {
	r.AvgGasPriceMarch2025 = avgGasPriceMarch2025
}

func (resultCircuit *ResultCircuit) SetResults(issuer issuer.Issuer) {
	resultCircuit.SetNumberOfTokensInCircuit(issuer.NumberOfTokensInCircuit)
	resultCircuit.SetGrothCCSTime(issuer.Result.GrothCCSTime)
	resultCircuit.SetGrothCCSNumberOfConstraints(issuer.Result.GrothCCSNumberOfConstraints)
	resultCircuit.SetGrothCCSSize(issuer.Result.GrothCCSSize)
	resultCircuit.SetGrothCCSGas(issuer.Result.GrothCCSGas)
	resultCircuit.SetGrothCCSCost(issuer.Result.GrothCCSCost)
}
func (result *ResultSetup) SetResults(issuer issuer.Issuer) {

	resultCircuit := ResultCircuit{}
	resultCircuit.SetResults(issuer)
	result.ZKPCircuitResults = append(result.ZKPCircuitResults, resultCircuit)

	result.SetGrothVerificationKeySize(issuer.Result.GrothVerificationKeySize)
	result.SetGrothVerificationKeyGas(issuer.Result.GrothVerificationKeyGas)
	result.SetGrothVerificationKeyCost(issuer.Result.GrothVerificationKeyCost)
	result.SetGrothProvingKeySize(issuer.Result.GrothProvingKeySize)

	result.SetEDDSAPublicKeySize(issuer.Result.EDDSAPublicKeySize)
	result.SetEDDSAPublicKeyGas(issuer.Result.EDDSAPublicKeyGas)
	result.SetEDDSAPublicKeyCost(issuer.Result.EDDSAPublicKeyCost)
	result.SetEDDSAKeyGenTime(issuer.Result.EDDSAKeyGenTime)

	result.SetAvgGasPriceMarch2025(issuer.Result.AvgGasPriceMarch2025)
	result.SetContractDeploymentCost(issuer.Result.ContractDeploymentCost)
	result.SetContractDeploymentGas(issuer.Result.ContractDeploymentGas)
}

func ComputeAverageResultSetup(results []ResultSetup) *ResultSetup {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_setup.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)
		for i := 0; i < len(results); i++ {
			var zkpResults []ResultCircuit

			for j := 0; j < len(results[i].ZKPCircuitResults); j++ {
				jsonData, _ := json.Marshal(results[i].ZKPCircuitResults[j])
				var zkpRes ResultCircuit
				json.Unmarshal(jsonData, &zkpRes)

				zkpResults = append(zkpResults, zkpRes)
			}
			results[i].ZKPCircuitResults = zkpResults
		}
	}

	result := &ResultSetup{}

	zkpResults := make(map[int]*ResultCircuit)
	firstResult := true
	for _, res := range results {
		result.GrothVerificationKeySize += res.GrothVerificationKeySize
		result.GrothVerificationKeyGas += res.GrothVerificationKeyGas
		result.GrothProvingKeySize += res.GrothProvingKeySize

		result.EDDSAKeyGenTime += res.EDDSAKeyGenTime
		result.EDDSAPublicKeySize += res.EDDSAPublicKeySize
		result.EDDSAPublicKeyGas += res.EDDSAPublicKeyGas
		result.EDDSAPrivateKeySize += res.EDDSAPrivateKeySize

		result.AvgGasPriceMarch2025 += res.AvgGasPriceMarch2025
		result.ContractDeploymentGas += res.ContractDeploymentGas

		if firstResult {
			for j := 0; j < len(res.ZKPCircuitResults); j++ {
				zkpResults[res.ZKPCircuitResults[j].NumberOfTokensInCircuit] = &ResultCircuit{
					NumberOfTokensInCircuit:     res.ZKPCircuitResults[j].NumberOfTokensInCircuit,
					GrothCCSSize:                res.ZKPCircuitResults[j].GrothCCSSize,
					GrothCCSCost:                res.ZKPCircuitResults[j].GrothCCSCost,
					GrothCCSGas:                 res.ZKPCircuitResults[j].GrothCCSGas,
					GrothCCSNumberOfConstraints: res.ZKPCircuitResults[j].GrothCCSNumberOfConstraints,
					GrothCCSTime:                res.ZKPCircuitResults[j].GrothCCSTime,
				}

			}
			result.GrothVerificationKeyCost = res.GrothVerificationKeyCost

			result.EDDSAPublicKeyCost = res.EDDSAPublicKeyCost
			result.ContractDeploymentCost = res.ContractDeploymentCost

		}

		if firstResult == false {

			result.GrothVerificationKeyCost.Add(result.GrothVerificationKeyCost, res.GrothVerificationKeyCost)

			result.EDDSAPublicKeyCost.Add(result.EDDSAPublicKeyCost, res.EDDSAPublicKeyCost)
			result.ContractDeploymentCost.Add(res.ContractDeploymentCost, res.ContractDeploymentCost)

			result.GrothVerificationKeySize = res.GrothVerificationKeySize / 2
			result.GrothVerificationKeyCost.Div(res.GrothVerificationKeyCost, big.NewInt(2))
			result.GrothVerificationKeyGas = res.GrothVerificationKeyGas / 2
			result.GrothProvingKeySize = res.GrothProvingKeySize / 2

			result.EDDSAPublicKeyCost.Div(result.EDDSAPublicKeyCost, big.NewInt(2))
			result.EDDSAKeyGenTime = res.EDDSAKeyGenTime / 2
			result.EDDSAPublicKeySize = res.EDDSAPublicKeySize / 2
			result.EDDSAPublicKeyGas = res.EDDSAPublicKeyGas / 2
			result.EDDSAPrivateKeySize = res.EDDSAPrivateKeySize / 2

			result.AvgGasPriceMarch2025 = res.AvgGasPriceMarch2025 / 2
			result.ContractDeploymentCost.Div(res.ContractDeploymentCost, big.NewInt(2))
			result.ContractDeploymentGas = res.ContractDeploymentGas / 2

			for j := 0; j < len(result.ZKPCircuitResults); j++ {
				temp := result.ZKPCircuitResults[j]
				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSSize += result.ZKPCircuitResults[j].GrothCCSSize
				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSSize = zkpResults[temp.NumberOfTokensInCircuit].GrothCCSSize / 2

				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSGas += result.ZKPCircuitResults[j].GrothCCSGas
				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSGas = zkpResults[temp.NumberOfTokensInCircuit].GrothCCSGas / 2

				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSNumberOfConstraints += result.ZKPCircuitResults[j].GrothCCSNumberOfConstraints
				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSNumberOfConstraints = zkpResults[temp.NumberOfTokensInCircuit].GrothCCSNumberOfConstraints / 2

				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSTime += result.ZKPCircuitResults[j].GrothCCSTime
				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSTime = zkpResults[temp.NumberOfTokensInCircuit].GrothCCSTime / 2

				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSCost.Add(zkpResults[temp.NumberOfTokensInCircuit].GrothCCSCost, result.ZKPCircuitResults[j].GrothCCSCost)
				zkpResults[temp.NumberOfTokensInCircuit].GrothCCSCost.Div(zkpResults[temp.NumberOfTokensInCircuit].GrothCCSCost, big.NewInt(2))
			}

		}
		firstResult = false
	}
	for _, zkpResult := range zkpResults {
		res := ResultCircuit{
			NumberOfTokensInCircuit:     zkpResult.NumberOfTokensInCircuit,
			GrothCCSSize:                zkpResult.GrothCCSSize,
			GrothCCSCost:                zkpResult.GrothCCSCost,
			GrothCCSGas:                 zkpResult.GrothCCSGas,
			GrothCCSNumberOfConstraints: zkpResult.GrothCCSNumberOfConstraints,
			GrothCCSTime:                zkpResult.GrothCCSTime,
		}
		result.ZKPCircuitResults = append(result.ZKPCircuitResults, res)
	}

	return result
}

func WriteResultSetupToFile(result ResultSetup, isavg bool) {

	var results []ResultSetup
	filename := fmt.Sprintf("benchmark/results/result_setup.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_setup_avg.json")
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

func (r *ResultSetup) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultSetup(jsonObj []byte) *ResultSetup {
	res := ResultSetup{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func (r ResultSetup) String() string {
	var response string

	response = response + "Groth16 proving Key size (in bytes): " + fmt.Sprintf("%d", r.GrothProvingKeySize) + "\n"
	response = response + "Groth16 verification Key size (in bytes): " + fmt.Sprintf("%d", r.GrothVerificationKeySize) + "\n"
	response = response + "Gas used to store Groth16 verification key:" + fmt.Sprintf("%d", r.GrothVerificationKeyGas) + "\n"
	response = response + "Cost to store Groth16 verification key (in wei):" + fmt.Sprintf("%s", r.GrothVerificationKeyCost.String()) + "\n"

	response = response + "EDDSA Public Key size (in bytes): " + fmt.Sprintf("%d", r.EDDSAPublicKeySize) + "\n"
	response = response + "Cost to store EDDSA public key (in wei):" + fmt.Sprintf("%s", r.EDDSAPublicKeyCost.String()) + "\n"
	response = response + "Gas used to store EDDSA public key:" + fmt.Sprintf("%d", r.EDDSAPublicKeyGas) + "\n"
	response = response + "EDDSA Private Key size (in bytes): " + fmt.Sprintf("%d", r.EDDSAPrivateKeySize) + "\n"

	response = response + "Cost to deploy smart contract (in wei):" + fmt.Sprintf("%s", r.ContractDeploymentCost.String()) + "\n"
	response = response + "Gas used to deploy smart contract:" + fmt.Sprintf("%d", r.ContractDeploymentGas) + "\n"

	for i := 0; i < len(r.ZKPCircuitResults); i++ {
		response = response + "{ \n Number of tokens in circuit: " + fmt.Sprintf("%d", r.ZKPCircuitResults[i].NumberOfTokensInCircuit) + "\n "
		response = response + "GrothCCCSSize (in bytes): " + fmt.Sprintf("%d", r.ZKPCircuitResults[i].GrothCCSSize) + "\n"
		response = response + "Cost to store Groth16 constraint system (in wei):" + fmt.Sprintf("%s", r.ZKPCircuitResults[i].GrothCCSCost.String()) + "\n"
		response = response + "Gas to store Groth16 constraint system (in wei):" + fmt.Sprintf("%d", r.ZKPCircuitResults[i].GrothCCSGas) + "\n"
		response = response + "Time to generate Groth16 constraint system (in micro seconds):" + fmt.Sprintf("%d", r.ZKPCircuitResults[i].GrothCCSTime) + "\n"
		response = response + "number of constraints in groth16 ccs:" + fmt.Sprintf("%d", r.ZKPCircuitResults[i].GrothCCSNumberOfConstraints) + "} \n"
	}

	return response
}

func ResetSetupFiles() {
	filename1 := fmt.Sprintf("benchmark/results/result_setup.json")
	filename2 := fmt.Sprintf("benchmark/results/result_setup_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}
