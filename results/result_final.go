package results

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"os"
)

type ResultFinal struct {
	TotalVCs                int `json:"totalVCs"`
	NumberOfRevokedVCs      int `json:"numberOfRevokedVCs"`
	TotalEpochs             int `json:"totalEpochs"`
	EpochDuration           int `json:"epochDuration"`
	NumberOfTokensInCircuit int `json:"numberOfTokensInCircuit"`

	VCGenerationTime            []int `json:"vcGenerationTime"`
	VCGenerationTimeAvg         int   `json:"vcGenerationTimeAvg"`
	VPGenerationTime            []int `json:"vPGenerationTime"`
	VPGenerationTimeAvg         int   `json:"vPGenerationTimeAvg"`
	GrothProofGenerationTime    []int `json:"grothproofGenerationTime"`
	GrothProofGenerationTimeAvg int   `json:"grothproofGenerationTimeAvg"`

	VPVerificationMetrics []VPVerificationMetrics `json:"vPVerificationMetrics"`

	VCSize            []int              `json:"-"`
	VCSizeAvg         int                `json:"vcSizeAvg"`
	VPSizeMetrics     []VPRelatedMetrics `json:"vpRelatedMetrics"`
	GrothProofSize    []int              `json:"-"`
	GrothProofSizeAvg int                `json:"grothproofSizeAvg"`

	TokenGenerationTime            []int `json:"tokenGenerationTime"`
	TokenGenerationTimeAvg         int   `json:"tokenGenerationTimeAvg"`
	TokenGenerationTimeAtHolder    []int `json:"tokenGenerationTimeAtHolder"`
	TokenGenerationTimeAtHolderAvg int   `json:"tokenGenerationTimeAtHolderAvg"`

	SizeOfRevocationList int `json:"sizeOfRevocationList"`

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

	EDDSAPrivateKeySize    int      `json:"eddsaPrivateKeySize"`
	EDDSAPublicKeySize     int      `json:"eddsaPublicKeySize"`
	EDDSAPublicKeyCost     *big.Int `json:"eddsaPublicKeyCost"`
	EDDSAPublicKeyGas      uint64   `json:"eddsaPublicKeyGas"`
	ContractDeploymentCost *big.Int `json:"contractDeploymentCost"`
	ContractDeploymentGas  uint64   `json:"contractDeploymentGas"`
	AvgGasPriceMarch2025   int      `json:"avgGasPriceMarch2025"`
}

func (r ResultFinal) String() string {
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
	response = response + "Gas used to deploy smart contract:" + fmt.Sprintf("%s", r.ContractDeploymentGas) + "\n"

	response = response + "Time to generate tokens at holder (in micro seconds):" + fmt.Sprintf("%v", r.TokenGenerationTimeAtHolder) + "\n"
	response = response + "Average time to generate tokens at holder (in micro seconds):" + fmt.Sprintf("%v", r.TokenGenerationTimeAtHolderAvg) + "\n"
	response = response + "Time to generate VPs (in micro seconds):" + fmt.Sprintf("%v", r.VPGenerationTime) + "\n"
	response = response + "Average time to generate VPs (in micro seconds):" + fmt.Sprintf("%v", r.VPGenerationTimeAvg) + "\n"

	response = response + "Groth16 proof size (in bytes): " + fmt.Sprintf("%v", r.GrothProofSize) + "\n"
	response = response + "Average time to generate groth16 proof (in micro seconds):" + fmt.Sprintf("%v", r.GrothProofGenerationTimeAvg) + "\n"

	return response
}

func ComputeFinalResult(result_issuer *ResultIssuer, result_verifier *ResultVerifier, results_holder []*ResultHolder) *ResultFinal {

	result := &ResultFinal{}
	result.TotalVCs = result_issuer.TotalVCs
	result.NumberOfRevokedVCs = result_issuer.NumberOfRevokedVCs
	result.TotalEpochs = result_issuer.TotalEpochs
	result.EpochDuration = result_issuer.EpochDuration
	result.NumberOfTokensInCircuit = result_issuer.NumberOfTokensInCircuit

	result.VCGenerationTime = result_issuer.VCGenerationTime
	result.VCGenerationTimeAvg = result_issuer.VCGenerationTimeAvg
	result.VCSize = result_issuer.VCSize
	result.VCSizeAvg = result_issuer.VCSizeAvg
	result.TokenGenerationTime = result_issuer.TokenGenerationTime
	result.TokenGenerationTimeAvg = result_issuer.TokenGenerationTimeAvg

	var token_gen_time_at_holder []int
	var token_gen_time_at_holder_avg int

	var vp_gen_time []int
	var vp_gen_time_avg int
	var vp_storage_metrics []VPRelatedMetrics

	var groth_proof_gen_time []int
	var groth_proof_gen_time_avg int
	var groth_proof_size []int
	var groth_proof_size_avg int

	i := 0
	for i = 0; i < len(results_holder); i++ {
		token_gen_time_at_holder = append(token_gen_time_at_holder, results_holder[i].TokenGenerationTime...)

		vp_storage_metrics = append(vp_storage_metrics, results_holder[i].VPSizeMetrics...)
		groth_proof_gen_time = append(groth_proof_gen_time, results_holder[i].GrothProofGenerationTime...)
		groth_proof_size = append(groth_proof_size, results_holder[i].GrothProofSize...)

		if i == 0 {
			token_gen_time_at_holder_avg = results_holder[i].TokenGenerationTimeAvg

			groth_proof_gen_time_avg = results_holder[i].GrothProofGenerationTimeAvg
			groth_proof_size_avg = results_holder[i].GrothProofSizeAvg

		} else {
			token_gen_time_at_holder_avg = (token_gen_time_at_holder_avg + results_holder[i].TokenGenerationTimeAvg) / 2

			groth_proof_gen_time_avg = (groth_proof_gen_time_avg + results_holder[i].GrothProofGenerationTimeAvg) / 2
			groth_proof_size_avg = (groth_proof_size_avg + results_holder[i].GrothProofSizeAvg) / 2

		}
	}

	result.TokenGenerationTimeAtHolder = token_gen_time_at_holder
	result.TokenGenerationTimeAtHolderAvg = token_gen_time_at_holder_avg
	result.VPGenerationTime = vp_gen_time
	result.VPGenerationTimeAvg = vp_gen_time_avg
	result.VPSizeMetrics = vp_storage_metrics
	result.GrothProofGenerationTime = groth_proof_gen_time
	result.GrothProofGenerationTimeAvg = groth_proof_gen_time_avg
	result.GrothProofSize = groth_proof_size
	result.GrothProofSizeAvg = groth_proof_size_avg

	result.SizeOfRevocationList = result_verifier.SizeOfRevocationList
	result.VPVerificationMetrics = result_verifier.VPMetrics
	result.GrothCCSSize = result_issuer.GrothCCSSize
	result.GrothCCSCost = result_issuer.GrothCCSCost
	result.GrothCCSGas = result_issuer.GrothCCSGas
	result.GrothCCSNumberOfConstraints = result_issuer.GrothCCSNumberOfConstraints
	result.GrothCCSTime = result_issuer.GrothCCSTime
	result.GrothCCSHashCost = result_issuer.GrothCCSHashCost
	result.GrothProvingKeySize = result_issuer.GrothProvingKeySize
	result.GrothVerificationKeySize = result_issuer.GrothVerificationKeySize
	result.GrothVerificationKeyCost = result_issuer.GrothVerificationKeyCost
	result.GrothVerificationKeyGas = result_issuer.GrothVerificationKeyGas

	result.EDDSAPublicKeySize = result_issuer.EDDSAPublicKeySize
	result.EDDSAPublicKeyCost = result_issuer.EDDSAPublicKeyCost
	result.EDDSAPublicKeyGas = result_issuer.EDDSAPublicKeyGas

	result.ContractDeploymentCost = result_issuer.ContractDeploymentCost
	result.ContractDeploymentGas = result_issuer.ContractDeploymentGas
	result.AvgGasPriceMarch2025 = result_issuer.AvgGasPriceMarch2025
	return result
}

func WriteResultToFile(result ResultFinal) {

	var results []ResultFinal
	//filename := fmt.Sprintf("results/results_computed.json")
	filename := fmt.Sprintf("results/result_%d_%d.json", result.TotalVCs, result.EpochDuration)
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

func (r *ResultFinal) Json() ([]byte, error) {
	return json.MarshalIndent(r, "***", "----")
	//return json.Marshal(r)
}
