/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package gocaptcha

import (
	"image"
	"path"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
	"github.com/wenlng/go-captcha/v2/click"
)

// ClickCaptInstance .
type ClickCaptInstance struct {
	ResourcesVersion string
	Version          string
	Instance         click.Captcha
}

// genClickOptions .
func genClickOptions(conf config.ClickConfig) ([]click.Option, error) {
	options := make([]click.Option, 0)

	// Master image
	if conf.Master.ImageSize.Height > 0 && conf.Master.ImageSize.Width > 0 {
		options = append(options, click.WithImageSize(conf.Master.ImageSize))
	}

	if conf.Master.RangeLength.Min >= 0 && conf.Master.RangeLength.Max > 0 {
		options = append(options, click.WithRangeLen(conf.Master.RangeLength))
	}

	if conf.Master.RangeAngles != nil && len(conf.Master.RangeAngles) > 0 {
		options = append(options, click.WithRangeAnglePos(conf.Master.RangeAngles))
	}

	if conf.Master.RangeSize.Min >= 0 && conf.Master.RangeSize.Max > 0 {
		options = append(options, click.WithRangeSize(conf.Master.RangeSize))
	}

	if conf.Master.RangeColors != nil && len(conf.Master.RangeColors) > 0 {
		options = append(options, click.WithRangeColors(conf.Master.RangeColors))
	}

	if conf.Master.DisplayShadow != false {
		options = append(options, click.WithDisplayShadow(conf.Master.DisplayShadow))
	}

	if conf.Master.ShadowColor != "" {
		options = append(options, click.WithShadowColor(conf.Master.ShadowColor))
	}

	if conf.Master.ShadowPoint.X > -999 &&
		conf.Master.ShadowPoint.Y > -999 &&
		conf.Master.ShadowPoint.X < 999 &&
		conf.Master.ShadowPoint.Y < 999 {
		options = append(options, click.WithShadowPoint(conf.Master.ShadowPoint))
	}

	if conf.Master.ImageAlpha > 0 {
		options = append(options, click.WithImageAlpha(conf.Master.ImageAlpha))
	}

	if conf.Master.UseShapeOriginalColor != false {
		options = append(options, click.WithUseShapeOriginalColor(conf.Master.UseShapeOriginalColor))
	}

	// Thumb image
	if conf.Thumb.ImageSize.Height != 0 && conf.Thumb.ImageSize.Width != 0 {
		options = append(options, click.WithRangeThumbImageSize(conf.Thumb.ImageSize))
	}

	if conf.Thumb.RangeVerifyLength.Min != 0 && conf.Thumb.RangeVerifyLength.Max != 0 {
		options = append(options, click.WithRangeVerifyLen(conf.Thumb.RangeVerifyLength))
	}

	if conf.Thumb.DisabledRangeVerifyLength != false {
		options = append(options, click.WithDisabledRangeVerifyLen(conf.Thumb.DisabledRangeVerifyLength))
	}

	if conf.Thumb.RangeTextSize.Min > 0 && conf.Thumb.RangeTextSize.Max > 0 {
		options = append(options, click.WithRangeThumbSize(conf.Thumb.RangeTextSize))
	}

	if conf.Thumb.RangeTextColors != nil && len(conf.Thumb.RangeTextColors) > 0 {
		options = append(options, click.WithRangeThumbColors(conf.Thumb.RangeTextColors))
	}

	if conf.Thumb.RangeBackgroundColors != nil && len(conf.Thumb.RangeBackgroundColors) > 0 {
		options = append(options, click.WithRangeThumbBgColors(conf.Thumb.RangeBackgroundColors))
	}

	if conf.Thumb.BackgroundDistort > 0 {
		options = append(options, click.WithRangeThumbBgDistort(conf.Thumb.BackgroundDistort))
	}

	if conf.Thumb.BackgroundDistortAlpha > 0 {
		options = append(options, click.WithThumbDisturbAlpha(conf.Thumb.BackgroundDistortAlpha))
	}

	if conf.Thumb.BackgroundCirclesNum > 0 {
		options = append(options, click.WithRangeThumbBgCirclesNum(conf.Thumb.BackgroundCirclesNum))
	}

	if conf.Thumb.BackgroundSlimLineNum > 0 {
		options = append(options, click.WithRangeThumbBgSlimLineNum(conf.Thumb.BackgroundSlimLineNum))
	}

	if conf.Thumb.IsThumbNonDeformAbility != false {
		options = append(options, click.WithIsThumbNonDeformAbility(conf.Thumb.IsThumbNonDeformAbility))
	}

	return options, nil
}

// GetMixinAlphaChars .
func GetMixinAlphaChars() []string {
	var ret = make([]string, 0)
	letterArr := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	numArr := strings.Split("0123456789", "")

	for _, s := range letterArr {
		for _, n := range numArr {
			ret = append(ret, s+n)
		}
	}

	for _, s := range numArr {
		for _, n := range letterArr {
			ret = append(ret, s+n)
		}
	}

	return ret
}

// genClickResources.
func genClickResources(conf config.ClickConfig, resources config.ResourceConfig) ([]click.Resource, error) {
	newResources := make([]click.Resource, 0)

	// Set chars resources
	if newChars, ok := resources.Char.Languages[conf.Language]; ok && len(newChars) > 0 {
		newResources = append(newResources, click.WithChars(newChars))
	} else {
		if conf.Language == LanguageNameChinese {
			newResources = append(newResources, click.WithChars(chars.GetChineseChars()))
		} else {
			newResources = append(newResources, click.WithChars(GetMixinAlphaChars()))
		}
	}

	// Set fonts resources
	if len(resources.Font.FileMaps) > 0 {
		var newFonts = make([]*truetype.Font, 0)
		for _, file := range resources.Font.FileMaps {
			resourcesPath := helper.GetResourceDirAbsPath()
			rootDir := resources.Font.FileDir
			filepath := path.Join(resourcesPath, rootDir, file)
			stream, err := helper.ReadFileStream(filepath)
			if err != nil {
				return nil, err
			}

			font, err := freetype.ParseFont(stream)
			if err != nil {
				return nil, err
			}

			newFonts = append(newFonts, font)
		}

		newResources = append(newResources, click.WithFonts(newFonts))
	} else {
		fonts, err := fzshengsksjw.GetFont()
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, click.WithFonts([]*truetype.Font{fonts}))
	}

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
		newResources = append(newResources, click.WithBackgrounds(newImages))
	} else {
		imgs, err := images.GetImages()
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, click.WithBackgrounds(imgs))
	}

	// Set Thumb images resources
	if len(resources.ThumbImage.FileMaps) > 0 {
		var newImages = make([]image.Image, 0)
		for _, file := range resources.ThumbImage.FileMaps {
			resourcesPath := helper.GetResourceDirAbsPath()
			rootDir := resources.ThumbImage.FileDir
			filepath := path.Join(resourcesPath, rootDir, file)

			img, err := helper.LoadImageData(filepath)
			if err != nil {
				return nil, err
			}
			newImages = append(newImages, img)
		}
		newResources = append(newResources, click.WithThumbBackgrounds(newImages))
	} else {
		//imgs, err := thumbs.GetThumbs()
		//if err != nil {
		//	return nil, err
		//}
		//newResources = append(newResources, click.WithThumbBackgrounds(imgs))
	}

	return newResources, nil
}

// setupClickCapt
func setupClickCapt(conf config.ClickConfig, resources config.ResourceConfig) (capt click.Captcha, err error) {
	newOptions, err := genClickOptions(conf)
	if err != nil {
		return nil, err
	}

	newResources, err := genClickResources(conf, resources)
	if err != nil {
		return nil, err
	}

	// builder
	builder := click.NewBuilder(newOptions...)
	builder.SetResources(newResources...)

	return builder.Make(), nil
}

func setupClickShapeCapt(conf config.ClickConfig, resources config.ResourceConfig) (capt click.Captcha, err error) {
	newOptions, err := genClickOptions(conf)
	if err != nil {
		return nil, err
	}

	newResources, err := genClickResources(conf, resources)
	if err != nil {
		return nil, err
	}

	// Set Shape images resources
	if len(resources.ShapeImage.FileMaps) > 0 {
		var newImageMaps = make(map[string]image.Image, 0)
		for name, file := range resources.ShapeImage.FileMaps {
			resourcesPath := helper.GetResourceDirAbsPath()
			rootDir := resources.ShapeImage.FileDir
			filepath := path.Join(resourcesPath, rootDir, file)

			img, err := helper.LoadImageData(filepath)
			if err != nil {
				return nil, err
			}
			newImageMaps[name] = img
		}
		newResources = append(newResources, click.WithShapes(newImageMaps))
	} else {
		imgs, err := shapes.GetShapes()
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, click.WithShapes(imgs))
	}

	// builder
	builder := click.NewBuilder(newOptions...)
	builder.SetResources(newResources...)

	return builder.MakeWithShape(), nil
}
