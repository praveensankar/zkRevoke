package holder

import (
	"fmt"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/common"
	"sync"
	"time"
	"zkrevoke/config"
	"zkrevoke/crypto2"
	"zkrevoke/model"
	"zkrevoke/results"
	"zkrevoke/utils"
)

/*
Duration - epoch duration in seconds
InitialTimeStamp - time stamp of the first epoch (epoch 0)
*/
type Holder struct {
	Duration                     int
	InitialTimeStamp             time.Time
	ccs                          constraint.ConstraintSystem
	eddsaPublicKey               signature.PublicKey
	zkpProvingKey                groth16.ProvingKey
	verfiableCredentials         []model.VerifiableCredential
	blockchainRPCEndpoint        string
	smartContractAddress         common.Address
	ID                           string
	holder_PrivateKey            signature.Signer
	Holder_PublicKey             signature.PublicKey
	NumberOfTokensInCircuit      int
	Result                       *results.ResultHolder
	IsActive                     bool
	SelectiveDisclosureExtension bool
	sync.RWMutex
}

func NewHolder(id int) *Holder {
	holder := Holder{}
	holder.ID = fmt.Sprintf("HOLDER#%d", id)
	holder.holder_PrivateKey, holder.Holder_PublicKey = crypto2.Generate_EDDSA_Keypairs()
	holder.Result = &results.ResultHolder{}
	holder.IsActive = false
	holder.SelectiveDisclosureExtension = false
	return &holder
}

func (holder *Holder) SetDuration(duration int) {
	holder.Duration = duration
	holder.Result.SetEpochDuration(duration)
}

func (holder *Holder) SetNumberOfTokensInCircuit(numberOfTokensInCircuit int) {
	holder.NumberOfTokensInCircuit = numberOfTokensInCircuit
	holder.Result.SetNumberOfTokensInVP(numberOfTokensInCircuit)
}

func (holder *Holder) SetInitialTimeStamp(initialTimeStamp time.Time) {
	holder.InitialTimeStamp = initialTimeStamp
}

func (holder *Holder) EnableSelectiveDisclosureExtension() {
	holder.SelectiveDisclosureExtension = true
}

func (holder *Holder) DisableSelectiveDisclosureExtension() {
	holder.SelectiveDisclosureExtension = false
}

func (holder *Holder) InitCryptoKeys(pki *utils.PublicParams) error {
	holder.SetCCS(pki.Ccs)
	holder.SetEddsaPublicKey(pki.EddsaPublicKey)
	holder.SetZKPProvingKey(pki.ZkpProvingKey)
	return nil
}

func (holder *Holder) SetUpBlockchainConnection(conf config.Config) {
	holder.blockchainRPCEndpoint = conf.BlockchainRpcEndpoint
	holder.smartContractAddress = common.HexToAddress(conf.SmartContractAddress)
	//zap.S().Infoln("***", holder.ID, "***\t blockchain rpc endpoint: ", holder.blockchainRPCEndpoint,
	//	"\t smart contract address: ", holder.smartContractAddress)
}

func (holder *Holder) SetCCS(cs constraint.ConstraintSystem) {
	holder.ccs = cs
}

func (holder *Holder) SetEddsaPublicKey(eddsaPublicKey signature.PublicKey) {
	holder.eddsaPublicKey = eddsaPublicKey
}

func (holder *Holder) SetZKPProvingKey(provingKey groth16.ProvingKey) {
	holder.zkpProvingKey = provingKey
}

func (holder *Holder) ReceiveVC(vc model.VerifiableCredential) {
	holder.Lock()
	holder.verfiableCredentials = append(holder.verfiableCredentials, vc)
	holder.Unlock()
}

/*
GetCurrentEpoch returns the current epoch
*/
func (holder Holder) GetCurrentEpoch() int {
	currentTime := time.Now()
	epoch := (int(currentTime.Sub(holder.InitialTimeStamp).Seconds())) / holder.Duration
	return epoch
}

func (holder *Holder) FinalizeResults() {
	holder.Result.ComputeAvgGrothProofSize()
	holder.Result.ComputeAvgGrothProofGenerationTime()
	holder.Result.ComputeAvgTokenGenerationTime()
}
