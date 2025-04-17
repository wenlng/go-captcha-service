package gocaptcha

import (
	"sync"

	"github.com/wenlng/go-captcha-service/internal/consts"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
)

const (
	LanguageNameChinese = "chinese"
	LanguageNameEnglish = "english"
)

// GoCaptcha .
type GoCaptcha struct {
	DynamicCnf *config.DynamicCaptchaConfig

	clickInstanceMaps      map[string]*ClickCaptInstance
	clickShapeInstanceMaps map[string]*ClickCaptInstance
	slideInstanceMaps      map[string]*SlideCaptInstance
	dragInstanceMaps       map[string]*SlideCaptInstance
	rotateInstanceMaps     map[string]*RotateCaptInstance
	keyMaps                map[string]int

	clickInstanceMutex      sync.RWMutex
	clickShapeInstanceMutex sync.RWMutex
	slideInstanceMutex      sync.RWMutex
	dragInstanceMutex       sync.RWMutex
	rotateInstanceMutex     sync.RWMutex
}

// newGoCaptcha .
func newGoCaptcha() *GoCaptcha {
	return &GoCaptcha{
		clickInstanceMaps:      make(map[string]*ClickCaptInstance, 0),
		clickShapeInstanceMaps: make(map[string]*ClickCaptInstance, 0),
		slideInstanceMaps:      make(map[string]*SlideCaptInstance, 0),
		dragInstanceMaps:       make(map[string]*SlideCaptInstance, 0),
		rotateInstanceMaps:     make(map[string]*RotateCaptInstance, 0),
		keyMaps:                make(map[string]int),
	}
}

// GetCaptTypeWithKey .
func (gc *GoCaptcha) GetCaptTypeWithKey(key string) int {
	t, ok := gc.keyMaps[key]
	if ok {
		return t
	}

	return consts.GoCaptchaTypeUnknown
}

// GetClickInstanceWithKey .
func (gc *GoCaptcha) GetClickInstanceWithKey(key string) *ClickCaptInstance {
	gc.clickInstanceMutex.RLock()
	defer gc.clickInstanceMutex.RUnlock()
	return gc.clickInstanceMaps[key]
}

// GetClickShapeInstanceWithKey .
func (gc *GoCaptcha) GetClickShapeInstanceWithKey(key string) *ClickCaptInstance {
	gc.clickInstanceMutex.RLock()
	defer gc.clickInstanceMutex.RUnlock()
	return gc.clickShapeInstanceMaps[key]
}

// GetSlideInstanceWithKey .
func (gc *GoCaptcha) GetSlideInstanceWithKey(key string) *SlideCaptInstance {
	gc.clickInstanceMutex.RLock()
	defer gc.clickInstanceMutex.RUnlock()
	return gc.slideInstanceMaps[key]
}

// GetDragInstanceWithKey .
func (gc *GoCaptcha) GetDragInstanceWithKey(key string) *SlideCaptInstance {
	gc.clickInstanceMutex.RLock()
	defer gc.clickInstanceMutex.RUnlock()
	return gc.dragInstanceMaps[key]
}

// GetRotateInstanceWithKey .
func (gc *GoCaptcha) GetRotateInstanceWithKey(key string) *RotateCaptInstance {
	gc.clickInstanceMutex.RLock()
	defer gc.clickInstanceMutex.RUnlock()
	return gc.rotateInstanceMaps[key]
}

// UpdateClickInstance .
func (gc *GoCaptcha) UpdateClickInstance(configMaps map[string]config.ClickConfig, resources config.ResourceConfig) error {
	gc.clickInstanceMutex.Lock()
	defer gc.clickInstanceMutex.Unlock()

	for key, cnf := range configMaps {
		ci, ok := gc.clickInstanceMaps[key]

		if !ok || ci.ResourcesVersion != resources.Version || ci.Version != cnf.Version {
			instance, err := setupClickCapt(cnf, resources)
			if err != nil {
				return err
			}
			gc.clickInstanceMaps[key] = &ClickCaptInstance{
				ResourcesVersion: resources.Version,
				Version:          cnf.Version,
				Instance:         instance,
			}
			gc.keyMaps[key] = consts.GoCaptchaTypeClick
		}
	}
	return nil
}

// UpdateClickShapeInstance .
func (gc *GoCaptcha) UpdateClickShapeInstance(configMaps map[string]config.ClickConfig, resources config.ResourceConfig) error {
	gc.clickShapeInstanceMutex.Lock()
	defer gc.clickShapeInstanceMutex.Unlock()

	for key, cnf := range configMaps {
		ci, ok := gc.clickShapeInstanceMaps[key]

		if !ok || ci.ResourcesVersion != resources.Version || ci.Version != cnf.Version {
			instance, err := setupClickShapeCapt(cnf, resources)
			if err != nil {
				return err
			}
			gc.clickShapeInstanceMaps[key] = &ClickCaptInstance{
				ResourcesVersion: resources.Version,
				Version:          cnf.Version,
				Instance:         instance,
			}
			gc.keyMaps[key] = consts.GoCaptchaTypeClickShape
		}
	}
	return nil
}

// UpdateSlideInstance .
func (gc *GoCaptcha) UpdateSlideInstance(configMaps map[string]config.SlideConfig, resources config.ResourceConfig) error {
	gc.slideInstanceMutex.Lock()
	defer gc.slideInstanceMutex.Unlock()

	for key, cnf := range configMaps {
		ci, ok := gc.slideInstanceMaps[key]

		if !ok || ci.ResourcesVersion != resources.Version || ci.Version != cnf.Version {
			instance, err := setupSlideCapt(cnf, resources)
			if err != nil {
				return err
			}
			gc.slideInstanceMaps[key] = &SlideCaptInstance{
				ResourcesVersion: resources.Version,
				Version:          cnf.Version,
				Instance:         instance,
			}
			gc.keyMaps[key] = consts.GoCaptchaTypeSlide
		}
	}
	return nil
}

// UpdateDragInstance .
func (gc *GoCaptcha) UpdateDragInstance(configMaps map[string]config.SlideConfig, resources config.ResourceConfig) error {
	gc.dragInstanceMutex.Lock()
	defer gc.dragInstanceMutex.Unlock()

	for key, cnf := range configMaps {
		ci, ok := gc.dragInstanceMaps[key]

		if !ok || ci.ResourcesVersion != resources.Version || ci.Version != cnf.Version {
			instance, err := setupDragCapt(cnf, resources)
			if err != nil {
				return err
			}
			gc.dragInstanceMaps[key] = &SlideCaptInstance{
				ResourcesVersion: resources.Version,
				Version:          cnf.Version,
				Instance:         instance,
			}
			gc.keyMaps[key] = consts.GoCaptchaTypeDrag
		}
	}
	return nil
}

// UpdateRotateInstance .
func (gc *GoCaptcha) UpdateRotateInstance(configMaps map[string]config.RotateConfig, resources config.ResourceConfig) error {
	gc.rotateInstanceMutex.Lock()
	defer gc.rotateInstanceMutex.Unlock()

	for key, cnf := range configMaps {
		ci, ok := gc.rotateInstanceMaps[key]

		if !ok || ci.ResourcesVersion != resources.Version || ci.Version != cnf.Version {
			instance, err := setupRotateCapt(cnf, resources)
			if err != nil {
				return err
			}
			gc.rotateInstanceMaps[key] = &RotateCaptInstance{
				ResourcesVersion: resources.Version,
				Version:          cnf.Version,
				Instance:         instance,
			}
			gc.keyMaps[key] = consts.GoCaptchaTypeRotate
		}
	}
	return nil
}

// HotUpdate .
func (gc *GoCaptcha) HotUpdate(dyCnf *config.DynamicCaptchaConfig) error {
	cnf := dyCnf.Get()

	var err error
	err = gc.UpdateClickInstance(cnf.Builder.ClickConfigMaps, cnf.Resources)
	if err != nil {
		return err
	}

	err = gc.UpdateClickShapeInstance(cnf.Builder.ClickShapeConfigMaps, cnf.Resources)
	if err != nil {
		return err
	}

	err = gc.UpdateSlideInstance(cnf.Builder.SlideConfigMaps, cnf.Resources)
	if err != nil {
		return err
	}

	err = gc.UpdateDragInstance(cnf.Builder.DragConfigMaps, cnf.Resources)
	if err != nil {
		return err
	}

	err = gc.UpdateRotateInstance(cnf.Builder.RotateConfigMaps, cnf.Resources)
	if err != nil {
		return err
	}

	return nil
}

// Setup initializes the captcha
func Setup(dyCnf *config.DynamicCaptchaConfig) (*GoCaptcha, error) {
	gc := newGoCaptcha()
	err := gc.HotUpdate(dyCnf)
	if err != nil {
		return nil, err
	}

	return gc, nil
}
