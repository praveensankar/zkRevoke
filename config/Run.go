package config

type Run struct {
	Flow               bool `json:"flow"`
	IRMATest           bool `json:"irma_test"`
	ComputeFinalResult bool `json:"compute_final_result"`
	CircuitTest        bool `json:"circuit"`
	VCTest             bool `json:"VCTest"`
	ZKPTest            bool `json:"zkp_test"`
	CryptoTest         bool `json:"crypto_test"`
	GenerateVCs        bool
}

func (run Run) String() string {
	var response string
	response = response + "\n\n--------------------------------------------------------"
	if run.Flow {
		response = response + " Running Flow "
	}
	if run.ComputeFinalResult {
		response = response + " Running Computer Final Result "
	}

	if run.IRMATest {
		response = response + " Running IRMA Test "
	}
	if run.CircuitTest {
		response = response + " Testing ZKP Circuit "
	}
	if run.VCTest {
		response = response + " Running VC Test "
	}
	if run.ZKPTest {
		response = response + " Running ZKP Test "
	}
	if run.CryptoTest {
		response = response + " Running Crypto Test"
	}
	if run.GenerateVCs {
		response = response + " Generating VC data set "
	}
	response = response + "--------------------------------------------------"
	return response
}
