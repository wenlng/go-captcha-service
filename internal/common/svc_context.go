package common

import (
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
)

type SvcContext struct {
	Cache         cache.Cache
	DynamicConfig *config.DynamicConfig
	Logger        *zap.Logger
	Captcha       *gocaptcha.GoCaptcha
}

func NewSvcContext() *SvcContext {
	return &SvcContext{}
}
