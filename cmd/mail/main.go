package main

import (
	"context"
	"fmt"
	"weather-subscriptions/internal/advises"
	"weather-subscriptions/internal/config"
	"weather-subscriptions/internal/db"
	"weather-subscriptions/internal/integrations/openai"
	"weather-subscriptions/internal/mail"
	"weather-subscriptions/internal/state"

	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	cfg, err := config.Read()
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	database, err := db.Connect(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	set := state.NewState(database)

	openaiService := openai.NewOpenAIService(cfg)
	advisesService := advises.NewAdvisesService(set, openaiService)
	mailerService := mail.NewMailerService(cfg)
	mailManager := mail.NewManager(
		context.Background(),
		cfg,
		set,
		set,
		set,
		set,
		mailerService,
		advisesService,
	)

	err = mailManager.SendDaily()
	if err != nil {
		zap.L().Error("failed to send daily mail", zap.Error(err))
	}
}
