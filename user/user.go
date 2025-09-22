package user

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func CreateUser(users []User) (bool, User) {
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
func LoginUser(userToLogin *User, users []User) error {
	if userToLogin.Name == "" {
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
					fmt.Printf("Login successful for %v\n", userToLogin.Name)
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

func LogoutUser(userToLogin *User) error {
	*userToLogin = User{}
	return nil
}
