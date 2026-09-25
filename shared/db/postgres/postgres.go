package postgres

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/monolit-network/xlayer-indexer/shared/db"
	"github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries"
	"github.com/monolit-network/xlayer-indexer/shared/models"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
)

type PostgresClient struct {
	pool          *pgxpool.Pool
	backendCipher cipher.AEAD
}

func NewPostgresClient(connectionString string, maxConns int) (*PostgresClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.MaxConns = int32(maxConns)

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	backendCipher, err := initBackendCipher()
	if err != nil {
		return nil, err
	}

	return &PostgresClient{pool: pool, backendCipher: backendCipher}, nil
}

func NewPostgresClientFromEnv() (*PostgresClient, error) {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	maxConnsStr := os.Getenv("DATABASE_MAX_CONNS")
	var maxConns int
	if maxConnsStr == "" {
		maxConns = 20
	} else {
		var err error
		maxConns, err = strconv.Atoi(maxConnsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DATABASE_MAX_CONNS: %w", err)
		}
	}
	client, err := NewPostgresClient(connectionString, maxConns)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}
	return client, nil
}

func (c *PostgresClient) Querier() db.DB {
	return &PostgresDB{pool: c.pool, backendCipher: c.backendCipher}
}

func (c *PostgresClient) Transaction(ctx context.Context, fn func(ctx context.Context, tx db.DB) error) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
	}()

	querier := &PostgresDB{tx: tx, backendCipher: c.backendCipher}
	if err := fn(ctx, querier); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rollbackErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *PostgresClient) WithRetry(ctx context.Context, maxRetries int, fn func(ctx context.Context, q db.DB) error) error {
	return c.withRetry(ctx, maxRetries, func(ctx context.Context) error {
		return fn(ctx, c.Querier())
	})
}

func (c *PostgresClient) TransactionWithRetry(ctx context.Context, maxRetries int, fn func(ctx context.Context, tx db.DB) error) error {
	return c.withRetry(ctx, maxRetries, func(ctx context.Context) error {
		return c.Transaction(ctx, fn)
	})
}

func (c *PostgresClient) withRetry(ctx context.Context, maxRetries int, fn func(ctx context.Context) error) error {
	delay := 200 * time.Millisecond
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(ctx); err != nil {
			if !isRetryablePostgresError(err) {
				return err
			}
			if maxRetries >= 0 && attempt >= maxRetries {
				return err
			}

			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				return ctx.Err()
			}

			if delay < 5*time.Second {
				delay *= 2
				if delay > 5*time.Second {
					delay = 5 * time.Second
				}
			}
			continue
		}

		return nil
	}
}

func isRetryablePostgresError(err error) bool {
	if err == nil || errors.Is(err, pgx.ErrNoRows) {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40001", // serialization_failure
			"40P01", // deadlock_detected
			"55P03", // lock_not_available
			"53300", // too_many_connections
			"57P01", // admin_shutdown
			"57P02", // crash_shutdown
			"57P03": // cannot_connect_now
			return true
		default:
			return strings.HasPrefix(pgErr.Code, "08")
		}
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection reset by peer") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "server closed the connection") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "timeout: context deadline exceeded") ||
		strings.Contains(msg, "failed to connect") ||
		strings.Contains(msg, "too many connections") ||
		strings.Contains(msg, "deadlock detected") ||
		strings.Contains(msg, "serialization failure")
}

func (c *PostgresClient) Close() error {
	c.pool.Close()
	return nil
}

func initBackendCipher() (cipher.AEAD, error) {
	rawKey := os.Getenv("BACKEND_AES_KEY")
	if rawKey == "" {
		return nil, fmt.Errorf("BACKEND_AES_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BACKEND_AES_KEY: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("BACKEND_AES_KEY must decode to 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AES-GCM: %w", err)
	}
	return gcm, nil
}

func generateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

type pgxQuerier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults
}

type PostgresDB struct {
	pool          *pgxpool.Pool
	tx            pgx.Tx
	backendCipher cipher.AEAD
}

func (q *PostgresDB) querier() pgxQuerier {
	if q.tx != nil {
		return q.tx
	}
	return q.pool
}

func (q *PostgresDB) CreateSocialPostsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateSocialPostsTableSQL)
	return err
}

func (q *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return q.querier().Query(ctx, query, args...)
}

func (q *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return q.querier().QueryRow(ctx, query, args...)
}

func (q *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return q.querier().Exec(ctx, query, args...)
}

func (q *PostgresDB) getNextAssetID(ctx context.Context) (models.VerifiedTokenAssetID, error) {
	var assetID models.VerifiedTokenAssetID
	row := q.querier().QueryRow(ctx, queries.SelectNextDropstabInfoAssetIDSQL)
	if err := row.Scan(&assetID); err != nil {
		return 0, err
	}

	return assetID, nil
}

func (q *PostgresDB) UpsertSocialPost(ctx context.Context, post models.SocialPost) (bool, error) {
	if post.Relations == nil {
		post.Relations = []string{}
	}
	if post.LinksMentioned == nil {
		post.LinksMentioned = []string{}
	}
	if post.TokensMentioned == nil {
		post.TokensMentioned = []string{}
	}

	var inserted bool
	err := q.querier().QueryRow(ctx, queries.UpsertSocialPostSQL,
		post.Link, post.Text, post.Date, string(post.Source), post.Author,
		post.Relations, post.LinksMentioned, post.TokensMentioned).Scan(&inserted)
	return inserted, err
}

func (q *PostgresDB) GetSocialPostByLink(ctx context.Context, link string) (*models.SocialPost, error) {
	var post models.SocialPost
	err := q.querier().QueryRow(ctx, queries.SelectSocialPostByLinkSQL, link).Scan(
		&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
		&post.Relations, &post.LinksMentioned, &post.TokensMentioned)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (q *PostgresDB) GetSocialPostsByLink(ctx context.Context, link string) ([]models.SocialPost, error) {
	rows, err := q.querier().Query(ctx, queries.SelectSocialPostByLinkSQL, link)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.SocialPost, 0, 8)
	for rows.Next() {
		var post models.SocialPost
		if err := rows.Scan(&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
			&post.Relations, &post.LinksMentioned, &post.TokensMentioned); err != nil {
			return nil, err
		}
		res = append(res, post)
	}
	return res, rows.Err()
}

func (q *PostgresDB) GetSocialPostsByLinkBeginning(ctx context.Context, beginning string) ([]models.SocialPost, error) {
	rows, err := q.querier().Query(ctx, queries.SelectSocialPostsByLinkBeginningSQL, beginning)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.SocialPost, 0, 8)
	for rows.Next() {
		var post models.SocialPost
		if err := rows.Scan(&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
			&post.Relations, &post.LinksMentioned, &post.TokensMentioned); err != nil {
			return nil, err
		}
		res = append(res, post)
	}
	return res, rows.Err()
}

func (q *PostgresDB) GetPostsReferencingLink(ctx context.Context, link string, limit int) ([]models.SocialPost, error) {
	rows, err := q.querier().Query(ctx, queries.SelectPostsReferencingLinkSQL, link, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.SocialPost, 0, limit)
	for rows.Next() {
		var post models.SocialPost
		if err := rows.Scan(&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
			&post.Relations, &post.LinksMentioned, &post.TokensMentioned); err != nil {
			return nil, err
		}
		res = append(res, post)
	}
	return res, rows.Err()
}

func (q *PostgresDB) CreateSocialSourcesTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, queries.CreateSocialSourcesTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (q *PostgresDB) UpsertSocialSource(ctx context.Context, source models.SocialSource) error {
	_, err := q.querier().Exec(ctx, queries.UpsertSocialSourceSQL, source.Value, source.Source)
	if err != nil {
		return err
	}

	return nil
}

func (q *PostgresDB) GetSocialSources(ctx context.Context, tag models.SocialSourceTag) ([]string, error) {
	rows, err := q.querier().Query(ctx, queries.SelectSocialSourcesSQL, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func (q *PostgresDB) GetOldestPostForAuthor(ctx context.Context, source models.SocialSourceTag, author string) (*models.SocialPost, error) {
	rows, err := q.querier().Query(ctx, queries.SelectOldestPostSQL, source, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var post models.SocialPost
		if err := rows.Scan(&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
			&post.Relations, &post.LinksMentioned, &post.TokensMentioned); err != nil {
			return nil, err
		}
		return &post, nil
	}
	return nil, nil
}

func (q *PostgresDB) GetNewestPostForAuthor(ctx context.Context, source models.SocialSourceTag, author string) (*models.SocialPost, error) {
	rows, err := q.querier().Query(ctx, queries.SelectNewestPostSQL, source, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var post models.SocialPost
		if err := rows.Scan(&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
			&post.Relations, &post.LinksMentioned, &post.TokensMentioned); err != nil {
			return nil, err
		}
		return &post, nil
	}
	return nil, nil
}

func (q *PostgresDB) GetSocialPostsByAuthor(ctx context.Context, author string) ([]models.SocialPost, error) {
	rows, err := q.querier().Query(ctx, queries.SelectSocialPostsByAuthorSQL, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.SocialPost, 0, 8)
	for rows.Next() {
		var post models.SocialPost
		if err := rows.Scan(&post.Link, &post.Text, &post.Date, &post.Source, &post.Author,
			&post.Relations, &post.LinksMentioned, &post.TokensMentioned); err != nil {
			return nil, err
		}
		res = append(res, post)
	}
	return res, rows.Err()
}

func (q *PostgresDB) UpsertSourceFeed(ctx context.Context, source string, feed string) error {
	_, err := q.querier().Exec(ctx, queries.UpsertSourceFeedSQL, source, feed)
	if err != nil {
		return err
	}

	return nil
}

func (q *PostgresDB) SelectArticles(ctx context.Context, keywords []string, source string, timeFrom time.Time, timeTo time.Time) ([]models.Article, error) {
	if len(keywords) == 0 && source == "" {
		return nil, fmt.Errorf("must provide at least one keyword or source")
	}

	args := make([]interface{}, 0)
	conditions := make([]string, 0)

	if source != "" {
		conditions = append(conditions, "source = $1")
		args = append(args, source)
	}

	if len(keywords) > 0 {
		keywordsQuery := strings.Join(keywords, " | ")
		conditions = append(conditions, fmt.Sprintf("to_tsvector('simple', text) @@ to_tsquery('simple', $%d)", len(args)+1))
		args = append(args, keywordsQuery)
	}

	if !timeFrom.IsZero() && !timeTo.IsZero() {
		conditions = append(conditions, fmt.Sprintf("date BETWEEN $%d AND $%d", len(args)+1, len(args)+2))
		args = append(args, timeFrom, timeTo)
	} else if !timeFrom.IsZero() {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", len(args)+1))
		args = append(args, timeFrom)
	} else if !timeTo.IsZero() {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", len(args)+1))
		args = append(args, timeTo)
	}

	query := queries.SelectArticlesBaseSQL
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY date DESC"

	rows, err := q.querier().Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Link, &article.Title, &article.Text, &article.Authors, &article.Date, &article.Source); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (q *PostgresDB) GetArticlesByIDs(ctx context.Context, ids []int) ([]models.Article, error) {
	if len(ids) == 0 {
		return []models.Article{}, nil
	}

	rows, err := q.querier().Query(ctx, queries.SelectArticlesByIDsSQL, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Link, &article.Title, &article.Text, &article.Authors, &article.Date, &article.Source); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}
