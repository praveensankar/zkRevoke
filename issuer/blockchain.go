package issuer

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
	"strconv"
	"time"
	"zkrevoke/blockchain-hardhat/contracts"
	"zkrevoke/utils"
	"zkrevoke/zkp"
)

func (issuer Issuer) getAuth() *bind.TransactOpts {
	// step 1: connect to a blockchain node using RPC endpoint
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(issuer.privateKey)
	if err != nil {
		zap.S().Fatalln("ISSUER: auth error: ", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		zap.S().Fatalln("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//zap.S().Infof("***ISSUER***: \t address: %s", fromAddress.String())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		zap.S().Fatalln(err)
	}

	chainID, err := client.ChainID(context.Background())

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = issuer.gasLimit
	//auth.GasPrice = issuer.gasPrice
	//auth.GasPrice, _ = client.SuggestGasPrice(context.Background())

	return auth
}

/*
Issues publishes revoked tokens at the start of each epoch to smart contract
returns
1) uint64 - gas used (reported from the ganache)
2) string - cost to store tokens (calculated the difference in account balance)
3) int - token size in bytes
4) int - number of transactions published
5) int - number of revoked tokens published
*/

func (issuer *Issuer) PublishRevokedTokensToBlockchain() (uint64, string, int, int, int) {
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}

	epoch := big.NewInt(int64(issuer.GetCurrentEpoch()))

	tokens := issuer.GetRevokedTokens()

	cost := big.NewInt(0)
	gasUsed := uint64(0)
	txCount := 0
	done := false

	for i := 0; i < len(tokens); i++ {
		var tokensBytes [][32]byte
		if done == true {
			break
		}
		numberOfTokensInaBlock := 1000
		for j := i * numberOfTokensInaBlock; j < (i+1)*numberOfTokensInaBlock; j++ {
			if j == len(tokens) {
				done = true
				break
			}
			if j == (i+1)*numberOfTokensInaBlock {
				break
			}
			byteRepr := [32]byte{}
			copy(byteRepr[:], tokens[j][:])
			tokensBytes = append(tokensBytes, byteRepr)
		}
		auth := issuer.getAuth()
		startBalance, _ := client.BalanceAt(context.Background(), issuer.account, nil)
		_, err = revocationList.RefreshRevokedTokens(auth, epoch, tokensBytes)
		endBalance, _ := client.BalanceAt(context.Background(), issuer.account, nil)
		txCount = txCount + 1
		header, err := client.HeaderByNumber(context.Background(), nil)
		co := new(big.Int).Sub(startBalance, endBalance)
		cost = cost.Add(cost, co)
		gasUsed = gasUsed + header.GasUsed
		if err != nil {
			zap.S().Fatalln("***ISSUER***: failed to publish revoked tokens", err)
		}
	}

	zap.S().Infoln("***ISSUER***: Published revoked tokens:  ", len(tokens), "\t cost in wei: ", cost.String(),
		"\t gas used: ", gasUsed)
	return gasUsed, cost.String(), len(tokens[0]), txCount, len(tokens)
}

/*
PublishHashOfCCS cryptography hash of the groth16 constraint system
Returns:
 1. uint64 - gas used (reported from the ganache)
 2. float64 - cost to store hash of ccs
*/
func (issuer *Issuer) PublishHashOfCCS() (uint64, string) {
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}
	auth := issuer.getAuth()

	ccs := issuer.ccs

	var buf bytes.Buffer
	_, err = ccs.WriteTo(&buf)
	if err != nil {
		// handle error
	}

	// Get byte slice
	ccsBytes := buf.Bytes()
	hash := sha256.Sum256(ccsBytes)

	var compressedCCSBytes bytes.Buffer
	w := zlib.NewWriter(&compressedCCSBytes)
	w.Write(ccsBytes)
	w.Close()
	zap.S().Infoln("***ISSUER***: ccs size (in KB): ", len(ccsBytes)/1000,
		"\t compressed ccs size (in KB): ", compressedCCSBytes.Len()/1000)

	startBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	_, err = revocationList.PublishCCSHash(auth, hash[:])
	endBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	cost := new(big.Int).Sub(startBalance, endBalance).String()
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish ccs hash", err)
	}

	hashStr := fmt.Sprintf("%s", hash)
	ccsBytes1, err := revocationList.RetrieveCCSHash(nil)
	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to fetch ccs hash", err)
	}
	ccsStr := fmt.Sprintf("%s", ccsBytes1)
	if hashStr == ccsStr {
		zap.S().Infof("***ISSUER***: verified validity of ccs")
	}

	zap.S().Infoln("***ISSUER***: Published hash of ccs used in groth16 zkp scheme:  \t cost: ", cost)
	return header.GasUsed, cost
}

/*
PublishCCS publishes  groth16 constraint system
Returns:
1) uint64 - amount of gas used
2) string - cost to deploy the contract in wei
*/
func (issuer *Issuer) PublishCCS() (uint64, string) {
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}

	ccs := issuer.ccs

	var buf bytes.Buffer
	_, err = ccs.WriteTo(&buf)
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
	zap.S().Infoln("***ISSUER***: ccs size (in KB): ", len(ccsBytes)/1000,
		"\t compressed ccs size (in KB): ", len(compressedCCSBytes)/1000)

	txCount := 0
	numberOfBytesInaTx := 1000
	done := false
	cost := big.NewInt(0)
	gasUsed := uint64(0)
	for i := 0; i < len(compressedCCSBytes); i++ {

		if done == true {
			break
		}
		j := i * numberOfBytesInaTx
		for ; j < (i+1)*numberOfBytesInaTx; j++ {
			if j == len(compressedCCSBytes) {
				done = true
				break
			}
			if j == (i+1)*numberOfBytesInaTx {
				break
			}
		}
		ccsSlice := compressedCCSBytes[i*numberOfBytesInaTx : j]
		auth := issuer.getAuth()
		startBalance, _ := client.BalanceAt(context.Background(), issuer.account, nil)
		txCount = txCount + 1
		_, err = revocationList.PublishCCS(auth, ccsSlice[:])
		if err != nil {
			zap.S().Fatalln("***ISSUER***: failed to publish ccs", err)
		}
		endBalance, _ := client.BalanceAt(context.Background(), issuer.account, nil)
		header, err := client.HeaderByNumber(context.Background(), nil)
		co := new(big.Int).Sub(startBalance, endBalance)
		cost = cost.Add(cost, co)
		gasUsed = gasUsed + header.GasUsed

		if err != nil {
			zap.S().Fatalln("***HOLDER***: failed to retrieve CCS", err)
		}
	}
	//retrievedCompressedCCS, err := revocationList.RetrieveCCS(nil)
	//zap.S().Infoln("***ISSUER***: Retrieved ccs bytes: size: ", len(retrievedCompressedCCS))
	//if err != nil {
	//	zap.S().Fatalln("***HOLDER***: failed to retrieve CCS", err)
	//}
	//var retrievedDeCompressedccsBytes bytes.Buffer
	//b := bytes.NewReader(retrievedCompressedCCS)
	//r, err := zlib.NewReader(b)
	//io.Copy(&retrievedDeCompressedccsBytes, r)
	//r.Close()
	//
	//// Re-instantiate with correct curve
	//retrievedccs := groth16.NewCS(ecc.BN254)
	//_, err = retrievedccs.ReadFrom(bytes.NewReader(retrievedDeCompressedccsBytes.Bytes()))
	//
	//if err != nil {
	//	zap.S().Fatalln("***ISSUER***: failed to retrieve CCS", err)
	//}
	zap.S().Infoln("***ISSUER***: Published ccs used in groth16 zkp scheme: \t number of transactions: ", txCount,
		"\t gas used: ", gasUsed, "\t cost: ", cost.String())

	return gasUsed, cost.String()
}

/*
PublishZKPVerifyingKey publishes zkp verifying key to smart contract
Returns:
1) uint64 - gas used (reported from the ganache)
2) float64 - cost to store zkp public key (calculated the difference in account balance)
*/
func (issuer *Issuer) PublishZKPVerifyingKey() (uint64, string) {
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}
	auth := issuer.getAuth()

	zkpVerifingKey := issuer.zkpVerifyingKey
	zkpVerifingKeyBytes, _ := zkp.GrothVerifyingKeyToBytes(zkpVerifingKey)

	startBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	_, err = revocationList.PublishZKPVerifyingKey(auth, zkpVerifingKeyBytes)

	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish zkp verifying key:", err)
	}
	endBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	header, err := client.HeaderByNumber(context.Background(), nil)
	cost := new(big.Int).Sub(startBalance, endBalance).String()
	zap.S().Infoln("***ISSUER***: Published zkp verifying key: size: ", len(zkpVerifingKeyBytes),
		"\t gas used: ", header.GasUsed,
		"\t cost: ", cost)

	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish zkp verifying key", err)
	}

	//vk, err := revocationList.RetrieveZKPVerifyingKey(nil)
	//
	//if err != nil {
	//	zap.S().Fatalln("***ISSUER***: Failed to retrieve zkp verifying key: %v", err)
	//}
	//zap.S().Infoln("***ISSUER***: Retrieved zkp verifying key: \t size: ", len(vk))

	return header.GasUsed, cost
}

/*
PublishEDDSAPublicKey publishes eddsa public key to smart contract
Returns:
1) uint64 - gas used (reported from the ganache)
2) float64 - cost to store eddsa public key (calculated the difference in account balance)
*/
func (issuer *Issuer) PublishEDDSAPublicKey() (uint64, string) {
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}
	auth := issuer.getAuth()

	publicKey := issuer.eddsaPublicKey

	startBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	_, err = revocationList.PublishEDDSAPublicKey(auth, publicKey.Bytes())

	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish  eddsa public key:", err)
	}
	endBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	header, err := client.HeaderByNumber(context.Background(), nil)
	cost := new(big.Int).Sub(startBalance, endBalance).String()
	zap.S().Infoln("***ISSUER***: Published eddsa public key: size (bytes): ", len(publicKey.Bytes()),
		"\t gas used: ", header.GasUsed,
		"\t cost: ", cost)

	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish eddsa public key", err)
	}

	pkBytes, err := revocationList.RetrieveEDDSAPublicKey(nil)

	if err != nil {
		zap.S().Fatalln("***ISSUER***: Failed to retrieve eddsa public key: %v", err)
	}
	zap.S().Infoln("***ISSUER***: Retrieved eddsa public key: \t size: ", len(pkBytes))
	pkRes := &eddsa.PublicKey{}
	// Deserialize the binary representation into the public key
	_, err = pkRes.SetBytes(pkBytes)
	if err != nil {
		zap.S().Fatalln("***ISSUER***: Failed to set public key: %v", err)
	}
	if pkRes.Equal(issuer.eddsaPublicKey) {
		zap.S().Infoln("***ISSUER***: EdDSA public key at the smart contract is the same as the public key")
	}
	return header.GasUsed, cost
}

/*
PublishEpochConfigurations publishes epoch duration and intial timestamp to smart contract
Returns:
1) uint64 - gas used (reported from the ganache)
2) float64 - cost to store epoch duration and initial timestamp (calculated the difference in account balance)
*/
func (issuer *Issuer) PublishEpochConfigurations() (uint64, string) {
	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}
	auth := issuer.getAuth()

	epochDuration := issuer.Duration
	initialTimestamp, _ := issuer.InitialTimeStamp.MarshalBinary()

	startBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	_, err = revocationList.PublishEpochConfigurations(auth, big.NewInt(int64(epochDuration)), initialTimestamp)

	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish epoch configurations:", err)
	}
	endBalance, err := client.BalanceAt(context.Background(), issuer.account, nil)
	header, err := client.HeaderByNumber(context.Background(), nil)
	cost := new(big.Int).Sub(startBalance, endBalance).String()
	zap.S().Infoln("***ISSUER***: Published epoch duration and initial timestamp: ",
		"\t epoch duration: ", epochDuration,
		"\t intial timestamp: ", issuer.InitialTimeStamp,
		"\t gas used: ", header.GasUsed,
		"\t cost: ", cost)

	if err != nil {
		zap.S().Fatalln("***ISSUER***: failed to publish eddsa public key", err)
	}
	//
	//duration, err := revocationList.RetrieveEpochDuration(nil)
	//timestampBytes, err := revocationList.RetrieveInitialTimeStamp(nil)
	//
	//timeStampRetrieved := time.Time{}
	//timeStampRetrieved.UnmarshalBinary(timestampBytes)
	//if err != nil {
	//	zap.S().Fatalln("***ISSUER***: Failed to retrieve eddsa public key: %v", err)
	//}
	//zap.S().Infoln("***ISSUER***: Retrieved epoch configurations: ",
	//	"\t epoch duration: ", duration.Int64(),
	//	"\t initial timestamp: ", timeStampRetrieved)

	return header.GasUsed, cost
}

/*
Issues computes the cost to store tokens by publishing random tokens
returns
1) uint64 - gas used (reported from the ganache)
2) *big.Int - cost to store tokens (calculated the difference in account balance)
3) int - token size in bytes
4) int - number of transactions published
5) time.Duration - time to refresh the list locally
*/
func (issuer *Issuer) ComputeTokenStorageCost(count int, epoch int) (uint64, *big.Int, int, int, time.Duration) {

	client, err := ethclient.Dial(issuer.blockchainRPCEndpoint)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to connect to the Ethereum client: %v", err)
	}
	revocationList, err := contracts.NewRevocationList(issuer.smartContractAddress, client)
	if err != nil {
		zap.S().Infof("***ISSUER***: Failed to instantiate Storage contract: %v", err)
	}

	// generate vc IDs
	var vcIDs []string
	counter := 1
	for i := 0; i < count; i++ {
		vcId := fmt.Sprintf("vcID#%s", strconv.Itoa(counter))
		counter = counter + 1
		vcIDs = append(vcIDs, vcId)
		seed := rand.Text()
		issuer.seedStore[vcId] = seed
	}
	var tokens [][]byte

	startListComputation := time.Now()
	for i := 0; i < count; i++ {
		seed := issuer.seedStore[vcIDs[i]]
		token := utils.ComputeToken(epoch, seed)
		tokens = append(tokens, token)
	}
	endListComputation := time.Since(startListComputation)

	cost := big.NewInt(0)
	gasUsed := uint64(0)
	txCount := 0
	done := false

	for i := 0; i < len(tokens); i++ {
		var tokensBytes [][32]byte
		if done == true {
			break
		}
		numberOfTokensInaBlock := 1000
		for j := i * numberOfTokensInaBlock; j < (i+1)*numberOfTokensInaBlock; j++ {
			if j == len(tokens) {
				done = true
				break
			}
			if j == (i+1)*numberOfTokensInaBlock {
				break
			}
			byteRepr := [32]byte{}
			copy(byteRepr[:], tokens[j][:])
			tokensBytes = append(tokensBytes, byteRepr)
		}
		auth := issuer.getAuth()
		startBalance, _ := client.BalanceAt(context.Background(), issuer.account, nil)
		_, err = revocationList.RefreshRevokedTokens(auth, big.NewInt(int64(epoch)), tokensBytes)
		endBalance, _ := client.BalanceAt(context.Background(), issuer.account, nil)
		txCount = txCount + 1
		header, err := client.HeaderByNumber(context.Background(), nil)
		co := new(big.Int).Sub(startBalance, endBalance)
		cost = cost.Add(cost, co)
		gasUsed = gasUsed + header.GasUsed
		if err != nil {
			zap.S().Fatalln("***ISSUER***: failed to publish revoked tokens", err)
		}
	}

	zap.S().Infoln("***ISSUER***: Published revoked tokens:  ", len(tokens), "\t cost in wei: ", cost.String(),
		"\t gas used: ", gasUsed)
	return gasUsed, cost, len(tokens[0]), txCount, endListComputation
}
