package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"gotodo/todo"
	"gotodo/user"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type credential struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func (c credential) MongoURI() string {
	escapedUser := url.QueryEscape(c.User)
	escapedPass := url.QueryEscape(c.Pass)

	return fmt.Sprintf("mongodb://%s:%s@%s:%d",
		escapedUser, escapedPass,
		c.Ip, c.Port)
}

func main() {
	//import json file and read it
	fileName := "credentials.json"
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	fmt.Printf("File \"%v\" opened successfully\n", fileName)

	//pass the json data into the struct
	byteValue, _ := io.ReadAll(jsonFile)
	mongoCredentials := credential{}
	json.Unmarshal(byteValue, &mongoCredentials)

	//create the context and connect to the server using the credentials
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(mongoCredentials.MongoURI()))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	// create database and collections
	db := client.Database("todoapp")
	usersCollection := db.Collection("users")

	// Some initial variables
	userToLogin := user.User{}
	users := []user.User{}
	todos := []todo.Todo{}

	//main loop
	for {
		//present initial menu
		//modify this to check the number of documents in the database
		menuOption := StartMenu(users, userToLogin, ctx, usersCollection)

		//go to login, create user, logout or exit depending on chosen option
		err := HandleMainMenu(menuOption, &users, &userToLogin, ctx, usersCollection)
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
