package mock

import "context"

type GeneratorMock struct {
	generateStructuredFuncMock func(ctx context.Context, prompt string, dest any) error
}

func NewGeneratorMock(
	generateStructuredFuncMock func(ctx context.Context, prompt string, dest any) error,
) *GeneratorMock {
	return &GeneratorMock{
		generateStructuredFuncMock: generateStructuredFuncMock,
	}
}

func (g *GeneratorMock) SetStructuredFunc(
	generateStructuredFunc func(ctx context.Context, prompt string, dest any) error,
) {
	g.generateStructuredFuncMock = generateStructuredFunc
}

func (g *GeneratorMock) GeneratePlainResponse(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (g *GeneratorMock) GenerateStructuredResponse(ctx context.Context, prompt string, dest any) error {
	return g.generateStructuredFuncMock(ctx, prompt, dest)
}
