package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/evmclient"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/processors/polymarketProcessor"
	"github.com/monolit-network/xlayer-indexer/shared/db/clickhouse"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres"
	"github.com/monolit-network/xlayer-indexer/util/promise"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func initRpcs(chain string) (*ethclient.Client, *ethclient.Client, *rpc.Client, error) {
	rpcURL := os.Getenv(fmt.Sprintf("EVM_RPC_URL_%s", strings.ToUpper(chain)))
	if rpcURL == "" {
		return nil, nil, nil, errors.New("EVM_RPC_URL is not set")
	}
	wsURL := os.Getenv(fmt.Sprintf("EVM_WS_URL_%s", strings.ToUpper(chain)))
	if wsURL == "" {
		return nil, nil, nil, errors.New("EVM_WS_URL is not set")
	}
	t := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{
		Transport: t,
		Timeout:   300 * time.Second,
	}

	rpcClient, err := rpc.DialOptions(context.Background(), rpcURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	defer ethClient.Close()

	wsEthClient, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, nil, nil, err
	}
	defer wsEthClient.Close()

	return ethClient, wsEthClient, rpcClient, nil
}

type txWithFee struct {
	TxHash common.Hash
	Fee    *big.Int
}

func main() {
	godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	dbClient, err := postgres.NewPostgresClientFromEnv()
	if err != nil {
		logger.Fatal("error creating database client", zap.Error(err))
	}
	defer dbClient.Close()

	clickClient, err := clickhouse.NewClickhouseClientFromEnv(logger)
	if err != nil {
		logger.Fatal("error creating clickhouse client", zap.Error(err))
	}
	defer clickClient.Close()

	chain := models.ChainPolygon

	ethClient, wsEthClient, rpcClient, err := initRpcs(string(chain))
	if err != nil {
		logger.Fatal("error initializing rpc clients", zap.Error(err))
	}
	defer ethClient.Close()
	defer wsEthClient.Close()
	defer rpcClient.Close()

	evmClient := evmclient.NewClient(ethClient, wsEthClient, rpcClient, logger, models.Chain(chain), evmclient.EvmClientOptionsFromEnv(models.Chain(chain))...)

	processor := polymarketProcessor.NewProcessor(evmClient, clickClient, nil, dbClient, logger)
	defer processor.Stop()

	defer logger.Info("finished")

	txsWithFees := make(chan txWithFee)
	// rows, err := clickClient.Conn.Query(context.Background(), "SELECT tx_hash, sum(fee) as tfee FROM evm.polymarket_order_events where fee > 0 group by tx_hash order by max(block_number) asc")
	// if err != nil {
	// 	logger.Fatal("error querying tx hashes with fees", zap.Error(err))
	// }
	logger.Info("querying tx hashes with fees done")
	go func() {
		defer close(txsWithFees)

		txsWithFees <- txWithFee{TxHash: common.HexToHash("0x6d530bd8b0cd9b201c26c89851cae14b486fb6d080c4396714b54886864d6df5"), Fee: big.NewInt(100)}

		// defer rows.Close()
		// for rows.Next() {
		// 	var txHashRaw string
		// 	var tfee *big.Int
		// 	err := rows.Scan(&txHashRaw, &tfee)
		// 	if err != nil {
		// 		logger.Fatal("error scanning tx hashes with fees", zap.Error(err))
		// 	}
		// 	txsWithFees <- txWithFee{TxHash: common.HexToHash(txHashRaw), Fee: tfee}
		// }
	}()

	receipts := make(chan promise.Promise[*types.Receipt], 128)
	go func() {
		defer close(receipts)
		for txWithFee := range txsWithFees {
			prom := promise.NewPromise[*types.Receipt]()
			receipts <- prom
			go func() {
				receipt, err := evmClient.GetReceipt(txWithFee.TxHash)
				if err != nil {
					prom.Reject(err)
					return
				}
				prom.Resolve(receipt)
			}()
		}
	}()

	processed := 0
	start := time.Now()
	for receiptProm := range receipts {
		receipt, err := receiptProm.Wait()
		if err != nil {
			logger.Error("error waiting for receipt", zap.Error(err))
			break
		}
		err = processor.ProcessLogsFromTx(receipt.Logs)
		if err != nil {
			logger.Error("error processing logs from transaction", zap.Error(err))
			break
		}
		processed++
		if processed%1 == 0 {
			logger.Info("processed", zap.Int("processed", processed), zap.Duration("time", time.Since(start)))
		}
		if processed%1000 == 0 {
			processor.LogTimingsCurrent()
		}
	}
	logger.Info("processed", zap.Int("processed", processed))
}
