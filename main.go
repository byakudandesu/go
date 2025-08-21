package main

import (
	"goapi/handlers"
	"goapi/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The struct that all the data will follow - now with GORM tags

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
	db.AutoMigrate(&models.Team{}, &models.User{})

	h := &handlers.Handler{DB: db}

	// All the methods possible and the functions linked with those methods
	router := gin.Default()

	router.GET("/test", func(request *gin.Context) {
		request.JSON(200, gin.H{"test": "okokokok"})
	})

	router.GET("/users", h.GetUsers)
	router.GET("/users/:id", h.GetUser)
	router.POST("/users", h.AddUser)
	router.PUT("/users/:id", h.ReplaceUser)
	router.PATCH("/users/:id", h.UpdateUser)
	router.DELETE("/users/:id", h.DeleteUser)

	router.POST("/teams", handlers.AdminOnly(), h.CreateTeam)
	router.GET("/teams", h.GetTeams)
	router.DELETE("/teams/:team_id", handlers.AdminOnly(), h.DeleteTeam)

	router.POST("/teams/:team_id/:user_id", h.AddUserToTeam)
	router.DELETE("/teams/:team_id/:user_id", h.RemoveUserFromTeam)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
