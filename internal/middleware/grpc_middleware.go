/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package middleware

import (
	"context"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/wenlng/go-captcha-service/internal/config"
)

// UnaryServerInterceptor implements gRPC unary interceptor
func UnaryServerInterceptor(dc *config.DynamicConfig, logger *zap.Logger, breaker *gobreaker.CircuitBreaker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Validate API key
		cfg := dc.Get()
		apiKeyMap := make(map[string]struct{})
		for _, key := range cfg.APIKeys {
			apiKeyMap[key] = struct{}{}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Warn("[GrpcMiddleware] Missing metadata")
			return nil, status.Error(codes.Unauthenticated, "missing API Key")
		}
		apiKeys := md.Get("x-api-key")
		if len(apiKeys) == 0 {
			logger.Warn("[GrpcMiddleware] Missing API Key")
			return nil, status.Error(codes.Unauthenticated, "missing API Key")
		}
		if _, exists := apiKeyMap[apiKeys[0]]; !exists {
			logger.Warn("[GrpcMiddleware] Invalid API Key", zap.String("key", apiKeys[0]))
			return nil, status.Error(codes.Unauthenticated, "invalid API Key")
		}

		// Apply circuit breaker
		var resp interface{}
		var err error
		_, cbErr := breaker.Execute(func() (interface{}, error) {
			resp, err = handler(ctx, req)
			return nil, nil
		})
		if cbErr == gobreaker.ErrOpenState || cbErr == gobreaker.ErrTooManyRequests {
			logger.Warn("[GrpcMiddleware] gRPC circuit breaker tripped", zap.Error(cbErr))
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}
		if cbErr != nil {
			logger.Error("[GrpcMiddleware] gRPC circuit breaker error", zap.Error(cbErr))
			return nil, status.Error(codes.Internal, "internal server error")
		}

		// Log request
		logger.Info("[GrpcMiddleware] gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}
