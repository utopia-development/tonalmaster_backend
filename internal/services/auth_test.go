package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/utopia-development/tonalmaster_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	users    map[string]repository.User
	sessions map[string]struct {
		userID    pgtype.UUID
		expiresAt time.Time
		revoked   bool
	}
}

func newMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{
		users: make(map[string]repository.User),
		sessions: make(map[string]struct {
			userID    pgtype.UUID
			expiresAt time.Time
			revoked   bool
		}),
	}
}

func (m *mockAuthRepo) CreateUser(_ context.Context, email, username, passwordHash, role string) (repository.User, error) {
	if _, ok := m.users[email]; ok {
		return repository.User{}, errors.New("duplicate")
	}
	var id pgtype.UUID
	_ = id.Scan("11111111-1111-1111-1111-111111111111")
	u := repository.User{ID: id, Email: email, Username: username, PasswordHash: passwordHash, Role: role}
	m.users[email] = u
	return u, nil
}

func (m *mockAuthRepo) FindUserByEmail(_ context.Context, email string) (repository.User, error) {
	u, ok := m.users[email]
	if !ok {
		return repository.User{}, repository.ErrNotFound
	}
	return u, nil
}

func (m *mockAuthRepo) CreateSession(_ context.Context, userID pgtype.UUID, tokenHash []byte, expiresAt time.Time) error {
	m.sessions[string(tokenHash)] = struct {
		userID    pgtype.UUID
		expiresAt time.Time
		revoked   bool
	}{userID: userID, expiresAt: expiresAt}
	return nil
}

func (m *mockAuthRepo) FindSessionUser(_ context.Context, tokenHash []byte) (repository.User, repository.Session, error) {
	s, ok := m.sessions[string(tokenHash)]
	if !ok || s.revoked || time.Now().After(s.expiresAt) {
		return repository.User{}, repository.Session{}, repository.ErrNotFound
	}
	for _, u := range m.users {
		if u.ID == s.userID {
			return u, repository.Session{UserID: s.userID, ExpiresAt: s.expiresAt}, nil
		}
	}
	return repository.User{}, repository.Session{}, repository.ErrNotFound
}

func (m *mockAuthRepo) RevokeSession(_ context.Context, tokenHash []byte) error {
	s, ok := m.sessions[string(tokenHash)]
	if !ok {
		return nil
	}
	s.revoked = true
	m.sessions[string(tokenHash)] = s
	return nil
}

func TestAuthRegisterLoginMeLogout(t *testing.T) {
	repo := newMockAuthRepo()
	svc := NewAuthService(repo)
	ctx := context.Background()

	user, token, err := svc.Register(ctx, "a@example.com", "alice", "password12", "test-code", "test-code")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.Email != "a@example.com" || user.Username != "alice" || user.Role != "reader" {
		t.Fatalf("user = %+v", user)
	}
	if token == "" {
		t.Fatal("expected session token")
	}

	me, err := svc.CurrentUser(ctx, token)
	if err != nil || me.Email != "a@example.com" {
		t.Fatalf("me after register: user=%+v err=%v", me, err)
	}

	if _, _, err := svc.Login(ctx, "a@example.com", "wrong-password"); !errors.Is(err, ErrInvalidRegistrationData) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}

	user2, token2, err := svc.Login(ctx, "a@example.com", "password12")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if user2.Username != "alice" || token2 == "" {
		t.Fatalf("login result user=%+v token empty=%v", user2, token2 == "")
	}

	if err := svc.Logout(ctx, token2); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := svc.CurrentUser(ctx, token2); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected unauthorized after logout, got %v", err)
	}
}

func TestAuthRegisterValidation(t *testing.T) {
	svc := NewAuthService(newMockAuthRepo())
	ctx := context.Background()
	if _, _, err := svc.Register(ctx, "", "x", "password12", "test-code", "test-code"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("empty email: %v", err)
	}
	if _, _, err := svc.Register(ctx, "a@b.c", "x", "short", "test-code", "test-code"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("short password: %v", err)
	}
}

func TestPasswordIsHashed(t *testing.T) {
	repo := newMockAuthRepo()
	svc := NewAuthService(repo)
	_, _, err := svc.Register(context.Background(), "b@example.com", "bob", "password12", "test-code", "test-code")
	if err != nil {
		t.Fatal(err)
	}
	u := repo.users["b@example.com"]
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("password12")); err != nil {
		t.Fatalf("stored hash does not match password: %v", err)
	}
	if u.PasswordHash == "password12" {
		t.Fatal("password stored in plaintext")
	}
}
