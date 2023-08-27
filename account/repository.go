package account

import (
	"context"
	"database/sql"
	"errors"

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

	account.SetId("")

	_, err = r.sql.ExecContext(ctx, query, account.Id, account.Username, password, account.Email)

	return err
}

func (r *Repository) FindByEmail(ctx context.Context, account value.AccountEntity) (value.AccountEntity, error) {
	query := `SELECT id, username, email, password FROM accounts WHERE email = $1`

	err := r.sql.
		QueryRowContext(ctx, query, account.Email).
		Scan(&account.Id, &account.Username, &account.Email, &account.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return account, ErrAccountNotFound
	}

	if err != nil {
		return account, err
	}

	return account, nil
}
