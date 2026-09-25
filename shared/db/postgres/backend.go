package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/monolit-network/xlayer-indexer/shared/apis/thirdweb"
	backendqueries "github.com/monolit-network/xlayer-indexer/shared/db/postgres/queries/backend"
	"github.com/monolit-network/xlayer-indexer/shared/models"
	"github.com/google/uuid"
)

func (q *PostgresDB) CreateBackendDatabase(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateBackendDatabaseSQL)
	return err
}

func (q *PostgresDB) CreateBackendUsersTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendUsersSQL)
	return err
}

func (q *PostgresDB) CreateBackendReferralRewardsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendReferralRewardsSQL)
	return err
}

func (q *PostgresDB) CreateBackendSessionsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendSessionsSQL)
	return err
}

func (q *PostgresDB) CreateBackendChatHistoryTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendChatHistorySQL)
	return err
}

func (q *PostgresDB) CreateBackendChatFoldersTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendChatFoldersSQL)
	return err
}

func (q *PostgresDB) CreateBackendChatFeedbackTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendChatFeedbackSQL)
	return err
}

func (q *PostgresDB) CreateBackendVizPayloadsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendVizPayloadsSQL)
	return err
}

func (q *PostgresDB) CreateBackendVizDataRefsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendVizDataRefsSQL)
	return err
}

func (q *PostgresDB) CreateAPIEndpointBlacklistTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableAPIEndpointBlacklistSQL)
	return err
}

func (q *PostgresDB) CreateBackendSignalsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendSignalsSQL)
	return err
}

func (q *PostgresDB) CreateBackendSignalEventsTable(ctx context.Context) error {
	_, err := q.querier().Exec(ctx, backendqueries.CreateTableBackendSignalEventsSQL)
	return err
}

func (q *PostgresDB) GetBackendUserIDByWallet(ctx context.Context, wallet string) (string, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectBackendUserIDByWalletSQL, wallet)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (q *PostgresDB) InsertBackendUserWithWallet(ctx context.Context, userID, wallet string) error {
	profiles, err := q.encryptBackendBlob([]byte("[]"))
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx, backendqueries.InsertBackendUserWithWalletSQL, userID, profiles, wallet)
	return err
}

func (q *PostgresDB) GetChatUserByChatID(ctx context.Context, chatID string) (string, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectChatUserByChatIDSQL, chatID)
	var userID string
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get chat user: %w", err)
	}
	return userID, nil
}

// GetMostRecentEmptyChat returns the chat_id of the user's most recent
// reserved-but-untouched chat (empty history, inactive, not deleted).
// Empty string + nil err means the user has no such chat — caller
// should INSERT a fresh row. Used by HandleCreateChat to dedupe
// retry-loops that would otherwise spam "New Chat" rows.
func (q *PostgresDB) GetMostRecentEmptyChat(ctx context.Context, userID string) (string, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectMostRecentEmptyChatSQL, userID)
	var chatID string
	if err := row.Scan(&chatID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get empty chat: %w", err)
	}
	return chatID, nil
}

func (q *PostgresDB) InsertChatFolder(ctx context.Context, folder models.BackendChatFolder) error {
	_, err := q.querier().Exec(ctx,
		backendqueries.InsertChatFolderSQL,
		folder.ID,
		folder.UserID,
		folder.Name,
		folder.CreatedAt,
		folder.UpdatedAt,
	)
	return err
}

func (q *PostgresDB) GetChatFoldersByUserID(ctx context.Context, userID string) ([]models.BackendChatFolder, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectChatFoldersByUserIDSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []models.BackendChatFolder
	for rows.Next() {
		folder, err := scanChatFolder(rows)
		if err != nil {
			return nil, err
		}
		folders = append(folders, *folder)
	}
	return folders, rows.Err()
}

func (q *PostgresDB) GetChatFolderByID(ctx context.Context, folderID string) (*models.BackendChatFolder, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectChatFolderByIDSQL, folderID)
	return scanChatFolder(row)
}

func (q *PostgresDB) UpdateChatFolderName(ctx context.Context, folderID, userID, name string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.UpdateChatFolderNameSQL, folderID, userID, name, time.Now())
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) DeleteChatFolder(ctx context.Context, folderID, userID string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.DeleteChatFolderSQL, folderID, userID, time.Now())
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) ClearChatFolderAssignments(ctx context.Context, folderID string) error {
	_, err := q.querier().Exec(ctx, backendqueries.ClearChatFolderAssignmentsSQL, folderID, time.Now())
	return err
}

func (q *PostgresDB) GetBackendUserByID(ctx context.Context, id string) (*models.BackendUser, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectBackendUserByIDSQL, id)
	return q.scanBackendUser(row)
}

func (q *PostgresDB) GetBackendUserByReferralCode(ctx context.Context, referralCode string) (*models.BackendUser, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectBackendUserByReferralCodeSQL, referralCode)
	return q.scanBackendUser(row)
}

func (q *PostgresDB) UpdateBackendUser(ctx context.Context, user models.BackendUser) error {
	profiles, err := marshalBackendProfiles(user.Profiles)
	if err != nil {
		return err
	}
	encrypted, err := q.encryptBackendBlob(profiles)
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx, backendqueries.UpdateBackendUserSQL, user.ThirdwebUserID, encrypted, user.ID)
	return err
}

func (q *PostgresDB) SetBackendUserReferralCode(ctx context.Context, userID, referralCode string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.SetBackendUserReferralCodeSQL, referralCode, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) SetBackendUserReferrerCode(ctx context.Context, userID, referrerCode string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.SetBackendUserReferrerCodeSQL, referrerCode, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) SetBackendUserReferralWallet(ctx context.Context, userID, referralWallet string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.SetBackendUserReferralWalletSQL, referralWallet, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) InsertBackendReferralRewardForPaidOrder(ctx context.Context, order models.BillingOrder) (bool, error) {
	res, err := q.querier().Exec(
		ctx,
		backendqueries.InsertBackendReferralRewardForPaidOrderSQL,
		uuid.New().String(),
		order.ID,
		order.AmountUSDCents,
		models.RefbackBps,
		models.BackendReferralRewardStatusAvailable,
		time.Now().UTC(),
		order.UserID,
	)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) GetBackendUserReferralInfo(ctx context.Context, userID string) (models.BackendReferralBalance, error) {
	row := q.querier().QueryRow(
		ctx,
		backendqueries.SelectBackendUserReferralInfoSQL,
		userID,
		models.BackendReferralRewardStatusPaidOut,
		models.BackendReferralRewardStatusAvailable,
	)
	var balance models.BackendReferralBalance
	var referralCode sql.NullString
	var referralWallet sql.NullString
	if err := row.Scan(
		&referralCode,
		&referralWallet,
		&balance.ReferralCount,
		&balance.TotalUSDCents,
		&balance.PaidOutUSDCents,
		&balance.AvailableUSDCents,
	); err != nil {
		return models.BackendReferralBalance{}, err
	}
	if referralCode.Valid {
		balance.ReferralCode = &referralCode.String
	}
	if referralWallet.Valid {
		balance.ReferralWallet = referralWallet.String
	}
	return balance, nil
}

func (q *PostgresDB) GetBackendUserReferrals(ctx context.Context, userID string, limit, offset int) ([]models.BackendReferral, int64, error) {
	rows, err := q.querier().Query(
		ctx,
		backendqueries.SelectBackendUserReferralsSQL,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	referrals := make([]models.BackendReferral, 0)
	var total int64
	for rows.Next() {
		var referral models.BackendReferral
		if err := rows.Scan(&referral.WalletPrefix, &referral.ReferralEarningsUSDCents, &total); err != nil {
			return nil, 0, err
		}
		referrals = append(referrals, referral)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return referrals, total, nil
}

func (q *PostgresDB) DeleteBackendUserByID(ctx context.Context, id string) error {
	_, err := q.querier().Exec(ctx, backendqueries.DeleteBackendUserByIDSQL, id)
	return err
}

func (q *PostgresDB) ListActiveAPIEndpointBlacklistRules(ctx context.Context) ([]models.APIEndpointBlacklistRule, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectActiveAPIEndpointBlacklistRulesSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query API endpoint blacklist rules: %w", err)
	}
	defer rows.Close()

	var rules []models.APIEndpointBlacklistRule
	for rows.Next() {
		var rule models.APIEndpointBlacklistRule
		var method sql.NullString
		var reason sql.NullString
		if err := rows.Scan(&rule.ID, &method, &rule.MatchType, &rule.Pattern, &reason, &rule.DisabledUntil); err != nil {
			return nil, fmt.Errorf("failed to scan API endpoint blacklist rule: %w", err)
		}
		if method.Valid {
			rule.Method = method.String
		}
		if reason.Valid {
			rule.Reason = reason.String
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate API endpoint blacklist rules: %w", err)
	}

	return rules, nil
}

func (q *PostgresDB) InsertUserSession(ctx context.Context, session models.BackendUserSession) error {
	_, err := q.querier().Exec(ctx,
		backendqueries.InsertUserSessionSQL,
		session.UserID,
		session.SessionID,
		session.CreatedAt,
		session.ExpiresAt,
	)
	return err
}

func (q *PostgresDB) GetUserSessionByUserIDAndSessionID(ctx context.Context, userID, sessionID string) (*models.BackendUserSession, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectUserSessionByUserIDAndSessionIDSQL, userID, sessionID)
	var session models.BackendUserSession
	if err := row.Scan(&session.UserID, &session.SessionID, &session.CreatedAt, &session.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (q *PostgresDB) GetUserSessionsByUserID(ctx context.Context, userID string) ([]models.BackendUserSession, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectUserSessionByUserIDSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.BackendUserSession
	for rows.Next() {
		var session models.BackendUserSession
		if err := rows.Scan(&session.UserID, &session.SessionID, &session.CreatedAt, &session.ExpiresAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

func (q *PostgresDB) GetUserSessionBySessionID(ctx context.Context, sessionID string) (*models.BackendUserSession, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectUserSessionBySessionIDSQL, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session models.BackendUserSession
		if err := rows.Scan(&session.UserID, &session.SessionID, &session.CreatedAt, &session.ExpiresAt); err != nil {
			return nil, err
		}
		return &session, nil
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (q *PostgresDB) UpdateUserSession(ctx context.Context, session models.BackendUserSession) error {
	_, err := q.querier().Exec(ctx,
		backendqueries.UpdateUserSessionSQL,
		session.ExpiresAt,
		session.UserID,
		session.SessionID,
	)
	return err
}

func (q *PostgresDB) DeleteUserSession(ctx context.Context, userID, sessionID string) error {
	_, err := q.querier().Exec(ctx, backendqueries.DeleteUserSessionByUserIDAndSessionIDSQL, userID, sessionID)
	return err
}

func (q *PostgresDB) InsertChat(ctx context.Context, history models.BackendChat) error {
	payload := normalizeChatHistoryPayload(history.History)
	encrypted, err := q.encryptBackendBlob(payload)
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx,
		backendqueries.InsertChatSQL,
		history.ChatID,
		history.UserID,
		history.ShortName,
		history.FolderID,
		history.CreatedAt,
		history.UpdatedAt,
		encrypted,
		history.Active,
	)
	return err
}

func (q *PostgresDB) GetChatByChatID(ctx context.Context, chatID string) (*models.BackendChat, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectChatByChatIDSQL, chatID)
	return q.scanChat(row)
}

func (q *PostgresDB) UpsertChatFeedback(ctx context.Context, feedback models.BackendChatFeedback) error {
	var mark sql.NullInt64
	if feedback.Mark != nil {
		mark.Valid = true
		mark.Int64 = int64(*feedback.Mark)
	}
	var commentTags []string
	var commentText sql.NullString
	if feedback.Comment != nil {
		commentTags = feedback.Comment.Tags
		if commentTags == nil {
			commentTags = []string{}
		}
		commentText.Valid = true
		commentText.String = feedback.Comment.Text
	}
	_, err := q.querier().Exec(ctx, backendqueries.UpsertChatFeedbackSQL, feedback.ChatID, feedback.MessageID, mark, commentTags, commentText)
	return err
}

func (q *PostgresDB) DeleteChatFeedback(ctx context.Context, chatID, messageID string) error {
	_, err := q.querier().Exec(ctx, backendqueries.DeleteChatFeedbackSQL, chatID, messageID)
	return err
}

func (q *PostgresDB) GetChatFeedbackByChatID(ctx context.Context, chatID string) (map[string]models.BackendChatFeedbackValue, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectChatFeedbackByChatIDSQL, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedback := make(map[string]models.BackendChatFeedbackValue)
	for rows.Next() {
		var messageID string
		var mark sql.NullInt64
		var commentTags []string
		var commentText sql.NullString
		if err := rows.Scan(&messageID, &mark, &commentTags, &commentText); err != nil {
			return nil, err
		}
		var comment *models.BackendChatFeedbackComment
		if commentTags != nil || commentText.Valid {
			if commentTags == nil {
				commentTags = []string{}
			}
			comment = &models.BackendChatFeedbackComment{
				Tags: commentTags,
				Text: commentText.String,
			}
		}
		if mark.Valid {
			markValue := int(mark.Int64)
			feedback[messageID] = models.BackendChatFeedbackValue{Mark: &markValue, Comment: comment}
		} else {
			feedback[messageID] = models.BackendChatFeedbackValue{Comment: comment}
		}
	}
	return feedback, rows.Err()
}

func (q *PostgresDB) InsertVizPayload(ctx context.Context, payload models.BackendVizPayload) error {
	if payload.Status == "" {
		payload.Status = models.BackendVizPayloadStatusCompleted
	}
	encrypted, err := q.encryptBackendBlob(normalizeChatHistoryPayload(payload.Payload))
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx,
		backendqueries.InsertVizPayloadSQL,
		payload.ID,
		payload.UserID,
		payload.ChatID,
		payload.MessageID,
		payload.RefID,
		encrypted,
		payload.Status,
		payload.Error,
		payload.CreatedAt,
		payload.UpdatedAt,
	)
	return err
}

func (q *PostgresDB) UpdateVizPayloadStatus(ctx context.Context, payloadID, userID, status, errorMessage string, updatedAt time.Time) error {
	tag, err := q.querier().Exec(ctx,
		backendqueries.UpdateVizPayloadStatusSQL,
		payloadID,
		userID,
		status,
		errorMessage,
		updatedAt,
	)
	if err == nil && tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (q *PostgresDB) UpdateVizPayloadResult(ctx context.Context, payload models.BackendVizPayload) error {
	if payload.Status == "" {
		payload.Status = models.BackendVizPayloadStatusCompleted
	}
	encrypted, err := q.encryptBackendBlob(normalizeChatHistoryPayload(payload.Payload))
	if err != nil {
		return err
	}
	tag, err := q.querier().Exec(ctx,
		backendqueries.UpdateVizPayloadResultSQL,
		payload.ID,
		payload.UserID,
		payload.RefID,
		encrypted,
		payload.Status,
		payload.Error,
		payload.UpdatedAt,
	)
	if err == nil && tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (q *PostgresDB) GetVizPayloadByID(ctx context.Context, payloadID string) (*models.BackendVizPayload, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectVizPayloadByIDSQL, payloadID)
	return q.scanVizPayload(row)
}

func (q *PostgresDB) ListVizPayloadsByUser(
	ctx context.Context,
	userID string,
	before *time.Time,
	limit int,
) ([]models.VizPayloadListItem, error) {
	rows, err := q.querier().Query(ctx, backendqueries.ListVizPayloadsByUserSQL, userID, before, limit)
	if err != nil {
		return nil, fmt.Errorf("list viz payloads: %w", err)
	}
	defer rows.Close()

	out := make([]models.VizPayloadListItem, 0, limit)
	for rows.Next() {
		var (
			id, userIDRow, messageID, refID string
			chatID                          sql.NullString
			encrypted                       []byte
			status, errStr                  string
			createdAt, updatedAt            time.Time
		)
		if err := rows.Scan(&id, &userIDRow, &chatID, &messageID, &refID,
			&encrypted, &status, &errStr, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan viz payload list row: %w", err)
		}
		decrypted, err := q.decryptBackendBlob(encrypted)
		if err != nil {
			return nil, fmt.Errorf("decrypt viz payload %s: %w", id, err)
		}
		var probe struct {
			Type  string `json:"type"`
			Title string `json:"title"`
		}
		_ = json.Unmarshal(decrypted, &probe)
		if probe.Title == "" {
			probe.Title = "untitled"
		}
		item := models.VizPayloadListItem{
			ID:        id,
			MessageID: messageID,
			RefID:     refID,
			Title:     probe.Title,
			VizType:   probe.Type,
			CreatedAt: createdAt,
		}
		if chatID.Valid {
			item.ChatID = &chatID.String
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate viz payload list rows: %w", err)
	}
	return out, nil
}

func (q *PostgresDB) InsertVizDataRef(ctx context.Context, ref models.BackendVizDataRef) error {
	encrypted, err := q.encryptBackendBlob(normalizeChatHistoryPayload(ref.Spec))
	if err != nil {
		return err
	}
	_, err = q.querier().Exec(ctx,
		backendqueries.InsertVizDataRefSQL,
		ref.ID,
		ref.UserID,
		ref.ChatID,
		ref.MessageID,
		ref.RefID,
		ref.Path,
		encrypted,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

func (q *PostgresDB) GetVizDataRefByID(ctx context.Context, refID string) (*models.BackendVizDataRef, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectVizDataRefByIDSQL, refID)
	return q.scanVizDataRef(row)
}

func (q *PostgresDB) UpdateChatHistory(ctx context.Context, history models.BackendChat) (bool, error) {
	payload := normalizeChatHistoryPayload(history.History)
	encrypted, err := q.encryptBackendBlob(payload)
	if err != nil {
		return false, err
	}
	res, err := q.querier().Exec(ctx,
		backendqueries.UpdateChatHistorySQL,
		time.Now(),
		encrypted,
		history.ChatID,
	)
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return false, nil
	}
	return true, nil
}

func (q *PostgresDB) UpdateChatShortName(ctx context.Context, chatID, shortName string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.UpdateChatShortNameSQL, shortName, chatID)
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return false, nil
	}
	return true, nil
}

func (q *PostgresDB) UpdateInactiveChatShortName(ctx context.Context, chatID, shortName string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.UpdateInactiveChatShortNameSQL, shortName, time.Now(), chatID)
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	return rows > 0, nil
}

func (q *PostgresDB) UpdateChatFolderID(ctx context.Context, chatID string, folderID *string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.UpdateChatFolderIDSQL, chatID, folderID, time.Now())
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (q *PostgresDB) UpdateChatShared(ctx context.Context, chatID string, shared bool) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.UpdateChatSharedSQL, shared, chatID)
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	return rows > 0, nil
}

func (q *PostgresDB) LockChat(ctx context.Context, chatID string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.LockChatSQL, chatID)
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return false, nil
	}
	return true, nil
}

func (q *PostgresDB) UnlockChat(ctx context.Context, chatID string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.UnlockChatSQL, chatID)
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return false, nil
	}
	return true, nil
}

func (q *PostgresDB) GetChatsByUserID(ctx context.Context, userID string, limit, offset int) ([]models.BackendUserChatSummary, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectChatsByUserIDSQL, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.BackendUserChatSummary
	for rows.Next() {
		var chat models.BackendUserChatSummary
		var folderID sql.NullString
		if err := rows.Scan(&chat.ChatID, &chat.ShortName, &folderID, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			return nil, err
		}
		if folderID.Valid {
			chat.FolderID = &folderID.String
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}

func (q *PostgresDB) GetChatCountByUserID(ctx context.Context, userID string) (int64, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectChatCountByUserIDSQL, userID)
	var total int64
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (q *PostgresDB) DeleteChat(ctx context.Context, chatID string) (bool, error) {
	res, err := q.querier().Exec(ctx, backendqueries.DeleteChatSQL, chatID, time.Now())
	if err != nil {
		return false, err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return false, nil
	}
	return true, nil
}

func (q *PostgresDB) GetChatShareInfo(ctx context.Context, chatID string) (*models.ChatShareInfo, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectShareInfoByChatIDSQL, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shareInfo := models.ChatShareInfo{}
	var hasRows bool
	for rows.Next() {
		var ownerID string
		var chatShared bool
		var chatActive bool
		if err := rows.Scan(&ownerID, &chatShared, &chatActive); err != nil {
			return nil, err
		}
		hasRows = true
		shareInfo.OwnerID = ownerID
		shareInfo.Shared = chatShared
		shareInfo.Active = chatActive
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !hasRows {
		return nil, nil
	}
	return &shareInfo, nil
}

func (q *PostgresDB) scanBackendUser(row pgx.Row) (*models.BackendUser, error) {
	var user models.BackendUser
	var profiles []byte
	var profileExtraRaw []byte
	var thirdwebUserID sql.NullString
	var source sql.NullString
	var referralCode sql.NullString
	var referrerCode sql.NullString
	var referralWallet sql.NullString
	if err := row.Scan(
		&user.ID,
		&thirdwebUserID,
		&source,
		&profiles,
		&referralCode,
		&referrerCode,
		&referralWallet,
		&profileExtraRaw,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	user.ThirdwebUserID = thirdwebUserID.String
	user.Source = source.String
	if referralCode.Valid {
		user.ReferralCode = &referralCode.String
	}
	if referrerCode.Valid {
		user.ReferrerCode = &referrerCode.String
	}
	if referralWallet.Valid {
		user.ReferralWallet = &referralWallet.String
	}
	var plaintext []byte
	if len(profiles) > 0 {
		var err error
		plaintext, err = q.decryptBackendBlob(profiles)
		if err != nil {
			return nil, err
		}
	}
	decoded, err := unmarshalBackendProfiles(plaintext)
	if err != nil {
		return nil, err
	}
	user.Profiles = decoded
	if len(profileExtraRaw) == 0 {
		user.ProfileExtra = models.EmptyProfileExtra()
	} else if err := json.Unmarshal(profileExtraRaw, &user.ProfileExtra); err != nil {
		return nil, fmt.Errorf("unmarshal profile_extra: %w", err)
	}
	return &user, nil
}

func (q *PostgresDB) scanChat(row pgx.Row) (*models.BackendChat, error) {
	var history models.BackendChat
	var folderID sql.NullString
	if err := row.Scan(
		&history.ChatID,
		&history.UserID,
		&history.ShortName,
		&folderID,
		&history.CreatedAt,
		&history.UpdatedAt,
		&history.History,
		&history.Active,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if folderID.Valid {
		history.FolderID = &folderID.String
	}
	if len(history.History) == 0 {
		history.History = json.RawMessage("[]")
		return &history, nil
	}
	decrypted, err := q.decryptBackendBlob([]byte(history.History))
	if err != nil {
		return nil, err
	}
	if len(decrypted) == 0 {
		decrypted = []byte("[]")
	}
	history.History = json.RawMessage(decrypted)
	return &history, nil
}

func scanChatFolder(row pgx.Row) (*models.BackendChatFolder, error) {
	var folder models.BackendChatFolder
	if err := row.Scan(
		&folder.ID,
		&folder.UserID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &folder, nil
}

func (q *PostgresDB) scanVizPayload(row pgx.Row) (*models.BackendVizPayload, error) {
	var payload models.BackendVizPayload
	var chatID sql.NullString
	var encrypted []byte
	if err := row.Scan(
		&payload.ID,
		&payload.UserID,
		&chatID,
		&payload.MessageID,
		&payload.RefID,
		&encrypted,
		&payload.Status,
		&payload.Error,
		&payload.CreatedAt,
		&payload.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if chatID.Valid {
		payload.ChatID = &chatID.String
	}
	decrypted, err := q.decryptBackendBlob(encrypted)
	if err != nil {
		return nil, err
	}
	payload.Payload = json.RawMessage(decrypted)
	return &payload, nil
}

func (q *PostgresDB) scanVizDataRef(row pgx.Row) (*models.BackendVizDataRef, error) {
	var ref models.BackendVizDataRef
	var chatID sql.NullString
	var encrypted []byte
	if err := row.Scan(
		&ref.ID,
		&ref.UserID,
		&chatID,
		&ref.MessageID,
		&ref.RefID,
		&ref.Path,
		&encrypted,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if chatID.Valid {
		ref.ChatID = &chatID.String
	}
	decrypted, err := q.decryptBackendBlob(encrypted)
	if err != nil {
		return nil, err
	}
	ref.Spec = json.RawMessage(decrypted)
	return &ref, nil
}

func marshalBackendProfiles(profiles []thirdweb.UserInfoProfile) ([]byte, error) {
	if profiles == nil {
		profiles = []thirdweb.UserInfoProfile{}
	}
	return json.Marshal(profiles)
}

func unmarshalBackendProfiles(payload []byte) ([]thirdweb.UserInfoProfile, error) {
	if len(payload) == 0 {
		return []thirdweb.UserInfoProfile{}, nil
	}
	var profiles []thirdweb.UserInfoProfile
	if err := json.Unmarshal(payload, &profiles); err != nil {
		return []thirdweb.UserInfoProfile{}, err
	}
	return profiles, nil
}

func normalizeChatHistoryPayload(payload json.RawMessage) []byte {
	if len(payload) == 0 {
		return []byte("[]")
	}
	return payload
}

func (q *PostgresDB) encryptBackendBlob(plaintext []byte) ([]byte, error) {
	if q.backendCipher == nil {
		return nil, errors.New("backend encryption cipher is not configured")
	}
	nonce, err := generateNonce(q.backendCipher.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := q.backendCipher.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (q *PostgresDB) decryptBackendBlob(blob []byte) ([]byte, error) {
	if q.backendCipher == nil {
		return nil, errors.New("backend encryption cipher is not configured")
	}
	nonceSize := q.backendCipher.NonceSize()
	if len(blob) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: %d", len(blob))
	}
	nonce := blob[:nonceSize]
	ciphertext := blob[nonceSize:]
	plaintext, err := q.backendCipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt backend payload: %w", err)
	}
	return plaintext, nil
}

func (q *PostgresDB) GetChatActive(ctx context.Context, chatID string) (bool, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectChatActiveSQL, chatID)
	var active bool
	if err := row.Scan(&active); err != nil {
		return false, err
	}
	return active, nil
}

func (q *PostgresDB) GetProfileExtra(ctx context.Context, userID string) (models.ProfileExtra, error) {
	var raw []byte
	err := q.querier().QueryRow(ctx, backendqueries.SelectProfileExtraByUserIDSQL, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EmptyProfileExtra(), pgx.ErrNoRows
		}
		return models.EmptyProfileExtra(), fmt.Errorf("query profile_extra: %w", err)
	}
	var p models.ProfileExtra
	if err := json.Unmarshal(raw, &p); err != nil {
		return models.EmptyProfileExtra(), fmt.Errorf("unmarshal profile_extra: %w", err)
	}
	return p, nil
}

func (q *PostgresDB) GetProfileExtraForUpdate(ctx context.Context, userID string) (models.ProfileExtra, error) {
	var raw []byte
	err := q.querier().QueryRow(ctx, backendqueries.SelectProfileExtraByUserIDForUpdateSQL, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EmptyProfileExtra(), pgx.ErrNoRows
		}
		return models.EmptyProfileExtra(), fmt.Errorf("query profile_extra for update: %w", err)
	}
	var p models.ProfileExtra
	if err := json.Unmarshal(raw, &p); err != nil {
		return models.EmptyProfileExtra(), fmt.Errorf("unmarshal profile_extra: %w", err)
	}
	return p, nil
}

func (q *PostgresDB) UpdateProfileExtra(ctx context.Context, userID string, profile models.ProfileExtra) error {
	if err := models.ValidateProfileExtra(profile); err != nil {
		return fmt.Errorf("validate profile_extra: %w", err)
	}
	raw, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("marshal profile_extra: %w", err)
	}
	tag, err := q.querier().Exec(ctx, backendqueries.UpdateProfileExtraSQL, userID, raw)
	if err != nil {
		return fmt.Errorf("exec update profile_extra: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// scanSignal deserialises a signals row. It accepts both pgx.Row (single-row
// QueryRow) and pgx.Rows (used inside rows.Next() loops), since pgx.Rows
// satisfies the pgx.Row interface.
func scanSignal(row pgx.Row) (*models.Signal, error) {
	var s models.Signal
	var lastValue, spec []byte
	var lastErr *string
	if err := row.Scan(&s.ID, &s.UserID, &s.Kind, &s.Status, &s.NLText, &spec, &s.HumanReadable,
		&s.CadenceMinutes, &s.CooldownMinutes, &lastValue, &s.LastEvaluatedAt, &s.LastFiredAt,
		&s.NextEvalAt, &s.ErrorCount, &lastErr, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	s.Spec = spec
	s.LastValue = lastValue
	s.LastError = lastErr
	return &s, nil
}

func (q *PostgresDB) InsertSignal(ctx context.Context, s models.Signal) error {
	now := time.Now().UTC()
	_, err := q.querier().Exec(ctx, backendqueries.InsertSignalSQL,
		s.ID, s.UserID, s.Kind, s.Status, s.NLText, []byte(s.Spec), s.HumanReadable,
		s.CadenceMinutes, s.CooldownMinutes, now, now)
	return err
}

// InsertSignalWithCap atomically enforces the per-kind active-signal cap: an
// advisory xact lock serializes concurrent creates for the same (user, kind),
// then count + insert run inside the same transaction. Returns false (and
// inserts nothing) when the cap is already reached — closes the TOCTOU window
// of a separate CountActiveSignals + InsertSignal pair.
func (q *PostgresDB) InsertSignalWithCap(ctx context.Context, s models.Signal, maxActive int) (bool, error) {
	if q.tx != nil {
		// Already inside a transaction — compose into it; the caller commits.
		return insertSignalWithCapTx(ctx, q.tx, s, maxActive)
	}
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	inserted, err := insertSignalWithCapTx(ctx, tx, s, maxActive)
	if err != nil || !inserted {
		return inserted, err
	}
	return true, tx.Commit(ctx)
}

func insertSignalWithCapTx(ctx context.Context, tx pgx.Tx, s models.Signal, maxActive int) (bool, error) {
	if _, err := tx.Exec(ctx, backendqueries.AcquireSignalCapLockSQL, s.UserID, s.Kind); err != nil {
		return false, err
	}
	var count int64
	if err := tx.QueryRow(ctx, backendqueries.CountActiveSignalsSQL, s.UserID, s.Kind).Scan(&count); err != nil {
		return false, err
	}
	if count >= int64(maxActive) {
		return false, nil
	}
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, backendqueries.InsertSignalSQL,
		s.ID, s.UserID, s.Kind, s.Status, s.NLText, []byte(s.Spec), s.HumanReadable,
		s.CadenceMinutes, s.CooldownMinutes, now, now); err != nil {
		return false, err
	}
	return true, nil
}

// LeaseSignal pushes next_eval_at forward as a short claim lease so the
// evaluator can release its row locks before calling Python. Unlike
// UpdateSignalAfterEvaluation it does not reset error_count or stamp
// last_evaluated_at. If the signal was deleted meanwhile the UPDATE affects
// zero rows, which is fine.
func (q *PostgresDB) LeaseSignal(ctx context.Context, id string, until time.Time) error {
	_, err := q.querier().Exec(ctx, backendqueries.LeaseSignalSQL, id, until)
	return err
}

func (q *PostgresDB) GetSignalsByUserID(ctx context.Context, userID string) ([]models.Signal, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectSignalsByUserSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	signals := []models.Signal{}
	for rows.Next() {
		s, err := scanSignal(rows)
		if err != nil {
			return nil, err
		}
		signals = append(signals, *s)
	}
	return signals, rows.Err()
}

func (q *PostgresDB) GetSignalByID(ctx context.Context, id, userID string) (*models.Signal, error) {
	row := q.querier().QueryRow(ctx, backendqueries.SelectSignalByIDSQL, id, userID)
	return scanSignal(row)
}

func (q *PostgresDB) UpdateSignal(ctx context.Context, id, userID string, status *string, cadenceMinutes *int) (bool, error) {
	tag, err := q.querier().Exec(ctx, backendqueries.UpdateSignalSQL, id, userID, status, cadenceMinutes)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (q *PostgresDB) DeleteSignal(ctx context.Context, id, userID string) (bool, error) {
	tag, err := q.querier().Exec(ctx, backendqueries.DeleteSignalSQL, id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (q *PostgresDB) CountActiveSignals(ctx context.Context, userID, kind string) (int64, error) {
	row := q.querier().QueryRow(ctx, backendqueries.CountActiveSignalsSQL, userID, kind)
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (q *PostgresDB) SelectDueSignalsForUpdate(ctx context.Context, limit int) ([]models.Signal, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectDueSignalsSQL, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []models.Signal
	for rows.Next() {
		s, err := scanSignal(rows)
		if err != nil {
			return nil, err
		}
		signals = append(signals, *s)
	}
	return signals, rows.Err()
}

func (q *PostgresDB) UpdateSignalAfterEvaluation(ctx context.Context, id string, fired bool, lastValue json.RawMessage, nextEvalAt time.Time, evalErr *string) error {
	var lastValueBytes []byte
	if len(lastValue) > 0 {
		lastValueBytes = []byte(lastValue)
	}
	_, err := q.querier().Exec(ctx, backendqueries.UpdateSignalAfterEvalSQL,
		id, fired, lastValueBytes, nextEvalAt, evalErr)
	return err
}

func (q *PostgresDB) MarkSignalError(ctx context.Context, id string, lastError string) error {
	_, err := q.querier().Exec(ctx, backendqueries.MarkSignalErrorSQL, id, lastError)
	return err
}

func (q *PostgresDB) InsertSignalEvent(ctx context.Context, e models.SignalEvent) error {
	var valueSnapshot []byte
	if len(e.ValueSnapshot) > 0 {
		valueSnapshot = []byte(e.ValueSnapshot)
	}
	_, err := q.querier().Exec(ctx, backendqueries.InsertSignalEventSQL,
		e.ID, e.SignalID, e.UserID, e.FiredAt, e.Title, e.Body, valueSnapshot, e.SeedQuery, e.Channel)
	return err
}

func (q *PostgresDB) GetSignalEventsByUserID(ctx context.Context, userID string, since *time.Time, limit int) ([]models.SignalEvent, error) {
	rows, err := q.querier().Query(ctx, backendqueries.SelectSignalEventsSQL, userID, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.SignalEvent{}
	for rows.Next() {
		var e models.SignalEvent
		var valueSnapshot []byte
		if err := rows.Scan(&e.ID, &e.SignalID, &e.UserID, &e.FiredAt, &e.Title, &e.Body,
			&valueSnapshot, &e.SeedQuery, &e.SeenAt, &e.Channel); err != nil {
			return nil, err
		}
		e.ValueSnapshot = json.RawMessage(valueSnapshot)
		events = append(events, e)
	}
	return events, rows.Err()
}

func (q *PostgresDB) MarkSignalEventSeen(ctx context.Context, id, userID string) (bool, error) {
	tag, err := q.querier().Exec(ctx, backendqueries.MarkSignalEventSeenSQL, id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (q *PostgresDB) PruneSignalEvents(ctx context.Context, userID string, keep int) error {
	_, err := q.querier().Exec(ctx, backendqueries.PruneSignalEventsSQL, userID, keep)
	return err
}

func (q *PostgresDB) ResetSignalErrorState(ctx context.Context, id, userID string) error {
	_, err := q.querier().Exec(ctx, backendqueries.ResetSignalErrorStateSQL, id, userID)
	return err
}

func (q *PostgresDB) CountDueSignals(ctx context.Context) (int64, error) {
	row := q.querier().QueryRow(ctx, backendqueries.CountDueSignalsSQL)
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
