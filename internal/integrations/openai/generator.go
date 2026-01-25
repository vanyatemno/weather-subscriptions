package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"go.uber.org/zap"
)

const generationModel = "openai/gpt-5.2:online"

func (o *OpenAIService) GeneratePlainResponse(ctx context.Context, prompt string) (string, error) {
	res, err := o.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: generationModel,
		},
	)
	if err != nil {
		zap.L().Error("failed to generate plain response", zap.Error(err))
		return "", err
	}

	if len(res.Choices) == 0 {
		zap.L().Error("No choices found", zap.String("prompt", prompt))
		return "", errors.New("generation request returned no choices")
	}

	return res.Choices[0].Message.Content, nil
}

func (o *OpenAIService) GenerateStructuredResponse(ctx context.Context, prompt string, dest any) error {
	responseFormat, err := getResponseFormat(dest)
	if err != nil {
		zap.L().Error("Failed to get response format", zap.Error(err))
	}

	chatParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:          generationModel,
		ResponseFormat: *responseFormat,
	}

	completion, err := o.client.Chat.Completions.New(ctx, chatParams)
	if err != nil {
		zap.L().Error("Failed to create completion", zap.Error(err))
		return fmt.Errorf("openai completion error: %w", err)
	}

	choice := completion.Choices[0]
	if choice.Message.Refusal != "" {
		zap.L().Error("Failed to create completion", zap.Any("choice", choice))
		return fmt.Errorf("model refused to generate response: %s", choice.Message.Refusal)
	}

	content := choice.Message.Content
	if content == "" {
		zap.L().Error("Failed to create completion", zap.Any("choice", choice))
		return errors.New("received empty content from model")
	}

	if err = json.Unmarshal([]byte(content), dest); err != nil {
		zap.L().Error("Failed to create completion", zap.Any("choice", choice), zap.Error(err))
		return fmt.Errorf("failed to unmarshal structured response: %w", err)
	}

	return nil
}

func getResponseFormat(dest any) (*openai.ChatCompletionNewParamsResponseFormatUnion, error) {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr && !val.IsNil() && val.Elem().Kind() == reflect.Struct {
		return nil, errors.New("dest must be a non-nil pointer to a struct")
	}

	elemType := val.Type().Elem()
	schemaName := elemType.Name()
	if schemaName == "" {
		schemaName = "StructuredResponse"
	}

	params := &schemaParams{
		Type:        elemType,
		Name:        schemaName,
		Description: fmt.Sprintf("Generate a valid %s object based on the prompt.", schemaName),
	}

	schema := generateSchema(params.Type)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        params.Name,
		Description: openai.String(params.Description),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	return &openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
	}, nil
}

func generateSchema(t reflect.Type) interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	v := reflect.New(t).Elem().Interface()
	schema := reflector.Reflect(v)
	return schema
}
