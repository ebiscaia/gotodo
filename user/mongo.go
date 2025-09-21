package user

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

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

