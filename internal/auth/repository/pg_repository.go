package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
	"url-shortener/internal/auth"
	"url-shortener/internal/models"
)

var (
	UserNotFound = errors.New("user not found")
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepository{db: db}
}

func (r *authRepository) Create(ctx context.Context, user *models.User) (uuid.UUID, error) {
	res := r.db.QueryRowxContext(ctx, createUser, user.Email, user.Password)
	var uid uuid.UUID
	if err := res.Scan(&uid); err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}

func (r *authRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.GetContext(ctx, &user, getUserByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.GetContext(ctx, &user, getUserByEmail, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error {
	_, err := r.db.ExecContext(ctx, updateLastLogin, loginTime, id)
	return err
}
