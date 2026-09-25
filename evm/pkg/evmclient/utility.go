package evmclient

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/erc20"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (c *Client) GetDecimals(address common.Address) (uint8, error) {
	var decimals uint8
	var err error
	var fromCache bool
	ctx := context.Background()
	err = c.withRetry(ctx, "GetDecimals", func() error {
		var e error
		decimals, fromCache, e = c.getDecimals(address)
		if e != nil && !fromCache {
			c.logger.Warn("error getting decimals", zap.String("token", address.Hex()), zap.Error(e))
		}
		return e
	})
	return decimals, err
}

func (c *Client) getDecimals(address common.Address) (resDecimals uint8, fromCache bool, err error) {
	if address == models.ETHTokenAddress {
		return models.ETHTokensDecimals, true, nil
	}
	if decimals, ok := c.decimalsCache.Get(address); ok {
		return decimals, true, nil
	}
	if _, ok := c.decimalsErrorCache.Get(address); ok {
		return 0, false, nil
	}
	instance, err := erc20.NewErc20(address, c.ethClient)
	if err != nil {
		return 0, false, err
	}
	if err := c.waitRequests(context.Background(), 1); err != nil {
		return 0, false, err
	}

	decimals, err := instance.Decimals(nil)
	if err != nil {
		noSuchMethodError := isNoSuchMethodError(err)
		if noSuchMethodError {
			c.decimalsErrorCache.Add(address, err)
		}
		return 0, noSuchMethodError, err
	}
	c.decimalsCache.Add(address, decimals)
	return decimals, false, nil
}

func isNoSuchMethodError(err error) bool {
	errStr := strings.ToLower(err.Error())
	possibleErrorMessages := []string{
		"execution reverted",
		"invalid opcode",
		"abi: attempting to unmarshal an empty string while arguments are expected",
		"no contract code at given address",
		"invalid jump destination",
	}
	for _, errMsg := range possibleErrorMessages {
		if strings.Contains(errStr, errMsg) {
			return true
		}
	}
	return false
}

func (c *Client) GetFactory(address common.Address) (common.Address, error) {
	var factory common.Address
	var err error
	ctx := context.Background()
	err = c.withRetry(ctx, "GetFactory", func() error {
		var e error
		factory, e = c.getFactory(address)
		if e != nil {
			c.logger.Warn("error getting factory", zap.String("token", address.Hex()), zap.Error(e))
		}
		return e
	})
	return factory, err
}

func (c *Client) getFactory(address common.Address) (common.Address, error) {
	if factory, ok := c.factoryCache.Get(address); ok {
		return factory, nil
	}

	methods := []func(common.Address) (common.Address, error){
		c.uniswapV2GetFactory,
		c.uniswapV1GetFactory,
	}

	for _, method := range methods {
		factory, err := method(address)
		if err != nil {
			continue
		}
		c.factoryCache.Add(address, factory)
		return factory, nil
	}

	return common.Address{}, fmt.Errorf("no factory found")
}

func (c *Client) uniswapV2GetFactory(address common.Address) (common.Address, error) {
	const abiJSON = `[{"constant":true,"inputs":[],"name":"factory","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return common.Address{}, err
	}

	data, err := parsedABI.Pack("factory")
	if err != nil {
		return common.Address{}, err
	}

	if err := c.waitRequests(context.Background(), 1); err != nil {
		return common.Address{}, err
	}

	response, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: data,
	}, nil)
	if err != nil {
		return common.Address{}, err
	}

	result, err := parsedABI.Unpack("factory", response)
	if err != nil {
		return common.Address{}, err
	}
	resCasted, ok := result[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("result is not a common.Address")
	}
	return resCasted, nil
}

func (c *Client) uniswapV1GetFactory(address common.Address) (common.Address, error) {
	const abiJSON = `[{"constant":true,"inputs":[],"name":"factoryAddress","outputs":[{"name":"out","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return common.Address{}, err
	}

	data, err := parsedABI.Pack("factoryAddress")
	if err != nil {
		return common.Address{}, err
	}

	if err := c.waitRequests(context.Background(), 1); err != nil {
		return common.Address{}, err
	}

	response, err := c.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: data,
	}, nil)
	if err != nil {
		return common.Address{}, err
	}

	result, err := parsedABI.Unpack("factoryAddress", response)
	if err != nil {
		return common.Address{}, err
	}

	resCasted, ok := result[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("result is not a common.Address")
	}
	return resCasted, nil
}

func (c *Client) GetIsContract(address common.Address) (bool, error) {
	var isContract bool
	var err error
	ctx := context.Background()
	err = c.withRetry(ctx, "GetIsContract", func() error {
		var e error
		isContract, e = c.getIsContract(address)
		if e != nil {
			c.logger.Warn("error getting is contract", zap.String("address", address.Hex()), zap.Error(e))
		}
		return e
	})
	return isContract, err
}

func (c *Client) getIsContract(address common.Address) (bool, error) {
	if isContract, ok := c.isContractCache.Get(address); ok {
		return isContract, nil
	}
	var isContract bool
	var err error
	ctx := context.Background()
	if err := c.waitRequests(ctx, 1); err != nil {
		return false, err
	}
	bytecode, err := c.ethClient.CodeAt(ctx, address, nil)
	if err != nil {
		return false, err
	}
	isContract = len(bytecode) > 0
	c.isContractCache.Add(address, isContract)
	return isContract, err
}

func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	blockNumber, err := c.ethClient.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(blockNumber)), nil
}

func (c *Client) GetTransaction(tx common.Hash) (*types.Transaction, error) {
	var transaction *types.Transaction
	var err error
	ctx := context.Background()
	err = c.withRetry(ctx, "GetTransaction", func() error {
		var e error
		transaction, e = c.getTransaction(tx)
		return e
	})
	return transaction, err
}

func (c *Client) getTransaction(tx common.Hash) (*types.Transaction, error) {
	ctx := context.Background()
	if err := c.waitRequests(ctx, 1); err != nil {
		return nil, err
	}
	transaction, _, err := c.ethClient.TransactionByHash(ctx, tx)
	if err != nil {
		return nil, err
	}
	return transaction, err
}

func (c *Client) GetReceipt(hash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	ctx := context.Background()
	err = c.withRetry(ctx, "GetReceipt", func() error {
		var e error
		receipt, e = c.getReceipt(hash)
		if e != nil {
			c.logger.Warn("error getting receipt", zap.String("hash", hash.Hex()), zap.Error(e))
		}
		return e
	})
	return receipt, err
}

func (c *Client) getReceipt(hash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()

	if err := c.waitRequests(ctx, 1); err != nil {
		return nil, err
	}

	receipt, err := c.ethClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (c *Client) GetBlockHeader(hash common.Hash) (*types.Header, error) {
	var blockHeader *types.Header
	var err error
	ctx := context.Background()
	err = c.withRetry(ctx, "GetBlockHeader", func() error {
		var e error
		blockHeader, e = c.getBlockHeader(hash)
		if e != nil {
			c.logger.Warn("error getting block header", zap.String("hash", hash.Hex()), zap.Error(e))
		}
		return e
	})
	return blockHeader, err
}

func (c *Client) getBlockHeader(hash common.Hash) (*types.Header, error) {
	ctx := context.Background()
	if c.blocksSem != nil {
		if err := c.blocksSem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		defer c.blocksSem.Release(1)
	}
	if err := c.waitRequests(ctx, 1); err != nil {
		return nil, err
	}
	blockHeader, err := c.ethClient.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return blockHeader, nil
}

func (c *Client) DebugTraceTransaction(txHash common.Hash) (*models.FullTraceResult, error) {
	var result *models.FullTraceResult
	var err error
	ctx := context.Background()
	err = c.withRetry(ctx, "DebugTraceTransaction", func() error {
		var e error
		result, e = c.debugTraceTransaction(txHash)
		if e != nil {
			c.logger.Warn("error getting trace",
				zap.String("tx_hash", txHash.Hex()),
				zap.Error(e),
			)
		}
		return e
	})
	return result, err
}

func (c *Client) debugTraceTransaction(txHash common.Hash) (*models.FullTraceResult, error) {
	ctx := context.Background()
	if err := c.waitRequests(ctx, 1); err != nil {
		return nil, err
	}

	traceFrame := models.FullTraceResult{}

	var raw json.RawMessage
	err := c.rpcClient.CallContext(ctx, &raw, "debug_traceTransaction", txHash.Hex(), map[string]interface{}{"tracer": "callTracer"})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &traceFrame)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("got trace", zap.Any("trace", traceFrame))

	return &traceFrame, nil
}
