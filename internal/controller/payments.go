package controller

import (
	"fmt"
	"net/http"

	// third party
	"github.com/gin-gonic/gin"

	// external
	"github.com/Shevchenkko/payment_system/pkg/logger"

	// internal
	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
)

// paymentRoutes - represents payment service router.
type paymentRoutes struct {
	service service.Services
	repos   service.Repositories
	logger  logger.Interface
}

// newPaymentRoutes - implements new payment service routes.
func newPaymentRoutes(handler *gin.RouterGroup, s service.Services, l logger.Interface, repo service.Repositories) {
	r := &paymentRoutes{s, repo, l}
	h := handler.Group("/payment")
	{
		// routes
		h.GET("/search", newAuthMiddleware(s, l), r.searchPayment)
		h.POST("/create", newAuthMiddleware(s, l), r.createPayment)
		h.PATCH("/sent", newAuthMiddleware(s, l), r.sentPayment)
	}
}

// searchPaymentRequestQuery - represents search payments request query.
type searchPaymentRequestQuery struct {
	Filter domain.Filter `form:"filter"`
}

// searchPaymentResponse - represents search payments response.
type searchPaymentResponse struct {
	Data       []service.PaymentOutput `json:"data"`
	Pagination *domain.Pagination      `json:"pagination"`

	Error *service.Error `json:"error,omitempty"`
}

func (r *paymentRoutes) searchPayment(c *gin.Context) {
	logger := r.logger.Named("searchPayment")

	filter, err := getFilterFromQuery(c.Request)
	if err != nil {
		logger.Error("failed to parse query params", "err", err)
		errorResponse(c, http.StatusBadRequest, "failed to parse query params")
		return
	}

	// parse request query
	var query searchPaymentRequestQuery
	logger.Info("parsing request query")
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("failed to parse request query", "err", err)
		errorResponse(c, http.StatusBadRequest, "failed to parse request query")
		return
	}

	// get client
	client, err := r.repos.Users.GetUserByID(c.Request.Context(), c.GetInt("clientID"))
	if err != nil {
		return
	}
	if client.Status == "LOCK" {
		errorResponse(c, http.StatusInternalServerError, "Your account is blocked! Please, turn to the nearest branch of our bank")
		return
	}

	response, err := r.service.Payments.SearchPayments(c.Request.Context(), filter, client.FullName)
	if err != nil {
		logger.Error("failed to search payments", "err", err)
		// get service error
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, searchPaymentResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to search payments")
		return
	}
	logger = logger.With("search payment", response)
	logger.Debug("got payments")

	logger.Info("successfully search payments")
	c.JSON(http.StatusOK, searchPaymentResponse{
		Data:       response.Data,
		Pagination: response.Pagination,
	})
}

// createPaymentRequestBody - represents createPayment request body.
type createPaymentRequestBody struct {
	FromClientIBAN  string  `json:"fromClientIban" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	ToClientIBAN    string  `json:"toClientIban" binding:"required"`
	ToClient        string  `json:"toClient" binding:"required"`
	OperationAmount float64 `json:"operationAmount" binding:"required"`
}

// createPaymentResponse - represents createPayment response.
type createPaymentResponse struct {
	PaymentID            int64          `json:"paymentId"`
	PaymentStatus        string         `json:"paymentStatus"`
	FromClient           string         `json:"fromClient"`
	FromClientITN        int64          `json:"fromClientItn"`
	FromClientIBAN       string         `json:"fromClientIban"`
	FromClientCardNumber int64          `json:"fromClientCardNumber"`
	Description          string         `json:"description"`
	ToClientIBAN         string         `json:"toClientIban"`
	ToClient             string         `json:"toClient"`
	OperationAmount      float64        `json:"operationAmount"`
	Error                *service.Error `json:"error,omitempty"`
}

func (r *paymentRoutes) createPayment(c *gin.Context) {
	logger := r.logger.Named("createPayment")

	// parse request body
	logger.Debug("parsing request body")
	var body createPaymentRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// check user
	userStatus, err := r.repos.Users.GetUserByID(c.Request.Context(), c.GetInt("clientID"))
	if err != nil {
		return
	}
	if userStatus.Status == "LOCK" {
		errorResponse(c, http.StatusInternalServerError, "Your account is blocked! Please, turn to the nearest branch of our bank")
		return
	}

	// check pperation amount
	var operationAmount float64
	if body.OperationAmount > 0 {
		operationAmount = body.OperationAmount
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to top up the bank account, please enter the correct amount")
		return
	}

	// create payment for client
	logger.Debug("creating payment for client")
	data, err := r.service.CreatePayment(c.Request.Context(), c.GetInt("clientID"),
		&service.PaymentInput{
			FromClientIBAN:  body.FromClientIBAN,
			Description:     body.Description,
			ToClientIBAN:    body.ToClientIBAN,
			ToClient:        body.ToClient,
			OperationAmount: operationAmount,
		})
	if err != nil {
		logger.Error("failed to create payment", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, createPaymentResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to create payment")
		return
	}

	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: fmt.Sprintf("Successfully created payment from %s to %s", data.FromClient, data.ToClient),
		})
	if err != nil {
		return
	}

	logger = logger.With("create data", data)
	logger.Info("successfully created payment")
	c.JSON(http.StatusOK, createPaymentResponse{
		PaymentID:            data.ID,
		PaymentStatus:        data.PaymentStatus,
		FromClient:           data.FromClient,
		FromClientITN:        data.FromClientITN,
		FromClientIBAN:       data.FromClientIBAN,
		FromClientCardNumber: data.FromClientCardNumber,
		Description:          data.Description,
		ToClientIBAN:         data.ToClientIBAN,
		ToClient:             data.ToClient,
		OperationAmount:      data.OperationAmount,
	})
}

// sentPaymentRequestBody - represents setPayment request body.
type sentPaymentRequestBody struct {
	PaymentID   int64  `json:"paymentId" binding:"required"`
	SecretValue string `json:"secretValue" binding:"required"`
}

// sentPaymentResponse - represents sentPayment response.
type sentPaymentResponse struct {
	Status string         `json:"status"`
	Error  *service.Error `json:"error,omitempty"`
}

func (r *paymentRoutes) sentPayment(c *gin.Context) {
	logger := r.logger.Named("sentPayment")

	// parse request body
	logger.Debug("parsing request body")
	var body sentPaymentRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// check user
	userStatus, err := r.repos.Users.GetUserByID(c.Request.Context(), c.GetInt("clientID"))
	if err != nil {
		return
	}
	if userStatus.Status == "LOCK" {
		errorResponse(c, http.StatusInternalServerError, "Your account is blocked! Please, turn to the nearest branch of our bank")
		return
	}

	// get payment
	payment, err := r.repos.Payments.GetPaymentByID(c.Request.Context(), body.PaymentID)
	if err != nil {
		return
	}

	// check bank account
	bank, err := r.repos.Banks.GetInfoByIBAN(c.Request.Context(), payment.FromClientIBAN)
	if err != nil {
		return
	}

	// check card status
	if bank.Status == "LOCK" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to sent payment, please unlock your bank account")
		return
	}

	var cardBalance float64
	if bank.Balance-payment.OperationAmount > 0 {
		cardBalance = bank.Balance - payment.OperationAmount
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to sent payment, please top up your balance")
		return
	}

	// sent payment for client
	logger.Debug("senting payment for client")
	data, err := r.service.SentPayment(c.Request.Context(), body.PaymentID, body.SecretValue, cardBalance)
	if err != nil {
		logger.Error("failed to create payment", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, sentPaymentResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to create payment")
		return
	}

	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: fmt.Sprintf("Successfully sent payment #%d", body.PaymentID),
		})
	if err != nil {
		return
	}

	logger = logger.With("sent data", data)
	logger.Info("successfully send payment")
	c.JSON(http.StatusOK, sentPaymentResponse{Status: data})
}
