package model

import (
	"encoding/json"
	"fmt"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/signature"
	"strconv"
	"time"
)

type EmploymentClaims struct {
	EmployeeID          string `json:"employee_id"`
	EmployeeName        string `json:"employee_name"`
	EmployerName        string `json:"employer_name"`
	EmployeeDesignation string `json:"employee_designation"`
	Salary              int    `json:"salary"`
}

func (claims EmploymentClaims) GenerateHashDigest() []byte {
	f := bn254_mimc.NewMiMC()
	_, _ = f.Write([]byte(claims.EmployeeID))
	_, _ = f.Write([]byte(claims.EmployeeName))
	_, _ = f.Write([]byte(claims.EmployerName))
	_, _ = f.Write([]byte(claims.EmployeeDesignation))
	_, _ = f.Write([]byte(strconv.Itoa(claims.Salary)))
	return f.Sum(nil)
}

/*
CreateEmploymentProofVC creates an employment proof VC
Inputs:

	(EmploymentClaims) claims -  set of claims about an employee.
	(skEDDSA) signature.Signer - private key for eddsa

	(validUntil) DateTimeStamp - expiry time of the VC

Output:

	(VerifiableCredential) - a new employment proof vc
	time to sign seed using eddsa scheme (time.Duration)
*/
func CreateEmploymentProofVC(vcID string, seed string, pkHolder []byte, skEDDSA signature.Signer, validFrom string, validUntil string, selectiveDisclosureExtension bool) (*VerifiableCredential, time.Duration, error) {

	myClaims := EmploymentClaims{
		EmployeeName:        "Bob",
		EmployeeID:          "employee#1",
		EmployerName:        "UiO",
		EmployeeDesignation: "PhD Research Fellow",
		Salary:              500000,
	}

	var claimsSet []Claims
	claimsSet = append(claimsSet, myClaims)
	startEDDSA := time.Now()
	signature_eddsa := SignVC_EDDSA(skEDDSA, seed, validUntil, myClaims, selectiveDisclosureExtension)
	endEDDSA := time.Since(startEDDSA)
	proof_eddsa := Proof{
		Type:         string(ProofTypeEDDSA),
		ProofPurpose: "ZKP Circuit",
		Cryptosuite:  "eddsa",
		Created:      "",
		Expires:      "",
		Nonce:        "",
		ProofValue:   signature_eddsa,
	}

	var vcType []string
	vcType = append(vcType, "Employment Proof")

	employmentProofVC := VerifiableCredential{
		Id:               vcID,
		Seed:             seed,
		Holder_PublicKey: pkHolder,
		Metadata: Metadata{
			Contexts:   "Verifiable Credential",
			Type:       vcType,
			Issuer:     "University",
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
		CredentialSubject: claimsSet,
		Proofs:            []Proof{proof_eddsa},
	}
	return &employmentProofVC, endEDDSA, nil
}

func (e EmploymentClaims) String() string {
	var response string
	response = response + "employee name: " + e.EmployeeName + "\t"
	response = response + "employee id: " + fmt.Sprintf("%v", e.EmployeeID) + "\t"
	response = response + "employer name: " + e.EmployerName + "\t"
	response = response + "designation: " + e.EmployeeDesignation + "\t"
	response = response + "salary: " + fmt.Sprintf("%v", e.Salary)
	return response
}

func JsonToEmploymentProofVC(jsonObj []byte) *VerifiableCredential {
	credential := VerifiableCredential{}
	//var claimsSet []Claims
	json.Unmarshal(jsonObj, &credential)
	var claimSet []Claims
	for i := 0; i < len(credential.CredentialSubject); i++ {
		jsonData, _ := json.Marshal(credential.CredentialSubject[i])
		var claims EmploymentClaims
		json.Unmarshal(jsonData, &claims)
		claimSet = append(claimSet, claims)
	}
	var proofs []Proof
	for i := 0; i < len(credential.Proofs); i++ {
		jsonData, _ := json.Marshal(credential.Proofs[i])
		var proof Proof
		json.Unmarshal(jsonData, &proof)
		proofs = append(proofs, proof)
	}
	credential.CredentialSubject = claimSet
	credential.Proofs = proofs
	return &credential
}
