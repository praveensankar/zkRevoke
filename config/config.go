/*
config package sets up necessary configurations to benchmark the system

*/

package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"math/big"
	"strconv"
	"time"
	"zkrevoke/utils"
)

/*
Todo: Remove PKI from here once smart contracts are ready
*/
type Config struct {
	SmartContractAddress         string
	Account                      string
	PrivateKey                   string
	BlockchainRpcEndpoint        string
	BlockchainWebSocketEndPoint  string
	GasLimit                     uint64
	GasPrice                     *big.Int
	LoggerType                   string
	LoggerFile                   string
	LoggerOutputMode             string
	DEBUG                        bool
	TokenStorageCostExps         uint
	TokenStorageNumberOfTokens   []int
	PKI                          utils.PublicParams
	UsePreGeneratedKeysAndVCs    bool
	Run                          Run
	Benchmark                    Benchmark
	Params                       Params
	StartingEpoch                int
	InitialTimestamp             time.Time
	SelectiveDisclosureExtension bool
}

func (config Config) printConfig() {
	zap.L().Info(config.Run.String())
	zap.L().Info("\n\n--------------------------------------------------------Epoch related configuration--------------------------------------------------")
	zap.L().Info("initial timestamp: " + config.InitialTimestamp.String())
	zap.L().Info("********************************************************************************************************************************\n")
	zap.L().Info("\n\n--------------------------------------------------------Logger related configuration--------------------------------------------------")
	zap.L().Info("logger environment: " + config.LoggerType)
	zap.L().Info("logger output file name: " + config.LoggerFile)
	zap.L().Info("logger output mode: " + config.LoggerOutputMode)
	zap.L().Info("********************************************************************************************************************************\n")

	zap.L().Info("\n\n--------------------------------------------------------blockchain related configuration--------------------------------------------------")
	zap.L().Info("account:" + config.Account)
	zap.L().Info("smart contract address:" + config.SmartContractAddress)
	zap.L().Info("blockchain rpc endpoint: " + config.BlockchainRpcEndpoint)
	zap.L().Info("gas limit: " + string(config.GasLimit))
	zap.L().Info("gas price: " + config.GasPrice.String())
	zap.L().Info("********************************************************************************************************************************\n")
	zap.L().Info("\n\n--------------------------------------------------------printing issuer configuration--------------------------------------------------")

	zap.L().Info("Starting epoch (since VCs are pregenerated, the starting epoch should add the amount of time elapsed from the issuance time till current time):" + strconv.Itoa(int(config.StartingEpoch)))
	zap.L().Info("********************************************************************************************************************************\n")
	zap.L().Info("--------------------------------------------------------printing verifier configuration--------------------------------------------------")

	zap.L().Info("********************************************************************************************************************************\n")
	zap.L().Info("--------------------------------------------------------token storage cost calculationconfiguration--------------------------------------------------")
	zap.S().Infoln("number of tokens to calculate cost:", config.TokenStorageNumberOfTokens)
	zap.S().Infoln("number of experiments:", config.TokenStorageCostExps)
	zap.L().Info("********************************************************************************************************************************\n")

	zap.L().Info("********************************************************************************************************************************\n")

}

/*
setupConfig sets up the config file, config file type and config file path
*/
func setupConfig() {
	viper.SetConfigFile("config.json") // name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
}

/*
ParseConfig reads the configuration inputs from the config.json file
and sets the configuration
*/
func ParseConfig() (Config, error) {

	zap.S().Info("Loading configuration file")

	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.json") // name of config file (without extension)
	config := Config{}
	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Fatalln("error reading config file for viper.\t", err)
	}

	config.SmartContractAddress = viper.GetString("contract.address")
	config.PrivateKey = viper.GetString("blockchain.privateKey")
	config.Account = viper.GetString("blockchain.account")
	config.GasLimit = viper.GetUint64("contract.gasLimit")
	config.GasPrice = big.NewInt(int64(viper.GetUint64("contract.gasPrice")))

	config.BlockchainRpcEndpoint = viper.GetString("blockchain.rpcEndpoint")
	config.BlockchainWebSocketEndPoint = viper.GetString("blockchain.wsEndPoint")
	config.LoggerType = viper.GetString("logger.env")
	config.LoggerOutputMode = viper.GetString("logger.output")
	config.LoggerFile = viper.GetString("logger.filename")

	config.DEBUG = viper.GetBool("mode.debug")

	config.TokenStorageNumberOfTokens = viper.GetIntSlice("token.numberOfTokensToStore")
	config.TokenStorageCostExps = viper.GetUint("token.numberOfExps")

	config.UsePreGeneratedKeysAndVCs = viper.GetBool("issuer.usePreGeneratedKeysAndVCs")

	config.SelectiveDisclosureExtension = viper.GetBool("selective_disclosure_extension")
	config.StartingEpoch = 0
	config.Params = Params{
		NumberOfExperiments:      viper.GetInt("params.number_of_experiments"),
		TotalVCs:                 viper.GetInt("params.total_vcs"),
		ExpirationPeriod:         viper.GetInt("params.expiration_period"),
		VerificationPeriod:       viper.GetInt("params.verification_period"),
		NumberOfTokensPerCircuit: viper.GetInt("params.number_of_tokens_per_circuit"),
		RevocationRateBase:       viper.GetInt("params.revocation_rate_base"),
		RevocationRateStep:       viper.GetInt("params.revocation_rate_step"),
		RevocationRateEnd:        viper.GetInt("params.revocation_rate_end"),
		EpochDuration:            viper.GetInt("params.epoch_duration"),
		IRMAWitnessUpdateMessageWithoutRepetition: viper.GetBool("params.irma_witness_update_message_without_repetition"),
	}
	config.Run = Run{
		Flow:               viper.GetBool("run.flow"),
		IRMATest:           viper.GetBool("run.irma_test"),
		ComputeFinalResult: viper.GetBool("run.compute_final_result"),
		CircuitTest:        viper.GetBool("run.circuit"),
		VCTest:             viper.GetBool("run.vc_test"),
		ZKPTest:            viper.GetBool("run.zkp_test"),
		CryptoTest:         viper.GetBool("run.crypto_test"),
		GenerateVCs:        viper.GetBool("run.generate_vcs"),
	}

	config.Benchmark = Benchmark{
		Setup:                           viper.GetBool("benchmark.setup"),
		Issaunce:                        viper.GetBool("benchmark.issuance"),
		Revocation:                      viper.GetBool("benchmark.revocation"),
		Refresh:                         viper.GetBool("benchmark.refresh"),
		Presentation_Verification:       viper.GetBool("benchmark.presentation_verification"),
		IRMASetup:                       viper.GetBool("benchmark.irma_setup"),
		IRMAIssuance:                    viper.GetBool("benchmark.irma_issuance"),
		IRMARevocation:                  viper.GetBool("benchmark.irma_revocation"),
		IRMAPresentationAndVerification: viper.GetBool("benchmark.irma_presentation_and_verification"),
		CircuitConstraints:              viper.GetBool("benchmark.circuit_constraints"),
		ListCommitment:                  viper.GetBool("benchmark.list_commitment"),
	}
	config.printConfig()
	return config, nil
}
