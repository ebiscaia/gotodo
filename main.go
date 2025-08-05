package main

import "fmt"

type Todo struct {
	name string
	user string
}

func main() {
	todo := Todo{name: "take the puppy for a lap", user: "eddie"}
	fmt.Printf("Todo: %v\n", todo.name)
	//mockup list of todos
	todos := []Todo{}
}
