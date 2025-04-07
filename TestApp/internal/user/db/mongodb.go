package db

import (
	apperror "TestApp/internal/apperror"
	"TestApp/internal/user"
	"TestApp/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user *user.User) (string, error) {
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user %v", err)
	}
	oid, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted id to object id %v", err)

	}
	d.logger.Tracef("inserted user with id %s", oid.String())
	return oid.Hex(), nil
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to convert string to object id %v", err)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		//TODO  404
		return u, fmt.Errorf("failed to find user %v", result.Err())
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user %v", err)
	}

	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	result, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users %v", err)
	}
	if result.Err() != nil {
		//TODO  404
		return u, fmt.Errorf("failed to find users %v", result.Err())
	}
	if err = result.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to decode users %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	//oid, err := bson.ObjectIDFromHex(user.ID)
	//if err != nil {
	//	return fmt.Errorf("failed to convert string to object id %v", err)
	//}
	filter := bson.M{"_id": user.ID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user %v", err)
	}

	var updateUserObj bson.M
	if err = bson.Unmarshal(userBytes, &updateUserObj); err != nil {
		return fmt.Errorf("failed to unmarshal user %v", err)
	}
	delete(updateUserObj, "_id")

	update := bson.M{"$set": updateUserObj}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("matched %d documents, updated %d docu,ents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert string to object id %v", err)
	}
	filter := bson.M{"_id": oid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user %v", err)
	}

	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
