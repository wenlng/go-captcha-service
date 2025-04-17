package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
)

// CommonLogic .
type CommonLogic struct {
	svcCtx *common.SvcContext

	cache      cache.Cache
	dynamicCfg *config.DynamicConfig
	logger     *zap.Logger
	captcha    *gocaptcha.GoCaptcha
}

// NewCommonLogic .
func NewCommonLogic(svcCtx *common.SvcContext) *CommonLogic {
	return &CommonLogic{
		svcCtx:     svcCtx,
		cache:      svcCtx.Cache,
		dynamicCfg: svcCtx.DynamicConfig,
		logger:     svcCtx.Logger,
		captcha:    svcCtx.Captcha,
	}
}

// CheckStatus .
func (cl *CommonLogic) CheckStatus(ctx context.Context, key string) (ret bool, err error) {
	if key == "" {
		return false, fmt.Errorf("invalid key")
	}

	cacheData, err := cl.cache.GetCache(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to get cache: %v", err)
	}

	if cacheData == "" {
		return false, nil
	}

	var captData *cache.CaptCacheData
	err = json.Unmarshal([]byte(cacheData), &captData)
	if err != nil {
		return false, fmt.Errorf("failed to json unmarshal: %v", err)
	}

	return captData.Status == 1, nil
}

// GetStatusInfo .
func (cl *CommonLogic) GetStatusInfo(ctx context.Context, key string) (data *cache.CaptCacheData, err error) {
	if key == "" {
		return nil, fmt.Errorf("invalid key")
	}

	captData := &cache.CaptCacheData{}

	cacheData, err := cl.cache.GetCache(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %v", err)
	}

	if cacheData == "" {
		captData.Data = struct{}{}
		return captData, nil
	}

	err = json.Unmarshal([]byte(cacheData), &captData)
	if err != nil {
		return nil, fmt.Errorf("failed to json unmarshal: %v", err)
	}

	return captData, nil
}

// DelStatusInfo .
func (cl *CommonLogic) DelStatusInfo(ctx context.Context, key string) (ret bool, err error) {
	if key == "" {
		return false, fmt.Errorf("invalid key")
	}

	err = cl.cache.DeleteCache(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to delete cache: %v", err)
	}

	return true, nil
}
