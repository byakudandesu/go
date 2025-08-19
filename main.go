package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// The struct that all the data will follow
type allUsers struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Org string `json:"org"`
}
// Starting data (Assuming they are data extracted from databases)
var all_users = []allUsers{
	{Name: "Byark", Age: 18, Org: "Jark"},
	{Name: "Mikey", Age: 19, Org: "Jark"},
	{Name: "Dylan", Age: 18, Org: "Jark"},
}

func main() {
	// All the methods possible and the functions linked with those methods
	router := gin.Default()

	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", addUser)
	router.PUT("/users/:id", replaceUser)
	router.PATCH("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

// Returns everything stored 
func getUsers(request *gin.Context) {
	request.JSON(200, all_users)
}

// Returns single user taking in id as parameter which is just the order of the items in the all users slice
func getUser(request *gin.Context) {
	// Check if id can be taken in as a parameter, if cannot convert to int, error. 
	strid := request.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	//If id is out of range, error
	if id < 0 || id >= len(all_users) {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	request.JSON(200, all_users[id])
}

// adding user to databases
func addUser(request *gin.Context) {
	// need to take given json from user and format so the code can read
	var newUser allUsers

	// if the inputted data can be formatted to the struct we made
	if err := request.ShouldBindJSON(&newUser); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// add new user to list
	all_users = append(all_users, newUser)
	request.JSON(201, newUser)
}

func replaceUser(request *gin.Context) {
	var userToUpdate allUsers

	strid := request.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	if id < 0 || id >= len(all_users) {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if err := request.ShouldBindJSON(&userToUpdate); err != nil {
			request.JSON(400, gin.H{"error": err.Error()})
			return
	}

	// completely replacing, so assigning a new slice into an existing slice
	all_users[id] = userToUpdate
	request.JSON(200, all_users[id])
}

func updateUser(request *gin.Context) {
	// interface instead of the struct, since we may have values that are missing
	var updates map[string]interface{}

	strid := request.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	if id < 0 || id >= len(all_users) {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := request.ShouldBindJSON(&updates); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check interface if each item exists. if ok is false, it doesn't exist.
	// we need to tell the code that it is of some type, since it comes from an interface, we don't know the type
	if name, ok := updates["name"]; ok {
		all_users[id].Name = name.(string)
	}
	if age, ok := updates["age"]; ok {
		all_users[id].Age = int(age.(float64))
	}
	if org, ok := updates["org"]; ok {
		all_users[id].Org = org.(string)
	}
	request.JSON(200, all_users[id])
	log.Printf("%v" ,all_users[id])
}

func deleteUser(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("id"))
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	if id < 0 || id >= len(all_users) {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// interesting logic. Appending everything but the user to delete.
	all_users = append(all_users[:id] , all_users[id+1:]...)
	request.JSON(204, nil)
}