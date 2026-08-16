package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// SECTION - AuthService.Register()
func TestAuthServiceRegister(t *testing.T) {
	mockRepository := &mockUserRepository{
		createFunc: func(
			ctx context.Context,
			user *model.User,
		) error {
			user.ID = bson.NewObjectID()
			return nil
		},
	}

	userService := NewUserService(mockRepository)

	authService := NewAuthService(
		mockRepository,
		userService,
		nil,
	)

	user, err := authService.Register(
		context.Background(),
		RegisterInput{
			Name:     "  Palm  ",
			Email:    "  PALM@EXAMPLE.COM  ",
			Password: "password123",
		},
	)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.ID.IsZero() {
		t.Fatal("Register() user ID must not be zero")
	}

	if user.Name != "Palm" {
		t.Errorf(
			"Register() Name = %q, want %q",
			user.Name,
			"Palm",
		)
	}

	if user.Email != "palm@example.com" {
		t.Errorf(
			"Register() Email = %q, want %q",
			user.Email,
			"palm@example.com",
		)
	}

	if user.Password == "password123" {
		t.Fatal("Register() password must be hashed")
	}

	if !CheckPassword(user.Password, "password123") {
		t.Fatal("Register() password hash is invalid")
	}
}

func TestAuthServiceRegisterDuplicateEmail(t *testing.T) {
	mockRepository := &mockUserRepository{
		createFunc: func(
			ctx context.Context,
			user *model.User,
		) error {
			return repository.ErrEmailAlreadyExists
		},
	}

	userService := NewUserService(mockRepository)

	authService := NewAuthService(
		mockRepository,
		userService,
		nil,
	)

	user, err := authService.Register(
		context.Background(),
		RegisterInput{
			Name:     "Palm",
			Email:    "palm@example.com",
			Password: "password123",
		},
	)

	if user != nil {
		t.Errorf("Register() user = %v, want nil", user)
	}

	if !errors.Is(err, repository.ErrEmailAlreadyExists) {
		t.Fatalf(
			"Register() error = %v, want ErrEmailAlreadyExists",
			err,
		)
	}
}

//!SECTION

// SECTION - AuthService.Login()
func TestAuthServiceLogin(t *testing.T) {
	id := bson.NewObjectID()

	hashedPassword, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	expectedUser := &model.User{
		ID:       id,
		Name:     "Palm",
		Email:    "palm@example.com",
		Password: hashedPassword,
	}

	mockRepository := &mockUserRepository{
		findByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			if email != expectedUser.Email {
				t.Errorf(
					"FindByEmail() email = %q, want %q",
					email,
					expectedUser.Email,
				)
			}

			return expectedUser, nil
		},
	}

	tokenService := NewTokenService(
		"test-secret",
		time.Hour,
	)

	authService := NewAuthService(
		mockRepository,
		nil,
		tokenService,
	)

	token, err := authService.Login(
		context.Background(),
		LoginInput{
			Email:    "  PALM@EXAMPLE.COM  ",
			Password: "password123",
		},
	)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if token == "" {
		t.Fatal("Login() token must not be empty")
	}

	claims, err := tokenService.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if claims.Subject != id.Hex() {
		t.Errorf(
			"token subject = %q, want %q",
			claims.Subject,
			id.Hex(),
		)
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
	hashedPassword, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	mockRepository := &mockUserRepository{
		findByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return &model.User{
				ID:       bson.NewObjectID(),
				Email:    "palm@example.com",
				Password: hashedPassword,
			}, nil
		},
	}

	authService := NewAuthService(
		mockRepository,
		nil,
		NewTokenService("test-secret", time.Hour),
	)

	token, err := authService.Login(
		context.Background(),
		LoginInput{
			Email:    "palm@example.com",
			Password: "wrong-password",
		},
	)

	if token != "" {
		t.Errorf("Login() token = %q, want empty", token)
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"Login() error = %v, want ErrInvalidCredentials",
			err,
		)
	}
}
func TestAuthServiceLoginUserNotFound(t *testing.T) {
	mockRepository := &mockUserRepository{
		findByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}

	authService := NewAuthService(
		mockRepository,
		nil,
		NewTokenService("test-secret", time.Hour),
	)

	token, err := authService.Login(
		context.Background(),
		LoginInput{
			Email:    "unknown@example.com",
			Password: "password123",
		},
	)

	if token != "" {
		t.Errorf("Login() token = %q, want empty", token)
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"Login() error = %v, want ErrInvalidCredentials",
			err,
		)
	}
}
func TestAuthServiceLoginInvalidInput(t *testing.T) {
	authService := NewAuthService(
		&mockUserRepository{},
		nil,
		NewTokenService("test-secret", time.Hour),
	)

	tests := []struct {
		name  string
		input LoginInput
	}{
		{
			name: "empty email",
			input: LoginInput{
				Password: "password123",
			},
		},
		{
			name: "empty password",
			input: LoginInput{
				Email: "palm@example.com",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := authService.Login(
				context.Background(),
				test.input,
			)

			if token != "" {
				t.Errorf(
					"Login() token = %q, want empty",
					token,
				)
			}

			if !errors.Is(err, ErrInvalidCredentials) {
				t.Fatalf(
					"Login() error = %v, want ErrInvalidCredentials",
					err,
				)
			}
		})
	}
}

//!SECTION
