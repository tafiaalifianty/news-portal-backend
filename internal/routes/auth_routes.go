package routes

import (
	"final-project-backend/internal/handlers"
	"final-project-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup, h *handlers.Handler) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/refresh", middlewares.AuthorizeRefreshToken, h.RefreshToken)
}
