package queries

const (
	CreateTableBackendChatHistorySQL = `CREATE TABLE IF NOT EXISTS backend.chat_history (
		chat_id uuid PRIMARY KEY,
		user_id uuid NOT NULL REFERENCES backend.users(id),
		short_name TEXT NOT NULL,
		folder_id uuid NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted boolean NOT NULL DEFAULT FALSE,
		active boolean NOT NULL DEFAULT FALSE,
		shared boolean NOT NULL DEFAULT FALSE,
		history BYTEA NOT NULL
	);

	ALTER TABLE backend.chat_history ADD COLUMN IF NOT EXISTS folder_id uuid NULL;
	
	CREATE INDEX IF NOT EXISTS idx_backend_chat_history_user_id ON backend.chat_history (user_id);
	CREATE INDEX IF NOT EXISTS idx_backend_chat_history_folder_id ON backend.chat_history (folder_id);`

	CreateTableBackendChatFoldersSQL = `CREATE TABLE IF NOT EXISTS backend.chat_folders (
		id uuid PRIMARY KEY,
		user_id uuid NOT NULL REFERENCES backend.users(id),
		name TEXT NOT NULL,
		deleted boolean NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	ALTER TABLE backend.chat_history ADD COLUMN IF NOT EXISTS folder_id uuid NULL;
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_constraint
			WHERE conname = 'backend_chat_history_folder_id_fkey'
		) THEN
			ALTER TABLE backend.chat_history
				ADD CONSTRAINT backend_chat_history_folder_id_fkey
				FOREIGN KEY (folder_id) REFERENCES backend.chat_folders(id) ON DELETE SET NULL;
		END IF;
	END $$;

	CREATE INDEX IF NOT EXISTS idx_backend_chat_folders_user_id ON backend.chat_folders (user_id);
	CREATE INDEX IF NOT EXISTS idx_backend_chat_history_folder_id ON backend.chat_history (folder_id);`

	CreateTableBackendChatFeedbackSQL = `CREATE TABLE IF NOT EXISTS backend.chat_feedback (
		chat_id uuid NOT NULL REFERENCES backend.chat_history(chat_id) ON DELETE CASCADE,
		message_id TEXT NOT NULL,
		mark INTEGER NULL,
		comment_tags TEXT[] NULL,
		comment_text TEXT NULL,
		PRIMARY KEY (chat_id, message_id),
		CONSTRAINT backend_chat_feedback_mark_check CHECK (mark IS NULL OR mark IN (0, 1)),
		CONSTRAINT backend_chat_feedback_comment_tags_count_check CHECK (comment_tags IS NULL OR cardinality(comment_tags) <= 20),
		CONSTRAINT backend_chat_feedback_comment_text_length_check CHECK (comment_text IS NULL OR char_length(comment_text) <= 300)
	);`

	CreateTableBackendVizPayloadsSQL = `CREATE TABLE IF NOT EXISTS backend.viz_payloads (
		id uuid PRIMARY KEY,
		user_id uuid NOT NULL REFERENCES backend.users(id),
		chat_id uuid NULL REFERENCES backend.chat_history(chat_id) ON DELETE SET NULL,
		message_id TEXT NOT NULL DEFAULT '',
		ref_id TEXT NOT NULL DEFAULT '',
		payload BYTEA NOT NULL,
		status TEXT NOT NULL DEFAULT 'completed',
		error TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	ALTER TABLE backend.viz_payloads ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'completed';
	ALTER TABLE backend.viz_payloads ADD COLUMN IF NOT EXISTS error TEXT NOT NULL DEFAULT '';

	CREATE INDEX IF NOT EXISTS idx_backend_viz_payloads_user_id ON backend.viz_payloads (user_id);
	CREATE INDEX IF NOT EXISTS idx_backend_viz_payloads_chat_id ON backend.viz_payloads (chat_id);
	CREATE INDEX IF NOT EXISTS idx_backend_viz_payloads_message_id ON backend.viz_payloads (message_id);`

	CreateTableBackendVizDataRefsSQL = `CREATE TABLE IF NOT EXISTS backend.viz_data_refs (
		id uuid PRIMARY KEY,
		user_id uuid NOT NULL REFERENCES backend.users(id),
		chat_id uuid NULL REFERENCES backend.chat_history(chat_id) ON DELETE SET NULL,
		message_id TEXT NOT NULL DEFAULT '',
		ref_id TEXT NOT NULL DEFAULT '',
		path TEXT NOT NULL DEFAULT '',
		spec BYTEA NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_backend_viz_data_refs_user_id ON backend.viz_data_refs (user_id);
	CREATE INDEX IF NOT EXISTS idx_backend_viz_data_refs_chat_id ON backend.viz_data_refs (chat_id);`

	InsertChatSQL = `INSERT INTO backend.chat_history (chat_id, user_id, short_name, folder_id, created_at, updated_at, history, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	SelectChatByChatIDSQL     = `SELECT chat_id, user_id, short_name, folder_id, created_at, updated_at, history, active FROM backend.chat_history WHERE chat_id = $1 AND deleted = FALSE;`
	SelectChatUserByChatIDSQL = `SELECT user_id FROM backend.chat_history WHERE chat_id = $1 AND deleted = FALSE;`

	// Most-recent untouched chat for a user: empty history, never
	// activated, not deleted. Lets HandleCreateChat reuse the same
	// reserved chat_id when the frontend retries upload-creation,
	// instead of spamming "New Chat" rows.
	//
	// The literal `'[]'` MUST match the Go-side ``emptyChatHistoryJSON``
	// constant in backend/pkg/controllers/chat/controller.go — both
	// sides describe the same "no messages yet" marker.
	SelectMostRecentEmptyChatSQL = `SELECT chat_id FROM backend.chat_history
		WHERE user_id = $1 AND history = '[]' AND active = FALSE AND deleted = FALSE
		ORDER BY created_at DESC LIMIT 1;`

	InsertChatFolderSQL = `INSERT INTO backend.chat_folders (id, user_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5);`

	SelectChatFoldersByUserIDSQL = `SELECT id, user_id, name, created_at, updated_at
		FROM backend.chat_folders
		WHERE user_id = $1 AND deleted = FALSE
		ORDER BY updated_at DESC, name ASC;`

	SelectChatFolderByIDSQL = `SELECT id, user_id, name, created_at, updated_at
		FROM backend.chat_folders
		WHERE id = $1 AND deleted = FALSE;`

	UpdateChatFolderNameSQL = `UPDATE backend.chat_folders
		SET name = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2 AND deleted = FALSE;`

	DeleteChatFolderSQL = `UPDATE backend.chat_folders
		SET deleted = TRUE, updated_at = $3
		WHERE id = $1 AND user_id = $2 AND deleted = FALSE;`

	ClearChatFolderAssignmentsSQL = `UPDATE backend.chat_history
		SET folder_id = NULL, updated_at = $2
		WHERE folder_id = $1 AND deleted = FALSE;`

	UpsertChatFeedbackSQL = `INSERT INTO backend.chat_feedback (chat_id, message_id, mark, comment_tags, comment_text) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id, message_id) DO UPDATE SET mark = EXCLUDED.mark, comment_tags = EXCLUDED.comment_tags, comment_text = EXCLUDED.comment_text;`

	DeleteChatFeedbackSQL = `DELETE FROM backend.chat_feedback WHERE chat_id = $1 AND message_id = $2;`

	SelectChatFeedbackByChatIDSQL = `SELECT message_id, mark, comment_tags, comment_text FROM backend.chat_feedback WHERE chat_id = $1;`

	InsertVizPayloadSQL = `INSERT INTO backend.viz_payloads (id, user_id, chat_id, message_id, ref_id, payload, status, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	SelectVizPayloadByIDSQL = `SELECT id, user_id, chat_id, message_id, ref_id, payload, status, error, created_at, updated_at
		FROM backend.viz_payloads WHERE id = $1;`

	ListVizPayloadsByUserSQL = `SELECT id, user_id, chat_id, message_id, ref_id, payload, status, error, created_at, updated_at
		FROM backend.viz_payloads
		WHERE user_id = $1
		  AND status = 'completed'
		  AND ($2::timestamp IS NULL OR created_at < $2)
		ORDER BY created_at DESC
		LIMIT $3;`

	UpdateVizPayloadStatusSQL = `UPDATE backend.viz_payloads
		SET status = $3, error = $4, updated_at = $5
		WHERE id = $1 AND user_id = $2;`

	UpdateVizPayloadResultSQL = `UPDATE backend.viz_payloads
		SET ref_id = $3, payload = $4, status = $5, error = $6, updated_at = $7
		WHERE id = $1 AND user_id = $2;`

	InsertVizDataRefSQL = `INSERT INTO backend.viz_data_refs (id, user_id, chat_id, message_id, ref_id, path, spec, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	SelectVizDataRefByIDSQL = `SELECT id, user_id, chat_id, message_id, ref_id, path, spec, created_at, updated_at
		FROM backend.viz_data_refs WHERE id = $1;`

	SelectChatsByUserIDSQL     = `SELECT chat_id, short_name, folder_id, created_at, updated_at FROM backend.chat_history WHERE user_id = $1 AND deleted = FALSE ORDER BY updated_at DESC LIMIT $2 OFFSET $3;`
	SelectChatCountByUserIDSQL = `SELECT COUNT(*) FROM backend.chat_history WHERE user_id = $1 AND deleted = FALSE;`

	SelectChatActiveSQL = `SELECT active FROM backend.chat_history WHERE chat_id = $1 AND deleted = FALSE;`

	UpdateChatHistorySQL = `UPDATE backend.chat_history SET updated_at = $1, history = $2 WHERE chat_id = $3 AND deleted = FALSE RETURNING *;`

	UpdateChatShortNameSQL = `UPDATE backend.chat_history SET short_name = $1 WHERE chat_id = $2 AND deleted = FALSE RETURNING *;`

	UpdateInactiveChatShortNameSQL = `UPDATE backend.chat_history SET short_name = $1, updated_at = $2 WHERE chat_id = $3 AND active = FALSE AND deleted = FALSE RETURNING *;`

	UpdateChatFolderIDSQL = `UPDATE backend.chat_history SET folder_id = $2, updated_at = $3 WHERE chat_id = $1 AND deleted = FALSE;`

	DeleteChatSQL = `UPDATE backend.chat_history SET deleted = TRUE, deleted_at = $2 WHERE chat_id = $1 AND deleted = FALSE RETURNING *;`

	LockChatSQL = `UPDATE backend.chat_history SET active = TRUE WHERE chat_id = $1 AND active = FALSE AND deleted = FALSE RETURNING *;`

	UnlockChatSQL = `UPDATE backend.chat_history SET active = FALSE WHERE chat_id = $1 AND active = TRUE AND deleted = FALSE RETURNING *;`

	UpdateChatSharedSQL = `UPDATE backend.chat_history SET shared = $1 WHERE chat_id = $2 AND deleted = FALSE RETURNING *;`

	SelectShareInfoByChatIDSQL = `   
			SELECT 
				ch.user_id,
				ch.shared AS chat_shared,
				ch.active AS chat_active
			FROM backend.chat_history ch
			WHERE ch.chat_id = $1
			AND ch.deleted = FALSE;`
)
