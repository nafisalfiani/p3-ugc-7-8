package domain

import (
	"context"

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
}

type UserInterface interface {
	List(ctx context.Context) ([]entity.User, error)
	Get(ctx context.Context, filter entity.User) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, user entity.User) error
}

// initUser creates user domain
func initUser(logger *logrus.Logger, db *mongo.Collection) UserInterface {
	return &user{
		logger:     logger,
		collection: db,
	}
}

// List returns list of users
func (s *user) List(ctx context.Context) ([]entity.User, error) {
	users := []entity.User{}
	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return users, errorAlias(err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil {
		return users, errorAlias(err)
	}

	return users, nil
}

// Get returns specific user by email
func (s *user) Get(ctx context.Context, req entity.User) (entity.User, error) {
	s.logger.Debug(req)
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

	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, errorAlias(err)
	}

	return user, nil
}

// Create creates new data
func (s *user) Create(ctx context.Context, user entity.User) (entity.User, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return user, errorAlias(err)
	}

	newUser, err := s.Get(ctx, entity.User{Id: res.InsertedID.(primitive.ObjectID)})
	if err != nil {
		return newUser, errorAlias(err)
	}

	return newUser, nil
}

// Update updates existing data
func (s *user) Update(ctx context.Context, user entity.User) (entity.User, error) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": user}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return user, errorAlias(err)
	}

	newUser, err := s.Get(ctx, entity.User{Id: user.Id})
	if err != nil {
		return newUser, errorAlias(err)
	}

	return newUser, nil
}

// Delete deletes existing data
func (s *user) Delete(ctx context.Context, user entity.User) error {
	filter := bson.M{"_id": user.Id}

	res, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errorAlias(err)
	}

	if res.DeletedCount < 1 {
		return errors.ErrNotFound
	}

	return nil
}
