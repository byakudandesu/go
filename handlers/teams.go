package handlers

import (
	"goapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// team stuff
func (h *Handler) CreateTeam(request *gin.Context) {
	var team models.Team

	if err := request.ShouldBindJSON(&team); err != nil {
		request.JSON(400, gin.H{"error": err.Error()})
		return
	}
	team.UsedBudget = 0

	if err := h.TeamRepo.Create(&team); err != nil {
		request.JSON(500, gin.H{"error": "Failed to create team"})
		return
	}

	request.JSON(201, team)
}

func (h *Handler) GetTeams(request *gin.Context) {
	// Use the repository method that preloads users
	teams, err := h.TeamRepo.GetAllWithUsers()
	if err != nil {
		request.JSON(500, gin.H{"error": "Failed to fetch teams"})
		return
	}
	request.JSON(200, teams)
}

func (h *Handler) AddUserToTeam(request *gin.Context) {
	// which team? which user?
	teamID := request.Param("team_id")
	userID := request.Param("user_id")

	// Convert string IDs to uint
	teamIDUint, err := strconv.ParseUint(teamID, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid team ID"})
		return
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// finding requested teams and users using repositories
	team, err := h.TeamRepo.GetByID(uint(teamIDUint))
	if err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}
	user, err := h.UserRepo.GetByID(uint(userIDUint))
	if err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// check if user is not in a team yet
	if user.TeamID != nil {
		request.JSON(400, gin.H{"error": "User already belongs in team"})
		return
	}

	// need to check if budget of team can afford user
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

	// assign team ID to user
	teamIDValue := uint(teamIDUint)
	user.TeamID = &teamIDValue

	team.UsedBudget = newBudgetUsed

	// Update user and team using repositories
	_, err = h.UserRepo.Replace(user, user.ID)
	if err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}
	err = h.TeamRepo.Update(team)
	if err != nil {
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
func (h *Handler) RemoveUserFromTeam(request *gin.Context) {
	teamID := request.Param("team_id")
	userID := request.Param("user_id")

	// Convert string IDs to uint
	teamIDUint, err := strconv.ParseUint(teamID, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid team ID"})
		return
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get team and user using repositories
	team, err := h.TeamRepo.GetByID(uint(teamIDUint))
	if err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}

	user, err := h.UserRepo.GetByID(uint(userIDUint))
	if err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// changing the teamid for said user to nil so there is no relationship
	user.TeamID = nil
	team.UsedBudget -= user.Salary

	// Update user and team using repositories
	_, err = h.UserRepo.Replace(user, user.ID)
	if err != nil {
		request.JSON(500, gin.H{"error": "Error saving to database"})
		return
	}
	err = h.TeamRepo.Update(team)
	if err != nil {
		request.JSON(500, gin.H{"error": "Error saving to database"})
		return
	}

	request.JSON(200, gin.H{"message": "User removed from team"})
}

func (h *Handler) DeleteTeam(request *gin.Context) {
	teamID := request.Param("team_id")

	// Convert string ID to uint
	teamIDUint, err := strconv.ParseUint(teamID, 10, 32)
	if err != nil {
		request.JSON(400, gin.H{"error": "Invalid team ID"})
		return
	}

	// Check if team exists
	team, err := h.TeamRepo.GetByID(uint(teamIDUint))
	if err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}

	// Remove all users from the team before deletion
	err = h.TeamRepo.RemoveAllUsersFromTeam(uint(teamIDUint))
	if err != nil {
		request.JSON(500, gin.H{"error": "Failed to remove users from team"})
		return
	}

	// Delete the team
	err = h.TeamRepo.Delete(uint(teamIDUint))
	if err != nil {
		request.JSON(500, gin.H{"error": "Failed to delete team"})
		return
	}

	request.JSON(200, gin.H{
		"message": "Team deleted and users freed",
		"team":    team.Name,
	})

}
