package handlers

import (
	"goapi/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Returns everything stored
func (h *Handler) GetUsers(request *gin.Context) {
	users, err := h.UserRepo.GetAll()
	if err != nil {
		request.JSON(400, gin.H{"error": "Could not find users"})
		return
	}
	request.JSON(200, users)
}

// Returns single user by ID from database
func (h *Handler) GetUser(request *gin.Context) {
	id := request.Param("id")
	
	// Convert string ID to uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	
	user, err := h.UserRepo.GetByID(uint(userID))
	if err != nil {
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
	// add new user to database using repository
	if err := h.UserRepo.Create(&newUser); err != nil {
		request.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	request.JSON(201, newUser)
}

func (h *Handler) ReplaceUser(request *gin.Context) {
	var userToUpdate models.User
	id := request.Param("id")
	
	// Convert string ID to uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	existingUser, err := h.UserRepo.GetByID(uint(userID))
	if err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := request.ShouldBindJSON(&userToUpdate); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Keep the same ID
	userToUpdate.ID = existingUser.ID

	// Replace user using repository
	updatedUser, err := h.UserRepo.Replace(&userToUpdate, uint(userID))
	if err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	request.JSON(200, updatedUser)
}

func (h *Handler) UpdateUser(request *gin.Context) {
	// interface instead of the struct, since we may have values that are missing
	var updates map[string]interface{}
	id := request.Param("id")
	
	// Convert string ID to uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	_, err = h.UserRepo.GetByID(uint(userID))
	if err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := request.ShouldBindJSON(&updates); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update using repository
	updatedUser, err := h.UserRepo.Update(updates, uint(userID))
	if err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	request.JSON(200, updatedUser)
	log.Printf("Updated user: %v", updatedUser)
}

func (h *Handler) DeleteUser(request *gin.Context) {
	id := request.Param("id")
	
	// Convert string ID to uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete user using repository
	err = h.UserRepo.Delete(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			request.JSON(404, gin.H{"error": "User not found"})
		} else {
			request.JSON(500, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	request.JSON(204, nil)
}
