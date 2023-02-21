package controller

import (
	"fmt"
	"net/http"

	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/service"
	"github.com/Shevchenkko/payment_system/pkg/logger"
	"github.com/gin-gonic/gin"
)

// bankAccountRoutes - represents bank account service router.
type bankAccountRoutes struct {
	service service.Services
	repos   service.Repositories
	logger  logger.Interface
}

// newBankAccountRoutes - implements new bank account service routes.
func newBankAccountRoutes(handler *gin.RouterGroup, s service.Services, l logger.Interface, repo service.Repositories) {
	r := &bankAccountRoutes{s, repo, l}
	h := handler.Group("/bank_account")
	{
		// routes
		h.GET("/search", newAuthMiddleware(s, l), r.searchBankAccount)
		h.POST("/create", newAuthMiddleware(s, l), r.createBankAccount)
		h.PATCH("/top_up", newAuthMiddleware(s, l), r.topUpBankAccount)
		h.PATCH("/lock", newAuthMiddleware(s, l), r.lockBankAccount)
		h.PATCH("/unlock", newAuthMiddleware(s, l), r.unlockBankAccount)
	}
}

// searchBankAccountRequestQuery - represents search bank account request query.
type searchBankAccountRequestQuery struct {
	Filter domain.Filter `form:"filter"`
}

// searchBankAccountResponse - represents search bank account response.
type searchBankAccountResponse struct {
	Data       []service.BankAccountOutput `json:"data"`
	Pagination *domain.Pagination          `json:"pagination"`

	Error *service.Error `json:"error,omitempty"`
}

func (r *bankAccountRoutes) searchBankAccount(c *gin.Context) {
	logger := r.logger.Named("searchBankAccount")

	filter, err := getFilterFromQuery(c.Request)
	if err != nil {
		logger.Error("failed to parse query params", "err", err)
		errorResponse(c, http.StatusBadRequest, "failed to parse query params")
		return
	}

	// parse request query
	var query searchBankAccountRequestQuery
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

	response, err := r.service.BankAccounts.SearchBankAccounts(c.Request.Context(), filter, client.FullName, c.GetString("userRole"))
	if err != nil {
		logger.Error("failed to search bank accounts", "err", err)
		// get service error
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, searchBankAccountResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to search bank accounts")
		return
	}
	logger = logger.With("search bank account", response)
	logger.Debug("got bank account")

	logger.Info("successfully search bank account")
	c.JSON(http.StatusOK, searchBankAccountResponse{
		Data:       response.Data,
		Pagination: response.Pagination,
	})
}

// createBankAccountRequestBody - represents createBankAccount request body.
type createBankAccountRequestBody struct {
	ITN         int64  `json:"itn" binding:"required"`
	SecretValue string `json:"secretValue" binding:"required"`
}

// createBankAccountResponse - represents createBankAccount response.
type createBankAccountResponse struct {
	ID         int            `json:"id"`
	Client     string         `json:"client"`
	CardNumber int64          `json:"cardNumber"`
	IBAN       string         `json:"iban"`
	Balance    float64        `json:"balance"`
	Error      *service.Error `json:"error,omitempty"`
}

func (r *bankAccountRoutes) createBankAccount(c *gin.Context) {
	logger := r.logger.Named("createBankAccount")

	// parse request body
	logger.Debug("parsing request body")
	var body createBankAccountRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// create bank account for client
	logger.Debug("creating bank account for client")
	data, err := r.service.CreateBankAccount(c.Request.Context(), c.GetInt("clientID"),
		&service.BankAccountInput{
			ITN:         body.ITN,
			SecretValue: body.SecretValue,
		})
	if err != nil {
		logger.Error("failed to create bank account", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, createBankAccountResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to create bank account")
		return
	}

	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: fmt.Sprintf("Successfully created bank account %d", data.CardNumber),
		})
	if err != nil {
		return
	}

	logger = logger.With("create data", data)
	logger.Info("successfully created bank account")
	c.JSON(http.StatusOK, createBankAccountResponse{
		ID:         data.ID,
		Client:     data.Client,
		CardNumber: data.CardNumber,
		IBAN:       data.IBAN,
		Balance:    data.Balance,
	})
}

// topUpBankAccountRequestBody - represents topUpBankAccount request body.
type topUpBankAccountRequestBody struct {
	CardNumber      int64   `json:"cardNumber" binding:"required"`
	OperationAmount float64 `json:"operationAmount" binding:"required"`
}

// topUpBankAccountResponse - represents topUpBankAccount response.
type topUpBankAccountResponse struct {
	Client     string         `json:"client"`
	CardNumber int64          `json:"cardNumber"`
	IBAN       string         `json:"iban"`
	Balance    float64        `json:"balance"`
	Error      *service.Error `json:"error,omitempty"`
}

func (r *bankAccountRoutes) topUpBankAccount(c *gin.Context) {
	logger := r.logger.Named("topUpBankAccount")

	// parse request body
	logger.Debug("parsing request body")
	var body topUpBankAccountRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	var cardBalance float64
	if body.OperationAmount > 0 {
		cardBalance = body.OperationAmount
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to top up the bank account, please enter the correct amount")
		return
	}

	// check card status
	status, err := r.repos.Banks.CheckCreditCard(c.Request.Context(), body.CardNumber)
	if err != nil {
		return
	}

	if status.Status == "LOCK" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to top up the bank account, please unlock your bank account")
		return
	}

	// top up bank account for client
	logger.Debug("top up bank account for client")
	data, err := r.service.TopUpBankAccount(c.Request.Context(), c.GetInt("clientID"),
		&service.TopUpBankAccountInput{
			OperationAmount: cardBalance,
			CardNumber:      body.CardNumber,
		})
	if err != nil {
		logger.Error("failed to top up bank account", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, topUpBankAccountResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to toping up bank account")
		return
	}

	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: fmt.Sprintf("Successfully toping up bank account %d on operation amount %0.2f", data.CardNumber, body.OperationAmount),
		})
	if err != nil {
		return
	}

	logger = logger.With("top up data", data)
	logger.Info("successfully toping up bank account")
	c.JSON(http.StatusOK, topUpBankAccountResponse{
		Client:     data.Client,
		CardNumber: data.CardNumber,
		IBAN:       data.IBAN,
		Balance:    data.Balance,
	})
}

// lockBankAccountRequestBody - represents lockBankAccount request body.
type lockBankAccountRequestBody struct {
	CardNumber  int64  `json:"cardNumber" binding:"required"`
	SecretValue string `json:"secretValue"`
}

// lockBankAccountResponse - represents lockBankAccount response.
type lockBankAccountResponse struct {
	Status *string        `json:"status,omitempty"`
	Error  *service.Error `json:"error,omitempty"`
}

func (r *bankAccountRoutes) lockBankAccount(c *gin.Context) {
	logger := r.logger.Named("blockBankAccount")

	// parse request body
	logger.Debug("parsing request body")
	var body lockBankAccountRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// get client
	client, err := r.repos.Users.GetUserByID(c.Request.Context(), c.GetInt("clientID"))
	if err != nil {
		return
	}

	// lock bank account
	logger.Debug("bank account blocking")
	status, err := r.service.BlockBankAccount(c.Request.Context(), client.FullName, c.GetString("userRole"),
		&service.ChangeBankAccountInput{
			CardNumber:  body.CardNumber,
			SecretValue: body.SecretValue,
		})
	if err != nil {
		logger.Error("failed to block bank account", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, resetPasswordResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to block account")
		return
	}

	logger.Info("successfully block bank account")
	c.JSON(http.StatusOK, lockBankAccountResponse{Status: &status})

	var mess string
	if status == "ACTIVE" {
		mess = fmt.Sprintf("Successfully change status to %s for bank account %d", status, body.CardNumber)
	} else {
		mess = fmt.Sprintf("%s for bank account %d", status, body.CardNumber)
	}
	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: mess,
		})
	if err != nil {
		return
	}
}

// unlockBankAccountRequestBody - represents unlockBankAccount request body.
type unlockBankAccountRequestBody struct {
	CardNumber  int64  `json:"cardNumber" binding:"required"`
	SecretValue string `json:"secretValue"`
}

// unlockBankAccountResponse - represents unlockBankAccount response.
type unlockBankAccountResponse struct {
	Status *string        `json:"status,omitempty"`
	Error  *service.Error `json:"error,omitempty"`
}

func (r *bankAccountRoutes) unlockBankAccount(c *gin.Context) {
	logger := r.logger.Named("unlockBankAccount")

	// parse request body
	logger.Debug("parsing request body")
	var body unlockBankAccountRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Error("failed to parse body", "err", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	logger = logger.With("body", body)

	// get client
	client, err := r.repos.Users.GetUserByID(c.Request.Context(), c.GetInt("clientID"))
	if err != nil {
		return
	}

	// unlock bank account
	logger.Debug("bank account unlocking")
	status, err := r.service.UnlockBankAccount(c.Request.Context(), client.FullName, c.GetString("userRole"),
		&service.ChangeBankAccountInput{
			CardNumber:  body.CardNumber,
			SecretValue: body.SecretValue,
		})
	if err != nil {
		logger.Error("failed to unlock bank account", "err", err)
		err, ok := err.(*service.Error)
		if ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, resetPasswordResponse{Error: err})
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to unlock account")
		return
	}

	logger.Info("successfully unlock bank account")
	c.JSON(http.StatusOK, unlockBankAccountResponse{Status: &status})

	var mess string
	if status == "LOCK" {
		mess = fmt.Sprintf("Successfully change status to %s for bank account %d", status, body.CardNumber)
	} else {
		mess = fmt.Sprintf("%s for bank account %d", status, body.CardNumber)
	}
	_, err = r.service.MessageLogs.CreateMessageLog(c.Request.Context(), c.GetInt("clientID"),
		&service.MessageLogInput{
			MessageLog: mess,
		})
	if err != nil {
		return
	}
}

// getFilterFromQuery - returns filter from query.
func getFilterFromQuery(r *http.Request) (*domain.Filter, error) {
	filter, err := domain.GetFilterFromQuery(r)
	if err != nil {
		return nil, err
	}
	return filter, nil
}
