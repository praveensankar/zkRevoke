package results

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

type ResultIssuer struct {
	TotalVCs                int `json:"totalVCs"`
	NumberOfRevokedVCs      int `json:"numberOfRevokedVCs"`
	TotalEpochs             int `json:"totalEpochs"`
	EpochDuration           int `json:"epochDuration"`
	NumberOfTokensInCircuit int `json:"numberOfTokensInCircuit"`

	VCGenerationTime       []int `json:"vcGenerationTime"`
	VCGenerationTimeAvg    int   `json:"vcGenerationTimeAvg"`
	VCSize                 []int `json:"vcSize"`
	VCSizeAvg              int   `json:"vcSizeAvg"`
	TokenGenerationTime    []int `json:"tokenGenerationTime"`
	TokenGenerationTimeAvg int   `json:"tokenGenerationTimeAvg"`

	GrothCCSSize                int      `json:"grothCcsSize"`
	GrothCCSCost                *big.Int `json:"GrothccsCost"`
	GrothCCSGas                 uint64   `json:"grothCcsGas"`
	GrothCCSNumberOfConstraints int64    `json:"grothccsNumberOfConstraints"`
	GrothCCSTime                int      `json:"grothccsTime"`
	GrothCCSHashCost            *big.Int `json:"GrothccsHashCost"`

	GrothProvingKeySize      int      `json:"grothProvingKeySize"`
	GrothVerificationKeySize int      `json:"grothVerificationKeySize"`
	GrothVerificationKeyCost *big.Int `json:"grothVerificationKeyCost"`
	GrothVerificationKeyGas  uint64   `json:"grothVerificationKeyGas"`

	EDDSAPrivateKeySize int      `json:"eddsaPrivateKeySize"`
	EDDSAPublicKeySize  int      `json:"eddsaPublicKeySize"`
	EDDSAPublicKeyCost  *big.Int `json:"eddsaPublicKeyCost"`
	EDDSAPublicKeyGas   uint64   `json:"eddsaPublicKeyGas"`
	EDDSAKeyGenTime     int      `json:"eddsaKeyGenTime"`
	EDDSASignTime       int      `json:"eddsa_sign_time"`

	ContractDeploymentCost *big.Int `json:"contractDeploymentCost"`
	ContractDeploymentGas  uint64   `json:"contractDeploymentGas"`
	AvgGasPriceMarch2025   int      `json:"avgGasPriceMarch2025"`
}

func (r *ResultIssuer) SetTotalVCs(totalVCs int) {
	r.TotalVCs = totalVCs
}

func (r *ResultIssuer) SetRevokedVCs(revokedVCs int) {
	r.NumberOfRevokedVCs = revokedVCs
}

func (r *ResultIssuer) SetTotalEpochs(totalEpochs int) {
	r.TotalEpochs = totalEpochs
}

func (r *ResultIssuer) SetEpochDuration(epochDuration int) {
	r.EpochDuration = epochDuration
}

func (r *ResultIssuer) SetGrothCCSSize(grothCCSSize int) {
	r.GrothCCSSize = grothCCSSize
}

func (r *ResultIssuer) SetGrothCCSGas(gas uint64) {
	r.GrothCCSGas = gas
}

func (r *ResultIssuer) SetGrothCCSTime(grothCCSTime time.Duration) {
	r.GrothCCSTime = int(grothCCSTime.Microseconds())
}

func (r *ResultIssuer) SetGrothCCSHashCost(grothCCSSHashCost string) {
	res, _ := new(big.Int).SetString(grothCCSSHashCost, 10)
	r.GrothCCSHashCost = res
}

func (r *ResultIssuer) SetGrothCCSCost(grothCCSSCost string) {
	res, _ := new(big.Int).SetString(grothCCSSCost, 10)
	r.GrothCCSCost = res
}

func (r *ResultIssuer) SetGrothCCSNumberOfConstraints(grothCCSNumberOfConstraints int64) {
	r.GrothCCSNumberOfConstraints = grothCCSNumberOfConstraints
}

func (r *ResultIssuer) SetEDDSASignTime(signTime time.Duration) {
	r.EDDSASignTime = int(signTime.Microseconds())
}

func (r *ResultIssuer) SetEDDSAPrivateKeySize(eddsAPrivateKeySize int) {
	r.EDDSAPrivateKeySize = eddsAPrivateKeySize
}

func (r *ResultIssuer) SetEDDSAPublicKeySize(eddsAPublicKeySize int) {
	r.EDDSAPublicKeySize = eddsAPublicKeySize
}

func (r *ResultIssuer) SetEDDSAPublicKeyCost(eddsaPublicKeyCost string) {
	res, _ := new(big.Int).SetString(eddsaPublicKeyCost, 10)
	r.EDDSAPublicKeyCost = res
}

func (r *ResultIssuer) SetEDDSAKeyGenTime(eddsaKeyGentime time.Duration) {
	r.EDDSAKeyGenTime = int(eddsaKeyGentime.Microseconds())
}

func (r *ResultIssuer) SetEDDSAPublicKeyGas(gas uint64) {
	r.EDDSAPublicKeyGas = gas
}

func (r *ResultIssuer) SetGrothProvingKeySize(grothProvingKeySize int) {
	r.GrothProvingKeySize = grothProvingKeySize
}

func (r *ResultIssuer) SetGrothVerificationKeySize(grothVerificationKeySize int) {
	r.GrothVerificationKeySize = grothVerificationKeySize
}

func (r *ResultIssuer) SetGrothVerificationKeyCost(grothVerificationKeyCost string) {
	res, _ := new(big.Int).SetString(grothVerificationKeyCost, 10)
	r.GrothVerificationKeyCost = res
}

func (r *ResultIssuer) SetGrothVerificationKeyGas(gas uint64) {
	r.GrothVerificationKeyGas = gas
}

func (r *ResultIssuer) SetContractDeploymentCost(contractDeploymentCost string) {
	res, _ := new(big.Int).SetString(contractDeploymentCost, 10)
	r.ContractDeploymentCost = res
}
func (r *ResultIssuer) SetContractDeploymentGas(gasUsed uint64) {
	r.ContractDeploymentGas = uint64(gasUsed)
}
func (r *ResultIssuer) SetAvgGasPriceMarch2025(avgGasPriceMarch2025 int) {
	r.AvgGasPriceMarch2025 = avgGasPriceMarch2025
}

func (r *ResultIssuer) SetNumberOfTokensInCircuit(numberOfTokensInCircuit int) {
	r.NumberOfTokensInCircuit = numberOfTokensInCircuit
}

func (r *ResultIssuer) AddTokenGenerationTime(time time.Duration) {
	r.TokenGenerationTime = append(r.TokenGenerationTime, int(time.Microseconds()))
}

func (r *ResultIssuer) ComputeAvgTokenGenerationTime() {
	res := 0

	if len(r.TokenGenerationTime) > 0 {
		for i := 0; i < len(r.TokenGenerationTime); i++ {
			res += r.TokenGenerationTime[i]
		}
		res = res / len(r.TokenGenerationTime)
	}

	r.TokenGenerationTimeAvg = res
}

func (r *ResultIssuer) AddVCGenerationTime(time time.Duration) {
	r.VCGenerationTime = append(r.VCGenerationTime, int(time.Microseconds()))
}

func (r *ResultIssuer) ComputeAvgVCGenerationTime() {
	res := 0
	if len(r.VCGenerationTime) > 0 {
		for i := 0; i < len(r.VCGenerationTime); i++ {
			res += r.VCGenerationTime[i]
		}
		res = res / len(r.VCGenerationTime)
	}
	r.VCGenerationTimeAvg = res
}

func (r *ResultIssuer) AddVCSize(size int) {
	r.VCSize = append(r.VCSize, size)
}

func (r *ResultIssuer) ComputeAvgVCSize() {
	res := 0
	for i := 0; i < len(r.VCSize); i++ {
		res += r.VCSize[i]
	}
	res = res / len(r.VCSize)
	r.VCSizeAvg = res
}
func WriteResultIssuerToFile(result ResultIssuer) {

	var results []ResultIssuer
	//filename := fmt.Sprintf("results/results_computed.json")
	filename := fmt.Sprintf("results/result_%d_%d_%d.json", result.TotalVCs, result.NumberOfRevokedVCs, result.EpochDuration)
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
	zap.S().Info(result.String())
}

func (r *ResultIssuer) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultIssuer(jsonObj []byte) *ResultIssuer {
	res := ResultIssuer{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func (r ResultIssuer) String() string {
	var response string
	response = response + "Total VCs: " + fmt.Sprintf("%d", r.TotalVCs) + "\n"
	response = response + "Revoked VCs: " + fmt.Sprintf("%d", r.NumberOfRevokedVCs) + "\n"
	response = response + "Total Epochs: " + fmt.Sprintf("%d", r.TotalEpochs) + "\n"
	response = response + "Epoch Duration (in seconds): " + fmt.Sprintf("%d", r.EpochDuration) + "\n"

	response = response + "GrothCCCSSize (in bytes): " + fmt.Sprintf("%d", r.GrothCCSSize) + "\n"
	response = response + "Cost to store Groth16 constraint system (in wei):" + fmt.Sprintf("%s", r.GrothCCSCost.String()) + "\n"
	response = response + "Gas to store Groth16 constraint system (in wei):" + fmt.Sprintf("%d", r.GrothCCSGas) + "\n"
	response = response + "Time to generate Groth16 constraint system (in micro seconds):" + fmt.Sprintf("%d", r.GrothCCSTime) + "\n"
	response = response + "Cost to store Groth16 hash of constraint system (in wei):" + fmt.Sprintf("%s", r.GrothCCSHashCost.String()) + "\n"
	response = response + "number of constraints in groth16 ccs:" + fmt.Sprintf("%d", r.GrothCCSNumberOfConstraints) + "\n"
	response = response + "number of constraints in groth16 ccs:" + fmt.Sprintf("%d", r.GrothCCSNumberOfConstraints) + "\n"

	response = response + "Groth16 proving Key size (in bytes): " + fmt.Sprintf("%d", r.GrothProvingKeySize) + "\n"
	response = response + "Groth16 verification Key size (in bytes): " + fmt.Sprintf("%d", r.GrothVerificationKeySize) + "\n"
	response = response + "Gas used to store Groth16 verification key:" + fmt.Sprintf("%d", r.GrothVerificationKeyGas) + "\n"
	response = response + "Cost to store Groth16 verification key (in wei):" + fmt.Sprintf("%s", r.GrothVerificationKeyCost.String()) + "\n"

	response = response + "EDDSA Public Key size (in bytes): " + fmt.Sprintf("%d", r.EDDSAPublicKeySize) + "\n"
	response = response + "Cost to store EDDSA public key (in wei):" + fmt.Sprintf("%s", r.EDDSAPublicKeyCost.String()) + "\n"
	response = response + "Gas used to store EDDSA public key:" + fmt.Sprintf("%d", r.EDDSAPublicKeyGas) + "\n"
	response = response + "EDDSA Private Key size (in bytes): " + fmt.Sprintf("%d", r.EDDSAPrivateKeySize) + "\n"

	response = response + "Time to generate tokens (in micro seconds):" + fmt.Sprintf("%v", r.TokenGenerationTime) + "\n"
	response = response + "Average time to generate tokens (in micro seconds):" + fmt.Sprintf("%v", r.TokenGenerationTimeAvg) + "\n"
	response = response + "Time to generate VCs (in micro seconds):" + fmt.Sprintf("%v", r.VCGenerationTime) + "\n"
	response = response + "Average time to generate VCs (in micro seconds):" + fmt.Sprintf("%v", r.VCGenerationTimeAvg) + "\n"
	response = response + "VC size (in bytes): " + fmt.Sprintf("%v", r.VCSize) + "\n"
	response = response + "Average VC size (in bytes): " + fmt.Sprintf("%v", r.VCSizeAvg) + "\n"

	response = response + "Cost to deploy smart contract (in wei):" + fmt.Sprintf("%s", r.ContractDeploymentCost.String()) + "\n"
	response = response + "Gas used to deploy smart contract:" + fmt.Sprintf("%d", r.ContractDeploymentGas) + "\n"
	return response
}
