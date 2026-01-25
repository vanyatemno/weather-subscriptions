package openai

import (
	"weather-subscriptions/internal/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const openrouterProviderURL = "https://openrouter.ai/api/v1"

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(config *config.Config) *OpenAIService {
	client := openai.NewClient(
		option.WithBaseURL(openrouterProviderURL),
		option.WithAPIKey(config.OpenAI.OpenrouterAPIKey),
	)
	return &OpenAIService{
		client: &client,
	}
}
