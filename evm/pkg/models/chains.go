package models

import "github.com/ethereum/go-ethereum/common"

type Chain string

const (
	ChainETH       Chain = "eth"
	ChainBSC       Chain = "bsc"
	ChainPolygon   Chain = "polygon"
	ChainArbitrum  Chain = "arbitrum"
	ChainOptimism  Chain = "optimism"
	ChainBase      Chain = "base"
	ChainAvalanche Chain = "avalanche"
	ChainFantom    Chain = "fantom"
	ChainCronos    Chain = "cronos"
	ChainMoonbeam  Chain = "moonbeam"
	ChainRobinhood Chain = "robinhood"
	ChainXlayer    Chain = "xlayer"
)

func ValidateChain(chain string) (Chain, bool) {
	switch Chain(chain) {
	case ChainETH,
		ChainBSC,
		ChainPolygon,
		ChainArbitrum,
		ChainOptimism,
		ChainBase,
		ChainAvalanche,
		ChainFantom,
		ChainCronos,
		ChainMoonbeam,
		ChainRobinhood,
		ChainXlayer:
		return Chain(chain), true
	default:
		return "", false
	}
}

type ChainId uint64

const (
	ChainIdETH       ChainId = 1
	ChainIdBSC       ChainId = 56
	ChainIdPolygon   ChainId = 137
	ChainIdArbitrum  ChainId = 42161
	ChainIdOptimism  ChainId = 10
	ChainIdBase      ChainId = 8453
	ChainIdAvalanche ChainId = 43114
	ChainIdRobinhood ChainId = 4663
	ChainIdXlayer    ChainId = 196
)

var IdToChain = map[ChainId]Chain{
	ChainIdETH:       ChainETH,
	ChainIdBSC:       ChainBSC,
	ChainIdPolygon:   ChainPolygon,
	ChainIdArbitrum:  ChainArbitrum,
	ChainIdOptimism:  ChainOptimism,
	ChainIdBase:      ChainBase,
	ChainIdAvalanche: ChainAvalanche,
	ChainIdRobinhood: ChainRobinhood,
	ChainIdXlayer:    ChainXlayer,
}

var ChainToID = map[Chain]ChainId{
	ChainETH:       ChainIdETH,
	ChainBSC:       ChainIdBSC,
	ChainPolygon:   ChainIdPolygon,
	ChainArbitrum:  ChainIdArbitrum,
	ChainOptimism:  ChainIdOptimism,
	ChainBase:      ChainIdBase,
	ChainAvalanche: ChainIdAvalanche,
	ChainRobinhood: ChainIdRobinhood,
	ChainXlayer:    ChainIdXlayer,
}

var ChainIDToWrappedNativeToken = map[ChainId]common.Address{
	ChainIdXlayer: common.HexToAddress("0xe538905cf8410324e03a5a23c1c177a474d59b2b"), // WOKB
	ChainIdETH:       WETHTokenAddress,
	ChainIdPolygon:   WPOLTokenAddress,
	ChainIdArbitrum:  WETHTokenAddress,
	ChainIdOptimism:  WETHTokenAddress,
	ChainIdBase:      WBaseTokenAddress,
	ChainIdBSC:       WBNBTokenAddress,
	ChainIdRobinhood: WRobinhoodTokenAddress,
}
