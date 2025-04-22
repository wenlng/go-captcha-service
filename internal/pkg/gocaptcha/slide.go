/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package gocaptcha

import (
	"image"
	"path"

	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
	"github.com/wenlng/go-captcha/v2/base/codec"
	"github.com/wenlng/go-captcha/v2/slide"
)

// SlideCaptInstance .
type SlideCaptInstance struct {
	ResourcesVersion string
	Version          string
	Instance         slide.Captcha
}

// genSlideOptions .
func genSlideOptions(conf config.SlideConfig) ([]slide.Option, error) {
	options := make([]slide.Option, 0)

	// Master image
	if conf.Master.ImageSize.Height > 0 && conf.Master.ImageSize.Width > 0 {
		options = append(options, slide.WithImageSize(conf.Master.ImageSize))
	}

	if conf.Master.ImageAlpha > 0 {
		options = append(options, slide.WithImageAlpha(conf.Master.ImageAlpha))
	}

	// Thumb image
	if conf.Thumb.RangeGraphSizes.Min >= 0 && conf.Thumb.RangeGraphSizes.Max > 0 {
		options = append(options, slide.WithRangeGraphSize(conf.Thumb.RangeGraphSizes))
	}

	if conf.Thumb.RangeGraphAngles != nil && len(conf.Thumb.RangeGraphAngles) > 0 {
		options = append(options, slide.WithRangeGraphAnglePos(conf.Thumb.RangeGraphAngles))
	}

	if conf.Thumb.GenerateGraphNumber > 0 {
		options = append(options, slide.WithGenGraphNumber(conf.Thumb.GenerateGraphNumber))
	}

	if conf.Thumb.EnableGraphVerticalRandom != false {
		options = append(options, slide.WithEnableGraphVerticalRandom(conf.Thumb.EnableGraphVerticalRandom))
	}

	if conf.Thumb.RangeDeadZoneDirections != nil && len(conf.Thumb.RangeDeadZoneDirections) > 0 {
		var list = make([]slide.DeadZoneDirectionType, 0, len(conf.Thumb.RangeDeadZoneDirections))

		for _, direction := range conf.Thumb.RangeDeadZoneDirections {
			if direction == "left" {
				list = append(list, slide.DeadZoneDirectionTypeLeft)
			} else if direction == "right" {
				list = append(list, slide.DeadZoneDirectionTypeRight)
			} else if direction == "top" {
				list = append(list, slide.DeadZoneDirectionTypeTop)
			} else if direction == "bottom" {
				list = append(list, slide.DeadZoneDirectionTypeBottom)
			}
		}

		options = append(options, slide.WithRangeDeadZoneDirections(list))
	}

	return options, nil
}

// genSlideResources.
func genSlideResources(conf config.SlideConfig, resources config.ResourceConfig) ([]slide.Resource, error) {
	newResources := make([]slide.Resource, 0)

	// Set Background images resources
	if len(resources.MasterImage.FileMaps) > 0 {
		var newImages = make([]image.Image, 0)
		for _, file := range resources.MasterImage.FileMaps {
			resourcesPath := helper.GetResourceDirAbsPath()
			rootDir := resources.MasterImage.FileDir
			filepath := path.Join(resourcesPath, rootDir, file)
			stream, err := helper.ReadFileStream(filepath)
			if err != nil {
				return nil, err
			}

			if path.Ext(file) == ".png" {
				png, err := codec.DecodeByteToPng(stream)
				if err != nil {
					return nil, err
				}
				newImages = append(newImages, png)
			} else {
				jpeg, err := codec.DecodeByteToJpeg(stream)
				if err != nil {
					return nil, err
				}
				newImages = append(newImages, jpeg)
			}
		}
		newResources = append(newResources, slide.WithBackgrounds(newImages))
	} else {
		imgs, err := images.GetImages()
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, slide.WithBackgrounds(imgs))
	}

	// Set Tile images resources
	var isSetGraphsDonne bool
	if len(resources.TileImage.FileMaps) > 0 {
		var newGraphs = make([]*slide.GraphImage, 0, len(resources.TileImage.FileMaps))
		for name, file := range resources.TileImage.FileMaps {
			resourcesPath := helper.GetResourceDirAbsPath()
			rootDir := resources.TileImage.FileDir
			overlayImageFilepath := path.Join(resourcesPath, rootDir, file)

			shadowImageFilepath, ok := resources.TileImage.FileMaps02[name]
			if !ok {
				break
			}
			shadowImageFilepath = path.Join(resourcesPath, rootDir, shadowImageFilepath)

			maskImageFilepath, ok := resources.TileImage.FileMaps03[name]
			if !ok {
				break
			}
			maskImageFilepath = path.Join(resourcesPath, rootDir, maskImageFilepath)

			graph := &slide.GraphImage{}

			// OverlayImage
			overlayImage, err := helper.LoadImageData(overlayImageFilepath)
			if err != nil {
				return nil, err
			}
			graph.OverlayImage = overlayImage

			// ShadowImage
			shadowImage, err := helper.LoadImageData(shadowImageFilepath)
			if err != nil {
				return nil, err
			}
			graph.ShadowImage = shadowImage

			// MaskImage
			maskImage, err := helper.LoadImageData(maskImageFilepath)
			if err != nil {
				return nil, err
			}
			graph.MaskImage = maskImage

			newGraphs = append(newGraphs, graph)
		}

		if len(newGraphs) > 0 {
			isSetGraphsDonne = true
			newResources = append(newResources, slide.WithGraphImages(newGraphs))
		}
	}

	if !isSetGraphsDonne {
		graphs, err := tiles.GetTiles()
		if err != nil {
			return nil, err
		}

		var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
		for i := 0; i < len(graphs); i++ {
			graph := graphs[i]
			newGraphs = append(newGraphs, &slide.GraphImage{
				OverlayImage: graph.OverlayImage,
				MaskImage:    graph.MaskImage,
				ShadowImage:  graph.ShadowImage,
			})
		}
		newResources = append(newResources, slide.WithGraphImages(newGraphs))
	}

	return newResources, nil
}

// setupSlideCapt
func setupSlideCapt(conf config.SlideConfig, resources config.ResourceConfig) (capt slide.Captcha, err error) {
	newOptions, err := genSlideOptions(conf)
	if err != nil {
		return nil, err
	}

	newResources, err := genSlideResources(conf, resources)
	if err != nil {
		return nil, err
	}

	// builder
	builder := slide.NewBuilder(newOptions...)
	builder.SetResources(newResources...)

	return builder.Make(), nil
}

func setupDragCapt(conf config.SlideConfig, resources config.ResourceConfig) (capt slide.Captcha, err error) {
	newOptions, err := genSlideOptions(conf)
	if err != nil {
		return nil, err
	}

	newResources, err := genSlideResources(conf, resources)
	if err != nil {
		return nil, err
	}

	// builder
	builder := slide.NewBuilder(newOptions...)
	builder.SetResources(newResources...)

	return builder.MakeWithRegion(), nil
}
