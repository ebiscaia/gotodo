package todo

//
import (
	"context"
	"fmt"
	"time"

	"gotodo/user"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func RemoveTodoAtIndex(col *mongo.Collection, usrTodos []bson.ObjectID, index int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, bson.M{"_id": usrTodos[index]})

	if err != nil {
		return err
	}

	return nil
}

func ChangeTodoDB(col *mongo.Collection, usrTodos []bson.ObjectID, index int, td string) error {
	// retrieve todo based on its id
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// _, err := col.DeleteOne(ctx, bson.M{"_id": usrTodos[index]})
	_, err := col.UpdateOne(ctx, bson.M{"_id": usrTodos[index]}, bson.D{{"$set", bson.D{{"td", td}}}})

	if err != nil {
		return err
	}

	return nil
}

func ChangeStatus(col *mongo.Collection, usrTodos []bson.ObjectID, index int, status bool) error {
	// retrieve todo based on its id
	println(status)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// _, err := col.DeleteOne(ctx, bson.M{"_id": usrTodos[index]})
	_, err := col.UpdateOne(ctx, bson.M{"_id": usrTodos[index]}, bson.D{{"$set", bson.D{{"isDone", status}}}})

	if err != nil {
		return err
	}

	return nil
}

func GetTodo(col *mongo.Collection, usrTodos []bson.ObjectID, index int) (Todo, error) {
	// retrieve todo based on its id
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	td := Todo{}
	err := col.FindOne(ctx, bson.M{"_id": usrTodos[index]}).Decode(&td)
	if err != nil {
		return td, err
	}

	return td, nil
}

func setFilter(usr user.User, printAll bool) bson.M {
	filter := bson.M{"user": usr.ID}
	if !printAll {
		filter = bson.M{"user": usr.ID, "isDone": false}
	}
	return filter
}

func CheckHasTodos(col *mongo.Collection, userToLogin user.User, allTodos bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ftr := setFilter(userToLogin, allTodos)

	usrHasTodos, err := col.CountDocuments(ctx, ftr)
	if err != nil {
		return false, err
	}
	if usrHasTodos == 0 {
		return false, nil
	}
	return true, nil
}

func setStatusText(td Todo) string {
	if td.IsDone {
		return "done"
	}
	return "undone"
}

func PrintTodos(col *mongo.Collection, userToLogin user.User, allTodos bool, ind bool) ([]bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// set an slice of ids that will be returned to
	// delete, change and change status functions
	ids := []bson.ObjectID{}

	// set an index variable that will be printed out when
	// the functions mentioned above are called. the index will
	// be the reference of which to do will be deleted/modified
	index := 0

	// prepare the cursor, first setting the options then the filter
	opts := options.Find()
	opts.SetSort(bson.D{{"isDone", 1}})

	ftr := setFilter(userToLogin, allTodos)

	sortedCursor, err := col.Find(ctx, ftr, opts)
	if err != nil {
		return []bson.ObjectID{}, err
	}

	// loop the cursor, create a Todo type of variable
	// append its ID to the ids slice, print accordingly
	// with the conditions (with index or not and showing
	// or not completed tasks)
	for sortedCursor.Next(ctx) {
		td := Todo{}
		err := sortedCursor.Decode(&td)
		if err != nil {
			return ids, err
		}
		ids = append(ids, td.ID)

		// set the text to represent the status instead of showing the boolean
		// value
		statusText := setStatusText(td)

		if ind {
			if allTodos {
				fmt.Printf("%v - %v (%v)\n", index+1, td.Td, statusText)
			} else {
				fmt.Printf("%v - %v\n", index+1, td.Td)
			}
		} else {
			if allTodos {
				fmt.Printf("%v (%v)\n", td.Td, statusText)
			} else {
				fmt.Printf("%v\n", td.Td)
			}
		}
		index += 1
	}

	return ids, nil
}

func FindTodoDBCI(col *mongo.Collection, usr user.User, td string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	todoToFind := Todo{}
	err := col.FindOne(ctx,
		bson.M{"user": usr.ID,
			"isDone": false,
			"td":     bson.M{"$regex": td, "$options": "i"}}).Decode(&todoToFind)
	return err == nil
}

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
