// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polymarketUmaCtfAdapterV2

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

// AncillaryDataUpdate is an auto generated low-level Go binding around an user-defined struct.
type AncillaryDataUpdate struct {
	Timestamp *big.Int
	Update    []byte
}

// QuestionData is an auto generated low-level Go binding around an user-defined struct.
type QuestionData struct {
	RequestTimestamp             *big.Int
	Reward                       *big.Int
	ProposalBond                 *big.Int
	Liveness                     *big.Int
	EmergencyResolutionTimestamp *big.Int
	Resolved                     bool
	Paused                       bool
	Reset                        bool
	RewardToken                  common.Address
	Creator                      common.Address
	AncillaryData                []byte
}

// PolymarketUmaCtfAdapterV2MetaData contains all meta data concerning the PolymarketUmaCtfAdapterV2 contract.
var PolymarketUmaCtfAdapterV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ctf\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_finder\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Flagged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Initialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAncillaryData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidOOPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayouts\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotFlagged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOptimisticOracle\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotReadyToResolve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Paused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PriceNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Resolved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SafetyPeriodNotPassed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsupportedToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"}],\"name\":\"AncillaryDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdminAddress\",\"type\":\"address\"}],\"name\":\"NewAdmin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payouts\",\"type\":\"uint256[]\"}],\"name\":\"QuestionEmergencyResolved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"QuestionFlagged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalBond\",\"type\":\"uint256\"}],\"name\":\"QuestionInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"QuestionPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"QuestionReset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"settledPrice\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payouts\",\"type\":\"uint256[]\"}],\"name\":\"QuestionResolved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"QuestionUnpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"removedAdmin\",\"type\":\"address\"}],\"name\":\"RemovedAdmin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralWhitelist\",\"outputs\":[{\"internalType\":\"contractIAddressWhitelist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ctf\",\"outputs\":[{\"internalType\":\"contractIConditionalTokens\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"payouts\",\"type\":\"uint256[]\"}],\"name\":\"emergencyResolve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emergencySafetyPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"flag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"getExpectedPayouts\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getLatestUpdate\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"}],\"internalType\":\"structAncillaryDataUpdate\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"getQuestion\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proposalBond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liveness\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"emergencyResolutionTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"resolved\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reset\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"internalType\":\"structQuestionData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getUpdates\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"}],\"internalType\":\"structAncillaryDataUpdate[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proposalBond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liveness\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"isFlagged\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"isInitialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAncillaryData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"optimisticOracle\",\"outputs\":[{\"internalType\":\"contractIOptimisticOracleV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"}],\"name\":\"postUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"priceDisputed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"questions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proposalBond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liveness\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"emergencyResolutionTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"resolved\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reset\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"ancillaryData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"ready\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"resolve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"questionID\",\"type\":\"bytes32\"}],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"updates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"update\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"yesOrNoIdentifier\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PolymarketUmaCtfAdapterV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use PolymarketUmaCtfAdapterV2MetaData.ABI instead.
var PolymarketUmaCtfAdapterV2ABI = PolymarketUmaCtfAdapterV2MetaData.ABI

// PolymarketUmaCtfAdapterV2 is an auto generated Go binding around an Ethereum contract.
type PolymarketUmaCtfAdapterV2 struct {
	PolymarketUmaCtfAdapterV2Caller     // Read-only binding to the contract
	PolymarketUmaCtfAdapterV2Transactor // Write-only binding to the contract
	PolymarketUmaCtfAdapterV2Filterer   // Log filterer for contract events
}

// PolymarketUmaCtfAdapterV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type PolymarketUmaCtfAdapterV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolymarketUmaCtfAdapterV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type PolymarketUmaCtfAdapterV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolymarketUmaCtfAdapterV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolymarketUmaCtfAdapterV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolymarketUmaCtfAdapterV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolymarketUmaCtfAdapterV2Session struct {
	Contract     *PolymarketUmaCtfAdapterV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// PolymarketUmaCtfAdapterV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolymarketUmaCtfAdapterV2CallerSession struct {
	Contract *PolymarketUmaCtfAdapterV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// PolymarketUmaCtfAdapterV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolymarketUmaCtfAdapterV2TransactorSession struct {
	Contract     *PolymarketUmaCtfAdapterV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// PolymarketUmaCtfAdapterV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type PolymarketUmaCtfAdapterV2Raw struct {
	Contract *PolymarketUmaCtfAdapterV2 // Generic contract binding to access the raw methods on
}

// PolymarketUmaCtfAdapterV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolymarketUmaCtfAdapterV2CallerRaw struct {
	Contract *PolymarketUmaCtfAdapterV2Caller // Generic read-only contract binding to access the raw methods on
}

// PolymarketUmaCtfAdapterV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolymarketUmaCtfAdapterV2TransactorRaw struct {
	Contract *PolymarketUmaCtfAdapterV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewPolymarketUmaCtfAdapterV2 creates a new instance of PolymarketUmaCtfAdapterV2, bound to a specific deployed contract.
func NewPolymarketUmaCtfAdapterV2(address common.Address, backend bind.ContractBackend) (*PolymarketUmaCtfAdapterV2, error) {
	contract, err := bindPolymarketUmaCtfAdapterV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2{PolymarketUmaCtfAdapterV2Caller: PolymarketUmaCtfAdapterV2Caller{contract: contract}, PolymarketUmaCtfAdapterV2Transactor: PolymarketUmaCtfAdapterV2Transactor{contract: contract}, PolymarketUmaCtfAdapterV2Filterer: PolymarketUmaCtfAdapterV2Filterer{contract: contract}}, nil
}

// NewPolymarketUmaCtfAdapterV2Caller creates a new read-only instance of PolymarketUmaCtfAdapterV2, bound to a specific deployed contract.
func NewPolymarketUmaCtfAdapterV2Caller(address common.Address, caller bind.ContractCaller) (*PolymarketUmaCtfAdapterV2Caller, error) {
	contract, err := bindPolymarketUmaCtfAdapterV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2Caller{contract: contract}, nil
}

// NewPolymarketUmaCtfAdapterV2Transactor creates a new write-only instance of PolymarketUmaCtfAdapterV2, bound to a specific deployed contract.
func NewPolymarketUmaCtfAdapterV2Transactor(address common.Address, transactor bind.ContractTransactor) (*PolymarketUmaCtfAdapterV2Transactor, error) {
	contract, err := bindPolymarketUmaCtfAdapterV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2Transactor{contract: contract}, nil
}

// NewPolymarketUmaCtfAdapterV2Filterer creates a new log filterer instance of PolymarketUmaCtfAdapterV2, bound to a specific deployed contract.
func NewPolymarketUmaCtfAdapterV2Filterer(address common.Address, filterer bind.ContractFilterer) (*PolymarketUmaCtfAdapterV2Filterer, error) {
	contract, err := bindPolymarketUmaCtfAdapterV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2Filterer{contract: contract}, nil
}

// bindPolymarketUmaCtfAdapterV2 binds a generic wrapper to an already deployed contract.
func bindPolymarketUmaCtfAdapterV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PolymarketUmaCtfAdapterV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolymarketUmaCtfAdapterV2.Contract.PolymarketUmaCtfAdapterV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.PolymarketUmaCtfAdapterV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.PolymarketUmaCtfAdapterV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolymarketUmaCtfAdapterV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.contract.Transact(opts, method, params...)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) Admins(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "admins", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Admins(arg0 common.Address) (*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Admins(&_PolymarketUmaCtfAdapterV2.CallOpts, arg0)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Admins(&_PolymarketUmaCtfAdapterV2.CallOpts, arg0)
}

// CollateralWhitelist is a free data retrieval call binding the contract method 0xe4ee614a.
//
// Solidity: function collateralWhitelist() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) CollateralWhitelist(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "collateralWhitelist")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralWhitelist is a free data retrieval call binding the contract method 0xe4ee614a.
//
// Solidity: function collateralWhitelist() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) CollateralWhitelist() (common.Address, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.CollateralWhitelist(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// CollateralWhitelist is a free data retrieval call binding the contract method 0xe4ee614a.
//
// Solidity: function collateralWhitelist() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) CollateralWhitelist() (common.Address, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.CollateralWhitelist(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// Ctf is a free data retrieval call binding the contract method 0x22a9339f.
//
// Solidity: function ctf() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) Ctf(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "ctf")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ctf is a free data retrieval call binding the contract method 0x22a9339f.
//
// Solidity: function ctf() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Ctf() (common.Address, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Ctf(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// Ctf is a free data retrieval call binding the contract method 0x22a9339f.
//
// Solidity: function ctf() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) Ctf() (common.Address, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Ctf(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// EmergencySafetyPeriod is a free data retrieval call binding the contract method 0xc66d4c6c.
//
// Solidity: function emergencySafetyPeriod() view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) EmergencySafetyPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "emergencySafetyPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EmergencySafetyPeriod is a free data retrieval call binding the contract method 0xc66d4c6c.
//
// Solidity: function emergencySafetyPeriod() view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) EmergencySafetyPeriod() (*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.EmergencySafetyPeriod(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// EmergencySafetyPeriod is a free data retrieval call binding the contract method 0xc66d4c6c.
//
// Solidity: function emergencySafetyPeriod() view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) EmergencySafetyPeriod() (*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.EmergencySafetyPeriod(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// GetExpectedPayouts is a free data retrieval call binding the contract method 0x34e5e28e.
//
// Solidity: function getExpectedPayouts(bytes32 questionID) view returns(uint256[])
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) GetExpectedPayouts(opts *bind.CallOpts, questionID [32]byte) ([]*big.Int, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "getExpectedPayouts", questionID)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetExpectedPayouts is a free data retrieval call binding the contract method 0x34e5e28e.
//
// Solidity: function getExpectedPayouts(bytes32 questionID) view returns(uint256[])
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) GetExpectedPayouts(questionID [32]byte) ([]*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetExpectedPayouts(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// GetExpectedPayouts is a free data retrieval call binding the contract method 0x34e5e28e.
//
// Solidity: function getExpectedPayouts(bytes32 questionID) view returns(uint256[])
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) GetExpectedPayouts(questionID [32]byte) ([]*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetExpectedPayouts(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// GetLatestUpdate is a free data retrieval call binding the contract method 0xc0cab0a2.
//
// Solidity: function getLatestUpdate(bytes32 questionID, address owner) view returns((uint256,bytes))
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) GetLatestUpdate(opts *bind.CallOpts, questionID [32]byte, owner common.Address) (AncillaryDataUpdate, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "getLatestUpdate", questionID, owner)

	if err != nil {
		return *new(AncillaryDataUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(AncillaryDataUpdate)).(*AncillaryDataUpdate)

	return out0, err

}

// GetLatestUpdate is a free data retrieval call binding the contract method 0xc0cab0a2.
//
// Solidity: function getLatestUpdate(bytes32 questionID, address owner) view returns((uint256,bytes))
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) GetLatestUpdate(questionID [32]byte, owner common.Address) (AncillaryDataUpdate, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetLatestUpdate(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID, owner)
}

// GetLatestUpdate is a free data retrieval call binding the contract method 0xc0cab0a2.
//
// Solidity: function getLatestUpdate(bytes32 questionID, address owner) view returns((uint256,bytes))
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) GetLatestUpdate(questionID [32]byte, owner common.Address) (AncillaryDataUpdate, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetLatestUpdate(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID, owner)
}

// GetQuestion is a free data retrieval call binding the contract method 0x58c039cd.
//
// Solidity: function getQuestion(bytes32 questionID) view returns((uint256,uint256,uint256,uint256,uint256,bool,bool,bool,address,address,bytes))
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) GetQuestion(opts *bind.CallOpts, questionID [32]byte) (QuestionData, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "getQuestion", questionID)

	if err != nil {
		return *new(QuestionData), err
	}

	out0 := *abi.ConvertType(out[0], new(QuestionData)).(*QuestionData)

	return out0, err

}

// GetQuestion is a free data retrieval call binding the contract method 0x58c039cd.
//
// Solidity: function getQuestion(bytes32 questionID) view returns((uint256,uint256,uint256,uint256,uint256,bool,bool,bool,address,address,bytes))
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) GetQuestion(questionID [32]byte) (QuestionData, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetQuestion(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// GetQuestion is a free data retrieval call binding the contract method 0x58c039cd.
//
// Solidity: function getQuestion(bytes32 questionID) view returns((uint256,uint256,uint256,uint256,uint256,bool,bool,bool,address,address,bytes))
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) GetQuestion(questionID [32]byte) (QuestionData, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetQuestion(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// GetUpdates is a free data retrieval call binding the contract method 0x555c56fc.
//
// Solidity: function getUpdates(bytes32 questionID, address owner) view returns((uint256,bytes)[])
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) GetUpdates(opts *bind.CallOpts, questionID [32]byte, owner common.Address) ([]AncillaryDataUpdate, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "getUpdates", questionID, owner)

	if err != nil {
		return *new([]AncillaryDataUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new([]AncillaryDataUpdate)).(*[]AncillaryDataUpdate)

	return out0, err

}

// GetUpdates is a free data retrieval call binding the contract method 0x555c56fc.
//
// Solidity: function getUpdates(bytes32 questionID, address owner) view returns((uint256,bytes)[])
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) GetUpdates(questionID [32]byte, owner common.Address) ([]AncillaryDataUpdate, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetUpdates(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID, owner)
}

// GetUpdates is a free data retrieval call binding the contract method 0x555c56fc.
//
// Solidity: function getUpdates(bytes32 questionID, address owner) view returns((uint256,bytes)[])
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) GetUpdates(questionID [32]byte, owner common.Address) ([]AncillaryDataUpdate, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.GetUpdates(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID, owner)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) IsAdmin(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "isAdmin", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) IsAdmin(addr common.Address) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.IsAdmin(&_PolymarketUmaCtfAdapterV2.CallOpts, addr)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) IsAdmin(addr common.Address) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.IsAdmin(&_PolymarketUmaCtfAdapterV2.CallOpts, addr)
}

// IsFlagged is a free data retrieval call binding the contract method 0xbf2dde38.
//
// Solidity: function isFlagged(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) IsFlagged(opts *bind.CallOpts, questionID [32]byte) (bool, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "isFlagged", questionID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFlagged is a free data retrieval call binding the contract method 0xbf2dde38.
//
// Solidity: function isFlagged(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) IsFlagged(questionID [32]byte) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.IsFlagged(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// IsFlagged is a free data retrieval call binding the contract method 0xbf2dde38.
//
// Solidity: function isFlagged(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) IsFlagged(questionID [32]byte) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.IsFlagged(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// IsInitialized is a free data retrieval call binding the contract method 0xf7b637bb.
//
// Solidity: function isInitialized(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) IsInitialized(opts *bind.CallOpts, questionID [32]byte) (bool, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "isInitialized", questionID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsInitialized is a free data retrieval call binding the contract method 0xf7b637bb.
//
// Solidity: function isInitialized(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) IsInitialized(questionID [32]byte) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.IsInitialized(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// IsInitialized is a free data retrieval call binding the contract method 0xf7b637bb.
//
// Solidity: function isInitialized(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) IsInitialized(questionID [32]byte) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.IsInitialized(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// MaxAncillaryData is a free data retrieval call binding the contract method 0x6b5acc63.
//
// Solidity: function maxAncillaryData() view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) MaxAncillaryData(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "maxAncillaryData")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAncillaryData is a free data retrieval call binding the contract method 0x6b5acc63.
//
// Solidity: function maxAncillaryData() view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) MaxAncillaryData() (*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.MaxAncillaryData(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// MaxAncillaryData is a free data retrieval call binding the contract method 0x6b5acc63.
//
// Solidity: function maxAncillaryData() view returns(uint256)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) MaxAncillaryData() (*big.Int, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.MaxAncillaryData(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// OptimisticOracle is a free data retrieval call binding the contract method 0x22302922.
//
// Solidity: function optimisticOracle() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) OptimisticOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "optimisticOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OptimisticOracle is a free data retrieval call binding the contract method 0x22302922.
//
// Solidity: function optimisticOracle() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) OptimisticOracle() (common.Address, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.OptimisticOracle(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// OptimisticOracle is a free data retrieval call binding the contract method 0x22302922.
//
// Solidity: function optimisticOracle() view returns(address)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) OptimisticOracle() (common.Address, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.OptimisticOracle(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// Questions is a free data retrieval call binding the contract method 0x95addb90.
//
// Solidity: function questions(bytes32 ) view returns(uint256 requestTimestamp, uint256 reward, uint256 proposalBond, uint256 liveness, uint256 emergencyResolutionTimestamp, bool resolved, bool paused, bool reset, address rewardToken, address creator, bytes ancillaryData)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) Questions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	RequestTimestamp             *big.Int
	Reward                       *big.Int
	ProposalBond                 *big.Int
	Liveness                     *big.Int
	EmergencyResolutionTimestamp *big.Int
	Resolved                     bool
	Paused                       bool
	Reset                        bool
	RewardToken                  common.Address
	Creator                      common.Address
	AncillaryData                []byte
}, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "questions", arg0)

	outstruct := new(struct {
		RequestTimestamp             *big.Int
		Reward                       *big.Int
		ProposalBond                 *big.Int
		Liveness                     *big.Int
		EmergencyResolutionTimestamp *big.Int
		Resolved                     bool
		Paused                       bool
		Reset                        bool
		RewardToken                  common.Address
		Creator                      common.Address
		AncillaryData                []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestTimestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reward = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ProposalBond = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Liveness = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.EmergencyResolutionTimestamp = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Resolved = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.Paused = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.Reset = *abi.ConvertType(out[7], new(bool)).(*bool)
	outstruct.RewardToken = *abi.ConvertType(out[8], new(common.Address)).(*common.Address)
	outstruct.Creator = *abi.ConvertType(out[9], new(common.Address)).(*common.Address)
	outstruct.AncillaryData = *abi.ConvertType(out[10], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Questions is a free data retrieval call binding the contract method 0x95addb90.
//
// Solidity: function questions(bytes32 ) view returns(uint256 requestTimestamp, uint256 reward, uint256 proposalBond, uint256 liveness, uint256 emergencyResolutionTimestamp, bool resolved, bool paused, bool reset, address rewardToken, address creator, bytes ancillaryData)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Questions(arg0 [32]byte) (struct {
	RequestTimestamp             *big.Int
	Reward                       *big.Int
	ProposalBond                 *big.Int
	Liveness                     *big.Int
	EmergencyResolutionTimestamp *big.Int
	Resolved                     bool
	Paused                       bool
	Reset                        bool
	RewardToken                  common.Address
	Creator                      common.Address
	AncillaryData                []byte
}, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Questions(&_PolymarketUmaCtfAdapterV2.CallOpts, arg0)
}

// Questions is a free data retrieval call binding the contract method 0x95addb90.
//
// Solidity: function questions(bytes32 ) view returns(uint256 requestTimestamp, uint256 reward, uint256 proposalBond, uint256 liveness, uint256 emergencyResolutionTimestamp, bool resolved, bool paused, bool reset, address rewardToken, address creator, bytes ancillaryData)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) Questions(arg0 [32]byte) (struct {
	RequestTimestamp             *big.Int
	Reward                       *big.Int
	ProposalBond                 *big.Int
	Liveness                     *big.Int
	EmergencyResolutionTimestamp *big.Int
	Resolved                     bool
	Paused                       bool
	Reset                        bool
	RewardToken                  common.Address
	Creator                      common.Address
	AncillaryData                []byte
}, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Questions(&_PolymarketUmaCtfAdapterV2.CallOpts, arg0)
}

// Ready is a free data retrieval call binding the contract method 0xfcac49a2.
//
// Solidity: function ready(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) Ready(opts *bind.CallOpts, questionID [32]byte) (bool, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "ready", questionID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Ready is a free data retrieval call binding the contract method 0xfcac49a2.
//
// Solidity: function ready(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Ready(questionID [32]byte) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Ready(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// Ready is a free data retrieval call binding the contract method 0xfcac49a2.
//
// Solidity: function ready(bytes32 questionID) view returns(bool)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) Ready(questionID [32]byte) (bool, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Ready(&_PolymarketUmaCtfAdapterV2.CallOpts, questionID)
}

// Updates is a free data retrieval call binding the contract method 0x89ab0871.
//
// Solidity: function updates(bytes32 , uint256 ) view returns(uint256 timestamp, bytes update)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) Updates(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (struct {
	Timestamp *big.Int
	Update    []byte
}, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "updates", arg0, arg1)

	outstruct := new(struct {
		Timestamp *big.Int
		Update    []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Update = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Updates is a free data retrieval call binding the contract method 0x89ab0871.
//
// Solidity: function updates(bytes32 , uint256 ) view returns(uint256 timestamp, bytes update)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Updates(arg0 [32]byte, arg1 *big.Int) (struct {
	Timestamp *big.Int
	Update    []byte
}, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Updates(&_PolymarketUmaCtfAdapterV2.CallOpts, arg0, arg1)
}

// Updates is a free data retrieval call binding the contract method 0x89ab0871.
//
// Solidity: function updates(bytes32 , uint256 ) view returns(uint256 timestamp, bytes update)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) Updates(arg0 [32]byte, arg1 *big.Int) (struct {
	Timestamp *big.Int
	Update    []byte
}, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Updates(&_PolymarketUmaCtfAdapterV2.CallOpts, arg0, arg1)
}

// YesOrNoIdentifier is a free data retrieval call binding the contract method 0xdddb4680.
//
// Solidity: function yesOrNoIdentifier() view returns(bytes32)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Caller) YesOrNoIdentifier(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PolymarketUmaCtfAdapterV2.contract.Call(opts, &out, "yesOrNoIdentifier")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// YesOrNoIdentifier is a free data retrieval call binding the contract method 0xdddb4680.
//
// Solidity: function yesOrNoIdentifier() view returns(bytes32)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) YesOrNoIdentifier() ([32]byte, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.YesOrNoIdentifier(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// YesOrNoIdentifier is a free data retrieval call binding the contract method 0xdddb4680.
//
// Solidity: function yesOrNoIdentifier() view returns(bytes32)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2CallerSession) YesOrNoIdentifier() ([32]byte, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.YesOrNoIdentifier(&_PolymarketUmaCtfAdapterV2.CallOpts)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) AddAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "addAdmin", admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.AddAdmin(&_PolymarketUmaCtfAdapterV2.TransactOpts, admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.AddAdmin(&_PolymarketUmaCtfAdapterV2.TransactOpts, admin)
}

// EmergencyResolve is a paid mutator transaction binding the contract method 0x9ce7c0e0.
//
// Solidity: function emergencyResolve(bytes32 questionID, uint256[] payouts) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) EmergencyResolve(opts *bind.TransactOpts, questionID [32]byte, payouts []*big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "emergencyResolve", questionID, payouts)
}

// EmergencyResolve is a paid mutator transaction binding the contract method 0x9ce7c0e0.
//
// Solidity: function emergencyResolve(bytes32 questionID, uint256[] payouts) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) EmergencyResolve(questionID [32]byte, payouts []*big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.EmergencyResolve(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID, payouts)
}

// EmergencyResolve is a paid mutator transaction binding the contract method 0x9ce7c0e0.
//
// Solidity: function emergencyResolve(bytes32 questionID, uint256[] payouts) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) EmergencyResolve(questionID [32]byte, payouts []*big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.EmergencyResolve(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID, payouts)
}

// Flag is a paid mutator transaction binding the contract method 0x78165a48.
//
// Solidity: function flag(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) Flag(opts *bind.TransactOpts, questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "flag", questionID)
}

// Flag is a paid mutator transaction binding the contract method 0x78165a48.
//
// Solidity: function flag(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Flag(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Flag(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Flag is a paid mutator transaction binding the contract method 0x78165a48.
//
// Solidity: function flag(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) Flag(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Flag(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Initialize is a paid mutator transaction binding the contract method 0x185d1646.
//
// Solidity: function initialize(bytes ancillaryData, address rewardToken, uint256 reward, uint256 proposalBond, uint256 liveness) returns(bytes32 questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) Initialize(opts *bind.TransactOpts, ancillaryData []byte, rewardToken common.Address, reward *big.Int, proposalBond *big.Int, liveness *big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "initialize", ancillaryData, rewardToken, reward, proposalBond, liveness)
}

// Initialize is a paid mutator transaction binding the contract method 0x185d1646.
//
// Solidity: function initialize(bytes ancillaryData, address rewardToken, uint256 reward, uint256 proposalBond, uint256 liveness) returns(bytes32 questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Initialize(ancillaryData []byte, rewardToken common.Address, reward *big.Int, proposalBond *big.Int, liveness *big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Initialize(&_PolymarketUmaCtfAdapterV2.TransactOpts, ancillaryData, rewardToken, reward, proposalBond, liveness)
}

// Initialize is a paid mutator transaction binding the contract method 0x185d1646.
//
// Solidity: function initialize(bytes ancillaryData, address rewardToken, uint256 reward, uint256 proposalBond, uint256 liveness) returns(bytes32 questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) Initialize(ancillaryData []byte, rewardToken common.Address, reward *big.Int, proposalBond *big.Int, liveness *big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Initialize(&_PolymarketUmaCtfAdapterV2.TransactOpts, ancillaryData, rewardToken, reward, proposalBond, liveness)
}

// Pause is a paid mutator transaction binding the contract method 0xed56531a.
//
// Solidity: function pause(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) Pause(opts *bind.TransactOpts, questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "pause", questionID)
}

// Pause is a paid mutator transaction binding the contract method 0xed56531a.
//
// Solidity: function pause(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Pause(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Pause(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Pause is a paid mutator transaction binding the contract method 0xed56531a.
//
// Solidity: function pause(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) Pause(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Pause(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// PostUpdate is a paid mutator transaction binding the contract method 0x072d1259.
//
// Solidity: function postUpdate(bytes32 questionID, bytes update) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) PostUpdate(opts *bind.TransactOpts, questionID [32]byte, update []byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "postUpdate", questionID, update)
}

// PostUpdate is a paid mutator transaction binding the contract method 0x072d1259.
//
// Solidity: function postUpdate(bytes32 questionID, bytes update) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) PostUpdate(questionID [32]byte, update []byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.PostUpdate(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID, update)
}

// PostUpdate is a paid mutator transaction binding the contract method 0x072d1259.
//
// Solidity: function postUpdate(bytes32 questionID, bytes update) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) PostUpdate(questionID [32]byte, update []byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.PostUpdate(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID, update)
}

// PriceDisputed is a paid mutator transaction binding the contract method 0x0d8f2372.
//
// Solidity: function priceDisputed(bytes32 , uint256 , bytes ancillaryData, uint256 ) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) PriceDisputed(opts *bind.TransactOpts, arg0 [32]byte, arg1 *big.Int, ancillaryData []byte, arg3 *big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "priceDisputed", arg0, arg1, ancillaryData, arg3)
}

// PriceDisputed is a paid mutator transaction binding the contract method 0x0d8f2372.
//
// Solidity: function priceDisputed(bytes32 , uint256 , bytes ancillaryData, uint256 ) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) PriceDisputed(arg0 [32]byte, arg1 *big.Int, ancillaryData []byte, arg3 *big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.PriceDisputed(&_PolymarketUmaCtfAdapterV2.TransactOpts, arg0, arg1, ancillaryData, arg3)
}

// PriceDisputed is a paid mutator transaction binding the contract method 0x0d8f2372.
//
// Solidity: function priceDisputed(bytes32 , uint256 , bytes ancillaryData, uint256 ) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) PriceDisputed(arg0 [32]byte, arg1 *big.Int, ancillaryData []byte, arg3 *big.Int) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.PriceDisputed(&_PolymarketUmaCtfAdapterV2.TransactOpts, arg0, arg1, ancillaryData, arg3)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) RemoveAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "removeAdmin", admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.RemoveAdmin(&_PolymarketUmaCtfAdapterV2.TransactOpts, admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.RemoveAdmin(&_PolymarketUmaCtfAdapterV2.TransactOpts, admin)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x8bad0c0a.
//
// Solidity: function renounceAdmin() returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) RenounceAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "renounceAdmin")
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x8bad0c0a.
//
// Solidity: function renounceAdmin() returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) RenounceAdmin() (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.RenounceAdmin(&_PolymarketUmaCtfAdapterV2.TransactOpts)
}

// RenounceAdmin is a paid mutator transaction binding the contract method 0x8bad0c0a.
//
// Solidity: function renounceAdmin() returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) RenounceAdmin() (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.RenounceAdmin(&_PolymarketUmaCtfAdapterV2.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xed3c7d40.
//
// Solidity: function reset(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) Reset(opts *bind.TransactOpts, questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "reset", questionID)
}

// Reset is a paid mutator transaction binding the contract method 0xed3c7d40.
//
// Solidity: function reset(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Reset(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Reset(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Reset is a paid mutator transaction binding the contract method 0xed3c7d40.
//
// Solidity: function reset(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) Reset(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Reset(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Resolve is a paid mutator transaction binding the contract method 0x5c23bdf5.
//
// Solidity: function resolve(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) Resolve(opts *bind.TransactOpts, questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "resolve", questionID)
}

// Resolve is a paid mutator transaction binding the contract method 0x5c23bdf5.
//
// Solidity: function resolve(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Resolve(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Resolve(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Resolve is a paid mutator transaction binding the contract method 0x5c23bdf5.
//
// Solidity: function resolve(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) Resolve(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Resolve(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Unpause is a paid mutator transaction binding the contract method 0x2f4dae9f.
//
// Solidity: function unpause(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Transactor) Unpause(opts *bind.TransactOpts, questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.contract.Transact(opts, "unpause", questionID)
}

// Unpause is a paid mutator transaction binding the contract method 0x2f4dae9f.
//
// Solidity: function unpause(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Session) Unpause(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Unpause(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// Unpause is a paid mutator transaction binding the contract method 0x2f4dae9f.
//
// Solidity: function unpause(bytes32 questionID) returns()
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2TransactorSession) Unpause(questionID [32]byte) (*types.Transaction, error) {
	return _PolymarketUmaCtfAdapterV2.Contract.Unpause(&_PolymarketUmaCtfAdapterV2.TransactOpts, questionID)
}

// PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator is returned from FilterAncillaryDataUpdated and is used to iterate over the raw logs and unpacked data for AncillaryDataUpdated events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator struct {
	Event *PolymarketUmaCtfAdapterV2AncillaryDataUpdated // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2AncillaryDataUpdated)
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
		it.Event = new(PolymarketUmaCtfAdapterV2AncillaryDataUpdated)
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
func (it *PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2AncillaryDataUpdated represents a AncillaryDataUpdated event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2AncillaryDataUpdated struct {
	QuestionID [32]byte
	Owner      common.Address
	Update     []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAncillaryDataUpdated is a free log retrieval operation binding the contract event 0x0059e11815211969c0c4aaf3f498b52b6c2f2d14f286275d0862d70de22a836b.
//
// Solidity: event AncillaryDataUpdated(bytes32 indexed questionID, address indexed owner, bytes update)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterAncillaryDataUpdated(opts *bind.FilterOpts, questionID [][32]byte, owner []common.Address) (*PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "AncillaryDataUpdated", questionIDRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2AncillaryDataUpdatedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "AncillaryDataUpdated", logs: logs, sub: sub}, nil
}

// WatchAncillaryDataUpdated is a free log subscription operation binding the contract event 0x0059e11815211969c0c4aaf3f498b52b6c2f2d14f286275d0862d70de22a836b.
//
// Solidity: event AncillaryDataUpdated(bytes32 indexed questionID, address indexed owner, bytes update)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchAncillaryDataUpdated(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2AncillaryDataUpdated, questionID [][32]byte, owner []common.Address) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "AncillaryDataUpdated", questionIDRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2AncillaryDataUpdated)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "AncillaryDataUpdated", log); err != nil {
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

// ParseAncillaryDataUpdated is a log parse operation binding the contract event 0x0059e11815211969c0c4aaf3f498b52b6c2f2d14f286275d0862d70de22a836b.
//
// Solidity: event AncillaryDataUpdated(bytes32 indexed questionID, address indexed owner, bytes update)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseAncillaryDataUpdated(log types.Log) (*PolymarketUmaCtfAdapterV2AncillaryDataUpdated, error) {
	event := new(PolymarketUmaCtfAdapterV2AncillaryDataUpdated)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "AncillaryDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2NewAdminIterator is returned from FilterNewAdmin and is used to iterate over the raw logs and unpacked data for NewAdmin events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2NewAdminIterator struct {
	Event *PolymarketUmaCtfAdapterV2NewAdmin // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2NewAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2NewAdmin)
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
		it.Event = new(PolymarketUmaCtfAdapterV2NewAdmin)
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
func (it *PolymarketUmaCtfAdapterV2NewAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2NewAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2NewAdmin represents a NewAdmin event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2NewAdmin struct {
	Admin           common.Address
	NewAdminAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewAdmin is a free log retrieval operation binding the contract event 0xf9ffabca9c8276e99321725bcb43fb076a6c66a54b7f21c4e8146d8519b417dc.
//
// Solidity: event NewAdmin(address indexed admin, address indexed newAdminAddress)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterNewAdmin(opts *bind.FilterOpts, admin []common.Address, newAdminAddress []common.Address) (*PolymarketUmaCtfAdapterV2NewAdminIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var newAdminAddressRule []interface{}
	for _, newAdminAddressItem := range newAdminAddress {
		newAdminAddressRule = append(newAdminAddressRule, newAdminAddressItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "NewAdmin", adminRule, newAdminAddressRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2NewAdminIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "NewAdmin", logs: logs, sub: sub}, nil
}

// WatchNewAdmin is a free log subscription operation binding the contract event 0xf9ffabca9c8276e99321725bcb43fb076a6c66a54b7f21c4e8146d8519b417dc.
//
// Solidity: event NewAdmin(address indexed admin, address indexed newAdminAddress)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchNewAdmin(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2NewAdmin, admin []common.Address, newAdminAddress []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var newAdminAddressRule []interface{}
	for _, newAdminAddressItem := range newAdminAddress {
		newAdminAddressRule = append(newAdminAddressRule, newAdminAddressItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "NewAdmin", adminRule, newAdminAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2NewAdmin)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "NewAdmin", log); err != nil {
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
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseNewAdmin(log types.Log) (*PolymarketUmaCtfAdapterV2NewAdmin, error) {
	event := new(PolymarketUmaCtfAdapterV2NewAdmin)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "NewAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator is returned from FilterQuestionEmergencyResolved and is used to iterate over the raw logs and unpacked data for QuestionEmergencyResolved events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionEmergencyResolved // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionEmergencyResolved)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionEmergencyResolved)
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
func (it *PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionEmergencyResolved represents a QuestionEmergencyResolved event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionEmergencyResolved struct {
	QuestionID [32]byte
	Payouts    []*big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterQuestionEmergencyResolved is a free log retrieval operation binding the contract event 0x6edb5841a476c9c29c34a652d1a44f785fe71a6157a3da9a6a6a589a1bd2945a.
//
// Solidity: event QuestionEmergencyResolved(bytes32 indexed questionID, uint256[] payouts)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionEmergencyResolved(opts *bind.FilterOpts, questionID [][32]byte) (*PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionEmergencyResolved", questionIDRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionEmergencyResolvedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionEmergencyResolved", logs: logs, sub: sub}, nil
}

// WatchQuestionEmergencyResolved is a free log subscription operation binding the contract event 0x6edb5841a476c9c29c34a652d1a44f785fe71a6157a3da9a6a6a589a1bd2945a.
//
// Solidity: event QuestionEmergencyResolved(bytes32 indexed questionID, uint256[] payouts)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionEmergencyResolved(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionEmergencyResolved, questionID [][32]byte) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionEmergencyResolved", questionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionEmergencyResolved)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionEmergencyResolved", log); err != nil {
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

// ParseQuestionEmergencyResolved is a log parse operation binding the contract event 0x6edb5841a476c9c29c34a652d1a44f785fe71a6157a3da9a6a6a589a1bd2945a.
//
// Solidity: event QuestionEmergencyResolved(bytes32 indexed questionID, uint256[] payouts)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionEmergencyResolved(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionEmergencyResolved, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionEmergencyResolved)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionEmergencyResolved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionFlaggedIterator is returned from FilterQuestionFlagged and is used to iterate over the raw logs and unpacked data for QuestionFlagged events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionFlaggedIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionFlagged // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionFlaggedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionFlagged)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionFlagged)
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
func (it *PolymarketUmaCtfAdapterV2QuestionFlaggedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionFlaggedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionFlagged represents a QuestionFlagged event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionFlagged struct {
	QuestionID [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterQuestionFlagged is a free log retrieval operation binding the contract event 0x2435a0347185933b12027c6f394a5fd9c03646dba233e956f50658719dfc0b35.
//
// Solidity: event QuestionFlagged(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionFlagged(opts *bind.FilterOpts, questionID [][32]byte) (*PolymarketUmaCtfAdapterV2QuestionFlaggedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionFlagged", questionIDRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionFlaggedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionFlagged", logs: logs, sub: sub}, nil
}

// WatchQuestionFlagged is a free log subscription operation binding the contract event 0x2435a0347185933b12027c6f394a5fd9c03646dba233e956f50658719dfc0b35.
//
// Solidity: event QuestionFlagged(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionFlagged(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionFlagged, questionID [][32]byte) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionFlagged", questionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionFlagged)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionFlagged", log); err != nil {
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

// ParseQuestionFlagged is a log parse operation binding the contract event 0x2435a0347185933b12027c6f394a5fd9c03646dba233e956f50658719dfc0b35.
//
// Solidity: event QuestionFlagged(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionFlagged(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionFlagged, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionFlagged)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionFlagged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionInitializedIterator is returned from FilterQuestionInitialized and is used to iterate over the raw logs and unpacked data for QuestionInitialized events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionInitializedIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionInitialized // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionInitialized)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionInitialized)
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
func (it *PolymarketUmaCtfAdapterV2QuestionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionInitialized represents a QuestionInitialized event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionInitialized struct {
	QuestionID       [32]byte
	RequestTimestamp *big.Int
	Creator          common.Address
	AncillaryData    []byte
	RewardToken      common.Address
	Reward           *big.Int
	ProposalBond     *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterQuestionInitialized is a free log retrieval operation binding the contract event 0xeee0897acd6893adcaf2ba5158191b3601098ab6bece35c5d57874340b64c5b7.
//
// Solidity: event QuestionInitialized(bytes32 indexed questionID, uint256 indexed requestTimestamp, address indexed creator, bytes ancillaryData, address rewardToken, uint256 reward, uint256 proposalBond)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionInitialized(opts *bind.FilterOpts, questionID [][32]byte, requestTimestamp []*big.Int, creator []common.Address) (*PolymarketUmaCtfAdapterV2QuestionInitializedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}
	var requestTimestampRule []interface{}
	for _, requestTimestampItem := range requestTimestamp {
		requestTimestampRule = append(requestTimestampRule, requestTimestampItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionInitialized", questionIDRule, requestTimestampRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionInitializedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionInitialized", logs: logs, sub: sub}, nil
}

// WatchQuestionInitialized is a free log subscription operation binding the contract event 0xeee0897acd6893adcaf2ba5158191b3601098ab6bece35c5d57874340b64c5b7.
//
// Solidity: event QuestionInitialized(bytes32 indexed questionID, uint256 indexed requestTimestamp, address indexed creator, bytes ancillaryData, address rewardToken, uint256 reward, uint256 proposalBond)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionInitialized(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionInitialized, questionID [][32]byte, requestTimestamp []*big.Int, creator []common.Address) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}
	var requestTimestampRule []interface{}
	for _, requestTimestampItem := range requestTimestamp {
		requestTimestampRule = append(requestTimestampRule, requestTimestampItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionInitialized", questionIDRule, requestTimestampRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionInitialized)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionInitialized", log); err != nil {
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

// ParseQuestionInitialized is a log parse operation binding the contract event 0xeee0897acd6893adcaf2ba5158191b3601098ab6bece35c5d57874340b64c5b7.
//
// Solidity: event QuestionInitialized(bytes32 indexed questionID, uint256 indexed requestTimestamp, address indexed creator, bytes ancillaryData, address rewardToken, uint256 reward, uint256 proposalBond)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionInitialized(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionInitialized, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionInitialized)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionPausedIterator is returned from FilterQuestionPaused and is used to iterate over the raw logs and unpacked data for QuestionPaused events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionPausedIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionPaused // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionPaused)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionPaused)
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
func (it *PolymarketUmaCtfAdapterV2QuestionPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionPaused represents a QuestionPaused event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionPaused struct {
	QuestionID [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterQuestionPaused is a free log retrieval operation binding the contract event 0x6ded7250a9d5f79aef5add44600fc20a74a0af6f4730baa4fc4ab87bf484b812.
//
// Solidity: event QuestionPaused(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionPaused(opts *bind.FilterOpts, questionID [][32]byte) (*PolymarketUmaCtfAdapterV2QuestionPausedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionPaused", questionIDRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionPausedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionPaused", logs: logs, sub: sub}, nil
}

// WatchQuestionPaused is a free log subscription operation binding the contract event 0x6ded7250a9d5f79aef5add44600fc20a74a0af6f4730baa4fc4ab87bf484b812.
//
// Solidity: event QuestionPaused(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionPaused(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionPaused, questionID [][32]byte) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionPaused", questionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionPaused)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionPaused", log); err != nil {
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

// ParseQuestionPaused is a log parse operation binding the contract event 0x6ded7250a9d5f79aef5add44600fc20a74a0af6f4730baa4fc4ab87bf484b812.
//
// Solidity: event QuestionPaused(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionPaused(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionPaused, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionPaused)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionResetIterator is returned from FilterQuestionReset and is used to iterate over the raw logs and unpacked data for QuestionReset events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionResetIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionReset // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionResetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionReset)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionReset)
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
func (it *PolymarketUmaCtfAdapterV2QuestionResetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionResetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionReset represents a QuestionReset event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionReset struct {
	QuestionID [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterQuestionReset is a free log retrieval operation binding the contract event 0x7981b5832932948db4e32a4a16a0f44b2ce7ff088574afb9364b313f70f82e8f.
//
// Solidity: event QuestionReset(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionReset(opts *bind.FilterOpts, questionID [][32]byte) (*PolymarketUmaCtfAdapterV2QuestionResetIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionReset", questionIDRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionResetIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionReset", logs: logs, sub: sub}, nil
}

// WatchQuestionReset is a free log subscription operation binding the contract event 0x7981b5832932948db4e32a4a16a0f44b2ce7ff088574afb9364b313f70f82e8f.
//
// Solidity: event QuestionReset(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionReset(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionReset, questionID [][32]byte) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionReset", questionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionReset)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionReset", log); err != nil {
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

// ParseQuestionReset is a log parse operation binding the contract event 0x7981b5832932948db4e32a4a16a0f44b2ce7ff088574afb9364b313f70f82e8f.
//
// Solidity: event QuestionReset(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionReset(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionReset, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionReset)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionReset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionResolvedIterator is returned from FilterQuestionResolved and is used to iterate over the raw logs and unpacked data for QuestionResolved events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionResolvedIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionResolved // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionResolved)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionResolved)
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
func (it *PolymarketUmaCtfAdapterV2QuestionResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionResolved represents a QuestionResolved event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionResolved struct {
	QuestionID   [32]byte
	SettledPrice *big.Int
	Payouts      []*big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterQuestionResolved is a free log retrieval operation binding the contract event 0x566c3fbdd12dd86bb341787f6d531f79fd7ad4ce7e3ae2d15ac0ca1b601af9df.
//
// Solidity: event QuestionResolved(bytes32 indexed questionID, int256 indexed settledPrice, uint256[] payouts)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionResolved(opts *bind.FilterOpts, questionID [][32]byte, settledPrice []*big.Int) (*PolymarketUmaCtfAdapterV2QuestionResolvedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}
	var settledPriceRule []interface{}
	for _, settledPriceItem := range settledPrice {
		settledPriceRule = append(settledPriceRule, settledPriceItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionResolved", questionIDRule, settledPriceRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionResolvedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionResolved", logs: logs, sub: sub}, nil
}

// WatchQuestionResolved is a free log subscription operation binding the contract event 0x566c3fbdd12dd86bb341787f6d531f79fd7ad4ce7e3ae2d15ac0ca1b601af9df.
//
// Solidity: event QuestionResolved(bytes32 indexed questionID, int256 indexed settledPrice, uint256[] payouts)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionResolved(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionResolved, questionID [][32]byte, settledPrice []*big.Int) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}
	var settledPriceRule []interface{}
	for _, settledPriceItem := range settledPrice {
		settledPriceRule = append(settledPriceRule, settledPriceItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionResolved", questionIDRule, settledPriceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionResolved)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionResolved", log); err != nil {
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

// ParseQuestionResolved is a log parse operation binding the contract event 0x566c3fbdd12dd86bb341787f6d531f79fd7ad4ce7e3ae2d15ac0ca1b601af9df.
//
// Solidity: event QuestionResolved(bytes32 indexed questionID, int256 indexed settledPrice, uint256[] payouts)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionResolved(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionResolved, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionResolved)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionResolved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2QuestionUnpausedIterator is returned from FilterQuestionUnpaused and is used to iterate over the raw logs and unpacked data for QuestionUnpaused events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionUnpausedIterator struct {
	Event *PolymarketUmaCtfAdapterV2QuestionUnpaused // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2QuestionUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2QuestionUnpaused)
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
		it.Event = new(PolymarketUmaCtfAdapterV2QuestionUnpaused)
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
func (it *PolymarketUmaCtfAdapterV2QuestionUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2QuestionUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2QuestionUnpaused represents a QuestionUnpaused event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2QuestionUnpaused struct {
	QuestionID [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterQuestionUnpaused is a free log retrieval operation binding the contract event 0x92d28918c5574e7fc0f4f948c39502682c81cfb4089b07b83f95b3264e5e5e06.
//
// Solidity: event QuestionUnpaused(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterQuestionUnpaused(opts *bind.FilterOpts, questionID [][32]byte) (*PolymarketUmaCtfAdapterV2QuestionUnpausedIterator, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "QuestionUnpaused", questionIDRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2QuestionUnpausedIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "QuestionUnpaused", logs: logs, sub: sub}, nil
}

// WatchQuestionUnpaused is a free log subscription operation binding the contract event 0x92d28918c5574e7fc0f4f948c39502682c81cfb4089b07b83f95b3264e5e5e06.
//
// Solidity: event QuestionUnpaused(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchQuestionUnpaused(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2QuestionUnpaused, questionID [][32]byte) (event.Subscription, error) {

	var questionIDRule []interface{}
	for _, questionIDItem := range questionID {
		questionIDRule = append(questionIDRule, questionIDItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "QuestionUnpaused", questionIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2QuestionUnpaused)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionUnpaused", log); err != nil {
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

// ParseQuestionUnpaused is a log parse operation binding the contract event 0x92d28918c5574e7fc0f4f948c39502682c81cfb4089b07b83f95b3264e5e5e06.
//
// Solidity: event QuestionUnpaused(bytes32 indexed questionID)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseQuestionUnpaused(log types.Log) (*PolymarketUmaCtfAdapterV2QuestionUnpaused, error) {
	event := new(PolymarketUmaCtfAdapterV2QuestionUnpaused)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "QuestionUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolymarketUmaCtfAdapterV2RemovedAdminIterator is returned from FilterRemovedAdmin and is used to iterate over the raw logs and unpacked data for RemovedAdmin events raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2RemovedAdminIterator struct {
	Event *PolymarketUmaCtfAdapterV2RemovedAdmin // Event containing the contract specifics and raw log

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
func (it *PolymarketUmaCtfAdapterV2RemovedAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolymarketUmaCtfAdapterV2RemovedAdmin)
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
		it.Event = new(PolymarketUmaCtfAdapterV2RemovedAdmin)
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
func (it *PolymarketUmaCtfAdapterV2RemovedAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolymarketUmaCtfAdapterV2RemovedAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolymarketUmaCtfAdapterV2RemovedAdmin represents a RemovedAdmin event raised by the PolymarketUmaCtfAdapterV2 contract.
type PolymarketUmaCtfAdapterV2RemovedAdmin struct {
	Admin        common.Address
	RemovedAdmin common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemovedAdmin is a free log retrieval operation binding the contract event 0x787a2e12f4a55b658b8f573c32432ee11a5e8b51677d1e1e937aaf6a0bb5776e.
//
// Solidity: event RemovedAdmin(address indexed admin, address indexed removedAdmin)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) FilterRemovedAdmin(opts *bind.FilterOpts, admin []common.Address, removedAdmin []common.Address) (*PolymarketUmaCtfAdapterV2RemovedAdminIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var removedAdminRule []interface{}
	for _, removedAdminItem := range removedAdmin {
		removedAdminRule = append(removedAdminRule, removedAdminItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.FilterLogs(opts, "RemovedAdmin", adminRule, removedAdminRule)
	if err != nil {
		return nil, err
	}
	return &PolymarketUmaCtfAdapterV2RemovedAdminIterator{contract: _PolymarketUmaCtfAdapterV2.contract, event: "RemovedAdmin", logs: logs, sub: sub}, nil
}

// WatchRemovedAdmin is a free log subscription operation binding the contract event 0x787a2e12f4a55b658b8f573c32432ee11a5e8b51677d1e1e937aaf6a0bb5776e.
//
// Solidity: event RemovedAdmin(address indexed admin, address indexed removedAdmin)
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) WatchRemovedAdmin(opts *bind.WatchOpts, sink chan<- *PolymarketUmaCtfAdapterV2RemovedAdmin, admin []common.Address, removedAdmin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}
	var removedAdminRule []interface{}
	for _, removedAdminItem := range removedAdmin {
		removedAdminRule = append(removedAdminRule, removedAdminItem)
	}

	logs, sub, err := _PolymarketUmaCtfAdapterV2.contract.WatchLogs(opts, "RemovedAdmin", adminRule, removedAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolymarketUmaCtfAdapterV2RemovedAdmin)
				if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "RemovedAdmin", log); err != nil {
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
func (_PolymarketUmaCtfAdapterV2 *PolymarketUmaCtfAdapterV2Filterer) ParseRemovedAdmin(log types.Log) (*PolymarketUmaCtfAdapterV2RemovedAdmin, error) {
	event := new(PolymarketUmaCtfAdapterV2RemovedAdmin)
	if err := _PolymarketUmaCtfAdapterV2.contract.UnpackLog(event, "RemovedAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
