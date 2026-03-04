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

/*
HolderBandwidth: The total witness messages received by a holder + the total non-revocation proofs send by the holder
***************** for i:1,...,m: res = res + witness_size[i]+non_rev_proof_size[i]
*/
type ResultPresentationAndVerification struct {
	TotalNumberOfIssuedVCs                       int `json:"total_number_of_issued_vcs"`
	RevocationRate                               int `json:"revocation_rate"`
	TotalRevokedVCsPerEpoch                      int `json:"total_revoked_vcs_per_epoch"`
	TimeToUpdateWitness                          int `json:"time_to_update_witness"`
	TotalNumberOfEpochs                          int `json:"total_number_of_epochs"`
	VPValidityPeriod                             int `json:"vp_validity_period"`
	TimeToCreateDisclosureProofWithNonRevocation int `json:"time_to_create_disclosure_proof_with_non_revocation"`
	CurrentEpoch                                 int `json:"current_epoch"`
	HolderBandwidth                              int `json:"holder_bandwidth"`
	DisclosureProofSize                          int `json:"disclosure_proof_size"`
	NonRevocationProofSize                       int `json:"non_revocation_proof_size"`
	TotalWitnessUpdateMessagesReceived           int `json:"total_witness_update_messages_received"`
	ProofVerificationTime                        int `json:"proof_verification_time"`
}

type ResultPresentationAndVerificationList struct {
	Results []ResultPresentationAndVerification
	mutex   sync.Mutex
}

func (list *ResultPresentationAndVerificationList) Add(result ResultPresentationAndVerification, witnessUpdateMessageWithoutRepetition bool) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.Results = append(list.Results, result)
	WriteResultPresentationAndVerificationToFile(result, false, witnessUpdateMessageWithoutRepetition)
}

func (list *ResultPresentationAndVerificationList) Get(index int) ResultPresentationAndVerification {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	return list.Results[index]
}

func ComputeAverageResultPresentationAndVerification(results []ResultPresentationAndVerification, witnessUpdateMessageWithoutRepetition bool) []*ResultPresentationAndVerification {

	if len(results) == 0 {
		filename := fmt.Sprintf("irma/benchmark/results/result_presentation_verification.json")
		if witnessUpdateMessageWithoutRepetition == true {
			filename = fmt.Sprintf("irma/benchmark/results/result_presentation_verification_witness_update_without_repetition.json")
		}
		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	//key: TotalRevokedVCsPerEpoch
	AvgResults := make(map[string]*ResultPresentationAndVerification)

	//key: TotalRevokedVCsPerEpoch
	isExists := make(map[string]bool)

	for _, res := range results {

		id := fmt.Sprintf(strconv.Itoa(res.TotalNumberOfIssuedVCs), "#", strconv.Itoa(res.RevocationRate),
			"#", strconv.Itoa(res.TotalNumberOfEpochs), "#", strconv.Itoa(res.CurrentEpoch))
		_, okay := isExists[id]
		if okay {
			AvgResults[id].TimeToUpdateWitness = (AvgResults[id].TimeToUpdateWitness + res.TimeToUpdateWitness) / 2
			AvgResults[id].TimeToCreateDisclosureProofWithNonRevocation = (AvgResults[id].TimeToCreateDisclosureProofWithNonRevocation + res.TimeToCreateDisclosureProofWithNonRevocation) / 2
			AvgResults[id].DisclosureProofSize = (AvgResults[id].DisclosureProofSize + res.DisclosureProofSize) / 2
			AvgResults[id].NonRevocationProofSize = (AvgResults[id].NonRevocationProofSize + res.NonRevocationProofSize) / 2
			AvgResults[id].ProofVerificationTime = (AvgResults[id].ProofVerificationTime + res.ProofVerificationTime) / 2
			AvgResults[id].TotalNumberOfEpochs = (AvgResults[id].TotalNumberOfEpochs + res.TotalNumberOfEpochs) / 2
			AvgResults[id].HolderBandwidth = (AvgResults[id].HolderBandwidth + res.HolderBandwidth) / 2
			AvgResults[id].TotalWitnessUpdateMessagesReceived = (AvgResults[id].TotalWitnessUpdateMessagesReceived + res.TotalWitnessUpdateMessagesReceived) / 2
		} else {
			result := &ResultPresentationAndVerification{}
			result.TotalNumberOfIssuedVCs = res.TotalNumberOfIssuedVCs
			result.VPValidityPeriod = res.VPValidityPeriod
			result.RevocationRate = res.RevocationRate
			result.TotalRevokedVCsPerEpoch = res.TotalRevokedVCsPerEpoch
			result.HolderBandwidth = res.HolderBandwidth
			result.TimeToUpdateWitness = res.TimeToUpdateWitness
			result.TimeToCreateDisclosureProofWithNonRevocation = res.TimeToCreateDisclosureProofWithNonRevocation
			result.DisclosureProofSize = res.DisclosureProofSize
			result.NonRevocationProofSize = res.NonRevocationProofSize
			result.ProofVerificationTime = res.ProofVerificationTime
			result.TotalNumberOfEpochs = res.TotalNumberOfEpochs
			result.TotalWitnessUpdateMessagesReceived = res.TotalWitnessUpdateMessagesReceived
			result.CurrentEpoch = res.CurrentEpoch
			AvgResults[id] = result
			isExists[id] = true
		}

	}

	var finalResults []*ResultPresentationAndVerification
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func WriteResultPresentationAndVerificationToFile(result ResultPresentationAndVerification, isavg bool, witnessUpdateMessageWithoutRepetition bool) {

	var results []ResultPresentationAndVerification
	filename := fmt.Sprintf("irma/benchmark/results/result_presentation_verification.json")

	if isavg == true && witnessUpdateMessageWithoutRepetition == false {
		filename = fmt.Sprintf("irma/benchmark/results/result_presentation_verification_avg.json")
	}
	if isavg == false && witnessUpdateMessageWithoutRepetition == true {
		filename = fmt.Sprintf("irma/benchmark/results/result_presentation_verification_witness_update_without_repetition.json")
	}
	if isavg == true && witnessUpdateMessageWithoutRepetition == true {
		filename = fmt.Sprintf("irma/benchmark/results/result_presentation_verification_avg_witness_update_without_repetition.json")
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
}

func ResetPresentationAndVerificationFiles(witnessUpdateMessageWithoutRepetition bool) {

	if witnessUpdateMessageWithoutRepetition == false {
		filename1 := fmt.Sprintf("irma/benchmark/results/result_presentation_verification.json")
		filename2 := fmt.Sprintf("irma/benchmark/results/result_presentation_verification_avg.json")
		os.Remove(filename1)
		os.Remove(filename2)
	}

	if witnessUpdateMessageWithoutRepetition == true {
		filename3 := fmt.Sprintf("irma/benchmark/results/result_presentation_verification_witness_update_without_repetition.json")
		filename4 := fmt.Sprintf("irma/benchmark/results/result_presentation_verification_avg_witness_update_without_repetition.json")
		os.Remove(filename3)
		os.Remove(filename4)
	}
}
