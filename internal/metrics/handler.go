package metrics

import (
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns a Fiber handler for Prometheus metrics
func Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Set("Content-Type", "text/plain")
		
		// Create a test request to capture the metrics
		req := httptest.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		
		promhttp.Handler().ServeHTTP(w, req)
		
		return c.Send(w.Body.Bytes())
	}
}
