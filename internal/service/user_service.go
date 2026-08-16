package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Name  *string
	Email *string
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Create(
	ctx context.Context,
	input CreateUserInput,
) (*model.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	passwordLength := len([]byte(input.Password))

	if name == "" ||
		email == "" ||
		passwordLength < 8 ||
		passwordLength > 72 {
		return nil, ErrInvalidInput
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByID(
	ctx context.Context,
	id string,
) (*model.User, error) {
	objectID, err := parseUserID(id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return user, nil
}

func (s *UserService) List(
	ctx context.Context,
) ([]model.User, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}

func (s *UserService) Update(
	ctx context.Context,
	id string,
	input UpdateUserInput,
) (*model.User, error) {
	objectID, err := parseUserID(id)
	if err != nil {
		return nil, err
	}

	update := model.UserUpdate{}

	//NOTE - case update name
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrInvalidInput
		}

		update.Name = &name
	}

	//NOTE - case update email
	if input.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*input.Email))
		if email == "" {
			return nil, ErrInvalidInput
		}

		update.Email = &email
	}

	//NOTE - ถ้าไม่มีการอัปเดตทั้ง name และ email ErrInvalidInput
	if update.Name == nil && update.Email == nil {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepository.Update(
		ctx,
		objectID,
		update,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (s *UserService) Delete(
	ctx context.Context,
	id string,
) error {
	objectID, err := parseUserID(id)
	if err != nil {
		return err
	}

	if err := s.userRepository.Delete(ctx, objectID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (s *UserService) Count(
	ctx context.Context,
) (int64, error) {
	count, err := s.userRepository.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

// NOTE - parseUserID แปลง string เป็น bson.ObjectID
func parseUserID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrInvalidUserID
	}

	return objectID, nil
}
