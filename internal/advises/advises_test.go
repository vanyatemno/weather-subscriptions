package advises

import (
	"context"
	"errors"
	"strings"
	"testing"

	"weather-subscriptions/internal/db/models"
)

type advisesStateMock struct {
	getAdvisesFunc  func(weatherID string) ([]*models.Advise, error)
	saveAdvisesFunc func(advises []*models.Advise) error
	savedAdvises    []*models.Advise
}

func (m *advisesStateMock) GetAdvises(weatherID string) ([]*models.Advise, error) {
	if m.getAdvisesFunc == nil {
		return nil, nil
	}
	return m.getAdvisesFunc(weatherID)
}

func (m *advisesStateMock) SaveAdvises(advises []*models.Advise) error {
	m.savedAdvises = advises
	if m.saveAdvisesFunc == nil {
		return nil
	}
	return m.saveAdvisesFunc(advises)
}

func (m *advisesStateMock) RemoveAdvise(_ *models.Advise) error {
	return nil
}

type generatorMock struct {
	generateStructuredFunc func(ctx context.Context, prompt string, dest any) error
	lastPrompt             string
}

func (m *generatorMock) GeneratePlainResponse(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (m *generatorMock) GenerateStructuredResponse(ctx context.Context, prompt string, dest any) error {
	m.lastPrompt = prompt
	if m.generateStructuredFunc == nil {
		return nil
	}
	return m.generateStructuredFunc(ctx, prompt, dest)
}

func TestGetAdviseFromCache(t *testing.T) {
	ctx := context.Background()
	weather := &models.Weather{ID: "weather", City: models.City{Name: "Warsaw"}}
	cached := []*models.Advise{{Name: "Spot", Description: "desc", Link: "link", WeatherID: weather.ID}}

	advisesState := &advisesStateMock{
		getAdvisesFunc: func(weatherID string) ([]*models.Advise, error) {
			return cached, nil
		},
	}
	generator := &generatorMock{}
	service := NewAdvisesService(advisesState, generator)

	res, err := service.GetAdvise(ctx, weather)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || len(res.Places) != 1 {
		t.Fatalf("expected 1 place, got %v", res)
	}
	if generator.lastPrompt != "" {
		t.Fatalf("expected generator not to be called")
	}
	if advisesState.savedAdvises != nil {
		t.Fatalf("expected advises not to be saved when cached")
	}
}

func TestGetAdviseErrors(t *testing.T) {
	ctx := context.Background()
	weather := &models.Weather{ID: "weather", City: models.City{Name: "Warsaw"}}

	tests := []struct {
		name         string
		advisesState *advisesStateMock
		generator    *generatorMock
		wantError    bool
	}{
		{
			name: "get advises error",
			advisesState: &advisesStateMock{
				getAdvisesFunc: func(weatherID string) ([]*models.Advise, error) {
					return nil, errors.New("get")
				},
			},
			generator: &generatorMock{},
			wantError: true,
		},
		{
			name: "generate error",
			advisesState: &advisesStateMock{
				getAdvisesFunc: func(weatherID string) ([]*models.Advise, error) {
					return nil, nil
				},
			},
			generator: &generatorMock{
				generateStructuredFunc: func(ctx context.Context, prompt string, dest any) error {
					return errors.New("generate")
				},
			},
			wantError: true,
		},
		{
			name: "save error",
			advisesState: &advisesStateMock{
				getAdvisesFunc: func(weatherID string) ([]*models.Advise, error) {
					return nil, nil
				},
				saveAdvisesFunc: func(advises []*models.Advise) error {
					return errors.New("save")
				},
			},
			generator: &generatorMock{
				generateStructuredFunc: func(ctx context.Context, prompt string, dest any) error {
					advise := dest.(*Advise)
					advise.Places = []Place{{Name: "Spot", Description: "desc", Link: "link"}}
					return nil
				},
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewAdvisesService(test.advisesState, test.generator)
			_, err := service.GetAdvise(ctx, weather)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAdviseGenerateSuccess(t *testing.T) {
	ctx := context.Background()
	weather := &models.Weather{
		ID:          "weather",
		City:        models.City{Name: "Warsaw"},
		Temperature: 10.5,
		Description: "Sunny",
	}

	advisesState := &advisesStateMock{
		getAdvisesFunc: func(weatherID string) ([]*models.Advise, error) {
			return nil, nil
		},
	}
	generator := &generatorMock{
		generateStructuredFunc: func(ctx context.Context, prompt string, dest any) error {
			advise := dest.(*Advise)
			advise.Places = []Place{
				{Name: "Spot", Description: "Nice ([link](http://example.com))", Link: "http://example.com"},
			}
			return nil
		},
	}
	service := NewAdvisesService(advisesState, generator)

	res, err := service.GetAdvise(ctx, weather)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || len(res.Places) != 1 {
		t.Fatalf("expected generated advise")
	}
	if strings.Contains(res.Places[0].Description, "http") {
		t.Fatalf("expected links to be removed from description")
	}
	if advisesState.savedAdvises == nil || len(advisesState.savedAdvises) != 1 {
		t.Fatalf("expected advises to be saved")
	}
	if generator.lastPrompt == "" {
		t.Fatalf("expected prompt to be generated")
	}
}

func TestBuildPrompt(t *testing.T) {
	weather := &models.Weather{City: models.City{Name: "Paris"}, Temperature: 5.5, Description: "Cloudy"}
	res := buildPrompt(weather)
	if !strings.Contains(res, "Paris") {
		t.Fatalf("expected prompt to include city name")
	}
	if !strings.Contains(res, "5.500000°C") {
		t.Fatalf("expected prompt to include temperature")
	}
	if !strings.Contains(res, "Cloudy") {
		t.Fatalf("expected prompt to include description")
	}
}

func TestCleanupResponse(t *testing.T) {
	advise := &Advise{Places: []Place{{Description: "A ([link](http://example.com)) place"}}}
	cleanupResponse(advise)
	if strings.Contains(advise.Places[0].Description, "http") {
		t.Fatalf("expected links to be removed")
	}
}
