package handlers

import (
	"goapi/models"
	"log"

	"github.com/gin-gonic/gin"
)

// Returns everything stored
func (h *Handler) GetUsers(request *gin.Context) {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		request.JSON(400, gin.H{"error": "Could not find users"})
	}
	request.JSON(200, users)

}

// Returns single user by ID from database
func (h *Handler) GetUser(request *gin.Context) {
	var user models.User
	id := request.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	request.JSON(200, user)
}

// adding user to databases
func (h *Handler) AddUser(request *gin.Context) {
	// need to take given json from user and format so the code can read
	var newUser models.User

	// if the inputted data can be formatted to the struct we made
	if err := request.ShouldBindJSON(&newUser); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// this would ignore the salary calculation, so the user must be added through the addUserToTeam function
	newUser.TeamID = nil
	// add new user to database
	if err := h.DB.Create(&newUser).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	request.JSON(201, newUser)
}

func (h *Handler) ReplaceUser(request *gin.Context) {
	var userToUpdate models.User
	id := request.Param("id")

	// Check if user exists
	var existingUser models.User
	if err := h.DB.First(&existingUser, id).Error; err != nil {
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
	if err := h.DB.Save(&userToUpdate).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	request.JSON(200, userToUpdate)
}

func (h *Handler) UpdateUser(request *gin.Context) {
	// interface instead of the struct, since we may have values that are missing
	var updates map[string]interface{}
	id := request.Param("id")

	// Check if user exists
	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := request.ShouldBindJSON(&updates); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	// Fetch updated user
	h.DB.First(&user, id)
	request.JSON(200, user)
	log.Printf("Updated user: %v", user)
}

func (h *Handler) DeleteUser(request *gin.Context) {
	id := request.Param("id")

	// Delete user from database
	result := h.DB.Delete(&models.User{}, id)
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
