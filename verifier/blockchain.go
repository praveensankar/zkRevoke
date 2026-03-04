package verifier

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"io"
	"time"
	"zkrevoke/blockchain-hardhat/contracts"
	"zkrevoke/zkp"
)

/*
Verifier fetches revoked tokens from the smart contract
Returns
int: epoch
[]string: revoked tokens
*/

func (verifier *Verifier) GetRevokedTokens() (int, []string) {
	client, err := ethclient.Dial(verifier.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to connect to the Ethereum client: %v", err)
	}
	//zap.S().Infoln("***VERIFIER***: smart contraact address: ", verifier.smartContractAddress)
	revocationList, err := contracts.NewRevocationList(verifier.smartContractAddress, client)

	tokenList, err := revocationList.GetTokens(nil)

	if err != nil {
		zap.S().Fatalln("***VERIFIER***: failed to retrieve revoked tokens", err)
	}

	var tokens []string

	for i := 0; i < len(tokenList.Tokens); i++ {
		tokens = append(tokens, hex.EncodeToString(tokenList.Tokens[i][:]))
	}

	return int(tokenList.Epoch.Int64()), tokens
}

func (verifier *Verifier) RetrieveHashOfCCS() bool {
	client, err := ethclient.Dial(verifier.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(verifier.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to instantiate Storage contract: %v", err)
	}

	ccs := verifier.ccs
	var buf bytes.Buffer
	_, err = ccs.WriteTo(&buf)
	// Get byte slice
	storedCCS := buf.Bytes()
	hash := sha256.Sum256(storedCCS)
	hashStr := fmt.Sprintf("%s", hash)
	ccsBytes, err := revocationList.RetrieveCCSHash(nil)
	ccsStr := fmt.Sprintf("%s", ccsBytes)
	if hashStr == ccsStr {
		zap.S().Infof("***VERIFIER***: verified validity of ccs")
		return true
	}

	return false
}

func (verifier *Verifier) RetrieveCCS() constraint.ConstraintSystem {
	client, err := ethclient.Dial(verifier.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(verifier.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to instantiate Storage contract: %v", err)
	}

	retrievedCompressedCCS, err := revocationList.RetrieveCCS(nil)
	zap.S().Infoln("***VERIFIER***: Retrieved ccs bytes: size: ", len(retrievedCompressedCCS))
	if err != nil {
		zap.S().Fatalln("***VERIFIER***: failed to retrieve CCS", err)
	}
	var retrievedDeCompressedccsBytes bytes.Buffer
	b := bytes.NewReader(retrievedCompressedCCS)
	r, err := zlib.NewReader(b)
	io.Copy(&retrievedDeCompressedccsBytes, r)
	r.Close()

	// Re-instantiate with correct curve
	retrievedccs := groth16.NewCS(ecc.BN254)
	_, err = retrievedccs.ReadFrom(bytes.NewReader(retrievedDeCompressedccsBytes.Bytes()))

	if err != nil {
		zap.S().Fatalln("***VERIFIER***: failed to retrieve CCS", err)
	}
	return retrievedccs
}

func (verifier *Verifier) RetrieveZKPVerifyingKey() groth16.VerifyingKey {
	client, err := ethclient.Dial(verifier.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(verifier.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to instantiate Storage contract: %v", err)
	}

	vkBytes, err := revocationList.RetrieveZKPVerifyingKey(nil)

	if err != nil {
		zap.S().Fatalln("***VERIFIER***: Failed to retrieve zkp verifying key: %v", err)
	}
	zap.S().Infoln("***VERIFIER***: Retrieved zkp verifying key: \t size: ", len(vkBytes))
	zkpVerifyingKey, _ := zkp.BytesToGrothVerifyingKey(vkBytes)
	return zkpVerifyingKey
}

func (verifier *Verifier) RetrieveEDDSAPublicKey() signature.PublicKey {
	client, err := ethclient.Dial(verifier.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(verifier.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to instantiate Storage contract: %v", err)
	}

	pkBytes, err := revocationList.RetrieveEDDSAPublicKey(nil)

	if err != nil {
		zap.S().Fatalln("***VERIFIER***: Failed to retrieve eddsa public key: %v", err)
	}
	//zap.S().Infoln("***VERIFIER***: Retrieved eddsa public key: \t size: ", len(pkBytes))

	publicKey := &eddsa.PublicKey{}
	// Deserialize the binary representation into the public key
	_, err = publicKey.SetBytes(pkBytes)
	if err != nil {
		zap.S().Fatalln("***VERIFIER***: Failed to set public key: %v", err)
	}

	return publicKey
}

/*
RetrieveEpochConfigurations retrieves
1) int64 - epoch duration
2) intitial timestamp
*/
func (verifier *Verifier) RetrieveEpochConfigurations() (int64, time.Time) {
	client, err := ethclient.Dial(verifier.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(verifier.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***VERIFIER***: Failed to instantiate Storage contract: %v", err)
	}

	duration, err := revocationList.RetrieveEpochDuration(nil)
	timestampBytes, err := revocationList.RetrieveInitialTimeStamp(nil)

	timeStampRetrieved := time.Time{}
	timeStampRetrieved.UnmarshalBinary(timestampBytes)
	if err != nil {
		zap.S().Fatalln("***VERIFIER***: Failed to retrieve epoch configuration: %v", err)
	}
	zap.S().Infoln("***VERIFIER***: Retrieved epoch configurations: ",
		"\t epoch duration: ", duration.Int64(),
		"\t initial timestamp: ", timeStampRetrieved)

	return duration.Int64(), timeStampRetrieved
}
