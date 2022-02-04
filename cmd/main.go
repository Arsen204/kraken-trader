package main

import (
	"course-project/configs"
	"course-project/internal/delivery/rest"
	"course-project/internal/domain"
	"course-project/internal/repository"
	"course-project/internal/server"
	"course-project/internal/service"
	pkgpostgres "course-project/pkg/postgres"
	tg "course-project/pkg/telegram"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 30 * time.Second

func main() {
	// Logger init
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	// logger.Formatter = &logrus.JSONFormatter{}

	// Templates init
	err := domain.InitTemplate()
	if err != nil {
		logger.Error(err)
	}

	// Configs init
	config, err := configs.NewAppConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Repository
	pool, err := pkgpostgres.NewPool(config.DB.DNS(), logger)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)

	// Dependencies
	telegramService := tg.NewTelegramService(config.Telegram)

	publishService := service.NewPublishService(logger, repo, telegramService)
	exchangeService := service.NewExchangeService(logger, config.Kraken, publishService)
	indicatorService := service.NewIndicatorService()
	algorithmService := service.NewAlgorithmService(logger, exchangeService, indicatorService)
	services := service.NewService(algorithmService, exchangeService)

	// HTTP Server
	srv := new(server.Server)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	handlers := rest.NewHandler(serverCtx, serverStopCtx, logger, services)

	logger.Print("app running")

	go func() {
		err := srv.Run(config.Port, handlers.InitRoutes())
		if err != http.ErrServerClosed {
			logger.Fatalf("error occurred while running rest server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quit
		serverStopCtx()
	}()

	<-serverCtx.Done()

	logger.Print("app shutting down")

	shutdownCtx, shutdownStopCtx := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownStopCtx()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}
