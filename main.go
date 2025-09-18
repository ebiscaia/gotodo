package main

import (
	"fmt"
	"os"

	"gotodo/todo"
	"gotodo/user"
)

func main() {
	// Some initial variables
	userToLogin := user.User{}
	users := []user.User{}
	todos := []todo.Todo{}

	//main loop
	for {
		//present initial menu
		menuOption := StartMenu(users, userToLogin)

		//go to login, create user, logout or exit depending on chosen option
		err := HandleMainMenu(menuOption, &users, &userToLogin)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		// Do not allow an empty user to deal with todos
		if userToLogin.Name == "" {
			continue
		}

		// present menu with todo options after a user is logged in
		for {
			menuTodoOption := TodoMenu()
			result, err := HandleTodoMenu(userToLogin, menuTodoOption, &todos)
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
