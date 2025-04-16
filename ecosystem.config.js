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
      CACHE_TYPE: 'redis',
      CACHE_TTL: '60',
      CACHE_CLEANUP_INTERVAL: '10',
    },
    env_production: {
      CONFIG: '/etc/go-captcha-service/config.json',
      CACHE_TYPE: 'etcd',
      CACHE_TTL: '30',
      CACHE_CLEANUP_INTERVAL: '5',
    }
  }]
};