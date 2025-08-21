package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The struct that all the data will follow - now with GORM tags
type User struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Org  string `json:"org"`
}

var db *gorm.DB

func main() {
	// Connect to PostgreSQL
	var dsn string
	if os.Getenv("DB_HOST") != "" {
		// Running in Docker
		dsn = "host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
			" password=" + os.Getenv("DB_PASSWORD") +
			" dbname=" + os.Getenv("DB_NAME") +
			" port=" + os.Getenv("DB_PORT") +
			" sslmode=disable"
	} else {
		// Running locally
		dsn = "host=localhost user=postgres password=mysecretpassword dbname=goapi_db port=5432 sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&User{})

	// Add initial data if table is empty
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		initialUsers := []User{
			{Name: "Byark", Age: 18, Org: "Jark"},
			{Name: "Mikey", Age: 19, Org: "Jark"},
			{Name: "Dylan", Age: 18, Org: "Jark"},
		}
		db.Create(&initialUsers)
		log.Println("Added initial users to database")
	}

	// All the methods possible and the functions linked with those methods
	router := gin.Default()

	router.GET("/test", func(request *gin.Context) {
		request.JSON(200, gin.H{"test": "okokokok"})
	})

	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", addUser)
	router.PUT("/users/:id", replaceUser)
	router.PATCH("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

// Returns everything stored
func getUsers(request *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		request.JSON(400, gin.H{"error": "Could not find users"})
	}
	request.JSON(200, users)

}

// Returns single user by ID from database
func getUser(request *gin.Context) {
	var user User
	id := request.Param("id")

	if err := db.First(&user, id).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	request.JSON(200, user)
}

// adding user to databases
func addUser(request *gin.Context) {
	// need to take given json from user and format so the code can read
	var newUser User

	// if the inputted data can be formatted to the struct we made
	if err := request.ShouldBindJSON(&newUser); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// add new user to database
	if err := db.Create(&newUser).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	request.JSON(201, newUser)
}

func replaceUser(request *gin.Context) {
	var userToUpdate User
	id := request.Param("id")

	// Check if user exists
	var existingUser User
	if err := db.First(&existingUser, id).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := request.ShouldBindJSON(&userToUpdate); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Keep the same ID
	userToUpdate.ID = existingUser.ID

	// Save replaces all fields
	if err := db.Save(&userToUpdate).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	request.JSON(200, userToUpdate)
}

func updateUser(request *gin.Context) {
	// interface instead of the struct, since we may have values that are missing
	var updates map[string]interface{}
	id := request.Param("id")

	// Check if user exists
	var user User
	if err := db.First(&user, id).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := request.ShouldBindJSON(&updates); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	// Fetch updated user
	db.First(&user, id)
	request.JSON(200, user)
	log.Printf("Updated user: %v", user)
}

func deleteUser(request *gin.Context) {
	id := request.Param("id")

	// Delete user from database
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		request.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}
	if result.RowsAffected == 0 {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	request.JSON(204, nil)
}
