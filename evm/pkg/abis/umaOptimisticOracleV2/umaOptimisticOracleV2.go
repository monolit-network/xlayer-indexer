// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package umaOptimisticOracleV2

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

// OptimisticOracleV2InterfaceRequest is an auto generated low-level Go binding around an user-defined struct.
type OptimisticOracleV2InterfaceRequest struct {
	Proposer        common.Address
	Disputer        common.Address
	Currency        common.Address
	Settled         bool
	RequestSettings OptimisticOracleV2InterfaceRequestSettings
	ProposedPrice   *big.Int
	ResolvedPrice   *big.Int
	ExpirationTime  *big.Int
	Reward          *big.Int
	FinalFee        *big.Int
}

// OptimisticOracleV2InterfaceRequestSettings is an auto generated low-level Go binding around an user-defined struct.
type OptimisticOracleV2InterfaceRequestSettings struct {
	EventBased              bool
	RefundOnDispute         bool
	CallbackOnPriceProposed bool
	CallbackOnPriceDisputed bool
	CallbackOnPriceSettled  bool
	Bond                    *big.Int
	CustomLiveness          *big.Int
}

// UmaOptimisticOracleV2MetaData contains all meta data concerning the UmaOptimisticOracleV2 contract.
var UmaOptimisticOracleV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_liveness\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_finderAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_timerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"proposedPrice\",\"type\":\"int256\"}],\"name\":\"DisputePrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"proposedPrice\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expirationTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"currency\",\"type\":\"address\"}],\"name\":\"ProposePrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"currency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalFee\",\"type\":\"uint256\"}],\"name\":\"RequestPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payout\",\"type\":\"uint256\"}],\"name\":\"Settle\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"OO_ANCILLARY_DATA_LIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TOO_EARLY_RESPONSE\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ancillaryBytesLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultLiveness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"disputePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBond\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"disputePriceFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBond\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finder\",\"outputs\":[{\"internalType\":\"contractFinderInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"getRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"settled\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"eventBased\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"refundOnDispute\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceProposed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceDisputed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceSettled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customLiveness\",\"type\":\"uint256\"}],\"internalType\":\"structOptimisticOracleV2Interface.RequestSettings\",\"name\":\"requestSettings\",\"type\":\"tuple\"},{\"internalType\":\"int256\",\"name\":\"proposedPrice\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"resolvedPrice\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"finalFee\",\"type\":\"uint256\"}],\"internalType\":\"structOptimisticOracleV2Interface.Request\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"enumOptimisticOracleV2Interface.State\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"hasPrice\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"proposedPrice\",\"type\":\"int256\"}],\"name\":\"proposePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBond\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"proposedPrice\",\"type\":\"int256\"}],\"name\":\"proposePriceFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBond\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"requestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBond\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"requests\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"settled\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"eventBased\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"refundOnDispute\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceProposed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceDisputed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceSettled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customLiveness\",\"type\":\"uint256\"}],\"internalType\":\"structOptimisticOracleV2Interface.RequestSettings\",\"name\":\"requestSettings\",\"type\":\"tuple\"},{\"internalType\":\"int256\",\"name\":\"proposedPrice\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"resolvedPrice\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"finalFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"}],\"name\":\"setBond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalBond\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceProposed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceDisputed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"callbackOnPriceSettled\",\"type\":\"bool\"}],\"name\":\"setCallbacks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"setCurrentTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"customLiveness\",\"type\":\"uint256\"}],\"name\":\"setCustomLiveness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"setEventBased\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"setRefundOnDispute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"settle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"payout\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"name\":\"settleAndGetPrice\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"stampAncillaryData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// UmaOptimisticOracleV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use UmaOptimisticOracleV2MetaData.ABI instead.
var UmaOptimisticOracleV2ABI = UmaOptimisticOracleV2MetaData.ABI

// UmaOptimisticOracleV2 is an auto generated Go binding around an Ethereum contract.
type UmaOptimisticOracleV2 struct {
	UmaOptimisticOracleV2Caller     // Read-only binding to the contract
	UmaOptimisticOracleV2Transactor // Write-only binding to the contract
	UmaOptimisticOracleV2Filterer   // Log filterer for contract events
}

// UmaOptimisticOracleV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UmaOptimisticOracleV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UmaOptimisticOracleV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UmaOptimisticOracleV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UmaOptimisticOracleV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UmaOptimisticOracleV2Session struct {
	Contract     *UmaOptimisticOracleV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UmaOptimisticOracleV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UmaOptimisticOracleV2CallerSession struct {
	Contract *UmaOptimisticOracleV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// UmaOptimisticOracleV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UmaOptimisticOracleV2TransactorSession struct {
	Contract     *UmaOptimisticOracleV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// UmaOptimisticOracleV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type UmaOptimisticOracleV2Raw struct {
	Contract *UmaOptimisticOracleV2 // Generic contract binding to access the raw methods on
}

// UmaOptimisticOracleV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV2CallerRaw struct {
	Contract *UmaOptimisticOracleV2Caller // Generic read-only contract binding to access the raw methods on
}

// UmaOptimisticOracleV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV2TransactorRaw struct {
	Contract *UmaOptimisticOracleV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewUmaOptimisticOracleV2 creates a new instance of UmaOptimisticOracleV2, bound to a specific deployed contract.
func NewUmaOptimisticOracleV2(address common.Address, backend bind.ContractBackend) (*UmaOptimisticOracleV2, error) {
	contract, err := bindUmaOptimisticOracleV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2{UmaOptimisticOracleV2Caller: UmaOptimisticOracleV2Caller{contract: contract}, UmaOptimisticOracleV2Transactor: UmaOptimisticOracleV2Transactor{contract: contract}, UmaOptimisticOracleV2Filterer: UmaOptimisticOracleV2Filterer{contract: contract}}, nil
}

// NewUmaOptimisticOracleV2Caller creates a new read-only instance of UmaOptimisticOracleV2, bound to a specific deployed contract.
func NewUmaOptimisticOracleV2Caller(address common.Address, caller bind.ContractCaller) (*UmaOptimisticOracleV2Caller, error) {
	contract, err := bindUmaOptimisticOracleV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2Caller{contract: contract}, nil
}

// NewUmaOptimisticOracleV2Transactor creates a new write-only instance of UmaOptimisticOracleV2, bound to a specific deployed contract.
func NewUmaOptimisticOracleV2Transactor(address common.Address, transactor bind.ContractTransactor) (*UmaOptimisticOracleV2Transactor, error) {
	contract, err := bindUmaOptimisticOracleV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2Transactor{contract: contract}, nil
}

// NewUmaOptimisticOracleV2Filterer creates a new log filterer instance of UmaOptimisticOracleV2, bound to a specific deployed contract.
func NewUmaOptimisticOracleV2Filterer(address common.Address, filterer bind.ContractFilterer) (*UmaOptimisticOracleV2Filterer, error) {
	contract, err := bindUmaOptimisticOracleV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2Filterer{contract: contract}, nil
}

// bindUmaOptimisticOracleV2 binds a generic wrapper to an already deployed contract.
func bindUmaOptimisticOracleV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UmaOptimisticOracleV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UmaOptimisticOracleV2.Contract.UmaOptimisticOracleV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.UmaOptimisticOracleV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.UmaOptimisticOracleV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UmaOptimisticOracleV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.contract.Transact(opts, method, params...)
}

// OOANCILLARYDATALIMIT is a free data retrieval call binding the contract method 0x4ccb56f5.
//
// Solidity: function OO_ANCILLARY_DATA_LIMIT() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) OOANCILLARYDATALIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "OO_ANCILLARY_DATA_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OOANCILLARYDATALIMIT is a free data retrieval call binding the contract method 0x4ccb56f5.
//
// Solidity: function OO_ANCILLARY_DATA_LIMIT() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) OOANCILLARYDATALIMIT() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.OOANCILLARYDATALIMIT(&_UmaOptimisticOracleV2.CallOpts)
}

// OOANCILLARYDATALIMIT is a free data retrieval call binding the contract method 0x4ccb56f5.
//
// Solidity: function OO_ANCILLARY_DATA_LIMIT() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) OOANCILLARYDATALIMIT() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.OOANCILLARYDATALIMIT(&_UmaOptimisticOracleV2.CallOpts)
}

// TOOEARLYRESPONSE is a free data retrieval call binding the contract method 0xff8c1a8c.
//
// Solidity: function TOO_EARLY_RESPONSE() view returns(int256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) TOOEARLYRESPONSE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "TOO_EARLY_RESPONSE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOOEARLYRESPONSE is a free data retrieval call binding the contract method 0xff8c1a8c.
//
// Solidity: function TOO_EARLY_RESPONSE() view returns(int256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) TOOEARLYRESPONSE() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.TOOEARLYRESPONSE(&_UmaOptimisticOracleV2.CallOpts)
}

// TOOEARLYRESPONSE is a free data retrieval call binding the contract method 0xff8c1a8c.
//
// Solidity: function TOO_EARLY_RESPONSE() view returns(int256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) TOOEARLYRESPONSE() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.TOOEARLYRESPONSE(&_UmaOptimisticOracleV2.CallOpts)
}

// AncillaryBytesLimit is a free data retrieval call binding the contract method 0xc371dda7.
//
// Solidity: function ancillaryBytesLimit() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) AncillaryBytesLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "ancillaryBytesLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AncillaryBytesLimit is a free data retrieval call binding the contract method 0xc371dda7.
//
// Solidity: function ancillaryBytesLimit() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) AncillaryBytesLimit() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.AncillaryBytesLimit(&_UmaOptimisticOracleV2.CallOpts)
}

// AncillaryBytesLimit is a free data retrieval call binding the contract method 0xc371dda7.
//
// Solidity: function ancillaryBytesLimit() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) AncillaryBytesLimit() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.AncillaryBytesLimit(&_UmaOptimisticOracleV2.CallOpts)
}

// DefaultLiveness is a free data retrieval call binding the contract method 0xfe4e1983.
//
// Solidity: function defaultLiveness() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) DefaultLiveness(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "defaultLiveness")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultLiveness is a free data retrieval call binding the contract method 0xfe4e1983.
//
// Solidity: function defaultLiveness() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) DefaultLiveness() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.DefaultLiveness(&_UmaOptimisticOracleV2.CallOpts)
}

// DefaultLiveness is a free data retrieval call binding the contract method 0xfe4e1983.
//
// Solidity: function defaultLiveness() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) DefaultLiveness() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.DefaultLiveness(&_UmaOptimisticOracleV2.CallOpts)
}

// Finder is a free data retrieval call binding the contract method 0xb9a3c84c.
//
// Solidity: function finder() view returns(address)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) Finder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "finder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Finder is a free data retrieval call binding the contract method 0xb9a3c84c.
//
// Solidity: function finder() view returns(address)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) Finder() (common.Address, error) {
	return _UmaOptimisticOracleV2.Contract.Finder(&_UmaOptimisticOracleV2.CallOpts)
}

// Finder is a free data retrieval call binding the contract method 0xb9a3c84c.
//
// Solidity: function finder() view returns(address)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) Finder() (common.Address, error) {
	return _UmaOptimisticOracleV2.Contract.Finder(&_UmaOptimisticOracleV2.CallOpts)
}

// GetCurrentTime is a free data retrieval call binding the contract method 0x29cb924d.
//
// Solidity: function getCurrentTime() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) GetCurrentTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "getCurrentTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentTime is a free data retrieval call binding the contract method 0x29cb924d.
//
// Solidity: function getCurrentTime() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) GetCurrentTime() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.GetCurrentTime(&_UmaOptimisticOracleV2.CallOpts)
}

// GetCurrentTime is a free data retrieval call binding the contract method 0x29cb924d.
//
// Solidity: function getCurrentTime() view returns(uint256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) GetCurrentTime() (*big.Int, error) {
	return _UmaOptimisticOracleV2.Contract.GetCurrentTime(&_UmaOptimisticOracleV2.CallOpts)
}

// GetRequest is a free data retrieval call binding the contract method 0xa9904f9b.
//
// Solidity: function getRequest(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns((address,address,address,bool,(bool,bool,bool,bool,bool,uint256,uint256),int256,int256,uint256,uint256,uint256))
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) GetRequest(opts *bind.CallOpts, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (OptimisticOracleV2InterfaceRequest, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "getRequest", requester, identifier, timestamp, ancillaryData)

	if err != nil {
		return *new(OptimisticOracleV2InterfaceRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(OptimisticOracleV2InterfaceRequest)).(*OptimisticOracleV2InterfaceRequest)

	return out0, err

}

// GetRequest is a free data retrieval call binding the contract method 0xa9904f9b.
//
// Solidity: function getRequest(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns((address,address,address,bool,(bool,bool,bool,bool,bool,uint256,uint256),int256,int256,uint256,uint256,uint256))
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) GetRequest(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (OptimisticOracleV2InterfaceRequest, error) {
	return _UmaOptimisticOracleV2.Contract.GetRequest(&_UmaOptimisticOracleV2.CallOpts, requester, identifier, timestamp, ancillaryData)
}

// GetRequest is a free data retrieval call binding the contract method 0xa9904f9b.
//
// Solidity: function getRequest(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns((address,address,address,bool,(bool,bool,bool,bool,bool,uint256,uint256),int256,int256,uint256,uint256,uint256))
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) GetRequest(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (OptimisticOracleV2InterfaceRequest, error) {
	return _UmaOptimisticOracleV2.Contract.GetRequest(&_UmaOptimisticOracleV2.CallOpts, requester, identifier, timestamp, ancillaryData)
}

// GetState is a free data retrieval call binding the contract method 0xba4b930c.
//
// Solidity: function getState(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns(uint8)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) GetState(opts *bind.CallOpts, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (uint8, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "getState", requester, identifier, timestamp, ancillaryData)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetState is a free data retrieval call binding the contract method 0xba4b930c.
//
// Solidity: function getState(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns(uint8)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) GetState(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (uint8, error) {
	return _UmaOptimisticOracleV2.Contract.GetState(&_UmaOptimisticOracleV2.CallOpts, requester, identifier, timestamp, ancillaryData)
}

// GetState is a free data retrieval call binding the contract method 0xba4b930c.
//
// Solidity: function getState(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns(uint8)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) GetState(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (uint8, error) {
	return _UmaOptimisticOracleV2.Contract.GetState(&_UmaOptimisticOracleV2.CallOpts, requester, identifier, timestamp, ancillaryData)
}

// HasPrice is a free data retrieval call binding the contract method 0xbc58ccaa.
//
// Solidity: function hasPrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns(bool)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) HasPrice(opts *bind.CallOpts, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (bool, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "hasPrice", requester, identifier, timestamp, ancillaryData)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPrice is a free data retrieval call binding the contract method 0xbc58ccaa.
//
// Solidity: function hasPrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns(bool)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) HasPrice(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (bool, error) {
	return _UmaOptimisticOracleV2.Contract.HasPrice(&_UmaOptimisticOracleV2.CallOpts, requester, identifier, timestamp, ancillaryData)
}

// HasPrice is a free data retrieval call binding the contract method 0xbc58ccaa.
//
// Solidity: function hasPrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) view returns(bool)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) HasPrice(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (bool, error) {
	return _UmaOptimisticOracleV2.Contract.HasPrice(&_UmaOptimisticOracleV2.CallOpts, requester, identifier, timestamp, ancillaryData)
}

// Requests is a free data retrieval call binding the contract method 0x9d866985.
//
// Solidity: function requests(bytes32 ) view returns(address proposer, address disputer, address currency, bool settled, (bool,bool,bool,bool,bool,uint256,uint256) requestSettings, int256 proposedPrice, int256 resolvedPrice, uint256 expirationTime, uint256 reward, uint256 finalFee)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) Requests(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Proposer        common.Address
	Disputer        common.Address
	Currency        common.Address
	Settled         bool
	RequestSettings OptimisticOracleV2InterfaceRequestSettings
	ProposedPrice   *big.Int
	ResolvedPrice   *big.Int
	ExpirationTime  *big.Int
	Reward          *big.Int
	FinalFee        *big.Int
}, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "requests", arg0)

	outstruct := new(struct {
		Proposer        common.Address
		Disputer        common.Address
		Currency        common.Address
		Settled         bool
		RequestSettings OptimisticOracleV2InterfaceRequestSettings
		ProposedPrice   *big.Int
		ResolvedPrice   *big.Int
		ExpirationTime  *big.Int
		Reward          *big.Int
		FinalFee        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Proposer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Disputer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Currency = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Settled = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.RequestSettings = *abi.ConvertType(out[4], new(OptimisticOracleV2InterfaceRequestSettings)).(*OptimisticOracleV2InterfaceRequestSettings)
	outstruct.ProposedPrice = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ResolvedPrice = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.ExpirationTime = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Reward = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.FinalFee = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Requests is a free data retrieval call binding the contract method 0x9d866985.
//
// Solidity: function requests(bytes32 ) view returns(address proposer, address disputer, address currency, bool settled, (bool,bool,bool,bool,bool,uint256,uint256) requestSettings, int256 proposedPrice, int256 resolvedPrice, uint256 expirationTime, uint256 reward, uint256 finalFee)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) Requests(arg0 [32]byte) (struct {
	Proposer        common.Address
	Disputer        common.Address
	Currency        common.Address
	Settled         bool
	RequestSettings OptimisticOracleV2InterfaceRequestSettings
	ProposedPrice   *big.Int
	ResolvedPrice   *big.Int
	ExpirationTime  *big.Int
	Reward          *big.Int
	FinalFee        *big.Int
}, error) {
	return _UmaOptimisticOracleV2.Contract.Requests(&_UmaOptimisticOracleV2.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x9d866985.
//
// Solidity: function requests(bytes32 ) view returns(address proposer, address disputer, address currency, bool settled, (bool,bool,bool,bool,bool,uint256,uint256) requestSettings, int256 proposedPrice, int256 resolvedPrice, uint256 expirationTime, uint256 reward, uint256 finalFee)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) Requests(arg0 [32]byte) (struct {
	Proposer        common.Address
	Disputer        common.Address
	Currency        common.Address
	Settled         bool
	RequestSettings OptimisticOracleV2InterfaceRequestSettings
	ProposedPrice   *big.Int
	ResolvedPrice   *big.Int
	ExpirationTime  *big.Int
	Reward          *big.Int
	FinalFee        *big.Int
}, error) {
	return _UmaOptimisticOracleV2.Contract.Requests(&_UmaOptimisticOracleV2.CallOpts, arg0)
}

// StampAncillaryData is a free data retrieval call binding the contract method 0xaf5d2f39.
//
// Solidity: function stampAncillaryData(bytes ancillaryData, address requester) pure returns(bytes)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) StampAncillaryData(opts *bind.CallOpts, ancillaryData []byte, requester common.Address) ([]byte, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "stampAncillaryData", ancillaryData, requester)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// StampAncillaryData is a free data retrieval call binding the contract method 0xaf5d2f39.
//
// Solidity: function stampAncillaryData(bytes ancillaryData, address requester) pure returns(bytes)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) StampAncillaryData(ancillaryData []byte, requester common.Address) ([]byte, error) {
	return _UmaOptimisticOracleV2.Contract.StampAncillaryData(&_UmaOptimisticOracleV2.CallOpts, ancillaryData, requester)
}

// StampAncillaryData is a free data retrieval call binding the contract method 0xaf5d2f39.
//
// Solidity: function stampAncillaryData(bytes ancillaryData, address requester) pure returns(bytes)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) StampAncillaryData(ancillaryData []byte, requester common.Address) ([]byte, error) {
	return _UmaOptimisticOracleV2.Contract.StampAncillaryData(&_UmaOptimisticOracleV2.CallOpts, ancillaryData, requester)
}

// TimerAddress is a free data retrieval call binding the contract method 0x1c39c38d.
//
// Solidity: function timerAddress() view returns(address)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Caller) TimerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV2.contract.Call(opts, &out, "timerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TimerAddress is a free data retrieval call binding the contract method 0x1c39c38d.
//
// Solidity: function timerAddress() view returns(address)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) TimerAddress() (common.Address, error) {
	return _UmaOptimisticOracleV2.Contract.TimerAddress(&_UmaOptimisticOracleV2.CallOpts)
}

// TimerAddress is a free data retrieval call binding the contract method 0x1c39c38d.
//
// Solidity: function timerAddress() view returns(address)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2CallerSession) TimerAddress() (common.Address, error) {
	return _UmaOptimisticOracleV2.Contract.TimerAddress(&_UmaOptimisticOracleV2.CallOpts)
}

// DisputePrice is a paid mutator transaction binding the contract method 0xfba7f1e3.
//
// Solidity: function disputePrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) DisputePrice(opts *bind.TransactOpts, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "disputePrice", requester, identifier, timestamp, ancillaryData)
}

// DisputePrice is a paid mutator transaction binding the contract method 0xfba7f1e3.
//
// Solidity: function disputePrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) DisputePrice(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.DisputePrice(&_UmaOptimisticOracleV2.TransactOpts, requester, identifier, timestamp, ancillaryData)
}

// DisputePrice is a paid mutator transaction binding the contract method 0xfba7f1e3.
//
// Solidity: function disputePrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) DisputePrice(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.DisputePrice(&_UmaOptimisticOracleV2.TransactOpts, requester, identifier, timestamp, ancillaryData)
}

// DisputePriceFor is a paid mutator transaction binding the contract method 0x76c7823f.
//
// Solidity: function disputePriceFor(address disputer, address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) DisputePriceFor(opts *bind.TransactOpts, disputer common.Address, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "disputePriceFor", disputer, requester, identifier, timestamp, ancillaryData)
}

// DisputePriceFor is a paid mutator transaction binding the contract method 0x76c7823f.
//
// Solidity: function disputePriceFor(address disputer, address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) DisputePriceFor(disputer common.Address, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.DisputePriceFor(&_UmaOptimisticOracleV2.TransactOpts, disputer, requester, identifier, timestamp, ancillaryData)
}

// DisputePriceFor is a paid mutator transaction binding the contract method 0x76c7823f.
//
// Solidity: function disputePriceFor(address disputer, address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) DisputePriceFor(disputer common.Address, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.DisputePriceFor(&_UmaOptimisticOracleV2.TransactOpts, disputer, requester, identifier, timestamp, ancillaryData)
}

// ProposePrice is a paid mutator transaction binding the contract method 0xb8b4f908.
//
// Solidity: function proposePrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) ProposePrice(opts *bind.TransactOpts, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, proposedPrice *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "proposePrice", requester, identifier, timestamp, ancillaryData, proposedPrice)
}

// ProposePrice is a paid mutator transaction binding the contract method 0xb8b4f908.
//
// Solidity: function proposePrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) ProposePrice(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, proposedPrice *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.ProposePrice(&_UmaOptimisticOracleV2.TransactOpts, requester, identifier, timestamp, ancillaryData, proposedPrice)
}

// ProposePrice is a paid mutator transaction binding the contract method 0xb8b4f908.
//
// Solidity: function proposePrice(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) ProposePrice(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, proposedPrice *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.ProposePrice(&_UmaOptimisticOracleV2.TransactOpts, requester, identifier, timestamp, ancillaryData, proposedPrice)
}

// ProposePriceFor is a paid mutator transaction binding the contract method 0x7c82288f.
//
// Solidity: function proposePriceFor(address proposer, address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) ProposePriceFor(opts *bind.TransactOpts, proposer common.Address, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, proposedPrice *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "proposePriceFor", proposer, requester, identifier, timestamp, ancillaryData, proposedPrice)
}

// ProposePriceFor is a paid mutator transaction binding the contract method 0x7c82288f.
//
// Solidity: function proposePriceFor(address proposer, address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) ProposePriceFor(proposer common.Address, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, proposedPrice *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.ProposePriceFor(&_UmaOptimisticOracleV2.TransactOpts, proposer, requester, identifier, timestamp, ancillaryData, proposedPrice)
}

// ProposePriceFor is a paid mutator transaction binding the contract method 0x7c82288f.
//
// Solidity: function proposePriceFor(address proposer, address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) ProposePriceFor(proposer common.Address, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, proposedPrice *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.ProposePriceFor(&_UmaOptimisticOracleV2.TransactOpts, proposer, requester, identifier, timestamp, ancillaryData, proposedPrice)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x11df92f1.
//
// Solidity: function requestPrice(bytes32 identifier, uint256 timestamp, bytes ancillaryData, address currency, uint256 reward) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) RequestPrice(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, currency common.Address, reward *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "requestPrice", identifier, timestamp, ancillaryData, currency, reward)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x11df92f1.
//
// Solidity: function requestPrice(bytes32 identifier, uint256 timestamp, bytes ancillaryData, address currency, uint256 reward) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) RequestPrice(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, currency common.Address, reward *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.RequestPrice(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, currency, reward)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x11df92f1.
//
// Solidity: function requestPrice(bytes32 identifier, uint256 timestamp, bytes ancillaryData, address currency, uint256 reward) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) RequestPrice(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, currency common.Address, reward *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.RequestPrice(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, currency, reward)
}

// SetBond is a paid mutator transaction binding the contract method 0xad5a755a.
//
// Solidity: function setBond(bytes32 identifier, uint256 timestamp, bytes ancillaryData, uint256 bond) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SetBond(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, bond *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "setBond", identifier, timestamp, ancillaryData, bond)
}

// SetBond is a paid mutator transaction binding the contract method 0xad5a755a.
//
// Solidity: function setBond(bytes32 identifier, uint256 timestamp, bytes ancillaryData, uint256 bond) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SetBond(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, bond *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetBond(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, bond)
}

// SetBond is a paid mutator transaction binding the contract method 0xad5a755a.
//
// Solidity: function setBond(bytes32 identifier, uint256 timestamp, bytes ancillaryData, uint256 bond) returns(uint256 totalBond)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SetBond(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, bond *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetBond(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, bond)
}

// SetCallbacks is a paid mutator transaction binding the contract method 0xf327b075.
//
// Solidity: function setCallbacks(bytes32 identifier, uint256 timestamp, bytes ancillaryData, bool callbackOnPriceProposed, bool callbackOnPriceDisputed, bool callbackOnPriceSettled) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SetCallbacks(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, callbackOnPriceProposed bool, callbackOnPriceDisputed bool, callbackOnPriceSettled bool) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "setCallbacks", identifier, timestamp, ancillaryData, callbackOnPriceProposed, callbackOnPriceDisputed, callbackOnPriceSettled)
}

// SetCallbacks is a paid mutator transaction binding the contract method 0xf327b075.
//
// Solidity: function setCallbacks(bytes32 identifier, uint256 timestamp, bytes ancillaryData, bool callbackOnPriceProposed, bool callbackOnPriceDisputed, bool callbackOnPriceSettled) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SetCallbacks(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, callbackOnPriceProposed bool, callbackOnPriceDisputed bool, callbackOnPriceSettled bool) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetCallbacks(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, callbackOnPriceProposed, callbackOnPriceDisputed, callbackOnPriceSettled)
}

// SetCallbacks is a paid mutator transaction binding the contract method 0xf327b075.
//
// Solidity: function setCallbacks(bytes32 identifier, uint256 timestamp, bytes ancillaryData, bool callbackOnPriceProposed, bool callbackOnPriceDisputed, bool callbackOnPriceSettled) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SetCallbacks(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, callbackOnPriceProposed bool, callbackOnPriceDisputed bool, callbackOnPriceSettled bool) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetCallbacks(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, callbackOnPriceProposed, callbackOnPriceDisputed, callbackOnPriceSettled)
}

// SetCurrentTime is a paid mutator transaction binding the contract method 0x22f8e566.
//
// Solidity: function setCurrentTime(uint256 time) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SetCurrentTime(opts *bind.TransactOpts, time *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "setCurrentTime", time)
}

// SetCurrentTime is a paid mutator transaction binding the contract method 0x22f8e566.
//
// Solidity: function setCurrentTime(uint256 time) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SetCurrentTime(time *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetCurrentTime(&_UmaOptimisticOracleV2.TransactOpts, time)
}

// SetCurrentTime is a paid mutator transaction binding the contract method 0x22f8e566.
//
// Solidity: function setCurrentTime(uint256 time) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SetCurrentTime(time *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetCurrentTime(&_UmaOptimisticOracleV2.TransactOpts, time)
}

// SetCustomLiveness is a paid mutator transaction binding the contract method 0x473c45fe.
//
// Solidity: function setCustomLiveness(bytes32 identifier, uint256 timestamp, bytes ancillaryData, uint256 customLiveness) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SetCustomLiveness(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte, customLiveness *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "setCustomLiveness", identifier, timestamp, ancillaryData, customLiveness)
}

// SetCustomLiveness is a paid mutator transaction binding the contract method 0x473c45fe.
//
// Solidity: function setCustomLiveness(bytes32 identifier, uint256 timestamp, bytes ancillaryData, uint256 customLiveness) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SetCustomLiveness(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, customLiveness *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetCustomLiveness(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, customLiveness)
}

// SetCustomLiveness is a paid mutator transaction binding the contract method 0x473c45fe.
//
// Solidity: function setCustomLiveness(bytes32 identifier, uint256 timestamp, bytes ancillaryData, uint256 customLiveness) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SetCustomLiveness(identifier [32]byte, timestamp *big.Int, ancillaryData []byte, customLiveness *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetCustomLiveness(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData, customLiveness)
}

// SetEventBased is a paid mutator transaction binding the contract method 0x120698af.
//
// Solidity: function setEventBased(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SetEventBased(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "setEventBased", identifier, timestamp, ancillaryData)
}

// SetEventBased is a paid mutator transaction binding the contract method 0x120698af.
//
// Solidity: function setEventBased(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SetEventBased(identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetEventBased(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData)
}

// SetEventBased is a paid mutator transaction binding the contract method 0x120698af.
//
// Solidity: function setEventBased(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SetEventBased(identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetEventBased(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData)
}

// SetRefundOnDispute is a paid mutator transaction binding the contract method 0x91f58dcb.
//
// Solidity: function setRefundOnDispute(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SetRefundOnDispute(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "setRefundOnDispute", identifier, timestamp, ancillaryData)
}

// SetRefundOnDispute is a paid mutator transaction binding the contract method 0x91f58dcb.
//
// Solidity: function setRefundOnDispute(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SetRefundOnDispute(identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetRefundOnDispute(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData)
}

// SetRefundOnDispute is a paid mutator transaction binding the contract method 0x91f58dcb.
//
// Solidity: function setRefundOnDispute(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns()
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SetRefundOnDispute(identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SetRefundOnDispute(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData)
}

// Settle is a paid mutator transaction binding the contract method 0x5e9a79a9.
//
// Solidity: function settle(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 payout)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) Settle(opts *bind.TransactOpts, requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "settle", requester, identifier, timestamp, ancillaryData)
}

// Settle is a paid mutator transaction binding the contract method 0x5e9a79a9.
//
// Solidity: function settle(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 payout)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) Settle(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.Settle(&_UmaOptimisticOracleV2.TransactOpts, requester, identifier, timestamp, ancillaryData)
}

// Settle is a paid mutator transaction binding the contract method 0x5e9a79a9.
//
// Solidity: function settle(address requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(uint256 payout)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) Settle(requester common.Address, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.Settle(&_UmaOptimisticOracleV2.TransactOpts, requester, identifier, timestamp, ancillaryData)
}

// SettleAndGetPrice is a paid mutator transaction binding the contract method 0x53b59239.
//
// Solidity: function settleAndGetPrice(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(int256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Transactor) SettleAndGetPrice(opts *bind.TransactOpts, identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.contract.Transact(opts, "settleAndGetPrice", identifier, timestamp, ancillaryData)
}

// SettleAndGetPrice is a paid mutator transaction binding the contract method 0x53b59239.
//
// Solidity: function settleAndGetPrice(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(int256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Session) SettleAndGetPrice(identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SettleAndGetPrice(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData)
}

// SettleAndGetPrice is a paid mutator transaction binding the contract method 0x53b59239.
//
// Solidity: function settleAndGetPrice(bytes32 identifier, uint256 timestamp, bytes ancillaryData) returns(int256)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2TransactorSession) SettleAndGetPrice(identifier [32]byte, timestamp *big.Int, ancillaryData []byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV2.Contract.SettleAndGetPrice(&_UmaOptimisticOracleV2.TransactOpts, identifier, timestamp, ancillaryData)
}

// UmaOptimisticOracleV2DisputePriceIterator is returned from FilterDisputePrice and is used to iterate over the raw logs and unpacked data for DisputePrice events raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2DisputePriceIterator struct {
	Event *UmaOptimisticOracleV2DisputePrice // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV2DisputePriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV2DisputePrice)
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
		it.Event = new(UmaOptimisticOracleV2DisputePrice)
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
func (it *UmaOptimisticOracleV2DisputePriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV2DisputePriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV2DisputePrice represents a DisputePrice event raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2DisputePrice struct {
	Requester     common.Address
	Proposer      common.Address
	Disputer      common.Address
	Identifier    [32]byte
	Timestamp     *big.Int
	AncillaryData []byte
	ProposedPrice *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDisputePrice is a free log retrieval operation binding the contract event 0x5165909c3d1c01c5d1e121ac6f6d01dda1ba24bc9e1f975b5a375339c15be7f3.
//
// Solidity: event DisputePrice(address indexed requester, address indexed proposer, address indexed disputer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) FilterDisputePrice(opts *bind.FilterOpts, requester []common.Address, proposer []common.Address, disputer []common.Address) (*UmaOptimisticOracleV2DisputePriceIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}
	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.FilterLogs(opts, "DisputePrice", requesterRule, proposerRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2DisputePriceIterator{contract: _UmaOptimisticOracleV2.contract, event: "DisputePrice", logs: logs, sub: sub}, nil
}

// WatchDisputePrice is a free log subscription operation binding the contract event 0x5165909c3d1c01c5d1e121ac6f6d01dda1ba24bc9e1f975b5a375339c15be7f3.
//
// Solidity: event DisputePrice(address indexed requester, address indexed proposer, address indexed disputer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) WatchDisputePrice(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV2DisputePrice, requester []common.Address, proposer []common.Address, disputer []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}
	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.WatchLogs(opts, "DisputePrice", requesterRule, proposerRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV2DisputePrice)
				if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "DisputePrice", log); err != nil {
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

// ParseDisputePrice is a log parse operation binding the contract event 0x5165909c3d1c01c5d1e121ac6f6d01dda1ba24bc9e1f975b5a375339c15be7f3.
//
// Solidity: event DisputePrice(address indexed requester, address indexed proposer, address indexed disputer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) ParseDisputePrice(log types.Log) (*UmaOptimisticOracleV2DisputePrice, error) {
	event := new(UmaOptimisticOracleV2DisputePrice)
	if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "DisputePrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV2ProposePriceIterator is returned from FilterProposePrice and is used to iterate over the raw logs and unpacked data for ProposePrice events raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2ProposePriceIterator struct {
	Event *UmaOptimisticOracleV2ProposePrice // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV2ProposePriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV2ProposePrice)
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
		it.Event = new(UmaOptimisticOracleV2ProposePrice)
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
func (it *UmaOptimisticOracleV2ProposePriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV2ProposePriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV2ProposePrice represents a ProposePrice event raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2ProposePrice struct {
	Requester           common.Address
	Proposer            common.Address
	Identifier          [32]byte
	Timestamp           *big.Int
	AncillaryData       []byte
	ProposedPrice       *big.Int
	ExpirationTimestamp *big.Int
	Currency            common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterProposePrice is a free log retrieval operation binding the contract event 0x6e51dd00371aabffa82cd401592f76ed51e98a9ea4b58751c70463a2c78b5ca1.
//
// Solidity: event ProposePrice(address indexed requester, address indexed proposer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice, uint256 expirationTimestamp, address currency)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) FilterProposePrice(opts *bind.FilterOpts, requester []common.Address, proposer []common.Address) (*UmaOptimisticOracleV2ProposePriceIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.FilterLogs(opts, "ProposePrice", requesterRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2ProposePriceIterator{contract: _UmaOptimisticOracleV2.contract, event: "ProposePrice", logs: logs, sub: sub}, nil
}

// WatchProposePrice is a free log subscription operation binding the contract event 0x6e51dd00371aabffa82cd401592f76ed51e98a9ea4b58751c70463a2c78b5ca1.
//
// Solidity: event ProposePrice(address indexed requester, address indexed proposer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice, uint256 expirationTimestamp, address currency)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) WatchProposePrice(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV2ProposePrice, requester []common.Address, proposer []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.WatchLogs(opts, "ProposePrice", requesterRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV2ProposePrice)
				if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "ProposePrice", log); err != nil {
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

// ParseProposePrice is a log parse operation binding the contract event 0x6e51dd00371aabffa82cd401592f76ed51e98a9ea4b58751c70463a2c78b5ca1.
//
// Solidity: event ProposePrice(address indexed requester, address indexed proposer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 proposedPrice, uint256 expirationTimestamp, address currency)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) ParseProposePrice(log types.Log) (*UmaOptimisticOracleV2ProposePrice, error) {
	event := new(UmaOptimisticOracleV2ProposePrice)
	if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "ProposePrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV2RequestPriceIterator is returned from FilterRequestPrice and is used to iterate over the raw logs and unpacked data for RequestPrice events raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2RequestPriceIterator struct {
	Event *UmaOptimisticOracleV2RequestPrice // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV2RequestPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV2RequestPrice)
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
		it.Event = new(UmaOptimisticOracleV2RequestPrice)
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
func (it *UmaOptimisticOracleV2RequestPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV2RequestPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV2RequestPrice represents a RequestPrice event raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2RequestPrice struct {
	Requester     common.Address
	Identifier    [32]byte
	Timestamp     *big.Int
	AncillaryData []byte
	Currency      common.Address
	Reward        *big.Int
	FinalFee      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRequestPrice is a free log retrieval operation binding the contract event 0xf1679315ff325c257a944e0ca1bfe7b26616039e9511f9610d4ba3eca851027b.
//
// Solidity: event RequestPrice(address indexed requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, address currency, uint256 reward, uint256 finalFee)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) FilterRequestPrice(opts *bind.FilterOpts, requester []common.Address) (*UmaOptimisticOracleV2RequestPriceIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.FilterLogs(opts, "RequestPrice", requesterRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2RequestPriceIterator{contract: _UmaOptimisticOracleV2.contract, event: "RequestPrice", logs: logs, sub: sub}, nil
}

// WatchRequestPrice is a free log subscription operation binding the contract event 0xf1679315ff325c257a944e0ca1bfe7b26616039e9511f9610d4ba3eca851027b.
//
// Solidity: event RequestPrice(address indexed requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, address currency, uint256 reward, uint256 finalFee)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) WatchRequestPrice(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV2RequestPrice, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.WatchLogs(opts, "RequestPrice", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV2RequestPrice)
				if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "RequestPrice", log); err != nil {
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

// ParseRequestPrice is a log parse operation binding the contract event 0xf1679315ff325c257a944e0ca1bfe7b26616039e9511f9610d4ba3eca851027b.
//
// Solidity: event RequestPrice(address indexed requester, bytes32 identifier, uint256 timestamp, bytes ancillaryData, address currency, uint256 reward, uint256 finalFee)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) ParseRequestPrice(log types.Log) (*UmaOptimisticOracleV2RequestPrice, error) {
	event := new(UmaOptimisticOracleV2RequestPrice)
	if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "RequestPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV2SettleIterator is returned from FilterSettle and is used to iterate over the raw logs and unpacked data for Settle events raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2SettleIterator struct {
	Event *UmaOptimisticOracleV2Settle // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV2SettleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV2Settle)
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
		it.Event = new(UmaOptimisticOracleV2Settle)
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
func (it *UmaOptimisticOracleV2SettleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV2SettleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV2Settle represents a Settle event raised by the UmaOptimisticOracleV2 contract.
type UmaOptimisticOracleV2Settle struct {
	Requester     common.Address
	Proposer      common.Address
	Disputer      common.Address
	Identifier    [32]byte
	Timestamp     *big.Int
	AncillaryData []byte
	Price         *big.Int
	Payout        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSettle is a free log retrieval operation binding the contract event 0x3f384afb4bd9f0aef0298c80399950011420eb33b0e1a750b20966270247b9a0.
//
// Solidity: event Settle(address indexed requester, address indexed proposer, address indexed disputer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 price, uint256 payout)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) FilterSettle(opts *bind.FilterOpts, requester []common.Address, proposer []common.Address, disputer []common.Address) (*UmaOptimisticOracleV2SettleIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}
	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.FilterLogs(opts, "Settle", requesterRule, proposerRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV2SettleIterator{contract: _UmaOptimisticOracleV2.contract, event: "Settle", logs: logs, sub: sub}, nil
}

// WatchSettle is a free log subscription operation binding the contract event 0x3f384afb4bd9f0aef0298c80399950011420eb33b0e1a750b20966270247b9a0.
//
// Solidity: event Settle(address indexed requester, address indexed proposer, address indexed disputer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 price, uint256 payout)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) WatchSettle(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV2Settle, requester []common.Address, proposer []common.Address, disputer []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}
	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV2.contract.WatchLogs(opts, "Settle", requesterRule, proposerRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV2Settle)
				if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "Settle", log); err != nil {
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

// ParseSettle is a log parse operation binding the contract event 0x3f384afb4bd9f0aef0298c80399950011420eb33b0e1a750b20966270247b9a0.
//
// Solidity: event Settle(address indexed requester, address indexed proposer, address indexed disputer, bytes32 identifier, uint256 timestamp, bytes ancillaryData, int256 price, uint256 payout)
func (_UmaOptimisticOracleV2 *UmaOptimisticOracleV2Filterer) ParseSettle(log types.Log) (*UmaOptimisticOracleV2Settle, error) {
	event := new(UmaOptimisticOracleV2Settle)
	if err := _UmaOptimisticOracleV2.contract.UnpackLog(event, "Settle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
