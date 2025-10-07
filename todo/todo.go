package todo

import (
	"bufio"
	//"context"
	"fmt"
	"os"
	//"time"

	// "slices"
	// "strconv"

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

func DisplayTodos(col *mongo.Collection, userToLogin user.User, allTodos bool, index bool) []bson.ObjectID {

	hasTodos, err := CheckHasTodos(col, userToLogin, allTodos)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	if !hasTodos {
		if allTodos {
			fmt.Printf("User %v does not have any todos\n", userToLogin.Name)
			return []bson.ObjectID{}
		}
		fmt.Printf("User %v does not have any pending todos\n", userToLogin.Name)
		return []bson.ObjectID{}
	}

	//print and return ids as a slice
	ids, err := PrintTodos(col, userToLogin, allTodos, index)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return ids
}

func CreateTodo(usrLogin user.User, col *mongo.Collection) {
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter new todo:")

	//use db to save
	if scn.Scan() {
		err := AddTodoToDB(col, Todo{Td: scn.Text(), User: usrLogin.ID})
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		fmt.Println("There was an error with todo creation. Leaving...")
		os.Exit(1)
	}
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
		if index <= 0 || index > lenTodo {
			fmt.Println("Index is out of range. Please try again.")
			continue
		}
		index--
		return index
	}
}

func removeTodoAtIndex(usrLogin user.User, lTodos *[]Todo, todosUsr []Todo, index int) {
	for pos := range *lTodos {
		if (*lTodos)[pos].User != usrLogin.Name {
			continue
		}
		if (*lTodos)[pos].Name == todosUsr[index].Name {
			*lTodos = slices.Delete(*lTodos, pos, pos+1)
			break
		}
	}
}

func DeleteTodo(usrLogin user.User, lTodos *[]Todo) {
	DisplayTodos(usrLogin, lTodos, false, true)
	todosUsr := userTodos(lTodos, usrLogin, false)
	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "delete")
	removeTodoAtIndex(usrLogin, lTodos, todosUsr, index)
}

func changeTodoAtIndex(usrLogin user.User, lTodos *[]Todo, todosUsr []Todo, index int) {
	scn := bufio.NewScanner(os.Stdin)
	for pos := range *lTodos {
		if (*lTodos)[pos].User != usrLogin.Name {
			continue
		}
		if (*lTodos)[pos].Name == todosUsr[index].Name {
			fmt.Println("Enter new todo:")
			if scn.Scan() {
				(*lTodos)[pos].Name = scn.Text()
				break
			}
			fmt.Println("There was an error with todo creation. Leaving...")
			os.Exit(1)
		}
	}
}

func ChangeTodo(usrLogin user.User, lTodos *[]Todo) {
	DisplayTodos(usrLogin, lTodos, false, true)
	todosUsr := userTodos(lTodos, usrLogin, false)
	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "change")
	changeTodoAtIndex(usrLogin, lTodos, todosUsr, index)
}

func changeStatusAtIndex(usrLogin user.User, lTodos *[]Todo, todosUsr []Todo, index int) {
	scn := bufio.NewScanner(os.Stdin)
	statusStr := "not done"
	for pos := range *lTodos {
		if (*lTodos)[pos].User != usrLogin.Name {
			continue
		}
		if (*lTodos)[pos].Name == todosUsr[index].Name {
			fmt.Print("The current status of the task is: ")
			if (*lTodos)[pos].IsDone {
				statusStr = "done"
			}
			fmt.Printf("%v\n", statusStr)

			for {
				fmt.Println("Would you like to change it (y/n): ")
				if scn.Scan() {
					option := scn.Text()
					if option != "y" && option != "n" {
						fmt.Println("Please choose a proper option")
						continue
					}
					if option == "y" {
						(*lTodos)[pos].IsDone = !(*lTodos)[pos].IsDone
						break
					}
					if option == "n" {
						break
					}
					fmt.Println("There was an error with todo creation. Leaving...")
					os.Exit(1)
				}
			}
		}
	}
}

func ChangeStatusTodo(usrLogin user.User, lTodos *[]Todo) {
	DisplayTodos(usrLogin, lTodos, true, true)
	todosUsr := userTodos(lTodos, usrLogin, true)
	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "done")
	changeStatusAtIndex(usrLogin, lTodos, todosUsr, index)
}
