package service

import (
	"context"
	"errors"
	"testing"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// SECTION - UserService.Create
func TestUserServiceCreate(t *testing.T) {
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

	user, err := userService.Create(
		context.Background(),
		CreateUserInput{
			Name:     "  Palm  ",
			Email:    "  PALM@EXAMPLE.COM  ",
			Password: "password123",
		},
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.ID.IsZero() {
		t.Fatal("Create() user ID must not be zero")
	}

	if user.Name != "Palm" {
		t.Errorf(
			"Create() Name = %q, want %q",
			user.Name,
			"Palm",
		)
	}

	if user.Email != "palm@example.com" {
		t.Errorf(
			"Create() Email = %q, want %q",
			user.Email,
			"palm@example.com",
		)
	}

	if user.Password == "password123" {
		t.Fatal("Create() must hash password")
	}

	if !CheckPassword(user.Password, "password123") {
		t.Fatal("Create() password hash is invalid")
	}

	if user.CreatedAt.IsZero() {
		t.Fatal("Create() CreatedAt must be set")
	}
}

func TestUserServiceCreateDuplicateEmail(t *testing.T) {
	mockRepository := &mockUserRepository{
		createFunc: func(
			ctx context.Context,
			user *model.User,
		) error {
			return repository.ErrEmailAlreadyExists
		},
	}

	userService := NewUserService(mockRepository)

	_, err := userService.Create(
		context.Background(),
		CreateUserInput{
			Name:     "Palm",
			Email:    "palm@example.com",
			Password: "password123",
		},
	)

	if !errors.Is(err, repository.ErrEmailAlreadyExists) {
		t.Fatalf(
			"Create() error = %v, want ErrEmailAlreadyExists",
			err,
		)
	}
}

func TestUserServiceCreateInvalidInput(t *testing.T) {
	userService := NewUserService(&mockUserRepository{})

	tests := []struct {
		name  string
		input CreateUserInput
	}{
		{
			name: "empty name",
			input: CreateUserInput{
				Email:    "palm@example.com",
				Password: "password123",
			},
		},
		{
			name: "empty email",
			input: CreateUserInput{
				Name:     "Palm",
				Password: "password123",
			},
		},
		{
			name: "short password",
			input: CreateUserInput{
				Name:     "Palm",
				Email:    "palm@example.com",
				Password: "123",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := userService.Create(
				context.Background(),
				test.input,
			)

			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf(
					"Create() error = %v, want ErrInvalidInput",
					err,
				)
			}
		})
	}
}

//!SECTION

// SECTION - UserService.GetByID
func TestUserServiceGetByID(t *testing.T) {
	id := bson.NewObjectID()

	expectedUser := &model.User{
		ID:    id,
		Name:  "Palm",
		Email: "palm@example.com",
	}

	mockRepository := &mockUserRepository{
		findByIDFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
		) (*model.User, error) {
			if receivedID != id {
				t.Errorf(
					"FindByID() ID = %s, want %s",
					receivedID.Hex(),
					id.Hex(),
				)
			}

			return expectedUser, nil
		},
	}

	userService := NewUserService(mockRepository)

	user, err := userService.GetByID(
		context.Background(),
		id.Hex(),
	)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if user.ID != id {
		t.Errorf(
			"GetByID() ID = %s, want %s",
			user.ID.Hex(),
			id.Hex(),
		)
	}

	if user.Name != expectedUser.Name {
		t.Errorf(
			"GetByID() Name = %q, want %q",
			user.Name,
			expectedUser.Name,
		)
	}

	if user.Email != expectedUser.Email {
		t.Errorf(
			"GetByID() Email = %q, want %q",
			user.Email,
			expectedUser.Email,
		)
	}
}

func TestUserServiceGetByIDInvalidID(t *testing.T) {
	userService := NewUserService(&mockUserRepository{})

	_, err := userService.GetByID(
		context.Background(),
		"invalid-id",
	)

	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf(
			"GetByID() error = %v, want ErrInvalidUserID",
			err,
		)
	}
}

func TestUserServiceGetByIDNotFound(t *testing.T) {
	id := bson.NewObjectID()

	mockRepository := &mockUserRepository{
		findByIDFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
		) (*model.User, error) {
			if receivedID != id {
				t.Errorf(
					"FindByID() ID = %s, want %s",
					receivedID.Hex(),
					id.Hex(),
				)
			}

			return nil, repository.ErrUserNotFound
		},
	}

	userService := NewUserService(mockRepository)

	_, err := userService.GetByID(
		context.Background(),
		id.Hex(),
	)

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf(
			"GetByID() error = %v, want ErrUserNotFound",
			err,
		)
	}
}

//!SECTION

// SECTION - UserService.List
func TestUserServiceList(t *testing.T) {
	expectedUsers := []model.User{
		{
			ID:    bson.NewObjectID(),
			Name:  "Palm",
			Email: "palm@example.com",
		},
		{
			ID:    bson.NewObjectID(),
			Name:  "Friend",
			Email: "friend@example.com",
		},
	}

	mockRepository := &mockUserRepository{
		findAllFunc: func(
			ctx context.Context,
		) ([]model.User, error) {
			return expectedUsers, nil
		},
	}

	userService := NewUserService(mockRepository)

	users, err := userService.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(users) != len(expectedUsers) {
		t.Fatalf(
			"List() returned %d users, want %d",
			len(users),
			len(expectedUsers),
		)
	}

	if users[0].Email != expectedUsers[0].Email {
		t.Errorf(
			"List() first user email = %q, want %q",
			users[0].Email,
			expectedUsers[0].Email,
		)
	}

	if users[1].Email != expectedUsers[1].Email {
		t.Errorf(
			"List() second user email = %q, want %q",
			users[1].Email,
			expectedUsers[1].Email,
		)
	}
}

func TestUserServiceListRepositoryError(t *testing.T) {
	repositoryError := errors.New("database error")

	mockRepository := &mockUserRepository{
		findAllFunc: func(
			ctx context.Context,
		) ([]model.User, error) {
			return nil, repositoryError
		},
	}

	userService := NewUserService(mockRepository)

	users, err := userService.List(context.Background())

	if users != nil {
		t.Errorf("List() users = %v, want nil", users)
	}

	if !errors.Is(err, repositoryError) {
		t.Fatalf(
			"List() error = %v, want database error",
			err,
		)
	}
}

//!SECTION

// SECTION - UserService.Update
func TestUserServiceUpdate(t *testing.T) {
	id := bson.NewObjectID()
	name := "  Palm Updated  "
	email := "  PALM.NEW@EXAMPLE.COM  "

	expectedUser := &model.User{
		ID:    id,
		Name:  "Palm Updated",
		Email: "palm.new@example.com",
	}

	mockRepository := &mockUserRepository{
		updateFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
			update model.UserUpdate,
		) (*model.User, error) {
			if receivedID != id {
				t.Errorf(
					"Update() ID = %s, want %s",
					receivedID.Hex(),
					id.Hex(),
				)
			}

			if update.Name == nil {
				t.Fatal("Update() Name must not be nil")
			}

			if *update.Name != "Palm Updated" {
				t.Errorf(
					"Update() Name = %q, want %q",
					*update.Name,
					"Palm Updated",
				)
			}

			if update.Email == nil {
				t.Fatal("Update() Email must not be nil")
			}

			if *update.Email != "palm.new@example.com" {
				t.Errorf(
					"Update() Email = %q, want %q",
					*update.Email,
					"palm.new@example.com",
				)
			}

			return expectedUser, nil
		},
	}

	userService := NewUserService(mockRepository)

	user, err := userService.Update(
		context.Background(),
		id.Hex(),
		UpdateUserInput{
			Name:  &name,
			Email: &email,
		},
	)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf(
			"Update() user ID = %s, want %s",
			user.ID.Hex(),
			expectedUser.ID.Hex(),
		)
	}

	if user.Name != expectedUser.Name {
		t.Errorf(
			"Update() user Name = %q, want %q",
			user.Name,
			expectedUser.Name,
		)
	}

	if user.Email != expectedUser.Email {
		t.Errorf(
			"Update() user Email = %q, want %q",
			user.Email,
			expectedUser.Email,
		)
	}
}

func TestUserServiceUpdateInvalidID(t *testing.T) {
	name := "Palm Updated"
	userService := NewUserService(&mockUserRepository{})

	user, err := userService.Update(
		context.Background(),
		"invalid-id",
		UpdateUserInput{
			Name: &name,
		},
	)

	if user != nil {
		t.Errorf("Update() user = %v, want nil", user)
	}

	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf(
			"Update() error = %v, want ErrInvalidUserID",
			err,
		)
	}
}

func TestUserServiceUpdateNotFound(t *testing.T) {
	id := bson.NewObjectID()
	name := "Palm Updated"

	mockRepository := &mockUserRepository{
		updateFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
			update model.UserUpdate,
		) (*model.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}

	userService := NewUserService(mockRepository)

	user, err := userService.Update(
		context.Background(),
		id.Hex(),
		UpdateUserInput{
			Name: &name,
		},
	)

	if user != nil {
		t.Errorf("Update() user = %v, want nil", user)
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf(
			"Update() error = %v, want ErrUserNotFound",
			err,
		)
	}
}

func TestUserServiceUpdateDuplicateEmail(t *testing.T) {
	id := bson.NewObjectID()
	email := "existing@example.com"

	mockRepository := &mockUserRepository{
		updateFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
			update model.UserUpdate,
		) (*model.User, error) {
			return nil, repository.ErrEmailAlreadyExists
		},
	}

	userService := NewUserService(mockRepository)

	user, err := userService.Update(
		context.Background(),
		id.Hex(),
		UpdateUserInput{
			Email: &email,
		},
	)

	if user != nil {
		t.Errorf("Update() user = %v, want nil", user)
	}

	if !errors.Is(err, repository.ErrEmailAlreadyExists) {
		t.Fatalf(
			"Update() error = %v, want ErrEmailAlreadyExists",
			err,
		)
	}
}

func TestUserServiceUpdateEmptyInput(t *testing.T) {
	id := bson.NewObjectID()
	userService := NewUserService(&mockUserRepository{})

	user, err := userService.Update(
		context.Background(),
		id.Hex(),
		UpdateUserInput{},
	)

	if user != nil {
		t.Errorf("Update() user = %v, want nil", user)
	}

	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf(
			"Update() error = %v, want ErrInvalidInput",
			err,
		)
	}
}

//!SECTION

// SECTION - UserService.Delete
func TestUserServiceDelete(t *testing.T) {
	id := bson.NewObjectID()
	deleteCalled := false

	mockRepository := &mockUserRepository{
		deleteFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
		) error {
			deleteCalled = true

			if receivedID != id {
				t.Errorf(
					"Delete() ID = %s, want %s",
					receivedID.Hex(),
					id.Hex(),
				)
			}

			return nil
		},
	}

	userService := NewUserService(mockRepository)

	err := userService.Delete(
		context.Background(),
		id.Hex(),
	)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if !deleteCalled {
		t.Fatal("Delete() repository was not called")
	}
}

func TestUserServiceDeleteInvalidID(t *testing.T) {
	userService := NewUserService(&mockUserRepository{})

	err := userService.Delete(
		context.Background(),
		"invalid-id",
	)

	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf(
			"Delete() error = %v, want ErrInvalidUserID",
			err,
		)
	}
}

func TestUserServiceDeleteNotFound(t *testing.T) {
	id := bson.NewObjectID()

	mockRepository := &mockUserRepository{
		deleteFunc: func(
			ctx context.Context,
			receivedID bson.ObjectID,
		) error {
			return repository.ErrUserNotFound
		},
	}

	userService := NewUserService(mockRepository)

	err := userService.Delete(
		context.Background(),
		id.Hex(),
	)

	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf(
			"Delete() error = %v, want ErrUserNotFound",
			err,
		)
	}
}

//!SECTION

// SECTION - UserService.Count
func TestUserServiceCount(t *testing.T) {
	mockRepository := &mockUserRepository{
		countFunc: func(
			ctx context.Context,
		) (int64, error) {
			return 5, nil
		},
	}

	userService := NewUserService(mockRepository)

	count, err := userService.Count(context.Background())
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 5 {
		t.Errorf(
			"Count() count = %d, want %d",
			count,
			5,
		)
	}
}

func TestUserServiceCountRepositoryError(t *testing.T) {
	repositoryError := errors.New("database error")

	mockRepository := &mockUserRepository{
		countFunc: func(
			ctx context.Context,
		) (int64, error) {
			return 0, repositoryError
		},
	}

	userService := NewUserService(mockRepository)

	count, err := userService.Count(context.Background())

	if count != 0 {
		t.Errorf(
			"Count() count = %d, want 0",
			count,
		)
	}

	if !errors.Is(err, repositoryError) {
		t.Fatalf(
			"Count() error = %v, want database error",
			err,
		)
	}
}

//!SECTION
