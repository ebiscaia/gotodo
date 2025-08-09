package main

import "fmt"

type Todo struct {
	name string
	user string
}

type User struct {
	name string
	pass string
}

func inputUserPass(msg string) User {
	user := User{}
	fmt.Printf("%v", msg)
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

func main() {
	//mockup list of todos
	todos := []Todo{}
	todos = append(todos, Todo{name: "take the puppy for a lap", user: "eddie"})
	todos = append(todos, Todo{name: "take the rubbish out", user: "eddie"})

	//mockup list of todos
	users := []User{}
	users = append(users, User{name: "user1", pass: "pass1"}, User{name: "user2", pass: "pass2"})

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
