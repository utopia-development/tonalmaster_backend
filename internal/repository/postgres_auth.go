package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type PostgresAuthRepository struct { db *pgxpool.Pool }

func NewPostgresAuthRepository(db *pgxpool.Pool) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) CreateUser(ctx context.Context, email, username, passwordHash, role string) (User, error) {
	var u User
	err := r.db.QueryRow(ctx, "INSERT INTO users (email, username, password_hash, role) VALUES ($1,$2,$3,$4) RETURNING id,email,username,password_hash,role", email, username, passwordHash, role).Scan(&u.ID,&u.Email,&u.Username,&u.PasswordHash,&u.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { return User{}, ErrUserAlreadyExists }
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return u,nil
}

func (r *PostgresAuthRepository) FindUserByEmail(ctx context.Context, email string) (User,error) {
	var u User
	err:=r.db.QueryRow(ctx,"SELECT id,email,username,password_hash,role FROM users WHERE LOWER(email)=LOWER($1)",email).Scan(&u.ID,&u.Email,&u.Username,&u.PasswordHash,&u.Role)
	if errors.Is(err,pgx.ErrNoRows){return User{},ErrNotFound}
	if err!=nil{return User{},fmt.Errorf("find user: %w",err)}
	return u,nil
}

func (r *PostgresAuthRepository) CreateSession(ctx context.Context,userID pgtype.UUID,tokenHash []byte,expiresAt time.Time) error {
	_,err:=r.db.Exec(ctx,"INSERT INTO sessions (user_id,token_hash,expires_at) VALUES ($1,$2,$3)",userID,tokenHash,expiresAt)
	return err
}

func (r *PostgresAuthRepository) FindSessionUser(ctx context.Context,tokenHash []byte)(User,Session,error){
	var u User
	var s Session
	err:=r.db.QueryRow(ctx,"SELECT s.id,s.user_id,s.expires_at,u.id,u.email,u.username,u.password_hash,u.role FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>CURRENT_TIMESTAMP",tokenHash).Scan(&s.ID,&s.UserID,&s.ExpiresAt,&u.ID,&u.Email,&u.Username,&u.PasswordHash,&u.Role)
	if errors.Is(err,pgx.ErrNoRows){return User{},Session{},ErrNotFound}
	if err!=nil{return User{},Session{},fmt.Errorf("find session: %w",err)}
	return u,s,nil
}

func (r *PostgresAuthRepository) RevokeSession(ctx context.Context,tokenHash []byte) error {
	_,err:=r.db.Exec(ctx,"UPDATE sessions SET revoked_at=CURRENT_TIMESTAMP WHERE token_hash=$1 AND revoked_at IS NULL",tokenHash)
	return err
}
