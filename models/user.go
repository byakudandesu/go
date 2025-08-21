package models

type User struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Org    string `json:"org"`
	TeamID *uint  `json:"team_id"`
	Salary int    `json:"salary"`
}
