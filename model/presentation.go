package model

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"zkrevoke/utils"
)

type VerifiablePresentation struct {
	Messages                    interface{}                 `json:"claims"`
	TokenPresentations          []TokenPresentation         `json:"token_presentations"`
	Holder_randomness           []byte                      `json:"holder_randomness"`
	Hash1                       []byte                      `json:"hash1"`
	ClaimsHash                  []byte                      `json:"claims_hash"`
	ValidFrom                   string                      `json:"validFrom"`
	ValidUntil                  string                      `json:"validUntil"`
	SelectiveDisclosureElements SelectiveDisclosureElements `json:"selective_disclosure_elements"`
}

type SelectiveDisclosureElements struct {
	SelectClaims interface{} `json:"select_claims"`
	Indexes      []int       `json:"indexes"`
}

type TokenPresentation struct {
	Epochs   []uint   `json:"epochs"`
	Tokens   [][]byte `json:"tokens"`
	ZKPProof Proof    `json:"zkp_proof"`
}

func (vp VerifiablePresentation) String() string {

	var response string
	response = response + fmt.Sprintf("{claims: %v", vp.Messages) + "\t"
	response = response + "(holder randomness: " + utils.GetShortString(hex.EncodeToString(vp.Holder_randomness)) + "), \n "
	response = response + "(fresheness response: " + utils.GetShortString(hex.EncodeToString(vp.Hash1)) + "), "
	response = response + "(valid from: " + fmt.Sprintf("%v", vp.ValidFrom) + "), "
	response = response + "(valid until: " + fmt.Sprintf("%v", vp.ValidUntil) + "), "
	for i := 0; i < len(vp.TokenPresentations); i++ {
		response = response + "\n {---Tokens: ["
		for j := 0; j < len(vp.TokenPresentations[i].Epochs); j++ {
			response = response + fmt.Sprintf("( epoch: %d, ", vp.TokenPresentations[i].Epochs[j])
			response = response + " token: " + utils.GetShortString(hex.EncodeToString(vp.TokenPresentations[i].Tokens[j])) + ")"
		}
		response = response + "],  "
		response = response + "(zkp proof: " + utils.GetShortString(hex.EncodeToString(vp.TokenPresentations[i].ZKPProof.ProofValue)) + ")"
		response = response + "}"
	}

	return response
}

func (vp *VerifiablePresentation) Json() []byte {
	jsonObj, _ := json.MarshalIndent(vp, "", "    ")
	return jsonObj
}
