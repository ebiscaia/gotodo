package main

import "fmt"

type Todo struct {
	name string
	user string
}

func main() {
	fmt.Println("program started")
	//mockup user and password
	todo := Todo{name: "take the puppy for a lap", user: "eddie"}
}
