package issuer

import (
	"go.uber.org/zap"
	"zkrevoke/config"
)

func TestIssuer(conf *config.Config) {
	issuer := NewIssuer(conf)
	zap.S().Infoln("****ISSUER TEST*****")

	numberOfIssuedVCs := 10
	numberOfRevokedVCs := 4

	issuer.BulkIssueVCs(numberOfIssuedVCs)
	issuer.RevokeVCsRandomly(numberOfRevokedVCs)

	zap.S().Info("Issued VC IDs: ", issuer.GetIssuedVCIDs())
	zap.S().Info("Revoked VC IDs: ", issuer.RevokedVCIDs)
}
