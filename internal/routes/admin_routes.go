package routes

import (
	"net/http"

	"final-project-backend/internal/handlers"
	"final-project-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitAdminRoutes(r *gin.RouterGroup, h *handlers.Handler) {
	admin := r.Group("/admin", middlewares.AuthorizeAdmin)
	{
		posts := admin.Group("/posts")
		{
			posts.POST("", h.CreatePost)
			posts.DELETE("/:id", h.DeletePost)
		}

		invoices := admin.Group("invoices")
		{
			invoices.GET("", h.GetAllInvoices)
			invoices.PATCH("/:code", h.UpdateProcessedInvoice)
		}

		vouchers := admin.Group("vouchers")
		{
			vouchers.GET("", h.GetAllVouchers)
		}

		gifts := admin.Group("gifts")
		{
			gifts.GET("", h.GetAllGifts)
			gifts.GET("/users", h.GetAllUserGifts)
			gifts.PATCH("/users/:id", h.UpdateUserGiftStatus)
			gifts.PATCH("/stock/:id", h.UpdateGiftStock)
		}

		admin.GET("/ping-admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong admin",
			})
		})
	}
}
