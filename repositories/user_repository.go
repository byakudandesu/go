package repositories

import (
	"goapi/models"

	"gorm.io/gorm"
)

// all implamentations of this must have these methods
type UserRepository interface {
	// from the handler part, take the id, and do the search operation in db, return the user found.
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	Create(user *models.User) error
	Replace(user *models.User, id uint) (*models.User, error)
	Update(updates map[string]interface{}, id uint) (*models.User, error)
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

// Constructor function - this creates a new repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Now implement each method from the interface
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Replace(user *models.User, id uint) (*models.User, error) {
	user.ID = id
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(updates map[string]interface{}, id uint) (*models.User, error) {
	var user models.User
	// First find the user
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	// Then update it
	if err := r.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
