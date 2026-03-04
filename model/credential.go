package model

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	bn254_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/signature"
	"go.uber.org/zap"
	"strconv"
	"zkrevoke/crypto2"
	"zkrevoke/utils"
)

/*
VerifiableCredential data model consists of the following elements:
1) Metadata - consists of the different properties mentioned in the w3c vc specifiations
2) Claims - consists of one or more set of claims about subjects.
3) Proofs - consists of a signature: eddsa_signature
*/
type VerifiableCredential struct {
	Id                string `json:"id"`
	Seed              string `json:"seed"`
	Holder_PublicKey  []byte `json:"holder_public_key"`
	Metadata          Metadata
	CredentialSubject []Claims
	Proofs            []Proof
}

/*
URI represents URI (rfc3986). eg:- http://example.com
*/
type URI interface{}

type Metadata struct {
	Contexts         interface{}      `json:"@context"`
	Type             interface{}      `json:"type"`
	Issuer           URI              `json:"issuer"`
	ValidFrom        string           `json:"validFrom"`
	ValidUntil       string           `json:"validUntil"`
	CredentialStatus CredentialStatus `json:"credentialStatus"`
}

/*
CredentialStatus describes information related to discovering the status of a credential

	Id - time-bound token
	Type - TokenListEntry
*/
type CredentialStatus struct {
	Id   URI         `json:"id"`
	Type interface{} `json:"type"`
}

/*
CredentialSubject consists of claims about one or more subjects
*/
type CredentialSubject []Claims
type Claims interface{}

type Proof struct {
	Type         string `json:"type"`
	ProofPurpose string `json:"proofPurpose"`
	Cryptosuite  string `json:"cryptosuite"`
	Created      string `json:"created"`
	Expires      string `json:"expires"`
	Nonce        string `json:"nonce"`
	ProofValue   []byte `json:"proofValue"`
}

// ProofType represents the type of proof used for credential verification
type ProofType string

const (
	ProofTypeZKP   ProofType = "ZKP"
	ProofTypeEDDSA ProofType = "EDDSA"
)

func (cs CredentialStatus) String() string {
	var response string
	response = response + "[Id : " + fmt.Sprintf("%v", cs.Id) + "\t"
	response = response + "Type: " + fmt.Sprintf("%v", cs.Type) + "]\t"
	return response
}

func (proof Proof) String() string {
	var response string
	response = response + " \n"
	response = response + "type : " + fmt.Sprintf("%v", proof.Type) + "\t"
	response = response + "proof purpose : " + fmt.Sprintf("%v", proof.ProofPurpose) + "\t"
	response = response + "crypto suite: " + fmt.Sprintf("%v", proof.Cryptosuite) + "\t"
	response = response + "created at: " + fmt.Sprintf("%v", proof.Created) + "\t"
	response = response + "expires: " + fmt.Sprintf("%v", proof.Expires) + "\t"
	response = response + "nonce: " + fmt.Sprintf("%v", proof.Nonce) + "\t"
	response = response + "proof value : " + utils.GetShortString(hex.EncodeToString(proof.ProofValue)) + "\t"
	return response
}

func (metadata Metadata) String() string {

	var response string

	response = response + "---context: " + fmt.Sprintf("%v", metadata.Contexts) + "\t"
	response = response + "---type: " + fmt.Sprintf("%v", metadata.Type) + "\t"
	response = response + "---entities: " + fmt.Sprintf("%v", metadata.Issuer) + "\t"
	response = response + "---valid from: " + fmt.Sprintf("%v", metadata.ValidFrom) + "\t"
	response = response + "---valid until: " + fmt.Sprintf("%v", metadata.ValidUntil) + "\t"
	response = response + "---credential status: " + metadata.CredentialStatus.String()

	return response
}

func (vc VerifiableCredential) GetId() string {
	return fmt.Sprintf("%v", vc.Id)
}

func (vc VerifiableCredential) GetSeed() string {
	return fmt.Sprintf("%v", vc.Seed)
}

func (vc VerifiableCredential) String() string {

	var response string
	response = response + fmt.Sprintf("{---vc ID: %v", vc.Id) + "\t"
	response = response + fmt.Sprintf("\t---seed: %v", vc.Seed) + "\n"
	response = response + fmt.Sprintf("%v", vc.Metadata) + "\n"
	response = response + fmt.Sprintf("claims: %v", vc.CredentialSubject) + "\n"
	response = response + fmt.Sprintf("holder's public key: %s", utils.GetShortString(hex.EncodeToString(vc.Holder_PublicKey)))
	response = response + fmt.Sprintf("Proofs: %v", vc.Proofs) + "\n}"
	return response
}

/*
SignVC_EDDSA signs the vc_ID, seed, and holder's public key in a VC using eddsa algorithm.
Returns the digital signature
*/
func SignVC_EDDSA(privateKey signature.Signer, seed string, validUntil string, claims EmploymentClaims, selectiveDisclosureExtension bool) []byte {

	//x_point, y_point := crypto.Parse_PublicKey(publicKey_Holder)
	var claims_hash []byte
	if selectiveDisclosureExtension == false {
		claims_hash = claims.GenerateHashDigest()
	} else {
		f := bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployeeID))
		hash1 := f.Sum(nil)
		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployeeName))
		hash2 := f.Sum(nil)
		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployerName))
		hash3 := f.Sum(nil)
		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(claims.EmployeeDesignation))
		hash4 := f.Sum(nil)
		f = bn254_mimc.NewMiMC()
		_, _ = f.Write([]byte(strconv.Itoa(claims.Salary)))
		hash5 := f.Sum(nil)
		f = bn254_mimc.NewMiMC()
		_, _ = f.Write(hash1)
		_, _ = f.Write(hash2)
		_, _ = f.Write(hash3)
		_, _ = f.Write(hash4)
		_, _ = f.Write(hash5)
		claims_hash = f.Sum(nil)
	}

	g := bn254_mimc.NewMiMC()
	_, _ = g.Write([]byte(seed))
	_, _ = g.Write(claims_hash)
	_, _ = g.Write([]byte(validUntil))
	msg := g.Sum(nil)

	//zap.S().Infoln("***CREDENTIAL***: message", msg)
	eddsa_signature, err := crypto2.Sign_EDDSA(privateKey, []byte(msg))

	//zap.S().Infoln("***CREDENTIAL***: signature", eddsa_signature)
	if err != nil {
		zap.S().Infoln("***CREDENTIAL***: signing failed", err)
	} else {
		return eddsa_signature
	}
	return nil
}

func (vc *VerifiableCredential) Json() []byte {
	jsonObj, _ := json.MarshalIndent(vc, "", "    ")
	return jsonObj
}
