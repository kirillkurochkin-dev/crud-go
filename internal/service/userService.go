package service

import (
	"context"
	"crud-go/internal/entity"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByCredentials(ctx context.Context, email, password string) (entity.User, error)
}

type User struct {
	userRepository UsersRepository
	hasher         PasswordHasher

	hmacSecret []byte
	tokenTtl   time.Duration
}

func NewUser(userRepository UsersRepository, hasher PasswordHasher, hmacSecret []byte, tokenTtl time.Duration) *User {
	return &User{userRepository: userRepository, hasher: hasher, hmacSecret: hmacSecret, tokenTtl: tokenTtl}

}

func (u *User) SignUp(ctx context.Context, input entity.SignUpInput) error {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := entity.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.userRepository.Create(ctx, user)
}

func (u *User) SignIn(ctx context.Context, input entity.SignInInput) (string, error) {
	password, err := u.hasher.Hash(input.Password)
	if err != nil {
		return " ", err
	}

	user, err := u.userRepository.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		return " ", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	sig, _ := token.SignedString(u.hmacSecret)
	return sig, nil
}

func (s *User) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}
