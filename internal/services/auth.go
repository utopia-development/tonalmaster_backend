package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/utopia-development/tonalmaster_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct { repo repository.AuthRepository }

func NewAuthService(repo repository.AuthRepository) *AuthService { return &AuthService{repo: repo} }

func (s *AuthService) Register(ctx context.Context,email,username,password string)(repository.User,string,error){
	email=strings.TrimSpace(email); username=strings.TrimSpace(username)
	if email==""||username==""||len(password)<8{return repository.User{},"",ErrInvalidCredentials}
	hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost); if err!=nil{return repository.User{},"",err}
	user,err:=s.repo.CreateUser(ctx,email,username,string(hash),"reader"); if err!=nil{return repository.User{},"",err}
	token,err:=newToken(); if err!=nil{return repository.User{},"",err}
	if err:=s.repo.CreateSession(ctx,user.ID,hashToken(token),time.Now().Add(30*24*time.Hour));err!=nil{return repository.User{},"",err}
	return user,token,nil
}

func (s *AuthService) Login(ctx context.Context,email,password string)(repository.User,string,error){
	user,err:=s.repo.FindUserByEmail(ctx,strings.TrimSpace(email)); if err!=nil{return repository.User{},"",ErrInvalidCredentials}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(password))!=nil{return repository.User{},"",ErrInvalidCredentials}
	token,err:=newToken();if err!=nil{return repository.User{},"",err}
	if err:=s.repo.CreateSession(ctx,user.ID,hashToken(token),time.Now().Add(30*24*time.Hour));err!=nil{return repository.User{},"",err}
	return user,token,nil
}

func (s *AuthService) CurrentUser(ctx context.Context,token string)(repository.User,error){
	if token==""{return repository.User{},ErrInvalidCredentials}
	user,_,err:=s.repo.FindSessionUser(ctx,hashToken(token));if err!=nil{return repository.User{},ErrInvalidCredentials}
	return user,nil
}

func (s *AuthService) Logout(ctx context.Context,token string)error{
	if token==""{return nil}
	return s.repo.RevokeSession(ctx,hashToken(token))
}

func newToken()(string,error){b:=make([]byte,32);if _,err:=rand.Read(b);err!=nil{return "",err};return base64.RawURLEncoding.EncodeToString(b),nil}
func hashToken(token string)[]byte{sum:=sha256.Sum256([]byte(token));return sum[:]}
var _ = pgtype.UUID{}
