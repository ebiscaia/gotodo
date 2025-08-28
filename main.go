package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
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

func inputHidden() string {
	input, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return string(input)
}

func inputUserPass(msg string) (string, string) {
	userName := ""
	userPass := ""
	fmt.Printf("\n%v\n", msg)
	fmt.Print("Enter username: ")
	fmt.Scanf("%s", &userName)
	fmt.Print("Enter password: ")
	userPass = inputHidden()
	return userName, userPass
}

func checkUser(usr string, userSlice []User) (bool, User) {
	for _, user := range userSlice {
		if usr == user.name {
			return true, user
		}
	}
	return false, User{}
}

func checkPass(pass string, usr User) error {
	err := bcrypt.CompareHashAndPassword([]byte(usr.pass), []byte(pass))
	return err
}

func HashPassword(pass string, cost int) (string, error) {
	//default value for cost
	if cost == 0 {
		cost = 12
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func setNewUser(userName string, userPass string) User {
	newUser := User{}
	newUser.name = userName
	// password needs to be encripted
	hashed, err := HashPassword(userPass, 12)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	newUser.pass = hashed
	return newUser
}

func createUser(users []User) (bool, User) {
	for {
		tempUserName, tempUserPass := inputUserPass("Creating a new user")
		confPass := ""
		validUser, _ := checkUser(tempUserName, users)
		if validUser {
			fmt.Printf("User %v already exists. Try a different user name.\n", tempUserName)
		} else {
			fmt.Println()
			fmt.Print("Confirm password: ")
			confPass = inputHidden()
			fmt.Println()
			if tempUserPass == confPass {
				fmt.Printf("Success. Creating user %v\n", tempUserName)
				// set user with encryption
				userToCreate := setNewUser(tempUserName, tempUserPass)
				return true, userToCreate
			} else {
				fmt.Println("Passwords do not match. Try again")
			}
		}
	}
}

// needs some reworking as function has been modified
func loginUser(userToLogin *User, users []User) error {
	if userToLogin.name == "" {
		for {
			//input user name and pass
			userName, userPass := inputUserPass("Logging in")
			//check user
			validUser, userToCheck := checkUser(userName, users)
			fmt.Println()
			if validUser {
				validPass := checkPass(userPass, userToCheck)
				if validPass == nil {
					*userToLogin = userToCheck
					fmt.Printf("Login successful for %v\n", userToLogin.name)
					return nil
				}
				fmt.Println("Wrong password. Try again")
				continue
			}
			fmt.Println("Wrong user. Try again")
		}
	}
	return nil
}

func logoutUser(userToLogin *User) error {
	*userToLogin = User{}
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

func startMenu(users []User, curUser User) Menu {
	menuStart := []Menu{}
	menuStart = append(menuStart, Menu{message: "Create user", instruction: "create"})
	menuStart = append(menuStart, Menu{message: "Login", instruction: "login"})
	menuStart = append(menuStart, Menu{message: "Logout", instruction: "logout"})
	menuStart = append(menuStart, Menu{message: "Exit", instruction: "exit"})

	if len(users) == 0 {
		menuStart = slices.Delete(menuStart, 1, 3)
	} else {
		if curUser.name == "" {
			menuStart = slices.Delete(menuStart, 2, 3)
		} else {
			menuStart = slices.Delete(menuStart, 0, 2)
		}
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

func handleMainMenu(menuOption Menu, users *[]User, userToLogin *User) error {
	switch menuOption.instruction {
	case "create":
		successCreate, userToCreate := createUser(*users)
		if successCreate {
			*users = append(*users, userToCreate)
			*userToLogin = userToCreate
		}

		err := loginUser(userToLogin, *users)
		if err == nil {
			fmt.Printf("User %v is logged in\n", userToLogin.name)
		}
		return err

	case "login":
		err := loginUser(userToLogin, *users)
		if err == nil {
			fmt.Printf("User %v is logged in\n", userToLogin.name)
		}
		return err

	case "logout":
		err := logoutUser(userToLogin)
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
		//present initial menu
		menuOption := startMenu(users, userToLogin)

		//go to login, create user, logout or exit depending on chosen option
		err := handleMainMenu(menuOption, &users, &userToLogin)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		// Do not allow an empty user to deal with todos
		if userToLogin.name == "" {
			continue
		}

		// present menu with todo options after a user is logged in
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
