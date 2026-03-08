package webui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/gin-gonic/gin"
)

func TestFallbackServedWhenNextUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &config.Config{Frontend: config.FrontendConfig{NextInternalURL: ""}}
	RegisterRoutes(r, cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Frontend Next.js indisponível") {
		t.Fatalf("expected fallback HTML response")
	}
}
