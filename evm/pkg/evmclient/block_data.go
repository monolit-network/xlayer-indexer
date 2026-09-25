package evmclient

import (
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockReceipts map[common.Hash]*types.Receipt

type BlockTraces map[common.Hash]*models.FullTraceResult
