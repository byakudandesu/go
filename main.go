package main

import (
	"goapi/handlers"
	"goapi/models"
	"goapi/repositories"
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

	// Create repositories
	userRepo := repositories.NewUserRepository(db)
	teamRepo := repositories.NewTeamRepository(db)

	// Create handler with repositories
	h := &handlers.Handler{
		UserRepo: userRepo,
		TeamRepo: teamRepo,
	}

	// All the methods possible and the functions linked with those methods
	router := gin.Default()

	router.GET("/test", func(request *gin.Context) {
		request.JSON(200, gin.H{"test": "okokokok"})
	})

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", h.GetUsers)
			users.GET("/:id", h.GetUser)
			users.POST("/", h.AddUser)
			users.PUT("/:id", h.ReplaceUser)
			users.PATCH("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
		teams := api.Group("/teams")
		{
			teams.POST("/", handlers.AdminOnly(), h.CreateTeam)
			teams.GET("/", h.GetTeams)
			teams.DELETE("/:team_id", handlers.AdminOnly(), h.DeleteTeam)
			teams.POST("/:team_id/:user_id", h.AddUserToTeam)
			teams.DELETE("/:team_id/:user_id", h.RemoveUserFromTeam)
		}
	}

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
