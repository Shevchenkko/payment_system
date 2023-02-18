// Package controller implements application http delivery.
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

// userRoutes - represents user service router.
type userRoutes struct {
	service service.Users
	logger  logger.Interface
}

// newUserRoutes - implements new user service routes.
func newUserRoutes(handler *gin.RouterGroup, s service.Users, l logger.Interface) {
	r := &userRoutes{s, l}
	h := handler.Group("/users")
	{
		// routes
		h.POST("/register", r.registerUser)
		h.POST("/login", r.loginUser)
		h.POST("/sendemail", r.sendEmail)
		h.PATCH("/resetpassword", r.resetPassword)
	}
}

// registerUserRequestBody - represents registerUser request body.
type registerUserRequestBody struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// registerUserResponse - represents registerUser response.
type registerUserResponse struct {
	Token    string         `json:"token"`
	UserID   int            `json:"userId"`
	FullName string         `json:"fullName" binding:"required"`
	Email    string         `json:"email" binding:"required"`
	Error    *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) registerUser(c *gin.Context) {
	logger := r.logger.Named("registerUser")

	// parse request body
	logger.Debug("parsing request body")
	var body registerUserRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// register user
	logger.Debug("registering user")
	registerData, err := r.service.RegisterUser(c.Request.Context(),
		&service.RegisterUserInput{
			FullName: body.FullName,
			Email:    body.Email,
			Password: body.Password,
		})
	if err != nil {
		logger.Error("failed to register user", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, registerUserResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to register user")
		return
	}

	logger = logger.With("registerData", registerData)
	logger.Info("successfully registered in")
	c.JSON(http.StatusOK, registerUserResponse{
		Token:    registerData.Token,
		UserID:   registerData.UserID,
		FullName: registerData.FullName,
		Email:    registerData.Email,
	})
}

// loginUserRequestBody - represents loginUser request body.
type loginUserRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// loginUserResponse - represents login response.
type loginUserResponse struct {
	Token    string         `json:"token"`
	UserID   int            `json:"userId"`
	FullName string         `json:"fullName" binding:"required"`
	Email    string         `json:"email" binding:"required"`
	Error    *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) loginUser(c *gin.Context) {
	logger := r.logger.Named("loginUser")

	// parse request body
	logger.Debug("parsing request body")
	var body loginUserRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// login user
	logger.Debug("loginning user")
	loginData, err := r.service.LoginUser(c.Request.Context(),
		&service.LoginUserInput{
			Email:    body.Email,
			Password: body.Password,
		})
	if err != nil {
		logger.Error("failed to login user", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, loginUserResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to login user")
		return
	}

	logger.Info("successfully logged in")
	c.JSON(http.StatusOK, loginUserResponse{
		Token:    loginData.Token,
		UserID:   loginData.UserID,
		FullName: loginData.FullName,
		Email:    loginData.Email,
	})
}

// sendEmailRequestBody - represents sendEmail request body.
type sendEmailRequestBody struct {
	Email string `json:"email" binding:"required"`
}

// sendEmailResponse - represents send email response.
type sendEmailResponse struct {
	Error *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) sendEmail(c *gin.Context) {
	logger := r.logger.Named("sendEmail")

	// parse request body
	logger.Debug("parsing request body")
	var body sendEmailRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// send email
	logger.Debug("sending email")
	err = r.service.SendEmail(c.Request.Context(),
		&service.SendUserEmailInput{
			Email: body.Email,
		})
	if err != nil {
		logger.Error("failed to send email", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, sendEmailResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info("successfully sended email")
	c.JSON(http.StatusOK, sendEmailResponse{})
}

// resetPasswordRequestBody - represents resetPassword request body.
type resetPasswordRequestBody struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// resetPasswordResponse - represents resetPassword response.
type resetPasswordResponse struct {
	Error *service.Error `json:"error,omitempty"`
}

func (r *userRoutes) resetPassword(c *gin.Context) {
	logger := r.logger.Named("resetPassword")

	// parse request body
	logger.Debug("parsing request body")
	var body resetPasswordRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// reset password
	logger.Debug("resetting password")
	err = r.service.ResetPassword(c.Request.Context(),
		&service.ResetPasswordInput{
			Token:    body.Token,
			Password: body.Password,
		})
	if err != nil {
		logger.Error("failed to reset password", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, resetPasswordResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to reset password")
		return
	}

	logger.Info("successfully resetted password")
	c.JSON(http.StatusOK, resetPasswordResponse{})
}
