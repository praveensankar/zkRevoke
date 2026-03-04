package benchmark

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"zkrevoke/crypto2"
	"zkrevoke/utils"
	"zkrevoke/zkp"
)

type ResultCircuitConstraintsType string

const (
	SignatureVerificationCircuit ResultCircuitConstraintsType = "signature_verification"
	TokenVerificationCicuit      ResultCircuitConstraintsType = "token_verification"
	ChallengeVerificationCircuit ResultCircuitConstraintsType = "challenge_verification"
	EmptyCircuit                 ResultCircuitConstraintsType = "empty_circuit"
	CompleteCircuit              ResultCircuitConstraintsType = "complete_circuit"
)

type ResultCircuitConstraints struct {
	CircuitType           ResultCircuitConstraintsType `json:"circuit_type"`
	NumberOfConstraints   int                          `json:"number_of_constraints"`
	PrivateWitnessGenTime int                          `json:"private_witness_gen_time"`
	ProofGenTime          int                          `json:"proof_gen_time"`
	PublicWitnessGenTime  int                          `json:"public_witness_gen_time"`
	ProofVerifyTime       int                          `json:"proof_verify_time"`
}

func ComputeAverageResultCircuitConstraints(results []ResultCircuitConstraints) []*ResultCircuitConstraints {

	if len(results) == 0 {
		filename := fmt.Sprintf("benchmark/results/result_circuit_constraints.json")

		jsonFile, _ := os.Open(filename)
		resJson, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(resJson, &results)

	}

	// key: NumberOfRevokedVCs
	AvgResults := make(map[ResultCircuitConstraintsType]*ResultCircuitConstraints)

	for _, res := range results {
		circuitType := res.CircuitType
		_, exists := AvgResults[circuitType]
		if exists {
			AvgResults[res.CircuitType].NumberOfConstraints = (AvgResults[res.CircuitType].NumberOfConstraints + res.NumberOfConstraints) / 2
			AvgResults[res.CircuitType].ProofGenTime = (AvgResults[res.CircuitType].ProofGenTime + res.ProofGenTime) / 2
			AvgResults[res.CircuitType].PrivateWitnessGenTime = (AvgResults[res.CircuitType].PrivateWitnessGenTime + res.PrivateWitnessGenTime) / 2
			AvgResults[res.CircuitType].PublicWitnessGenTime = (AvgResults[res.CircuitType].PublicWitnessGenTime + res.PublicWitnessGenTime) / 2
			AvgResults[res.CircuitType].ProofVerifyTime = (AvgResults[res.CircuitType].ProofVerifyTime + res.ProofVerifyTime) / 2

		} else {
			result := &ResultCircuitConstraints{}
			result.CircuitType = res.CircuitType
			result.NumberOfConstraints = res.NumberOfConstraints
			result.ProofGenTime = res.ProofGenTime
			result.PrivateWitnessGenTime = res.PrivateWitnessGenTime
			result.PublicWitnessGenTime = res.PublicWitnessGenTime
			result.ProofVerifyTime = res.ProofVerifyTime
			AvgResults[res.CircuitType] = result
		}
	}

	var finalResults []*ResultCircuitConstraints
	for _, res := range AvgResults {
		finalResults = append(finalResults, res)
	}
	return finalResults
}

func (r *ResultCircuitConstraints) Json() ([]byte, error) {
	//return json.MarshalIndent(r, "","    ")
	return json.Marshal(r)
}

func JsonToResultCircuitConstraints(jsonObj []byte) *ResultCircuitConstraints {
	res := ResultCircuitConstraints{}
	json.Unmarshal(jsonObj, &res)
	return &res
}

func WriteResultCircuitConstraintsToFile(result ResultCircuitConstraints, isavg bool) {

	var results []ResultCircuitConstraints
	filename := fmt.Sprintf("benchmark/results/result_circuit_constraints.json")
	if isavg {
		filename = fmt.Sprintf("benchmark/results/result_circuit_constraints_avg.json")
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
	//zap.S().Errorln("RESULTS - successfully written to the file")
}

func (r ResultCircuitConstraints) String() string {
	var response string
	response = response + "Number of constraints in the circuit: " + fmt.Sprintf("%d", r.NumberOfConstraints) + "\t"
	response = response + "Private witness generation time (in micro seconds): " + fmt.Sprintf("%d", r.PrivateWitnessGenTime) + "\t"
	response = response + "Public witness generation time (in micro seconds): " + fmt.Sprintf("%d", r.PublicWitnessGenTime) + "\t"
	response = response + "Proof generation time (in micro seconds): " + fmt.Sprintf("%d", r.ProofGenTime) + "\t"
	response = response + "Proof Verification time (in micro seconds): " + fmt.Sprintf("%d", r.ProofVerifyTime) + "\n"
	return response
}

func ResetResultCircuitConstraints() {
	filename1 := fmt.Sprintf("benchmark/results/result_circuit_constraints.json")
	filename2 := fmt.Sprintf("benchmark/results/result_circuit_constraints_avg.json")
	os.Remove(filename1)
	os.Remove(filename2)
}

func BenchmarkEmptyCircuit() *ResultCircuitConstraints {

	res := &ResultCircuitConstraints{}
	res.CircuitType = EmptyCircuit
	ccs := zkp.NewCircuitEmpty()

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := zkp.SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	res.NumberOfConstraints = ccs.GetNbConstraints()

	// Holder generates a witness to prove the correctness of the index
	start := time.Now()

	witness := zkp.PrivateWitnessGenerationForEmptyCircuit()

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)

	res.PrivateWitnessGenTime = int(end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := zkp.ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	res.ProofGenTime = int(end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := zkp.PublicWitnessGenerationForEmptyCircuit()
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	res.PublicWitnessGenTime = int(end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := zkp.VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	res.ProofVerifyTime = int(end.Microseconds())
	if status == false {
		zap.S().Errorln("unable to verify the groth proof")
	}
	return res
}

func BenchmarkTokenVerificationCircuit() *ResultCircuitConstraints {

	res := &ResultCircuitConstraints{}
	res.CircuitType = TokenVerificationCicuit

	epoch := 1
	seed := rand.Text()
	token := utils.ComputeToken(epoch, seed)
	var epochs [][]byte
	var tokens [][]byte
	epochs = append(epochs, []byte(strconv.Itoa(epoch)))
	tokens = append(tokens, token)
	ccs := zkp.NewCircuitForTokenVerification(1)

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := zkp.SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)
	res.NumberOfConstraints = ccs.GetNbConstraints()
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := zkp.PrivateWitnessParameters{
		Seed: []byte(seed),
	}

	pubParams := zkp.PublicWitnessParameters{
		Epochs: epochs,
		Tokens: tokens,
	}
	witness := zkp.PrivateWitnessGenerationForTokenVerificationCircuit(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	res.PrivateWitnessGenTime = int(end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := zkp.ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	res.ProofGenTime = int(end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := zkp.PublicWitnessGenerationForTokenVerificationCircuit(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	res.PublicWitnessGenTime = int(end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := zkp.VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	if status == false {
		zap.S().Errorln("unable to verify the groth proof")
	}
	res.ProofVerifyTime = int(end.Microseconds())

	return res

}

func BenchmarkChallengeVerification() *ResultCircuitConstraints {

	res := &ResultCircuitConstraints{}
	res.CircuitType = ChallengeVerificationCircuit

	challenge := rand.Text()
	holder_randomness := rand.Text()
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write([]byte(challenge))
	_, _ = f.Write([]byte(holder_randomness))

	msg := f.Sum(nil)

	ccs := zkp.NewCircuitForChallengeVerification()

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := zkp.SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	res.NumberOfConstraints = ccs.GetNbConstraints()

	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := zkp.PrivateWitnessParameters{
		HolderRadomness: []byte(holder_randomness),
	}
	pubParams := zkp.PublicWitnessParameters{
		Hash1:     msg,
		Challenge: []byte(challenge),
	}
	witness := zkp.PrivateWitnessGenerationForChallengeVerification(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	res.PrivateWitnessGenTime = int(end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := zkp.ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	res.ProofGenTime = int(end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := zkp.PublicWitnessGenerationForChallengeVerification(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	res.PublicWitnessGenTime = int(end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := zkp.VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	res.ProofVerifyTime = int(end.Microseconds())
	if status == false {
		zap.S().Errorln("unable to verify the groth proof")
	}
	return res

}

func BenchmarkSignatureVerificationCircuit() *ResultCircuitConstraints {

	res := &ResultCircuitConstraints{}
	res.CircuitType = SignatureVerificationCircuit

	seed := rand.Text()
	validUntil := time.Now().Add(time.Duration(365) * 86400 * time.Second)
	validUntilStr := strconv.Itoa(int(validUntil.Unix()))

	var claims [][]byte
	claims = append(claims, []byte("employee_name:bob"))
	claims = append(claims, []byte("employee_id:employee#1"))
	claims = append(claims, []byte("employer_name:UiO"))
	claims = append(claims, []byte("employee_designation:PhD Research Fellow"))
	claims = append(claims, []byte("salary:500000"))
	privateKey, publicKey := crypto2.Generate_EDDSA_Keypairs()

	f := bn254_mimc.NewMiMC()
	_, _ = f.Write([]byte(seed))
	_, _ = f.Write(zkp.ComputeHashOnClaims(claims))
	_, _ = f.Write([]byte(validUntilStr))
	msg := f.Sum(nil)

	signature, _ := crypto2.Sign_EDDSA(privateKey, msg)

	//signature_Holder, _ := crypto.Sign_EDDSA(privateKey_Holder, msg2)
	ccs := zkp.NewCircuitSignatureVerification()

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := zkp.SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	res.NumberOfConstraints = ccs.GetNbConstraints()
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := zkp.PrivateWitnessParameters{
		Seed:      []byte(seed),
		Signature: signature,
	}
	pubParams := zkp.PublicWitnessParameters{
		PublicKey:  publicKey,
		ValidUntil: []byte(validUntilStr),
		ClaimsHash: zkp.ComputeHashOnClaims(claims),
	}
	witness := zkp.PrivateWitnessGenerationForCircuitSignatureVerification(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	res.PrivateWitnessGenTime = int(end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := zkp.ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	res.ProofGenTime = int(end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := zkp.PublicWitnessGenerationForCircuitSignatureVerification(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	res.PublicWitnessGenTime = int(end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := zkp.VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	res.ProofVerifyTime = int(end.Microseconds())
	if status == false {
		zap.S().Errorln("unable to verify the groth proof")
	}

	return res
}

func BenchmarkCompleteCircuit() *ResultCircuitConstraints {

	res := &ResultCircuitConstraints{}
	res.CircuitType = CompleteCircuit
	seed := rand.Text()
	validUntil := time.Now().Add(time.Duration(365) * 86400 * time.Second)
	validUntilStr := strconv.Itoa(int(validUntil.Unix()))

	var claims [][]byte
	claims = append(claims, []byte("employee_name:bob"))
	claims = append(claims, []byte("employee_id:employee#1"))
	claims = append(claims, []byte("employer_name:UiO"))
	claims = append(claims, []byte("employee_designation:PhD Research Fellow"))
	claims = append(claims, []byte("salary:500000"))
	privateKey, publicKey := crypto2.Generate_EDDSA_Keypairs()

	f := bn254_mimc.NewMiMC()
	_, _ = f.Write([]byte(seed))
	_, _ = f.Write(zkp.ComputeHashOnClaims(claims))
	_, _ = f.Write([]byte(validUntilStr))
	msg := f.Sum(nil)

	signature, _ := crypto2.Sign_EDDSA(privateKey, msg)

	challenge := rand.Text()
	holder_randomness := rand.Text()
	f1 := bn254_mimc.NewMiMC()
	_, _ = f1.Write([]byte(challenge))
	_, _ = f1.Write([]byte(holder_randomness))

	msg2 := f1.Sum(nil)

	epoch := 1
	token := utils.ComputeToken(epoch, seed)
	var epochs [][]byte
	var tokens [][]byte
	epochs = append(epochs, []byte(strconv.Itoa(epoch)))
	tokens = append(tokens, token)

	//signature_Holder, _ := crypto.Sign_EDDSA(privateKey_Holder, msg2)
	ccs := zkp.NewCircuit(1)

	// Issuer then creates a proving key and verification key for the circuit
	pk, vk := zkp.SetupGroth(ccs)
	buf := new(bytes.Buffer)
	pk.WriteTo(buf)

	buf = new(bytes.Buffer)
	vk.WriteTo(buf)

	res.NumberOfConstraints = ccs.GetNbConstraints()
	// Holder generates a witness to prove the correctness of the index
	start := time.Now()
	privParams := zkp.PrivateWitnessParameters{
		Seed:            []byte(seed),
		HolderRadomness: []byte(holder_randomness),
		Signature:       signature,
	}
	pubParams := zkp.PublicWitnessParameters{
		PublicKey:  publicKey,
		Hash1:      msg2,
		Challenge:  []byte(challenge),
		Epochs:     epochs,
		ValidUntil: []byte(validUntilStr),
		Tokens:     tokens,
		ClaimsHash: zkp.ComputeHashOnClaims(claims),
	}
	witness := zkp.PrivateWitnessGeneration(pubParams, privParams)

	end := time.Since(start)
	buf = new(bytes.Buffer)
	witness.WriteTo(buf)
	res.PrivateWitnessGenTime = int(end.Microseconds())

	// Holder creates a proof to prove the correctness of the index. Proof takes the witness as the input.
	start = time.Now()
	proof := zkp.ProveGroth(ccs, pk, witness)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	proof.WriteTo(buf)
	res.ProofGenTime = int(end.Microseconds())

	// Verifier creates a public witness based on the public inputs given by the holder
	start = time.Now()
	publicWitness := zkp.PublicWitnessGeneration(pubParams)
	end = time.Since(start)
	buf = new(bytes.Buffer)
	publicWitness.WriteTo(buf)
	res.PublicWitnessGenTime = int(end.Microseconds())

	// Verifier verify the proof using the verification key and the public witness
	start = time.Now()
	status := zkp.VerifyGroth(proof, vk, publicWitness)
	end = time.Since(start)
	res.ProofVerifyTime = int(end.Microseconds())
	if status == false {
		zap.S().Errorln("unable to verify the groth proof")
	}

	return res
}
