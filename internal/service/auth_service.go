package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService struct {
	userRepository repository.UserRepository
	userService    *UserService
	tokenService   *TokenService
}

func NewAuthService(
	userRepository repository.UserRepository,
	userService *UserService,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		userService:    userService,
		tokenService:   tokenService,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	input RegisterInput,
) (*model.User, error) {
	return s.userService.Create(
		ctx,
		CreateUserInput{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		},
	)
}

func (s *AuthService) Login(
	ctx context.Context,
	input LoginInput,
) (string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if email == "" || input.Password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", fmt.Errorf("find user for login: %w", err)
	}

	if !CheckPassword(user.Password, input.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenService.Generate(user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
