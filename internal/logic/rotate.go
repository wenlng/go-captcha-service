/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wenlng/go-captcha-service/internal/adapt"
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/consts"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"github.com/wenlng/go-captcha/v2/rotate"
	"go.uber.org/zap"
)

// RotateCaptLogic .
type RotateCaptLogic struct {
	svcCtx *common.SvcContext

	cacheMgr   *cache.CacheManager
	dynamicCfg *config.DynamicConfig
	logger     *zap.Logger
	captcha    *gocaptcha.GoCaptcha
}

// NewRotateCaptLogic .
func NewRotateCaptLogic(svcCtx *common.SvcContext) *RotateCaptLogic {
	return &RotateCaptLogic{
		svcCtx:     svcCtx,
		cacheMgr:   svcCtx.CacheMgr,
		dynamicCfg: svcCtx.DynamicConfig,
		logger:     svcCtx.Logger,
		captcha:    svcCtx.Captcha,
	}
}

// GetData .
func (cl *RotateCaptLogic) GetData(ctx context.Context, id string) (res *adapt.CaptData, err error) {
	res = &adapt.CaptData{}

	if id == "" {
		return nil, fmt.Errorf("missing id parameter")
	}

	var capt *gocaptcha.RotateCaptInstance
	ttype := cl.svcCtx.Captcha.GetCaptTypeWithKey(id)
	switch ttype {
	case consts.GoCaptchaTypeRotate:
		capt = cl.svcCtx.Captcha.GetRotateInstanceWithKey(id)
		break
	}
	if capt == nil || capt.Instance == nil {
		return nil, fmt.Errorf("missing captcha type")
	}

	captData, err := capt.Instance.Generate()
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
		Type:   ttype,
		Status: 0,
	}
	cacheDataByte, err := json.Marshal(cacheData)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal: %v", err)
	}

	key, err := helper.GenerateIDWithNode(cl.dynamicCfg.Get().ServiceNode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %v", err)
	}

	err = cl.cacheMgr.GetCache().SetCache(ctx, key, string(cacheDataByte))
	if err != nil {
		return res, fmt.Errorf("failed to write cache:: %v", err)
	}

	opts := capt.Instance.GetOptions()
	res.MasterWidth = int32(opts.GetImageSize())
	res.MasterHeight = int32(opts.GetImageSize())
	res.ThumbWidth = int32(data.Width)
	res.ThumbHeight = int32(data.Height)
	res.ThumbSize = int32(data.Width)
	res.CaptchaKey = key
	return res, nil
}

// CheckData .
func (cl *RotateCaptLogic) CheckData(ctx context.Context, key string, angle int) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("invalid key")
	}

	cacheData, err := cl.cacheMgr.GetCache().GetCache(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to get cache: %v", err)
	}

	if cacheData == "" {
		return false, nil
	}

	var cacheCaptData *cache.CaptCacheData
	err = json.Unmarshal([]byte(cacheData), &cacheCaptData)
	if err != nil {
		return false, fmt.Errorf("failed to json unmarshal: %v", err)
	}

	var dct *rotate.Block
	captDataStr, err := json.Marshal(cacheCaptData.Data)
	if err != nil {
		return false, fmt.Errorf("failed to json marshal: %v", err)
	}
	err = json.Unmarshal(captDataStr, &dct)
	if err != nil {
		return false, fmt.Errorf("failed to json unmarshal: %v", err)
	}

	if cacheCaptData.Status == 2 {
		return false, nil
	}

	ret := rotate.Validate(angle, dct.Angle, 2)

	if ret {
		cacheCaptData.Status = 1
	} else {
		cacheCaptData.Status = 2
	}

	cacheDataByte, err := json.Marshal(cacheCaptData)
	if err != nil {
		return ret, fmt.Errorf("failed to json marshal: %v", err)
	}

	err = cl.cacheMgr.GetCache().SetCache(ctx, key, string(cacheDataByte))
	if err != nil {
		return ret, fmt.Errorf("failed to update cache:: %v", err)
	}

	return ret, nil
}
