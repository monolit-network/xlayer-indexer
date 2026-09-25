package polymarket

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bn256"
)

var (
	// BN254 curve parameters
	P, _ = new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	B    = big.NewInt(3)

	// (P + 1) / 4 for sqrt calculation (works because P ≡ 3 mod 4)
	sqrtExp *big.Int

	one    = big.NewInt(1)
	two    = big.NewInt(2)
	bit254 *big.Int
)

func init() {
	sqrtExp = new(big.Int).Add(P, one)
	sqrtExp.Div(sqrtExp, big.NewInt(4))

	bit254 = new(big.Int).Lsh(one, 254)
}

// sqrt computes modular square root: x^((P+1)/4) mod P
func sqrt(x *big.Int) *big.Int {
	return new(big.Int).Exp(x, sqrtExp, P)
}

// mulmod computes (a * b) mod P
func mulmod(a, b *big.Int) *big.Int {
	result := new(big.Int).Mul(a, b)
	return result.Mod(result, P)
}

// addmod computes (a + b) mod P
func addmod(a, b *big.Int) *big.Int {
	result := new(big.Int).Add(a, b)
	return result.Mod(result, P)
}

func getCollectionId(parentCollectionId common.Hash, conditionId common.Hash, indexSet *big.Int) common.Hash {
	// x1 = uint(keccak256(abi.encodePacked(conditionId, indexSet)))
	indexSetBytes := common.LeftPadBytes(indexSet.Bytes(), 32)
	hash := crypto.Keccak256(append(conditionId.Bytes(), indexSetBytes...))
	x1 := new(big.Int).SetBytes(hash)

	// odd = x1 >> 255 != 0
	odd := new(big.Int).Rsh(x1, 255).Sign() != 0

	var y1, yy *big.Int
	for {
		x1 = addmod(x1, one)
		// yy = x1³ + B mod P
		x1Cubed := mulmod(x1, mulmod(x1, x1))
		yy = addmod(x1Cubed, B)
		y1 = sqrt(yy)

		if mulmod(y1, y1).Cmp(yy) == 0 {
			break
		}
	}

	// Adjust y1 sign based on odd flag
	y1IsEven := new(big.Int).Mod(y1, two).Sign() == 0
	if (odd && y1IsEven) || (!odd && !y1IsEven) {
		y1 = new(big.Int).Sub(P, y1)
	}

	// Handle parent collection
	x2 := new(big.Int).SetBytes(parentCollectionId.Bytes())
	if x2.Sign() != 0 {
		// odd = x2 >> 254 != 0
		odd = new(big.Int).Rsh(x2, 254).Sign() != 0

		// x2 = (x2 << 2) >> 2 - clear top 2 bits
		mask := new(big.Int).Sub(bit254, one)
		x2.And(x2, mask)

		// yy = x2³ + B mod P
		x2Cubed := mulmod(x2, mulmod(x2, x2))
		yy = addmod(x2Cubed, B)
		y2 := sqrt(yy)

		y2IsEven := new(big.Int).Mod(y2, two).Sign() == 0
		if (odd && y2IsEven) || (!odd && !y2IsEven) {
			y2 = new(big.Int).Sub(P, y2)
		}

		if mulmod(y2, y2).Cmp(yy) != 0 {
			panic("invalid parent collection ID")
		}

		// EC point addition
		x1, y1 = ecAdd(x1, y1, x2, y2)
	}

	// Set high bit based on y1 parity
	if new(big.Int).Mod(y1, two).Cmp(one) == 0 {
		x1.Xor(x1, bit254)
	}

	return common.BytesToHash(common.LeftPadBytes(x1.Bytes(), 32))
}

// ecAdd performs elliptic curve point addition on BN254 G1
func ecAdd(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	p1 := new(bn256.G1)
	p2 := new(bn256.G1)

	// Encode points as 64-byte slices (32 bytes X + 32 bytes Y)
	p1Bytes := make([]byte, 64)
	copy(p1Bytes[32-len(x1.Bytes()):32], x1.Bytes())
	copy(p1Bytes[64-len(y1.Bytes()):64], y1.Bytes())

	p2Bytes := make([]byte, 64)
	copy(p2Bytes[32-len(x2.Bytes()):32], x2.Bytes())
	copy(p2Bytes[64-len(y2.Bytes()):64], y2.Bytes())

	if _, err := p1.Unmarshal(p1Bytes); err != nil {
		panic("invalid point 1: " + err.Error())
	}
	if _, err := p2.Unmarshal(p2Bytes); err != nil {
		panic("invalid point 2: " + err.Error())
	}

	// Add points
	result := new(bn256.G1).Add(p1, p2)
	resultBytes := result.Marshal()

	resX := new(big.Int).SetBytes(resultBytes[:32])
	resY := new(big.Int).SetBytes(resultBytes[32:64])

	return resX, resY
}

func getTokenIdPartition(conditionId common.Hash, partition *big.Int, parentCollectionId common.Hash, collateralToken common.Address) *big.Int {
	collectionId := getCollectionId(parentCollectionId, conditionId, partition)
	data := append(collateralToken.Bytes(), collectionId.Bytes()...)
	hash := crypto.Keccak256(data)
	return new(big.Int).SetBytes(hash)
}

func GetTokenId(conditionId common.Hash, index int64, parentCollectionId common.Hash, collateralToken common.Address) *big.Int {
	partition := new(big.Int).Lsh(big.NewInt(1), uint(index))
	return getTokenIdPartition(conditionId, partition, parentCollectionId, collateralToken)
}

func getQuestionId(marketId common.Hash, index uint8) common.Hash {
	var questionId common.Hash
	copy(questionId[:], marketId[:])
	questionId[31] |= index
	return questionId
}

func getConditionId(oracle common.Address, questionId common.Hash) common.Hash {
	data := make([]byte, 0, 84)
	data = append(data, oracle.Bytes()...)
	data = append(data, questionId.Bytes()...)
	data = append(data, common.LeftPadBytes(big.NewInt(2).Bytes(), 32)...)
	return crypto.Keccak256Hash(data)
}
