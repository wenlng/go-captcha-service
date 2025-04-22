/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package logic

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
)

// ResourceLogic .
type ResourceLogic struct {
	svcCtx *common.SvcContext

	cacheMgr   *cache.CacheManager
	dynamicCfg *config.DynamicConfig
	logger     *zap.Logger
	captcha    *gocaptcha.GoCaptcha
}

// NewResourceLogic .
func NewResourceLogic(svcCtx *common.SvcContext) *ResourceLogic {
	return &ResourceLogic{
		svcCtx:     svcCtx,
		cacheMgr:   svcCtx.CacheMgr,
		dynamicCfg: svcCtx.DynamicConfig,
		logger:     svcCtx.Logger,
		captcha:    svcCtx.Captcha,
	}
}

// SaveResource .
func (cl *ResourceLogic) SaveResource(ctx context.Context, dirname string, files []*multipart.FileHeader) (ret, allDone bool, err error) {
	resourcePath := helper.GetResourceDirAbsPath()
	dirPath := filepath.Join(resourcePath, dirname)
	dirPath = filepath.Clean(dirPath)

	if !helper.IsSubPath(resourcePath, dirPath) {
		return false, false, fmt.Errorf("invalid dirpath")
	}

	err = helper.EnsureDir(dirPath)
	if err != nil {
		return false, false, err
	}

	var hasSkipFileSave bool
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return false, false, fmt.Errorf("failed to open file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		filename := filepath.Base(fileHeader.Filename)
		if filename == "" {
			return false, false, fmt.Errorf("invalid filename")
		}

		dstPath := filepath.Join(dirPath, filename)

		if helper.FileExists(dstPath) {
			hasSkipFileSave = true
			continue
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return false, false, fmt.Errorf("failed to create file %s: %v", filename, err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			return false, false, fmt.Errorf("failed to save file %s: %v", filename, err)
		}
	}

	if hasSkipFileSave {
		return true, false, fmt.Errorf("some files failed to be uploaded. check if they already exist")
	}

	return true, true, nil
}

// GetResourceList .
func (cl *ResourceLogic) GetResourceList(ctx context.Context, filepath string) ([]string, error) {
	resourcePath := helper.GetResourceDirAbsPath()
	filepath = path.Join(resourcePath, filepath)
	filepath = path.Clean(filepath)

	if !helper.IsSubPath(resourcePath, path.Dir(filepath)) {
		return nil, fmt.Errorf("invalid filepath")
	}

	fileList, err := helper.TraverseDir(filepath, resourcePath)
	if err != nil {
		return nil, nil
	}

	return fileList, nil
}

// DelResource .
func (cl *ResourceLogic) DelResource(ctx context.Context, filepath string) (ret bool, err error) {
	resourcePath := helper.GetResourceDirAbsPath()
	filepath = path.Join(resourcePath, filepath)
	filepath = path.Clean(filepath)

	if !helper.IsSubPath(resourcePath, path.Dir(filepath)) {
		return false, fmt.Errorf("invalid filepath")
	}

	if helper.FileExists(filepath) {
		err = helper.DeleteFile(filepath)
		if err != nil {
			cl.logger.Error("failed to delete resource, err: ", zap.Error(err))
			return false, err
		}
	} else {
		return false, nil
	}

	return true, nil
}
