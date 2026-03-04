// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RevocationListList is an auto generated low-level Go binding around an user-defined struct.
type RevocationListList struct {
	Epoch  *big.Int
	Tokens [][32]byte
}

// RevocationListMetaData contains all meta data concerning the RevocationList contract.
var RevocationListMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EDDSAPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EpochDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetTokens\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"tokens\",\"type\":\"bytes32[]\"}],\"internalType\":\"structRevocationList.List\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"InitialTimeStamp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bbsPublicKey\",\"type\":\"bytes\"}],\"name\":\"PublishBBSPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_publicKey\",\"type\":\"bytes\"}],\"name\":\"PublishEDDSAPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch_duration\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"intial_timestamp\",\"type\":\"bytes\"}],\"name\":\"PublishEpochConfigurations\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_verifyingKey\",\"type\":\"bytes\"}],\"name\":\"PublishZKPVerifyingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"_tokens\",\"type\":\"bytes32[]\"}],\"name\":\"RefreshRevokedTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveBBSPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveCCS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveCCSHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveEDDSAPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveEpochDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveInitialTimeStamp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RetrieveZKPVerifyingKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZKPVerifyingKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bbsPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ccs\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ccsHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"list\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_ccs\",\"type\":\"bytes\"}],\"name\":\"publishCCS\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_ccsHash\",\"type\":\"bytes\"}],\"name\":\"publishCCSHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registerIssuers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"test\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b503360025f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611a888061005c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610171575f3560e01c806367ebfdd5116100dc578063bd03b11311610095578063d03c7e661161006f578063d03c7e66146103c9578063da191722146103e7578063e6ba9a3c14610405578063f8a8fd6d1461042357610171565b8063bd03b11314610373578063c0de5cae1461038f578063c6a6c0e1146103ab57610171565b806367ebfdd5146102c15780638a8d3bde146102df57806395681733146102fd578063b5b1eadd1461031b578063b5fd3c5c14610339578063b8f1b1741461035557610171565b806339fcc15e1161012e57806339fcc15e146102135780634164db7d1461022f57806341ebf5a31461024d5780634c59672e146102695780635a9c658c146102875780635cbb46f6146102a357610171565b80630de54b85146101755780630f560cd71461017f5780631b717f041461019d5780632723e24b146101bb5780632b5d9cbe146101d957806330cc2120146101f5575b5f5ffd5b61017d61042d565b005b61018761042f565b6040516101949190611237565b60405180910390f35b6101a5610439565b6040516101b291906112c0565b60405180910390f35b6101c36104c9565b6040516101d091906112c0565b60405180910390f35b6101f360048036038101906101ee919061141d565b610559565b005b6101fd6105c4565b60405161020a919061156d565b60405180910390f35b61022d6004803603810190610228919061141d565b61063c565b005b61023761076c565b60405161024491906112c0565b60405180910390f35b6102676004803603810190610262919061141d565b6107fc565b005b61027161092c565b60405161027e91906112c0565b60405180910390f35b6102a1600480360381019061029c91906116a5565b6109bc565b005b6102ab610add565b6040516102b891906112c0565b60405180910390f35b6102c9610b69565b6040516102d691906112c0565b60405180910390f35b6102e7610bf5565b6040516102f491906112c0565b60405180910390f35b610305610c81565b6040516103129190611237565b60405180910390f35b610323610c87565b60405161033091906112c0565b60405180910390f35b610353600480360381019061034e919061141d565b610d17565b005b61035d610d82565b60405161036a91906112c0565b60405180910390f35b61038d600480360381019061038891906116ff565b610e0e565b005b6103a960048036038101906103a4919061141d565b610e81565b005b6103b3610eec565b6040516103c09190611237565b60405180910390f35b6103d1610ef5565b6040516103de91906112c0565b60405180910390f35b6103ef610f81565b6040516103fc91906112c0565b60405180910390f35b61040d61100d565b60405161041a91906112c0565b60405180910390f35b61042b61109d565b005b565b5f805f0154905081565b60606006805461044890611786565b80601f016020809104026020016040519081016040528092919081815260200182805461047490611786565b80156104bf5780601f10610496576101008083540402835291602001916104bf565b820191905f5260205f20905b8154815290600101906020018083116104a257829003601f168201915b5050505050905090565b6060600380546104d890611786565b80601f016020809104026020016040519081016040528092919081815260200182805461050490611786565b801561054f5780601f106105265761010080835404028352916020019161054f565b820191905f5260205f20905b81548152906001019060200180831161053257829003601f168201915b5050505050905090565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105b1575f5ffd5b80600590816105c09190611956565b5050565b6105cc611182565b5f6040518060400160405290815f82015481526020016001820180548060200260200160405190810160405280929190818152602001828054801561062e57602002820191905f5260205f20905b81548152602001906001019080831161061a575b505050505081525050905090565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610694575f5ffd5b5f5f90505b81518110156107685760048282815181106106b7576106b6611a25565b5b602001015160f81c60f81b90808054806106d090611786565b80601f81036106ed57835f5260205f2060ff1984168155603f9350505b5060028201835560018101925050506001900381546001161561071d57905f5260205f2090602091828204019190065b909190919091601f036101000a81548160ff021916907f0100000000000000000000000000000000000000000000000000000000000000840402179055508080600101915050610699565b5050565b60606008805461077b90611786565b80601f01602080910402602001604051908101604052809291908181526020018280546107a790611786565b80156107f25780601f106107c9576101008083540402835291602001916107f2565b820191905f5260205f20905b8154815290600101906020018083116107d557829003601f168201915b5050505050905090565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610854575f5ffd5b5f5f90505b815181101561092857600382828151811061087757610876611a25565b5b602001015160f81c60f81b908080548061089090611786565b80601f81036108ad57835f5260205f2060ff1984168155603f9350505b506002820183556001810192505050600190038154600116156108dd57905f5260205f2090602091828204019190065b909190919091601f036101000a81548160ff021916907f0100000000000000000000000000000000000000000000000000000000000000840402179055508080600101915050610859565b5050565b60606007805461093b90611786565b80601f016020809104026020016040519081016040528092919081815260200182805461096790611786565b80156109b25780601f10610989576101008083540402835291602001916109b2565b820191905f5260205f20905b81548152906001019060200180831161099557829003601f168201915b5050505050905090565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a14575f5ffd5b815f5f015403610a81575f5f90505b8151811015610a7b575f600101828281518110610a4357610a42611a25565b5b6020026020010151908060018154018082558091505060019003905f5260205f20015f90919091909150558080600101915050610a23565b50610ad9565b5f5f5f82015f9055600182015f610a98919061119b565b50506040518060400160405280838152602001828152505f5f820151815f01556020820151816001019080519060200190610ad49291906111b9565b509050505b5050565b60058054610aea90611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610b1690611786565b8015610b615780601f10610b3857610100808354040283529160200191610b61565b820191905f5260205f20905b815481529060010190602001808311610b4457829003601f168201915b505050505081565b60048054610b7690611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610ba290611786565b8015610bed5780601f10610bc457610100808354040283529160200191610bed565b820191905f5260205f20905b815481529060010190602001808311610bd057829003601f168201915b505050505081565b60088054610c0290611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610c2e90611786565b8015610c795780601f10610c5057610100808354040283529160200191610c79565b820191905f5260205f20905b815481529060010190602001808311610c5c57829003601f168201915b505050505081565b60095481565b606060048054610c9690611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610cc290611786565b8015610d0d5780601f10610ce457610100808354040283529160200191610d0d565b820191905f5260205f20905b815481529060010190602001808311610cf057829003601f168201915b5050505050905090565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610d6f575f5ffd5b8060079081610d7e9190611956565b5050565b60038054610d8f90611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610dbb90611786565b8015610e065780601f10610ddd57610100808354040283529160200191610e06565b820191905f5260205f20905b815481529060010190602001808311610de957829003601f168201915b505050505081565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610e66575f5ffd5b816009819055508060089081610e7c9190611956565b505050565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610ed9575f5ffd5b8060069081610ee89190611956565b5050565b5f600954905090565b60068054610f0290611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610f2e90611786565b8015610f795780601f10610f5057610100808354040283529160200191610f79565b820191905f5260205f20905b815481529060010190602001808311610f5c57829003601f168201915b505050505081565b60078054610f8e90611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610fba90611786565b80156110055780601f10610fdc57610100808354040283529160200191611005565b820191905f5260205f20905b815481529060010190602001808311610fe857829003601f168201915b505050505081565b60606005805461101c90611786565b80601f016020809104026020016040519081016040528092919081815260200182805461104890611786565b80156110935780601f1061106a57610100808354040283529160200191611093565b820191905f5260205f20905b81548152906001019060200180831161107657829003601f168201915b5050505050905090565b5f600267ffffffffffffffff8111156110b9576110b86112f9565b5b6040519080825280602002602001820160405280156110e75781602001602082028036833780820191505090505b5090507f68656c6c6f0000000000000000000000000000000000000000000000000000005f1b815f815181106111205761111f611a25565b5b6020026020010181815250507f68656c6c6f0000000000000000000000000000000000000000000000000000005f1b8160018151811061116357611162611a25565b5b6020026020010181815250505f6001905061117e81836109bc565b5050565b60405180604001604052805f8152602001606081525090565b5080545f8255905f5260205f20908101906111b69190611204565b50565b828054828255905f5260205f209081019282156111f3579160200282015b828111156111f25782518255916020019190600101906111d7565b5b5090506112009190611204565b5090565b5b8082111561121b575f815f905550600101611205565b5090565b5f819050919050565b6112318161121f565b82525050565b5f60208201905061124a5f830184611228565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61129282611250565b61129c818561125a565b93506112ac81856020860161126a565b6112b581611278565b840191505092915050565b5f6020820190508181035f8301526112d88184611288565b905092915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61132f82611278565b810181811067ffffffffffffffff8211171561134e5761134d6112f9565b5b80604052505050565b5f6113606112e0565b905061136c8282611326565b919050565b5f67ffffffffffffffff82111561138b5761138a6112f9565b5b61139482611278565b9050602081019050919050565b828183375f83830152505050565b5f6113c16113bc84611371565b611357565b9050828152602081018484840111156113dd576113dc6112f5565b5b6113e88482856113a1565b509392505050565b5f82601f830112611404576114036112f1565b5b81356114148482602086016113af565b91505092915050565b5f60208284031215611432576114316112e9565b5b5f82013567ffffffffffffffff81111561144f5761144e6112ed565b5b61145b848285016113f0565b91505092915050565b61146d8161121f565b82525050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f819050919050565b6114ae8161149c565b82525050565b5f6114bf83836114a5565b60208301905092915050565b5f602082019050919050565b5f6114e182611473565b6114eb818561147d565b93506114f68361148d565b805f5b8381101561152657815161150d88826114b4565b9750611518836114cb565b9250506001810190506114f9565b5085935050505092915050565b5f604083015f8301516115485f860182611464565b506020830151848203602086015261156082826114d7565b9150508091505092915050565b5f6020820190508181035f8301526115858184611533565b905092915050565b6115968161121f565b81146115a0575f5ffd5b50565b5f813590506115b18161158d565b92915050565b5f67ffffffffffffffff8211156115d1576115d06112f9565b5b602082029050602081019050919050565b5f5ffd5b6115ef8161149c565b81146115f9575f5ffd5b50565b5f8135905061160a816115e6565b92915050565b5f61162261161d846115b7565b611357565b90508083825260208201905060208402830185811115611645576116446115e2565b5b835b8181101561166e578061165a88826115fc565b845260208401935050602081019050611647565b5050509392505050565b5f82601f83011261168c5761168b6112f1565b5b813561169c848260208601611610565b91505092915050565b5f5f604083850312156116bb576116ba6112e9565b5b5f6116c8858286016115a3565b925050602083013567ffffffffffffffff8111156116e9576116e86112ed565b5b6116f585828601611678565b9150509250929050565b5f5f60408385031215611715576117146112e9565b5b5f611722858286016115a3565b925050602083013567ffffffffffffffff811115611743576117426112ed565b5b61174f858286016113f0565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061179d57607f821691505b6020821081036117b0576117af611759565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026118127fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826117d7565b61181c86836117d7565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61185761185261184d8461121f565b611834565b61121f565b9050919050565b5f819050919050565b6118708361183d565b61188461187c8261185e565b8484546117e3565b825550505050565b5f5f905090565b61189b61188c565b6118a6818484611867565b505050565b5b818110156118c9576118be5f82611893565b6001810190506118ac565b5050565b601f82111561190e576118df816117b6565b6118e8846117c8565b810160208510156118f7578190505b61190b611903856117c8565b8301826118ab565b50505b505050565b5f82821c905092915050565b5f61192e5f1984600802611913565b1980831691505092915050565b5f611946838361191f565b9150826002028217905092915050565b61195f82611250565b67ffffffffffffffff811115611978576119776112f9565b5b6119828254611786565b61198d8282856118cd565b5f60209050601f8311600181146119be575f84156119ac578287015190505b6119b6858261193b565b865550611a1d565b601f1984166119cc866117b6565b5f5b828110156119f3578489015182556001820191506020850194506020810190506119ce565b86831015611a105784890151611a0c601f89168261191f565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffdfea264697066735822122098678e864f9046ee1965dfaaf9483692751bcfc92e4a7260be5256badb56da9564736f6c634300081d0033",
}

// RevocationListABI is the input ABI used to generate the binding from.
// Deprecated: Use RevocationListMetaData.ABI instead.
var RevocationListABI = RevocationListMetaData.ABI

// RevocationListBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RevocationListMetaData.Bin instead.
var RevocationListBin = RevocationListMetaData.Bin

// DeployRevocationList deploys a new Ethereum contract, binding an instance of RevocationList to it.
func DeployRevocationList(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RevocationList, error) {
	parsed, err := RevocationListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RevocationListBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RevocationList{RevocationListCaller: RevocationListCaller{contract: contract}, RevocationListTransactor: RevocationListTransactor{contract: contract}, RevocationListFilterer: RevocationListFilterer{contract: contract}}, nil
}

// RevocationList is an auto generated Go binding around an Ethereum contract.
type RevocationList struct {
	RevocationListCaller     // Read-only binding to the contract
	RevocationListTransactor // Write-only binding to the contract
	RevocationListFilterer   // Log filterer for contract events
}

// RevocationListCaller is an auto generated read-only Go binding around an Ethereum contract.
type RevocationListCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevocationListTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RevocationListTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevocationListFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RevocationListFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevocationListSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RevocationListSession struct {
	Contract     *RevocationList   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RevocationListCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RevocationListCallerSession struct {
	Contract *RevocationListCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RevocationListTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RevocationListTransactorSession struct {
	Contract     *RevocationListTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RevocationListRaw is an auto generated low-level Go binding around an Ethereum contract.
type RevocationListRaw struct {
	Contract *RevocationList // Generic contract binding to access the raw methods on
}

// RevocationListCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RevocationListCallerRaw struct {
	Contract *RevocationListCaller // Generic read-only contract binding to access the raw methods on
}

// RevocationListTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RevocationListTransactorRaw struct {
	Contract *RevocationListTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRevocationList creates a new instance of RevocationList, bound to a specific deployed contract.
func NewRevocationList(address common.Address, backend bind.ContractBackend) (*RevocationList, error) {
	contract, err := bindRevocationList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RevocationList{RevocationListCaller: RevocationListCaller{contract: contract}, RevocationListTransactor: RevocationListTransactor{contract: contract}, RevocationListFilterer: RevocationListFilterer{contract: contract}}, nil
}

// NewRevocationListCaller creates a new read-only instance of RevocationList, bound to a specific deployed contract.
func NewRevocationListCaller(address common.Address, caller bind.ContractCaller) (*RevocationListCaller, error) {
	contract, err := bindRevocationList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RevocationListCaller{contract: contract}, nil
}

// NewRevocationListTransactor creates a new write-only instance of RevocationList, bound to a specific deployed contract.
func NewRevocationListTransactor(address common.Address, transactor bind.ContractTransactor) (*RevocationListTransactor, error) {
	contract, err := bindRevocationList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RevocationListTransactor{contract: contract}, nil
}

// NewRevocationListFilterer creates a new log filterer instance of RevocationList, bound to a specific deployed contract.
func NewRevocationListFilterer(address common.Address, filterer bind.ContractFilterer) (*RevocationListFilterer, error) {
	contract, err := bindRevocationList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RevocationListFilterer{contract: contract}, nil
}

// bindRevocationList binds a generic wrapper to an already deployed contract.
func bindRevocationList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RevocationListMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RevocationList *RevocationListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RevocationList.Contract.RevocationListCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RevocationList *RevocationListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevocationList.Contract.RevocationListTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RevocationList *RevocationListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RevocationList.Contract.RevocationListTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RevocationList *RevocationListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RevocationList.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RevocationList *RevocationListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevocationList.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RevocationList *RevocationListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RevocationList.Contract.contract.Transact(opts, method, params...)
}

// EDDSAPublicKey is a free data retrieval call binding the contract method 0xda191722.
//
// Solidity: function EDDSAPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCaller) EDDSAPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "EDDSAPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EDDSAPublicKey is a free data retrieval call binding the contract method 0xda191722.
//
// Solidity: function EDDSAPublicKey() view returns(bytes)
func (_RevocationList *RevocationListSession) EDDSAPublicKey() ([]byte, error) {
	return _RevocationList.Contract.EDDSAPublicKey(&_RevocationList.CallOpts)
}

// EDDSAPublicKey is a free data retrieval call binding the contract method 0xda191722.
//
// Solidity: function EDDSAPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) EDDSAPublicKey() ([]byte, error) {
	return _RevocationList.Contract.EDDSAPublicKey(&_RevocationList.CallOpts)
}

// EpochDuration is a free data retrieval call binding the contract method 0x95681733.
//
// Solidity: function EpochDuration() view returns(uint256)
func (_RevocationList *RevocationListCaller) EpochDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "EpochDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochDuration is a free data retrieval call binding the contract method 0x95681733.
//
// Solidity: function EpochDuration() view returns(uint256)
func (_RevocationList *RevocationListSession) EpochDuration() (*big.Int, error) {
	return _RevocationList.Contract.EpochDuration(&_RevocationList.CallOpts)
}

// EpochDuration is a free data retrieval call binding the contract method 0x95681733.
//
// Solidity: function EpochDuration() view returns(uint256)
func (_RevocationList *RevocationListCallerSession) EpochDuration() (*big.Int, error) {
	return _RevocationList.Contract.EpochDuration(&_RevocationList.CallOpts)
}

// GetTokens is a free data retrieval call binding the contract method 0x30cc2120.
//
// Solidity: function GetTokens() view returns((uint256,bytes32[]))
func (_RevocationList *RevocationListCaller) GetTokens(opts *bind.CallOpts) (RevocationListList, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "GetTokens")

	if err != nil {
		return *new(RevocationListList), err
	}

	out0 := *abi.ConvertType(out[0], new(RevocationListList)).(*RevocationListList)

	return out0, err

}

// GetTokens is a free data retrieval call binding the contract method 0x30cc2120.
//
// Solidity: function GetTokens() view returns((uint256,bytes32[]))
func (_RevocationList *RevocationListSession) GetTokens() (RevocationListList, error) {
	return _RevocationList.Contract.GetTokens(&_RevocationList.CallOpts)
}

// GetTokens is a free data retrieval call binding the contract method 0x30cc2120.
//
// Solidity: function GetTokens() view returns((uint256,bytes32[]))
func (_RevocationList *RevocationListCallerSession) GetTokens() (RevocationListList, error) {
	return _RevocationList.Contract.GetTokens(&_RevocationList.CallOpts)
}

// InitialTimeStamp is a free data retrieval call binding the contract method 0x8a8d3bde.
//
// Solidity: function InitialTimeStamp() view returns(bytes)
func (_RevocationList *RevocationListCaller) InitialTimeStamp(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "InitialTimeStamp")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// InitialTimeStamp is a free data retrieval call binding the contract method 0x8a8d3bde.
//
// Solidity: function InitialTimeStamp() view returns(bytes)
func (_RevocationList *RevocationListSession) InitialTimeStamp() ([]byte, error) {
	return _RevocationList.Contract.InitialTimeStamp(&_RevocationList.CallOpts)
}

// InitialTimeStamp is a free data retrieval call binding the contract method 0x8a8d3bde.
//
// Solidity: function InitialTimeStamp() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) InitialTimeStamp() ([]byte, error) {
	return _RevocationList.Contract.InitialTimeStamp(&_RevocationList.CallOpts)
}

// RetrieveBBSPublicKey is a free data retrieval call binding the contract method 0xe6ba9a3c.
//
// Solidity: function RetrieveBBSPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCaller) RetrieveBBSPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveBBSPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RetrieveBBSPublicKey is a free data retrieval call binding the contract method 0xe6ba9a3c.
//
// Solidity: function RetrieveBBSPublicKey() view returns(bytes)
func (_RevocationList *RevocationListSession) RetrieveBBSPublicKey() ([]byte, error) {
	return _RevocationList.Contract.RetrieveBBSPublicKey(&_RevocationList.CallOpts)
}

// RetrieveBBSPublicKey is a free data retrieval call binding the contract method 0xe6ba9a3c.
//
// Solidity: function RetrieveBBSPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) RetrieveBBSPublicKey() ([]byte, error) {
	return _RevocationList.Contract.RetrieveBBSPublicKey(&_RevocationList.CallOpts)
}

// RetrieveCCS is a free data retrieval call binding the contract method 0xb5b1eadd.
//
// Solidity: function RetrieveCCS() view returns(bytes)
func (_RevocationList *RevocationListCaller) RetrieveCCS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveCCS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RetrieveCCS is a free data retrieval call binding the contract method 0xb5b1eadd.
//
// Solidity: function RetrieveCCS() view returns(bytes)
func (_RevocationList *RevocationListSession) RetrieveCCS() ([]byte, error) {
	return _RevocationList.Contract.RetrieveCCS(&_RevocationList.CallOpts)
}

// RetrieveCCS is a free data retrieval call binding the contract method 0xb5b1eadd.
//
// Solidity: function RetrieveCCS() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) RetrieveCCS() ([]byte, error) {
	return _RevocationList.Contract.RetrieveCCS(&_RevocationList.CallOpts)
}

// RetrieveCCSHash is a free data retrieval call binding the contract method 0x2723e24b.
//
// Solidity: function RetrieveCCSHash() view returns(bytes)
func (_RevocationList *RevocationListCaller) RetrieveCCSHash(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveCCSHash")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RetrieveCCSHash is a free data retrieval call binding the contract method 0x2723e24b.
//
// Solidity: function RetrieveCCSHash() view returns(bytes)
func (_RevocationList *RevocationListSession) RetrieveCCSHash() ([]byte, error) {
	return _RevocationList.Contract.RetrieveCCSHash(&_RevocationList.CallOpts)
}

// RetrieveCCSHash is a free data retrieval call binding the contract method 0x2723e24b.
//
// Solidity: function RetrieveCCSHash() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) RetrieveCCSHash() ([]byte, error) {
	return _RevocationList.Contract.RetrieveCCSHash(&_RevocationList.CallOpts)
}

// RetrieveEDDSAPublicKey is a free data retrieval call binding the contract method 0x4c59672e.
//
// Solidity: function RetrieveEDDSAPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCaller) RetrieveEDDSAPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveEDDSAPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RetrieveEDDSAPublicKey is a free data retrieval call binding the contract method 0x4c59672e.
//
// Solidity: function RetrieveEDDSAPublicKey() view returns(bytes)
func (_RevocationList *RevocationListSession) RetrieveEDDSAPublicKey() ([]byte, error) {
	return _RevocationList.Contract.RetrieveEDDSAPublicKey(&_RevocationList.CallOpts)
}

// RetrieveEDDSAPublicKey is a free data retrieval call binding the contract method 0x4c59672e.
//
// Solidity: function RetrieveEDDSAPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) RetrieveEDDSAPublicKey() ([]byte, error) {
	return _RevocationList.Contract.RetrieveEDDSAPublicKey(&_RevocationList.CallOpts)
}

// RetrieveEpochDuration is a free data retrieval call binding the contract method 0xc6a6c0e1.
//
// Solidity: function RetrieveEpochDuration() view returns(uint256)
func (_RevocationList *RevocationListCaller) RetrieveEpochDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveEpochDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RetrieveEpochDuration is a free data retrieval call binding the contract method 0xc6a6c0e1.
//
// Solidity: function RetrieveEpochDuration() view returns(uint256)
func (_RevocationList *RevocationListSession) RetrieveEpochDuration() (*big.Int, error) {
	return _RevocationList.Contract.RetrieveEpochDuration(&_RevocationList.CallOpts)
}

// RetrieveEpochDuration is a free data retrieval call binding the contract method 0xc6a6c0e1.
//
// Solidity: function RetrieveEpochDuration() view returns(uint256)
func (_RevocationList *RevocationListCallerSession) RetrieveEpochDuration() (*big.Int, error) {
	return _RevocationList.Contract.RetrieveEpochDuration(&_RevocationList.CallOpts)
}

// RetrieveInitialTimeStamp is a free data retrieval call binding the contract method 0x4164db7d.
//
// Solidity: function RetrieveInitialTimeStamp() view returns(bytes)
func (_RevocationList *RevocationListCaller) RetrieveInitialTimeStamp(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveInitialTimeStamp")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RetrieveInitialTimeStamp is a free data retrieval call binding the contract method 0x4164db7d.
//
// Solidity: function RetrieveInitialTimeStamp() view returns(bytes)
func (_RevocationList *RevocationListSession) RetrieveInitialTimeStamp() ([]byte, error) {
	return _RevocationList.Contract.RetrieveInitialTimeStamp(&_RevocationList.CallOpts)
}

// RetrieveInitialTimeStamp is a free data retrieval call binding the contract method 0x4164db7d.
//
// Solidity: function RetrieveInitialTimeStamp() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) RetrieveInitialTimeStamp() ([]byte, error) {
	return _RevocationList.Contract.RetrieveInitialTimeStamp(&_RevocationList.CallOpts)
}

// RetrieveZKPVerifyingKey is a free data retrieval call binding the contract method 0x1b717f04.
//
// Solidity: function RetrieveZKPVerifyingKey() view returns(bytes)
func (_RevocationList *RevocationListCaller) RetrieveZKPVerifyingKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "RetrieveZKPVerifyingKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RetrieveZKPVerifyingKey is a free data retrieval call binding the contract method 0x1b717f04.
//
// Solidity: function RetrieveZKPVerifyingKey() view returns(bytes)
func (_RevocationList *RevocationListSession) RetrieveZKPVerifyingKey() ([]byte, error) {
	return _RevocationList.Contract.RetrieveZKPVerifyingKey(&_RevocationList.CallOpts)
}

// RetrieveZKPVerifyingKey is a free data retrieval call binding the contract method 0x1b717f04.
//
// Solidity: function RetrieveZKPVerifyingKey() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) RetrieveZKPVerifyingKey() ([]byte, error) {
	return _RevocationList.Contract.RetrieveZKPVerifyingKey(&_RevocationList.CallOpts)
}

// ZKPVerifyingKey is a free data retrieval call binding the contract method 0xd03c7e66.
//
// Solidity: function ZKPVerifyingKey() view returns(bytes)
func (_RevocationList *RevocationListCaller) ZKPVerifyingKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "ZKPVerifyingKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ZKPVerifyingKey is a free data retrieval call binding the contract method 0xd03c7e66.
//
// Solidity: function ZKPVerifyingKey() view returns(bytes)
func (_RevocationList *RevocationListSession) ZKPVerifyingKey() ([]byte, error) {
	return _RevocationList.Contract.ZKPVerifyingKey(&_RevocationList.CallOpts)
}

// ZKPVerifyingKey is a free data retrieval call binding the contract method 0xd03c7e66.
//
// Solidity: function ZKPVerifyingKey() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) ZKPVerifyingKey() ([]byte, error) {
	return _RevocationList.Contract.ZKPVerifyingKey(&_RevocationList.CallOpts)
}

// BbsPublicKey is a free data retrieval call binding the contract method 0x5cbb46f6.
//
// Solidity: function bbsPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCaller) BbsPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "bbsPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// BbsPublicKey is a free data retrieval call binding the contract method 0x5cbb46f6.
//
// Solidity: function bbsPublicKey() view returns(bytes)
func (_RevocationList *RevocationListSession) BbsPublicKey() ([]byte, error) {
	return _RevocationList.Contract.BbsPublicKey(&_RevocationList.CallOpts)
}

// BbsPublicKey is a free data retrieval call binding the contract method 0x5cbb46f6.
//
// Solidity: function bbsPublicKey() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) BbsPublicKey() ([]byte, error) {
	return _RevocationList.Contract.BbsPublicKey(&_RevocationList.CallOpts)
}

// Ccs is a free data retrieval call binding the contract method 0x67ebfdd5.
//
// Solidity: function ccs() view returns(bytes)
func (_RevocationList *RevocationListCaller) Ccs(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "ccs")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Ccs is a free data retrieval call binding the contract method 0x67ebfdd5.
//
// Solidity: function ccs() view returns(bytes)
func (_RevocationList *RevocationListSession) Ccs() ([]byte, error) {
	return _RevocationList.Contract.Ccs(&_RevocationList.CallOpts)
}

// Ccs is a free data retrieval call binding the contract method 0x67ebfdd5.
//
// Solidity: function ccs() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) Ccs() ([]byte, error) {
	return _RevocationList.Contract.Ccs(&_RevocationList.CallOpts)
}

// CcsHash is a free data retrieval call binding the contract method 0xb8f1b174.
//
// Solidity: function ccsHash() view returns(bytes)
func (_RevocationList *RevocationListCaller) CcsHash(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "ccsHash")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// CcsHash is a free data retrieval call binding the contract method 0xb8f1b174.
//
// Solidity: function ccsHash() view returns(bytes)
func (_RevocationList *RevocationListSession) CcsHash() ([]byte, error) {
	return _RevocationList.Contract.CcsHash(&_RevocationList.CallOpts)
}

// CcsHash is a free data retrieval call binding the contract method 0xb8f1b174.
//
// Solidity: function ccsHash() view returns(bytes)
func (_RevocationList *RevocationListCallerSession) CcsHash() ([]byte, error) {
	return _RevocationList.Contract.CcsHash(&_RevocationList.CallOpts)
}

// List is a free data retrieval call binding the contract method 0x0f560cd7.
//
// Solidity: function list() view returns(uint256 epoch)
func (_RevocationList *RevocationListCaller) List(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RevocationList.contract.Call(opts, &out, "list")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// List is a free data retrieval call binding the contract method 0x0f560cd7.
//
// Solidity: function list() view returns(uint256 epoch)
func (_RevocationList *RevocationListSession) List() (*big.Int, error) {
	return _RevocationList.Contract.List(&_RevocationList.CallOpts)
}

// List is a free data retrieval call binding the contract method 0x0f560cd7.
//
// Solidity: function list() view returns(uint256 epoch)
func (_RevocationList *RevocationListCallerSession) List() (*big.Int, error) {
	return _RevocationList.Contract.List(&_RevocationList.CallOpts)
}

// PublishBBSPublicKey is a paid mutator transaction binding the contract method 0x2b5d9cbe.
//
// Solidity: function PublishBBSPublicKey(bytes _bbsPublicKey) returns()
func (_RevocationList *RevocationListTransactor) PublishBBSPublicKey(opts *bind.TransactOpts, _bbsPublicKey []byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "PublishBBSPublicKey", _bbsPublicKey)
}

// PublishBBSPublicKey is a paid mutator transaction binding the contract method 0x2b5d9cbe.
//
// Solidity: function PublishBBSPublicKey(bytes _bbsPublicKey) returns()
func (_RevocationList *RevocationListSession) PublishBBSPublicKey(_bbsPublicKey []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishBBSPublicKey(&_RevocationList.TransactOpts, _bbsPublicKey)
}

// PublishBBSPublicKey is a paid mutator transaction binding the contract method 0x2b5d9cbe.
//
// Solidity: function PublishBBSPublicKey(bytes _bbsPublicKey) returns()
func (_RevocationList *RevocationListTransactorSession) PublishBBSPublicKey(_bbsPublicKey []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishBBSPublicKey(&_RevocationList.TransactOpts, _bbsPublicKey)
}

// PublishEDDSAPublicKey is a paid mutator transaction binding the contract method 0xb5fd3c5c.
//
// Solidity: function PublishEDDSAPublicKey(bytes _publicKey) returns()
func (_RevocationList *RevocationListTransactor) PublishEDDSAPublicKey(opts *bind.TransactOpts, _publicKey []byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "PublishEDDSAPublicKey", _publicKey)
}

// PublishEDDSAPublicKey is a paid mutator transaction binding the contract method 0xb5fd3c5c.
//
// Solidity: function PublishEDDSAPublicKey(bytes _publicKey) returns()
func (_RevocationList *RevocationListSession) PublishEDDSAPublicKey(_publicKey []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishEDDSAPublicKey(&_RevocationList.TransactOpts, _publicKey)
}

// PublishEDDSAPublicKey is a paid mutator transaction binding the contract method 0xb5fd3c5c.
//
// Solidity: function PublishEDDSAPublicKey(bytes _publicKey) returns()
func (_RevocationList *RevocationListTransactorSession) PublishEDDSAPublicKey(_publicKey []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishEDDSAPublicKey(&_RevocationList.TransactOpts, _publicKey)
}

// PublishEpochConfigurations is a paid mutator transaction binding the contract method 0xbd03b113.
//
// Solidity: function PublishEpochConfigurations(uint256 epoch_duration, bytes intial_timestamp) returns()
func (_RevocationList *RevocationListTransactor) PublishEpochConfigurations(opts *bind.TransactOpts, epoch_duration *big.Int, intial_timestamp []byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "PublishEpochConfigurations", epoch_duration, intial_timestamp)
}

// PublishEpochConfigurations is a paid mutator transaction binding the contract method 0xbd03b113.
//
// Solidity: function PublishEpochConfigurations(uint256 epoch_duration, bytes intial_timestamp) returns()
func (_RevocationList *RevocationListSession) PublishEpochConfigurations(epoch_duration *big.Int, intial_timestamp []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishEpochConfigurations(&_RevocationList.TransactOpts, epoch_duration, intial_timestamp)
}

// PublishEpochConfigurations is a paid mutator transaction binding the contract method 0xbd03b113.
//
// Solidity: function PublishEpochConfigurations(uint256 epoch_duration, bytes intial_timestamp) returns()
func (_RevocationList *RevocationListTransactorSession) PublishEpochConfigurations(epoch_duration *big.Int, intial_timestamp []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishEpochConfigurations(&_RevocationList.TransactOpts, epoch_duration, intial_timestamp)
}

// PublishZKPVerifyingKey is a paid mutator transaction binding the contract method 0xc0de5cae.
//
// Solidity: function PublishZKPVerifyingKey(bytes _verifyingKey) returns()
func (_RevocationList *RevocationListTransactor) PublishZKPVerifyingKey(opts *bind.TransactOpts, _verifyingKey []byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "PublishZKPVerifyingKey", _verifyingKey)
}

// PublishZKPVerifyingKey is a paid mutator transaction binding the contract method 0xc0de5cae.
//
// Solidity: function PublishZKPVerifyingKey(bytes _verifyingKey) returns()
func (_RevocationList *RevocationListSession) PublishZKPVerifyingKey(_verifyingKey []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishZKPVerifyingKey(&_RevocationList.TransactOpts, _verifyingKey)
}

// PublishZKPVerifyingKey is a paid mutator transaction binding the contract method 0xc0de5cae.
//
// Solidity: function PublishZKPVerifyingKey(bytes _verifyingKey) returns()
func (_RevocationList *RevocationListTransactorSession) PublishZKPVerifyingKey(_verifyingKey []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishZKPVerifyingKey(&_RevocationList.TransactOpts, _verifyingKey)
}

// RefreshRevokedTokens is a paid mutator transaction binding the contract method 0x5a9c658c.
//
// Solidity: function RefreshRevokedTokens(uint256 _epoch, bytes32[] _tokens) returns()
func (_RevocationList *RevocationListTransactor) RefreshRevokedTokens(opts *bind.TransactOpts, _epoch *big.Int, _tokens [][32]byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "RefreshRevokedTokens", _epoch, _tokens)
}

// RefreshRevokedTokens is a paid mutator transaction binding the contract method 0x5a9c658c.
//
// Solidity: function RefreshRevokedTokens(uint256 _epoch, bytes32[] _tokens) returns()
func (_RevocationList *RevocationListSession) RefreshRevokedTokens(_epoch *big.Int, _tokens [][32]byte) (*types.Transaction, error) {
	return _RevocationList.Contract.RefreshRevokedTokens(&_RevocationList.TransactOpts, _epoch, _tokens)
}

// RefreshRevokedTokens is a paid mutator transaction binding the contract method 0x5a9c658c.
//
// Solidity: function RefreshRevokedTokens(uint256 _epoch, bytes32[] _tokens) returns()
func (_RevocationList *RevocationListTransactorSession) RefreshRevokedTokens(_epoch *big.Int, _tokens [][32]byte) (*types.Transaction, error) {
	return _RevocationList.Contract.RefreshRevokedTokens(&_RevocationList.TransactOpts, _epoch, _tokens)
}

// PublishCCS is a paid mutator transaction binding the contract method 0x39fcc15e.
//
// Solidity: function publishCCS(bytes _ccs) returns()
func (_RevocationList *RevocationListTransactor) PublishCCS(opts *bind.TransactOpts, _ccs []byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "publishCCS", _ccs)
}

// PublishCCS is a paid mutator transaction binding the contract method 0x39fcc15e.
//
// Solidity: function publishCCS(bytes _ccs) returns()
func (_RevocationList *RevocationListSession) PublishCCS(_ccs []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishCCS(&_RevocationList.TransactOpts, _ccs)
}

// PublishCCS is a paid mutator transaction binding the contract method 0x39fcc15e.
//
// Solidity: function publishCCS(bytes _ccs) returns()
func (_RevocationList *RevocationListTransactorSession) PublishCCS(_ccs []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishCCS(&_RevocationList.TransactOpts, _ccs)
}

// PublishCCSHash is a paid mutator transaction binding the contract method 0x41ebf5a3.
//
// Solidity: function publishCCSHash(bytes _ccsHash) returns()
func (_RevocationList *RevocationListTransactor) PublishCCSHash(opts *bind.TransactOpts, _ccsHash []byte) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "publishCCSHash", _ccsHash)
}

// PublishCCSHash is a paid mutator transaction binding the contract method 0x41ebf5a3.
//
// Solidity: function publishCCSHash(bytes _ccsHash) returns()
func (_RevocationList *RevocationListSession) PublishCCSHash(_ccsHash []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishCCSHash(&_RevocationList.TransactOpts, _ccsHash)
}

// PublishCCSHash is a paid mutator transaction binding the contract method 0x41ebf5a3.
//
// Solidity: function publishCCSHash(bytes _ccsHash) returns()
func (_RevocationList *RevocationListTransactorSession) PublishCCSHash(_ccsHash []byte) (*types.Transaction, error) {
	return _RevocationList.Contract.PublishCCSHash(&_RevocationList.TransactOpts, _ccsHash)
}

// RegisterIssuers is a paid mutator transaction binding the contract method 0x0de54b85.
//
// Solidity: function registerIssuers() returns()
func (_RevocationList *RevocationListTransactor) RegisterIssuers(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "registerIssuers")
}

// RegisterIssuers is a paid mutator transaction binding the contract method 0x0de54b85.
//
// Solidity: function registerIssuers() returns()
func (_RevocationList *RevocationListSession) RegisterIssuers() (*types.Transaction, error) {
	return _RevocationList.Contract.RegisterIssuers(&_RevocationList.TransactOpts)
}

// RegisterIssuers is a paid mutator transaction binding the contract method 0x0de54b85.
//
// Solidity: function registerIssuers() returns()
func (_RevocationList *RevocationListTransactorSession) RegisterIssuers() (*types.Transaction, error) {
	return _RevocationList.Contract.RegisterIssuers(&_RevocationList.TransactOpts)
}

// Test is a paid mutator transaction binding the contract method 0xf8a8fd6d.
//
// Solidity: function test() returns()
func (_RevocationList *RevocationListTransactor) Test(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevocationList.contract.Transact(opts, "test")
}

// Test is a paid mutator transaction binding the contract method 0xf8a8fd6d.
//
// Solidity: function test() returns()
func (_RevocationList *RevocationListSession) Test() (*types.Transaction, error) {
	return _RevocationList.Contract.Test(&_RevocationList.TransactOpts)
}

// Test is a paid mutator transaction binding the contract method 0xf8a8fd6d.
//
// Solidity: function test() returns()
func (_RevocationList *RevocationListTransactorSession) Test() (*types.Transaction, error) {
	return _RevocationList.Contract.Test(&_RevocationList.TransactOpts)
}
