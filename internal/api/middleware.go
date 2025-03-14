// internal/api/middleware.go
package api

import (
	"context"
	"net/http"
	"time"

	"github.com/vgeshiktor/nba-stats/pkg/errors"
	"github.com/vgeshiktor/nba-stats/pkg/logger"

	"github.com/google/uuid"
)

// contextKey is a custom type for storing values in context.
type contextKey string

// RequestIDKey is the context key for the unique request ID.
const RequestIDKey contextKey = "requestID"

// LoggingMiddleware logs the details of each incoming request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Retrieve request ID from context (if available)
		reqID, _ := r.Context().Value(RequestIDKey).(string)
		logger.Info("[RequestID: %s] Started %s %s", reqID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		logger.Info("[RequestID: %s] Completed %s %s in %v", reqID, r.Method, r.URL.Path, time.Since(start))
	})
}

// AuthenticationMiddleware checks for the presence of an Authorization header.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("checking authentication...")
		// Simple authentication: verify that the Authorization header is set.
		if r.Header.Get("Authorization") == "" {
			errors.WriteError(w, http.StatusUnauthorized, "Unauthorized: missing token")
			return
		}
		logger.Info("Authentication successful")
		// In a real-world scenario, add token validation logic here.
		next.ServeHTTP(w, r)
	})
}

// RequestTracingMiddleware assigns a unique request ID to each request.
func RequestTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		// Optionally, add the request ID to response headers for downstream tracing.
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ChainMiddleware applies a list of middleware functions to an http.Handler.
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Apply middleware in reverse order so that the first middleware
	// in the list is the outermost wrapper.
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
