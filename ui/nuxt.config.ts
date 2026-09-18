// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    'nuxt-charts'
  ],

  // Static single-page app: `nuxt generate` emits plain HTML/JS/CSS into
  // .output/public, served by nginx in production (see Dockerfile).
  ssr: false,

  devtools: {
    enabled: true
  },

  css: ['~/assets/css/main.css'],

  compatibilityDate: '2026-06-30',

  nitro: {
    // Dev-only: forward /api/* to the Go API so the browser stays same-origin.
    // In production nginx does this (see nginx/default.conf.template).
    // The proxy strips the matched prefix, so the target must include /api.
    devProxy: {
      '/api': {
        target: `${process.env.NUXT_API_URL || 'http://localhost:8080'}/api`,
        changeOrigin: true
      }
    }
  },

  eslint: {
    config: {
      stylistic: {
        commaDangle: 'never',
        braceStyle: '1tbs'
      }
    }
  }
})
