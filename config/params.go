package config

type Params struct {
	NumberOfExperiments                       int  `json:"number_of_experiments"`
	TotalVCs                                  int  `json:"total_vcs"`
	ExpirationPeriod                          int  `json:"expiration_period"`
	VerificationPeriod                        int  `json:"verification_period"`
	NumberOfTokensPerCircuit                  int  `json:"number_of_tokens_per_circuit"`
	RevocationRateBase                        int  `json:"revocation_rate_base"`
	RevocationRateStep                        int  `json:"revocation_rate_step"`
	RevocationRateEnd                         int  `json:"revocation_rate_end"`
	EpochDuration                             int  `json:"epoch_duration"`
	IRMAWitnessUpdateMessageWithoutRepetition bool `json:"irma_witness_update_message_without_repetition"`
}
