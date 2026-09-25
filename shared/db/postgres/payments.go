package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/jackc/pgx/v5"
)

func (q *PostgresDB) CreatePaymentWalletsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateTablePaymentWalletsSQL)
	return err
}

func (q *PostgresDB) CreatePaymentChainInfoTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateTablePaymentChainInfoSQL)
	return err
}

func scanWallet(row pgx.Row, cipher *PostgresDB) (*models.Wallet, error) {
	privateKey := []byte{}
	var wallet models.Wallet
	if err := row.Scan(&wallet.Address, &privateKey, &wallet.ReservedOrderID, &wallet.ReservedAt, &wallet.CooldownUntil); err != nil {
		return nil, err
	}
	privateKeyBytes, err := cipher.decryptBackendBlob(privateKey)
	if err != nil {
		return nil, err
	}
	wallet.PrivateKey = privateKeyBytes
	return &wallet, nil
}

func scanWalletRows(rows pgx.Rows, cipher *PostgresDB) ([]models.Wallet, error) {
	defer rows.Close()

	wallets := make([]models.Wallet, 0)
	for rows.Next() {
		privateKey := []byte{}
		var wallet models.Wallet
		if err := rows.Scan(&wallet.Address, &privateKey, &wallet.ReservedOrderID, &wallet.ReservedAt, &wallet.CooldownUntil); err != nil {
			return nil, err
		}
		privateKeyBytes, err := cipher.decryptBackendBlob(privateKey)
		if err != nil {
			return nil, err
		}
		wallet.PrivateKey = privateKeyBytes
		wallets = append(wallets, wallet)
	}

	return wallets, rows.Err()
}

func randomPaymentWalletAddressCursor() (string, error) {
	var bytes [20]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("generate payment wallet cursor: %w", err)
	}
	return "0x" + hex.EncodeToString(bytes[:]), nil
}

func (q *PostgresDB) GetAvailablePaymentWalletForUpdate(ctx context.Context) (*models.Wallet, error) {
	cursor, err := randomPaymentWalletAddressCursor()
	if err != nil {
		return nil, err
	}

	wallet, err := scanWallet(q.querier().QueryRow(ctx, queries.SelectAvailablePaymentWalletFromCursorForUpdateSQL, cursor), q)
	if err == nil {
		return wallet, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}

	wallet, err = scanWallet(q.querier().QueryRow(ctx, queries.SelectAvailablePaymentWalletBeforeCursorForUpdateSQL, cursor), q)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return wallet, nil
}

func (q *PostgresDB) GetPaymentWallets(ctx context.Context) ([]models.Wallet, error) {
	rows, err := q.querier().Query(ctx, queries.SelectPaymentWalletsSQL)
	if err != nil {
		return nil, err
	}
	return scanWalletRows(rows, q)
}

func (q *PostgresDB) InsertPaymentWallet(ctx context.Context, wallet models.Wallet) error {
	encryptedPrivateKey, err := q.encryptBackendBlob(wallet.PrivateKey)
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx, queries.InsertPaymentWalletSQL, wallet.Address, encryptedPrivateKey, wallet.ReservedOrderID, wallet.ReservedAt, wallet.CooldownUntil)
	return err
}

func (q *PostgresDB) ReservePaymentWallet(ctx context.Context, address string, orderID string, reservedAt time.Time) error {
	_, err := q.querier().Exec(ctx, queries.ReservePaymentWalletSQL, address, orderID, reservedAt)
	return err
}

func (q *PostgresDB) ReleasePaymentWalletReservation(ctx context.Context, address string) error {
	_, err := q.querier().Exec(ctx, queries.ReleasePaymentWalletReservationSQL, address)
	return err
}

func (q *PostgresDB) GetPaymentChainInfos(ctx context.Context) ([]models.ChainInfo, error) {
	rows, err := q.querier().Query(ctx, queries.SelectPaymentChainInfosSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chainInfos := []models.ChainInfo{}
	for rows.Next() {
		var info models.ChainInfo
		var tokenAddresses []string
		var decimals []int32
		var tokenNames []string
		var tokenImageURLs []string
		err := rows.Scan(&info.ChainName, &tokenAddresses, &info.LatestBlockOffset, &decimals, &tokenNames, &tokenImageURLs, &info.URL)
		if err != nil {
			return nil, err
		}
		info.Tokens, err = buildPaymentTokens(tokenAddresses, decimals, tokenNames, tokenImageURLs)
		if err != nil {
			return nil, fmt.Errorf("build payment tokens for chain %s: %w", info.ChainName, err)
		}
		chainInfos = append(chainInfos, info)
	}
	return chainInfos, nil
}

func buildPaymentTokens(tokenAddresses []string, decimals []int32, tokenNames []string, tokenImageURLs []string) ([]models.PaymentToken, error) {
	if len(tokenAddresses) != len(decimals) {
		return nil, fmt.Errorf("token addresses and decimals length mismatch: %d != %d", len(tokenAddresses), len(decimals))
	}
	if len(tokenAddresses) != len(tokenNames) {
		return nil, fmt.Errorf("token addresses and token names length mismatch: %d != %d", len(tokenAddresses), len(tokenNames))
	}
	if len(tokenAddresses) != len(tokenImageURLs) {
		return nil, fmt.Errorf("token addresses and token image urls length mismatch: %d != %d", len(tokenAddresses), len(tokenImageURLs))
	}

	tokens := make([]models.PaymentToken, len(tokenAddresses))
	for i := range tokenAddresses {
		tokens[i] = models.PaymentToken{
			Address:  strings.ToLower(tokenAddresses[i]),
			Decimal:  decimals[i],
			Name:     tokenNames[i],
			ImageURL: tokenImageURLs[i],
		}
	}

	return tokens, nil
}
