package gocaptcha

import (
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
)

type GoCaptcha struct {
	ClickCaptInstance      click.Captcha
	ClickShapeCaptInstance click.Captcha
	SlideCaptInstance      slide.Captcha
	DragCaptInstance       slide.Captcha
	RotateCaptInstance     rotate.Captcha
}

func Setup() (*GoCaptcha, error) {
	var gc = &GoCaptcha{}

	cc, err := setupClick()
	if err != nil {
		return nil, err
	}
	gc.ClickCaptInstance = cc

	ccs, err := setupClickShape()
	if err != nil {
		return nil, err
	}
	gc.ClickShapeCaptInstance = ccs

	sc, err := setupSlide()
	if err != nil {
		return nil, err
	}
	gc.SlideCaptInstance = sc

	scc, err := setupDrag()
	if err != nil {
		return nil, err
	}
	gc.DragCaptInstance = scc

	rc, err := setupRotate()
	if err != nil {
		return nil, err
	}
	gc.RotateCaptInstance = rc

	return gc, nil
}
