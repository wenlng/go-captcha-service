package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/wenlng/go-captcha-service/internal/helper"
)

// BuilderConfig .
type BuilderConfig struct {
	ClickConfigMaps      map[string]ClickConfig  `json:"click_config_maps"`
	ClickShapeConfigMaps map[string]ClickConfig  `json:"click_shape_config_maps"`
	SlideConfigMaps      map[string]SlideConfig  `json:"slide_config_maps"`
	DragConfigMaps       map[string]SlideConfig  `json:"drag_config_maps"`
	RotateConfigMaps     map[string]RotateConfig `json:"rotate_config_maps"`
}

// CaptchaConfig defines the configuration structure for the gocaptcha
type CaptchaConfig struct {
	Resources ResourceConfig `json:"resources"`
	Builder   BuilderConfig  `json:"builder"`
}

// DynamicCaptchaConfig .
type DynamicCaptchaConfig struct {
	Config      CaptchaConfig
	mu          sync.RWMutex
	hotCbsHooks map[string]HandleHotCallbackHookFnc
}

type HandleHotCallbackHookFnc = func(*DynamicCaptchaConfig)

// NewDynamicConfig .
func NewDynamicConfig(file string) (*DynamicCaptchaConfig, error) {
	cfg, err := Load(file)
	if err != nil {
		return nil, err
	}
	dc := &DynamicCaptchaConfig{Config: cfg, hotCbsHooks: make(map[string]HandleHotCallbackHookFnc)}
	go dc.watchFile(file)
	return dc, nil
}

// Get retrieves the current configuration
func (dc *DynamicCaptchaConfig) Get() CaptchaConfig {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.Config
}

// Update updates the configuration
func (dc *DynamicCaptchaConfig) Update(cfg CaptchaConfig) error {
	if err := Validate(cfg); err != nil {
		return err
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.Config = cfg
	return nil
}

// RegisterHotCallbackHook callback when updating configuration
func (dc *DynamicCaptchaConfig) RegisterHotCallbackHook(key string, callback HandleHotCallbackHookFnc) {
	if _, ok := dc.hotCbsHooks[key]; !ok {
		dc.hotCbsHooks[key] = callback
	}
}

// UnRegisterHotCallbackHook callback when updating configuration
func (dc *DynamicCaptchaConfig) UnRegisterHotCallbackHook(key string) {
	if _, ok := dc.hotCbsHooks[key]; !ok {
		delete(dc.hotCbsHooks, key)
	}
}

// HandleHotCallbackHook .
func (dc *DynamicCaptchaConfig) HandleHotCallbackHook() {
	for _, fnc := range dc.hotCbsHooks {
		if fnc != nil {
			fnc(dc)
		}
	}
}

// watchFile monitors the CaptchaConfig file for changes
func (dc *DynamicCaptchaConfig) watchFile(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	absPath, err := filepath.Abs(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		return
	}
	dir := filepath.Dir(absPath)

	if err := watcher.Add(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to watch directory: %v\n", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Name == absPath && (event.Op&fsnotify.Write == fsnotify.Write) {
				cfg, err := Load(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to reload CaptchaConfig: %v\n", err)
					continue
				}
				if err := dc.Update(cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update CaptchaConfig: %v\n", err)
					continue
				}

				// Instance update gocaptcha
				dc.HandleHotCallbackHook()

				fmt.Printf("GoCaptcha Configuration reloaded successfully\n")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}

// HotUpdate hot update configuration
func (dc *DynamicCaptchaConfig) HotUpdate(cfg CaptchaConfig) error {
	if err := dc.Update(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update CaptchaConfig: %v\n", err)
		return err
	}

	// Instance update gocaptcha
	dc.HandleHotCallbackHook()
	return nil
}

// Load reads the configuration from a file
func Load(file string) (CaptchaConfig, error) {
	var config CaptchaConfig
	data, err := os.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("failed to read GoCaptcha config file: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse GoCaptcha config file: %v", err)
	}
	return config, nil
}

// Validate checks the configuration for validity
func Validate(config CaptchaConfig) error {
	filepathList := make([]string, 0, 0)
	resourcePath := helper.GetResourceDirAbsPath()

	fontConfig := config.Resources.Font
	for _, f := range fontConfig.FileMaps {
		filepathList = append(filepathList, path.Join(resourcePath, fontConfig.FileDir, f))
	}

	shapeImageConfig := config.Resources.ShapeImage
	for _, f := range shapeImageConfig.FileMaps {
		filepathList = append(filepathList, path.Join(resourcePath, shapeImageConfig.FileDir, f))
	}

	MasterImageConfig := config.Resources.MasterImage
	for _, f := range MasterImageConfig.FileMaps {
		filepathList = append(filepathList, path.Join(resourcePath, MasterImageConfig.FileDir, f))
	}

	ThumbImageConfig := config.Resources.ThumbImage
	for _, f := range ThumbImageConfig.FileMaps {
		filepathList = append(filepathList, path.Join(resourcePath, ThumbImageConfig.FileDir, f))
	}

	TileImageConfig := config.Resources.TileImage
	for _, f := range TileImageConfig.FileMaps {
		filepathList = append(filepathList, path.Join(resourcePath, TileImageConfig.FileDir, f))
	}
	for _, f := range TileImageConfig.FileMaps02 {
		filepathList = append(filepathList, path.Join(resourcePath, TileImageConfig.FileDir, f))
	}
	for _, f := range TileImageConfig.FileMaps03 {
		filepathList = append(filepathList, path.Join(resourcePath, TileImageConfig.FileDir, f))
	}

	if err := isValidFileExist(filepathList); err != nil {
		return err
	}

	return nil
}

// isValidFileExist checks if the file is existed
func isValidFileExist(filePaths []string) error {
	for _, filePath := range filePaths {
		if ok := helper.FileExists(filePath); !ok {
			return fmt.Errorf("file not exist: %s", filePath)
		} else if ok = helper.IsFile(filePath); !ok {
			return fmt.Errorf("not file type: %s", filePath)
		}
	}
	return nil
}

// DefaultConfig .
func DefaultConfig() CaptchaConfig {
	return CaptchaConfig{
		Resources: ResourceConfig{
			Char:        ResourceChar{},
			Font:        ResourceFileConfig{},
			ShapeImage:  ResourceFileConfig{},
			MasterImage: ResourceFileConfig{},
			ThumbImage:  ResourceFileConfig{},
			TileImage:   ResourceMultiFileConfig{},
		},
		Builder: BuilderConfig{
			ClickConfigMaps: map[string]ClickConfig{
				"click_default_ch": {
					Language: "chinese",
					Master:   ClickMasterOption{},
					Thumb:    ClickThumbOption{},
				},
				"click_dark_ch": {
					Language: "chinese",
					Master:   ClickMasterOption{},
					Thumb:    ClickThumbOption{},
				},
				"click_default_en": {
					Language: "english",
					Master:   ClickMasterOption{},
					Thumb:    ClickThumbOption{},
				},
				"click_dark_en": {
					Language: "english",
					Master:   ClickMasterOption{},
					Thumb:    ClickThumbOption{},
				},
				"click_shape_light_default": {},
				"click_shape_dark_default":  {},
			},
			ClickShapeConfigMaps: map[string]ClickConfig{
				"click_shape_default": {},
			},
			SlideConfigMaps: map[string]SlideConfig{
				"slide_default": {},
			},
			DragConfigMaps: map[string]SlideConfig{
				"drag_default": {},
			},
			RotateConfigMaps: map[string]RotateConfig{
				"rotate_default": {},
			},
		},
	}
}
