package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
	"weather-subscriptions/api/routes"
	"weather-subscriptions/internal/mail"
	"weather-subscriptions/internal/mail/mailer_service"
	"weather-subscriptions/internal/state"

	"github.com/go-co-op/gocron"
	fiber "github.com/gofiber/fiber/v2"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db"
)

var (
	webApp *fiber.App
	appCtx context.Context
)

func main() {
	var cancel context.CancelFunc
	appCtx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	cfg, err := config.Read()
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	mailerService := mailer_service.New(cfg)

	database, err := db.Connect(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	set := state.NewState(database)

	scheduler := createScheduler(cfg, set, mailerService)

	scheduler.StartAsync()
	go createWebserver(cfg, set, mailerService)

	for range appCtx.Done() {
		_ = webApp.ShutdownWithContext(appCtx)
		return
	}
}

func createWebserver(cfg *config.Config, set state.Stateful, mailer mailer_service.MailerService) {
	webApp = fiber.New()

	webApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	routes.New(cfg, set, mailer).Setup(webApp)
	if err := webApp.Listen(":" + cfg.Port); err != nil {
		zap.L().Error("failed to start server: %v", zap.Error(err))
	}
}

func createScheduler(cfg *config.Config, state state.Stateful, mailer mailer_service.MailerService) *gocron.Scheduler {
	mailManager := mail.New(appCtx, cfg, state, mailer)
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(1).Hour().Do(mailManager.SendHourly)
	if err != nil {
		zap.L().Error("failed to start cron: %v", zap.Error(err))
	}

	_, err = scheduler.Every(1).Day().At("12:00").Do(mailManager.SendDaily)
	if err != nil {
		zap.L().Error("failed to start cron: %v", zap.Error(err))
	}

	return scheduler
}
