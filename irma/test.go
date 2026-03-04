package irma

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/privacybydesign/gabi"
	"github.com/privacybydesign/gabi/gabikeys"
	"github.com/privacybydesign/gabi/revocation"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"time"
	"zkrevoke/irma/internal/common"
)

func TestFullIssueAndShowWithRevocation() {
	Logger := logrus.StandardLogger()
	revocation.Logger = Logger

	sk, pk := Generate_keys()

	err := gabikeys.GenerateRevocationKeypair(sk, pk)
	if err != nil {
		zap.S().Error("Error generating revocation keypair", err)
	}
	update, _ := GenerateNewAccumulator(sk, pk)
	acc := update.SignedAccumulator.Accumulator
	events := update.Events

	factor1, _ := common.RandomPrimeInRange(rand.Reader, 3, revocation.Parameters.AttributeSize)

	witness, _ := GenerateWitnessForFactor(sk, acc, factor1)
	witness.SignedAccumulator = update.SignedAccumulator
	status := VerifyWitness(witness, pk, update.SignedAccumulator)
	proofForFactor1 := GenerateProofForZKP(witness, pk)
	if proofForFactor1 == nil {
		fmt.Println("Factor1: ", factor1, " zkp proof generation failed")
	} else {
		statusOfFactor1 := VerifyZKP(proofForFactor1, pk, update.SignedAccumulator)
		fmt.Println("Factor1: ", factor1, "\t zkp verification status: ", statusOfFactor1)
	}
	// Issuance
	context, err := common.RandomBigInt(pk.Params.Lh)
	if err != nil {
		zap.S().Errorln("***IRMA***: ", err)
	}
	nonce1, err := common.RandomBigInt(pk.Params.Lstatzk)

	nonce2, err := common.RandomBigInt(pk.Params.Lstatzk)

	secret, err := common.RandomBigInt(pk.Params.Lm)

	b, err := gabi.NewCredentialBuilder(pk, context, secret, nonce2, nil, nil)
	if err != nil {
		zap.S().Errorln("***IRMA***: error creating a new credential builder", err)
	}
	commitMsg, err := b.CommitToSecretAndProve(nonce1)

	issuer := gabi.NewIssuer(sk, pk, context)
	attrs := revocationAttrs(witness)

	msg, err := issuer.IssueSignature(commitMsg.U, attrs, witness, nonce2, nil)
	if err != nil {
		zap.S().Errorln("***IRMA***: error issuing signature: ", err)
	}

	cred, err := b.ConstructCredential(msg, attrs)
	zap.S().Infoln("***IRMA***: issuing new VC: ", cred.Attributes)
	zap.S().Infoln("***IRMA***: issuing new VC acc. witness: ", cred.NonRevocationWitness.U)
	if err != nil {
		zap.S().Errorln("***IRMA***: error constructing credential: ", err)
	}
	// Showing

	disclosedAttributes := []int{1, 3}

	start := time.Now()
	// prepare nonrevocation proof cache
	cred.NonrevPrepareCache()
	// show again, using the nonrevocation proof cache
	// nonce1s is provided by a verifier to a holder to prevent reply attacks (verifier verifies that the holder is not
	//replying older proofs)
	nonce1s, err := common.RandomBigInt(pk.Params.Lstatzk)
	proofd, err := cred.CreateDisclosureProof(disclosedAttributes, nil, true, context, nonce1s)
	end := time.Now()
	proofDJson, _ := json.Marshal(proofd)
	nonRevProofJson, _ := json.Marshal(proofd.NonRevocationProof)
	responseSize := 0
	for key, _ := range proofd.NonRevocationProof.Responses {
		responseSize = responseSize + len(proofd.NonRevocationProof.Responses[key].Bytes())
	}
	zap.S().Infoln("Time to create a disclosure proof (in milliseconds): ", end.Sub(start).Milliseconds())
	zap.S().Infoln("proof size (bytes): ", len(proofDJson), "\t non revocation proof size (bytes): ", len(nonRevProofJson))
	zap.S().Infoln("Non revocation proof: ", "\n Nu: ", len(proofd.NonRevocationProof.Nu.Bytes()),
		"\n Cu: ", len(proofd.NonRevocationProof.Cu.Bytes()),
		"\n Cr: ", len(proofd.NonRevocationProof.Cr.Bytes()),
		"\n challenges: ", len(proofd.NonRevocationProof.Challenge.Bytes()),
		"\n responses: ", responseSize,
		"\n signed accumulator: (1831 - previous sizes)")
	start = time.Now()
	status = proofd.Verify(pk, context, nonce1s, false)
	zap.S().Infoln("verification result: ", status)
	end = time.Now()
	zap.S().Infoln("Time to verify disclosure proof (in milliseconds): ", end.Sub(start).Milliseconds())

	start0 := time.Now()
	// simulate revocation of another credential
	w, err := revocation.RandomWitness(sk, acc)
	acc, event, err := acc.Remove(sk, w.E, update.Events[len(update.Events)-1])
	events = append(events, event)
	end0 := time.Now()
	zap.S().Infoln("Revoked 1 random VCs.  Time for revocation (in seconds): ", end0.Sub(start0).Seconds())

	update, err = revocation.NewUpdate(sk, acc, events)
	zap.S().Infoln("Witness update message after revoking 1 VC: number of events: ", len(update.Events))
	start00 := time.Now()
	cred.NonRevocationWitness.Update(pk, update)
	end00 := time.Now()
	cred.NonrevPrepareCache()
	nonce1s, err = common.RandomBigInt(pk.Params.Lstatzk)
	start000 := time.Now()
	proofd, err = cred.CreateDisclosureProof(disclosedAttributes, nil, true, context, nonce1s)
	end000 := time.Now()
	proofDJson, _ = json.Marshal(proofd)
	nonRevProofJson, _ = json.Marshal(proofd.NonRevocationProof)
	start0000 := time.Now()
	status = proofd.Verify(pk, context, nonce1s, false)
	end0000 := time.Now()
	zap.S().Infoln(" Disclosure proof size (bytes): ", len(proofDJson),
		"\t non revocation proof size (bytes): ", len(nonRevProofJson),
		"\n Time to update witness: ", end00.Sub(start00).Milliseconds(), " ms",
		"\t Time to create disclosure proof: ", end000.Sub(start000).Milliseconds(), " ms",
		"\t Time to verify: ", end0000.Sub(start0000).Milliseconds(), " ms")

	start1 := time.Now()
	for i := 0; i < 100; i++ {
		var event1 []*revocation.Event
		w, err = revocation.RandomWitness(sk, acc)
		acc, event, err = acc.Remove(sk, w.E, event)
		events = append(events, event)
		event1 = append(event1, event)
		update, err = revocation.NewUpdate(sk, acc, event1)
		zap.S().Infoln("Witness update message after revoking a VCs: number of events: ", len(update.Events))
		cred.NonRevocationWitness.Update(pk, update)
	}
	end1 := time.Now()
	zap.S().Infoln("Revoked 100 random VCs.  Time for revocation (in seconds): ", end1.Sub(start1).Seconds())
	//update, err = revocation.NewUpdate(sk, acc, events)

	if err != nil {
		zap.S().Errorln("***IRMA***: error creating update", err)
	}
	start2 := time.Now()
	// update witness and nonrevocation proof cache
	//cred.NonRevocationWitness.Update(pk, update)
	end2 := time.Now()
	cred.NonrevPrepareCache()
	nonce1s, err = common.RandomBigInt(pk.Params.Lstatzk)
	start22 := time.Now()
	proofd, err = cred.CreateDisclosureProof(disclosedAttributes, nil, true, context, nonce1s)
	end22 := time.Now()
	proofDJson, _ = json.Marshal(proofd)
	nonRevProofJson, _ = json.Marshal(proofd.NonRevocationProof)

	start3 := time.Now()
	status = proofd.Verify(pk, context, nonce1s, false)
	zap.S().Infoln("verification status: ", status)
	end3 := time.Now()
	zap.S().Infoln(" Disclosure proof size (bytes): ", len(proofDJson),
		"\t Non revocation proof size (bytes): ", len(nonRevProofJson),
		"\n Time to update witness: ", end2.Sub(start2).Milliseconds(), " ms",
		"\t Time to create disclosure proof: ", end22.Sub(start22).Milliseconds(), " ms",
		"\t Time to verify: ", end3.Sub(start3).Milliseconds(), " ms")

	start4 := time.Now()
	for i := 0; i < 1000; i++ {
		w, err = revocation.RandomWitness(sk, acc)
		acc, event, err = acc.Remove(sk, w.E, event)
		events = append(events, event)
	}
	end4 := time.Now()
	zap.S().Infoln("Revoked 1000 random VCs.  Time for revocation (in seconds): ", end4.Sub(start4).Seconds())
	update, err = revocation.NewUpdate(sk, acc, events)
	start5 := time.Now()
	cred.NonRevocationWitness.Update(pk, update)
	end5 := time.Now()
	cred.NonrevPrepareCache()
	nonce1s, err = common.RandomBigInt(pk.Params.Lstatzk)
	start55 := time.Now()
	proofd, err = cred.CreateDisclosureProof(disclosedAttributes, nil, true, context, nonce1s)
	end55 := time.Now()
	proofDJson, _ = json.Marshal(proofd)
	nonRevProofJson, _ = json.Marshal(proofd.NonRevocationProof)
	start6 := time.Now()
	status = proofd.Verify(pk, context, nonce1s, false)
	end6 := time.Now()
	zap.S().Infoln(" Disclosure proof size (bytes): ", len(proofDJson),
		"\t non revocation proof size (bytes): ", len(nonRevProofJson),
		"\n Time to update witness: ", end5.Sub(start5).Milliseconds(), " ms",
		"\t Time to create disclosure proof: ", end55.Sub(start55).Milliseconds(), " ms",
		"\t Time to verify: ", end6.Sub(start6).Milliseconds(), " ms")

	start7 := time.Now()
	for i := 0; i < 10000; i++ {
		w, err = revocation.RandomWitness(sk, acc)
		acc, event, err = acc.Remove(sk, w.E, event)
		events = append(events, event)
	}
	end7 := time.Now()
	zap.S().Infoln("Revoked 10000 random VCs.  Time for revocation (in seconds): ", end7.Sub(start7).Seconds())
	update, err = revocation.NewUpdate(sk, acc, events)
	start8 := time.Now()
	cred.NonRevocationWitness.Update(pk, update)
	end8 := time.Now()
	//cred.NonrevPrepareCache()
	nonce1s, err = common.RandomBigInt(pk.Params.Lstatzk)
	start88 := time.Now()
	proofd, err = cred.CreateDisclosureProof(disclosedAttributes, nil, true, context, nonce1s)
	end88 := time.Now()
	proofDJson, _ = json.Marshal(proofd)
	nonRevProofJson, _ = json.Marshal(proofd.NonRevocationProof)
	start9 := time.Now()
	status = proofd.Verify(pk, context, nonce1s, false)
	end9 := time.Now()
	zap.S().Infoln(" Disclosure proof size (bytes): ", len(proofDJson),
		"\t Non revocation proof size (bytes): ", len(nonRevProofJson),
		"\n Time to update witness: ", end8.Sub(start8).Milliseconds(), " ms",
		"\t Time to create disclosure proof: ", end88.Sub(start88).Milliseconds(), " ms",
		"\t Time to verify: ", end9.Sub(start9).Milliseconds(), " ms")
}
