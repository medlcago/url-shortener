package repository

const (
	createLink = `
		INSERT INTO links (alias, original_url, expires_at, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	getLinkByAlias = `
			SELECT * FROM links
			WHERE alias = $1 AND (expires_at IS NULL OR expires_at > $2)
		`

	existsLink = `SELECT EXISTS(SELECT 1 FROM links WHERE alias = $1)`
)
