// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  nitro: {
    devProxy: {
      '/api': { 
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  runtimeConfig:{
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE,
    }
  }
})
