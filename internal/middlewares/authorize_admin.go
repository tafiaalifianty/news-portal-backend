package middlewares

import (
	"net/http"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func AuthorizeAdmin(c *gin.Context) {
	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	var user dtos.JwtData = userContext.(dtos.JwtData)

	if user.Role != string(models.Admin) {
		helpers.SendErrorResponse(
			c,
			http.StatusForbidden,
			http.StatusText(http.StatusForbidden),
		)
		return
	} else {
		c.Next()
	}
}
