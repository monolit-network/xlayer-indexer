package polymarket

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketCTFExchangeV2"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/abis/polymarketFPMM"
	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	sharedModels "github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetQuestionID(t *testing.T) {
	p := NewCTFParser(nil, nil)

	hexStr := "713A207469746C653A204E464C2053756E6461793A20537465656C6572732076732E20526176656E732C206465736372697074696F6E3A20496E20746865207570636F6D696E67204E464C2067616D65207363686564756C656420666F72204A616E7561727920312C20323032333A0A0A496620746865205069747473627572676820537465656C6572732077696E2C2074686973206D61726B65742077696C6C207265736F6C766520746F2022537465656C657273222E0A0A4966207468652042616C74696D6F726520526176656E732077696E2C2074686973206D61726B65742077696C6C207265736F6C766520746F2022526176656E73222E0A0A496620746869732067616D6520656E647320696E2061207469652C2074686973206D61726B65742077696C6C207265736F6C76652035302D35302E20496620746869732067616D65206973206E6F7420636F6D706C65746564206279204A616E7561727920382C20323032332C2031313A35393A353920504D2045542C2074686973206D61726B65742077696C6C207265736F6C76652035302D35302E207265735F646174613A2070313A20302C2070323A20312C2070333A20302E352E20576865726520703120636F72726573706F6E647320746F20526176656E732C20703220746F20537465656C6572732C20703320746F20756E6B6E6F776E2F35302D35302C696E697469616C697A65723A39313433306361643264333937353736363439393731376661306436366137386438313465356335"
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("failed to decode hex string: %v", err)
	}
	questionId := p.questionIDFromAncillaryData(data)
	assert.Equal(t, "0x3b84150af09f630ff6224cfa94edaabba9bd7717bd7222bb53112842c32d922c", questionId.Hex())
}

func TestParseAncillaryData(t *testing.T) {
	hexStr := "713a207469746c653a205370726561643a2053756e7320282d312e35292c206465736372697074696f6e3a20496e20746865207570636f6d696e67204e42412067616d652c207363686564756c656420666f72204f63746f62657220323720617420393a303020504d2045543a0a0a54686973206d61726b65742077696c6c207265736f6c766520746f202253756e7322206966207468652053756e732077696e207468652067616d652062792032206f72206d6f726520706f696e74732e0a0a4f74686572776973652c2074686973206d61726b65742077696c6c207265736f6c766520746f20224a617a7a222e204966207468652067616d6520656e647320696e2061207469652c2074686973206d61726b65742077696c6c207265736f6c766520746f20224a617a7a222e0a0a4966207468652067616d6520697320706f7374706f6e65642c2074686973206d61726b65742077696c6c2072656d61696e206f70656e20756e74696c207468652067616d6520686173206265656e20636f6d706c657465642e204966207468652067616d652069732063616e63656c656420656e746972656c792c2077697468206e6f206d616b652d75702067616d652c2074686973206d61726b65742077696c6c207265736f6c76652035302d35302e206d61726b65745f69643a20363530383332207265735f646174613a2070313a20302c2070323a20312c2070333a20302e352e20576865726520703120636f72726573706f6e647320746f204a617a7a2c20703220746f2053756e732c20703320746f20756e6b6e6f776e2f35302d35302e2055706461746573206d61646520627920746865207175657374696f6e2063726561746f7220766961207468652062756c6c6574696e20626f61726420617420307836353037304245393134373734363044384137416545623934656639326665303536433266324137206173206465736372696265642062792068747470733a2f2f706f6c79676f6e7363616e2e636f6d2f74782f3078613134663031623131356334393133363234666333663530386639363066346465613235323735386537336332386635663037663865313964376263613036362073686f756c6420626520636f6e736964657265642e2c696e697469616c697a65723a39313433306361643264333937353736363439393731376661306436366137386438313465356335"
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("failed to decode hex string: %v", err)
	}

	p := NewCTFParser(nil, nil)

	parsedData := p.ParseAncillaryData(data)

	assert.Equal(t, "Spread: Suns (-1.5)", *parsedData.Question)

	assert.Equal(t, `In the upcoming NBA game, scheduled for October 27 at 9:00 PM ET:

This market will resolve to "Suns" if the Suns win the game by 2 or more points.

Otherwise, this market will resolve to "Jazz". If the game ends in a tie, this market will resolve to "Jazz".

If the game is postponed, this market will remain open until the game has been completed. If the game is canceled entirely, with no make-up game, this market will resolve 50-50`, *parsedData.Description)

	assert.Equal(t, `p1: 0, p2: 1, p3: 0.5. Where p1 corresponds to Jazz, p2 to Suns, p3 to unknown/50-50. Updates made by the question creator via the bulletin board at 0x65070BE91477460D8A7AeEb94ef92fe056C2f2A7 as described by https://polygonscan.com/tx/0xa14f01b115c4913624fc3f508f960f4dea252758e73c28f5f07f8e19d7bca066 should be considered`, *parsedData.ResData)

	assert.Equal(t, int64(650832), *parsedData.MarketID)

	assert.Equal(t, `91430cad2d3975766499717fa0d66a78d814e5c5`, *parsedData.Initializer)
}

func TestGetTokenId(t *testing.T) {
	for _, tc := range tokenCalculationTestCases {
		t.Run(fmt.Sprintf("conditionId: %s, parentCollectionId: %s, indexSet: %s, collateralToken: %s", tc.conditionId.Hex(), tc.parentCollectionId.Hex(), tc.indexSet.String(), tc.collateralToken.Hex()), func(t *testing.T) {
			collectionId := getCollectionId(tc.parentCollectionId, tc.conditionId, tc.indexSet)
			if tc.expectedCollectionId != (common.Hash{}) {
				assert.Equal(t, tc.expectedCollectionId.Hex(), collectionId.Hex())
			}
			positionId := getTokenIdPartition(tc.conditionId, tc.indexSet, tc.parentCollectionId, tc.collateralToken)
			if positionId.String() != tc.expectedPositionId {
				// try with wrapped collateral token
				positionId2 := getTokenIdPartition(tc.conditionId, tc.indexSet, tc.parentCollectionId, wrappedCollateralToken)
				if positionId2.String() == tc.expectedPositionId {
					assert.Equal(t, tc.expectedPositionId, positionId2.String())
					return
				}
			}
			assert.Equal(t, tc.expectedPositionId, positionId.String())
		})
	}
}

func TestCreateContract(t *testing.T) {
	contract, err := polymarketFPMM.NewPolymarketFPMM(common.Address{}, nil)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	assert.NotNil(t, contract)
}

func TestParseOrderFilledLogV2NormalizesBuySide(t *testing.T) {
	p := NewCTFParser(nil, nil)
	tokenID := big.NewInt(12345)
	log := newOrderFilledV2Log(t, 0, tokenID)

	orderFilled, err := p.ParseOrderFilledLog(log)

	require.NoError(t, err)
	require.NotNil(t, orderFilled)
	assert.Equal(t, "0", orderFilled.MakerAssetID.String())
	assert.Equal(t, tokenID.String(), orderFilled.TakerAssetID.String())
	assert.Equal(t, "100", orderFilled.MakerAmountFilled.String())
	assert.Equal(t, "200", orderFilled.TakerAmountFilled.String())
	assert.Equal(t, "3", orderFilled.Fee.String())
}

func TestParseOrderFilledLogV2NormalizesSellSide(t *testing.T) {
	p := NewCTFParser(nil, nil)
	tokenID := big.NewInt(67890)
	log := newOrderFilledV2Log(t, 1, tokenID)

	orderFilled, err := p.ParseOrderFilledLog(log)

	require.NoError(t, err)
	require.NotNil(t, orderFilled)
	assert.Equal(t, tokenID.String(), orderFilled.MakerAssetID.String())
	assert.Equal(t, "0", orderFilled.TakerAssetID.String())
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000a11"), orderFilled.Maker)
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000b22"), orderFilled.Taker)
}

func newOrderFilledV2Log(t *testing.T, side uint8, tokenID *big.Int) *types.Log {
	t.Helper()

	contractABI, err := polymarketCTFExchangeV2.PolymarketCTFExchangeV2MetaData.GetAbi()
	require.NoError(t, err)

	data, err := contractABI.Events["OrderFilled"].Inputs.NonIndexed().Pack(
		side,
		tokenID,
		big.NewInt(100),
		big.NewInt(200),
		big.NewInt(3),
		[32]byte{},
		[32]byte{},
	)
	require.NoError(t, err)

	maker := common.HexToAddress("0x0000000000000000000000000000000000000a11")
	taker := common.HexToAddress("0x0000000000000000000000000000000000000b22")
	return &types.Log{
		Address: models.PolymarketCTFExchangeAddressV2,
		Topics: []common.Hash{
			models.PolymarketOrderFilledEventV2SelectorHash,
			common.HexToHash("0x0123"),
			common.BytesToHash(maker.Bytes()),
			common.BytesToHash(taker.Bytes()),
		},
		Data: data,
	}
}

func TestParseNegRiskPositionConverted(t *testing.T) {
	p := NewCTFParser(nil, nil)
	raw := &sharedModels.PolymarketNegRiskPositionsConvertedRaw{
		LogAddress:  common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"),
		Stakeholder: common.HexToAddress("0x59d4a3B45519B25659549d513EadFc58F80897E7"),
		MarketId:    common.HexToHash("0xe3b1bc389210504ebcb9cffe4b0ed06ccac50561e0f24abb6379984cec030f00"),
		IndexSet:    big.NewInt(1),
		Amount:      big.NewInt(531000000),
	}

	marketDataHash := common.HexToHash("0x11010000000000000000000071523d0f655b41e805cec45b17163f528b59b820")
	marketData := ParseMarketData(marketDataHash)
	assert.Equal(t, marketData.QuestionCount, uint8(17))
	assert.Equal(t, marketData.Determined, true)
	assert.Equal(t, marketData.Result, uint8(0))
	assert.Equal(t, marketData.FeeBips, uint16(0))
	assert.Equal(t, marketData.Oracle, common.HexToAddress("0x71523d0f655B41E805Cec45b17163f528B59B820"))

	result := p.ParseNegRiskPositionConverted(raw, marketData)
	assert.NotNil(t, result)
	assert.Len(t, result.NoTokensGiven, 1)
	assert.Len(t, result.YesTokensReceived, 16)
	assert.Equal(t, result.CollateralOut.String(), "0")
	assert.Equal(t, result.AmountIn.String(), "531000000")
	assert.Equal(t, result.AmountOut.String(), "531000000")
	assert.Equal(t, result.NoTokensGiven[0].TokenID.String(), "48331043336612883890938759509493159234755048973500640148014422747788308965732")
}

func TestParseNegRiskPositionConverted2(t *testing.T) {
	p := NewCTFParser(nil, nil)
	raw := &sharedModels.PolymarketNegRiskPositionsConvertedRaw{
		LogAddress:  common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"),
		Stakeholder: common.HexToAddress("0x3057De5eFC5fce194b0B034F08Ad0beeA9e77a96"),
		MarketId:    common.HexToHash("0xdfb995970db1479de541ff75ce9b6bfeb6ba104867f6efdd996624284b86e600"),
		IndexSet:    big.NewInt(24),
		Amount:      big.NewInt(1234563),
	}

	marketDataHash := common.HexToHash("0x0b0000000000000000000000661992aebf6becf7ba5abb66f6b0bf62aa7a2e93")
	marketData := ParseMarketData(marketDataHash)
	assert.Equal(t, marketData.QuestionCount, uint8(11))
	assert.Equal(t, marketData.Determined, false)
	assert.Equal(t, marketData.Result, uint8(0))
	assert.Equal(t, marketData.FeeBips, uint16(0))
	assert.Equal(t, marketData.Oracle, common.HexToAddress("0x661992aebf6BecF7BA5abB66f6b0Bf62Aa7a2E93"))

	result := p.ParseNegRiskPositionConverted(raw, marketData)
	assert.NotNil(t, result)
	assert.Len(t, result.NoTokensGiven, 2)
	assert.Len(t, result.YesTokensReceived, 9)
	assert.Equal(t, result.CollateralOut.String(), "1234563")
	assert.Equal(t, result.AmountIn.String(), "1234563")
	assert.Equal(t, result.AmountOut.String(), "1234563")

	tokensIn := []string{
		"33199307015573997160127175404470082428284143318133448438933310663153719364725",
		"45222491279042081245407854458168351681052325467963981864630237227066398470730",
	}
	resTokensInStrings := []string{}
	for _, token := range result.NoTokensGiven {
		resTokensInStrings = append(resTokensInStrings, token.TokenID.String())
	}

	for _, token := range tokensIn {
		assert.Contains(t, resTokensInStrings, token)
	}
	for _, token := range result.NoTokensGiven {
		assert.Contains(t, tokensIn, token.TokenID.String())
	}

	tokensOut := []string{
		"78486200154613929876692928807891105873702544860723824991241378920152908263156",
		"35298931584333277389584810771117059758462428354136848157702028766863432403460",
		"35938186546229623759393467280134715531628158562905292092671665137818934121606",
		"30089169520467070760523489348420117980565835371148321235916179225385808899155",
		"8761070964184698699025585216982900385117372285952024983712590631481019444656",
		"44737098827247974321827418719741119757420788395060220700795281017930288867898",
		"15399260758602440881971021750462083112278813984199184535811533651712580262501",
		"70338904209343078745262473571492670059340827287527303927370022368045172526237",
		"4482710125732668075283135166385308531089145073765674922551403368866526624018",
	}

	resTokensOutStrings := make([]string, 0)
	for _, token := range result.YesTokensReceived {
		resTokensOutStrings = append(resTokensOutStrings, token.TokenID.String())
	}

	for _, token := range tokensOut {
		assert.Contains(t, resTokensOutStrings, token)
	}
	for _, token := range result.YesTokensReceived {
		assert.Contains(t, tokensOut, token.TokenID.String())
	}

}

var (
	wrappedCollateralToken    = models.PolymarketNegRiskWrappedCollateralAddress
	defaultParentCollectionId = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	defaultCollateralToken    = common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174")
)

var tokenCalculationTestCases = []struct {
	conditionId        common.Hash
	parentCollectionId common.Hash
	indexSet           *big.Int
	collateralToken    common.Address

	expectedCollectionId common.Hash
	expectedPositionId   string
}{
	{
		conditionId:        common.HexToHash("0x83d7fc5f59eaf73060a742c1734fdf9ba0e7dc07d27846d9210f07beb5f319e7"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedCollectionId: common.HexToHash("0x669a6171381de382fa50f38af3347b7c9a0c06eb059b75cfd2ce45766393751f"),
		expectedPositionId:   "84003763053781121077095547621708347905507570798446118351127377081446837298612",
	},
	{
		conditionId:        common.HexToHash("0x178179A1B62E29AD7D7FD2A07B2C9AE013553C15554D48024BE9793998D2AFDA"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    wrappedCollateralToken,

		expectedCollectionId: common.HexToHash("0x05211898ee51903596df4ea978878a9a0312d0eb00beb6a7130879b23d1a3191"),
		expectedPositionId:   "43718257888068854502235591051205617681318972030595749863641195355046154403333",
	},
	{
		conditionId:        common.HexToHash("0x178179A1B62E29AD7D7FD2A07B2C9AE013553C15554D48024BE9793998D2AFDA"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    wrappedCollateralToken,

		expectedCollectionId: common.HexToHash("0x0c1d2326e0fafa29a6c463374912293adf293dd6e8b1989ccae504827805330e"),
		expectedPositionId:   "31251112292863300032495665804626754516557214325398496754022883012056148398604",
	},
	{
		conditionId:        common.HexToHash("0x464EAC41FE4EE9E4317FA68B01D922D7B1E48EDD3584C3EC30A87308CD6C684B"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    wrappedCollateralToken,

		expectedCollectionId: common.HexToHash("0x593cb4061eaadb679edc3b6a1fcd62f6f4b21d1c8acc3c2e1519a17354555c0d"),
		expectedPositionId:   "110105962709742756399914071345068468355063301951920723493736436405184596241376",
	},

	{
		conditionId:        common.HexToHash("0x0002a45f7736686e98f5e6476a3d51dd48db232f49115312a07b047c5272eff6"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "25025646619520528368956414960932415270214002600335105407720414855152573043376",
	},
	{
		conditionId:        common.HexToHash("0x0002a45f7736686e98f5e6476a3d51dd48db232f49115312a07b047c5272eff6"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "23730767560780769504439203266635963984946402417129664205405764774491860662719",
	},
	{
		conditionId:        common.HexToHash("0x000388290ce64d1e5e98a246b96f6c3af5e731d6d0f27cdd88dd6e1791808c09"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "24360795570856515273548016662246023231381256879917916663180194532747151550114",
	},
	{
		conditionId:        common.HexToHash("0x000388290ce64d1e5e98a246b96f6c3af5e731d6d0f27cdd88dd6e1791808c09"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "59142414867305534799269650104152966959696800637607762359463856446288219390165",
	},
	{
		conditionId:        common.HexToHash("0x0003c55b045243989673c96d6df12e4f6a74ad6b9b561b7d5e2f8cbac3a03417"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "112960553284197531944859596497200846251766583439406059618502468183350551567495",
	},
	{
		conditionId:        common.HexToHash("0x0003c55b045243989673c96d6df12e4f6a74ad6b9b561b7d5e2f8cbac3a03417"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "8996758215026844954110546669110312646088271735158037541715840969582824253641",
	},
	{
		conditionId:        common.HexToHash("0x0004866a1cd8e94ac08e4ff562f038a8d36549cb148f0bb17f64b3e3b96d346d"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "115133372593718990661522741313796381972609948368753012572280437659997828364758",
	},
	{
		conditionId:        common.HexToHash("0x0004866a1cd8e94ac08e4ff562f038a8d36549cb148f0bb17f64b3e3b96d346d"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "112511989825845734182236324158203291540744630369333629530699046876204345790304",
	},
	{
		conditionId:        common.HexToHash("0x00048d28dda28b7d4834c26f8d250ad83fc993d5bf8250a27091c7b346895a50"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "58279118994270368017979235664956235643252109323245853871333019262731738194762",
	},
	{
		conditionId:        common.HexToHash("0x00048d28dda28b7d4834c26f8d250ad83fc993d5bf8250a27091c7b346895a50"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "63577656802562650859732341604986257493868112308163604485275703581629497058489",
	},
	{
		conditionId:        common.HexToHash("0x0004a8999a7b83148a9bf5e5ec66495d1e4c48e6b2b5a4abf5b7159130ffed76"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "63975626192865853852806303927209510927120760919439120621521948675554703046302",
	},
	{
		conditionId:        common.HexToHash("0x0004a8999a7b83148a9bf5e5ec66495d1e4c48e6b2b5a4abf5b7159130ffed76"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "99025044845069488765369925966520848396411322291347848197076450955124674801327",
	},
	{
		conditionId:        common.HexToHash("0x000693af411390222d3c950f6cdb7b539a9a736fffd714ee1473b93ff5810fd6"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "83334270654887718143913614153031906965899352252635583085904093029532833012895",
	},
	{
		conditionId:        common.HexToHash("0x000693af411390222d3c950f6cdb7b539a9a736fffd714ee1473b93ff5810fd6"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "72766638210312394065719819950882788377704751181947320312338076684864941949583",
	},
	{
		conditionId:        common.HexToHash("0x0007bac99743e5596e1bada6fd3545da8a65282d651040a6464900cb8aed423a"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "100029950543387699762270472155941054442488872584233492721277645390063322369423",
	},
	{
		conditionId:        common.HexToHash("0x0007bac99743e5596e1bada6fd3545da8a65282d651040a6464900cb8aed423a"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "32842970350214142420269990560222863329494355261191015863073628542773906026415",
	},
	{
		conditionId:        common.HexToHash("0x0007e75ef4dd9285b629ab37d40c60fe3f0a254b32d92d603591911904445c0c"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "86343292128005284538439541622778653390773202902200801065662693480138787837772",
	},
	{
		conditionId:        common.HexToHash("0x0007e75ef4dd9285b629ab37d40c60fe3f0a254b32d92d603591911904445c0c"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "111056697726309312129432695354648181874892480405971881106423680381525316430025",
	},
	{
		conditionId:        common.HexToHash("0x00086d72c4a9945a7c5aa100907f1e128b490e4fcb2b89ed00e954845e9757e3"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "48478842431397713783542307285176866680510695828362589487641742360349880697403",
	},
	{
		conditionId:        common.HexToHash("0x00086d72c4a9945a7c5aa100907f1e128b490e4fcb2b89ed00e954845e9757e3"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "17083468761061288801240004184915663332588187403474123332679480229145229604659",
	},
	{
		conditionId:        common.HexToHash("0x0008ddb7d2b82439e50a1d6045ea22187ae0423ac7ecb398af16d06b503e26cf"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "99483559583606847980726940653058587685037824929288162054599300449351319956399",
	},
	{
		conditionId:        common.HexToHash("0x0008ddb7d2b82439e50a1d6045ea22187ae0423ac7ecb398af16d06b503e26cf"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "110938874444136110150947059030552381612447075298449850706918772219817454290532",
	},
	{
		conditionId:        common.HexToHash("0x0009218995a1390d2bbeb4c200f4ab7f63ed545d66e9de8ac9b99d0e86a31488"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "59135531317763704917560233056977378997118894662084449769158659955410160753230",
	},
	{
		conditionId:        common.HexToHash("0x0009218995a1390d2bbeb4c200f4ab7f63ed545d66e9de8ac9b99d0e86a31488"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "44260380418820450919509330660134540605928569754396930837412747480588148896448",
	},
	{
		conditionId:        common.HexToHash("0x000a770743f6e080fb73381268582787a38542fc958d98db22177bd5bfca5923"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "19629885166776131323153550714452541306299721101808036257768303176033626116980",
	},
	{
		conditionId:        common.HexToHash("0x000a770743f6e080fb73381268582787a38542fc958d98db22177bd5bfca5923"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "62037735620973583294501131914921584321655287885346894009836073003294710464449",
	},
	{
		conditionId:        common.HexToHash("0x000ad402ae645de41d4bc23198cd4c84d6e38ab2c007995031c320248e2d55c8"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "50538009090894183534513796234343544850736507943814941810902663398384692629490",
	},
	{
		conditionId:        common.HexToHash("0x000ad402ae645de41d4bc23198cd4c84d6e38ab2c007995031c320248e2d55c8"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "97541415468380520230572445069967882665304478977722815361811428012449107381384",
	},
	{
		conditionId:        common.HexToHash("0x000aea148f4edc5a38eeab8fc0513a53c62f9cccd6a914dbdedd5b8f5ebcd09f"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "34788014837895920236754606175112731352237923879275352749797813182283395941003",
	},
	{
		conditionId:        common.HexToHash("0x000aea148f4edc5a38eeab8fc0513a53c62f9cccd6a914dbdedd5b8f5ebcd09f"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "83284846444114785263128461055176911372937341129597639707811816165677284016300",
	},
	{
		conditionId:        common.HexToHash("0x000d2622bf2bc49ffe1b9b609440017d09a75fa97be607e713b4c0045cfd1916"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "33151191092355688953218326855270416648851482779191712553202315435186399630367",
	},
	{
		conditionId:        common.HexToHash("0x000d2622bf2bc49ffe1b9b609440017d09a75fa97be607e713b4c0045cfd1916"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "88655575081374444143517995769308950927502700167063816868856021956053804446182",
	},
	{
		conditionId:        common.HexToHash("0x000e14d57f237efca90247b1e2e18193113827764a0d85d7364245d8c31f3979"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "2770452408239763521047486923569299631737370751332043111430237430641393831885",
	},
	{
		conditionId:        common.HexToHash("0x000e14d57f237efca90247b1e2e18193113827764a0d85d7364245d8c31f3979"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "108079031350675923472855411057456677190338254279089814193034424267515999515445",
	},
	{
		conditionId:        common.HexToHash("0x0013858672cb2d7dc7dc984f1c69fc8c9e435670ca5503463cc96a75dc0f946e"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "76996576790893615538173654178664990810663291489734666789347868428911845394391",
	},
	{
		conditionId:        common.HexToHash("0x0013858672cb2d7dc7dc984f1c69fc8c9e435670ca5503463cc96a75dc0f946e"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "55422526941827462525921523343395896917547042044614905158786966093073574667706",
	},
	{
		conditionId:        common.HexToHash("0x0013919c78694ad5672b874a9c80cf62b06c63e337baed9c6ca2e8826a3b2c5a"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "64473968592174318477911633885237974719502335364854551297727082341450134561074",
	},
	{
		conditionId:        common.HexToHash("0x0013919c78694ad5672b874a9c80cf62b06c63e337baed9c6ca2e8826a3b2c5a"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "67589392835585961573282863111213771542824012118471065006037776433608794289682",
	},
	{
		conditionId:        common.HexToHash("0x0013bc151aac6b55a5e9c4a566e2d908c13978f727e58aec12c6261527e88e1c"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "80749248541556595414531518370552126548765790825998600068991680371657865991191",
	},
	{
		conditionId:        common.HexToHash("0x0013bc151aac6b55a5e9c4a566e2d908c13978f727e58aec12c6261527e88e1c"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "72294595681961852298118793031109178305252083177848670168058808831097830787934",
	},
	{
		conditionId:        common.HexToHash("0x0014234cd11997bf2480618c43d5a9aaabc630fcf1eeace0957a86b28ac9bb4d"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "48190221594110879372375447173475696833172144893881806458982502549963268632006",
	},
	{
		conditionId:        common.HexToHash("0x0014234cd11997bf2480618c43d5a9aaabc630fcf1eeace0957a86b28ac9bb4d"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "93099205868419732769312762889136355803335732826012636046530078304868433840482",
	},
	{
		conditionId:        common.HexToHash("0x001543c3d454fc848ba12d0dd46df8e70b2011b148fb453584337f7609e2fd42"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "114006653432446068168081713788996744951989823612254193344163949337585747906060",
	},
	{
		conditionId:        common.HexToHash("0x001543c3d454fc848ba12d0dd46df8e70b2011b148fb453584337f7609e2fd42"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "98969975707367651949475542102161857006012308112241070109196339303430841381453",
	},
	{
		conditionId:        common.HexToHash("0x00169ae097007cf17a1e0368fc158d860abdff70d4f37a950f7dfd166abae51d"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "9802709259800486909006980505724589276528515909271599264190969786915067760545",
	},
	{
		conditionId:        common.HexToHash("0x00169ae097007cf17a1e0368fc158d860abdff70d4f37a950f7dfd166abae51d"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "1321623108118877726126977946753002357320609692075891740450064022681216347037",
	},
	{
		conditionId:        common.HexToHash("0x00170a925d97ede2ef8ac1c2dc4bc176ee011482b35653c01ef354469d7c76d2"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "25688525863651896091311688845391984382081303406270737509020234373417509853258",
	},
	{
		conditionId:        common.HexToHash("0x00170a925d97ede2ef8ac1c2dc4bc176ee011482b35653c01ef354469d7c76d2"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "67257856682749610235883335720925553383229948802386573020764891903084112145873",
	},
	{
		conditionId:        common.HexToHash("0x0018c37aaafb3423209c7da9180e1b6057d7db6f99411c8eb3902526f2161647"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "84542533479941776851117457250458123727795022449439110257960725327390137309897",
	},
	{
		conditionId:        common.HexToHash("0x0018c37aaafb3423209c7da9180e1b6057d7db6f99411c8eb3902526f2161647"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "61009578930611726282624963028509075473616387864095701672916226861361343301018",
	},
	{
		conditionId:        common.HexToHash("0x001aca5aaf7b83cda4ab497b09cf8c4cc30611891d4b7b7d7ba1e36c4ec29bcb"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "86451778375441556270690215381475430894730867991715668126141778434310238274670",
	},
	{
		conditionId:        common.HexToHash("0x001aca5aaf7b83cda4ab497b09cf8c4cc30611891d4b7b7d7ba1e36c4ec29bcb"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "25748801505173666762010329044249714682133590893801410970946154160084916617705",
	},

	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "66094911325210800217863651486280594209584449966157268713022829988400850712121",
	},
	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "93240038847443264846920713568005650733949478209913835452672370242435634036071",
	},
	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(4),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "28554598681269789630630461346941051685251459881220291882123400499594815321740",
	},
	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(3),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "10626229078873614401207563845509151617607399614537262588821987535770099165728",
	},
	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(6),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "40611390011953835360974788350986894448230405872316808715550176290245842404681",
	},
	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(5),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "8950064735805333458332827963851465213781579897345424868973334238040916251858",
	},
	{
		conditionId:        common.HexToHash("0x001b6faa35c7d18d7eabea1a599812f7f0e132ea624e793b3921320e5bea9f5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(7),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "18064625499127865019786358757959591864563905020663360913205034136078233062348",
	},
	{
		conditionId:        common.HexToHash("0x001cab05d90516664d1ffb16cbb0fbe557c7e8624f83ecd7277064b4c67ab413"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "78841022279336091385148591876581695360754103600839868235340288485241260884147",
	},
	{
		conditionId:        common.HexToHash("0x001cab05d90516664d1ffb16cbb0fbe557c7e8624f83ecd7277064b4c67ab413"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "6345431971106106712569955750747983606869823689647767457713242759218115955052",
	},
	{
		conditionId:        common.HexToHash("0x001d660d9690357fd9306b7035b2c7cb767318acd99c44e1122fed31290aaa0c"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "90906742398735094209796670358591438888310514256028455376243735294222196895542",
	},
	{
		conditionId:        common.HexToHash("0x001d660d9690357fd9306b7035b2c7cb767318acd99c44e1122fed31290aaa0c"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "39449958771804881639117985107365587885686235763526143169229363905185882387076",
	},
	{
		conditionId:        common.HexToHash("0x001d6eba6da8bad7c61e04ffef911a4578799630779d5075461e806e949bcbaa"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "74405383223403652004560869720382305856307433942032035449217863238194873167776",
	},
	{
		conditionId:        common.HexToHash("0x001d6eba6da8bad7c61e04ffef911a4578799630779d5075461e806e949bcbaa"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "69328036292707661970161512521584145997612878812462546104663955882479681047640",
	},
	{
		conditionId:        common.HexToHash("0x001dcf3b8f4f12bab33ea23363287582f8b9397487bfc8dca247681bce3b64c4"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "29025348385073691541035476870855826301901665939253469649292902474582268345157",
	},
	{
		conditionId:        common.HexToHash("0x001dcf3b8f4f12bab33ea23363287582f8b9397487bfc8dca247681bce3b64c4"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "24152605726101278305093299715026409235351664806941537558431285733277839816593",
	},
	{
		conditionId:        common.HexToHash("0x001dd4fa85859071f1a8006450edd47f834452faf0493831291955d5f45eefca"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "77996364035500375791043832426295084683435456320169058571624808291122807696573",
	},
	{
		conditionId:        common.HexToHash("0x001dd4fa85859071f1a8006450edd47f834452faf0493831291955d5f45eefca"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "45337682013912078007568496950305366244911930564789031287061349432961918231693",
	},
	{
		conditionId:        common.HexToHash("0x002040315aac5da345e7888973d850d795134cc90da490f38091d0ccd63aa0d2"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "73752630629951745214181232882666599928303233175552511988738296962155787614116",
	},
	{
		conditionId:        common.HexToHash("0x002040315aac5da345e7888973d850d795134cc90da490f38091d0ccd63aa0d2"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "3375453420431161957386126668990238092866619981359297940695770530899201243654",
	},
	{
		conditionId:        common.HexToHash("0x0020ae06e91a7cde13521bc68846f57eab20b471de0b21b01df0ed6290f86fd1"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "85462348394726151800790046698707263873680745465537494516466442778245631077986",
	},
	{
		conditionId:        common.HexToHash("0x0020ae06e91a7cde13521bc68846f57eab20b471de0b21b01df0ed6290f86fd1"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "91523706511577339537584403945599850152217281914282233341896655497236114580897",
	},
	{
		conditionId:        common.HexToHash("0x0021bb53d9c023585acdda01466d71e00108b4428694be006b896ab2fa6da08a"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "17250295014852180266426606852813266446819333537067584437713410642544241092588",
	},
	{
		conditionId:        common.HexToHash("0x0021bb53d9c023585acdda01466d71e00108b4428694be006b896ab2fa6da08a"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "33908473583097845547578369755252172764179940812456255098736289329613057753285",
	},
	{
		conditionId:        common.HexToHash("0x0024a17abbd10a57e6edf22cc175a3233c87f8c146424e6f347c3df5c4f313fb"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "8787355922532991421051938667071269850644645908116152927134946403819890710316",
	},
	{
		conditionId:        common.HexToHash("0x0024a17abbd10a57e6edf22cc175a3233c87f8c146424e6f347c3df5c4f313fb"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "31505503626295801369180788683836689798047780483619851919861574002734134691751",
	},
	{
		conditionId:        common.HexToHash("0x00262c7134bc0a149f8e1ec5ccd7b6480f035fd3ea2172377fa5c59867f7d5c0"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "110573615514198976468769723655200848000526107554892935793699238573252277717705",
	},
	{
		conditionId:        common.HexToHash("0x00262c7134bc0a149f8e1ec5ccd7b6480f035fd3ea2172377fa5c59867f7d5c0"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "48127515854966788601303879237416328794574941484858816073070171753240426994147",
	},
	{
		conditionId:        common.HexToHash("0x002745fef5825441fcf577cf22ce324b84e90fa7e4c7f7cb0fc74b73227912a5"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "22509072565654408091552118386753459562477284526491151173929011308575293411125",
	},
	{
		conditionId:        common.HexToHash("0x002745fef5825441fcf577cf22ce324b84e90fa7e4c7f7cb0fc74b73227912a5"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "35392191556880851387756105042587844322122651630785665747772255898928648170173",
	},
	{
		conditionId:        common.HexToHash("0x00289840c4263f68f5aadacb76922fce1665a2cce50d12d46c062e2da7778342"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "38527763944728826812491034054114067429135306060048733973351126417956219618854",
	},
	{
		conditionId:        common.HexToHash("0x00289840c4263f68f5aadacb76922fce1665a2cce50d12d46c062e2da7778342"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "89231678400856617435494573110006335130601284721945389636616002153230174824822",
	},
	{
		conditionId:        common.HexToHash("0x002b8635dd898a562fbf907aa81f6da78bfb586a4a6dbd08547df6cea7045500"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "97822749907046493262594195544625755958016893920008267934286408769106410807087",
	},
	{
		conditionId:        common.HexToHash("0x002b8635dd898a562fbf907aa81f6da78bfb586a4a6dbd08547df6cea7045500"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    defaultCollateralToken,

		expectedPositionId: "63431915373392643845104091954921439784447207875388426362613136644550795024251",
	},
	{
		conditionId:        common.HexToHash("0x002cbae411db59499a161c7746627c4bfd68632e9225f1d54edefc586379ce5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(1),
		collateralToken:    wrappedCollateralToken,

		expectedPositionId: "53246765615777256938923521745854953852658210782177348469325115884735501014069",
	},
	{
		conditionId:        common.HexToHash("0x002cbae411db59499a161c7746627c4bfd68632e9225f1d54edefc586379ce5b"),
		parentCollectionId: defaultParentCollectionId,
		indexSet:           big.NewInt(2),
		collateralToken:    wrappedCollateralToken,

		expectedPositionId: "25007517135845101427908720667301529150815968019672653679159966700099365720712",
	},
}
