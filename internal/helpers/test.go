package helpers

import (
	"encoding/json"
	"strings"

	"final-project-backend/internal/dtos"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func MakeRequestBody(dto interface{}) *strings.Reader {
	payload, _ := json.Marshal(dto)
	return strings.NewReader(string(payload))
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

func MiddlewareMockUser(mockUserContext dtos.JwtData) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user", mockUserContext)
		c.Next()
	}
}

func MiddlewareMockID(mockID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Params = []gin.Param{{Key: "id", Value: mockID}}
	}
}

func MiddlewareMockParams(mockKey string, mockValue string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Params = []gin.Param{{Key: mockKey, Value: mockValue}}
	}
}
