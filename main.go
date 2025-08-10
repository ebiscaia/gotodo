package main

import (
	"fmt"
	"slices"
)

type Todo struct {
	name string
	user string
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
			fmt.Print("\nConfirm password: ")
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

func startMenu(users []User) Menu {
	menu := []Menu{}
	menu = append(menu, Menu{message: "Create user", instruction: "create"})
	menu = append(menu, Menu{message: "Login", instruction: "login"})
	menu = append(menu, Menu{message: "Exit", instruction: "exit"})
	menuChosen := Menu{}

	if len(users) == 0 {
		menu = slices.Delete(menu, 1, 2)
	}

	for {
		promptMenu(menu)
		//Enter chosen menu option
		fmt.Print("Enter option: ")
		_, err := fmt.Scan(&menuChosen.index)
		if validateMenu(err, menuChosen, menu) {
			break
		}
	}

	//Get the right option based on the index
	for _, menuItem := range menu {
		if menuItem.index == menuChosen.index {
			menuChosen = menuItem
		}
	}
	return menuChosen
}

func main() {
	//mockup list of todos
	todos := []Todo{}
	todos = append(todos, Todo{name: "take the puppy for a lap", user: "eddie"})
	todos = append(todos, Todo{name: "take the rubbish out", user: "eddie"})

	//mockup list of users
	users := []User{}
	users = append(users, User{name: "user1", pass: "pass1"}, User{name: "user2", pass: "pass2"})

	//menu items
	menuOption := startMenu(users)

	//loop through menu but avoid login if there are no users

	//For now just print the result
	fmt.Printf("\nChosen option: %v", menuOption.instruction)

	//loop through the user list
	for _, user := range users {
		fmt.Printf("User: %v\n", user.name)
	}
	fmt.Println()

	// Empty userToLogin
	userToLogin := User{}

	//Simulate creating a user
	successCreate, userToCreate := createUser(users)
	if successCreate {
		users = append(users, userToCreate)
		userToLogin = userToCreate
	}

	fmt.Printf("%v\n", users)

	//Simulate a login using inputs
	successLogin, userToLogin := loginUser(userToLogin, users)
	if successLogin {
		fmt.Printf("User %v is logged in", userToLogin.name)
	}

	//loop through the todo list
	for _, todo := range todos {
		fmt.Printf("Todo: %v\n", todo.name)
	}
}
