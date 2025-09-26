package user

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
	Pass string        `bson:"pass"`
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
		if usr == user.Name {
			return true, user
		}
	}
	return false, User{}
}

func checkPass(pass string, usr User) error {
	err := bcrypt.CompareHashAndPassword([]byte(usr.Pass), []byte(pass))
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
	newUser.Name = userName
	// password needs to be encripted
	hashed, err := HashPassword(userPass, 12)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	newUser.Pass = hashed
	return newUser
}

func CreateUser(col *mongo.Collection, users []User) (bool, User) {
	for {
		tempUser := User{}
		tempUser.Name, tempUser.Pass = inputUserPass("Creating a new user")
		confPass := ""
		//check user in the db. if so repeat cicle
		// also, remove else statements

		if FindUserDBCI(col, tempUser) {
			fmt.Printf("User %v already exists. Usernames are case insensitive (user = User).\n", tempUser.Name)
		} else {
			fmt.Println()
			fmt.Print("Confirm password: ")
			confPass = inputHidden()
			fmt.Println()
			if tempUser.Pass == confPass {
				fmt.Printf("Success. Creating user %v\n", tempUser.Name)
				// set user with encryption
				userToCreate := setNewUser(tempUser.Name, tempUser.Pass)
				return true, userToCreate
			} else {
				fmt.Println("Passwords do not match. Try again")
			}
		}
	}
}

func LoginUserFromDB(col *mongo.Collection, userToLogin *User) error {
	// loop till user and pass match one in the database
	if userToLogin.Name == "" {
		for {
			userName, userPass := inputUserPass("Logging in")
			*userToLogin = User{Name: userName, Pass: userPass}
			fmt.Println()
			// check user
			// try to find the user from db
			userToCompare, err := FindUserDB(col, *userToLogin)
			if err != nil {
				fmt.Println("User not found. Try again")
				*userToLogin = userToCompare
				continue
			}
			validPass := checkPass(userToLogin.Pass, userToCompare)
			if validPass != nil {
				fmt.Println("Wrong password. Try again")
				continue
			}
			*userToLogin = userToCompare
			fmt.Printf("Login successful for %v\n", userToLogin.Name)
			return nil
		}
	}
	// or if it is a new user, update the user with the id set by the db
	userToCompare, _ := FindUserDB(col, *userToLogin)
	*userToLogin = userToCompare
	fmt.Printf("Login successful for %v\n", userToLogin.Name)
	return nil
}

func LogoutUser(userToLogin *User) error {
	*userToLogin = User{}
	return nil
}
