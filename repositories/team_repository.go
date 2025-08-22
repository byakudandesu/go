package repositories

import (
	"goapi/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	GetAll() ([]models.Team, error)
	GetAllWithUsers() ([]models.Team, error)
	GetByID(id uint) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
	RemoveAllUsersFromTeam(teamID uint) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{
		db: db,
	}
}

func (r *teamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) GetAll() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *teamRepository) GetAllWithUsers() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Preload("Users").Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *teamRepository) GetByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}

func (r *teamRepository) RemoveAllUsersFromTeam(teamID uint) error {
	return r.db.Model(&models.User{}).Where("team_id = ?", teamID).Update("team_id", nil).Error
}
