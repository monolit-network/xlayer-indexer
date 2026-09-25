package evmclient

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/monolit-network/xlayer-indexer/evm/pkg/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type rpcBlock struct {
	Hash         *common.Hash        `json:"hash"`
	Transactions []rpcTransaction    `json:"transactions"`
	UncleHashes  []common.Hash       `json:"uncles"`
	Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

type rpcReceipt struct {
	types.Receipt
	From common.Address  `json:"from"`
	To   *common.Address `json:"to"`
}

func isZeroPlaceholderTransaction(tx rpcTransaction) bool {
	if tx.tx == nil {
		return false
	}

	zeroAddr := common.Address{}
	if tx.From != nil && *tx.From != zeroAddr {
		return false
	}
	if tx.tx.To() != nil && *tx.tx.To() != zeroAddr {
		return false
	}
	if tx.tx.Nonce() != 0 || tx.tx.Gas() != 0 {
		return false
	}
	if tx.tx.Value() != nil && tx.tx.Value().Sign() != 0 {
		return false
	}
	if len(tx.tx.Data()) != 0 {
		return false
	}

	v, r, s := tx.tx.RawSignatureValues()
	return (v == nil || v.Sign() == 0) &&
		(r == nil || r.Sign() == 0) &&
		(s == nil || s.Sign() == 0)
}

func filterZeroPlaceholderTransactions(txs []rpcTransaction) []rpcTransaction {
	if len(txs) == 0 {
		return txs
	}

	filtered := make([]rpcTransaction, 0, len(txs))
	for _, tx := range txs {
		if isZeroPlaceholderTransaction(tx) {
			continue
		}
		filtered = append(filtered, tx)
	}
	return filtered
}

func isZeroPlaceholderReceipt(receipt rpcReceipt) bool {
	zeroAddr := common.Address{}
	if receipt.From != zeroAddr {
		return false
	}
	if receipt.To != nil && *receipt.To != zeroAddr {
		return false
	}
	if receipt.GasUsed != 0 {
		return false
	}
	if receipt.EffectiveGasPrice != nil && receipt.EffectiveGasPrice.Sign() != 0 {
		return false
	}
	return true
}

func filterZeroPlaceholderReceipts(receipts []rpcReceipt) BlockReceipts {
	if len(receipts) == 0 {
		return BlockReceipts{}
	}

	filtered := make(BlockReceipts, len(receipts))
	for _, receipt := range receipts {
		if isZeroPlaceholderReceipt(receipt) {
			continue
		}
		r := receipt.Receipt
		filtered[r.TxHash] = &r
	}
	return filtered
}

func isZeroPlaceholderTrace(trace *models.FullTraceResult) bool {
	if trace == nil {
		return true
	}
	if traceHasLogs(trace) {
		return false
	}
	zeroAddr := common.Address{}.Hex()
	if trace.From != zeroAddr && trace.From != "" {
		return false
	}
	if trace.To != zeroAddr && trace.To != "" {
		return false
	}
	return true
}

func traceHasLogs(trace *models.FullTraceResult) bool {
	if trace == nil {
		return false
	}
	if len(trace.Logs) > 0 {
		return true
	}
	for _, call := range trace.Calls {
		if traceHasLogs(call) {
			return true
		}
	}
	return false
}

func filterZeroPlaceholderTraces(traces BlockTraces) BlockTraces {
	if len(traces) == 0 {
		return BlockTraces{}
	}
	filtered := make(BlockTraces, len(traces))
	for txHash, trace := range traces {
		if isZeroPlaceholderTrace(trace) {
			continue
		}
		filtered[txHash] = trace
	}
	return filtered
}

// copy from go-ethereum/ethclient.getBlock() method
// does not support uncles
func blockFromRpcJson(raw json.RawMessage, withTransactions bool) (*types.Block, error) {
	// Decode header and transactions.
	var head *types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	// When the block is not found, the API returns JSON null.
	if head == nil {
		return nil, ethereum.NotFound
	}

	if !withTransactions {
		return types.NewBlockWithHeader(head), nil
	}

	var body rpcBlock
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Pending blocks don't return a block hash, compute it for sender caching.
	if body.Hash == nil {
		tmp := head.Hash()
		body.Hash = &tmp
	}
	body.Transactions = filterZeroPlaceholderTransactions(body.Transactions)

	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, errors.New("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == types.EmptyTxsHash && len(body.Transactions) > 0 {
		return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != types.EmptyTxsHash && len(body.Transactions) == 0 {
		return nil, errors.New("server returned empty transaction list but block header indicates transactions")
	}
	// // Load uncles because they are not included in the block response.
	// var uncles []*types.Header
	// if len(body.UncleHashes) > 0 {
	// 	uncles = make([]*types.Header, len(body.UncleHashes))
	// 	reqs := make([]rpc.BatchElem, len(body.UncleHashes))
	// 	for i := range reqs {
	// 		reqs[i] = rpc.BatchElem{
	// 			Method: "eth_getUncleByBlockHashAndIndex",
	// 			Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
	// 			Result: &uncles[i],
	// 		}
	// 	}
	// 	if err := ec.c.BatchCallContext(ctx, reqs); err != nil {
	// 		return nil, err
	// 	}
	// 	for i := range reqs {
	// 		if reqs[i].Error != nil {
	// 			return nil, reqs[i].Error
	// 		}
	// 		if uncles[i] == nil {
	// 			return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
	// 		}
	// 	}
	// }

	// Fill the sender cache of transactions in the block.
	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if tx.From != nil {
			setSenderFromServer(tx.tx, *tx.From, *body.Hash)
		}
		txs[i] = tx.tx
	}

	return types.NewBlockWithHeader(head).WithBody(
		types.Body{
			Transactions: txs,
			Withdrawals:  body.Withdrawals,
		}), nil
}

func setSenderFromServer(tx *types.Transaction, addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender(&senderFromServer{addr, block}, tx)
}

// senderFromServer is a types.Signer that remembers the sender address returned by the RPC
// server. It is stored in the transaction's sender address cache to avoid an additional
// request in TransactionSender.
type senderFromServer struct {
	addr      common.Address
	blockhash common.Hash
}

var errNotCached = errors.New("sender not cached")

func (s *senderFromServer) Equal(other types.Signer) bool {
	os, ok := other.(*senderFromServer)
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer) Sender(tx *types.Transaction) (common.Address, error) {
	if s.addr == (common.Address{}) {
		return common.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer) ChainID() *big.Int {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) Hash(tx *types.Transaction) common.Hash {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) SignatureValues(tx *types.Transaction, sig []byte) (R, S, V *big.Int, err error) {
	panic("can't sign with senderFromServer")
}
