package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wenlng/go-captcha-service/internal/adapt"
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"github.com/wenlng/go-captcha/v2/rotate"
	"go.uber.org/zap"
)

// RotateCaptLogic .
type RotateCaptLogic struct {
	svcCtx *common.SvcContext

	cache   cache.Cache
	config  *config.Config
	logger  *zap.Logger
	captcha *gocaptcha.GoCaptcha
}

// NewRotateCaptLogic .
func NewRotateCaptLogic(svcCtx *common.SvcContext) *RotateCaptLogic {
	return &RotateCaptLogic{
		svcCtx:  svcCtx,
		cache:   svcCtx.Cache,
		config:  svcCtx.Config,
		logger:  svcCtx.Logger,
		captcha: svcCtx.Captcha,
	}
}

// GetData .
func (cl *RotateCaptLogic) GetData(ctx context.Context, ctype, theme, lang int) (res *adapt.CaptData, err error) {
	res = &adapt.CaptData{}

	if ctype < 0 {
		return nil, fmt.Errorf("missing parameter")
	}

	captData, err := cl.captcha.RotateCaptInstance.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate captcha data failed: %v", err)
	}

	data := captData.GetData()
	if data == nil {
		return nil, fmt.Errorf("generate captcha data failed: %v", err)
	}

	res.MasterImageBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, fmt.Errorf("failed to convert base64 encoding: %v", err)
	}

	res.ThumbImageBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		return nil, fmt.Errorf("failed to convert base64 encoding: %v", err)
	}

	cacheData := &cache.CaptCacheData{
		Data:   data,
		Status: 0,
	}
	cacheDataByte, err := json.Marshal(cacheData)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal: %v", err)
	}

	key, err := helper.GenUniqueId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %v", err)
	}

	err = cl.cache.SetCache(ctx, key, string(cacheDataByte))
	if err != nil {
		return res, fmt.Errorf("failed to write cache:: %v", err)
	}

	opts := cl.captcha.ClickCaptInstance.GetOptions()
	res.MasterImageWidth = int32(opts.GetImageSize().Width)
	res.MasterImageHeight = int32(opts.GetImageSize().Height)
	res.ThumbImageWidth = int32(data.Width)
	res.ThumbImageHeight = int32(data.Height)
	res.ThumbImageSize = int32(data.Width)
	res.CaptchaKey = key
	return res, nil
}

// CheckData .
func (cl *RotateCaptLogic) CheckData(ctx context.Context, key string, angle int) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("invalid key")
	}

	cacheData, err := cl.cache.GetCache(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to get cache: %v", err)
	}

	var captData *cache.CaptCacheData
	err = json.Unmarshal([]byte(cacheData), &captData)
	if err != nil {
		return false, fmt.Errorf("failed to json unmarshal: %v", err)
	}

	dct, ok := captData.Data.(*rotate.Block)
	if !ok {
		return false, fmt.Errorf("cache data invalid: %v", err)
	}

	ret := rotate.CheckAngle(int64(angle), int64(dct.Angle), 2)

	if ret {
		captData.Status = 1
		cacheDataByte, err := json.Marshal(captData)
		if err != nil {
			return ret, fmt.Errorf("failed to json marshal: %v", err)
		}

		err = cl.cache.SetCache(ctx, key, string(cacheDataByte))
		if err != nil {
			return ret, fmt.Errorf("failed to update cache:: %v", err)
		}
	}

	return ret, nil
}
