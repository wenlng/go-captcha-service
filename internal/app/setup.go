package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/config"
	config2 "github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
	"github.com/wenlng/go-service-link/dynaconfig"
	"github.com/wenlng/go-service-link/dynaconfig/provider"
	"github.com/wenlng/go-service-link/foundation/common"
	"github.com/wenlng/go-service-link/servicediscovery"
	"go.uber.org/zap"
)

// setupServiceDiscovery ..
func setupServiceDiscovery(dCfg *config.DynamicConfig, logger *zap.Logger) (servicediscovery.ServiceDiscovery, error) {
	cfg := dCfg.Get()

	var discovery servicediscovery.ServiceDiscovery
	if !cfg.EnableServiceDiscovery {
		return nil, nil
	}

	var sdType servicediscovery.ServiceDiscoveryType = servicediscovery.ServiceDiscoveryTypeNone
	switch cfg.ServiceDiscoveryType {
	case ServiceDiscoveryTypeEtcd:
		sdType = servicediscovery.ServiceDiscoveryTypeEtcd
		break
	case ServiceDiscoveryTypeConsul:
		sdType = servicediscovery.ServiceDiscoveryTypeConsul
		break
	case ServiceDiscoveryTypeNacos:
		sdType = servicediscovery.ServiceDiscoveryTypeNacos
		break
	case ServiceDiscoveryTypeZookeeper:
		sdType = servicediscovery.ServiceDiscoveryTypeZookeeper
		break
	}

	sdCfg := servicediscovery.Config{
		Type:           sdType,
		Addrs:          cfg.ServiceDiscoveryAddrs,
		TTL:            time.Duration(cfg.ServiceDiscoveryTTL) * time.Second,
		KeepAlive:      time.Duration(cfg.ServiceDiscoveryKeepAlive) * time.Second,
		ServiceName:    cfg.ServiceName,
		MaxRetries:     cfg.ServiceDiscoveryMaxRetries,
		BaseRetryDelay: time.Duration(cfg.ServiceDiscoveryBaseRetryDelay) * time.Millisecond,
		Username:       cfg.ServiceDiscoveryUsername,
		Password:       cfg.ServiceDiscoveryPassword,
	}
	if cfg.ServiceDiscoveryTlsCertFile != "" && cfg.ServiceDiscoveryTlsKeyFile != "" && cfg.ServiceDiscoveryTlsCaFile != "" {
		sdCfg.TlsConfig = &common.TLSConfig{
			Address:    cfg.ServiceDiscoveryTlsAddress,
			CertFile:   cfg.ServiceDiscoveryTlsCertFile,
			KeyFile:    cfg.ServiceDiscoveryTlsKeyFile,
			CAFile:     cfg.ServiceDiscoveryTlsCaFile,
			ServerName: cfg.ServiceDiscoveryTlsServerName,
		}
	}

	var err error
	discovery, err = servicediscovery.NewServiceDiscovery(sdCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize service discovery: %v", err)
	}
	discovery.SetOutputLogCallback(func(logType servicediscovery.OutputLogType, message string) {
		if logType == servicediscovery.OutputLogTypeError {
			logger.Error("[AppSetup] Service discovery error: ", zap.String("message", message))
		} else if logType == servicediscovery.OutputLogTypeWarn {
			logger.Warn("[AppSetup] Service discovery warn: ", zap.String("message", message))
		} else {
			logger.Info("[AppSetup] Service discovery info: ", zap.String("message", message))
		}
	})
	return discovery, err
}

// setupDynamicConfig ..
func setupDynamicConfig(appDynaCfg *config.DynamicConfig, captDynaCfg *config2.DynamicCaptchaConfig, logger *zap.Logger) (*dynaconfig.ConfigManager, error) {
	appCfg := appDynaCfg.Get()
	captCfg := appDynaCfg.Get()

	if !appCfg.EnableDynamicConfig {
		return nil, nil
	}

	appConfigKey := "/config/go-captcha-service/app-config"
	appConfigName := "go-captcha-service-app-config"

	captConfigKey := "/config/go-captcha-service/captcha-config"
	captConfigName := "go-captcha-service-captcha-config"

	configs := make(map[string]*provider.Config)
	configs[appConfigKey] = &provider.Config{
		Name:    appConfigName,
		Version: appDynaCfg.Get().ConfigVersion,
		Content: appCfg,
		ValidateCallback: func(config *provider.Config) (skip bool, err error) {
			if config.Content == "" {
				return false, fmt.Errorf("contnet must be not empty")
			}
			return true, nil
		},
	}

	configs[captConfigKey] = &provider.Config{
		Name:    captConfigName,
		Version: captDynaCfg.Get().ConfigVersion,
		Content: captCfg,
		ValidateCallback: func(config *provider.Config) (skip bool, err error) {
			if config.Content == "" {
				return false, fmt.Errorf("contnet must be not empty")
			}
			return true, nil
		},
	}

	keys := make([]string, 0)
	for key, _ := range configs {
		keys = append(keys, key)
	}

	var sdType provider.ProviderType
	switch appCfg.ServiceDiscoveryType {
	case ServiceDiscoveryTypeEtcd:
		sdType = provider.ProviderTypeEtcd
		break
	case ServiceDiscoveryTypeConsul:
		sdType = provider.ProviderTypeConsul
		break
	case ServiceDiscoveryTypeNacos:
		sdType = provider.ProviderTypeNacos
		break
	case ServiceDiscoveryTypeZookeeper:
		sdType = provider.ProviderTypeZookeeper
		break
	}

	providerCfg := provider.ProviderConfig{
		Type:      sdType,
		Endpoints: strings.Split(appCfg.DynamicConfigAddrs, ","),
		Username:  appCfg.DynamicConfigUsername,
		Password:  appCfg.DynamicConfigPassword,
	}
	if appCfg.DynamicConfigTlsCertFile != "" && appCfg.DynamicConfigTlsKeyFile != "" && appCfg.DynamicConfigTlsCaFile != "" {
		providerCfg.TlsConfig = &common.TLSConfig{
			Address:    appCfg.DynamicConfigTlsAddress,
			CertFile:   appCfg.DynamicConfigTlsCertFile,
			KeyFile:    appCfg.DynamicConfigTlsKeyFile,
			CAFile:     appCfg.DynamicConfigTlsCaFile,
			ServerName: appCfg.DynamicConfigTlsServerName,
		}
	}

	manager, err := dynaconfig.NewConfigManager(dynaconfig.ConfigManagerParams{
		ProviderConfig: providerCfg,
		Configs:        configs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager, err: %v ", err)
	}

	manager.SetOutputLogCallback(func(logType dynaconfig.OutputLogType, message string) {
		if logType == dynaconfig.OutputLogTypeError {
			logger.Error("[AppSetup] " + message)
		} else if logType == dynaconfig.OutputLogTypeWarn {
			logger.Warn("[AppSetup] " + message)
		} else if logType == dynaconfig.OutputLogTypeDebug {
			logger.Debug("[AppSetup] " + message)
		} else {
			logger.Info("[AppSetup] " + message)
		}
	})

	manager.Subscribe(func(key string, pConf *provider.Config) error {
		if pConf.Content == "" {
			return nil
		}

		if key == appConfigKey {
			cConf := pConf.Content
			var newConf config.Config

			captCnfStr, err := json.Marshal(cConf)
			if err != nil {
				logger.Error("[AppSetup.ConfigManager] Filed to json marshal app config: ", zap.Any("config", cConf))
				return nil
			}

			if err = json.Unmarshal(captCnfStr, &newConf); err != nil {
				logger.Error("[AppSetup.ConfigManager] Filed to json unmarshal app config: ", zap.Any("config", cConf))
				return nil
			}

			err = appDynaCfg.Update(newConf)
			if err != nil {
				logger.Error("[AppSetup.ConfigManager] Filed to update app config", zap.Error(err))
				return nil
			}
			appDynaCfg.HandleHotCallback(config.HotCallbackTypeRemoteConfig)
		} else if key == captConfigKey {
			cConf := pConf.Content
			var newConf config2.CaptchaConfig

			captCnfStr, err := json.Marshal(cConf)
			if err != nil {
				logger.Error("[AppSetup.ConfigManager] Filed to json marshal app config: ", zap.Any("config", cConf))
				return nil
			}

			if err = json.Unmarshal(captCnfStr, &newConf); err != nil {
				logger.Error("[AppSetup.ConfigManager] Filed to json unmarshal app config: ", zap.Any("config", cConf))
				return nil
			}

			err = captDynaCfg.Update(newConf)
			if err != nil {
				logger.Error("[AppSetup.ConfigManager] Filed to update captcha config", zap.Error(err))
			}
			captDynaCfg.HandleHotCallback(config2.HotCallbackTypeRemoteConfig)
		}

		return nil
	})

	appDynaCfg.RegisterHotCallback("ASYNC_APP_CONFIG", func(dynamicConfig *config.DynamicConfig, callbackType config.HotCallbackType) {
		if callbackType == config.HotCallbackTypeLocalConfigFile {
			manager.RefreshConfig(context.Background(), appConfigKey, &provider.Config{
				Name:    appConfigName,
				Version: appDynaCfg.Get().ConfigVersion,
				Content: appDynaCfg.Get(),
				ValidateCallback: func(config *provider.Config) (skip bool, err error) {
					if config.Content == "" {
						return false, fmt.Errorf("contnet must be not empty")
					}
					return true, nil
				},
			})
		}
	})

	captDynaCfg.RegisterHotCallback("ASYNC_CAPTCHA_CONFIG", func(captchaConfig *config2.DynamicCaptchaConfig, callbackType config2.HotCallbackType) {
		if callbackType == config2.HotCallbackTypeLocalConfigFile {
			manager.RefreshConfig(context.Background(), captConfigKey, &provider.Config{
				Name:    captConfigName,
				Version: captDynaCfg.Get().ConfigVersion,
				Content: captDynaCfg.Get(),
				ValidateCallback: func(config *provider.Config) (skip bool, err error) {
					if config.Content == "" {
						return false, fmt.Errorf("contnet must be not empty")
					}
					return true, nil
				},
			})
		}
	})

	manager.ASyncConfig(context.Background())

	if err = manager.Watch(); err != nil {
		return nil, fmt.Errorf("failed to start watch: %v ", err)
	}

	//////////////////////// testing /////////////////////////
	// Testing read the configuration content in real time
	//go func() {
	//	for {
	//		time.Sleep(10 * time.Second)
	//		for _, key := range keys {
	//			c := manager.GetLocalConfig(key)
	//			fmt.Printf("+++++++ >>> Current config manager -> config for %s: %+v\n\n", key, c)
	//			dc := appDynaCfg.Get()
	//			fmt.Printf("------- >>> Current app local -> config for %s: %+v\n\n\n\n", key, dc)
	//		}
	//	}
	//}()
	/////////////////////////////////////////////////

	return manager, nil
}

// setupLoggerLevel setting the log Level
func setupLoggerLevel(logger *zap.Logger, level string) {
	switch level {
	case "error":
		logger.WithOptions(zap.IncreaseLevel(zap.ErrorLevel))
		break
	case "debug":
		logger.WithOptions(zap.IncreaseLevel(zap.DebugLevel))
		break
	case "warn":
		logger.WithOptions(zap.IncreaseLevel(zap.WarnLevel))
		break
	case "info":
		logger.WithOptions(zap.IncreaseLevel(zap.InfoLevel))
		break
	}
}

// setupHealthCheck performs a health check on HTTP and gRPC servers
func setupHealthCheck(httpAddr, grpcAddr string) error {
	resp, err := http.Get("http://localhost" + httpAddr + "/status/health")
	if err != nil || resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("HTTP health check failed: %v", err)
	}
	resp.Body.Close()

	conn, err := net.DialTimeout("tcp", "localhost"+grpcAddr, 1*time.Second)
	if err != nil {
		return fmt.Errorf("gRPC health check failed: %v", err)
	}
	conn.Close()

	return nil
}

// cfg ...
func setupCacheManager(dcfg *config.DynamicConfig, logger *zap.Logger) (*cache.CacheManager, error) {
	cfg := dcfg.Get()
	// Initialize cache
	ttl := time.Duration(cfg.CacheTTL) * time.Second
	cleanInt := time.Duration(10) * time.Second // MemoryCache cleanupInterval

	cacheMgr, err := cache.NewCacheManager(&cache.CacheMgrParams{
		Type:          cache.CacheType(cfg.CacheType),
		CacheAddrs:    cfg.CacheAddrs,
		CacheUsername: cfg.CacheUsername,
		CachePassword: cfg.CachePassword,
		KeyPrefix:     cfg.CacheKeyPrefix,
		Ttl:           ttl,
		CleanInt:      cleanInt,
	})

	if err != nil {
		return nil, err
	}
	dcfg.RegisterHotCallback("UPDATE_SETUP_CACHE", func(dnCfg *config.DynamicConfig, hotType config.HotCallbackType) {
		newCfg := dnCfg.Get()
		err = cacheMgr.Setup(&cache.CacheMgrParams{
			Type:          cache.CacheType(newCfg.CacheType),
			CacheAddrs:    cfg.CacheAddrs,
			CacheUsername: cfg.CacheUsername,
			CachePassword: cfg.CachePassword,
			KeyPrefix:     newCfg.CacheKeyPrefix,
			Ttl:           ttl,
			CleanInt:      cleanInt,
		})
		if err != nil {
			logger.Error("[AppSetup] Setup cache manager", zap.Error(err))
		}
	})

	return cacheMgr, err
}
