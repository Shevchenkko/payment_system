package controller

import (
	"net/http"

	// third party
	"github.com/gin-gonic/gin"

	// external
	"github.com/Shevchenkko/payment_system/pkg/logger"

	// internal
	"github.com/Shevchenkko/payment_system/internal/service"
)

// NewRouter - represents application router.
func NewRouter(handler *gin.Engine, s service.Services, l logger.Interface, r service.Repositories) {
	// options
	handler.Use(gin.Logger(), gin.Recovery(), corsMiddleware)

	// k8s probe
	handler.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	// routers
	h := handler.Group("/api/v1")
	{
		newUserRoutes(h, s, l, r)
		newBankAccountRoutes(h, s, l, r)
		newPaymentRoutes(h, s, l, r)
		newAdminRoutes(h, s, l, r)
	}
}
