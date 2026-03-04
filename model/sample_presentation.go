package model

import (
	"encoding/json"
	"fmt"
)

type SampleEmploymentProofPresentation struct {
	EmployeeDesignation string `json:"employeeDesignation"`
	Salary              int    `json:"salary"`
}

func (s SampleEmploymentProofPresentation) String() string {
	var response string
	response = response + fmt.Sprintf("[---designation : %s", s.EmployeeDesignation) + "\t"
	response = response + fmt.Sprintf("---salary: %v", s.Salary) + "]"
	return response
}

func JsonToEmploymentProofPresentation(jsonObj []byte) *VerifiablePresentation {
	vp := VerifiablePresentation{}
	//var claimsSet []Claims
	json.Unmarshal(jsonObj, &vp)
	var tokenPresentations []TokenPresentation
	for i := 0; i < len(vp.TokenPresentations); i++ {
		jsonData, _ := json.Marshal(vp.TokenPresentations[i])
		var token TokenPresentation
		json.Unmarshal(jsonData, &token)
		tokenPresentations = append(tokenPresentations, token)
	}
	return &vp
}
