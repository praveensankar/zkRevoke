package engine

import (
	"crypto/rand"
	"go.uber.org/zap"
	"math/big"
	"time"
	"zkrevoke/config"
	"zkrevoke/holder"
	"zkrevoke/issuer"
	"zkrevoke/verifier"
)

func RunSimpleScenario(conf config.Config) {
	conf.Params.NumberOfTokensPerCircuit = 1
	issuer := issuer.NewIssuer(&conf)
	//revocationRate := conf.Params.RevocationRateBase
	//NumberOfRevokedVCs := (int(conf.Params.TotalVCs) * revocationRate) / 100

	conf.InitialTimestamp = time.Now()
	issuer.InitialTimeStamp = conf.InitialTimestamp
	address := issuer.DeployContract(conf)
	conf.SmartContractAddress = address
	issuer.SetUpBlockchainConnection(conf)
	issuer.SetupResults(conf)
	issuer.PublishPublicParameters()

	verifier := verifier.NewVerifier()
	verifier.SetUpBlockchainConnection(conf)

	verifier.RetrievePublicParameters()
	verifier.InitialTimeStamp = conf.InitialTimestamp
	verifier.SelectiveDisclosureExtension = conf.SelectiveDisclosureExtension

	var holders []*holder.Holder

	maxObjs, _ := rand.Int(rand.Reader, big.NewInt(int64((conf.Params.TotalVCs))))
	maximumNumberOfHolderObjects := int(maxObjs.Int64())
	for i := 0; i < maximumNumberOfHolderObjects; i++ {
		holder := holder.NewHolder(1)
		vc, pki, duration, initialTimestamp, numberOfTokensInCircuit := issuer.GenerateVC(holder.Holder_PublicKey.Bytes())
		holder.InitCryptoKeys(pki)
		holder.SetUpBlockchainConnection(conf)
		holder.SetDuration(duration)
		holder.SetInitialTimeStamp(initialTimestamp)
		holder.SetNumberOfTokensInCircuit(numberOfTokensInCircuit)
		holder.SelectiveDisclosureExtension = conf.SelectiveDisclosureExtension
		// i-th holder is issued i-th vc
		holder.ReceiveVC(*vc)
		holders = append(holders, holder)
	}

	ticker := time.NewTicker(time.Duration(conf.Params.EpochDuration) * time.Second)

	// At the start of each epoch
	//		issuer refreshes the revoked tokens and publishes it
	//		verifier verifies the revocation status of valid tokens it has received
	//		verifier also acts maliciously and traces the tokens of VPs from previous epochs
	for epoch := 0; epoch < int(conf.Params.ExpirationPeriod); epoch++ {
		zap.S().Infoln("********************************************************************** Epoch: ",
			epoch, " **********************************************************************")
		// issuer revokes 1 VC per epoch
		issuer.RevokeVCsRandomly(2)
		issuer.PublishRevokedTokens()
		challenge := verifier.RequestVP()
		numberOfTokens, _ := rand.Int(rand.Reader, big.NewInt(int64(conf.Params.ExpirationPeriod)))
		numberOfTokensInInt := numberOfTokens.Int64()
		if numberOfTokensInInt == 0 {
			numberOfTokensInInt = 1
		}
		//zap.S().Info("\n ****HOLDER ", i, "\t number of tokens: ", numberOfTokensInt)
		vp := holders[epoch].GenerateVP(int(numberOfTokensInInt), []byte(challenge))

		verifier.ReceiveVP(challenge, *vp)
		verifier.VerifyTokens()
		<-ticker.C
	}

}
