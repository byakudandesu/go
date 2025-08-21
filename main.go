package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The struct that all the data will follow - now with GORM tags
type User struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Org    string `json:"org"`
	TeamID *uint  `json:"team_id"`
	Salary int    `json:"salary"`
}

type Team struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Budget     int    `json:"budget"`
	UsedBudget int    `json:"used_budget"`
	Users      []User `json:"users,omitempty" gorm:"foreignKey:TeamID"`
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
	db.AutoMigrate(&Team{}, &User{})

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

	router.POST("/teams", adminOnly(), createTeam)
	router.GET("/teams", getTeams)
	router.DELETE("/teams/:team_id", adminOnly(), deleteTeam)

	router.POST("/teams/:team_id/:user_id", addUserToTeam)
	router.DELETE("/teams/:team_id/:user_id", removeUserFromTeam)

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
	// this would ignore the salary calculation, so the user must be added through the addUserToTeam function
	newUser.TeamID = nil
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

// team stuff
func createTeam(request *gin.Context) {
	var team Team

	if err := request.ShouldBindJSON(&team); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}
	team.UsedBudget = 0

	if err := db.Create(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to create team"})
		return
	}

	request.JSON(201, team)
}

func getTeams(request *gin.Context) {
	var teams []Team
	if err := db.Preload("Users").Find(&teams).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to fetch teams"})
		return
	}
	request.JSON(200, teams)
}

func addUserToTeam(request *gin.Context) {
	// which team? which user?
	teamID := request.Param("team_id")
	userID := request.Param("user_id")

	// finding requested teams and users into go varibles, if they don't exist, error
	var team Team
	if err := db.First(&team, teamID).Error; err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// check if user is not in a team yet
	if user.TeamID != nil {
		request.JSON(400, gin.H{"error": "User already belongs in team"})
		return
	}

	// need to check if budget of team can affor user
	newBudgetUsed := team.UsedBudget + user.Salary
	if newBudgetUsed > team.Budget {
		request.JSON(400, gin.H{
			"error":        "Team cannot afford user",
			"budget":       team.Budget,
			"current_used": team.UsedBudget,
			"user_salary":  user.Salary,
			"would_need":   newBudgetUsed,
		})
		return
	}

	// convert params string into uint and assign it to user's teamID
	teamIDUint, _ := strconv.Atoi(teamID)
	teamIDValue := uint(teamIDUint)
	user.TeamID = &teamIDValue

	team.UsedBudget = newBudgetUsed

	if err := db.Save(&user).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}
	if err := db.Save(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update team budget"})
		return
	}

	request.JSON(200, gin.H{
		"message":               "User added to team successfully",
		"user":                  user.Name,
		"team":                  team.Name,
		"team_budget_remaining": team.Budget - team.UsedBudget,
	})
}
func removeUserFromTeam(request *gin.Context) {
	teamID := request.Param("team_id")
	userID := request.Param("user_id")

	var team Team
	if err := db.First(&team, teamID).Error; err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// changing the teamid for said user to nil so there is no relationship
	user.TeamID = nil
	team.UsedBudget -= user.Salary

	if err := db.Save(&user).Error; err != nil {
		request.JSON(500, gin.H{"error": "Error saving to database"})
		return
	}
	if err := db.Save(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Error saving to database"})
		return
	}

	request.JSON(200, gin.H{"message": "User removed from team"})
}

func deleteTeam(request *gin.Context) {
	teamID := request.Param("team_id")

	var team Team

	teamIDInt, _ := strconv.Atoi(teamID)
	theteamID := uint(teamIDInt)

	if err := db.First(&team, teamID).Error; err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}

	// find all the users where the team id is the teamid in params and set them to null before team deletion
	if err := db.Model(&User{}).Where("team_id = ?", theteamID).Update("team_id",
		nil).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to remove users from team"})
		return
	}

	if err := db.Delete(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to delete team"})
		return
	}

	request.JSON(200, gin.H{
		"message": "Team deleted and users freed",
		"team":    team.Name,
	})

}

func adminOnly() gin.HandlerFunc {
	return func(request *gin.Context) {
		adminKey := request.GetHeader("admin-key")

		if adminKey != "byakubyaku" {
			request.JSON(403, gin.H{"error": "Admin access required"})
			request.Abort()
			return
		}
		request.Next()
	}
}
