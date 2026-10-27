package usecase

import (
	"context"
	"errors"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func (uc *UseCase) Login(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	user, err := uc.users.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	if !uc.hasher.Compare(in.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := uc.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
	}, nil
}

func (uc *UseCase) Register(ctx context.Context, in RegisterInput) (*RegisterOutput, error) {
	passwordHash, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(in.Username, in.Email, domain.RoleUser, passwordHash)
	createdUser, err := uc.users.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := uc.jwt.GenerateToken(createdUser.ID, createdUser.Role)
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		Token: token,
	}, nil
}
