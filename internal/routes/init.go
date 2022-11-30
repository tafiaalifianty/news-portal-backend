package routes

import (
	"net/http"

	"final-project-backend/internal/handlers"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, h *handlers.Handler) {
	api := r.Group("/")
	{
		InitAuthRoutes(api, h)

		// Docs
		api.Static("/docs", "dist")

		// Ping Route
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		api.GET("/invoices/:code", h.GetWaitingInvoiceByCode)
		api.PATCH("/invoices/:code", h.UpdateWaitingInvoiceToProcessed)

		protected := api.Group("/", middlewares.AuthorizeJWT)
		{
			protected.GET("/ping-private", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong private",
				})
			})

			InitAdminRoutes(protected, h)
			InitMemberRoutes(protected, h)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		helpers.SendErrorResponse(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})

}
