package e2e

import (
	"testing"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func TestSubscriptions(t *testing.T) {
	webApp, state := newTestApp()
	app := adaptor.FiberApp(webApp)
}
