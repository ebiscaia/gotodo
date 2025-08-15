package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
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

func loginUser(usr User, users []User) (bool, User) {
	userToLogin := usr

	if userToLogin.name == "" {
		for {
			userToLogin = inputUserPass("Logging in")
			validUser, validPass := checkUserPass(userToLogin, users)
			if validUser {
				if validPass {
					fmt.Printf("Login successful for %v\n", userToLogin.name)
					return true, userToLogin
				} else {
					fmt.Printf("Wrong password. User %v needs to try again.\n", userToLogin.name)
				}
			} else {
				fmt.Printf("User %v does not exist. Create one before logging in.\n", userToLogin.name)
			}
		}
	} else {
		return true, userToLogin
	}
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

func handleMainMenu(menuOption Menu, users []User, userToLogin User) (User, error) {
	switch menuOption.instruction {
	case "create":
		successCreate, userToCreate := createUser(users)
		if successCreate {
			users = append(users, userToCreate)
			userToLogin = userToCreate
		}

		successLogin, userToLogin := loginUser(userToLogin, users)
		if successLogin {
			fmt.Printf("User %v is logged in\n", userToLogin.name)
		}
		return userToLogin, nil

	case "login":
		successLogin, userToLogin := loginUser(userToLogin, users)
		if successLogin {
			fmt.Printf("User %v is logged in\n", userToLogin.name)
		}
		return userToLogin, nil

	case "exit":
		fmt.Println("Exiting... ")
		os.Exit(0)

	default:
		fmt.Println("There is an issue with the application. Leaving...")
		os.Exit(1)
	}
	return User{}, errors.New("Menu Error")
}

func displayTodos(userToLogin User, listTodos *[]Todo, allTodos bool) {
	printItems := 0
	for _, todo := range *listTodos {
		if todo.user == userToLogin.name {
			if allTodos {
				fmt.Printf("%v\n", todo.name)
				printItems++
			} else {
				if !todo.isDone {
					fmt.Printf("%v\n", todo.name)
					printItems++
				}
			}
		}
	}
	if printItems == 0 {
		if allTodos {
			fmt.Printf("User %v does not have any todos\n", userToLogin.name)
		} else {
			fmt.Printf("User %v does not have any pending todos\n", userToLogin.name)
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

func handleTodoMenu(userToLogin User, menuOption Menu, listTodos *[]Todo) (string, error) {
	switch menuOption.instruction {
	case "create":
		createTodo(userToLogin, listTodos)
		return "continue", nil
	case "delete":
		fmt.Println("Delete todo logic")
		return "continue", nil
	case "change":
		fmt.Println("Change todo logic")
		return "continue", nil
	case "done":
		fmt.Println("Done todo logic")
		return "continue", nil
	case "list":
		displayTodos(userToLogin, listTodos, false)
		return "continue", nil
	case "listAll":
		displayTodos(userToLogin, listTodos, true)
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
	// Empty userToLogin
	userToLogin := User{}

	//mockup list of todos
	todos := []Todo{}
	todos = append(todos, Todo{name: "take the puppy for a lap", user: "eddie"})
	todos = append(todos, Todo{name: "take the rubbish out", user: "eddie", isDone: true})
	todos = append(todos, Todo{name: "do the tax return", user: "user2", isDone: true})
	todos = append(todos, Todo{name: "do the shopping", user: "user1"})

	//mockup list of users
	users := []User{}
	users = append(users, User{name: "user1", pass: "pass1"}, User{name: "user2", pass: "pass2"})

	//menu items
	menuOption := startMenu(users)

	//go to login, create user or exit depending on chosen option
	userToLogin, err := handleMainMenu(menuOption, users, userToLogin)
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

	//loop through the todo list
	for _, todo := range todos {
		fmt.Printf("Todo: %v\n", todo.name)
	}
}
