import { defineStore } from 'pinia'
import { request } from '@/utils'

export const useSystemConfigStore = defineStore('system-config', {
  state: () => ({
    siteName: import.meta.env.VITE_TITLE || 'Kite Admin',
    logo: '',
    favicon: '',
    copyright: '',
    loaded: false,
  }),

  actions: {
    async fetchConfig() {
      try {
        const { data } = await request.get('/auth/system/config')
        if (data) {
          this.siteName = data.siteName || this.siteName
          this.logo = data.logo || ''
          this.favicon = data.favicon || ''
          this.copyright = data.copyright || ''
          this.loaded = true
          this.applyFavicon()
          this.applyTitle()
        }
      } catch {
        // ignore, use defaults
      }
    },

    applyFavicon() {
      if (!this.favicon) return
      let link = document.querySelector('link[rel="icon"]')
      if (!link) {
        link = document.createElement('link')
        link.rel = 'icon'
        document.head.appendChild(link)
      }
      link.href = this.favicon
    },

    applyTitle() {
      if (this.siteName) {
        document.title = this.siteName
      }
    },
  },
})
