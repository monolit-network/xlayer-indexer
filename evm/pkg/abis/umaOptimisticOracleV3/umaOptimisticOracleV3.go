// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package umaOptimisticOracleV3

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

// OptimisticOracleV3InterfaceAssertion is an auto generated low-level Go binding around an user-defined struct.
type OptimisticOracleV3InterfaceAssertion struct {
	EscalationManagerSettings OptimisticOracleV3InterfaceEscalationManagerSettings
	Asserter                  common.Address
	AssertionTime             uint64
	Settled                   bool
	Currency                  common.Address
	ExpirationTime            uint64
	SettlementResolution      bool
	DomainId                  [32]byte
	Identifier                [32]byte
	Bond                      *big.Int
	CallbackRecipient         common.Address
	Disputer                  common.Address
}

// OptimisticOracleV3InterfaceEscalationManagerSettings is an auto generated low-level Go binding around an user-defined struct.
type OptimisticOracleV3InterfaceEscalationManagerSettings struct {
	ArbitrateViaEscalationManager bool
	DiscardOracle                 bool
	ValidateDisputers             bool
	AssertingCaller               common.Address
	EscalationManager             common.Address
}

// UmaOptimisticOracleV3MetaData contains all meta data concerning the UmaOptimisticOracleV3 contract.
var UmaOptimisticOracleV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractFinderInterface\",\"name\":\"_finder\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_defaultCurrency\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_defaultLiveness\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"defaultCurrency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"defaultLiveness\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnedBondPercentage\",\"type\":\"uint256\"}],\"name\":\"AdminPropertiesSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"}],\"name\":\"AssertionDisputed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"domainId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"claim\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asserter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"escalationManager\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"}],\"name\":\"AssertionMade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"bondRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"disputed\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"settlementResolution\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"settleCaller\",\"type\":\"address\"}],\"name\":\"AssertionSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"claim\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"asserter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"callbackRecipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"escalationManager\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"liveness\",\"type\":\"uint64\"},{\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"domainId\",\"type\":\"bytes32\"}],\"name\":\"assertTruth\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"claim\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"asserter\",\"type\":\"address\"}],\"name\":\"assertTruthWithDefaults\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"assertions\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"arbitrateViaEscalationManager\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"discardOracle\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"validateDisputers\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"assertingCaller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"escalationManager\",\"type\":\"address\"}],\"internalType\":\"structOptimisticOracleV3Interface.EscalationManagerSettings\",\"name\":\"escalationManagerSettings\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"asserter\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"assertionTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"settled\",\"type\":\"bool\"},{\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"settlementResolution\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"domainId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackRecipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnedBondPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"cachedCurrencies\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isWhitelisted\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"finalFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cachedIdentifiers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cachedOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultCurrency\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultIdentifier\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultLiveness\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"}],\"name\":\"disputeAssertion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finder\",\"outputs\":[{\"internalType\":\"contractFinderInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"}],\"name\":\"getAssertion\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"arbitrateViaEscalationManager\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"discardOracle\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"validateDisputers\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"assertingCaller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"escalationManager\",\"type\":\"address\"}],\"internalType\":\"structOptimisticOracleV3Interface.EscalationManagerSettings\",\"name\":\"escalationManagerSettings\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"asserter\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"assertionTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"settled\",\"type\":\"bool\"},{\"internalType\":\"contractIERC20\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"expirationTime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"settlementResolution\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"domainId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"callbackRecipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"}],\"internalType\":\"structOptimisticOracleV3Interface.Assertion\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"}],\"name\":\"getAssertionResult\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"currency\",\"type\":\"address\"}],\"name\":\"getMinimumBond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numericalTrue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_defaultCurrency\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_defaultLiveness\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"_burnedBondPercentage\",\"type\":\"uint256\"}],\"name\":\"setAdminProperties\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"}],\"name\":\"settleAndGetAssertionResult\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"}],\"name\":\"settleAssertion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assertionId\",\"type\":\"bytes32\"}],\"name\":\"stampAssertion\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"identifier\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"currency\",\"type\":\"address\"}],\"name\":\"syncUmaParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// UmaOptimisticOracleV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use UmaOptimisticOracleV3MetaData.ABI instead.
var UmaOptimisticOracleV3ABI = UmaOptimisticOracleV3MetaData.ABI

// UmaOptimisticOracleV3 is an auto generated Go binding around an Ethereum contract.
type UmaOptimisticOracleV3 struct {
	UmaOptimisticOracleV3Caller     // Read-only binding to the contract
	UmaOptimisticOracleV3Transactor // Write-only binding to the contract
	UmaOptimisticOracleV3Filterer   // Log filterer for contract events
}

// UmaOptimisticOracleV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UmaOptimisticOracleV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UmaOptimisticOracleV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UmaOptimisticOracleV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UmaOptimisticOracleV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UmaOptimisticOracleV3Session struct {
	Contract     *UmaOptimisticOracleV3 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UmaOptimisticOracleV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UmaOptimisticOracleV3CallerSession struct {
	Contract *UmaOptimisticOracleV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// UmaOptimisticOracleV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UmaOptimisticOracleV3TransactorSession struct {
	Contract     *UmaOptimisticOracleV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// UmaOptimisticOracleV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type UmaOptimisticOracleV3Raw struct {
	Contract *UmaOptimisticOracleV3 // Generic contract binding to access the raw methods on
}

// UmaOptimisticOracleV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV3CallerRaw struct {
	Contract *UmaOptimisticOracleV3Caller // Generic read-only contract binding to access the raw methods on
}

// UmaOptimisticOracleV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UmaOptimisticOracleV3TransactorRaw struct {
	Contract *UmaOptimisticOracleV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewUmaOptimisticOracleV3 creates a new instance of UmaOptimisticOracleV3, bound to a specific deployed contract.
func NewUmaOptimisticOracleV3(address common.Address, backend bind.ContractBackend) (*UmaOptimisticOracleV3, error) {
	contract, err := bindUmaOptimisticOracleV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3{UmaOptimisticOracleV3Caller: UmaOptimisticOracleV3Caller{contract: contract}, UmaOptimisticOracleV3Transactor: UmaOptimisticOracleV3Transactor{contract: contract}, UmaOptimisticOracleV3Filterer: UmaOptimisticOracleV3Filterer{contract: contract}}, nil
}

// NewUmaOptimisticOracleV3Caller creates a new read-only instance of UmaOptimisticOracleV3, bound to a specific deployed contract.
func NewUmaOptimisticOracleV3Caller(address common.Address, caller bind.ContractCaller) (*UmaOptimisticOracleV3Caller, error) {
	contract, err := bindUmaOptimisticOracleV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3Caller{contract: contract}, nil
}

// NewUmaOptimisticOracleV3Transactor creates a new write-only instance of UmaOptimisticOracleV3, bound to a specific deployed contract.
func NewUmaOptimisticOracleV3Transactor(address common.Address, transactor bind.ContractTransactor) (*UmaOptimisticOracleV3Transactor, error) {
	contract, err := bindUmaOptimisticOracleV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3Transactor{contract: contract}, nil
}

// NewUmaOptimisticOracleV3Filterer creates a new log filterer instance of UmaOptimisticOracleV3, bound to a specific deployed contract.
func NewUmaOptimisticOracleV3Filterer(address common.Address, filterer bind.ContractFilterer) (*UmaOptimisticOracleV3Filterer, error) {
	contract, err := bindUmaOptimisticOracleV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3Filterer{contract: contract}, nil
}

// bindUmaOptimisticOracleV3 binds a generic wrapper to an already deployed contract.
func bindUmaOptimisticOracleV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UmaOptimisticOracleV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UmaOptimisticOracleV3.Contract.UmaOptimisticOracleV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.UmaOptimisticOracleV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.UmaOptimisticOracleV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UmaOptimisticOracleV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.contract.Transact(opts, method, params...)
}

// Assertions is a free data retrieval call binding the contract method 0xd60715b5.
//
// Solidity: function assertions(bytes32 ) view returns((bool,bool,bool,address,address) escalationManagerSettings, address asserter, uint64 assertionTime, bool settled, address currency, uint64 expirationTime, bool settlementResolution, bytes32 domainId, bytes32 identifier, uint256 bond, address callbackRecipient, address disputer)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) Assertions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	EscalationManagerSettings OptimisticOracleV3InterfaceEscalationManagerSettings
	Asserter                  common.Address
	AssertionTime             uint64
	Settled                   bool
	Currency                  common.Address
	ExpirationTime            uint64
	SettlementResolution      bool
	DomainId                  [32]byte
	Identifier                [32]byte
	Bond                      *big.Int
	CallbackRecipient         common.Address
	Disputer                  common.Address
}, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "assertions", arg0)

	outstruct := new(struct {
		EscalationManagerSettings OptimisticOracleV3InterfaceEscalationManagerSettings
		Asserter                  common.Address
		AssertionTime             uint64
		Settled                   bool
		Currency                  common.Address
		ExpirationTime            uint64
		SettlementResolution      bool
		DomainId                  [32]byte
		Identifier                [32]byte
		Bond                      *big.Int
		CallbackRecipient         common.Address
		Disputer                  common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EscalationManagerSettings = *abi.ConvertType(out[0], new(OptimisticOracleV3InterfaceEscalationManagerSettings)).(*OptimisticOracleV3InterfaceEscalationManagerSettings)
	outstruct.Asserter = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.AssertionTime = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Settled = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.Currency = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.ExpirationTime = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.SettlementResolution = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.DomainId = *abi.ConvertType(out[7], new([32]byte)).(*[32]byte)
	outstruct.Identifier = *abi.ConvertType(out[8], new([32]byte)).(*[32]byte)
	outstruct.Bond = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	outstruct.CallbackRecipient = *abi.ConvertType(out[10], new(common.Address)).(*common.Address)
	outstruct.Disputer = *abi.ConvertType(out[11], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Assertions is a free data retrieval call binding the contract method 0xd60715b5.
//
// Solidity: function assertions(bytes32 ) view returns((bool,bool,bool,address,address) escalationManagerSettings, address asserter, uint64 assertionTime, bool settled, address currency, uint64 expirationTime, bool settlementResolution, bytes32 domainId, bytes32 identifier, uint256 bond, address callbackRecipient, address disputer)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) Assertions(arg0 [32]byte) (struct {
	EscalationManagerSettings OptimisticOracleV3InterfaceEscalationManagerSettings
	Asserter                  common.Address
	AssertionTime             uint64
	Settled                   bool
	Currency                  common.Address
	ExpirationTime            uint64
	SettlementResolution      bool
	DomainId                  [32]byte
	Identifier                [32]byte
	Bond                      *big.Int
	CallbackRecipient         common.Address
	Disputer                  common.Address
}, error) {
	return _UmaOptimisticOracleV3.Contract.Assertions(&_UmaOptimisticOracleV3.CallOpts, arg0)
}

// Assertions is a free data retrieval call binding the contract method 0xd60715b5.
//
// Solidity: function assertions(bytes32 ) view returns((bool,bool,bool,address,address) escalationManagerSettings, address asserter, uint64 assertionTime, bool settled, address currency, uint64 expirationTime, bool settlementResolution, bytes32 domainId, bytes32 identifier, uint256 bond, address callbackRecipient, address disputer)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) Assertions(arg0 [32]byte) (struct {
	EscalationManagerSettings OptimisticOracleV3InterfaceEscalationManagerSettings
	Asserter                  common.Address
	AssertionTime             uint64
	Settled                   bool
	Currency                  common.Address
	ExpirationTime            uint64
	SettlementResolution      bool
	DomainId                  [32]byte
	Identifier                [32]byte
	Bond                      *big.Int
	CallbackRecipient         common.Address
	Disputer                  common.Address
}, error) {
	return _UmaOptimisticOracleV3.Contract.Assertions(&_UmaOptimisticOracleV3.CallOpts, arg0)
}

// BurnedBondPercentage is a free data retrieval call binding the contract method 0x08e7c3e6.
//
// Solidity: function burnedBondPercentage() view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) BurnedBondPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "burnedBondPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BurnedBondPercentage is a free data retrieval call binding the contract method 0x08e7c3e6.
//
// Solidity: function burnedBondPercentage() view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) BurnedBondPercentage() (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.BurnedBondPercentage(&_UmaOptimisticOracleV3.CallOpts)
}

// BurnedBondPercentage is a free data retrieval call binding the contract method 0x08e7c3e6.
//
// Solidity: function burnedBondPercentage() view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) BurnedBondPercentage() (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.BurnedBondPercentage(&_UmaOptimisticOracleV3.CallOpts)
}

// CachedCurrencies is a free data retrieval call binding the contract method 0x70762157.
//
// Solidity: function cachedCurrencies(address ) view returns(bool isWhitelisted, uint256 finalFee)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) CachedCurrencies(opts *bind.CallOpts, arg0 common.Address) (struct {
	IsWhitelisted bool
	FinalFee      *big.Int
}, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "cachedCurrencies", arg0)

	outstruct := new(struct {
		IsWhitelisted bool
		FinalFee      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsWhitelisted = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.FinalFee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CachedCurrencies is a free data retrieval call binding the contract method 0x70762157.
//
// Solidity: function cachedCurrencies(address ) view returns(bool isWhitelisted, uint256 finalFee)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) CachedCurrencies(arg0 common.Address) (struct {
	IsWhitelisted bool
	FinalFee      *big.Int
}, error) {
	return _UmaOptimisticOracleV3.Contract.CachedCurrencies(&_UmaOptimisticOracleV3.CallOpts, arg0)
}

// CachedCurrencies is a free data retrieval call binding the contract method 0x70762157.
//
// Solidity: function cachedCurrencies(address ) view returns(bool isWhitelisted, uint256 finalFee)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) CachedCurrencies(arg0 common.Address) (struct {
	IsWhitelisted bool
	FinalFee      *big.Int
}, error) {
	return _UmaOptimisticOracleV3.Contract.CachedCurrencies(&_UmaOptimisticOracleV3.CallOpts, arg0)
}

// CachedIdentifiers is a free data retrieval call binding the contract method 0x530dd392.
//
// Solidity: function cachedIdentifiers(bytes32 ) view returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) CachedIdentifiers(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "cachedIdentifiers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CachedIdentifiers is a free data retrieval call binding the contract method 0x530dd392.
//
// Solidity: function cachedIdentifiers(bytes32 ) view returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) CachedIdentifiers(arg0 [32]byte) (bool, error) {
	return _UmaOptimisticOracleV3.Contract.CachedIdentifiers(&_UmaOptimisticOracleV3.CallOpts, arg0)
}

// CachedIdentifiers is a free data retrieval call binding the contract method 0x530dd392.
//
// Solidity: function cachedIdentifiers(bytes32 ) view returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) CachedIdentifiers(arg0 [32]byte) (bool, error) {
	return _UmaOptimisticOracleV3.Contract.CachedIdentifiers(&_UmaOptimisticOracleV3.CallOpts, arg0)
}

// CachedOracle is a free data retrieval call binding the contract method 0xa7af2d0f.
//
// Solidity: function cachedOracle() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) CachedOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "cachedOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CachedOracle is a free data retrieval call binding the contract method 0xa7af2d0f.
//
// Solidity: function cachedOracle() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) CachedOracle() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.CachedOracle(&_UmaOptimisticOracleV3.CallOpts)
}

// CachedOracle is a free data retrieval call binding the contract method 0xa7af2d0f.
//
// Solidity: function cachedOracle() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) CachedOracle() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.CachedOracle(&_UmaOptimisticOracleV3.CallOpts)
}

// DefaultCurrency is a free data retrieval call binding the contract method 0x20402e1d.
//
// Solidity: function defaultCurrency() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) DefaultCurrency(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "defaultCurrency")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultCurrency is a free data retrieval call binding the contract method 0x20402e1d.
//
// Solidity: function defaultCurrency() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) DefaultCurrency() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.DefaultCurrency(&_UmaOptimisticOracleV3.CallOpts)
}

// DefaultCurrency is a free data retrieval call binding the contract method 0x20402e1d.
//
// Solidity: function defaultCurrency() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) DefaultCurrency() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.DefaultCurrency(&_UmaOptimisticOracleV3.CallOpts)
}

// DefaultIdentifier is a free data retrieval call binding the contract method 0xd509b017.
//
// Solidity: function defaultIdentifier() view returns(bytes32)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) DefaultIdentifier(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "defaultIdentifier")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DefaultIdentifier is a free data retrieval call binding the contract method 0xd509b017.
//
// Solidity: function defaultIdentifier() view returns(bytes32)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) DefaultIdentifier() ([32]byte, error) {
	return _UmaOptimisticOracleV3.Contract.DefaultIdentifier(&_UmaOptimisticOracleV3.CallOpts)
}

// DefaultIdentifier is a free data retrieval call binding the contract method 0xd509b017.
//
// Solidity: function defaultIdentifier() view returns(bytes32)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) DefaultIdentifier() ([32]byte, error) {
	return _UmaOptimisticOracleV3.Contract.DefaultIdentifier(&_UmaOptimisticOracleV3.CallOpts)
}

// DefaultLiveness is a free data retrieval call binding the contract method 0xfe4e1983.
//
// Solidity: function defaultLiveness() view returns(uint64)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) DefaultLiveness(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "defaultLiveness")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DefaultLiveness is a free data retrieval call binding the contract method 0xfe4e1983.
//
// Solidity: function defaultLiveness() view returns(uint64)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) DefaultLiveness() (uint64, error) {
	return _UmaOptimisticOracleV3.Contract.DefaultLiveness(&_UmaOptimisticOracleV3.CallOpts)
}

// DefaultLiveness is a free data retrieval call binding the contract method 0xfe4e1983.
//
// Solidity: function defaultLiveness() view returns(uint64)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) DefaultLiveness() (uint64, error) {
	return _UmaOptimisticOracleV3.Contract.DefaultLiveness(&_UmaOptimisticOracleV3.CallOpts)
}

// Finder is a free data retrieval call binding the contract method 0xb9a3c84c.
//
// Solidity: function finder() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) Finder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "finder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Finder is a free data retrieval call binding the contract method 0xb9a3c84c.
//
// Solidity: function finder() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) Finder() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.Finder(&_UmaOptimisticOracleV3.CallOpts)
}

// Finder is a free data retrieval call binding the contract method 0xb9a3c84c.
//
// Solidity: function finder() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) Finder() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.Finder(&_UmaOptimisticOracleV3.CallOpts)
}

// GetAssertion is a free data retrieval call binding the contract method 0x88302884.
//
// Solidity: function getAssertion(bytes32 assertionId) view returns(((bool,bool,bool,address,address),address,uint64,bool,address,uint64,bool,bytes32,bytes32,uint256,address,address))
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) GetAssertion(opts *bind.CallOpts, assertionId [32]byte) (OptimisticOracleV3InterfaceAssertion, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "getAssertion", assertionId)

	if err != nil {
		return *new(OptimisticOracleV3InterfaceAssertion), err
	}

	out0 := *abi.ConvertType(out[0], new(OptimisticOracleV3InterfaceAssertion)).(*OptimisticOracleV3InterfaceAssertion)

	return out0, err

}

// GetAssertion is a free data retrieval call binding the contract method 0x88302884.
//
// Solidity: function getAssertion(bytes32 assertionId) view returns(((bool,bool,bool,address,address),address,uint64,bool,address,uint64,bool,bytes32,bytes32,uint256,address,address))
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) GetAssertion(assertionId [32]byte) (OptimisticOracleV3InterfaceAssertion, error) {
	return _UmaOptimisticOracleV3.Contract.GetAssertion(&_UmaOptimisticOracleV3.CallOpts, assertionId)
}

// GetAssertion is a free data retrieval call binding the contract method 0x88302884.
//
// Solidity: function getAssertion(bytes32 assertionId) view returns(((bool,bool,bool,address,address),address,uint64,bool,address,uint64,bool,bytes32,bytes32,uint256,address,address))
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) GetAssertion(assertionId [32]byte) (OptimisticOracleV3InterfaceAssertion, error) {
	return _UmaOptimisticOracleV3.Contract.GetAssertion(&_UmaOptimisticOracleV3.CallOpts, assertionId)
}

// GetAssertionResult is a free data retrieval call binding the contract method 0xe39dfd7f.
//
// Solidity: function getAssertionResult(bytes32 assertionId) view returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) GetAssertionResult(opts *bind.CallOpts, assertionId [32]byte) (bool, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "getAssertionResult", assertionId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAssertionResult is a free data retrieval call binding the contract method 0xe39dfd7f.
//
// Solidity: function getAssertionResult(bytes32 assertionId) view returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) GetAssertionResult(assertionId [32]byte) (bool, error) {
	return _UmaOptimisticOracleV3.Contract.GetAssertionResult(&_UmaOptimisticOracleV3.CallOpts, assertionId)
}

// GetAssertionResult is a free data retrieval call binding the contract method 0xe39dfd7f.
//
// Solidity: function getAssertionResult(bytes32 assertionId) view returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) GetAssertionResult(assertionId [32]byte) (bool, error) {
	return _UmaOptimisticOracleV3.Contract.GetAssertionResult(&_UmaOptimisticOracleV3.CallOpts, assertionId)
}

// GetCurrentTime is a free data retrieval call binding the contract method 0x29cb924d.
//
// Solidity: function getCurrentTime() view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) GetCurrentTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "getCurrentTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentTime is a free data retrieval call binding the contract method 0x29cb924d.
//
// Solidity: function getCurrentTime() view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) GetCurrentTime() (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.GetCurrentTime(&_UmaOptimisticOracleV3.CallOpts)
}

// GetCurrentTime is a free data retrieval call binding the contract method 0x29cb924d.
//
// Solidity: function getCurrentTime() view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) GetCurrentTime() (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.GetCurrentTime(&_UmaOptimisticOracleV3.CallOpts)
}

// GetMinimumBond is a free data retrieval call binding the contract method 0x4360af3d.
//
// Solidity: function getMinimumBond(address currency) view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) GetMinimumBond(opts *bind.CallOpts, currency common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "getMinimumBond", currency)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumBond is a free data retrieval call binding the contract method 0x4360af3d.
//
// Solidity: function getMinimumBond(address currency) view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) GetMinimumBond(currency common.Address) (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.GetMinimumBond(&_UmaOptimisticOracleV3.CallOpts, currency)
}

// GetMinimumBond is a free data retrieval call binding the contract method 0x4360af3d.
//
// Solidity: function getMinimumBond(address currency) view returns(uint256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) GetMinimumBond(currency common.Address) (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.GetMinimumBond(&_UmaOptimisticOracleV3.CallOpts, currency)
}

// NumericalTrue is a free data retrieval call binding the contract method 0xda03b36e.
//
// Solidity: function numericalTrue() view returns(int256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) NumericalTrue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "numericalTrue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumericalTrue is a free data retrieval call binding the contract method 0xda03b36e.
//
// Solidity: function numericalTrue() view returns(int256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) NumericalTrue() (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.NumericalTrue(&_UmaOptimisticOracleV3.CallOpts)
}

// NumericalTrue is a free data retrieval call binding the contract method 0xda03b36e.
//
// Solidity: function numericalTrue() view returns(int256)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) NumericalTrue() (*big.Int, error) {
	return _UmaOptimisticOracleV3.Contract.NumericalTrue(&_UmaOptimisticOracleV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) Owner() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.Owner(&_UmaOptimisticOracleV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) Owner() (common.Address, error) {
	return _UmaOptimisticOracleV3.Contract.Owner(&_UmaOptimisticOracleV3.CallOpts)
}

// StampAssertion is a free data retrieval call binding the contract method 0xafedba4f.
//
// Solidity: function stampAssertion(bytes32 assertionId) view returns(bytes)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Caller) StampAssertion(opts *bind.CallOpts, assertionId [32]byte) ([]byte, error) {
	var out []interface{}
	err := _UmaOptimisticOracleV3.contract.Call(opts, &out, "stampAssertion", assertionId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// StampAssertion is a free data retrieval call binding the contract method 0xafedba4f.
//
// Solidity: function stampAssertion(bytes32 assertionId) view returns(bytes)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) StampAssertion(assertionId [32]byte) ([]byte, error) {
	return _UmaOptimisticOracleV3.Contract.StampAssertion(&_UmaOptimisticOracleV3.CallOpts, assertionId)
}

// StampAssertion is a free data retrieval call binding the contract method 0xafedba4f.
//
// Solidity: function stampAssertion(bytes32 assertionId) view returns(bytes)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3CallerSession) StampAssertion(assertionId [32]byte) ([]byte, error) {
	return _UmaOptimisticOracleV3.Contract.StampAssertion(&_UmaOptimisticOracleV3.CallOpts, assertionId)
}

// AssertTruth is a paid mutator transaction binding the contract method 0x6457c979.
//
// Solidity: function assertTruth(bytes claim, address asserter, address callbackRecipient, address escalationManager, uint64 liveness, address currency, uint256 bond, bytes32 identifier, bytes32 domainId) returns(bytes32 assertionId)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) AssertTruth(opts *bind.TransactOpts, claim []byte, asserter common.Address, callbackRecipient common.Address, escalationManager common.Address, liveness uint64, currency common.Address, bond *big.Int, identifier [32]byte, domainId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "assertTruth", claim, asserter, callbackRecipient, escalationManager, liveness, currency, bond, identifier, domainId)
}

// AssertTruth is a paid mutator transaction binding the contract method 0x6457c979.
//
// Solidity: function assertTruth(bytes claim, address asserter, address callbackRecipient, address escalationManager, uint64 liveness, address currency, uint256 bond, bytes32 identifier, bytes32 domainId) returns(bytes32 assertionId)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) AssertTruth(claim []byte, asserter common.Address, callbackRecipient common.Address, escalationManager common.Address, liveness uint64, currency common.Address, bond *big.Int, identifier [32]byte, domainId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.AssertTruth(&_UmaOptimisticOracleV3.TransactOpts, claim, asserter, callbackRecipient, escalationManager, liveness, currency, bond, identifier, domainId)
}

// AssertTruth is a paid mutator transaction binding the contract method 0x6457c979.
//
// Solidity: function assertTruth(bytes claim, address asserter, address callbackRecipient, address escalationManager, uint64 liveness, address currency, uint256 bond, bytes32 identifier, bytes32 domainId) returns(bytes32 assertionId)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) AssertTruth(claim []byte, asserter common.Address, callbackRecipient common.Address, escalationManager common.Address, liveness uint64, currency common.Address, bond *big.Int, identifier [32]byte, domainId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.AssertTruth(&_UmaOptimisticOracleV3.TransactOpts, claim, asserter, callbackRecipient, escalationManager, liveness, currency, bond, identifier, domainId)
}

// AssertTruthWithDefaults is a paid mutator transaction binding the contract method 0x36b13af4.
//
// Solidity: function assertTruthWithDefaults(bytes claim, address asserter) returns(bytes32)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) AssertTruthWithDefaults(opts *bind.TransactOpts, claim []byte, asserter common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "assertTruthWithDefaults", claim, asserter)
}

// AssertTruthWithDefaults is a paid mutator transaction binding the contract method 0x36b13af4.
//
// Solidity: function assertTruthWithDefaults(bytes claim, address asserter) returns(bytes32)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) AssertTruthWithDefaults(claim []byte, asserter common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.AssertTruthWithDefaults(&_UmaOptimisticOracleV3.TransactOpts, claim, asserter)
}

// AssertTruthWithDefaults is a paid mutator transaction binding the contract method 0x36b13af4.
//
// Solidity: function assertTruthWithDefaults(bytes claim, address asserter) returns(bytes32)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) AssertTruthWithDefaults(claim []byte, asserter common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.AssertTruthWithDefaults(&_UmaOptimisticOracleV3.TransactOpts, claim, asserter)
}

// DisputeAssertion is a paid mutator transaction binding the contract method 0xa6a22b43.
//
// Solidity: function disputeAssertion(bytes32 assertionId, address disputer) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) DisputeAssertion(opts *bind.TransactOpts, assertionId [32]byte, disputer common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "disputeAssertion", assertionId, disputer)
}

// DisputeAssertion is a paid mutator transaction binding the contract method 0xa6a22b43.
//
// Solidity: function disputeAssertion(bytes32 assertionId, address disputer) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) DisputeAssertion(assertionId [32]byte, disputer common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.DisputeAssertion(&_UmaOptimisticOracleV3.TransactOpts, assertionId, disputer)
}

// DisputeAssertion is a paid mutator transaction binding the contract method 0xa6a22b43.
//
// Solidity: function disputeAssertion(bytes32 assertionId, address disputer) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) DisputeAssertion(assertionId [32]byte, disputer common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.DisputeAssertion(&_UmaOptimisticOracleV3.TransactOpts, assertionId, disputer)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) Multicall(data [][]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.Multicall(&_UmaOptimisticOracleV3.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.Multicall(&_UmaOptimisticOracleV3.TransactOpts, data)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) RenounceOwnership() (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.RenounceOwnership(&_UmaOptimisticOracleV3.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.RenounceOwnership(&_UmaOptimisticOracleV3.TransactOpts)
}

// SetAdminProperties is a paid mutator transaction binding the contract method 0x82762520.
//
// Solidity: function setAdminProperties(address _defaultCurrency, uint64 _defaultLiveness, uint256 _burnedBondPercentage) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) SetAdminProperties(opts *bind.TransactOpts, _defaultCurrency common.Address, _defaultLiveness uint64, _burnedBondPercentage *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "setAdminProperties", _defaultCurrency, _defaultLiveness, _burnedBondPercentage)
}

// SetAdminProperties is a paid mutator transaction binding the contract method 0x82762520.
//
// Solidity: function setAdminProperties(address _defaultCurrency, uint64 _defaultLiveness, uint256 _burnedBondPercentage) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) SetAdminProperties(_defaultCurrency common.Address, _defaultLiveness uint64, _burnedBondPercentage *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SetAdminProperties(&_UmaOptimisticOracleV3.TransactOpts, _defaultCurrency, _defaultLiveness, _burnedBondPercentage)
}

// SetAdminProperties is a paid mutator transaction binding the contract method 0x82762520.
//
// Solidity: function setAdminProperties(address _defaultCurrency, uint64 _defaultLiveness, uint256 _burnedBondPercentage) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) SetAdminProperties(_defaultCurrency common.Address, _defaultLiveness uint64, _burnedBondPercentage *big.Int) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SetAdminProperties(&_UmaOptimisticOracleV3.TransactOpts, _defaultCurrency, _defaultLiveness, _burnedBondPercentage)
}

// SettleAndGetAssertionResult is a paid mutator transaction binding the contract method 0x8ea2f2ab.
//
// Solidity: function settleAndGetAssertionResult(bytes32 assertionId) returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) SettleAndGetAssertionResult(opts *bind.TransactOpts, assertionId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "settleAndGetAssertionResult", assertionId)
}

// SettleAndGetAssertionResult is a paid mutator transaction binding the contract method 0x8ea2f2ab.
//
// Solidity: function settleAndGetAssertionResult(bytes32 assertionId) returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) SettleAndGetAssertionResult(assertionId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SettleAndGetAssertionResult(&_UmaOptimisticOracleV3.TransactOpts, assertionId)
}

// SettleAndGetAssertionResult is a paid mutator transaction binding the contract method 0x8ea2f2ab.
//
// Solidity: function settleAndGetAssertionResult(bytes32 assertionId) returns(bool)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) SettleAndGetAssertionResult(assertionId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SettleAndGetAssertionResult(&_UmaOptimisticOracleV3.TransactOpts, assertionId)
}

// SettleAssertion is a paid mutator transaction binding the contract method 0x4124beef.
//
// Solidity: function settleAssertion(bytes32 assertionId) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) SettleAssertion(opts *bind.TransactOpts, assertionId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "settleAssertion", assertionId)
}

// SettleAssertion is a paid mutator transaction binding the contract method 0x4124beef.
//
// Solidity: function settleAssertion(bytes32 assertionId) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) SettleAssertion(assertionId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SettleAssertion(&_UmaOptimisticOracleV3.TransactOpts, assertionId)
}

// SettleAssertion is a paid mutator transaction binding the contract method 0x4124beef.
//
// Solidity: function settleAssertion(bytes32 assertionId) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) SettleAssertion(assertionId [32]byte) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SettleAssertion(&_UmaOptimisticOracleV3.TransactOpts, assertionId)
}

// SyncUmaParams is a paid mutator transaction binding the contract method 0xa8655785.
//
// Solidity: function syncUmaParams(bytes32 identifier, address currency) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) SyncUmaParams(opts *bind.TransactOpts, identifier [32]byte, currency common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "syncUmaParams", identifier, currency)
}

// SyncUmaParams is a paid mutator transaction binding the contract method 0xa8655785.
//
// Solidity: function syncUmaParams(bytes32 identifier, address currency) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) SyncUmaParams(identifier [32]byte, currency common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SyncUmaParams(&_UmaOptimisticOracleV3.TransactOpts, identifier, currency)
}

// SyncUmaParams is a paid mutator transaction binding the contract method 0xa8655785.
//
// Solidity: function syncUmaParams(bytes32 identifier, address currency) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) SyncUmaParams(identifier [32]byte, currency common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.SyncUmaParams(&_UmaOptimisticOracleV3.TransactOpts, identifier, currency)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.TransferOwnership(&_UmaOptimisticOracleV3.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UmaOptimisticOracleV3.Contract.TransferOwnership(&_UmaOptimisticOracleV3.TransactOpts, newOwner)
}

// UmaOptimisticOracleV3AdminPropertiesSetIterator is returned from FilterAdminPropertiesSet and is used to iterate over the raw logs and unpacked data for AdminPropertiesSet events raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AdminPropertiesSetIterator struct {
	Event *UmaOptimisticOracleV3AdminPropertiesSet // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV3AdminPropertiesSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV3AdminPropertiesSet)
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
		it.Event = new(UmaOptimisticOracleV3AdminPropertiesSet)
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
func (it *UmaOptimisticOracleV3AdminPropertiesSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV3AdminPropertiesSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV3AdminPropertiesSet represents a AdminPropertiesSet event raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AdminPropertiesSet struct {
	DefaultCurrency      common.Address
	DefaultLiveness      uint64
	BurnedBondPercentage *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterAdminPropertiesSet is a free log retrieval operation binding the contract event 0xd0f09246d369018534c67fec6a6c3259c6f962ef82c5521c337ae0ccc104e4bd.
//
// Solidity: event AdminPropertiesSet(address defaultCurrency, uint64 defaultLiveness, uint256 burnedBondPercentage)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) FilterAdminPropertiesSet(opts *bind.FilterOpts) (*UmaOptimisticOracleV3AdminPropertiesSetIterator, error) {

	logs, sub, err := _UmaOptimisticOracleV3.contract.FilterLogs(opts, "AdminPropertiesSet")
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3AdminPropertiesSetIterator{contract: _UmaOptimisticOracleV3.contract, event: "AdminPropertiesSet", logs: logs, sub: sub}, nil
}

// WatchAdminPropertiesSet is a free log subscription operation binding the contract event 0xd0f09246d369018534c67fec6a6c3259c6f962ef82c5521c337ae0ccc104e4bd.
//
// Solidity: event AdminPropertiesSet(address defaultCurrency, uint64 defaultLiveness, uint256 burnedBondPercentage)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) WatchAdminPropertiesSet(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV3AdminPropertiesSet) (event.Subscription, error) {

	logs, sub, err := _UmaOptimisticOracleV3.contract.WatchLogs(opts, "AdminPropertiesSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV3AdminPropertiesSet)
				if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AdminPropertiesSet", log); err != nil {
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

// ParseAdminPropertiesSet is a log parse operation binding the contract event 0xd0f09246d369018534c67fec6a6c3259c6f962ef82c5521c337ae0ccc104e4bd.
//
// Solidity: event AdminPropertiesSet(address defaultCurrency, uint64 defaultLiveness, uint256 burnedBondPercentage)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) ParseAdminPropertiesSet(log types.Log) (*UmaOptimisticOracleV3AdminPropertiesSet, error) {
	event := new(UmaOptimisticOracleV3AdminPropertiesSet)
	if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AdminPropertiesSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV3AssertionDisputedIterator is returned from FilterAssertionDisputed and is used to iterate over the raw logs and unpacked data for AssertionDisputed events raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AssertionDisputedIterator struct {
	Event *UmaOptimisticOracleV3AssertionDisputed // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV3AssertionDisputedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV3AssertionDisputed)
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
		it.Event = new(UmaOptimisticOracleV3AssertionDisputed)
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
func (it *UmaOptimisticOracleV3AssertionDisputedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV3AssertionDisputedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV3AssertionDisputed represents a AssertionDisputed event raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AssertionDisputed struct {
	AssertionId [32]byte
	Caller      common.Address
	Disputer    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssertionDisputed is a free log retrieval operation binding the contract event 0x60133788b013c89f2a3756dbc47e3484997b87bd7e0af98a7d70232032c1ce2b.
//
// Solidity: event AssertionDisputed(bytes32 indexed assertionId, address indexed caller, address indexed disputer)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) FilterAssertionDisputed(opts *bind.FilterOpts, assertionId [][32]byte, caller []common.Address, disputer []common.Address) (*UmaOptimisticOracleV3AssertionDisputedIterator, error) {

	var assertionIdRule []interface{}
	for _, assertionIdItem := range assertionId {
		assertionIdRule = append(assertionIdRule, assertionIdItem)
	}
	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.FilterLogs(opts, "AssertionDisputed", assertionIdRule, callerRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3AssertionDisputedIterator{contract: _UmaOptimisticOracleV3.contract, event: "AssertionDisputed", logs: logs, sub: sub}, nil
}

// WatchAssertionDisputed is a free log subscription operation binding the contract event 0x60133788b013c89f2a3756dbc47e3484997b87bd7e0af98a7d70232032c1ce2b.
//
// Solidity: event AssertionDisputed(bytes32 indexed assertionId, address indexed caller, address indexed disputer)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) WatchAssertionDisputed(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV3AssertionDisputed, assertionId [][32]byte, caller []common.Address, disputer []common.Address) (event.Subscription, error) {

	var assertionIdRule []interface{}
	for _, assertionIdItem := range assertionId {
		assertionIdRule = append(assertionIdRule, assertionIdItem)
	}
	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.WatchLogs(opts, "AssertionDisputed", assertionIdRule, callerRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV3AssertionDisputed)
				if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AssertionDisputed", log); err != nil {
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

// ParseAssertionDisputed is a log parse operation binding the contract event 0x60133788b013c89f2a3756dbc47e3484997b87bd7e0af98a7d70232032c1ce2b.
//
// Solidity: event AssertionDisputed(bytes32 indexed assertionId, address indexed caller, address indexed disputer)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) ParseAssertionDisputed(log types.Log) (*UmaOptimisticOracleV3AssertionDisputed, error) {
	event := new(UmaOptimisticOracleV3AssertionDisputed)
	if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AssertionDisputed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV3AssertionMadeIterator is returned from FilterAssertionMade and is used to iterate over the raw logs and unpacked data for AssertionMade events raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AssertionMadeIterator struct {
	Event *UmaOptimisticOracleV3AssertionMade // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV3AssertionMadeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV3AssertionMade)
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
		it.Event = new(UmaOptimisticOracleV3AssertionMade)
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
func (it *UmaOptimisticOracleV3AssertionMadeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV3AssertionMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV3AssertionMade represents a AssertionMade event raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AssertionMade struct {
	AssertionId       [32]byte
	DomainId          [32]byte
	Claim             []byte
	Asserter          common.Address
	CallbackRecipient common.Address
	EscalationManager common.Address
	Caller            common.Address
	ExpirationTime    uint64
	Currency          common.Address
	Bond              *big.Int
	Identifier        [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAssertionMade is a free log retrieval operation binding the contract event 0xdb1513f0abeb57a364db56aa3eb52015cca5268f00fd67bc73aaf22bccab02b7.
//
// Solidity: event AssertionMade(bytes32 indexed assertionId, bytes32 domainId, bytes claim, address indexed asserter, address callbackRecipient, address escalationManager, address caller, uint64 expirationTime, address currency, uint256 bond, bytes32 indexed identifier)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) FilterAssertionMade(opts *bind.FilterOpts, assertionId [][32]byte, asserter []common.Address, identifier [][32]byte) (*UmaOptimisticOracleV3AssertionMadeIterator, error) {

	var assertionIdRule []interface{}
	for _, assertionIdItem := range assertionId {
		assertionIdRule = append(assertionIdRule, assertionIdItem)
	}

	var asserterRule []interface{}
	for _, asserterItem := range asserter {
		asserterRule = append(asserterRule, asserterItem)
	}

	var identifierRule []interface{}
	for _, identifierItem := range identifier {
		identifierRule = append(identifierRule, identifierItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.FilterLogs(opts, "AssertionMade", assertionIdRule, asserterRule, identifierRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3AssertionMadeIterator{contract: _UmaOptimisticOracleV3.contract, event: "AssertionMade", logs: logs, sub: sub}, nil
}

// WatchAssertionMade is a free log subscription operation binding the contract event 0xdb1513f0abeb57a364db56aa3eb52015cca5268f00fd67bc73aaf22bccab02b7.
//
// Solidity: event AssertionMade(bytes32 indexed assertionId, bytes32 domainId, bytes claim, address indexed asserter, address callbackRecipient, address escalationManager, address caller, uint64 expirationTime, address currency, uint256 bond, bytes32 indexed identifier)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) WatchAssertionMade(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV3AssertionMade, assertionId [][32]byte, asserter []common.Address, identifier [][32]byte) (event.Subscription, error) {

	var assertionIdRule []interface{}
	for _, assertionIdItem := range assertionId {
		assertionIdRule = append(assertionIdRule, assertionIdItem)
	}

	var asserterRule []interface{}
	for _, asserterItem := range asserter {
		asserterRule = append(asserterRule, asserterItem)
	}

	var identifierRule []interface{}
	for _, identifierItem := range identifier {
		identifierRule = append(identifierRule, identifierItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.WatchLogs(opts, "AssertionMade", assertionIdRule, asserterRule, identifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV3AssertionMade)
				if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AssertionMade", log); err != nil {
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

// ParseAssertionMade is a log parse operation binding the contract event 0xdb1513f0abeb57a364db56aa3eb52015cca5268f00fd67bc73aaf22bccab02b7.
//
// Solidity: event AssertionMade(bytes32 indexed assertionId, bytes32 domainId, bytes claim, address indexed asserter, address callbackRecipient, address escalationManager, address caller, uint64 expirationTime, address currency, uint256 bond, bytes32 indexed identifier)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) ParseAssertionMade(log types.Log) (*UmaOptimisticOracleV3AssertionMade, error) {
	event := new(UmaOptimisticOracleV3AssertionMade)
	if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AssertionMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV3AssertionSettledIterator is returned from FilterAssertionSettled and is used to iterate over the raw logs and unpacked data for AssertionSettled events raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AssertionSettledIterator struct {
	Event *UmaOptimisticOracleV3AssertionSettled // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV3AssertionSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV3AssertionSettled)
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
		it.Event = new(UmaOptimisticOracleV3AssertionSettled)
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
func (it *UmaOptimisticOracleV3AssertionSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV3AssertionSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV3AssertionSettled represents a AssertionSettled event raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3AssertionSettled struct {
	AssertionId          [32]byte
	BondRecipient        common.Address
	Disputed             bool
	SettlementResolution bool
	SettleCaller         common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterAssertionSettled is a free log retrieval operation binding the contract event 0xf4fa324b13daeb4e1aae736c2553632ae0fb16fb31f2d4da8ac99fd056313a13.
//
// Solidity: event AssertionSettled(bytes32 indexed assertionId, address indexed bondRecipient, bool disputed, bool settlementResolution, address settleCaller)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) FilterAssertionSettled(opts *bind.FilterOpts, assertionId [][32]byte, bondRecipient []common.Address) (*UmaOptimisticOracleV3AssertionSettledIterator, error) {

	var assertionIdRule []interface{}
	for _, assertionIdItem := range assertionId {
		assertionIdRule = append(assertionIdRule, assertionIdItem)
	}
	var bondRecipientRule []interface{}
	for _, bondRecipientItem := range bondRecipient {
		bondRecipientRule = append(bondRecipientRule, bondRecipientItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.FilterLogs(opts, "AssertionSettled", assertionIdRule, bondRecipientRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3AssertionSettledIterator{contract: _UmaOptimisticOracleV3.contract, event: "AssertionSettled", logs: logs, sub: sub}, nil
}

// WatchAssertionSettled is a free log subscription operation binding the contract event 0xf4fa324b13daeb4e1aae736c2553632ae0fb16fb31f2d4da8ac99fd056313a13.
//
// Solidity: event AssertionSettled(bytes32 indexed assertionId, address indexed bondRecipient, bool disputed, bool settlementResolution, address settleCaller)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) WatchAssertionSettled(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV3AssertionSettled, assertionId [][32]byte, bondRecipient []common.Address) (event.Subscription, error) {

	var assertionIdRule []interface{}
	for _, assertionIdItem := range assertionId {
		assertionIdRule = append(assertionIdRule, assertionIdItem)
	}
	var bondRecipientRule []interface{}
	for _, bondRecipientItem := range bondRecipient {
		bondRecipientRule = append(bondRecipientRule, bondRecipientItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.WatchLogs(opts, "AssertionSettled", assertionIdRule, bondRecipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV3AssertionSettled)
				if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AssertionSettled", log); err != nil {
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

// ParseAssertionSettled is a log parse operation binding the contract event 0xf4fa324b13daeb4e1aae736c2553632ae0fb16fb31f2d4da8ac99fd056313a13.
//
// Solidity: event AssertionSettled(bytes32 indexed assertionId, address indexed bondRecipient, bool disputed, bool settlementResolution, address settleCaller)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) ParseAssertionSettled(log types.Log) (*UmaOptimisticOracleV3AssertionSettled, error) {
	event := new(UmaOptimisticOracleV3AssertionSettled)
	if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "AssertionSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UmaOptimisticOracleV3OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3OwnershipTransferredIterator struct {
	Event *UmaOptimisticOracleV3OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *UmaOptimisticOracleV3OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UmaOptimisticOracleV3OwnershipTransferred)
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
		it.Event = new(UmaOptimisticOracleV3OwnershipTransferred)
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
func (it *UmaOptimisticOracleV3OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UmaOptimisticOracleV3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UmaOptimisticOracleV3OwnershipTransferred represents a OwnershipTransferred event raised by the UmaOptimisticOracleV3 contract.
type UmaOptimisticOracleV3OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UmaOptimisticOracleV3OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UmaOptimisticOracleV3OwnershipTransferredIterator{contract: _UmaOptimisticOracleV3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UmaOptimisticOracleV3OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UmaOptimisticOracleV3.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UmaOptimisticOracleV3OwnershipTransferred)
				if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UmaOptimisticOracleV3 *UmaOptimisticOracleV3Filterer) ParseOwnershipTransferred(log types.Log) (*UmaOptimisticOracleV3OwnershipTransferred, error) {
	event := new(UmaOptimisticOracleV3OwnershipTransferred)
	if err := _UmaOptimisticOracleV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
