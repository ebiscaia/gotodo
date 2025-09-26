package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// check whether the collection has at least one document
func CheckUsers(col *mongo.Collection) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	hasUsers, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	if hasUsers == 0 {
		return false, nil
	}
	return true, nil
}

func FindUserDB(col *mongo.Collection, usr User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userToFind := User{}
	// This function is for login purposes. Therefore, user name is to be a exact match
	err := col.FindOne(ctx,
		bson.M{"name": usr.Name}).Decode(&userToFind)
	if err != nil {
		return User{}, err
	}
	return userToFind, nil
}

func AddUserToDB(col *mongo.Collection, us User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, us)
	if err != nil {
		return err
	}
	return nil
}
