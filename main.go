package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Todo struct {
	name   string
	user   string
	isDone bool
}

type User struct {
	name string
	pass string
}

type Menu struct {
	message     string
	instruction string
	index       int
}

func inputUserPass(msg string) User {
	user := User{}
	fmt.Printf("\n%v\n", msg)
	fmt.Print("Enter username: ")
	fmt.Scanf("%s", &user.name)
	fmt.Print("Enter password: ")
	fmt.Scanf("%s", &user.pass)
	return user
}

func checkUserPass(usr User, userSlice []User) (bool, bool) {
	for _, user := range userSlice {
		if usr.name == user.name {
			if usr.pass == user.pass {
				return true, true
			} else {
				return true, false
			}
		}
	}
	return false, false
}

func createUser(users []User) (bool, User) {
	for {
		userToCreate := inputUserPass("Creating a new user")
		confPass := ""
		validUser, _ := checkUserPass(userToCreate, users)
		if validUser {
			fmt.Printf("User %v already exists. Try a different user name.\n", userToCreate.name)
		} else {
			fmt.Print("Confirm password: ")
			fmt.Scanf("%s", &confPass)
			if userToCreate.pass == confPass {
				fmt.Printf("Success. Creating user %v\n", userToCreate.name)
				return true, userToCreate
			} else {
				fmt.Println("Passwords do not match. Try again")
			}
		}
	}
}

func loginUser(userToLogin *User, users []User) error {

	if userToLogin.name == "" {
		for {
			*userToLogin = inputUserPass("Logging in")
			validUser, validPass := checkUserPass(*userToLogin, users)
			if validUser {
				if validPass {
					fmt.Printf("Login successful for %v\n", userToLogin.name)
					return nil
				}
				return errors.New("Wrong password")
			}
			return errors.New("User does not exist")
		}
	}
	return nil
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

func startMenu(users []User) Menu {
	menuStart := []Menu{}
	menuStart = append(menuStart, Menu{message: "Create user", instruction: "create"})
	menuStart = append(menuStart, Menu{message: "Login", instruction: "login"})
	menuStart = append(menuStart, Menu{message: "Exit", instruction: "exit"})

	if len(users) == 0 {
		menuStart = slices.Delete(menuStart, 1, 2)
	}

	menuChosen := inputMenu(menuStart)
	return menuChosen
}

func todoMenu() Menu {
	//create a menu with todo operations
	menuTodo := []Menu{}
	menuTodo = append(menuTodo, Menu{instruction: "create", message: "Create todo"})
	menuTodo = append(menuTodo, Menu{instruction: "delete", message: "Delete todo"})
	menuTodo = append(menuTodo, Menu{instruction: "change", message: "Change todo"})
	menuTodo = append(menuTodo, Menu{instruction: "done", message: "Mark as done"})
	menuTodo = append(menuTodo, Menu{instruction: "list", message: "List pending todos"})
	menuTodo = append(menuTodo, Menu{instruction: "listAll", message: "List all todos"})
	menuTodo = append(menuTodo, Menu{instruction: "previous", message: "Previous menu"})
	menuTodo = append(menuTodo, Menu{instruction: "exit", message: "Exit program"})
	menuChosen := inputMenu(menuTodo)
	return menuChosen
}

func handleMainMenu(menuOption Menu, users []User, userToLogin *User) error {
	switch menuOption.instruction {
	case "create":
		successCreate, userToCreate := createUser(users)
		if successCreate {
			users = append(users, userToCreate)
			*userToLogin = userToCreate
		}

		err := loginUser(userToLogin, users)
		if err == nil {
			fmt.Printf("User %v is logged in\n", userToLogin.name)
		}
		return err

	case "login":
		err := loginUser(userToLogin, users)
		if err == nil {
			fmt.Printf("User %v is logged in\n", userToLogin.name)
		}
		return err

	case "exit":
		fmt.Println("Exiting... ")
		os.Exit(0)

	default:
		return errors.New("there is an issue with the application")
	}
	return errors.New("there is an issue with the application")
}

func userTodos(listTodos *[]Todo, userToLogin User, allTodos bool) []Todo {
	todosUsr := []Todo{}
	for _, tdo := range *listTodos {
		if tdo.user == userToLogin.name {
			if allTodos {
				todosUsr = append(todosUsr, tdo)
			} else {
				if !tdo.isDone {
					todosUsr = append(todosUsr, tdo)
				}
			}
		}
	}
	return todosUsr
}

func displayTodos(userToLogin User, listTodos *[]Todo, allTodos bool, index bool) {
	todosUsr := userTodos(listTodos, userToLogin, allTodos)

	if len(todosUsr) == 0 {
		if allTodos {
			fmt.Printf("User %v does not have any todos\n", userToLogin.name)
			return
		} else {
			fmt.Printf("User %v does not have any pending todos\n", userToLogin.name)
			return
		}
	}

	if index {
		if allTodos {
			for ind, todo := range todosUsr {
				fmt.Printf("%v - %v %v\n", ind+1, todo.name, todo.isDone)
			}
		} else {
			for ind, todo := range todosUsr {
				fmt.Printf("%v - %v\n", ind+1, todo.name)
			}
		}
	} else {
		if allTodos {
			for _, todo := range todosUsr {
				fmt.Printf("%v %v\n", todo.name, todo.isDone)
			}
		} else {
			for _, todo := range todosUsr {
				fmt.Printf("%v %v\n", todo.name, todo.isDone)
			}
		}
	}
}

func createTodo(usrLogin User, lTodos *[]Todo) {
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter new todo:")
	if scn.Scan() {
		*lTodos = append(*lTodos, Todo{name: scn.Text(), user: usrLogin.name})

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

func removeTodoAtIndex(usrLogin User, lTodos *[]Todo, todosUsr []Todo, index int) {
	for pos := range *lTodos {
		if (*lTodos)[pos].user != usrLogin.name {
			continue
		}
		if (*lTodos)[pos].name == todosUsr[index].name {
			*lTodos = slices.Delete(*lTodos, pos, pos+1)
			break
		}
	}
}

func deleteTodo(usrLogin User, lTodos *[]Todo) {
	displayTodos(usrLogin, lTodos, false, true)
	todosUsr := userTodos(lTodos, usrLogin, false)
	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "delete")
	removeTodoAtIndex(usrLogin, lTodos, todosUsr, index)
}

func changeTodoAtIndex(usrLogin User, lTodos *[]Todo, todosUsr []Todo, index int) {
	scn := bufio.NewScanner(os.Stdin)
	for pos := range *lTodos {
		if (*lTodos)[pos].user != usrLogin.name {
			continue
		}
		if (*lTodos)[pos].name == todosUsr[index].name {
			fmt.Println("Enter new todo:")
			if scn.Scan() {
				(*lTodos)[pos].name = scn.Text()
				break
			}
			fmt.Println("There was an error with todo creation. Leaving...")
			os.Exit(1)
		}
	}
}

func changeTodo(usrLogin User, lTodos *[]Todo) {
	displayTodos(usrLogin, lTodos, false, true)
	todosUsr := userTodos(lTodos, usrLogin, false)
	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "change")
	changeTodoAtIndex(usrLogin, lTodos, todosUsr, index)
}

func changeStatusAtIndex(usrLogin User, lTodos *[]Todo, todosUsr []Todo, index int) {
	scn := bufio.NewScanner(os.Stdin)
	statusStr := "not done"
	for pos := range *lTodos {
		if (*lTodos)[pos].user != usrLogin.name {
			continue
		}
		if (*lTodos)[pos].name == todosUsr[index].name {
			fmt.Print("The current status of the task is: ")
			if (*lTodos)[pos].isDone {
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
						(*lTodos)[pos].isDone = !(*lTodos)[pos].isDone
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

func changeStatusTodo(usrLogin User, lTodos *[]Todo) {
	displayTodos(usrLogin, lTodos, true, true)
	todosUsr := userTodos(lTodos, usrLogin, true)
	if len(todosUsr) == 0 {
		return
	}
	index := inputIndex(len(todosUsr), "done")
	changeStatusAtIndex(usrLogin, lTodos, todosUsr, index)
}

func handleTodoMenu(userToLogin User, menuOption Menu, listTodos *[]Todo) (string, error) {
	switch menuOption.instruction {
	case "create":
		createTodo(userToLogin, listTodos)
		return "continue", nil
	case "delete":
		deleteTodo(userToLogin, listTodos)
		return "continue", nil
	case "change":
		changeTodo(userToLogin, listTodos)
		return "continue", nil
	case "done":
		changeStatusTodo(userToLogin, listTodos)
		return "continue", nil
	case "list":
		displayTodos(userToLogin, listTodos, false, false)
		return "continue", nil
	case "listAll":
		displayTodos(userToLogin, listTodos, true, false)
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

func main() {
	// Some initial variables
	userToLogin := User{}
	users := []User{}
	todos := []Todo{}

	//main loop
	for {
		menuOption := startMenu(users)

		//go to login, create user or exit depending on chosen option
		err := handleMainMenu(menuOption, users, &userToLogin)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		// present menu with todo options
		for {
			menuTodoOption := todoMenu()
			result, err := handleTodoMenu(userToLogin, menuTodoOption, &todos)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			if result == "previous" {
				break
			}
		}
	}
}
