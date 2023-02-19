package controller

import (
	"fmt"
	"net/http"

	"github.com/Shevchenkko/payment_system/internal/service"
	"github.com/Shevchenkko/payment_system/pkg/logger"
	"github.com/gin-gonic/gin"
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
		h.POST("/create", newAuthMiddleware(s, l), r.createPayment)
		h.PATCH("/sent", newAuthMiddleware(s, l), r.sentPayment)
	}
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
	data, err := r.service.CreatePayment(c.Request.Context(), c.GetInt("client"),
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

	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("client"),
		&service.MessageLogInput{
			MessageLog: fmt.Sprintf("Successfully created payment from %s to %s", data.FromClient, data.ToClient),
		})
	if err != nil {
		return
	}

	logger = logger.With("create data", data)
	logger.Info("successfully created payment")
	c.JSON(http.StatusOK, createPaymentResponse{
		PaymentID:            data.PaymentID,
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
	PaymentID   int    `json:"paymentId" binding:"required"`
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

	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("client"),
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
