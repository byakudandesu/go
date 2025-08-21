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

	if err := h.DB.Create(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to create team"})
		return
	}

	request.JSON(201, team)
}

func (h *Handler) GetTeams(request *gin.Context) {
	var teams []models.Team
	if err := h.DB.Preload("Users").Find(&teams).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to fetch teams"})
		return
	}
	request.JSON(200, teams)
}

func (h *Handler) AddUserToTeam(request *gin.Context) {
	// which team? which user?
	teamID := request.Param("team_id")
	userID := request.Param("user_id")

	// finding requested teams and users into go varibles, if they don't exist, error
	var team models.Team
	if err := h.DB.First(&team, teamID).Error; err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
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

	if err := h.DB.Save(&user).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}
	if err := h.DB.Save(&team).Error; err != nil {
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

	var team models.Team
	if err := h.DB.First(&team, teamID).Error; err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		request.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// changing the teamid for said user to nil so there is no relationship
	user.TeamID = nil
	team.UsedBudget -= user.Salary

	if err := h.DB.Save(&user).Error; err != nil {
		request.JSON(500, gin.H{"error": "Error saving to database"})
		return
	}
	if err := h.DB.Save(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Error saving to database"})
		return
	}

	request.JSON(200, gin.H{"message": "User removed from team"})
}

func (h *Handler) DeleteTeam(request *gin.Context) {
	teamID := request.Param("team_id")

	var team models.Team

	teamIDInt, _ := strconv.Atoi(teamID)
	theteamID := uint(teamIDInt)

	if err := h.DB.First(&team, teamID).Error; err != nil {
		request.JSON(404, gin.H{"error": "Team not found"})
		return
	}

	// find all the users where the team id is the teamid in params and set them to null before team deletion
	if err := h.DB.Model(&models.User{}).Where("team_id = ?", theteamID).Update("team_id",
		nil).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to remove users from team"})
		return
	}

	if err := h.DB.Delete(&team).Error; err != nil {
		request.JSON(500, gin.H{"error": "Failed to delete team"})
		return
	}

	request.JSON(200, gin.H{
		"message": "Team deleted and users freed",
		"team":    team.Name,
	})

}
