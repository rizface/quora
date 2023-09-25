package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rizface/quora/account/value"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	tracer trace.Tracer
	sql    *sql.DB
}

func NewRepository(sql *sql.DB, tracer trace.Tracer) *Repository {
	return &Repository{sql: sql, tracer: tracer}
}

func emailIsUsed(ctx context.Context, sql *sql.DB, email string) (bool, error) {
	span := trace.SpanFromContext(ctx)

	var (
		count int
		err   = sql.
			QueryRowContext(ctx, `SELECT COUNT(id) as count FROM accounts WHERE email = $1`, email).
			Scan(&count)
		emailIsUsed = count > 0
	)

	span.AddEvent("check email availability", trace.WithAttributes(
		attribute.KeyValue{
			Key:   "email",
			Value: attribute.StringValue(email),
		},
		attribute.KeyValue{
			Key:   "isUsed",
			Value: attribute.BoolValue(emailIsUsed),
		},
	))

	return emailIsUsed, err
}

func usernameIsUsed(ctx context.Context, sql *sql.DB, username string) (bool, error) {
	span := trace.SpanFromContext(ctx)

	var (
		count int
		err   = sql.
			QueryRowContext(ctx, `SELECT COUNT(id) as count FROM accounts WHERE username = $1`, username).
			Scan(&count)
		usernameIsUsed = count > 0
	)

	span.AddEvent("check username availability", trace.WithAttributes(
		attribute.KeyValue{
			Key:   "username",
			Value: attribute.StringValue(username),
		},
		attribute.KeyValue{
			Key:   "isUsed",
			Value: attribute.BoolValue(usernameIsUsed),
		},
	))

	return usernameIsUsed, err
}

func (r *Repository) Create(ctx context.Context, account value.AccountEntity) error {
	ctx, span := r.tracer.Start(ctx, "account.Repository.Create")
	defer span.End()

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
	ctx, span := r.tracer.Start(ctx, "account.Repository.FindByEmail")
	defer span.End()

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
