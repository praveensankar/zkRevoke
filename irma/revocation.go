/*
The programs under the package "irma" are referenced and reused from
"https://github.com/privacybydesign/gabi/blob/master/LICENSE".

Privacy by Design Foundation, Maarten Everts reserves copyright for
the programs under the package "irma".

Redistribution and use of programs under the package "irma" in source and binary forms,
with or without modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of TNO, the IRMA project nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.
*/

package irma

import (
	"errors"
	"fmt"
	"github.com/privacybydesign/gabi/big"
	"github.com/privacybydesign/gabi/gabikeys"
	"github.com/privacybydesign/gabi/revocation"
	"time"
	"zkrevoke/irma/internal/common"
)

//
//func generateKeys() (*gabikeys.PrivateKey, *gabikeys.PublicKey) {
//	N, pprime, qprime, _ := generateGroup()
//
//	ecdsa, _ := signed.GenerateKey()
//
//	sk := &gabikeys.PrivateKey{
//		Counter: 0,
//		ECDSA:   ecdsa,
//		PPrime:  pprime,
//		QPrime:  qprime,
//		N:       N,
//	}
//	sk.Order = new(big.Int).Mul(sk.PPrime, sk.QPrime)
//	pk := &gabikeys.PublicKey{
//		Counter: 0,
//		ECDSA:   &ecdsa.PublicKey,
//		N:       N,
//		G:       common.RandomQR(N),
//		H:       common.RandomQR(N),
//	}
//
//	return sk, pk
//}
//
//func generateGroup() (*big.Int, *big.Int, *big.Int, error) {
//	p, err := safeprime.Generate(32, nil)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	q, err := safeprime.Generate(32, nil)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	n := new(big.Int).Mul(p, q)
//
//	p.Rsh(p, 1)
//	q.Rsh(q, 1)
//
//	return n, p, q, nil
//}

func GenerateNewAccumulator(sk *gabikeys.PrivateKey, pk *gabikeys.PublicKey) (*revocation.Update, error) {
	update, err := revocation.NewAccumulator(sk)
	_, err = update.Verify(pk)
	if err != nil {
		fmt.Println("failed to verify newly generated accumulator")
	}
	return update, err
}

func GenerateWitnessForFactor(sk *gabikeys.PrivateKey, acc *revocation.Accumulator, factor *big.Int) (*revocation.Witness, error) {
	eInverse, ok := common.ModInverse(factor, sk.Order)
	if !ok {
		return nil, errors.New("failed to compute modular inverse")
	}
	u := new(big.Int).Exp(acc.Nu, eInverse, sk.N)
	return &revocation.Witness{U: u, E: factor}, nil
}

func VerifyWitness(witness *revocation.Witness, pk *gabikeys.PublicKey, signedAcc *revocation.SignedAccumulator) bool {
	witness.SignedAccumulator = signedAcc
	err := witness.Verify(pk)
	if err == nil {
		fmt.Println("Factor: ", witness.E, " witness verification successful")
		return true
	} else {
		fmt.Println("Factor: ", witness.E, " witness verification failed")
	}
	return false
}

func RevokeFactor(factor *big.Int, sk *gabikeys.PrivateKey, acc *revocation.Accumulator, update *revocation.Update) (*revocation.Accumulator, *revocation.Event, error) {
	parentevent := update.Events[len(update.Events)-1]

	newAcc, event, err := acc.Remove(sk, factor, parentevent)
	if err != nil {
		fmt.Println("failed to remove factor: ", factor, " from accumulator")
	}
	fmt.Println("removed factor: ", factor, "\t new accumulator: ", newAcc.Nu, "\t event: ", event)
	return newAcc, event, err
}

func GenerateProofForZKP(witness *revocation.Witness, pk *gabikeys.PublicKey) *revocation.Proof {
	start := time.Now()
	randomizer := revocation.NewProofRandomizer()
	list, commit, err := revocation.NewProofCommit(pk, witness, randomizer)
	if err != nil {
		return nil
	}
	challenge := common.HashCommit(list, false)
	proof := commit.BuildProof(challenge)
	proof.Challenge = challenge
	end := time.Now()
	fmt.Println("zkp proof generation time (micro seconds): ", end.Sub(start).Microseconds())
	return proof
}

func VerifyZKP(proof *revocation.Proof, pk *gabikeys.PublicKey, acc *revocation.SignedAccumulator) bool {
	start := time.Now()
	proof.SignedAccumulator = acc
	status := proof.VerifyWithChallenge(pk, proof.Challenge)
	end := time.Now()
	fmt.Println("zkp proof verification: status: ", status, "\t time (nano seconds): ", end.Sub(start).Nanoseconds())
	return status
}
func verifyWitness(u, e *big.Int, acc *revocation.Accumulator, grp *gabikeys.PublicKey) bool {
	return new(big.Int).Exp(u, e, grp.N).Cmp(acc.Nu) == 0
}

type accumulator revocation.Accumulator

func (b accumulator) Base(name string) *big.Int {
	if name == "nu" {
		return b.Nu
	}
	return nil
}

func (b accumulator) Exp(ret *big.Int, name string, exp, n *big.Int) bool {
	if name == "nu" {
		ret.Exp(b.Nu, exp, n)
		return true
	}
	return false
}

func (b accumulator) Names() []string {
	return []string{"nu"}
}
