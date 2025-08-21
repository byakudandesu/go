package models

type Team struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Budget     int    `json:"budget"`
	UsedBudget int    `json:"used_budget"`
	Users      []User `json:"users,omitempty" gorm:"foreignKey:TeamID"`
}
