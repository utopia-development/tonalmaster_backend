package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID pgtype.UUID
	Email string
	Username string
	PasswordHash string
	Role string
}

type Session struct {
	ID pgtype.UUID
	UserID pgtype.UUID
	ExpiresAt time.Time
}

type AuthRepository interface {
	CreateUser(ctx context.Context, email, username, passwordHash, role string) (User, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	CreateSession(ctx context.Context, userID pgtype.UUID, tokenHash []byte, expiresAt time.Time) error
	FindSessionUser(ctx context.Context, tokenHash []byte) (User, Session, error)
	RevokeSession(ctx context.Context, tokenHash []byte) error
}
