package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/middleware"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"

	"github.com/wenlng/go-captcha-service/internal/cache"
)

func TestHTTPHandlers(t *testing.T) {
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

	svcCtx := &base.SvcContext{
		Cache:   cacheClient,
		Config:  &cnf,
		Logger:  logger,
		Captcha: captcha,
	}
	handlers := NewHTTPHandlers(svcCtx)

	t.Run("GetDataHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get-data?key=key1", nil)
		rr := httptest.NewRecorder()
		handlers.GetDataHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &resp)
		fmt.Println(resp)
		assert.Equal(t, "value1", resp["value"])
	})

	t.Run("ReadHandler_Miss", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/read?key=nonexistent", nil)
		rr := httptest.NewRecorder()
		handlers.GetDataHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		var resp middleware.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Contains(t, resp.Message, "cache miss")
	})

	t.Run("CheckDataHandler", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{"key": "key2", "value": "value2"})
		req := httptest.NewRequest(http.MethodPost, "/write", bytes.NewReader(reqBody))
		rr := httptest.NewRecorder()
		handlers.CheckDataHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "success", resp["status"])

		value, err := cacheClient.GetCache(context.Background(), "key2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)
	})

	t.Run("WriteHandler_InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/write", bytes.NewReader([]byte("invalid")))
		rr := httptest.NewRecorder()
		handlers.CheckDataHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var resp middleware.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "invalid request body", resp.Message)
	})
}
