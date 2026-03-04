package config

type Benchmark struct {
	Setup                     bool `json:"setup"`
	Issaunce                  bool `json:"issaunce"`
	Revocation                bool `json:"revocation"`
	Refresh                   bool `json:"refresh"`
	Presentation_Verification bool `json:"presentation_verification"`

	IRMASetup                       bool `json:"irma_setup"`
	IRMAIssuance                    bool `json:"irma_issuance"`
	IRMARevocation                  bool `json:"irma_revocation"`
	IRMAPresentationAndVerification bool `json:"irma_presentation_and_verification"`
	CircuitConstraints              bool `json:"circuit_constraints"`
	ListCommitment                  bool `json:"list_commitment"`
}

func (benchmark Benchmark) String() string {
	var response string
	response = response + "\n\n--------------------------------------------------------"
	if benchmark.Setup {
		response = response + "\n  Benchark Setup "
	}
	if benchmark.Issaunce {
		response = response + " \n Benchark Issuance "
	}
	if benchmark.Refresh {
		response = response + " \n Benchark Refresh "
	}
	if benchmark.Presentation_Verification {
		response = response + " \n Benchark Verification "
	}
	response = response + "--------------------------------------------------"
	return response
}
