package model

import (
	"crypto/rand"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zkrevoke/config"
	crypto "zkrevoke/crypto2"
)

func TestVC(config config.Config) {

	//vcID := make([]byte, 10)
	//rand.Read(vcID)

	//seed := make([]byte, 20)
	//rand.Read(seed)
	vcID := rand.Text()
	seed := rand.Text()

	zap.S().Infoln("***CREDENTIAL***: vcID: ", vcID, "\t seed: ", seed)
	skEDDSA, _ := crypto.Generate_EDDSA_Keypairs()
	_, pkHolder := crypto.Generate_EDDSA_Keypairs()
	zap.S().Infoln("***CREDENTIAL***: eddsa private key: ", skEDDSA)

	validFrom := time.Now()
	validUntilStr := strconv.Itoa(int(validFrom.Add(time.Duration(100) * time.Hour).Unix()))
	validFromStr := strconv.Itoa(int(validFrom.Unix()))
	myVC, _, err := CreateEmploymentProofVC(vcID, seed, pkHolder.Bytes(), skEDDSA, validFromStr, validUntilStr, false)

	if err != nil {
		zap.S().Errorln("err:", err)
	}
	zap.S().Infoln("new VC:", myVC.String())
	myVCJson := myVC.Json()
	zap.S().Infoln("***CREDENTIAL***: json format:", string(myVCJson))
	vc1 := JsonToEmploymentProofVC(myVCJson)
	zap.S().Infoln("***CREDENTIAL***: json to vc:", vc1)

	vp := VerifiablePresentation{}
	var selectClaims SampleEmploymentProofPresentation
	claimSet := vc1.CredentialSubject
	claims := claimSet[0].(EmploymentClaims)
	selectClaims.EmployeeDesignation = claims.EmployeeDesignation
	selectClaims.Salary = claims.Salary
	vp.Messages = selectClaims
	vp.ValidFrom = vc1.Metadata.ValidFrom
	vp.ValidUntil = vc1.Metadata.ValidUntil

	zap.S().Infoln("VP: ", vp.String())

	vpJson := vp.Json()
	zap.S().Infoln("VP json format:", string(vpJson))

	vp1 := JsonToEmploymentProofPresentation(vpJson)
	zap.S().Infoln("VP from json:", vp1)

}
