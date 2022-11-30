package handlers

import "final-project-backend/internal/services"

type Handler struct {
	services *services.Services
}

func New(s *services.Services) *Handler {
	return &Handler{services: s}
}
