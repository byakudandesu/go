package usecases

import (
	"goapi/models"
	"goapi/repositories"
)

type TeamUseCase interface {
	GetAllTeamsIncludingFreeAgents() ([]models.Team, error)
}

type teamUseCase struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func NewTeamUseCase(teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) TeamUseCase {
	return &teamUseCase{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (uc *teamUseCase) GetAllTeamsIncludingFreeAgents() ([]models.Team, error) {
	teams, err := uc.teamRepo.GetAllWithUsers()
	if err != nil {
		return nil, err
	}

	freeAgents, err := uc.userRepo.FreeAgents()
	if err != nil {
		return nil, err
	}

	freeAgentTeam := models.Team{
		ID:         0,
		Name:       "Free Agents",
		Budget:     0,
		UsedBudget: 0,
		Users:      freeAgents,
	}
	teams = append(teams, freeAgentTeam)
	return teams, nil
}
