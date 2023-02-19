// Package app represents application.
package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	// third party
	"github.com/gin-gonic/gin"

	// external

	"github.com/Shevchenkko/payment_system/pkg/httpserver"
	"github.com/Shevchenkko/payment_system/pkg/logger"
	"github.com/Shevchenkko/payment_system/pkg/mysql"

	// internal
	"github.com/Shevchenkko/payment_system/internal/api/emails"
	"github.com/Shevchenkko/payment_system/internal/controller"
	"github.com/Shevchenkko/payment_system/internal/domain"
	"github.com/Shevchenkko/payment_system/internal/repository"
	"github.com/Shevchenkko/payment_system/internal/service"
)

// Run - initializes and runs application.
func Run() {

	// init logger
	l := logger.New(os.Getenv("LOG_LEVEL"))

	// init repository
	sql, err := mysql.New(mysql.MySQLConfig{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Database: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		l.Fatal("failed to connect to mysql", "err", err)
	}

	err = sql.DB.AutoMigrate(
		&domain.User{},
		&domain.BankAccount{},
	)

	if err != nil {
		l.Fatal("automigration failed", "err", err)
	}

	// init apis
	apis := service.APIs{
		Emails: emails.New(),
	}

	// init repositories
	repositories := service.Repositories{
		Users: repository.NewUsersRepo(sql),
		Banks: repository.NewBankAccountsRepo(sql),
	}

	services := service.Services{
		Users: service.NewUserService(
			repositories,
			apis,
		),
		BankAccounts: service.NewBankAccountService(
			repositories,
		),
	}

	// init framework of choice
	handler := gin.New()

	// init router
	controller.NewRouter(handler, services, l)

	// init and run http server
	httpServer := httpserver.New(handler, httpserver.Port(os.Getenv("HTTP_PORT")), httpserver.ReadTimeout(60*time.Second), httpserver.WriteTimeout(60*time.Second))

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify", "err", err)
	}

	// shutdown http server
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
