/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/wenlng/go-captcha-service/internal/config"
)

// ErrorResponse defines the standard error response format
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HandlerFunc .
type HandlerFunc func(http.ResponseWriter, *http.Request)

// HTTPMiddleware defines the middleware function signature
type HTTPMiddleware func(HandlerFunc) HandlerFunc

// MwChain .
type MwChain struct {
	middlewares []HTTPMiddleware
}

// NewChainHTTP .
func NewChainHTTP(mws ...HTTPMiddleware) *MwChain {
	return &MwChain{middlewares: mws}
}

// AppendMiddleware new middlewares
func (c *MwChain) AppendMiddleware(final HTTPMiddleware) *MwChain {
	c.middlewares = append(c.middlewares, final)

	return c
}

// Then chains multiple middlewares
func (c *MwChain) Then(final HandlerFunc) http.HandlerFunc {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		if c.middlewares[i] != nil {
			final = c.middlewares[i](final)
		}
	}

	return http.HandlerFunc(final)
}

// APIKeyMiddleware validates API keys
func APIKeyMiddleware(dc *config.DynamicConfig, logger *zap.Logger) HTTPMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cfg := dc.Get()

			// Auth API
			authApisMap := cfg.GetAuthAPIs()
			if _, exists := authApisMap[r.URL.Path]; !exists {
				next(w, r)
				return
			}

			apiKeyMap := cfg.GetAPIKeys()
			if len(apiKeyMap) == 0 {
				logger.Warn("[HttpMiddleware] Missing API Key")
				WriteError(w, http.StatusUnauthorized, "missing API Key")
				return
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				logger.Warn("[HttpMiddleware] Missing API Key")
				WriteError(w, http.StatusUnauthorized, "missing API Key")
				return
			}
			if _, exists := apiKeyMap[apiKey]; !exists {
				logger.Warn("[HttpMiddleware] Invalid API Key", zap.String("key", apiKey))
				WriteError(w, http.StatusUnauthorized, "invalid API Key")
				return
			}
			next(w, r)
		}
	}
}

// RateLimitMiddleware enforces rate limiting
func RateLimitMiddleware(limiter *DynamicLimiter, logger *zap.Logger) HTTPMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := limiter.Wait(r.Context()); err != nil {
				logger.Warn("[HttpMiddleware] Rate limit exceeded", zap.String("client", r.RemoteAddr))
				WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next(w, r)
		}
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *zap.Logger) HTTPMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			logger.Info("[HttpMiddleware] HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("client", r.RemoteAddr),
				zap.Duration("duration", time.Since(start)),
			)
		}
	}
}

// CircuitBreakerMiddleware implements circuit breaking
func CircuitBreakerMiddleware(breaker *gobreaker.CircuitBreaker, logger *zap.Logger) HTTPMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := breaker.Execute(func() (interface{}, error) {
				next(w, r)
				return nil, nil
			})
			if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
				logger.Warn("[HttpMiddleware] Circuit breaker tripped", zap.Error(err))
				WriteError(w, http.StatusServiceUnavailable, "service unavailable")
				return
			}
			if err != nil {
				logger.Error("[HttpMiddleware] Circuit breaker error", zap.Error(err))
				WriteError(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}
	}
}

// CORSMiddleware implements cross-origin resource sharing
func CORSMiddleware(dc *config.DynamicConfig, logger *zap.Logger) HTTPMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !dc.Get().EnableCors {
				next(w, r)
				return
			}

			origin := r.Header.Get("Origin")

			allowedOrigins := []string{"*"}
			allowOrigin := "*"
			for _, allowed := range allowedOrigins {
				if allowed == "*" || allowed == origin {
					allowOrigin = origin
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			if requestedHeaders := r.Header.Get("Access-Control-Request-Headers"); requestedHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Custom-Header")
			}

			w.Header().Set("Access-Control-Max-Age", "86400")          // CacheMgr preflight response for 24 hours
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials if needed

			// Handle preflight (OPTIONS) requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Processing precheck requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}
}

// DynamicLimiter manages dynamic rate limiting
type DynamicLimiter struct {
	limiter *rate.Limiter
	mu      sync.RWMutex
}

// NewDynamicLimiter creates a new dynamic rate limiter
func NewDynamicLimiter(qps, burst int) *DynamicLimiter {
	return &DynamicLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
	}
}

// Wait checks the rate limit
func (d *DynamicLimiter) Wait(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.limiter.Wait(ctx)
}

// Update updates rate limit parameters
func (d *DynamicLimiter) Update(qps, burst int) {
	if qps <= 0 || burst <= 0 {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.limiter.SetLimit(rate.Limit(qps))
	d.limiter.SetBurst(burst)
}

// RateLimitHandler handles dynamic rate limit updates
func RateLimitHandler(limiter *DynamicLimiter, logger *zap.Logger) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var params struct {
			QPS   int `json:"qps"`
			Burst int `json:"burst"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			logger.Warn("[HttpMiddleware] Invalid rate limit params", zap.Error(err))
			WriteError(w, http.StatusBadRequest, "invalid parameters")
			return
		}
		if params.QPS <= 0 || params.Burst <= 0 {
			WriteError(w, http.StatusBadRequest, "qps and burst must be positive")
			return
		}
		limiter.Update(params.QPS, params.Burst)
		logger.Info("[HttpMiddleware] Rate limit updated",
			zap.Int("qps", params.QPS),
			zap.Int("burst", params.Burst),
		)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// WriteError sends an error response
func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}
