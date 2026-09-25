package models

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	UniswapV1FactoryAddressEth = common.HexToAddress("0xc0a47dfe034b400b47bdad5fecda2621de6c4d95")
	UniswapV2FactoryAddressEth = common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	UniswapV3FactoryAddressEth = common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
)
var (
	UniswapV2FactoryAddresses = map[Chain]common.Address{
		ChainETH:       common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"),
		ChainArbitrum:  common.HexToAddress("0xf1D7CC64Fb4452F05c498126312eBE29f30Fbcf9"),
		ChainAvalanche: common.HexToAddress("0x9e5A52f57b3038F1B8EeE45F28b3C1967e22799C"),
		ChainBSC:       common.HexToAddress("0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6"),
		ChainBase:      common.HexToAddress("0x8909Dc15e40173Ff4699343b6eB8132c65e18eC6"),
		ChainOptimism:  common.HexToAddress("0x0c3c1c532F1e39EdF36BE9Fe0bE1410313E074Bf"),
		ChainPolygon:   common.HexToAddress("0x9e5A52f57b3038F1B8EeE45F28b3C1967e22799C"),
	}
	UniswapV3FactoryAddresses = map[Chain]common.Address{
		ChainETH:       common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984"),
		ChainArbitrum:  common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984"),
		ChainAvalanche: common.HexToAddress("0x740b1c1de25031C31FF4fC9A62f554A55cdC1baD"),
		ChainBSC:       common.HexToAddress("0xdB1d10011AD0Ff90774D0C6Bb92e5C5c8b4461F7"),
		ChainBase:      common.HexToAddress("0x33128a8fC17869897dcE68Ed026d694621f6FDfD"),
		ChainOptimism:  common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984"),
		ChainPolygon:   common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984"),
	}
)

const (
	ETHTokensDecimals uint8 = 18
)

var (
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

	ETHTokenAddress    = common.HexToAddress("0x0000000000000000000000000000000000000000")
	POLLogTokenAddress = common.HexToAddress("0x0000000000000000000000000000000000001010")
	WETHTokenAddress   = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

	WPOLTokenAddress = common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270")
	WBNBTokenAddress = common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")

	CurveETHTokenAddress   = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	WBaseTokenAddress      = common.HexToAddress("0x4200000000000000000000000000000000000006")
	WRobinhoodTokenAddress = common.HexToAddress("0x0Bd7D308f8E1639FAb988df18A8011f41EAcAD73")

	PolymarketNegRiskWrappedCollateralAddress = common.HexToAddress("0x3a3bd7bb9528e159577f7c2e685cc81a765002e2")
)

var (
	PolymarketCTFExchangeAddress          = common.HexToAddress("0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E")
	PolymarketCTFExchangeAddressV2        = common.HexToAddress("0xE111180000d2663C0091e4f400237545B87B996B")
	PolymarketConditionalTokensAddress    = common.HexToAddress("0x4d97dcd97ec945f40cf65f87097ace5ea0476045")
	PolymarketNegRiskCTFExchangeAddress   = common.HexToAddress("0xC5d563A36AE78145C45a50134d48A1215220f80a")
	PolymarketNegRiskCTFExchangeAddressV2 = common.HexToAddress("0xe2222d279d744050d28e00520010520000310F59")
	PolymarketNegRiskAdapterAddress       = common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296")

	PolymarketCtfCollateralAdapterAddress         = common.HexToAddress("0xADa100874d00e3331D00F2007a9c336a65009718")
	PolymarketNegRiskCtfCollateralAdapterAddress  = common.HexToAddress("0xAdA200001000ef00D07553cEE7006808F895c6F1")
	PolymarketNegRiskCtfCollateralAdapterAddress2 = common.HexToAddress("0xada2005600dec949baf300f4c6120000bdb6eaab")

	PolymarketUmaCtfAdapterV2Address = common.HexToAddress("0x65070BE91477460D8A7AeEb94ef92fe056C2f2A7")

	UmaOptimisticOracleV2Address                  = common.HexToAddress("0xee3afe347d5c74317041e2618c49534daf887c24")
	UmaOptimisticOracleV3Address                  = common.HexToAddress("0x5953f2538F613E05bAED8A5AeFa8e6622467AD3D")
	UmaPolymarketManagedOptimisticOracleV2Address = common.HexToAddress("0x2C0367a9DB231dDeBd88a94b4f6461a6e47C58B1")

	PolymarketFeeModuleAddress        = common.HexToAddress("0xE3f18aCc55091e2c48d883fc8C8413319d4Ab7b0")
	PolymarketNegRiskFeeModuleAddress = common.HexToAddress("0x78769D50Be1763ed1CA0D5E878D93f05aabff29e")
)

var PolymarketCTFAdaptersAddresses = []common.Address{
	common.HexToAddress("0x71392e133063cc0d16f40e1f9b60227404bc03f7"), // 3.0?
	common.HexToAddress("0x6A9D222616C90FcA5754cd1333cFD9b7fb6a4F74"), // 2?
	common.HexToAddress("0x2f5e3684cb1f318ec51b00edba38d79ac2c0aa9d"), // V3 on polyscan
	common.HexToAddress("0x157Ce2d672854c848c9b79C49a8Cc6cc89176a49"), // 3.1
	common.HexToAddress("0x65070BE91477460D8A7AeEb94ef92fe056C2f2A7"), // no idea, primary
	common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"), // Neg Risk Adapter
	PolymarketCtfCollateralAdapterAddress,
	PolymarketNegRiskCtfCollateralAdapterAddress,
	PolymarketNegRiskCtfCollateralAdapterAddress2,
}

var PolymarketCTFAdaptersAddressesMap = map[common.Address]struct{}{
	common.HexToAddress("0x71392e133063cc0d16f40e1f9b60227404bc03f7"): {},
	common.HexToAddress("0x6A9D222616C90FcA5754cd1333cFD9b7fb6a4F74"): {},
	common.HexToAddress("0x2f5e3684cb1f318ec51b00edba38d79ac2c0aa9d"): {},
	common.HexToAddress("0x157Ce2d672854c848c9b79C49a8Cc6cc89176a49"): {},
	common.HexToAddress("0x65070BE91477460D8A7AeEb94ef92fe056C2f2A7"): {},
	common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"): {},
	PolymarketCtfCollateralAdapterAddress:                             {},
	PolymarketNegRiskCtfCollateralAdapterAddress:                      {},
	PolymarketNegRiskCtfCollateralAdapterAddress2:                     {},
}

var UniswapV4PoolManagerAddresses = map[Chain]common.Address{
	ChainRobinhood: common.HexToAddress("0x8366a39cc670b4001a1121b8f6a443a643e40951"),
	ChainXlayer:    common.HexToAddress("0x360e68faccca8ca495c1b759fd9eee466db9fb32"),
}

var UniswapV4SwapEventTopic = common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f")

var UniswapV3SwapEventTopic = common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67")

var UniswapV2SwapEventTopic = common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822")

// ERC-4337 EntryPoint UserOperationEvent(userOpHash, sender, paymaster, ...).
// Matched by topic only - agnostic to EntryPoint address/version.
// iZiSwap (discretized-liquidity AMM, X Layer flagship DEX) Swap event.
var IziSwapEventTopic = common.HexToHash("0xcd3829a3813dc3cdd188fd3d01dcf3268c16be2fdd2dd21d0665418816e46062")

// DODO route OrderHistory(fromToken, toToken, sender, fromAmount, returnAmount) - X Layer aggregator/router.
var DodoRouteEventTopic = common.HexToHash("0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db")

// Curve TokenExchange(buyer, sold_id, tokens_sold, bought_id, tokens_bought).
var CurveExchangeEventTopic = common.HexToHash("0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140")

var UserOperationEventTopic = common.HexToHash("0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f")
