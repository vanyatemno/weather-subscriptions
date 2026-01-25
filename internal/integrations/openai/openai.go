package openai

import (
	"github.com/openai/openai-go"

	"github.com/openai/openai-go/option"
	"weather-subscriptions/internal/config"
)

const providerURL = "https://openrouter.ai/api/v1"

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(config *config.Config) *OpenAIService {
	client := openai.NewClient(
		option.WithBaseURL(providerURL),
		option.WithAPIKey(config.OpenAI.OpenrouterAPIKey),
	)
	return &OpenAIService{
		client: &client,
	}
}
