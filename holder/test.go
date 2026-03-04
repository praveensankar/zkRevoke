package holder

import (
	"crypto/rand"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zkrevoke/config"
	"zkrevoke/crypto2"
	"zkrevoke/model"
	"zkrevoke/zkp"
)

func TestHolder(conf config.Config) {
	holder := NewHolder(0)

	holder.Duration = int(conf.Params.EpochDuration)
	holder.InitialTimeStamp = conf.InitialTimestamp
	ccs := zkp.NewCircuit(int(conf.Params.NumberOfTokensPerCircuit))
	zkpProvingKey, _ := zkp.SetupGroth(ccs)
	eddsaPrivateKey, eddsaPublicKey := crypto2.Generate_EDDSA_Keypairs()

	_, pkHolder := crypto2.Generate_EDDSA_Keypairs()

	holder.SetCCS(ccs)
	holder.SetEddsaPublicKey(eddsaPublicKey)

	holder.SetZKPProvingKey(zkpProvingKey)

	vcID := rand.Text()
	seed := rand.Text()
	validFrom := time.Now()
	validFromStr := strconv.Itoa(int(validFrom.Unix()))
	validUntilStr := strconv.Itoa(int(validFrom.Add(time.Duration(100) * time.Hour).Unix()))
	newVC, _, _ := model.CreateEmploymentProofVC(vcID, seed, pkHolder.Bytes(), eddsaPrivateKey, validFromStr, validUntilStr, false)

	zap.S().Infoln("****HOLDER****: received new VC:", newVC)

	numberOfEpochs := 1
	challenge := []byte("Challenge#123")
	vp := holder.GenerateVP(numberOfEpochs, challenge)
	zap.S().Infoln("****HOLDER****: generated new VP:", vp.String())

	numberOfEpochs = 5
	vp2 := holder.GenerateVP(numberOfEpochs, challenge)
	zap.S().Infoln("****HOLDER****: generated new VP:", vp2.String())

}
