import { defineStore } from 'pinia'
import api from '@/views/message/api'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    unreadCount: 0,
    messages: [],
    abortCtrl: null,
  }),

  actions: {
    async fetchUnreadCount() {
      const { data } = await api.getUnreadCount()
      this.unreadCount = data?.count ?? 0
    },

    async fetchRecentMessages() {
      const { data } = await api.getMyMessages({ pageNo: 1, pageSize: 5 })
      this.messages = data?.pageData ?? []
    },

    async markAsRead(id) {
      await api.markRead(id)
      if (this.unreadCount > 0) this.unreadCount--
      const msg = this.messages.find(m => m.id === id)
      if (msg) msg.isRead = true
    },

    async markAllAsRead() {
      await api.markAllRead()
      this.unreadCount = 0
      this.messages.forEach(m => { m.isRead = true })
    },

    connectSSE() {
      if (this.abortCtrl) return

      const token = localStorage.getItem('vue-naivue-admin_auth')
      let accessToken = ''
      try {
        accessToken = JSON.parse(token)?.accessToken || ''
      } catch { /* ignore */ }

      if (!accessToken) return

      const baseURL = import.meta.env.VITE_AXIOS_BASE_URL || '/api'
      const url = `${baseURL}/message/sse`
      const abortCtrl = new AbortController()

      const connect = async () => {
        try {
          const res = await fetch(url, {
            headers: { Authorization: `Bearer ${accessToken}` },
            signal: abortCtrl.signal,
          })

          const reader = res.body.getReader()
          const decoder = new TextDecoder()
          let buffer = ''

          while (true) {
            const { done, value } = await reader.read()
            if (done) break

            buffer += decoder.decode(value, { stream: true })
            const lines = buffer.split('\n')
            buffer = lines.pop()

            let eventType = ''
            let eventData = ''

            for (const line of lines) {
              if (line.startsWith('event: ')) {
                eventType = line.slice(7)
              } else if (line.startsWith('data: ')) {
                eventData = line.slice(6)
              } else if (line === '' && eventType && eventData) {
                this._handleSSEEvent(eventType, eventData)
                eventType = ''
                eventData = ''
              }
            }
          }
        } catch (err) {
          if (err.name !== 'AbortError') {
            // Reconnect after 5s
            setTimeout(() => {
              if (this.abortCtrl) connect()
            }, 5000)
          }
        }
      }

      this.abortCtrl = abortCtrl
      connect()
    },

    _handleSSEEvent(type, data) {
      try {
        const parsed = JSON.parse(data)
        if (type === 'init') {
          this.unreadCount = parsed.unreadCount ?? 0
        } else if (type === 'message') {
          this.unreadCount++
          window.$notification.info({
            title: '新消息',
            content: parsed.title || '您有一条新消息',
            duration: 3000,
          })
        }
      } catch { /* ignore */ }
    },

    disconnectSSE() {
      if (this.abortCtrl) {
        this.abortCtrl.abort()
        this.abortCtrl = null
      }
    },
  },
})
