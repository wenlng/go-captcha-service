package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/wenlng/go-captcha-service/internal/adapt"
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"github.com/wenlng/go-captcha/v2/slide"
	"go.uber.org/zap"
)

// SlideCaptLogic .
type SlideCaptLogic struct {
	svcCtx *common.SvcContext

	cache   cache.Cache
	config  *config.Config
	logger  *zap.Logger
	captcha *gocaptcha.GoCaptcha
}

// NewSlideCaptLogic .
func NewSlideCaptLogic(svcCtx *common.SvcContext) *SlideCaptLogic {
	return &SlideCaptLogic{
		svcCtx:  svcCtx,
		cache:   svcCtx.Cache,
		config:  svcCtx.Config,
		logger:  svcCtx.Logger,
		captcha: svcCtx.Captcha,
	}
}

// GetData .
func (cl *SlideCaptLogic) GetData(ctx context.Context, ctype, theme, lang int) (res *adapt.CaptData, err error) {
	res = &adapt.CaptData{}

	if ctype < 0 {
		return nil, fmt.Errorf("missing parameter")
	}

	captData, err := cl.captcha.SlideCaptInstance.Generate()
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

	res.ThumbImageBase64, err = captData.GetTileImage().ToBase64()
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
	res.DisplayX = int32(data.TileX)
	res.DisplayY = int32(data.TileY)
	res.CaptchaKey = key
	return res, nil
}

// CheckData .
func (cl *SlideCaptLogic) CheckData(ctx context.Context, key string, dots string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("invalid key")
	}

	cacheData, err := cl.cache.GetCache(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to get cache: %v", err)
	}

	src := strings.Split(dots, ",")

	var captData *cache.CaptCacheData
	err = json.Unmarshal([]byte(cacheData), &captData)
	if err != nil {
		return false, fmt.Errorf("failed to json unmarshal: %v", err)
	}

	dct, ok := captData.Data.(*slide.Block)
	if !ok {
		return false, fmt.Errorf("cache data invalid: %v", err)
	}

	ret := false
	if 2 == len(src) {
		sx, _ := strconv.ParseInt(src[0], 10, 64)
		sy, _ := strconv.ParseInt(src[1], 10, 64)
		ret = slide.CheckPoint(sx, sy, int64(dct.X), int64(dct.Y), 4)
	}

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
