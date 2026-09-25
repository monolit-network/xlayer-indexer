package queries

const (
	// Social posts
	CreateSocialPostsTableSQL = `
		CREATE TABLE IF NOT EXISTS social_posts (
			link TEXT PRIMARY KEY,
			text TEXT NOT NULL,
			date TIMESTAMP NOT NULL,
			source VARCHAR(32) NOT NULL,
			author TEXT NOT NULL,
			relations TEXT[] NOT NULL DEFAULT '{}',
			links_mentioned TEXT[] NOT NULL DEFAULT '{}',
			tokens_mentioned TEXT[] NOT NULL DEFAULT '{}'
		);

		CREATE INDEX IF NOT EXISTS social_posts_date_idx ON social_posts (date DESC);
		CREATE INDEX IF NOT EXISTS social_posts_source_idx ON social_posts (LOWER(source));
		CREATE INDEX IF NOT EXISTS social_posts_author_idx ON social_posts (LOWER(author));
		CREATE INDEX IF NOT EXISTS social_posts_relations_gin_idx ON social_posts USING GIN (relations);
		CREATE INDEX IF NOT EXISTS social_posts_link_prefix_ci_idx ON social_posts (LOWER(link) text_pattern_ops);
		`

	UpsertSocialPostSQL = `
		INSERT INTO social_posts (link, text, date, source, author, relations, links_mentioned, tokens_mentioned)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (link)
		DO UPDATE SET
			text = EXCLUDED.text,
			date = EXCLUDED.date,
			source = EXCLUDED.source,
			author = EXCLUDED.author,
			relations = EXCLUDED.relations,
			links_mentioned = EXCLUDED.links_mentioned,
			tokens_mentioned = EXCLUDED.tokens_mentioned
		RETURNING xmax = 0 as inserted`

	SelectSocialPostByLinkSQL = `
		SELECT link, text, date, source, author, relations, links_mentioned, tokens_mentioned
		FROM social_posts
		WHERE link = $1`

	SelectSocialPostsByLinkBeginningSQL = `
		SELECT link, text, date, source, author, relations, links_mentioned, tokens_mentioned
		FROM social_posts
		WHERE LOWER(link) LIKE LOWER($1) || '%'
		ORDER BY date DESC`

	SelectPostsReferencingLinkSQL = `
		SELECT link, text, date, source, author, relations, links_mentioned, tokens_mentioned
		FROM social_posts
		WHERE $1 = ANY(relations)
		ORDER BY date DESC
		LIMIT $2`

	CreateSocialSourcesTableSQL = `
		CREATE TABLE IF NOT EXISTS social_sources (
			value VARCHAR(255) PRIMARY KEY,
			source VARCHAR(32) NOT NULL
		);`

	UpsertSocialSourceSQL = `
		INSERT INTO social_sources (value, source)
		VALUES ($1, $2)
		ON CONFLICT (value)
		DO UPDATE SET 
			source = EXCLUDED.source`

	SelectSocialSourcesSQL = `
		SELECT value
		FROM social_sources
		WHERE source = $1`

	SelectOldestPostSQL = `
		SELECT link, text, date, source, author, relations, links_mentioned, tokens_mentioned
		FROM social_posts
		WHERE LOWER(source) = LOWER($1) AND LOWER(author) = LOWER($2)
		ORDER BY date ASC
		LIMIT 1`

	SelectNewestPostSQL = `
		SELECT link, text, date, source, author, relations, links_mentioned, tokens_mentioned
		FROM social_posts
		WHERE LOWER(source) = LOWER($1) AND LOWER(author) = LOWER($2)
		ORDER BY date DESC
		LIMIT 1`

	SelectSocialPostsByAuthorSQL = `
		SELECT link, text, date, source, author, relations, links_mentioned, tokens_mentioned
		FROM social_posts
		WHERE LOWER(author) = LOWER($1)
		ORDER BY date DESC`
)
