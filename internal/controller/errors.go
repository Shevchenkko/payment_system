package controller

import (
	// third party
	"github.com/gin-gonic/gin"
)

// response - represent error response.
type response struct {
	Error string `json:"error" example:"message"`
}

// errorResponse - responds with status code and error.
func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}
