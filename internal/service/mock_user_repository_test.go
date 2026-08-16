package service

import (
	"context"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockUserRepository struct {
	createFunc func(
		ctx context.Context,
		user *model.User,
	) error

	findByIDFunc func(
		ctx context.Context,
		id bson.ObjectID,
	) (*model.User, error)

	findByEmailFunc func(
		ctx context.Context,
		email string,
	) (*model.User, error)

	findAllFunc func(
		ctx context.Context,
	) ([]model.User, error)

	updateFunc func(
		ctx context.Context,
		id bson.ObjectID,
		update model.UserUpdate,
	) (*model.User, error)

	deleteFunc func(
		ctx context.Context,
		id bson.ObjectID,
	) error

	countFunc func(
		ctx context.Context,
	) (int64, error)
}

func (m *mockUserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	if m.createFunc == nil {
		panic("unexpected call to Create")
	}

	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.User, error) {
	if m.findByIDFunc == nil {
		panic("unexpected call to FindByID")
	}

	return m.findByIDFunc(ctx, id)
}

func (m *mockUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	if m.findByEmailFunc == nil {
		panic("unexpected call to FindByEmail")
	}

	return m.findByEmailFunc(ctx, email)
}

func (m *mockUserRepository) FindAll(
	ctx context.Context,
) ([]model.User, error) {
	if m.findAllFunc == nil {
		panic("unexpected call to FindAll")
	}

	return m.findAllFunc(ctx)
}

func (m *mockUserRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	update model.UserUpdate,
) (*model.User, error) {
	if m.updateFunc == nil {
		panic("unexpected call to Update")
	}

	return m.updateFunc(ctx, id, update)
}

func (m *mockUserRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	if m.deleteFunc == nil {
		panic("unexpected call to Delete")
	}

	return m.deleteFunc(ctx, id)
}

func (m *mockUserRepository) Count(
	ctx context.Context,
) (int64, error) {
	if m.countFunc == nil {
		panic("unexpected call to Count")
	}

	return m.countFunc(ctx)
}
