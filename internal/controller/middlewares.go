// Package controller implements application http delivery.
package controller

import (
	"net/http"
	"strings"

	// third party
	"github.com/Shevchenkko/payment_system/internal/service"
	"github.com/Shevchenkko/payment_system/pkg/logger"
	"github.com/gin-gonic/gin"
)

// corsMiddleware - used to allow incoming cross-origin requests.
func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

// newAuthMiddleware is used to get auth token from request headers and validate it.
func newAuthMiddleware(services service.Services, l logger.Interface) gin.HandlerFunc {
	logger := l.Named("authMiddleware")

	return func(c *gin.Context) {
		// get token and check if empty ("Bearer token")
		tokenStringRaw := c.GetHeader("Authorization")
		if tokenStringRaw == "" {
			logger.Debug("empty Authorization header", "tokenStringRaw", tokenStringRaw)
			errorResponse(c, http.StatusUnauthorized, "empty auth token")
			return
		}

		// split Bearer and token
		tokenStringArr := strings.Split(tokenStringRaw, " ")
		if len(tokenStringArr) != 2 {
			logger.Debug("malformed auth token", "tokenStringArr", tokenStringArr)
			errorResponse(c, http.StatusUnauthorized, "malformed auth token")
			return
		}

		// get token
		tokenString := tokenStringArr[1]
		valid, client, userRole := services.Users.VerifyAccessToken(c.Request.Context(), tokenString)
		if !valid {
			logger.Debug("invalid auth token", "tokenStringArr", tokenStringArr)
			errorResponse(c, http.StatusUnauthorized, "invalid auth token")
			return
		}

		// set user id to context
		c.Set("clientID", client)

		// set user role to context
		c.Set("userRole", userRole)
	}
}
