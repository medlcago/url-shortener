package repository

const (
	createUser = `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING *
	`

	getUserByID = `SELECT * FROM users WHERE id = $1`

	getUserByEmail = `SELECT * FROM users WHERE email = $1`

	updateLastLogin = `UPDATE users SET login_date = $1 WHERE id = $2`
)
