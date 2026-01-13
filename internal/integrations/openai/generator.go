package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
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
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("generation request returned no choices")
	}

	return res.Choices[0].Message.Content, nil
}

func (o *OpenAIService) GenerateStructuredResponse(ctx context.Context, prompt string, dest any) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("dest must be a non-nil pointer")
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

	responseFormat := getResponseFormat(params)

	chatParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:          generationModel,
		ResponseFormat: responseFormat,
	}

	// 4. Execute the request
	completion, err := o.client.Chat.Completions.New(ctx, chatParams)
	if err != nil {
		return fmt.Errorf("openai completion error: %w", err)
	}

	// 5. Check for refusal or empty content
	choice := completion.Choices[0]
	if choice.Message.Refusal != "" {
		return fmt.Errorf("model refused to generate response: %s", choice.Message.Refusal)
	}

	content := choice.Message.Content
	if content == "" {
		return fmt.Errorf("received empty content from model")
	}

	// 6. Unmarshal result into dest
	// The model is guaranteed to return valid JSON matching the schema because Strict: true is set
	if err := json.Unmarshal([]byte(content), dest); err != nil {
		return fmt.Errorf("failed to unmarshal structured response: %w", err)
	}

	return nil
}

func getResponseFormat(params *schemaParams) openai.ChatCompletionNewParamsResponseFormatUnion {
	schema := generateSchema(params.Type)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        params.Name,
		Description: openai.String(params.Description),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	return openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
	}
}

func generateSchema(t reflect.Type) interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	v := reflect.New(t).Elem().Interface()
	schema := reflector.Reflect(v)
	return schema
}
