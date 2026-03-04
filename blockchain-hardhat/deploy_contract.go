package blockchain_hardhat

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
	"zkrevoke/blockchain-hardhat/contracts"
	"zkrevoke/config"
)

/*
DeployContract deploys smart contrat and returns
1) string - deployed contract address
2) uint64 - gas used
3) int64 - cost to deploy the contract in wei
*/
func DeployContract(config config.Config) (string, uint64, string, error) {
	client, err := ethclient.Dial(config.BlockchainRpcEndpoint)
	if err != nil {
		zap.S().Fatalln("ERROR in deploying contract", err)
	}

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		zap.S().Fatalln(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		zap.S().Fatalln("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		zap.S().Fatalln(err)
	}

	gasLimit := config.GasLimit

	chainID, err := client.ChainID(context.Background())
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit

	account := common.HexToAddress(config.Account)
	startBalance, err := client.BalanceAt(context.Background(), account, nil)
	addresss, tx, _, err := contracts.DeployRevocationList(auth, client)
	header, err := client.HeaderByNumber(context.Background(), nil)
	endBalance, err := client.BalanceAt(context.Background(), account, nil)
	cost := new(big.Int).Sub(startBalance, endBalance).String()

	if err != nil {
		zap.S().Infof("Failed to deploy contract: %v", err)
	}

	if config.DEBUG == true {
		zap.L().Info("\n\n------------------------------------------------------- deploying smart contract --------------------------------------------------")
		zap.S().Infoln("BLOCKCHAIN - \t chain id: ", chainID)
		zap.S().Infoln("BLOCKCHAIN- \t  smart contract address: ", addresss.String())
		zap.S().Infoln("BLOCKCHAIN - \t tx gas price: ", tx.GasPrice().Int64())
		zap.L().Info("********************************************************************************************************************************\n")
	}
	return addresss.String(), header.GasUsed, cost, err
}
