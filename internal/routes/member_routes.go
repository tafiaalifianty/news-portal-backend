package routes

import (
	"final-project-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func InitMemberRoutes(r *gin.RouterGroup, h *handlers.Handler) {
	users := r.Group("/users")
	{
		users.GET("/profile", h.GetUserProfile)
		users.GET("/histories", h.GetUserHistories)
		users.GET("/subscriptions", h.GetUserSubscriptions)
		users.GET("/invoices", h.GetUserInvoices)
		users.GET("/gifts", h.GetUserGifts)
		users.GET("/referrals", h.GetUserReferrals)
		users.GET("/vouchers", h.GetUserVouchers)
		users.PATCH("/profile", h.UpdateUser)
	}
	posts := r.Group("/posts")
	{
		posts.GET("", h.GetAllPosts)
		posts.GET("/:id", h.GetPostByID)
		posts.GET("/recommendations", h.GetRecommendedPost)
		posts.GET("/trending", h.GetTrendingPosts)
		posts.GET("/types", h.GetAllTypes)
		posts.GET("/categories", h.GetAllCategories)
		posts.PATCH("/like/:id", h.LikePost)
		posts.PATCH("/share/:id", h.SharePost)
	}
	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.GET("", h.GetAllSubscriptions)
		subscriptions.POST("", h.AddUserSubscription)
	}
	invoices := r.Group("/invoices")
	{
		invoices.POST("", h.CreateInvoice)
		invoices.GET("/user/:code", h.GetWaitingInvoiceByCodeProtected)
	}
}
