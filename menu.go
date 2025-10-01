package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"gotodo/todo"
	"gotodo/user"
)

type Menu struct {
	message     string
	instruction string
	index       int
}

func promptMenu(fullMenu []Menu) {
	fmt.Println("Choose one of the following options: ")
	fmt.Println()
	for pos := range fullMenu {
		fullMenu[pos].index = pos + 1
		if fullMenu[pos].instruction == "exit" {
			fullMenu[pos].index = 0
		}
		fmt.Printf("%v - %v\n", fullMenu[pos].index, fullMenu[pos].message)
	}
}

func validateMenu(err error, menuChosen Menu, menu []Menu) bool {
	if err != nil {
		fmt.Println("Error while inputing data.")
		panic(1)
	}

	if menuChosen.index < len(menu) {
		return true
	} else {
		fmt.Println("Wrong option chosen. Try again")

	}
	return false
}

func inputMenu(menuItems []Menu) Menu {
	menuChosen := Menu{}
	for {
		promptMenu(menuItems)
		//Enter chosen menu option
		fmt.Print("Enter option: ")
		_, err := fmt.Scan(&menuChosen.index)
		if validateMenu(err, menuChosen, menuItems) {
			break
		}
	}

	//Get the right option based on the index
	for _, menuItem := range menuItems {
		if menuItem.index == menuChosen.index {
			menuChosen = menuItem
		}
	}
	return menuChosen
}

func StartMenu(curUser user.User, col *mongo.Collection) Menu {
	menuStart := []Menu{}
	menuStart = append(menuStart, Menu{message: "Create user", instruction: "create"})
	menuStart = append(menuStart, Menu{message: "Login", instruction: "login"})
	menuStart = append(menuStart, Menu{message: "Logout", instruction: "logout"})
	menuStart = append(menuStart, Menu{message: "Return to todo menu", instruction: "todo"})
	menuStart = append(menuStart, Menu{message: "Exit", instruction: "exit"})

	// replace len(users) by the number of documents in the users collection
	toLogin, err := user.CheckUsers(col)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if !toLogin {
		menuStart = slices.Delete(menuStart, 1, 4)
	} else {
		if curUser.Name == "" {
			menuStart = slices.Delete(menuStart, 2, 4)
		} else {
			menuStart = slices.Delete(menuStart, 0, 2)
		}
	}

	menuChosen := inputMenu(menuStart)
	return menuChosen
}

func TodoMenu() Menu {
	//create a menu with todo operations
	menuTodo := []Menu{}
	menuTodo = append(menuTodo, Menu{instruction: "create", message: "Create todo"})
	menuTodo = append(menuTodo, Menu{instruction: "delete", message: "Delete todo"})
	menuTodo = append(menuTodo, Menu{instruction: "change", message: "Change todo"})
	menuTodo = append(menuTodo, Menu{instruction: "done", message: "Mark as done/undone"})
	menuTodo = append(menuTodo, Menu{instruction: "list", message: "List pending todos"})
	menuTodo = append(menuTodo, Menu{instruction: "listAll", message: "List all todos"})
	menuTodo = append(menuTodo, Menu{instruction: "previous", message: "Previous menu"})
	menuTodo = append(menuTodo, Menu{instruction: "exit", message: "Exit program"})
	menuChosen := inputMenu(menuTodo)
	return menuChosen
}

func HandleMainMenu(menuOption Menu, userToLogin *user.User, col *mongo.Collection) error {
	switch menuOption.instruction {
	case "create":
		successCreate, userToCreate := user.CreateUser(col)
		if successCreate {
			// add here to the db
			err := user.AddUserToDB(col, userToCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			*userToLogin = userToCreate
		}

		err := user.LoginUserFromDB(col, userToLogin)
		if err == nil {
			fmt.Printf("User %v is logged in\n", userToLogin.Name)
		}
		return err

	case "login":
		err := user.LoginUserFromDB(col, userToLogin)
		if err == nil {
			fmt.Printf("User %v is logged in\n", userToLogin.Name)
		}
		return err

	case "logout":
		err := user.LogoutUser(userToLogin)
		return err

	case "todo":
		return nil

	case "exit":
		fmt.Println("Exiting... ")
		os.Exit(0)

	default:
		return errors.New("there is an issue with the application")
	}
	return errors.New("there is an issue with the application")
}

func HandleTodoMenu(userToLogin user.User, menuOption Menu, listTodos *[]todo.Todo, col *mongo.Collection) (string, error) {
	switch menuOption.instruction {
	case "create":
	case "delete":
		todo.DeleteTodo(userToLogin, listTodos)
		return "continue", nil
	case "change":
		todo.ChangeTodo(userToLogin, listTodos)
		return "continue", nil
	case "done":
		todo.ChangeStatusTodo(userToLogin, listTodos)
		return "continue", nil
	case "list":
		todo.DisplayTodos(userToLogin, listTodos, false, false)
		return "continue", nil
	case "listAll":
		todo.DisplayTodos(userToLogin, listTodos, true, false)
		todo.CreateTodo(userToLogin, col)
		return "continue", nil

	case "previous":
		return "previous", nil

	case "exit":
		fmt.Println("Exiting... ")
		os.Exit(0)

	default:
		fmt.Println("There is an issue with the application. Leaving...")
		os.Exit(1)
	}
	return "", errors.New("Menu Error")
}
