package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureHash(t *testing.T) {
	hash := SelectorFromSignature("swap(uint256,uint256,address,bytes)")
	fmt.Println(hash)
	assert.Equal(t, hash, UniswapV2SwapFuncSelectorHash)

	fmt.Println(UniswapV3SwapFuncSelectorHash)
}

func TestTest(t *testing.T) {
	fmt.Println(PolymarketFPMMFundingRemovedEventSelectorHash)
}

func TestTest2(t *testing.T) {
	fmt.Println(UmaOptimisticOracleV2InvalidPrice)
}
