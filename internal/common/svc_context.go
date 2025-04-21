/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package common

import (
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
)

// SvcContext service context
type SvcContext struct {
	CacheMgr      *cache.CacheManager
	DynamicConfig *config.DynamicConfig
	Logger        *zap.Logger
	Captcha       *gocaptcha.GoCaptcha
}

// NewSvcContext ..
func NewSvcContext() *SvcContext {
	return &SvcContext{}
}
