package issuer

import (
	"github.com/ethereum/go-ethereum/common"
	blockchain_hardhat "zkrevoke/blockchain-hardhat"
	"zkrevoke/config"
)

/*
DeployContract deploys smart contrat and returns the smart contract address
*/

func (issuer *Issuer) DeployContract(conf config.Config) string {
	address, gasUsed, cost, _ := blockchain_hardhat.DeployContract(conf)
	issuer.Result.SetContractDeploymentCost(cost)
	issuer.Result.SetContractDeploymentGas(gasUsed)
	issuer.smartContractAddress = common.HexToAddress(address)
	return address
}
