package iterators

import (
	"iter"
	"math/big"
)

// BigRange returns a sequence of big.Ints [start, end]
func BigRange(start *big.Int, end *big.Int) iter.Seq2[int, *big.Int] {
	bigOne := big.NewInt(1)

	return func(yield func(int, *big.Int) bool) {
		current := new(big.Int).Set(start)
		i := 0
		for current.Cmp(end) <= 0 {
			if !yield(i, new(big.Int).Set(current)) {
				return
			}
			current.Add(current, bigOne)
			i++
		}
	}
}

// BigRangeBatch returns a sequence of batches of big.Ints [start, end]
func BigRangeBatch(start *big.Int, end *big.Int, batchSize int) iter.Seq2[int, []*big.Int] {
	if batchSize <= 0 {
		return func(yield func(int, []*big.Int) bool) {}
	}
	bigOne := big.NewInt(1)

	return func(yield func(int, []*big.Int) bool) {
		current := new(big.Int).Set(start)
		i := 0
		for current.Cmp(end) <= 0 {
			batch := make([]*big.Int, 0, batchSize)
			for j := 0; j < batchSize && current.Cmp(end) <= 0; j++ {
				batch = append(batch, new(big.Int).Set(current))
				current.Add(current, bigOne)
			}
			if !yield(i, batch) {
				return
			}
			i++
		}
	}
}

type BigSubrange struct {
	Start *big.Int
	End   *big.Int
}

// BigRangeSubranges splits original range into subranges of size subrangeSize and returns a sequence of BigSubranges
// Example: start = 1, end = 30, subrangeSize = 10 will return:
// [{1, 10}, {11, 20}, {21, 30}]
func BigRangeSubranges(start *big.Int, end *big.Int, subrangeSize int) iter.Seq2[int, BigSubrange] {
	if subrangeSize <= 0 {
		return func(yield func(int, BigSubrange) bool) {}
	}

	bigOne := big.NewInt(1)

	return func(yield func(int, BigSubrange) bool) {
		current := new(big.Int).Set(start)
		i := 0
		for current.Cmp(end) <= 0 {
			subrangeStart := new(big.Int).Set(current)
			subrangeEnd := new(big.Int).Add(subrangeStart, big.NewInt(int64(subrangeSize-1)))
			if subrangeEnd.Cmp(end) > 0 {
				subrangeEnd.Set(end)
			}
			if !yield(i, BigSubrange{
				Start: new(big.Int).Set(subrangeStart),
				End:   new(big.Int).Set(subrangeEnd),
			}) {
				return
			}
			i++
			current.Add(subrangeEnd, bigOne)
		}
	}
}
