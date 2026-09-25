package iterators

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigRange(t *testing.T) {
	start := big.NewInt(0)
	end := big.NewInt(10)

	cnt := 0
	for i, num := range BigRange(start, end) {
		require.Equal(t, cnt, i)
		require.Equal(t, int64(cnt), num.Int64())
		cnt++
	}
	require.Equal(t, 11, cnt)
}

func TestBigRangeSubranges(t *testing.T) {
	tcs := []struct {
		start             *big.Int
		end               *big.Int
		subrangeSize      int
		expectedSubranges []BigSubrange
	}{
		{start: big.NewInt(0), end: big.NewInt(29), subrangeSize: 10, expectedSubranges: []BigSubrange{
			{Start: big.NewInt(0), End: big.NewInt(9)},
			{Start: big.NewInt(10), End: big.NewInt(19)},
			{Start: big.NewInt(20), End: big.NewInt(29)},
		}},
		{start: big.NewInt(0), end: big.NewInt(30), subrangeSize: 10, expectedSubranges: []BigSubrange{
			{Start: big.NewInt(0), End: big.NewInt(9)},
			{Start: big.NewInt(10), End: big.NewInt(19)},
			{Start: big.NewInt(20), End: big.NewInt(29)},
			{Start: big.NewInt(30), End: big.NewInt(30)},
		}},
		{start: big.NewInt(0), end: big.NewInt(3), subrangeSize: 1, expectedSubranges: []BigSubrange{
			{Start: big.NewInt(0), End: big.NewInt(0)},
			{Start: big.NewInt(1), End: big.NewInt(1)},
			{Start: big.NewInt(2), End: big.NewInt(2)},
			{Start: big.NewInt(3), End: big.NewInt(3)},
		}},
		{start: big.NewInt(0), end: big.NewInt(0), subrangeSize: 0, expectedSubranges: []BigSubrange{}},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("start: %s, end: %s, subrangeSize: %d", tc.start.String(), tc.end.String(), tc.subrangeSize), func(t *testing.T) {
			cnt := 0
			for i, subrange := range BigRangeSubranges(tc.start, tc.end, tc.subrangeSize) {
				require.Equal(t, i, cnt)
				require.Equal(t, tc.expectedSubranges[i].Start.Int64(), subrange.Start.Int64())
				require.Equal(t, tc.expectedSubranges[i].End.Int64(), subrange.End.Int64())
				cnt++
			}
			require.Equal(t, len(tc.expectedSubranges), cnt)
		})
	}
}

func TestBigRangeBatch(t *testing.T) {
	tcs := []struct {
		start           *big.Int
		end             *big.Int
		batchSize       int
		expectedBatches [][]*big.Int
	}{
		{start: big.NewInt(0), end: big.NewInt(5), batchSize: 2, expectedBatches: [][]*big.Int{
			{big.NewInt(0), big.NewInt(1)},
			{big.NewInt(2), big.NewInt(3)},
			{big.NewInt(4), big.NewInt(5)},
		}},
		{start: big.NewInt(0), end: big.NewInt(5), batchSize: 3, expectedBatches: [][]*big.Int{
			{big.NewInt(0), big.NewInt(1), big.NewInt(2)},
			{big.NewInt(3), big.NewInt(4), big.NewInt(5)},
		}},
		{start: big.NewInt(0), end: big.NewInt(5), batchSize: 4, expectedBatches: [][]*big.Int{
			{big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3)},
			{big.NewInt(4), big.NewInt(5)},
		}},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("start: %s, end: %s, batchSize: %d", tc.start.String(), tc.end.String(), tc.batchSize), func(t *testing.T) {
			cnt := 0
			for i, batch := range BigRangeBatch(tc.start, tc.end, tc.batchSize) {
				require.Equal(t, i, cnt)
				require.Equal(t, len(tc.expectedBatches[i]), len(batch))
				for j, num := range batch {
					require.Equal(t, tc.expectedBatches[i][j].Int64(), num.Int64())
				}
				cnt++
			}
			require.Equal(t, len(tc.expectedBatches), cnt)
		})
	}
}
