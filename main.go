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

func main() {
	//mockup list of todos
	todos := []Todo{}
	todos = append(todos, Todo{name: "take the puppy for a lap", user: "eddie"})
	todos = append(todos, Todo{name: "take the rubbish out", user: "eddie"})

	//mockup list of todos
	users := []User{}
	users = append(users, User{name: "user1", pass: "pass1"}, User{name: "user2", pass: "pass2"})

	//loop through the todo list
	for _, todo := range todos {
		fmt.Printf("Todo: %v\n", todo.name)
	}

	//loop through the user list
	fmt.Println()
	for _, user := range users {
		fmt.Printf("Todo: %v\n", user.name)
	}
	fmt.Println()

	//Simulate a login using inputs
	// user := User{}
	// fmt.Print("Enter username: ")
	// fmt.Scanf("%s", &user.name)
	// fmt.Print("Enter password: ")
	// fmt.Scanf("%s", &user.pass)

	// validUser, validPass := checkUserPass(user, users)
	// if validUser {
	// 	if validPass {
	// 		fmt.Printf("Login successful for %v\n", user.name)
	// 	} else {
	// 		fmt.Printf("Wrong password. User %v needs to try again.\n", user.name)
	// 	}
	// } else {
	// 	fmt.Printf("User %v does not exist. Create one before logging in.\n", user.name)
	// }

	//Simulate creating a user
	userToCreate := User{}
	confPass := ""
	fmt.Print("Enter username: ")
	fmt.Scanf("%s", &userToCreate.name)
	fmt.Print("Enter password: ")
	fmt.Scanf("%s", &userToCreate.pass)
	validUser, _ := checkUserPass(userToCreate, users)
	if validUser {
		fmt.Printf("User %v already exists. Try a different user name.\n", userToCreate.name)
	} else {
		fmt.Print("\nConfirm password: ")
		fmt.Scanf("%s", &confPass)
		if userToCreate.pass == confPass {
			fmt.Printf("No user associated with %v. Creating new user\n", userToCreate.name)
		} else {
			fmt.Println("Passwords do not match. Try again")
		}
	}

}
