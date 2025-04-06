package repository

const (
	createLink = `
		INSERT INTO links (alias, original_url, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (alias) DO UPDATE
		SET original_url = $2, expires_at = $3
		RETURNING id
	`

	getLinkByAlias = `
			SELECT * FROM links
			WHERE alias = $1 AND (expires_at IS NULL OR expires_at > $2)
		`

	existsLink = `SELECT EXISTS(SELECT 1 FROM links WHERE alias = $1)`
)
