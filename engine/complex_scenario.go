package engine

import (
	"crypto/rand"
	"go.uber.org/zap"
	"math/big"
	"time"
	"zkrevoke/config"
	"zkrevoke/holder"
	"zkrevoke/issuer"
	"zkrevoke/results"
	"zkrevoke/utils"
	"zkrevoke/verifier"
)

func RunComplexScenario(conf config.Config) {

	issuer := issuer.NewIssuer(&conf)
	revocationRate := conf.Params.RevocationRateBase

	NumberOfRevokedVCs := (int(conf.Params.TotalVCs) * revocationRate) / 100

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
		holder := holder.NewHolder(i)
		if conf.UsePreGeneratedKeysAndVCs == true {
			vc, pki, duration, _, numberOfTokensInCircuit := issuer.RequestVC()
			holder.InitCryptoKeys(pki)
			holder.SetUpBlockchainConnection(conf)
			holder.SetDuration(duration)
			holder.SetInitialTimeStamp(conf.InitialTimestamp)
			holder.SetNumberOfTokensInCircuit(numberOfTokensInCircuit)
			holder.SelectiveDisclosureExtension = conf.SelectiveDisclosureExtension
			holder.ReceiveVC(*vc)
		} else {
			vc, pki, duration, initialTimestamp, numberOfTokensInCircuit := issuer.GenerateVC(holder.Holder_PublicKey.Bytes())
			holder.InitCryptoKeys(pki)
			holder.SetUpBlockchainConnection(conf)
			holder.SetDuration(duration)
			holder.SetInitialTimeStamp(initialTimestamp)
			holder.SetNumberOfTokensInCircuit(numberOfTokensInCircuit)
			// i-th holder is issued i-th vc
			holder.ReceiveVC(*vc)
		}
		holders = append(holders, holder)
	}

	// issuer revokes number of vcs specified in the config
	issuer.RevokeVCsRandomly(NumberOfRevokedVCs)

	var holder_results []*results.ResultHolder

	maxHolders := 0
	if maximumNumberOfHolderObjects/4 >= 0 {
		maxHolders = maximumNumberOfHolderObjects / 4
	}
	randomNumberOfHolders := utils.GenerateRandomHoldersUniform(1, maximumNumberOfHolderObjects, maxHolders)
	zap.S().Infoln("Random Number of Holders chosen to share VP with a verifier: ", randomNumberOfHolders)
	for i := 0; i < maxHolders; i++ {
		challenge := verifier.RequestVP()
		numberOfTokens, _ := rand.Int(rand.Reader, big.NewInt(int64(conf.Params.ExpirationPeriod)))
		numberOfTokensInt := int(numberOfTokens.Int64()) % (int(conf.Params.NumberOfTokensPerCircuit))
		if numberOfTokensInt == 0 {
			numberOfTokensInt = 1
		}
		//numberOfTokensInt = int64(conf.NumberOfTokensInaVP)
		//zap.S().Info("\n ****HOLDER ", i, "\t number of tokens: ", numberOfTokensInt)
		vp := holders[randomNumberOfHolders[i]].GenerateVP(int(numberOfTokensInt), []byte(challenge))
		holders[randomNumberOfHolders[i]].IsActive = true
		zap.S().Infoln("****HOLDER ", randomNumberOfHolders[i], "****: generated new VP:", vp.String())

		verifier.ReceiveVP(challenge, *vp)
	}

	ticker := time.NewTicker(time.Duration(conf.Params.EpochDuration) * time.Second)

	zap.S().Infoln("*****REVOKED ", NumberOfRevokedVCs, " VCs *****")

	// At the start of each epoch
	//		issuer refreshes the revoked tokens and publishes it
	//		verifier verifies the revocation status of valid tokens it has received
	//		verifier also acts maliciously and traces the tokens of VPs from previous epochs
	for epoch := 0; epoch < int(conf.Params.ExpirationPeriod); epoch++ {
		zap.S().Infoln("********************************************************************** Epoch: ",
			epoch, " **********************************************************************")
		issuer.PublishRevokedTokens()
		verifier.VerifyTokens()
		randomNumberOfHolders := utils.GenerateRandomHoldersUniform(0, maximumNumberOfHolderObjects, maxHolders/2)

		for y := 0; y < maxHolders/2; y++ {
			challenge := verifier.RequestVP()
			numberOfTokens, _ := rand.Int(rand.Reader, big.NewInt(int64(conf.Params.ExpirationPeriod)))
			numberOfTokensInt := int(numberOfTokens.Int64()) % (int(conf.Params.NumberOfTokensPerCircuit))
			if numberOfTokensInt == 0 {
				numberOfTokensInt = 1
			}
			//numberOfTokensInt = int64(conf.NumberOfTokensInaVP)
			//zap.S().Info("\n ****HOLDER ", i, "\t number of tokens: ", numberOfTokensInt)
			vp := holders[randomNumberOfHolders[y]].GenerateVP(int(numberOfTokensInt), []byte(challenge))
			holders[randomNumberOfHolders[y]].IsActive = true
			zap.S().Infoln("****HOLDER ", randomNumberOfHolders[y], "****: generated new VP:", vp.String())

			verifier.ReceiveVP(challenge, *vp)
		}
		<-ticker.C
	}

	for z := 0; z < maximumNumberOfHolderObjects; z++ {
		if holders[z].IsActive == true {
			holders[z].FinalizeResults()
			holder_results = append(holder_results, holders[z].Result)
		}
	}

	zap.S().Infoln("********************************************************************** Final Results",
		" **********************************************************************")

}
