package main

import "fmt"
import "github.com/nudelfabrik/GoRadius/database"

func main() {
	db := database.NewDatabase()

	fmt.Println(db.GetUser("Test"))
}
