package issuer

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
	"zkrevoke/config"
	"zkrevoke/crypto2"
	"zkrevoke/model"
	"zkrevoke/results"
	"zkrevoke/utils"
	"zkrevoke/zkp"
)

/*
Duration - epoch duration in seconds
InitialTimeStamp - time stamp of the first epoch (epoch 0)
seedStore - key: vc ID, value: seed
*/
type Issuer struct {
	Duration                     int
	InitialTimeStamp             time.Time
	ccs                          constraint.ConstraintSystem
	eddsaPrivateKey              signature.Signer
	eddsaPublicKey               signature.PublicKey
	zkpProvingKey                groth16.ProvingKey
	zkpVerifyingKey              groth16.VerifyingKey
	verfiableCredentials         map[string]model.VerifiableCredential
	TotalNumberOfEpocs           int
	seedStore                    map[string]string
	RevokedVCIDs                 []string
	blockchainRPCEndpoint        string
	account                      common.Address
	smartContractAddress         common.Address
	privateKey                   string
	gasLimit                     uint64
	gasPrice                     *big.Int
	vcCounter                    int
	NumberOfTokensInCircuit      int
	MinimumExpiryDurationOfVC    int
	Result                       *results.ResultIssuer
	UsePregeneratedVCs           bool
	SelectiveDisclosureExtension bool
	sync.RWMutex
}

/*
NewIssuer creates a new issuer instance.
This function also stores the key pairs in the config
Currently only one pairs of cryptography key pairs are used. Upgrade to multiple keypairs in future if needed.
*/
func NewIssuer(conf *config.Config) *Issuer {
	issuer := &Issuer{}
	issuer.Result = &results.ResultIssuer{}
	issuer.Duration = int(conf.Params.EpochDuration)
	issuer.InitialTimeStamp = conf.InitialTimestamp
	issuer.vcCounter = 0

	issuer.Setup(conf)

	issuer.TotalNumberOfEpocs = int(conf.Params.TotalVCs)
	issuer.MinimumExpiryDurationOfVC = int(conf.Params.ExpirationPeriod)

	issuer.seedStore = make(map[string]string)
	issuer.verfiableCredentials = make(map[string]model.VerifiableCredential)
	pki := utils.PublicParams{}
	pki.Ccs = issuer.ccs
	pki.EddsaPublicKey = issuer.eddsaPublicKey
	pki.ZkpProvingKey = issuer.zkpProvingKey
	pki.ZkpVerifyingKey = issuer.zkpVerifyingKey
	conf.PKI = pki

	if conf.UsePreGeneratedKeysAndVCs == true {
		issuer.LoadStoredVCs(*conf)
	}

	issuer.SelectiveDisclosureExtension = conf.SelectiveDisclosureExtension
	return issuer
}

func (issuer *Issuer) Setup(conf *config.Config) {

	// generate eddsa key
	if conf.UsePreGeneratedKeysAndVCs == true || conf.Run.GenerateVCs == true {

		eddsaPrivateKey, _ := hex.DecodeString("e9370d0522d74b1dc8915b65886028240f4e70edfdd5600af3e5caed28e49c8279a5703413d4dcdd99b5f8b1759c19f7a86e474eb3d0eb6c72cb45b272c41eb0b093d80d38356e4bd1fd955d5a0c7b7ba86f00d7bb7112d14c6cf0a9be3009da")
		issuer.eddsaPrivateKey, issuer.eddsaPublicKey = crypto2.BytesToEDDSAKeys(eddsaPrivateKey)

	} else {
		startEDDSA := time.Now()
		issuer.eddsaPrivateKey, issuer.eddsaPublicKey = crypto2.Generate_EDDSA_Keypairs()
		endEDDSA := time.Since(startEDDSA)
		issuer.Result.SetEDDSAKeyGenTime(endEDDSA)

	}

	startCCS := time.Now()
	// generate zkp circuit
	issuer.NumberOfTokensInCircuit = int(conf.Params.NumberOfTokensPerCircuit)
	// hardcoded 5 claims
	issuer.ccs = zkp.NewCircuit(issuer.NumberOfTokensInCircuit)
	issuer.zkpProvingKey, issuer.zkpVerifyingKey = zkp.SetupGroth(issuer.ccs)
	endCCS := time.Since(startCCS)
	issuer.Result.SetGrothCCSTime(endCCS)

}

func (issuer *Issuer) SetUpBlockchainConnection(conf config.Config) {
	issuer.blockchainRPCEndpoint = conf.BlockchainRpcEndpoint
	issuer.privateKey = conf.PrivateKey
	issuer.gasLimit = conf.GasLimit
	issuer.gasPrice = conf.GasPrice
	issuer.account = common.HexToAddress(conf.Account)
	if conf.DEBUG == true {
		zap.S().Infoln("***ISSUER***: blockchain rpc endpoint: ", issuer.blockchainRPCEndpoint, "\t smart"+
			" contract address: ", issuer.smartContractAddress)
	}
}

func (issuer *Issuer) PublishPublicParameters() {

	_, cost := issuer.PublishHashOfCCS()
	issuer.Result.SetGrothCCSHashCost(cost)

	ccsGasUsed, ccsCost := issuer.PublishCCS()
	issuer.Result.SetGrothCCSCost(ccsCost)
	issuer.Result.SetGrothCCSGas(ccsGasUsed)

	zkpGasUsed, zkpCost := issuer.PublishZKPVerifyingKey()
	issuer.Result.SetGrothVerificationKeyCost(zkpCost)
	issuer.Result.SetGrothVerificationKeyGas(zkpGasUsed)

	eddsaPKGasUsed, eddsaPKCost := issuer.PublishEDDSAPublicKey()
	issuer.Result.SetEDDSAPublicKeyCost(eddsaPKCost)
	issuer.Result.SetEDDSAPublicKeyGas(eddsaPKGasUsed)

	issuer.PublishEpochConfigurations()
}

func (issuer *Issuer) SetupResults(conf config.Config) {

	issuer.Result.SetTotalVCs(int(conf.Params.TotalVCs))

	issuer.Result.SetTotalEpochs(int(conf.Params.ExpirationPeriod))
	issuer.Result.SetEpochDuration(int(conf.Params.EpochDuration))
	issuer.Result.SetNumberOfTokensInCircuit(issuer.NumberOfTokensInCircuit)

	var buf bytes.Buffer
	ccs := issuer.ccs
	_, err := ccs.WriteTo(&buf)
	if err != nil {
		// handle error
	}

	// Get byte slice
	ccsBytes := buf.Bytes()

	var buf1 bytes.Buffer
	w := zlib.NewWriter(&buf1)
	w.Write(ccsBytes)
	w.Close()
	compressedCCSBytes := buf1.Bytes()

	issuer.Result.SetGrothCCSSize(len(compressedCCSBytes))

	zkpProvingKey := issuer.zkpProvingKey
	zkpProvingKeyBytes, _ := zkp.GrothProvingKeyToBytes(zkpProvingKey)
	issuer.Result.SetGrothProvingKeySize(len(zkpProvingKeyBytes))

	zkpVerifyingKey := issuer.zkpVerifyingKey
	zkpVerifyingKeyBytes, _ := zkp.GrothVerifyingKeyToBytes(zkpVerifyingKey)
	issuer.Result.SetGrothVerificationKeySize(len(zkpVerifyingKeyBytes))

	numberOfConstraints := issuer.ccs.GetNbConstraints()
	issuer.Result.SetGrothCCSNumberOfConstraints(int64(numberOfConstraints))

	issuer.Result.SetEDDSAPublicKeySize(len(issuer.eddsaPublicKey.Bytes()))
	issuer.Result.SetEDDSAPrivateKeySize(len(issuer.eddsaPrivateKey.Bytes()))

	//issuer.Result.SetAvgGasPriceMarch2025(int(utils.ReadMarch2025AverageGasPriceFromCSV().Int64()))
}

func (issuer *Issuer) FinalizeResults() {
	issuer.Result.ComputeAvgTokenGenerationTime()
	issuer.Result.ComputeAvgVCGenerationTime()
	issuer.Result.ComputeAvgVCSize()
}

func (issuer *Issuer) Reset() {
	issuer.vcCounter = 0
	issuer.RevokedVCIDs = make([]string, 0)
}

func (issuer *Issuer) ResetRevocationStorage() {
	issuer.RevokedVCIDs = make([]string, 0)
}
