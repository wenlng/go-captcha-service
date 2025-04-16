package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"github.com/wenlng/go-captcha/v2/click"
	"go.uber.org/zap"
)

func TestCacheLogic(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	ttl := time.Duration(10) * time.Second
	cleanInt := time.Duration(30) * time.Second
	cacheClient := cache.NewMemoryCache("TEST_CAPTCHA_DATA:", ttl, cleanInt)
	defer cacheClient.Close()

	dc := &config.DynamicConfig{Config: config.DefaultConfig()}
	cnf := dc.Get()

	logger, err := zap.NewProduction()
	assert.NoError(t, err)

	captcha, err := gocaptcha.Setup()
	assert.NoError(t, err)

	svcCtx := &common.SvcContext{
		Cache:   cacheClient,
		Config:  &cnf,
		Logger:  logger,
		Captcha: captcha,
	}
	logic := NewClickCaptLogic(svcCtx)

	t.Run("GetData", func(t *testing.T) {
		_, err := logic.GetData(context.Background(), 0, 1, 1)
		assert.NoError(t, err)
	})

	t.Run("GetData_Miss", func(t *testing.T) {
		_, err := logic.GetData(context.Background(), -1, 1, 1)
		assert.Error(t, err)
	})

	t.Run("CheckData", func(t *testing.T) {
		data, err := logic.GetData(context.Background(), 1, 1, 1)
		assert.NoError(t, err)

		cacheData, err := svcCtx.Cache.GetCache(context.Background(), data.CaptchaKey)
		assert.NoError(t, err)

		var dct map[int]*click.Dot
		err = json.Unmarshal([]byte(cacheData), &dct)
		assert.NoError(t, err)

		var dots []string
		for i := 0; i < len(dct); i++ {
			dot := dct[i]
			dots = append(dots, strconv.Itoa(dot.X), strconv.Itoa(dot.Y))
		}

		dotStr := strings.Join(dots, ",")
		result, err := logic.CheckData(context.Background(), data.CaptchaKey, dotStr)
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("CheckData_MISS", func(t *testing.T) {
		data, err := logic.GetData(context.Background(), 1, 1, 1)
		assert.NoError(t, err)

		cacheData, err := svcCtx.Cache.GetCache(context.Background(), data.CaptchaKey)
		assert.NoError(t, err)

		var dct map[int]*click.Dot
		err = json.Unmarshal([]byte(cacheData), &dct)
		assert.NoError(t, err)

		var dots = []string{
			"111",
			"222",
		}
		dotStr := strings.Join(dots, ",")
		result, err := logic.CheckData(context.Background(), data.CaptchaKey, dotStr)
		assert.NoError(t, err)
		assert.Equal(t, false, result)
	})
}
