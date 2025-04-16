package gocaptcha

import (
	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
)

func setupRotate() (capt rotate.Captcha, err error) {
	builder := rotate.NewBuilder(
		rotate.WithRangeAnglePos([]option.RangeVal{
			{Min: 20, Max: 330},
		}),
	)

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		return nil, err
	}

	// set resources
	builder.SetResources(
		rotate.WithImages(imgs),
	)

	return builder.Make(), nil
}
