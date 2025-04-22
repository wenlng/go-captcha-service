package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/wenlng/go-captcha-service/internal/config"
)

func TestAPIKeyMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	dc := &config.DynamicConfig{
		Config: config.Config{
			APIKeys: []string{"valid-key"},
		},
	}

	mw := APIKeyMiddleware(dc, logger)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	t.Run("ValidKey", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", "valid-key")
		rr := httptest.NewRecorder()
		mw(handler)(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("MissingKey", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()
		mw(handler)(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var resp ErrorResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Equal(t, "missing API Key", resp.Message)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", "invalid-key")
		rr := httptest.NewRecorder()
		mw(handler)(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var resp ErrorResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Equal(t, "invalid API Key", resp.Message)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	limiter := NewDynamicLimiter(1, 1) // 1 QPS, 1 burst
	mw := RateLimitMiddleware(limiter, logger)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	t.Run("WithinLimit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()
		mw(handler)(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("ExceedLimit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()
		mw(handler)(rr, req) // First request
		rr = httptest.NewRecorder()
		mw(handler)(rr, req) // Second request exceeds limit
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		var resp ErrorResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Equal(t, "rate limit exceeded", resp.Message)
	})
}

func TestLoggingMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mw := LoggingMiddleware(logger)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	mw(handler)(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCircuitBreakerMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "test",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 1
		},
	})
	mw := CircuitBreakerMiddleware(breaker, logger)
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusInternalServerError)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	mw(handler)(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Simulate multiple failures to trip breaker
	for i := 0; i < 2; i++ {
		rr = httptest.NewRecorder()
		mw(handler)(rr, req)
	}

	// Breaker should be open
	rr = httptest.NewRecorder()
	mw(handler)(rr, req)
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	var resp ErrorResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, "service unavailable", resp.Message)
}

func TestRateLimitHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	limiter := NewDynamicLimiter(1000, 1000)
	handler := RateLimitHandler(limiter, logger)

	t.Run("ValidUpdate", func(t *testing.T) {
		body := strings.NewReader(`{"qps":10,"burst":10}`)
		req := httptest.NewRequest("POST", "/rate-limit", body)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]string
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Equal(t, "success", resp["status"])
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rate-limit", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		var resp ErrorResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Equal(t, "method not allowed", resp.Message)
	})

	t.Run("InvalidParams", func(t *testing.T) {
		body := strings.NewReader(`{"qps":0,"burst":0}`)
		req := httptest.NewRequest("POST", "/rate-limit", body)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorResponse
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Equal(t, "qps and burst must be positive", resp.Message)
	})
}

func TestGRPCInterceptor(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	dc := &config.DynamicConfig{
		Config: config.Config{
			APIKeys: []string{"valid-key"},
		},
	}
	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "grpc-test",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     5 * time.Second,
	})
	interceptor := UnaryServerInterceptor(dc, logger, breaker)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	t.Run("ValidKey", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-api-key", "valid-key"))
		resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "TestMethod"}, handler)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("MissingKey", func(t *testing.T) {
		ctx := context.Background()
		_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "TestMethod"}, handler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("InvalidKey", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-api-key", "invalid-key"))
		_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "TestMethod"}, handler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})
}
