package irma

import (
	"crypto/rand"
	"encoding/json"
	"github.com/privacybydesign/gabi"
	"github.com/privacybydesign/gabi/gabikeys"
	"github.com/privacybydesign/gabi/revocation"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
	"zkrevoke/config"
	"zkrevoke/irma/benchmark"
	"zkrevoke/irma/internal/common"
)

func Benchmark_Setup(conf config.Config) {

	Logger := logrus.StandardLogger()
	revocation.Logger = Logger
	var results []benchmark.ResultSetup
	conf.Params.NumberOfExperiments = conf.Params.NumberOfExperiments + 1
	benchmark.ResetSetupFiles()

	for exp := 0; exp < conf.Params.NumberOfExperiments; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****IRMA: BENCHMARK_SETUP: \t experiment: ", exp+1, "******")
		startKeyGen := time.Now()
		sk, pk := Generate_keys()
		_ = gabikeys.GenerateRevocationKeypair(sk, pk)
		zap.S().Infoln("secret key: ", sk)
		endKeyGen := time.Since(startKeyGen)

		startAccGen := time.Now()
		update, _ := GenerateNewAccumulator(sk, pk)
		endAccGen := time.Since(startAccGen)

		acc := update.SignedAccumulator.Accumulator
		zap.S().Infoln("***IRMA: BENCHMARK SETUP: initial accumulator value:", acc.Nu,
			"\t size in Bytes: ", len(acc.Nu.Bytes()))
		zap.S().Infoln("***IRMA: BENCHMARK SETUP: accumulator events: factor: ", update.Events[0].E, "\t size: ", len(update.Events))

		skBytes, _ := json.Marshal(sk)
		pkBytes, _ := json.Marshal(pk)
		result := benchmark.ResultSetup{
			PrivateKeySize:     len(skBytes),
			PublicKeySize:      len(pkBytes),
			KeyGenTime:         int(endKeyGen.Microseconds()),
			AccumulatorGenTime: int(endAccGen.Microseconds()),
			AccumulatorSize:    len(acc.Nu.Bytes()),
		}
		if exp > 0 {
			results = append(results, result)
			benchmark.WriteResultSetupToFile(result, false)
		}
	}
	avgResult := benchmark.ComputeAverageResultSetup(results)
	benchmark.WriteResultSetupToFile(avgResult, true)

}

func Benchmark_Issuance(conf config.Config) {
	Logger := logrus.StandardLogger()
	revocation.Logger = Logger

	var results []benchmark.ResultIssuance
	conf.Params.NumberOfExperiments = conf.Params.NumberOfExperiments + 1
	numberOfVCs := conf.Params.TotalVCs
	benchmark.ResetIssuanceFiles()
	for exp := 0; exp < conf.Params.NumberOfExperiments; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****IRMA: BENCHMARK_ISSUANCE: \t experiment: ", exp+1, "******")

		sk, pk := Generate_keys()
		err := gabikeys.GenerateRevocationKeypair(sk, pk)
		if err != nil {
			zap.S().Error("Error generating revocation keypair", err)
		}

		context, err := common.RandomBigInt(pk.Params.Lh)
		if err != nil {
			zap.S().Errorln("***IRMA***: ", err)
		}

		update, _ := GenerateNewAccumulator(sk, pk)
		acc := update.SignedAccumulator.Accumulator

		avgTimeToGenWitness := 0
		totalTimeToGenWitness := 0

		avgTimeToGenVC := 0
		totalTimeToGenVC := 0
		for i := 0; i < numberOfVCs; i++ {

			startWitnessGen := time.Now()
			factor1, _ := common.RandomPrimeInRange(rand.Reader, 3, revocation.Parameters.AttributeSize)
			witness, _ := GenerateWitnessForFactor(sk, acc, factor1)
			witness.SignedAccumulator = update.SignedAccumulator
			EndWitnessGen := time.Since(startWitnessGen)
			totalTimeToGenWitness += int(EndWitnessGen.Microseconds())

			attrs := revocationAttrs(witness)
			startVCGen := time.Now()
			nonce1, err := common.RandomBigInt(pk.Params.Lstatzk)

			nonce2, err := common.RandomBigInt(pk.Params.Lstatzk)

			secret, err := common.RandomBigInt(pk.Params.Lm)

			b, err := gabi.NewCredentialBuilder(pk, context, secret, nonce2, nil, nil)
			if err != nil {
				zap.S().Errorln("***IRMA***: error creating a new credential builder", err)
			}
			commitMsg, err := b.CommitToSecretAndProve(nonce1)
			issuer := gabi.NewIssuer(sk, pk, context)
			msg, err := issuer.IssueSignature(commitMsg.U, attrs, witness, nonce2, nil)
			if err != nil {
				zap.S().Errorln("***IRMA***: error issuing signature: ", err)
			}
			b.ConstructCredential(msg, attrs)
			endVCGen := time.Since(startVCGen)
			totalTimeToGenVC += int(endVCGen.Microseconds())
		}
		avgTimeToGenWitness += int(totalTimeToGenWitness) / numberOfVCs
		avgTimeToGenVC += int(totalTimeToGenVC) / numberOfVCs
		result := benchmark.ResultIssuance{
			TotalVCs:                 numberOfVCs,
			TimeToGenerateOneVC:      avgTimeToGenVC,
			TimeToGenerateAllVCs:     totalTimeToGenVC,
			TimeToGenerateOneWitness: avgTimeToGenWitness,
			TimeToGenerateAllWitness: totalTimeToGenWitness,
		}
		zap.S().Infoln(result)
		if exp > 0 {
			results = append(results, result)
			benchmark.WriteResultIssuanceToFile(result, false)
		}

	}
	avgResult := benchmark.ComputeAverageResultIssuance(results)
	for _, result := range avgResult {
		benchmark.WriteResultIssuanceToFile(*result, true)
	}

}

func Benchmark_Revocation(conf config.Config) {
	Logger := logrus.StandardLogger()
	revocation.Logger = Logger

	resultsRevocation := benchmark.ResultRevocationList{}
	resultsRevocation.Results = make([]benchmark.ResultRevocation, 0)

	resultsPresentation := benchmark.ResultPresentationAndVerificationList{}
	resultsPresentation.Results = make([]benchmark.ResultPresentationAndVerification, 0)
	conf.Params.NumberOfExperiments = conf.Params.NumberOfExperiments + 1
	totalVCs := conf.Params.TotalVCs
	var wg sync.WaitGroup
	benchmark.ResetRevocationFiles(conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)

	//revocationRateEnd := conf.Params.RevocationRateEnd

	// this value is set manually for benchmarking purposes
	revocationRateEnd := 5
	TotalRuns := (conf.Params.NumberOfExperiments) * (((revocationRateEnd - conf.Params.RevocationRateBase) / conf.Params.RevocationRateStep) + 1) * conf.Params.ExpirationPeriod
	var opsCounter atomic.Uint64

	var witnessUpdateWitnessRepetitionLabel string
	if conf.Params.IRMAWitnessUpdateMessageWithoutRepetition == false {
		witnessUpdateWitnessRepetitionLabel = "-no repetition"
	}

	for exp := 0; exp < conf.Params.NumberOfExperiments; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****IRMA: BENCHMARK_REVOCATION", witnessUpdateWitnessRepetitionLabel, ": \t experiment: ", exp+1, "******")

		for revocationRate := conf.Params.RevocationRateBase; revocationRate <= revocationRateEnd; revocationRate = revocationRate + conf.Params.RevocationRateStep {

			wg.Add(1)
			go func(totalVCs int, revocationRate int) {
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
				sk, pk := Generate_keys()
				err := gabikeys.GenerateRevocationKeypair(sk, pk)
				if err != nil {
					zap.S().Error("Error generating revocation keypair", err)
				}

				update, _ := GenerateNewAccumulator(sk, pk)
				acc := update.SignedAccumulator.Accumulator
				events := update.Events
				event := events[0]

				context, err := common.RandomBigInt(pk.Params.Lh)
				nonce1, err := common.RandomBigInt(pk.Params.Lstatzk)
				nonce2, err := common.RandomBigInt(pk.Params.Lstatzk)
				secret, err := common.RandomBigInt(pk.Params.Lm)
				b, err := gabi.NewCredentialBuilder(pk, context, secret, nonce2, nil, nil)
				if err != nil {
					zap.S().Errorln("***IRMA***: error creating a new credential builder", err)
				}
				commitMsg, err := b.CommitToSecretAndProve(nonce1)

				issuer := gabi.NewIssuer(sk, pk, context)

				factor1, _ := common.RandomPrimeInRange(rand.Reader, 3, revocation.Parameters.AttributeSize)
				witness, err := GenerateWitnessForFactor(sk, acc, factor1)
				witness.SignedAccumulator = update.SignedAccumulator

				if err != nil {
					zap.S().Errorln("***IRMA***: error creating a new witness", err)
				}
				attrs := revocationAttrs(witness)

				msg, err := issuer.IssueSignature(commitMsg.U, attrs, witness, nonce2, nil)
				if err != nil {
					zap.S().Errorln("***IRMA***: error issuing signature: ", err)
				}

				_, err = b.ConstructCredential(msg, attrs)
				if err != nil {
					zap.S().Errorln("***IRMA***: error constructing credential: ", err)
				}

				// show again, using the nonrevocation proof cache
				// nonce1s is provided by a verifier to a holder to prevent reply attacks (verifier verifies that the holder is not
				//replying older proofs)

				revoked_vcs_map := make(map[string]bool)

				// after each epoch, an issuer sends witness update message to all holders with valid VCs
				// during each epoch, an issuer revokes VCs and the accumulator is updated after each epoch
				zap.S().Infoln("*******************************************************************************************************************************************")
				for x := 0; x < total_epochs; x++ {

					_, err = common.RandomBigInt(pk.Params.Lstatzk)
					revTime := 0

					for j := 0; j < revokedVCsPerEpoch; j++ {

						index, _ := common.RandomPrimeInRange(rand.Reader, 3, revocation.Parameters.AttributeSize)
						if revoked_vcs_map[index.String()] == true {
							j = j - 1
							continue
						}
						revoked_vcs_map[index.String()] = true
						startRevocation := time.Now()
						w, _ := GenerateWitnessForFactor(sk, acc, index)
						acc, event, err = acc.Remove(sk, w.E, event)
						endRevocation := time.Since(startRevocation)
						revTime = revTime + int(endRevocation.Microseconds())
						if err != nil {
							zap.S().Errorln("***IRMA***: error removing witness", err)
						}
						events = append(events, event)
					}

					startUpdate := time.Now()
					update, err = revocation.NewUpdate(sk, acc, events)
					if conf.Params.IRMAWitnessUpdateMessageWithoutRepetition == true {
						events = make([]*revocation.Event, 0)
					}
					endUpdate := time.Since(startUpdate)
					revTime = revTime + int(endUpdate.Microseconds())
					if err != nil {
						zap.S().Errorln("Error creating revocation update", err)
					}

					witnessUpdateSize, _ := update.MarshalJSON()
					numberOfValidHolders := totalVCs - ((x + 1) * revokedVCsPerEpoch)
					issuerBandwidth := numberOfValidHolders * len(witnessUpdateSize)

					if exp > 0 {
						resultRevocation := benchmark.ResultRevocation{
							TotalValidVCs:            totalVCs,
							RevocationRate:           revocationRate,
							IssuerBandwidth:          issuerBandwidth,
							TotalRevokedVCsPerEpoch:  revokedVCsPerEpoch,
							TotalHoldersWithValidVCs: numberOfValidHolders,
							TotalEpochs:              total_epochs,
							CurrentEpoch:             x + 1,
							TimeToRevokeVCs:          revTime,
							WitnessUpdateSize:        len(witnessUpdateSize),
						}
						resultsRevocation.Add(resultRevocation, conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)
					}

					opsCounter.Add(1)
					zap.S().Infoln(opsCounter.Load(), "/", TotalRuns, ":",
						"exp : ", exp+1,
						"Total VCs: ", totalVCs,
						"  total epochs:  ", total_epochs,
						"  rev. rate: ", revocationRate, "%",
						"  rev. VCs per epoch: ", revokedVCsPerEpoch,
						"  cur. Epoch: ", x+1,
						"  issuer bandwidth: ", issuerBandwidth,
						"  witness update size: ", len(witnessUpdateSize), " B",
						"  revocation time: ", revTime, " micro seconds")

				}
			}(totalVCs, revocationRate)
		}

		wg.Wait()
	}

	avgResultRevocation := benchmark.ComputeAverageResultRevocation(resultsRevocation.Results, conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)
	for _, result := range avgResultRevocation {
		benchmark.WriteResultRevocationToFile(*result, true, conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)
	}

}

func Benchmark_Presentation_And_Verification(conf config.Config) {
	Logger := logrus.StandardLogger()
	revocation.Logger = Logger
	conf.Params.NumberOfExperiments = conf.Params.NumberOfExperiments + 1

	results := benchmark.ResultPresentationAndVerificationList{}
	results.Results = make([]benchmark.ResultPresentationAndVerification, 0)
	totalVCs := conf.Params.TotalVCs
	var wg sync.WaitGroup

	revocationRateEnd := conf.Params.RevocationRateEnd

	TotalRuns := (conf.Params.NumberOfExperiments) * (((revocationRateEnd - conf.Params.RevocationRateBase) / conf.Params.RevocationRateStep) + 1) * conf.Params.VerificationPeriod

	var opsCounter atomic.Uint64
	benchmark.ResetPresentationAndVerificationFiles(conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)
	var witnessUpdateWitnessRepetitionLabel string
	if conf.Params.IRMAWitnessUpdateMessageWithoutRepetition == false {
		witnessUpdateWitnessRepetitionLabel = "-no repetition"
	}
	for exp := 0; exp < conf.Params.NumberOfExperiments; exp++ {
		zap.S().Infoln("\n \n *******************************************************************************************")
		zap.S().Infoln("\n  *****IRMA: BENCHMARK_PRESENTATION_AND_VERIFICATION", witnessUpdateWitnessRepetitionLabel, ": \t experiment: ", exp+1, "******")

		for revocationRate := conf.Params.RevocationRateBase; revocationRate <= revocationRateEnd; revocationRate = revocationRate + conf.Params.RevocationRateStep {

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

			sk, pk := Generate_keys()
			err := gabikeys.GenerateRevocationKeypair(sk, pk)
			if err != nil {
				zap.S().Error("Error generating revocation keypair", err)
			}

			update, _ := GenerateNewAccumulator(sk, pk)
			acc := update.SignedAccumulator.Accumulator
			events := update.Events
			event := events[0]

			context, err := common.RandomBigInt(pk.Params.Lh)
			nonce1, err := common.RandomBigInt(pk.Params.Lstatzk)
			nonce2, err := common.RandomBigInt(pk.Params.Lstatzk)
			secret, err := common.RandomBigInt(pk.Params.Lm)
			b, err := gabi.NewCredentialBuilder(pk, context, secret, nonce2, nil, nil)
			if err != nil {
				zap.S().Errorln("***IRMA***: error creating a new credential builder", err)
			}
			commitMsg, err := b.CommitToSecretAndProve(nonce1)

			issuer := gabi.NewIssuer(sk, pk, context)

			factor1, _ := common.RandomPrimeInRange(rand.Reader, 3, revocation.Parameters.AttributeSize)
			witness, err := GenerateWitnessForFactor(sk, acc, factor1)
			witness.SignedAccumulator = update.SignedAccumulator

			if err != nil {
				zap.S().Errorln("***IRMA***: error creating a new witness", err)
			}
			attrs := revocationAttrs(witness)

			msg, err := issuer.IssueSignature(commitMsg.U, attrs, witness, nonce2, nil)
			if err != nil {
				zap.S().Errorln("***IRMA***: error issuing signature: ", err)
			}

			cred, err := b.ConstructCredential(msg, attrs)
			if err != nil {
				zap.S().Errorln("***IRMA***: error constructing credential: ", err)
			}

			// show again, using the nonrevocation proof cache
			// nonce1s is provided by a verifier to a holder to prevent reply attacks (verifier verifies that the holder is not
			//replying older proofs)

			timeToUpdateWitnesses := 0
			sizeOfTotalWitnessUpdateMessagesReceived := 0
			timeToCreateDisclosureProof := 0
			timeToVerifyProof := 0
			totalDisclosureProofSize := 0
			totalNonRevProofSize := 0
			revoked_vcs_map := make(map[string]bool)

			// after each epoch, a holder receives a witness update message from an issuer
			// during each epoch, an issuer revokes VCs and the accumulator is updated after each epoch
			zap.S().Infoln("*******************************************************************************************************************************************")
			for x := 0; x < conf.Params.VerificationPeriod; x++ {

				nonce1s, err := common.RandomBigInt(pk.Params.Lstatzk)

				for j := 0; j < revokedVCsPerEpoch; j++ {

					index, _ := common.RandomPrimeInRange(rand.Reader, 3, revocation.Parameters.AttributeSize)
					if revoked_vcs_map[index.String()] == true {
						j = j - 1
						continue
					}
					revoked_vcs_map[index.String()] = true
					w, _ := GenerateWitnessForFactor(sk, acc, index)
					acc, event, err = acc.Remove(sk, w.E, event)
					if err != nil {
						zap.S().Errorln("***IRMA***: error removing witness", err)
					}
					events = append(events, event)
				}

				update, err = revocation.NewUpdate(sk, acc, events)
				if conf.Params.IRMAWitnessUpdateMessageWithoutRepetition == true {
					events = make([]*revocation.Event, 0)
				}

				if err != nil {
					zap.S().Errorln("Error creating revocation update", err)
				}

				//********************* holder receiving witness update from the issuer and updating the local witness ************************
				witnessUpdateMessageSize, _ := update.MarshalJSON()
				sizeOfTotalWitnessUpdateMessagesReceived = sizeOfTotalWitnessUpdateMessagesReceived + len(witnessUpdateMessageSize)
				WitnessUpdateStart := time.Now()
				err = cred.NonRevocationWitness.Update(pk, update)
				if err != nil {
					zap.S().Errorln("***IRMA***: error updating non revocation witness", err)
				}
				cred.NonrevPrepareCache()
				WitnessUpdateEnd := time.Since(WitnessUpdateStart)
				timeToUpdateWitnesses = timeToUpdateWitnesses + int(WitnessUpdateEnd.Microseconds())

				ProofWithNonRevocationStart := time.Now()
				proofd, err := cred.CreateDisclosureProof(nil, nil, true, context, nonce1s)
				if err != nil {
					zap.S().Errorln("Error creating disclosure proof", err)
				}
				ProofWithNonRevocationEnd := time.Since(ProofWithNonRevocationStart)
				timeToCreateDisclosureProof = timeToCreateDisclosureProof + int(ProofWithNonRevocationEnd.Microseconds())

				proofDJson, _ := json.Marshal(proofd)
				totalDisclosureProofSize = totalDisclosureProofSize + len(proofDJson)

				nonRevProofJson, _ := json.Marshal(proofd.NonRevocationProof)
				totalNonRevProofSize = totalNonRevProofSize + len(nonRevProofJson)

				proofVerificationStart := time.Now()
				status := proofd.Verify(pk, context, nonce1s, false)
				proofVerificationEnd := time.Since(proofVerificationStart)
				timeToVerifyProof = timeToVerifyProof + int(proofVerificationEnd.Microseconds())

				if status == false {
					zap.S().Errorln("***IRMA***: proof verification failure: at epoch: ", x+1)
				}
				opsCounter.Add(1)

				zap.S().Infoln(opsCounter.Load(), "/", TotalRuns, ":",
					"exp : ", exp+1,
					"Total VCs: ", totalVCs,
					" total epochs:  ", total_epochs,
					" max VP validity: ", conf.Params.VerificationPeriod,
					" rev. rate: ", revocationRate,
					" revoked VCs per epoch: ", revokedVCsPerEpoch,
					" cur. epoch: ", x+1,
					" non. rev. proof size: ", len(nonRevProofJson), "B",
					"wit. update message size: ", len(witnessUpdateMessageSize), "B",
					" holder bandwidth: ", sizeOfTotalWitnessUpdateMessagesReceived+totalNonRevProofSize, "B",
					" total wit. messages: ", sizeOfTotalWitnessUpdateMessagesReceived,
					" total proofs sent:", totalNonRevProofSize,
					" total full proofs sent: ", totalDisclosureProofSize,
					" wit. update time: ", WitnessUpdateEnd.Microseconds(), "micro secs.",
					"total wit. update time: ", timeToUpdateWitnesses, "micro secs.",
					" proof ver. time: ", proofVerificationEnd.Microseconds(), "micro secs.")

				if exp > 0 {
					result := benchmark.ResultPresentationAndVerification{
						TotalNumberOfIssuedVCs:  totalVCs,
						RevocationRate:          revocationRate,
						VPValidityPeriod:        conf.Params.VerificationPeriod,
						TotalNumberOfEpochs:     total_epochs,
						CurrentEpoch:            x + 1,
						TotalRevokedVCsPerEpoch: revokedVCsPerEpoch,
						TimeToUpdateWitness:     timeToUpdateWitnesses,
						TimeToCreateDisclosureProofWithNonRevocation: timeToCreateDisclosureProof,
						DisclosureProofSize:                          totalDisclosureProofSize,
						NonRevocationProofSize:                       totalNonRevProofSize,
						ProofVerificationTime:                        timeToVerifyProof,
						HolderBandwidth:                              sizeOfTotalWitnessUpdateMessagesReceived + totalNonRevProofSize,
						TotalWitnessUpdateMessagesReceived:           sizeOfTotalWitnessUpdateMessagesReceived,
					}
					results.Add(result, conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)
				}
			}

		}

		wg.Wait()
	}

	//avgResult := benchmark.ComputeAverageResultPresentationAndVerification(results.Results)
	for _, result := range results.Results {
		benchmark.WriteResultPresentationAndVerificationToFile(result, false, conf.Params.IRMAWitnessUpdateMessageWithoutRepetition)
	}
}
