package controller

import (
	"fmt"
	"net/http"

	// third party
	"github.com/gin-gonic/gin"

	// external
	"github.com/Shevchenkko/payment_system/pkg/logger"

	// internal
	"github.com/Shevchenkko/payment_system/internal/service"
)

// adminRoutes - represents admin service router.
type adminRoutes struct {
	service service.Services
	repos   service.Repositories
	logger  logger.Interface
}

// newAdminRoutes - implements new admin service routes.
func newAdminRoutes(handler *gin.RouterGroup, s service.Services, l logger.Interface, repo service.Repositories) {
	r := &adminRoutes{s, repo, l}
	h := handler.Group("/admin")
	{
		// routes
		h.PATCH("/lock_user", newAuthMiddleware(s, l), r.lockUser)
		h.PATCH("/unlock_user", newAuthMiddleware(s, l), r.unlockUser)
	}
}

// lockUserRequestBody - represents lockUser request body.
type lockUserRequestBody struct {
	UserID int64 `json:"userId" binding:"required"`
}

// lockUserResponse - represents lockUser response.
type lockUserResponse struct {
	Status *string        `json:"status,omitempty"`
	Error  *service.Error `json:"error,omitempty"`
}

func (r *adminRoutes) lockUser(c *gin.Context) {
	logger := r.logger.Named("blockUser")

	// parse request body
	logger.Debug("parsing request body")
	var body lockUserRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// check role
	var status string
	if c.GetString("userRole") == "admin" {
		status, err = r.service.LockUser(c.Request.Context(), body.UserID, c.GetString("userRole"))
		if err != nil {
			logger.Error("failed to block user", "err", err)
			err, ok := err.(*service.Error)
			if ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, lockUserResponse{Error: err})
				return
			}
			errorResponse(c, http.StatusInternalServerError, "failed to block user")
			return
		}
	} else {
		logger.Error("failed to block user", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, lockUserResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusForbidden, "you need admin rights!")
		return
	}

	// get client
	client, err := r.repos.Users.GetUserByID(c.Request.Context(), int(body.UserID))
	if err != nil {
		return
	}

	logger.Info("successfully block user")
	c.JSON(http.StatusOK, lockUserResponse{Status: &status})

	var mess string
	if status == "ACTIVE" {
		mess = fmt.Sprintf("Successfully change status to %s for user %s", status, client.FullName)
	} else {
		mess = fmt.Sprintf("%s for user %s", status, client.FullName)
	}
	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: mess,
		})
	if err != nil {
		return
	}
}

// unlockUserRequestBody - represents unlockUserrequest body.
type unlockUserRequestBody struct {
	UserID int64 `json:"userId" binding:"required"`
}

// unlockUserResponse - represents unlockUser response.
type unlockUserResponse struct {
	Status *string        `json:"status,omitempty"`
	Error  *service.Error `json:"error,omitempty"`
}

func (r *adminRoutes) unlockUser(c *gin.Context) {
	logger := r.logger.Named("unlockUser")

	// parse request body
	logger.Debug("parsing request body")
	var body unlockUserRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// check role
	var status string
	if c.GetString("userRole") == "admin" {
		status, err = r.service.UnlockUser(c.Request.Context(), body.UserID, c.GetString("userRole"))
		if err != nil {
			logger.Error("failed to unlock user", "err", err)
			err, ok := err.(*service.Error)
			if ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, unlockUserResponse{Error: err})
				return
			}
			errorResponse(c, http.StatusInternalServerError, "failed to unlock user")
			return
		}
	} else {
		logger.Error("failed to unlock user", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, unlockUserResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusForbidden, "you need admin rights!")
		return
	}

	// get client
	client, err := r.repos.Users.GetUserByID(c.Request.Context(), int(body.UserID))
	if err != nil {
		return
	}

	logger.Info("successfully unlock user")
	c.JSON(http.StatusOK, unlockUserResponse{Status: &status})

	var mess string
	if status == "LOCK" {
		mess = fmt.Sprintf("Successfully change status to %s for user %s", status, client.FullName)
	} else {
		mess = fmt.Sprintf("%s for user %s", status, client.FullName)
	}
	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: mess,
		})
	if err != nil {
		return
	}
}
