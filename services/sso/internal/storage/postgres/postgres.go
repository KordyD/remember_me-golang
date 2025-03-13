package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/kordyd/remember_me-golang/services/sso/internal/models"
	"github.com/kordyd/remember_me-golang/services/sso/internal/storage"
	"github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) GetApp(ctx context.Context, id string) (models.App, error) {
	const op = "postgres.GetApp"
	row := p.db.QueryRowContext(ctx, `SELECT id, name, secret FROM apps WHERE id = $1`, id)
	var app models.App
	err := row.Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

func (p *Postgres) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.GetUser"
	row := p.db.QueryRowContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1", email)
	var user models.User
	err := row.Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (p *Postgres) SaveUser(ctx context.Context, email string, passHash []byte) (string, error) {
	const op = "postgres.SaveUser"
	id := uuid.New()
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (id, email, pass_hash) VALUES ($1, $2, $3)", id.String(), email, string(passHash))
	var postgresErr *pq.Error
	if err != nil {
		if errors.As(err, &postgresErr) && postgresErr.Code == pgerrcode.UniqueViolation {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserAlreadyExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id.String(), nil
}

func New(storageCredentials string) (*Postgres, error) {
	const op = "storage.postgres.New"
	db, err := sql.Open("postgres", storageCredentials)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Postgres{db: db}, err
}

func (p *Postgres) Close() error {
	const op = "storage.postgres.Close"
	err := p.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
