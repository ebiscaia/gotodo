package todo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"gotodo/user"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Todo struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	User   bson.ObjectID `bson:"user"`
	Td     string        `bson:"td"`
	IsDone bool          `bson:"isDone"`
}

func DisplayTodos(col *mongo.Collection, userToLogin user.User) []bson.ObjectID {

	hasTodos, err := CheckHasTodos(col, userToLogin, allTodos)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	if !hasTodos {
		fmt.Printf("User %v does not have any todos\n", userToLogin.Name)
		return []bson.ObjectID{}
	}

	//print and return ids as a slice
	ids, err := PrintTodos(col, userToLogin, allTodos, index)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return ids
}

func CreateTodo(col *mongo.Collection, usrLogin user.User) {
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter new todo:")

	if scn.Scan() {
		// check for duplicates (case insensitive)
		if FindTodoDBCI(col, usrLogin, scn.Text()) {
			fmt.Println("There is already a pending todo with the same name")
			return
		}

		// add todo to the database
		err := AddTodoToDB(col, Todo{Td: scn.Text(), User: usrLogin.ID})
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return
	}
	fmt.Println("There was an error with todo creation. Leaving...")
	os.Exit(1)

}

func inputIndex(lenTodo int, funcParent string) int {
	scn := bufio.NewScanner(os.Stdin)
	message := ""
	switch funcParent {
	case "delete":
		message = "Enter index of todo to delete: "
	case "change":
		message = "Enter index of todo to be changed: "
	case "done":
		message = "Enter index of todo to have status changed: "
	}
	for {
		fmt.Println(message)
		if !scn.Scan() {
			fmt.Println("There is an internal error. Leaving...")
			os.Exit(1)
		}
		index, err := strconv.Atoi(scn.Text())
		if err != nil {
			fmt.Printf("The following error has occured: %v\n", err)
			fmt.Println("Leaving...")
			os.Exit(1)
		}
		if index < 0 || index > lenTodo {
			fmt.Println("Index is out of range. Please try again.")
			continue
		}
		index--
		return index
	}
}

func DeleteTodo(col *mongo.Collection, usrLogin user.User) {
	todosUsr := DisplayTodos(col, usrLogin, true, true)

	// Display zero as option to return to previous menu
	fmt.Println("0 - Return to previous menu")

	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "delete")

	// return to previous menu
	if index == -1 {
		return
	}

	// or display which todo is being deleted
	todo, err := GetTodo(col, todosUsr, index)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task deleted: %v\n", todo.Td)

	err = RemoveTodoAtIndex(col, todosUsr, index)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func promptChangeTodo() string {
	td := ""
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println("Modify todo: ")
	if scn.Scan() {
		td = scn.Text()
	}
	return td
}

func ChangeTodo(col *mongo.Collection, usrLogin user.User) {
	todosUsr := DisplayTodos(col, usrLogin, true, true)

	// Display zero as option to return to previous menu
	fmt.Println("0 - Return to previous menu")

	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "change")

	// return to previous menu
	if index == -1 {
		return
	}

	// or bring the todo at index
	todo, err := GetTodo(col, todosUsr, index)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task chosen: %v\n", todo.Td)
	newTodo := promptChangeTodo()

	// update the task with its new name
	err = ChangeTodoDB(col, todosUsr, index, newTodo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}

func promptChangeStatus() bool {
	scn := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Would you like to change it (y/n): ")
		if scn.Scan() {
			option := scn.Text()
			if option != "y" && option != "n" {
				fmt.Println("Please choose a proper option")
				continue
			}
			if option == "y" {
				return true
			}
			if option == "n" {
				return false
			}
			fmt.Println("There was an error with todo creation. Leaving...")
			os.Exit(1)
		}
	}
}

func ChangeStatusTodo(col *mongo.Collection, usrLogin user.User) {
	todosUsr := DisplayTodos(col, usrLogin, true, true)

	// Display zero as option to return to previous menu
	fmt.Println("0 - Return to previous menu")

	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "done")

	// return to previous menu
	if index == -1 {
		return
	}

	// or bring the todo at index
	todo, err := GetTodo(col, todosUsr, index)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	statusText := "undone"
	if todo.IsDone {
		statusText = "done"
	}

	fmt.Printf("Task chosen: %v\n", todo.Td)
	fmt.Printf("Status: %v\n", statusText)

	// prompt for change
	if promptChangeStatus() {
		err := ChangeStatus(col, todosUsr, index, !todo.IsDone)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
