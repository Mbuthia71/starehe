module.exports = {
  apps: [
    {
      name: 'starehian-api',
      script: '/root/starehian/starehian-api',
      env: {
        DATABASE_URL: 'postgres://starehe_user:changeme@localhost:5432/starehe_db',
        REDIS_URL: 'redis://localhost:6379',
        JWT_SECRET: 'starehian_jwt_secret_key_2024',
        REFRESH_TOKEN_SECRET: 'starehian_refresh_secret_key_2024',
        CENTRIFUGO_SECRET: 'starehian_centrifugo_secret_2024',
        CENTRIFUGO_API_KEY: 'starehian_centrifugo_api_key_2024',
        CENTRIFUGO_URL: 'http://localhost:8000/api',
        CENTRIFUGO_WS_URL: 'ws://localhost:8000/connection/websocket'
      }
    },
    {
      name: 'starehe-frontend',
      script: '/var/www/starehe-frontend/.output/server/index.mjs',
      cwd: '/var/www/starehe-frontend',
      env: {
        PORT: 3001,
        HOST: '0.0.0.0',
        NODE_ENV: 'production'
      },
      max_restarts: 10,
      min_uptime: '10s',
      restart_delay: 4000
    },
    {
      name: 'centrifugo',
      script: 'centrifugo',
      args: '--config=/root/starehian/configs/centrifugo.json',
      interpreter: 'none',
      env: {
        CENTRIFUGO_SECRET: 'starehian_centrifugo_secret_2024',
        CENTRIFUGO_API_KEY: 'starehian_centrifugo_api_key_2024'
      }
    }
  ]
}
