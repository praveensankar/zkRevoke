package verifier

import (
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/common"
	"sync"
	"time"
	"zkrevoke/config"
	"zkrevoke/model"
	"zkrevoke/results"
)

/*
Duration - epoch duration in seconds
InitialTimeStamp - time stamp of the first epoch (epoch 0)
Tokens - key: epoch, value: token. Tokens corresponding to epochs are stored.
*/
type Verifier struct {
	Duration         int
	InitialTimeStamp time.Time
	ccs              constraint.ConstraintSystem

	eddsaPublicKey               signature.PublicKey
	zkpVerifyingKey              groth16.VerifyingKey
	verifiablePresentations      []model.VerifiablePresentation
	Tokens                       map[int][]*TokenStorage
	blockchainRPCEndpoint        string
	smartContractAddress         common.Address
	VPCounter                    int
	Challenge                    map[string]int
	Result                       *results.ResultVerifier
	SelectiveDisclosureExtension bool
	sync.RWMutex
}

type TokenStorage struct {
	Epoch                   int
	Token                   []byte
	VCExpiryDate            int64
	VCValidFrom             int64
	NumberOfEpochsVCisValid int
}

func NewVerifier() *Verifier {
	verifier := Verifier{}
	verifier.Challenge = make(map[string]int)
	verifier.Tokens = make(map[int][]*TokenStorage)
	verifier.VPCounter = 0
	verifier.Result = &results.ResultVerifier{}
	verifier.SelectiveDisclosureExtension = false
	return &verifier
}

func (verifier *Verifier) EnableSelectiveDisclosureExtension() {
	verifier.SelectiveDisclosureExtension = true
}

func (verifier *Verifier) DisableSelectiveDisclosureExtension() {
	verifier.SelectiveDisclosureExtension = false
}

func (verifier *Verifier) SetUpBlockchainConnection(conf config.Config) {
	verifier.blockchainRPCEndpoint = conf.BlockchainRpcEndpoint
	verifier.smartContractAddress = common.HexToAddress(conf.SmartContractAddress)
}

func (verifier *Verifier) SetCCS(cs constraint.ConstraintSystem) {
	verifier.ccs = cs
}

func (verifier *Verifier) SetEddsaPublicKey(eddsaPublicKey signature.PublicKey) {
	verifier.eddsaPublicKey = eddsaPublicKey
}

func (verifier *Verifier) SetZKPVerifyingKey(zkpVerifyingKey groth16.VerifyingKey) {
	verifier.zkpVerifyingKey = zkpVerifyingKey
}

/*
GetCurrentEpoch returns the current epoch
*/
func (verifier Verifier) GetCurrentEpoch() int {
	currentTime := time.Now()
	epoch := (int(currentTime.Sub(verifier.InitialTimeStamp).Seconds())) / verifier.Duration
	return epoch
}

/*
Fetches public parameters needed for cryptography operations
*/
func (verifier *Verifier) RetrievePublicParameters() {

	zkpVerifyingKey := verifier.RetrieveZKPVerifyingKey()
	verifier.SetZKPVerifyingKey(zkpVerifyingKey)

	eddsaPublicKey := verifier.RetrieveEDDSAPublicKey()
	verifier.SetEddsaPublicKey(eddsaPublicKey)

	duration, time := verifier.RetrieveEpochConfigurations()
	verifier.Duration = int(duration)
	verifier.InitialTimeStamp = time

	ccs := verifier.RetrieveCCS()
	verifier.SetCCS(ccs)

}

func (verifier *Verifier) FinalizeResults() {
}
