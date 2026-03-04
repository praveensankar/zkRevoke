package test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"zkrevoke/model"
)

func TestCredential(t *testing.T) {
	zap.S().Infoln("Credential TEST")

	vcID := "id#1"
	seed := "seed#1"
	vadlidFromStr := "start_time"
	validUntilStr := "valid_until"

	myVC, _, err := model.CreateEmploymentProofVC(vcID, seed, nil, nil, vadlidFromStr, validUntilStr, false)
	assert.NoError(t, err)
	assert.NotNil(t, myVC)
	zap.S().Infof("Created Employment Proof VC")
	zap.S().Infoln("new VC: ", myVC)
}
