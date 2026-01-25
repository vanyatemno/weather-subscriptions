package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"weather-subscriptions/api/routes"
	_ "weather-subscriptions/docs"
	"weather-subscriptions/internal/advises"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db"
	"weather-subscriptions/internal/integrations/openai"
	"weather-subscriptions/internal/mail"
	"weather-subscriptions/internal/state"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

var (
	webApp *fiber.App
)

// @title Weather Subscriptions API
// @version 1.0
// @description API for managing weather subscriptions and fetching weather data.
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	appCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	cfg, err := config.Read()
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	mailerService := mail.NewMailerService(cfg)

	database, err := db.Connect(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	set := state.NewState(database)

	openaiService := openai.NewOpenAIService(cfg)
	advisesService := advises.NewAdvisesService(set, openaiService)

	scheduler := createScheduler(appCtx, cfg, set, mailerService, advisesService)

	scheduler.StartAsync()
	go createWebserver(cfg, set, mailerService)

	for range appCtx.Done() {
		_ = webApp.ShutdownWithContext(appCtx)
		return
	}
}

func createWebserver(cfg *config.Config, set *state.State, mailer mail.MailerService) {
	webApp = fiber.New()

	webApp.Use(
		cors.New(
			cors.Config{
				AllowOrigins: "*",
			},
		),
	)

	webApp.Get("/swagger/*", swagger.HandlerDefault)

	routes.New(cfg, set, mailer).Setup(webApp)
	if err := webApp.Listen(":" + cfg.Port); err != nil {
		zap.L().Error("failed to start server: %v", zap.Error(err))
	}
}

func createScheduler(
	ctx context.Context,
	cfg *config.Config,
	state *state.State,
	mailer mail.MailerService,
	advises *advises.AdvisesService,
) *gocron.Scheduler {
	mailManager := mail.New(
		ctx,
		cfg,
		state,
		state,
		state,
		state,
		mailer,
		advises,
	)
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
