package server

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/proto"
)

func TestCacheServer(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	ttl := time.Duration(10) * time.Second
	cleanInt := time.Duration(30) * time.Second
	cacheClient := cache.NewMemoryCache("TEST_CAPTCHA_DATA:", ttl, cleanInt)
	defer cacheClient.Close()

	dc := &config.DynamicConfig{Config: config.DefaultConfig()}
	cnf := dc.Get()
	captDCfg := &config2.DynamicCaptchaConfig{Config: config2.DefaultConfig()}

	logger, err := zap.NewProduction()
	assert.NoError(t, err)

	captcha, err := gocaptcha.Setup(captDCfg)
	assert.NoError(t, err)

	svcCtx := &common.SvcContext{
		Cache:   cacheClient,
		Config:  &cnf,
		Logger:  logger,
		Captcha: captcha,
	}
	server := NewGoCaptchaServer(svcCtx)

	t.Run("GetData", func(t *testing.T) {
		req := &proto.GetDataRequest{
			Type: proto.GoCaptchaType_GoCaptchaTypeClick,
		}
		resp, err := server.GetData(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("GetData_Miss", func(t *testing.T) {
		req := &proto.GetDataRequest{
			Type: -1,
		}
		_, err := server.GetData(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		//assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
