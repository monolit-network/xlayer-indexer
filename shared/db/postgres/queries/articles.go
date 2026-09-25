package queries

const (
	UpsertSourceFeedSQL = `
		INSERT INTO article_sources (source, feed)
		VALUES ($1, $2)
		ON CONFLICT (source)
		DO UPDATE SET 
			feed = EXCLUDED.feed`

	SelectArticlesBaseSQL = `
		SELECT id, link, title, text, authors, date, source
		FROM articles`

	SelectArticlesByIDsSQL = `
		SELECT id, link, title, text, authors, date, source
		FROM articles
		WHERE id = ANY($1)`
)
