// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polymarketNegRiskFeeModule

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

// Order is an auto generated low-level Go binding around an user-defined struct.
type Order struct {
	Salt          *big.Int
	Maker         common.Address
	Signer        common.Address
	Taker         common.Address
	TokenId       *big.Int
	MakerAmount   *big.Int
	TakerAmount   *big.Int
	Expiration    *big.Int
	Nonce         *big.Int
	FeeRateBps    *big.Int
	Side          uint8
	SignatureType uint8
	Signature     []byte
}

// PolymarketNegRiskFeeModuleMetaData contains all meta data concerning the PolymarketNegRiskFeeModule contract.
var PolymarketNegRiskFeeModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_negRiskCtfExchange\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_negRiskAdapter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_ctf\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"NotAdmin\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdminAddress\",\"type\":\"address\"}],\"name\":\"NewAdmin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"removedAdmin\",\"type\":\"address\"}],\"name\":\"RemovedAdmin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateral\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ctf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exchange\",\"outputs\":[{\"internalType\":\"contractIExchange\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeRateBps\",\"type\":\"uint256\"},{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"},{\"internalType\":\"enumSignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structOrder\",\"name\":\"takerOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeRateBps\",\"type\":\"uint256\"},{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"},{\"internalType\":\"enumSignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structOrder[]\",\"name\":\"makerOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"makerFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"makerFeeRate\",\"type\":\"uint256\"}],\"name\":\"matchOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155BatchReceived\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PolymarketNegRiskFeeModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use PolymarketNegRiskFeeModuleMetaData.ABI instead.
var PolymarketNegRiskFeeModuleABI = PolymarketNegRiskFeeModuleMetaData.ABI

// PolymarketNegRiskFeeModule is an auto generated Go binding around an Ethereum contract.
type PolymarketNegRiskFeeModule struct {
	PolymarketNegRiskFeeModuleCaller     // Read-only binding to the contract
	PolymarketNegRiskFeeModuleTransactor // Write-only binding to the contract
	PolymarketNegRiskFeeModuleFilterer   // Log filterer for contract events
}

// PolymarketNegRiskFeeModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolymarketNegRiskFeeModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolymarketNegRiskFeeModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolymarketNegRiskFeeModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolymarketNegRiskFeeModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolymarketNegRiskFeeModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolymarketNegRiskFeeModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolymarketNegRiskFeeModuleSession struct {
	Contract     *PolymarketNegRiskFeeModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// PolymarketNegRiskFeeModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolymarketNegRiskFeeModuleCallerSession struct {
	Contract *PolymarketNegRiskFeeModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// PolymarketNegRiskFeeModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolymarketNegRiskFeeModuleTransactorSession struct {
	Contract     *PolymarketNegRiskFeeModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// PolymarketNegRiskFeeModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolymarketNegRiskFeeModuleRaw struct {
	Contract *PolymarketNegRiskFeeModule // Generic contract binding to access the raw methods on
}

// PolymarketNegRiskFeeModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolymarketNegRiskFeeModuleCallerRaw struct {
	Contract *PolymarketNegRiskFeeModuleCaller // Generic read-only contract binding to access the raw methods on
}

// PolymarketNegRiskFeeModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolymarketNegRiskFeeModuleTransactorRaw struct {
	Contract *PolymarketNegRiskFeeModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolymarketNegRiskFeeModule creates a new instance of PolymarketNegRiskFeeModule, bound to a specific deployed contract.
func NewPolymarketNegRiskFeeModule(address common.Address, backend bind.ContractBackend) (*PolymarketNegRiskFeeModule, error) {
	contract, err := bindPolymarketNegRiskFeeModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModule{PolymarketNegRiskFeeModuleCaller: PolymarketNegRiskFeeModuleCaller{contract: contract}, PolymarketNegRiskFeeModuleTransactor: PolymarketNegRiskFeeModuleTransactor{contract: contract}, PolymarketNegRiskFeeModuleFilterer: PolymarketNegRiskFeeModuleFilterer{contract: contract}}, nil
}

// NewPolymarketNegRiskFeeModuleCaller creates a new read-only instance of PolymarketNegRiskFeeModule, bound to a specific deployed contract.
func NewPolymarketNegRiskFeeModuleCaller(address common.Address, caller bind.ContractCaller) (*PolymarketNegRiskFeeModuleCaller, error) {
	contract, err := bindPolymarketNegRiskFeeModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleCaller{contract: contract}, nil
}

// NewPolymarketNegRiskFeeModuleTransactor creates a new write-only instance of PolymarketNegRiskFeeModule, bound to a specific deployed contract.
func NewPolymarketNegRiskFeeModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*PolymarketNegRiskFeeModuleTransactor, error) {
	contract, err := bindPolymarketNegRiskFeeModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleTransactor{contract: contract}, nil
}

// NewPolymarketNegRiskFeeModuleFilterer creates a new log filterer instance of PolymarketNegRiskFeeModule, bound to a specific deployed contract.
func NewPolymarketNegRiskFeeModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*PolymarketNegRiskFeeModuleFilterer, error) {
	contract, err := bindPolymarketNegRiskFeeModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleFilterer{contract: contract}, nil
}

// bindPolymarketNegRiskFeeModule binds a generic wrapper to an already deployed contract.
func bindPolymarketNegRiskFeeModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PolymarketNegRiskFeeModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolymarketNegRiskFeeModule.Contract.PolymarketNegRiskFeeModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.PolymarketNegRiskFeeModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.PolymarketNegRiskFeeModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolymarketNegRiskFeeModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.contract.Transact(opts, method, params...)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(uint256)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCaller) Admins(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolymarketNegRiskFeeModule.contract.Call(opts, &out, "admins", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(uint256)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _PolymarketNegRiskFeeModule.Contract.Admins(&_PolymarketNegRiskFeeModule.CallOpts, arg0)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(uint256)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCallerSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _PolymarketNegRiskFeeModule.Contract.Admins(&_PolymarketNegRiskFeeModule.CallOpts, arg0)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCaller) Collateral(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolymarketNegRiskFeeModule.contract.Call(opts, &out, "collateral")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) Collateral() (common.Address, error) {
	return _PolymarketNegRiskFeeModule.Contract.Collateral(&_PolymarketNegRiskFeeModule.CallOpts)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCallerSession) Collateral() (common.Address, error) {
	return _PolymarketNegRiskFeeModule.Contract.Collateral(&_PolymarketNegRiskFeeModule.CallOpts)
}

// Ctf is a free data retrieval call binding the contract method 0x22a9339f.
//
// Solidity: function ctf() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCaller) Ctf(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolymarketNegRiskFeeModule.contract.Call(opts, &out, "ctf")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ctf is a free data retrieval call binding the contract method 0x22a9339f.
//
// Solidity: function ctf() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) Ctf() (common.Address, error) {
	return _PolymarketNegRiskFeeModule.Contract.Ctf(&_PolymarketNegRiskFeeModule.CallOpts)
}

// Ctf is a free data retrieval call binding the contract method 0x22a9339f.
//
// Solidity: function ctf() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCallerSession) Ctf() (common.Address, error) {
	return _PolymarketNegRiskFeeModule.Contract.Ctf(&_PolymarketNegRiskFeeModule.CallOpts)
}

// Exchange is a free data retrieval call binding the contract method 0xd2f7265a.
//
// Solidity: function exchange() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCaller) Exchange(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolymarketNegRiskFeeModule.contract.Call(opts, &out, "exchange")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Exchange is a free data retrieval call binding the contract method 0xd2f7265a.
//
// Solidity: function exchange() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) Exchange() (common.Address, error) {
	return _PolymarketNegRiskFeeModule.Contract.Exchange(&_PolymarketNegRiskFeeModule.CallOpts)
}

// Exchange is a free data retrieval call binding the contract method 0xd2f7265a.
//
// Solidity: function exchange() view returns(address)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCallerSession) Exchange() (common.Address, error) {
	return _PolymarketNegRiskFeeModule.Contract.Exchange(&_PolymarketNegRiskFeeModule.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCaller) IsAdmin(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _PolymarketNegRiskFeeModule.contract.Call(opts, &out, "isAdmin", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) IsAdmin(addr common.Address) (bool, error) {
	return _PolymarketNegRiskFeeModule.Contract.IsAdmin(&_PolymarketNegRiskFeeModule.CallOpts, addr)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleCallerSession) IsAdmin(addr common.Address) (bool, error) {
	return _PolymarketNegRiskFeeModule.Contract.IsAdmin(&_PolymarketNegRiskFeeModule.CallOpts, addr)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) AddAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "addAdmin", admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.AddAdmin(&_PolymarketNegRiskFeeModule.TransactOpts, admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.AddAdmin(&_PolymarketNegRiskFeeModule.TransactOpts, admin)
}

// MatchOrders is a paid mutator transaction binding the contract method 0xd2539b37.
//
// Solidity: function matchOrders((uint256,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint8,uint8,bytes) takerOrder, (uint256,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint8,uint8,bytes)[] makerOrders, uint256 takerFillAmount, uint256[] makerFillAmounts, uint256 makerFeeRate) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) MatchOrders(opts *bind.TransactOpts, takerOrder Order, makerOrders []Order, takerFillAmount *big.Int, makerFillAmounts []*big.Int, makerFeeRate *big.Int) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "matchOrders", takerOrder, makerOrders, takerFillAmount, makerFillAmounts, makerFeeRate)
}

// MatchOrders is a paid mutator transaction binding the contract method 0xd2539b37.
//
// Solidity: function matchOrders((uint256,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint8,uint8,bytes) takerOrder, (uint256,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint8,uint8,bytes)[] makerOrders, uint256 takerFillAmount, uint256[] makerFillAmounts, uint256 makerFeeRate) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) MatchOrders(takerOrder Order, makerOrders []Order, takerFillAmount *big.Int, makerFillAmounts []*big.Int, makerFeeRate *big.Int) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.MatchOrders(&_PolymarketNegRiskFeeModule.TransactOpts, takerOrder, makerOrders, takerFillAmount, makerFillAmounts, makerFeeRate)
}

// MatchOrders is a paid mutator transaction binding the contract method 0xd2539b37.
//
// Solidity: function matchOrders((uint256,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint8,uint8,bytes) takerOrder, (uint256,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint8,uint8,bytes)[] makerOrders, uint256 takerFillAmount, uint256[] makerFillAmounts, uint256 makerFeeRate) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) MatchOrders(takerOrder Order, makerOrders []Order, takerFillAmount *big.Int, makerFillAmounts []*big.Int, makerFeeRate *big.Int) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.MatchOrders(&_PolymarketNegRiskFeeModule.TransactOpts, takerOrder, makerOrders, takerFillAmount, makerFillAmounts, makerFeeRate)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.OnERC1155BatchReceived(&_PolymarketNegRiskFeeModule.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.OnERC1155BatchReceived(&_PolymarketNegRiskFeeModule.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.OnERC1155Received(&_PolymarketNegRiskFeeModule.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.OnERC1155Received(&_PolymarketNegRiskFeeModule.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) RemoveAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "removeAdmin", admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.RemoveAdmin(&_PolymarketNegRiskFeeModule.TransactOpts, admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.RemoveAdmin(&_PolymarketNegRiskFeeModule.TransactOpts, admin)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x8bad0c0a.
//
// Solidity: function renounceAdmin() returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) RenounceAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "renounceAdmin")
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x8bad0c0a.
//
// Solidity: function renounceAdmin() returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) RenounceAdmin() (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.RenounceAdmin(&_PolymarketNegRiskFeeModule.TransactOpts)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x8bad0c0a.
//
// Solidity: function renounceAdmin() returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) RenounceAdmin() (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.RenounceAdmin(&_PolymarketNegRiskFeeModule.TransactOpts)
}

// WithdrawFees is a paid mutator transaction binding the contract method 0x425c2096.
//
// Solidity: function withdrawFees(address to, uint256 id, uint256 amount) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactor) WithdrawFees(opts *bind.TransactOpts, to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.contract.Transact(opts, "withdrawFees", to, id, amount)
}

// WithdrawFees is a paid mutator transaction binding the contract method 0x425c2096.
//
// Solidity: function withdrawFees(address to, uint256 id, uint256 amount) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleSession) WithdrawFees(to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.WithdrawFees(&_PolymarketNegRiskFeeModule.TransactOpts, to, id, amount)
}

// WithdrawFees is a paid mutator transaction binding the contract method 0x425c2096.
//
// Solidity: function withdrawFees(address to, uint256 id, uint256 amount) returns()
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleTransactorSession) WithdrawFees(to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PolymarketNegRiskFeeModule.Contract.WithdrawFees(&_PolymarketNegRiskFeeModule.TransactOpts, to, id, amount)
}

// PolymarketNegRiskFeeModuleFeeRefundedIterator is returned from FilterFeeRefunded and is used to iterate over the raw logs and unpacked data for FeeRefunded events raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleFeeRefundedIterator struct {
	Event *PolymarketNegRiskFeeModuleFeeRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolymarketNegRiskFeeModuleFeeRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketNegRiskFeeModuleFeeRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolymarketNegRiskFeeModuleFeeRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolymarketNegRiskFeeModuleFeeRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketNegRiskFeeModuleFeeRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketNegRiskFeeModuleFeeRefunded represents a FeeRefunded event raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleFeeRefunded struct {
	Token  common.Address
	To     common.Address
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeRefunded is a free log retrieval operation binding the contract event 0x18fe0464eb77016dc4e227eb0d690e4002756d82b442143bbfb874548952b5f2.
//
// Solidity: event FeeRefunded(address token, address to, uint256 id, uint256 amount)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) FilterFeeRefunded(opts *bind.FilterOpts) (*PolymarketNegRiskFeeModuleFeeRefundedIterator, error) {

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.FilterLogs(opts, "FeeRefunded")
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleFeeRefundedIterator{contract: _PolymarketNegRiskFeeModule.contract, event: "FeeRefunded", logs: logs, sub: sub}, nil
}

// WatchFeeRefunded is a free log subscription operation binding the contract event 0x18fe0464eb77016dc4e227eb0d690e4002756d82b442143bbfb874548952b5f2.
//
// Solidity: event FeeRefunded(address token, address to, uint256 id, uint256 amount)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) WatchFeeRefunded(opts *bind.WatchOpts, sink chan<- *PolymarketNegRiskFeeModuleFeeRefunded) (event.Subscription, error) {

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.WatchLogs(opts, "FeeRefunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketNegRiskFeeModuleFeeRefunded)
				if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "FeeRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeRefunded is a log parse operation binding the contract event 0x18fe0464eb77016dc4e227eb0d690e4002756d82b442143bbfb874548952b5f2.
//
// Solidity: event FeeRefunded(address token, address to, uint256 id, uint256 amount)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) ParseFeeRefunded(log types.Log) (*PolymarketNegRiskFeeModuleFeeRefunded, error) {
	event := new(PolymarketNegRiskFeeModuleFeeRefunded)
	if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "FeeRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketNegRiskFeeModuleFeeWithdrawnIterator is returned from FilterFeeWithdrawn and is used to iterate over the raw logs and unpacked data for FeeWithdrawn events raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleFeeWithdrawnIterator struct {
	Event *PolymarketNegRiskFeeModuleFeeWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolymarketNegRiskFeeModuleFeeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketNegRiskFeeModuleFeeWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolymarketNegRiskFeeModuleFeeWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolymarketNegRiskFeeModuleFeeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketNegRiskFeeModuleFeeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketNegRiskFeeModuleFeeWithdrawn represents a FeeWithdrawn event raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleFeeWithdrawn struct {
	Token  common.Address
	To     common.Address
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeWithdrawn is a free log retrieval operation binding the contract event 0x6ce49f8691a80db5eb4f60cd55b14640529346a7ddf9bf8f77a423fa6a10bfdb.
//
// Solidity: event FeeWithdrawn(address token, address to, uint256 id, uint256 amount)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) FilterFeeWithdrawn(opts *bind.FilterOpts) (*PolymarketNegRiskFeeModuleFeeWithdrawnIterator, error) {

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.FilterLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleFeeWithdrawnIterator{contract: _PolymarketNegRiskFeeModule.contract, event: "FeeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFeeWithdrawn is a free log subscription operation binding the contract event 0x6ce49f8691a80db5eb4f60cd55b14640529346a7ddf9bf8f77a423fa6a10bfdb.
//
// Solidity: event FeeWithdrawn(address token, address to, uint256 id, uint256 amount)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) WatchFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *PolymarketNegRiskFeeModuleFeeWithdrawn) (event.Subscription, error) {

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.WatchLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketNegRiskFeeModuleFeeWithdrawn)
				if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeWithdrawn is a log parse operation binding the contract event 0x6ce49f8691a80db5eb4f60cd55b14640529346a7ddf9bf8f77a423fa6a10bfdb.
//
// Solidity: event FeeWithdrawn(address token, address to, uint256 id, uint256 amount)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) ParseFeeWithdrawn(log types.Log) (*PolymarketNegRiskFeeModuleFeeWithdrawn, error) {
	event := new(PolymarketNegRiskFeeModuleFeeWithdrawn)
	if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketNegRiskFeeModuleNewAdminIterator is returned from FilterNewAdmin and is used to iterate over the raw logs and unpacked data for NewAdmin events raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleNewAdminIterator struct {
	Event *PolymarketNegRiskFeeModuleNewAdmin // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolymarketNegRiskFeeModuleNewAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketNegRiskFeeModuleNewAdmin)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolymarketNegRiskFeeModuleNewAdmin)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolymarketNegRiskFeeModuleNewAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketNegRiskFeeModuleNewAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketNegRiskFeeModuleNewAdmin represents a NewAdmin event raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleNewAdmin struct {
	Admin           common.Address
	NewAdminAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewAdmin is a free log retrieval operation binding the contract event 0xf9ffabca9c8276e99321725bcb43fb076a6c66a54b7f21c4e8146d8519b417dc.
//
// Solidity: event NewAdmin(address indexed admin, address indexed newAdminAddress)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) FilterNewAdmin(opts *bind.FilterOpts, admin []common.Address, newAdminAddress []common.Address) (*PolymarketNegRiskFeeModuleNewAdminIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var newAdminAddressRule []interface{}
	for _, newAdminAddressItem := range newAdminAddress {
		newAdminAddressRule = append(newAdminAddressRule, newAdminAddressItem)
	}

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.FilterLogs(opts, "NewAdmin", adminRule, newAdminAddressRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleNewAdminIterator{contract: _PolymarketNegRiskFeeModule.contract, event: "NewAdmin", logs: logs, sub: sub}, nil
}

// WatchNewAdmin is a free log subscription operation binding the contract event 0xf9ffabca9c8276e99321725bcb43fb076a6c66a54b7f21c4e8146d8519b417dc.
//
// Solidity: event NewAdmin(address indexed admin, address indexed newAdminAddress)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) WatchNewAdmin(opts *bind.WatchOpts, sink chan<- *PolymarketNegRiskFeeModuleNewAdmin, admin []common.Address, newAdminAddress []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var newAdminAddressRule []interface{}
	for _, newAdminAddressItem := range newAdminAddress {
		newAdminAddressRule = append(newAdminAddressRule, newAdminAddressItem)
	}

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.WatchLogs(opts, "NewAdmin", adminRule, newAdminAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketNegRiskFeeModuleNewAdmin)
				if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "NewAdmin", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewAdmin is a log parse operation binding the contract event 0xf9ffabca9c8276e99321725bcb43fb076a6c66a54b7f21c4e8146d8519b417dc.
//
// Solidity: event NewAdmin(address indexed admin, address indexed newAdminAddress)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) ParseNewAdmin(log types.Log) (*PolymarketNegRiskFeeModuleNewAdmin, error) {
	event := new(PolymarketNegRiskFeeModuleNewAdmin)
	if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "NewAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketNegRiskFeeModuleRemovedAdminIterator is returned from FilterRemovedAdmin and is used to iterate over the raw logs and unpacked data for RemovedAdmin events raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleRemovedAdminIterator struct {
	Event *PolymarketNegRiskFeeModuleRemovedAdmin // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolymarketNegRiskFeeModuleRemovedAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketNegRiskFeeModuleRemovedAdmin)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolymarketNegRiskFeeModuleRemovedAdmin)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolymarketNegRiskFeeModuleRemovedAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketNegRiskFeeModuleRemovedAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketNegRiskFeeModuleRemovedAdmin represents a RemovedAdmin event raised by the PolymarketNegRiskFeeModule contract.
type PolymarketNegRiskFeeModuleRemovedAdmin struct {
	Admin        common.Address
	RemovedAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedAdmin is a free log retrieval operation binding the contract event 0x787a2e12f4a55b658b8f573c32432ee11a5e8b51677d1e1e937aaf6a0bb5776e.
//
// Solidity: event RemovedAdmin(address indexed admin, address indexed removedAdmin)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) FilterRemovedAdmin(opts *bind.FilterOpts, admin []common.Address, removedAdmin []common.Address) (*PolymarketNegRiskFeeModuleRemovedAdminIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var removedAdminRule []interface{}
	for _, removedAdminItem := range removedAdmin {
		removedAdminRule = append(removedAdminRule, removedAdminItem)
	}

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.FilterLogs(opts, "RemovedAdmin", adminRule, removedAdminRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketNegRiskFeeModuleRemovedAdminIterator{contract: _PolymarketNegRiskFeeModule.contract, event: "RemovedAdmin", logs: logs, sub: sub}, nil
}

// WatchRemovedAdmin is a free log subscription operation binding the contract event 0x787a2e12f4a55b658b8f573c32432ee11a5e8b51677d1e1e937aaf6a0bb5776e.
//
// Solidity: event RemovedAdmin(address indexed admin, address indexed removedAdmin)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) WatchRemovedAdmin(opts *bind.WatchOpts, sink chan<- *PolymarketNegRiskFeeModuleRemovedAdmin, admin []common.Address, removedAdmin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var removedAdminRule []interface{}
	for _, removedAdminItem := range removedAdmin {
		removedAdminRule = append(removedAdminRule, removedAdminItem)
	}

	logs, sub, err := _PolymarketNegRiskFeeModule.contract.WatchLogs(opts, "RemovedAdmin", adminRule, removedAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketNegRiskFeeModuleRemovedAdmin)
				if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "RemovedAdmin", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemovedAdmin is a log parse operation binding the contract event 0x787a2e12f4a55b658b8f573c32432ee11a5e8b51677d1e1e937aaf6a0bb5776e.
//
// Solidity: event RemovedAdmin(address indexed admin, address indexed removedAdmin)
func (_PolymarketNegRiskFeeModule *PolymarketNegRiskFeeModuleFilterer) ParseRemovedAdmin(log types.Log) (*PolymarketNegRiskFeeModuleRemovedAdmin, error) {
	event := new(PolymarketNegRiskFeeModuleRemovedAdmin)
	if err := _PolymarketNegRiskFeeModule.contract.UnpackLog(event, "RemovedAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
