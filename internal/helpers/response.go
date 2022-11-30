package helpers

import "github.com/gin-gonic/gin"

type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	IsError bool        `json:"is_error"`
}

func SendSuccessResponse(
	c *gin.Context,
	code int,
	message string,
	data interface{},
) {
	c.JSON(code, JsonResponse{
		Code:    code,
		Message: message,
		Data:    data,
		IsError: false,
	})
}

func SendErrorResponse(
	c *gin.Context,
	code int,
	message string,
) {
	c.AbortWithStatusJSON(code, JsonResponse{
		Code:    code,
		Message: message,
		IsError: true,
	})
}
