package domain

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/entity"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	logger     *logrus.Logger
	collection *mongo.Collection
	cache      *redis.Client
}

type UserInterface interface {
	List(ctx context.Context) ([]entity.User, error)
	Get(ctx context.Context, filter entity.User) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, user entity.User) error
}

// initUser creates user domain
func initUser(logger *logrus.Logger, db *mongo.Collection, cache *redis.Client) UserInterface {
	return &user{
		logger:     logger,
		collection: db,
		cache:      cache,
	}
}

// List returns list of users
func (u *user) List(ctx context.Context) ([]entity.User, error) {
	users := []entity.User{}
	cursor, err := u.collection.Find(ctx, bson.D{})
	if err != nil {
		return users, errorAlias(err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return users, errorAlias(err)
	}

	return users, nil
}

// Get returns specific user by email
func (u *user) Get(ctx context.Context, req entity.User) (entity.User, error) {
	u.logger.Debug(req)
	user := entity.User{}
	var filter any

	switch {
	case req.Email != "":
		filter = bson.M{"email": req.Email}
	case req.Id.String() != "":
		filter = bson.M{"_id": req.Id}
	case req.Name != "":
		filter = bson.M{"_id": req.Id}
	}

	// get from cache, if no error and user found, direct return
	user, err := u.getUserCache(ctx, req.Id.Hex())
	if err == nil && !user.Id.IsZero() {
		u.logger.Info(fmt.Sprintf("cache for user:%v found", req.Id.Hex()))
		return user, nil
	}

	if err := u.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, errorAlias(err)
	}

	// set user cache if result found from mongo
	if err := u.setUserCache(ctx, user); err != nil {
		u.logger.Error(fmt.Sprintf("cache for user:%v failed to be set", req.Id.Hex()))
	}

	return user, nil
}

// Create creates new data
func (u *user) Create(ctx context.Context, user entity.User) (entity.User, error) {
	res, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return user, errorAlias(err)
	}

	newUser, err := u.Get(ctx, entity.User{Id: res.InsertedID.(primitive.ObjectID)})
	if err != nil {
		return newUser, errorAlias(err)
	}

	return newUser, nil
}

// Update updates existing data
func (u *user) Update(ctx context.Context, user entity.User) (entity.User, error) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": user}

	_, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return user, errorAlias(err)
	}

	newUser, err := u.Get(ctx, entity.User{Id: user.Id})
	if err != nil {
		return newUser, errorAlias(err)
	}

	return newUser, nil
}

// Delete deletes existing data
func (u *user) Delete(ctx context.Context, user entity.User) error {
	filter := bson.M{"_id": user.Id}

	res, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errorAlias(err)
	}

	if res.DeletedCount < 1 {
		return errors.ErrNotFound
	}

	return nil
}
