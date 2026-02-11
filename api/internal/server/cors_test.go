package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/stretchr/testify/assert"
)

// newCORSRouter creates a chi router with the CORS middleware configured
// identically to how setupRoutes will configure it. This allows testing
// CORS behavior without requiring a full Server (DB, Redis, etc.).
func newCORSRouter(allowedOrigins string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(allowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return r
}

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		allowedOrigins  string
		requestOrigin   string
		method          string
		wantOrigin      string
		wantCredentials string
		wantMethods     string
		wantHeaders     string
		description     string
	}{
		{
			name:            "allowed origin GET returns Access-Control-Allow-Origin",
			allowedOrigins:  "http://localhost:5173",
			requestOrigin:   "http://localhost:5173",
			method:          "GET",
			wantOrigin:      "http://localhost:5173",
			wantCredentials: "true",
			wantMethods:     "",
			wantHeaders:     "",
			description:     "A GET from an allowed origin should echo back that origin in Access-Control-Allow-Origin. Check that the allowedOrigins config includes this origin.",
		},
		{
			name:            "disallowed origin GET does not return Access-Control-Allow-Origin",
			allowedOrigins:  "http://localhost:5173",
			requestOrigin:   "http://evil.example.com",
			method:          "GET",
			wantOrigin:      "",
			wantCredentials: "",
			wantMethods:     "",
			wantHeaders:     "",
			description:     "A GET from a disallowed origin should NOT have Access-Control-Allow-Origin header. Check that the cors middleware rejects unknown origins.",
		},
		{
			name:            "OPTIONS preflight returns correct CORS headers",
			allowedOrigins:  "http://localhost:5173",
			requestOrigin:   "http://localhost:5173",
			method:          "OPTIONS",
			wantOrigin:      "http://localhost:5173",
			wantCredentials: "true",
			wantMethods:     "GET",
			wantHeaders:     "",
			description:     "An OPTIONS preflight from an allowed origin should return Allow-Origin and Allow-Methods. Verify cors.Options.AllowedMethods includes the expected methods.",
		},
		{
			name:            "credentials header is set for allowed origin",
			allowedOrigins:  "http://localhost:5173",
			requestOrigin:   "http://localhost:5173",
			method:          "GET",
			wantOrigin:      "http://localhost:5173",
			wantCredentials: "true",
			wantMethods:     "",
			wantHeaders:     "",
			description:     "Access-Control-Allow-Credentials should be 'true' for allowed origins. Check that cors.Options.AllowCredentials is set to true.",
		},
		{
			name:            "multiple origins - first origin is allowed",
			allowedOrigins:  "http://localhost:5173,http://localhost:3000",
			requestOrigin:   "http://localhost:5173",
			method:          "GET",
			wantOrigin:      "http://localhost:5173",
			wantCredentials: "true",
			wantMethods:     "",
			wantHeaders:     "",
			description:     "When multiple origins are configured (comma-separated), the first origin should be allowed. Check strings.Split parsing of AllowedOrigins.",
		},
		{
			name:            "multiple origins - second origin is allowed",
			allowedOrigins:  "http://localhost:5173,http://localhost:3000",
			requestOrigin:   "http://localhost:3000",
			method:          "GET",
			wantOrigin:      "http://localhost:3000",
			wantCredentials: "true",
			wantMethods:     "",
			wantHeaders:     "",
			description:     "When multiple origins are configured (comma-separated), the second origin should also be allowed. Check strings.Split parsing of AllowedOrigins.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newCORSRouter(tt.allowedOrigins)

			var req *http.Request
			if tt.method == "OPTIONS" {
				req = httptest.NewRequest("OPTIONS", "/test", nil)
				req.Header.Set("Access-Control-Request-Method", "GET")
			} else {
				req = httptest.NewRequest(tt.method, "/test", nil)
			}
			req.Header.Set("Origin", tt.requestOrigin)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			gotOrigin := resp.Header.Get("Access-Control-Allow-Origin")
			assert.Equal(t, tt.wantOrigin, gotOrigin,
				"Access-Control-Allow-Origin mismatch. %s", tt.description)

			gotCredentials := resp.Header.Get("Access-Control-Allow-Credentials")
			assert.Equal(t, tt.wantCredentials, gotCredentials,
				"Access-Control-Allow-Credentials mismatch. %s", tt.description)

			if tt.wantMethods != "" {
				gotMethods := resp.Header.Get("Access-Control-Allow-Methods")
				assert.Contains(t, gotMethods, tt.wantMethods,
					"Access-Control-Allow-Methods should contain '%s'. %s", tt.wantMethods, tt.description)
			}
		})
	}
}
