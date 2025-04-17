package gocaptcha

import (
	"image"
	"path"

	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
	"github.com/wenlng/go-captcha/v2/rotate"
)

// RotateCaptInstance .
type RotateCaptInstance struct {
	ResourcesVersion string
	Version          string
	Instance         rotate.Captcha
}

// genRotateOptions .
func genRotateOptions(conf config.RotateConfig) ([]rotate.Option, error) {
	options := make([]rotate.Option, 0)

	// Master image
	if conf.Master.ImageSquareSize != 0 {
		options = append(options, rotate.WithImageSquareSize(conf.Master.ImageSquareSize))
	}

	// Thumb image
	if conf.Thumb.RangeAngles != nil && len(conf.Thumb.RangeAngles) > 0 {
		options = append(options, rotate.WithRangeAnglePos(conf.Thumb.RangeAngles))
	}

	if conf.Thumb.RangeImageSquareSizes != nil && len(conf.Thumb.RangeImageSquareSizes) > 0 {
		options = append(options, rotate.WithRangeThumbImageSquareSize(conf.Thumb.RangeImageSquareSizes))
	}

	if conf.Thumb.ImageAlpha > 0 {
		options = append(options, rotate.WithThumbImageAlpha(conf.Thumb.ImageAlpha))
	}

	return options, nil
}

// genRotateResources.
func genRotateResources(conf config.RotateConfig, resources config.ResourceConfig) ([]rotate.Resource, error) {
	newResources := make([]rotate.Resource, 0)

	// Set Background images resources
	if len(resources.MasterImage.FileMaps) > 0 {
		var newImages = make([]image.Image, 0)
		for _, file := range resources.MasterImage.FileMaps {
			resourcesPath := helper.GetResourceDirAbsPath()
			rootDir := resources.MasterImage.FileDir
			filepath := path.Join(resourcesPath, rootDir, file)

			img, err := helper.LoadImageData(filepath)
			if err != nil {
				return nil, err
			}
			newImages = append(newImages, img)
		}
		newResources = append(newResources, rotate.WithImages(newImages))
	} else {
		imgs, err := images.GetImages()
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, rotate.WithImages(imgs))
	}

	return newResources, nil
}

// setupRotateCapt
func setupRotateCapt(conf config.RotateConfig, resources config.ResourceConfig) (capt rotate.Captcha, err error) {
	newOptions, err := genRotateOptions(conf)
	if err != nil {
		return nil, err
	}

	newResources, err := genRotateResources(conf, resources)
	if err != nil {
		return nil, err
	}

	// builder
	builder := rotate.NewBuilder(newOptions...)
	builder.SetResources(newResources...)

	return builder.Make(), nil
}
