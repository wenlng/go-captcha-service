/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

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
	ConfigVersion int64          `json:"config_version"`
	Resources     ResourceConfig `json:"resources"`
	Builder       BuilderConfig  `json:"builder"`
}

// DynamicCaptchaConfig .
type DynamicCaptchaConfig struct {
	Config      CaptchaConfig
	mu          sync.RWMutex
	hotCbsHooks map[string]HandleHotCallbackFnc

	outputLogCbs helper.OutputLogCallback
}

type HotCallbackType int

const (
	HotCallbackTypeLocalConfigFile HotCallbackType = 1
	HotCallbackTypeRemoteConfig                    = 2
)

type HandleHotCallbackFnc = func(*DynamicCaptchaConfig, HotCallbackType)

// NewDynamicConfig .
func NewDynamicConfig(file string, hasWatchFile bool) (*DynamicCaptchaConfig, error) {
	cfg := DefaultConfig()
	var err error

	if file != "" {
		cfg, err = Load(file)
		if err != nil {
			return nil, err
		}
	}

	dc := &DynamicCaptchaConfig{Config: cfg, hotCbsHooks: make(map[string]HandleHotCallbackFnc)}

	if hasWatchFile {
		go dc.watchFile(file)
	}

	return dc, nil
}

// DefaultDynamicConfig .
func DefaultDynamicConfig() *DynamicCaptchaConfig {
	cfg := DefaultConfig()
	return &DynamicCaptchaConfig{Config: cfg, hotCbsHooks: make(map[string]HandleHotCallbackFnc)}
}

// SetOutputLogCallback Set the log out hook function
func (dc *DynamicCaptchaConfig) SetOutputLogCallback(outputLogCbs helper.OutputLogCallback) {
	dc.outputLogCbs = outputLogCbs
}

// outLog ..
func (dc *DynamicCaptchaConfig) outLog(logType helper.OutputLogType, message string) {
	if dc.outputLogCbs != nil {
		dc.outputLogCbs(logType, message)
	}
}

// Get ..
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

// MarshalConfig ..
func (dc *DynamicCaptchaConfig) MarshalConfig() (string, error) {
	dc.mu.RLock()
	cByte, err := json.Marshal(dc.Config)
	if err != nil {
		return "", err
	}
	dc.mu.RUnlock()

	return string(cByte), nil
}

// UnMarshalConfig ..
func (dc *DynamicCaptchaConfig) UnMarshalConfig(str string) error {
	var config CaptchaConfig
	err := json.Unmarshal([]byte(str), &config)
	if err != nil {
		return err
	}

	dc.mu.Lock()
	dc.Config = config
	dc.mu.Unlock()

	return nil
}

// RegisterHotCallback callback when updating configuration
func (dc *DynamicCaptchaConfig) RegisterHotCallback(key string, callback HandleHotCallbackFnc) {
	if _, ok := dc.hotCbsHooks[key]; !ok {
		dc.hotCbsHooks[key] = callback
	}
}

// UnRegisterHotCallback callback when updating configuration
func (dc *DynamicCaptchaConfig) UnRegisterHotCallback(key string) {
	if _, ok := dc.hotCbsHooks[key]; !ok {
		delete(dc.hotCbsHooks, key)
	}
}

// HandleHotCallback .
func (dc *DynamicCaptchaConfig) HandleHotCallback(hostType HotCallbackType) {
	for _, fnc := range dc.hotCbsHooks {
		if fnc != nil {
			fnc(dc, hostType)
		}
	}
}

// watchFile monitors the CaptchaConfig file for changes
func (dc *DynamicCaptchaConfig) watchFile(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		dc.outLog(helper.OutputLogTypeError, fmt.Sprintf("[CaptchaConfig] Failed to create watcher, err: %v", err))
		return
	}
	defer watcher.Close()

	absPath, err := filepath.Abs(file)
	if err != nil {
		dc.outLog(helper.OutputLogTypeError, fmt.Sprintf("[CaptchaConfig] Failed to get absolute path, err: %v", err))
		return
	}
	dir := filepath.Dir(absPath)

	if err := watcher.Add(dir); err != nil {
		dc.outLog(helper.OutputLogTypeError, fmt.Sprintf("[CaptchaConfig] Failed to watch directory, err: %v", err))
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
					dc.outLog(helper.OutputLogTypeError, fmt.Sprintf("[CaptchaConfig]Failed to reload Config, err: %v", err))
					continue
				}
				if err = dc.Update(cfg); err != nil {
					dc.outLog(helper.OutputLogTypeError, fmt.Sprintf("[CaptchaConfig] Failed to update Config, err: %v", err))
					continue
				}

				// Instance update gocaptcha
				dc.HandleHotCallback(HotCallbackTypeLocalConfigFile)
				dc.outLog(helper.OutputLogTypeInfo, "[CaptchaConfig] Configuration reloaded successfully")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			dc.outLog(helper.OutputLogTypeError, fmt.Sprintf("[CaptchaConfig] Failed to watcher, err: %v", err))
		}
	}
}

// HotUpdate hot update configuration
func (dc *DynamicCaptchaConfig) HotUpdate(cfg CaptchaConfig) error {
	if err := dc.Update(cfg); err != nil {
		return err
	}

	// Instance update gocaptcha
	dc.HandleHotCallback(HotCallbackTypeLocalConfigFile)
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
			Version: "0.0.1",
			Char: ResourceChar{
				Languages: map[string][]string{
					"chinese": make([]string, 0),
					"english": make([]string, 0),
				},
			},
			Font:        ResourceFileConfig{},
			ShapeImage:  ResourceFileConfig{},
			MasterImage: ResourceFileConfig{},
			ThumbImage:  ResourceFileConfig{},
			TileImage:   ResourceMultiFileConfig{},
		},
		Builder: BuilderConfig{
			ClickConfigMaps: map[string]ClickConfig{
				"click-default-ch": {
					Version:  "0.0.1",
					Language: "chinese",
					Master:   ClickMasterOption{},
					Thumb:    ClickThumbOption{},
				},
				"click-dark-ch": {
					Version:  "0.0.1",
					Language: "chinese",
					Master:   ClickMasterOption{},
					Thumb: ClickThumbOption{
						RangeTextColors: []string{
							"#4a85fb",
							"#d93ffb",
							"#56be01",
							"#ee2b2b",
							"#cd6904",
							"#b49b03",
							"#01ad90",
						},
					},
				},
				"click-default-en": {
					Version:  "0.0.1",
					Language: "english",
					Master:   ClickMasterOption{},
					Thumb:    ClickThumbOption{},
				},
				"click-dark-en": {
					Version:  "0.0.1",
					Language: "english",
					Master:   ClickMasterOption{},
					Thumb: ClickThumbOption{
						RangeTextColors: []string{
							"#4a85fb",
							"#d93ffb",
							"#56be01",
							"#ee2b2b",
							"#cd6904",
							"#b49b03",
							"#01ad90",
						},
					},
				},
				"click-shape-light-default": {
					Version: "0.0.1",
				},
				"click-shape-dark-default": {
					Version: "0.0.1",
				},
			},
			ClickShapeConfigMaps: map[string]ClickConfig{
				"click-shape-default": {
					Version: "0.0.1",
				},
			},
			SlideConfigMaps: map[string]SlideConfig{
				"slide-default": {
					Version: "0.0.1",
				},
			},
			DragConfigMaps: map[string]SlideConfig{
				"drag-default": {
					Version: "0.0.1",
				},
			},
			RotateConfigMaps: map[string]RotateConfig{
				"rotate-default": {
					Version: "0.0.1",
				},
			},
		},
	}
}
