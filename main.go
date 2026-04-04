/*
Zap library is used for Logging. The log is stored on file or shown
in the console depending on the setting given in the config file.
*/
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/big"
	"time"
	"zkrevoke/benchmark"
	"zkrevoke/config"
	"zkrevoke/crypto2"
	"zkrevoke/engine"

	"zkrevoke/irma"
	"zkrevoke/issuer"
	"zkrevoke/model"
	"zkrevoke/zkp"
	_ "zkrevoke/zkp"
)

/*
SetupLogger sets up the logger.
The logger output mode is retrieved from the config file.
If the output mode is console then the log is shown on console
If the output mode is file then the log is stored in a file
*/
func SetupLogger(conf config.Config) {

	var filename string

	if conf.LoggerOutputMode == "console" {
		filename = "stdout"
	} else {
		r, _ := rand.Int(rand.Reader, big.NewInt(10000))
		_, month, day := time.Now().Date()
		filename = fmt.Sprintf("logs/%v_%v_%v_%v.json", conf.LoggerFile, month, day, r)
	}

	//OutputPaths: []string{"stdout"},
	//OutputPaths: []string{filename},
	zapConfig := &zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{filename},
		ErrorOutputPaths: []string{filename},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:   "level",
			MessageKey: "***",
		},
	}
	if conf.LoggerType == "dev" {
		zap.ReplaceGlobals(zap.Must(zapConfig.Build()))
	} else if conf.LoggerType == "prod" {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	} else {
		zap.ReplaceGlobals(zap.Must(zapConfig.Build()))
	}
}

func main() {

	test := flag.Bool("test", false, "enable debug mode")
	flag.Parse()
	conf, err := config.ParseConfig(test)
	conf.InitialTimestamp = time.Now()
	if err != nil {
		zap.S().Infoln("error parsing config: ", err)
	}
	SetupLogger(conf)

	startTime := time.Now()
	if conf.Run.GenerateVCs {
		issuer := issuer.NewIssuer(&conf)
		issuer.GenerateAndStoreVCs(conf)
	}

	if conf.Run.Flow {
		engine.RunSimpleScenario(conf)
	}

	if conf.Run.IRMATest {
		//irma.Test()
		irma.TestFullIssueAndShowWithRevocation()
	}
	if conf.Run.CircuitTest {
		zkp.Test_Circuit()
	}

	if conf.Run.VCTest {
		model.TestVC(conf)
	}

	if conf.Run.ZKPTest {

		zkp.TestConstraintsInCircuit()
		//zkp.Test_Circuit()

		//zkp.Test_Circuit_V4()
	}
	if conf.Run.CryptoTest {
		crypto2.Test_EDDSA()
		crypto2.Test_LoadKeys()
	}

	if conf.Benchmark.Setup {
		benchmark.Benchmark_Setup(conf)
	}
	if conf.Benchmark.Issaunce {
		benchmark.Benchmark_Issuance(conf)
	}
	if conf.Benchmark.Revocation {
		benchmark.Benchmark_Revocation(conf)
	}
	if conf.Benchmark.Refresh {
		benchmark.Benchmark_Refresh(conf)
		benchmark.Benchmark_TokenStorageCost(conf)
	}
	if conf.Benchmark.Presentation_Verification {
		benchmark.Benchmark_Presentation_Verification(conf)
	}
	if conf.Benchmark.CircuitConstraints {
		benchmark.Benchmark_CircuitConstraints(conf)
	}
	if conf.Benchmark.ListCommitment {
		benchmark.Benchmark_ListCommitment(conf)
	}
	if conf.Benchmark.IRMASetup {
		irma.Benchmark_Setup(conf)
	}
	if conf.Benchmark.IRMAIssuance {
		irma.Benchmark_Issuance(conf)
	}
	if conf.Benchmark.IRMARevocation {
		conf.Params.IRMAWitnessUpdateMessageWithoutRepetition = false
		irma.Benchmark_Revocation(conf)
		conf.Params.IRMAWitnessUpdateMessageWithoutRepetition = true
		irma.Benchmark_Revocation(conf)
	}
	if conf.Benchmark.IRMAPresentationAndVerification {
		conf.Params.IRMAWitnessUpdateMessageWithoutRepetition = false
		irma.Benchmark_Presentation_And_Verification(conf)
		conf.Params.IRMAWitnessUpdateMessageWithoutRepetition = true
		irma.Benchmark_Presentation_And_Verification(conf)
	}

	endTime := time.Since(startTime)
	if endTime.Minutes() >= 60 {
		zap.S().Infoln("Total time: ", endTime.Hours(), "hours")
	} else {
		zap.S().Infoln("Total time: ", endTime.Minutes(), "minutes")
	}

}
