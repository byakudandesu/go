package handlers

import "gorm.io/gorm"

// Handler groups all HTTP handler methods together with their dependencies the dependency is the db
// Why? Before handlers were functions that couldn't access the database
// Now they're methods that can access h.DB
//
// Usage in main.go:
//
//	h := &handlers.Handler{DB: db}
//	router.GET("/users", h.GetUsers)
type Handler struct {
	DB *gorm.DB
}
