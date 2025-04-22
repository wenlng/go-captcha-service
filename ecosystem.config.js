module.exports = {
  apps: [{
    name: 'go-captcha-service',
    script: './build/go-captcha-service',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '1G',
    env: {
      CONFIG: 'config.json',
      GO_CAPTCHA_CONFIG: 'gocaptcha.json',
      SERVICE_NAME: 'go-captcha-service',
      CACHE_TYPE: 'redis',
      CACHE_ADDRS: 'localhost:6379',
    },
    env_production: {
      CONFIG: 'config.json',
      GO_CAPTCHA_CONFIG: 'gocaptcha.json',
      SERVICE_NAME: 'go-captcha-service',
      CACHE_TYPE: 'redis',
      CACHE_ADDRS: 'localhost:6379',
    }
  }]
};