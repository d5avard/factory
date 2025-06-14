package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestStatusRoute(t *testing.T) {
	logger := NewLogger("logs", "test.log")

	// Avoid zap.Sync() in tests
	// defer logger.Close()

	router := mux.NewRouter()
	routes := []RouteInjector{mockRoute{}}

	server := NewServer(logger, router, routes)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	res := httptest.NewRecorder()

	server.Router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}

	if !strings.Contains(res.Body.String(), "ok") {
		t.Errorf("unexpected response body: %s", res.Body.String())
	}
}

type mockRoute struct{}

func (m mockRoute) Register(r *mux.Router, logger *Logger) {
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Test route hit")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")
}
