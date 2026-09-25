package queries

const (
	CreateTableBackendSessionsSQL = `CREATE TABLE IF NOT EXISTS backend.sessions (
		user_id uuid NOT NULL REFERENCES backend.users(id),
		session_id uuid NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		PRIMARY KEY (user_id, session_id)
	);

	CREATE INDEX IF NOT EXISTS idx_backend_sessions_user_id ON backend.sessions (user_id);
	CREATE INDEX IF NOT EXISTS idx_backend_sessions_session_id ON backend.sessions (session_id);`

	InsertUserSessionSQL = `INSERT INTO backend.sessions (user_id, session_id, created_at, expires_at) VALUES ($1, $2, $3, $4);`

	SelectUserSessionByUserIDAndSessionIDSQL = `SELECT user_id, session_id, created_at, expires_at FROM backend.sessions WHERE user_id = $1 AND session_id = $2;`

	SelectUserSessionByUserIDSQL = `SELECT user_id, session_id, created_at, expires_at FROM backend.sessions WHERE user_id = $1;`

	SelectUserSessionBySessionIDSQL = `SELECT user_id, session_id, created_at, expires_at FROM backend.sessions WHERE session_id = $1;`

	UpdateUserSessionSQL = `UPDATE backend.sessions SET expires_at = $1 WHERE user_id = $2 AND session_id = $3;`

	DeleteUserSessionByUserIDAndSessionIDSQL = `DELETE FROM backend.sessions WHERE user_id = $1 AND session_id = $2;`
)
