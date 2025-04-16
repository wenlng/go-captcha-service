package common

import (
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
)

type SvcContext struct {
	Cache   cache.Cache
	Config  *config.Config
	Logger  *zap.Logger
	Captcha *gocaptcha.GoCaptcha
}

func NewSvcContext() *SvcContext {
	return &SvcContext{}
}
