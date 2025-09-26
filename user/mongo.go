package user

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// check whether the collection has at least one document
func CheckUsers(ctx context.Context, col *mongo.Collection) (bool, error) {
	hasUsers, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	if hasUsers == 0 {
		return false, nil
	}
	return true, nil
}

func FindUserDB(ctx context.Context, col *mongo.Collection, usr User) (User, error) {
	userToFind := User{}
	// This function is for login purposes. Therefore, user name is to be a exact match
	err := col.FindOne(ctx,
		bson.M{"name": usr.Name}).Decode(&userToFind)
	if err != nil {
		return User{}, err
	}
	return userToFind, nil
}

func AddUserToDB(ctx context.Context, col *mongo.Collection, us User) error {
	_, err := col.InsertOne(ctx, us)
	if err != nil {
		return err
	}
	return nil
}
