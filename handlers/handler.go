package handlers

import (
	"goapi/repositories"
)

type Handler struct {
	UserRepo repositories.UserRepository
	TeamRepo repositories.TeamRepository
}
