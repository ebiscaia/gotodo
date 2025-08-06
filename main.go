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

func checkUser(user string, userSlice []User) bool {
	for _, username := range userSlice {
		if user == username.name {
			return true
		}
	}
	return false
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

	for _, userToCheck := range []string{"user2", "user20"} {
		if checkUser(userToCheck, users) {
			fmt.Printf("User %v exists\n", userToCheck)
		} else {
			fmt.Printf("User %v does not exist\n", userToCheck)
		}
	}

}
