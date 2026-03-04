package benchmark

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"go.uber.org/zap"
	"math"
	"sync"
	"sync/atomic"
	"time"
	blockchain_hardhat "zkrevoke/blockchain-hardhat"
	"zkrevoke/config"
	"zkrevoke/holder"
	"zkrevoke/issuer"
	"zkrevoke/results"
	"zkrevoke/utils"
	"zkrevoke/verifier"
)

func Benchmark_Setup(conf config.Config) {

	conf.UsePreGeneratedKeysAndVCs = false
	conf.Run.GenerateVCs = false
	MaxTokensInCircuit := conf.Params.VerificationPeriod
	ResetSetupFiles()

	conf.Params.NumberOfExperiments = conf.Params.NumberOfExperiments + 1
	numberOfExperiments := conf.Params.NumberOfExperiments
	if numberOfExperiments > 3 {
		numberOfExperiments = 3
	}

	var finalResults []ResultSetup
	start := time.Now()
	for i := 0; i < numberOfExperiments; i++ {
		result := ResultSetup{}

		for numberOfTokensInCircuit := 1; numberOfTokensInCircuit <= MaxTokensInCircuit; numberOfTokensInCircuit++ {
			zap.S().Infoln("\n \n *******************************************************************************************")
			zap.S().Infoln("\n  *****BENCHMARK_SETUP: \t experiment: ", i+1, "\t number of tokens in a circuit: ", numberOfTokensInCircuit, "******")
			conf.Params.NumberOfTokensPerCircuit = int(uint(numberOfTokensInCircuit))

			issuer := &issuer.Issuer{}
			issuer.Result = &results.ResultIssuer{}
			issuer.Setup(&conf)

			address := issuer.DeployContract(conf)
			conf.SmartContractAddress = address
			issuer.SetUpBlockchainConnection(conf)
			issuer.PublishPublicParameters()
			issuer.SetupResults(conf)

			if numberOfTokensInCircuit == 1 {
				result.SetResults(*issuer)
			} else {
				resultCircuit := ResultCircuit{}
				resultCircuit.SetResults(*issuer)
				result.ZKPCircuitResults = append(result.ZKPCircuitResults, resultCircuit)
			}
		}
		if i > 0 {
			finalResults = append(finalResults, result)
			WriteResultSetupToFile(result, false)
		}
	}

	avgResults := ComputeAverageResultSetup(finalResults)
	WriteResultSetupToFile(*avgResults, true)
	end := time.Since(start)
	zap.S().Infoln("Benchmark_Setup took: ", end.Minutes(), " minutes")
}

/*
Benchmark_Issuance measures the sizes of: (a) seed, (b) sig
It also measures the computation time of: (a) sig

It does not measure the size of claims, validFrom and validUntil since they are not
generated as part of issuance, rather they are inputs to issuance.
*/
func Benchmark_Issuance(conf config.Config) {

	var finalResults []ResultIssuance
	ResetIssuanceFiles()
	start := time.Now()

	for i := 0; i < conf.Params.NumberOfExperiments+1; i++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****BENCHMARK_ISSUANCE: \t experiment: ", i+1)

		issuer := issuer.NewIssuer(&conf)
		vc, _, _, _, _ := issuer.GenerateVC(nil)
		eddsaSignatureSize := len(vc.Proofs[0].ProofValue)

		seedSize := len(vc.Seed)
		eddsaSignatureTime := issuer.Result.EDDSASignTime

		result := ResultIssuance{
			EDDSASignatureTime: eddsaSignatureTime,
			EDDSASignatureSize: eddsaSignatureSize,
			SeedSize:           seedSize,
		}
		if i > 0 {
			zap.S().Infoln(result)
			finalResults = append(finalResults, result)
			WriteResultIssuanceToFile(result, false)
		}

	}

	end := time.Since(start)

	avgResult := ComputeAverageResultIssuance(finalResults)
	WriteResultIssuanceToFile(avgResult, true)
	zap.S().Infoln("Benchark_issuance took: ", end.Seconds(), " seconds")
}

/*
Benchmark_Revocation benchmarks

	a) the size of revoked list of tokens
	b) the time to refresh the list
	c) the cost of publishing the revoked tokens in the smart contract

The number of revoked VCs per epoch is decided following a uniform distribution.
 1. number_of_revoked_vcs := revocation_rate * total_vcs / 100
 2. number_of_revoked_vcs_per_epoch := number_of_revoked_vcs / m
*/
func Benchmark_Revocation(conf config.Config) {

	start := time.Now()
	results := ResultRevocationList{}
	results.Results = make([]ResultRevocation, 0)
	conf.Params.NumberOfExperiments = conf.Params.NumberOfExperiments + 1

	//revocationRateEnd := conf.Params.RevocationRateEnd

	// this value is set manually for benchmarking purposes
	revocationRateEnd := 5
	var wg sync.WaitGroup
	TotalRuns := (conf.Params.NumberOfExperiments + 1) * (((revocationRateEnd - conf.Params.RevocationRateBase) / conf.Params.RevocationRateStep) + 1) * conf.Params.ExpirationPeriod

	var opsCounter atomic.Uint64

	conf.InitialTimestamp = time.Now()
	issuer2 := issuer.NewIssuer(&conf)
	issuer2.GenerateAndStoreVCsWithoutSignatures(conf)
	totalVCs := conf.Params.TotalVCs
	ResetRevocationFiles()
	for exp := 0; exp < conf.Params.NumberOfExperiments; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****zkRevoke: BENCHMARK_REVOCATION: \t experiment: ", exp+1, "******")

		for revocationRate := conf.Params.RevocationRateBase; revocationRate <= revocationRateEnd; revocationRate = revocationRate + conf.Params.RevocationRateStep {
			wg.Add(1)
			go func(totalVCs int, revocationRate int, issuer2 *issuer.Issuer) {
				defer wg.Done()
				// The number of revoked VCs per epoch is decided following a uniform distribution.
				//    1) number_of_revoked_vcs := revocation_rate * total_vcs / 100
				//	  2) number_of_revoked_vcs_per_epoch := number_of_revoked_vcs / m
				revokedVCs := revocationRate * totalVCs / 100

				total_epochs := conf.Params.ExpirationPeriod

				// computes the number of revoked VCs per epoch
				revokedVCsPerEpoch := revokedVCs / total_epochs

				// revoke atleast one VC per epoch for the purpose of untraceability
				if revokedVCsPerEpoch == 0 {
					revokedVCsPerEpoch = 1
				}

				issuer2.ResetRevocationStorage()

				for x := 0; x < total_epochs; x++ {

					revocationTime := issuer2.RevokeVCsRandomly(revokedVCsPerEpoch)

					totalNumberOfRevokedVCs := revokedVCsPerEpoch * (x + 1)
					tokenSize := issuer2.GetTokenSize()
					opsCounter.Add(1)
					zap.S().Infoln(opsCounter.Load(), "/", TotalRuns, ":",
						"exp : ", exp+1,
						"Total VCs: ", totalVCs,
						"  total_epochs:  ", total_epochs,
						"  rev. rate: ", revocationRate, "%",
						"  rev. VCs per epoch: ", revokedVCsPerEpoch,
						"  cur. Epoch: ", x+1,
						"  issuer bandwidth: ", totalNumberOfRevokedVCs*tokenSize,
						"  revocation time: ", revocationTime, " nano seconds")
					if exp > 0 {
						resultRevocation := ResultRevocation{
							TotalRevokedVCsPerEpoch: revokedVCsPerEpoch,
							TimeToRevokeVCs:         revocationTime,
							TotalValidVCs:           totalVCs,
							RevocationRate:          revocationRate,
							IssuerBandwidth:         totalNumberOfRevokedVCs * tokenSize,
							TotalEpochs:             total_epochs,
							CurrentEpoch:            x + 1,
							TokenSize:               tokenSize,
						}
						results.Add(resultRevocation)
					}
				}
			}(totalVCs, revocationRate, issuer2)
		}

		wg.Wait()
	}

	end := time.Since(start)
	avgResults := ComputeAverageResultRevocation(results.Results)
	for _, result := range avgResults {
		WriteResultRevocationToFile(*result, true)
	}
	zap.S().Infoln("Benchmark_Refresh took: ", end.Minutes(), " minutes")
}

/*
Benchmark_Refresh benchmarks
a) the size of revoked list of tokens
b) the time to refresh the list
c) the cost of publishing the revoked tokens in the smart contract
*/
func Benchmark_Refresh(conf config.Config) {
	start := time.Now()
	var results []ResultRefresh
	ResetRefreshFiles()
	var opsCounter atomic.Uint64
	revocationRateEnd := 1
	TotalRuns := (conf.Params.NumberOfExperiments + 1) * conf.Params.ExpirationPeriod

	conf.InitialTimestamp = time.Now()
	issuer := issuer.NewIssuer(&conf)
	issuer.GenerateAndStoreVCsWithoutSignatures(conf)

	for exp := 0; exp < conf.Params.NumberOfExperiments+1; exp++ {

		for revocationRate := conf.Params.RevocationRateBase; revocationRate <= revocationRateEnd; revocationRate = revocationRate + conf.Params.RevocationRateStep {

			totalVCs := conf.Params.TotalVCs
			revokedVCs := revocationRate * totalVCs / 100
			total_epochs := conf.Params.ExpirationPeriod
			// computes the number of revoked VCs per epoch
			revokedVCsPerEpoch := revokedVCs / total_epochs

			// revoke atleast one VC per epoch for the purpose of untraceability
			if revokedVCsPerEpoch == 0 {
				revokedVCsPerEpoch = 1
			}

			address, _, _, _ := blockchain_hardhat.DeployContract(conf)
			conf.SmartContractAddress = address
			conf.InitialTimestamp = time.Now()

			issuer.ResetRevocationStorage()

			for x := 0; x < total_epochs; x++ {
				issuer.RevokeVCsRandomly(revokedVCsPerEpoch)
				endTokenGen := issuer.CalculateTimeToComputeTokensGivenEpoch(x + 1)

				res := ResultRefresh{
					NumberOfRevokedVCs: revokedVCsPerEpoch * (x + 1),
					RevocationRate:     revocationRate,
					Time:               int(endTokenGen.Microseconds()),
					TotalEpochs:        total_epochs,
					CurrentEpoch:       x + 1,
				}
				opsCounter.Add(1)
				zap.S().Infoln(opsCounter.Load(), "/", TotalRuns, ":",
					"exp : ", exp+1,
					"  rev. rate: ", revocationRate, "%",
					"  cur. Epoch: ", x+1,
					"  total_epochs:  ", total_epochs,
					"Total VCs: ", totalVCs,
					" refresh time: ", endTokenGen.Microseconds(), " micro seconds")

				if exp > 0 {
					WriteResultRefreshToFile(res, false)
					results = append(results, res)
				}
			}
		}

	}
	//avgResults := ComputeAverageResultRefresh(results)
	//for _, result := range avgResults {
	//	WriteResultRefreshToFile(*result, true)
	//}

	end := time.Since(start)
	zap.S().Infoln("Benchmark_Refresh took: ", end.Minutes(), " minutes")
}

/*
Benchmark_Verification bencharks the time to verify ZKP proof.
There are two variable in the verification procedure:
(a) k - number of tokens that can be proven using a single ZKP.
(b) m - number of epochs that a holder wants a verifier to check revocation status
(b)
*/
func Benchmark_Presentation_Verification(conf config.Config) {
	start := time.Now()

	MaxTokensInCircuit := conf.Params.NumberOfTokensPerCircuit
	MaxVPValidity := conf.Params.VerificationPeriod
	var results []ResultVerification
	var resultsPresentation []ResultPresentation
	TotalRuns := (conf.Params.NumberOfExperiments + 2) * MaxVPValidity * (MaxTokensInCircuit + 1)
	counter := 0
	ResetVerificationFiles()
	ResetPresentationFiles()
	conf.Params.EpochDuration = 100000 // epoch duration is set to large number since verifiers need sufficient time to verify proofs due to huge volume
	for exp := 0; exp < conf.Params.NumberOfExperiments+2; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****ZKREVOKE: BENCHMARK_PRESENTATION_VERIFICATION: \t experiment: ", exp+1, "******")
		for numberOfTokensInCircuit := 1; numberOfTokensInCircuit <= int(math.Pow(2, float64(MaxTokensInCircuit))); numberOfTokensInCircuit = numberOfTokensInCircuit * 2 {
			conf.Params.NumberOfTokensPerCircuit = numberOfTokensInCircuit
			issuer := issuer.NewIssuer(&conf)
			holder := holder.NewHolder(numberOfTokensInCircuit)
			vc, pki, duration, initialTimeStamp, _ := issuer.GenerateVC(holder.Holder_PublicKey.Bytes())

			verifier := verifier.NewVerifier()
			verifier.SetUpBlockchainConnection(conf)
			verifier.SetCCS(pki.Ccs)
			verifier.SetEddsaPublicKey(pki.EddsaPublicKey)
			verifier.SetZKPVerifyingKey(pki.ZkpVerifyingKey)

			verifier.Duration = int(duration)
			verifier.InitialTimeStamp = initialTimeStamp
			verifier.SelectiveDisclosureExtension = conf.SelectiveDisclosureExtension
			holder.InitCryptoKeys(pki)

			holder.SetDuration(duration)
			holder.SetInitialTimeStamp(initialTimeStamp)
			holder.SetNumberOfTokensInCircuit(numberOfTokensInCircuit)
			holder.SelectiveDisclosureExtension = conf.SelectiveDisclosureExtension
			// i-th holder is issued i-th vc
			holder.ReceiveVC(*vc)

			for m := 1; m <= MaxVPValidity; m = m + 1 {

				challenge := rand.Text()
				vp := holder.GenerateVP(m, []byte(challenge))

				numberOfZKPProofs := len(vp.TokenPresentations)
				zkpProofSize := len(vp.TokenPresentations[0].ZKPProof.ProofValue)
				totalZKPProofSize := 0
				for i := 0; i < numberOfZKPProofs; i++ {
					totalZKPProofSize += len(vp.TokenPresentations[i].ZKPProof.ProofValue)
				}

				verifier.InitialTimeStamp = conf.InitialTimestamp
				verifier.VerifyVP(challenge, *vp)

				counter = counter + 1

				zap.S().Infoln(counter, "/", TotalRuns, ":",
					"exp : ", exp+1,
					" VP Validity:  ", m,
					" Number of tokens in a circuit: ", numberOfTokensInCircuit,
					" ZKP proof ver. time: ", verifier.Result.VPMetrics[m-1].GrothProofVerificationTimeTotal, "micro secs.",
				)

				if exp > 1 {
					res := ResultVerification{
						NumberOfTokensInCircuit: numberOfTokensInCircuit,
						VPValidityPeriod:        m,
						ZKPProofVerTime:         verifier.Result.VPMetrics[m-1].GrothProofVerificationTimeTotal,
					}
					results = append(results, res)

					resPresentation := ResultPresentation{
						VPValidityPeriod:           m,
						NumberOfTokensInCircuit:    numberOfTokensInCircuit,
						NumberOfZKPProofs:          numberOfZKPProofs,
						ZKPProofSize:               zkpProofSize,
						TotalZKPProofSize:          totalZKPProofSize,
						TimeToGenerateOneZKPProof:  holder.Result.VPSizeMetrics[m-1].SingleZKPProofGenTime,
						TimeToGenerateAllZKPProofs: holder.Result.VPSizeMetrics[m-1].AllZKPProofsGenTime,
					}

					resultsPresentation = append(resultsPresentation, resPresentation)
				}

			}

		}
	}
	//avgResults := ComputeAverageResultVerification(results)
	for _, result := range results {
		WriteResultVerificationToFile(result, false)
	}

	//avgResultsPresentation := ComputeAverageResultPresentation(resultsPresentation)
	for _, result := range resultsPresentation {
		WriteResultPresentationToFile(result, false)
	}
	end := time.Since(start)
	zap.S().Infoln("Benchmark_Verification took: ", end.Minutes(), " minutes")
}

func Benchmark_TokenStorageCost(conf config.Config) {
	start := time.Now()
	var results []ResultRefresh
	ResetTokenStorageResultFiles()
	for exp := 0; exp < conf.Params.NumberOfExperiments+1; exp++ {
		index := 0
		count := len(conf.TokenStorageNumberOfTokens)
		for i := 0; i < count; i++ {
			conf.Params.TotalVCs = conf.TokenStorageNumberOfTokens[index]
			numberOfTokensInSmartContractPerEpoch := conf.TokenStorageNumberOfTokens[index]
			epoch := exp*count + i
			address, _, _, _ := blockchain_hardhat.DeployContract(conf)
			conf.SmartContractAddress = address
			conf.InitialTimestamp = time.Now()

			issuer := issuer.NewIssuer(&conf)
			issuer.SetUpBlockchainConnection(conf)

			gasUsed, cost, tokenSize, txCount, time := issuer.ComputeTokenStorageCost(numberOfTokensInSmartContractPerEpoch, epoch)
			res := ResultRefresh{
				NumberOfRevokedVCs:   int(numberOfTokensInSmartContractPerEpoch),
				TokenSize:            tokenSize,
				Gas:                  gasUsed,
				Cost:                 cost,
				NumberOfTransactions: txCount,
				Time:                 int(time.Microseconds()),
			}
			index++
			if exp > 0 {
				WriteTokenStorageResultToFile(res, false)
				results = append(results, res)
			}
		}
	}
	avgResults := ComputeAverageResultRefresh(results)
	for _, result := range avgResults {
		WriteTokenStorageResultToFile(*result, true)
	}

	end := time.Since(start)
	zap.S().Infoln("Benchmark_Benchmark_TokenStorageCost took: ", end.Minutes(), " minutes")
}

/*
Benchmarks the performance and cost of invidual constraints in the ZKP circuit
*/
func Benchmark_CircuitConstraints(conf config.Config) {
	start := time.Now()
	var results []ResultCircuitConstraints
	ResetResultCircuitConstraints()
	for exp := 0; exp < conf.Params.NumberOfExperiments+1; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****ZKREVOKE: BENCHMARK_CIRCUIT_CONSTRAINTS: \t experiment: ", exp+1, "******")

		res := BenchmarkEmptyCircuit()
		res1 := BenchmarkTokenVerificationCircuit()
		res2 := BenchmarkChallengeVerification()
		res3 := BenchmarkSignatureVerificationCircuit()
		res4 := BenchmarkCompleteCircuit()

		if exp > 0 {
			WriteResultCircuitConstraintsToFile(*res, false)
			results = append(results, *res)
			WriteResultCircuitConstraintsToFile(*res1, false)
			results = append(results, *res1)
			WriteResultCircuitConstraintsToFile(*res2, false)
			results = append(results, *res2)
			WriteResultCircuitConstraintsToFile(*res3, false)
			results = append(results, *res3)
			WriteResultCircuitConstraintsToFile(*res4, false)
			results = append(results, *res4)
		}
	}
	avgResults := ComputeAverageResultCircuitConstraints(results)
	for _, result := range avgResults {
		WriteResultCircuitConstraintsToFile(*result, true)
	}

	end := time.Since(start)
	zap.S().Infoln("Benchmark circuit constraints took: ", end.Minutes(), " minutes")
}

func Benchmark_ListCommitment(conf config.Config) {
	results := ResultListCommitmentList{}
	results.Results = make([]ResultListCommitment, 0)
	ResetCommitmentFiles()
	benchmark_start := time.Now()
	for exp := 0; exp < conf.Params.NumberOfExperiments+1; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****ZKREVOKE: BENCHMARK_LIST_COMMITMENT: \t experiment: ", exp+1, "******")
		for revocationRate := conf.Params.RevocationRateBase; revocationRate <= conf.Params.RevocationRateEnd; revocationRate = revocationRate + conf.Params.RevocationRateStep {

			totalVCs := conf.Params.TotalVCs
			revokedVCs := revocationRate * totalVCs / 100
			total_epochs := conf.Params.ExpirationPeriod
			// computes the number of revoked VCs per epoch
			revokedVCsPerEpoch := revokedVCs / total_epochs

			// revoke atleast one VC per epoch for the purpose of untraceability
			if revokedVCsPerEpoch == 0 {
				revokedVCsPerEpoch = 1
			}
			var tokens [][]byte
			for x := 0; x < total_epochs; x++ {

				// compute a list of tokens
				for j := 0; j < revokedVCsPerEpoch; j++ {
					seed := rand.Text()
					epoch := x
					token := utils.ComputeToken(epoch, seed)
					tokens = append(tokens, token)
				}

				h := sha256.New()

				start := time.Now()
				// compute a list commitment
				for k := 0; k < len(tokens); k++ {
					h.Write([]byte(tokens[k]))
				}
				list_commitment := h.Sum(nil)
				commtiment_create_time := time.Since(start)

				h1 := sha256.New()
				start2 := time.Now()
				// compute a list commitment
				for k := 0; k < len(tokens); k++ {
					h1.Write([]byte(tokens[k]))
				}
				list_commitment2 := h1.Sum(nil)

				status := bytes.Equal(list_commitment, list_commitment2)
				commtiment_verification_time := time.Since(start2)

				if exp > 0 {
					resultListCommitment := ResultListCommitment{
						TotalRevokedVCsPerEpoch:        revokedVCsPerEpoch,
						TotalValidVCs:                  totalVCs,
						RevocationRate:                 revocationRate,
						TotalEpochs:                    total_epochs,
						CurrentEpoch:                   x + 1,
						TimeToCreateCommitment:         int(commtiment_create_time.Microseconds()),
						TimeToVerifyCommitment:         int(commtiment_verification_time.Microseconds()),
						SizeOfTheListAtTheCurrentEpoch: len(tokens),
					}
					results.Add(resultListCommitment)
					if status {
						zap.S().Infoln("Exp: ", exp, "\t Epoch: ", x+1, "Revocation rate: ", revocationRate, "% \t List size: ", len(tokens), "\t time to create (micro seconds):", commtiment_create_time.Microseconds(), "\t time to verify (micro seconds):", commtiment_verification_time.Microseconds())
					} else {
						zap.S().Infoln("Commitment list is not ok")
					}
				}
			}
		}

	}
	avgResults := ComputeAverageResultListCommitment(results.Results)
	for _, result := range avgResults {
		WriteResultListCommitmentToFile(*result, true)
	}

	benchmark_end := time.Since(benchmark_start)
	zap.S().Infoln("Benchmark list commitment took: ", benchmark_end.Minutes(), " minutes")

}
