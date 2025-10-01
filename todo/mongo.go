package todo

//
import (
	"context"
	"fmt"
	"time"
	//
	// 	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AddTodoToDB(col *mongo.Collection, td Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Printf("%v", td)
	_, err := col.InsertOne(ctx, td)
	if err != nil {
		return err
	}
	return nil
}
