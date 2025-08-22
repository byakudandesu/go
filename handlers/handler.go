package handlers

import (
	"goapi/repositories"
	"goapi/usecases"
)

type Handler struct {
	UserRepo    repositories.UserRepository
	TeamRepo    repositories.TeamRepository
	TeamUseCase usecases.TeamUseCase
}
