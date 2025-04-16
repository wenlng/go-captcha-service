package gocaptcha

import (
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/shapes"

	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
)

func setupClick() (capt click.Captcha, err error) {
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
	)

	// fonts
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		return nil, err
	}

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		return nil, err
	}

	// set resources
	builder.SetResources(
		click.WithChars(chars.GetChineseChars()),
		click.WithFonts([]*truetype.Font{fonts}),
		click.WithBackgrounds(imgs),
	)

	return builder.Make(), nil
}

func setupClickShape() (capt click.Captcha, err error) {
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 3, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 3}),
		click.WithRangeThumbBgDistort(1),
		click.WithIsThumbNonDeformAbility(true),
	)

	// shape
	shapeMaps, err := shapes.GetShapes()
	if err != nil {
		return nil, err
	}

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		return nil, err
	}

	// set resources
	builder.SetResources(
		click.WithShapes(shapeMaps),
		click.WithBackgrounds(imgs),
	)

	return builder.Make(), nil
}
