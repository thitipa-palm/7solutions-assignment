package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/thitipa-palm/7solutions-assignment/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const usersCollection = "users"

type MongoUserRepository struct {
	collection *mongo.Collection
}

var _ UserRepository = (*MongoUserRepository)(nil)

func NewMongoUserRepository(
	database *mongo.Database,
) *MongoUserRepository {
	return &MongoUserRepository{
		collection: database.Collection(usersCollection),
	}
}

func (r *MongoUserRepository) EnsureIndexes(
	ctx context.Context,
) error {
	index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetName("unique_email"),
	}

	if _, err := r.collection.Indexes().CreateOne(ctx, index); err != nil {
		return fmt.Errorf("create email index: %w", err)
	}

	return nil
}

func (r *MongoUserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrEmailAlreadyExists
		}

		return fmt.Errorf("insert user: %w", err)
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("unexpected inserted ID type")
	}

	user.ID = id

	return nil
}

func (r *MongoUserRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.User, error) {
	var user model.User

	err := r.collection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
	).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find user by ID: %w", err)
	}

	return &user, nil
}

func (r *MongoUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	var user model.User

	err := r.collection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
	).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

func (r *MongoUserRepository) FindAll(
	ctx context.Context,
) ([]model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}
	defer cursor.Close(ctx)

	users := make([]model.User, 0)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}

	return users, nil
}

func (r *MongoUserRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	update model.UserUpdate,
) (*model.User, error) {
	fields := bson.D{}

	if update.Name != nil {
		fields = append(fields, bson.E{
			Key:   "name",
			Value: *update.Name,
		})
	}

	if update.Email != nil {
		fields = append(fields, bson.E{
			Key:   "email",
			Value: *update.Email,
		})
	}

	if len(fields) == 0 {
		return r.FindByID(ctx, id)
	}

	filter := bson.D{{Key: "_id", Value: id}}
	updateDocument := bson.D{
		{
			Key:   "$set",
			Value: fields,
		},
	}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updatedUser model.User

	err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		updateDocument,
		opts,
	).Decode(&updatedUser)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}

	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrEmailAlreadyExists
	}

	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &updatedUser, nil
}

func (r *MongoUserRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	result, err := r.collection.DeleteOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
	)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *MongoUserRepository) Count(
	ctx context.Context,
) (int64, error) {
	opts := options.Count().SetHint("_id_")

	count, err := r.collection.CountDocuments(
		ctx,
		bson.D{},
		opts,
	)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}
