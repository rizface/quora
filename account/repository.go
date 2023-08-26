package account

import (
	"context"
	"database/sql"

	"github.com/rizface/quora/account/value"
)

type Repository struct {
	sql *sql.DB
}

func NewRepository(sql *sql.DB) *Repository {
	return &Repository{sql: sql}
}

func emailIsUsed(ctx context.Context, sql *sql.DB, email string) (bool, error) {
	var count int
	return count > 0, sql.
		QueryRowContext(ctx, `SELECT COUNT(id) as count FROM accounts WHERE email = $1`, email).
		Scan(&count)
}

func usernameIsUsed(ctx context.Context, sql *sql.DB, username string) (bool, error) {
	var count int

	return count > 0, sql.
		QueryRowContext(ctx, `SELECT COUNT(id) as count FROM accounts WHERE username = $1`, username).
		Scan(&count)
}

func (r *Repository) Create(ctx context.Context, account value.AccountEntity) error {
	if used, err := emailIsUsed(ctx, r.sql, account.Email); err != nil {
		return err
	} else if used {
		return ErrEmailIsUsed
	}

	if used, err := usernameIsUsed(ctx, r.sql, account.Username); err != nil {
		return err
	} else if used {
		return ErrUsernameIsUsed
	}

	query := `
		INSERT INTO accounts (id, username, password, email) VALUES ($1, $2, $3, $4)
	`

	password, err := account.GetPasswordHash()
	if err != nil {
		return err
	}

	_, err = r.sql.ExecContext(ctx, query, account.Id, account.Username, password, account.Email)

	return err
}
