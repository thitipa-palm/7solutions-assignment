package repository

import (
	"context"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error

	FindByID(
		ctx context.Context,
		id bson.ObjectID,
	) (*model.User, error)

	FindByEmail(
		ctx context.Context,
		email string,
	) (*model.User, error)

	FindAll(ctx context.Context) ([]model.User, error)

	Update(
		ctx context.Context,
		id bson.ObjectID,
		update model.UserUpdate,
	) (*model.User, error)

	Delete(
		ctx context.Context,
		id bson.ObjectID,
	) error

	Count(ctx context.Context) (int64, error)
}
