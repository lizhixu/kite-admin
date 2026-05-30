import { defineStore } from 'pinia'
import api from '@/views/message/api'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    unreadCount: 0,
    inboxMessages: [],
    inboxTotal: 0,
    inboxPage: 1,
    inboxPageSize: 15,
    abortCtrl: null,
    showInbox: false,
    detailMessage: null,
  }),

  getters: {
    recentMessages: state => state.inboxMessages.slice(0, 8),
  },

  actions: {
    async fetchUnreadCount() {
      const { data } = await api.getUnreadCount()
      this.unreadCount = data?.count ?? 0
    },

    async fetchInbox(params = {}) {
      const pageNo = params.pageNo ?? this.inboxPage
      const pageSize = params.pageSize ?? this.inboxPageSize
      const { data } = await api.getMyMessages({ pageNo, pageSize })
      this.inboxMessages = data?.pageData ?? []
      this.inboxTotal = data?.total ?? 0
      this.inboxPage = pageNo
      this.inboxPageSize = pageSize
      await this.fetchUnreadCount()
    },

    async markAsRead(id) {
      await api.markRead(id)
      this.inboxMessages.forEach((msg) => {
        if (msg.id === id && !msg.isRead) {
          msg.isRead = true
          msg.readAt = new Date().toISOString()
        }
      })
      if (this.detailMessage?.id === id && !this.detailMessage.isRead) {
        this.detailMessage.isRead = true
        this.detailMessage.readAt = new Date().toISOString()
      }
      await this.fetchUnreadCount()
    },

    async markAllAsRead() {
      await api.markAllRead()
      this.unreadCount = 0
      this.inboxMessages.forEach((msg) => {
        msg.isRead = true
        msg.readAt ||= new Date().toISOString()
      })
      if (this.detailMessage && !this.detailMessage.isRead) {
        this.detailMessage.isRead = true
        this.detailMessage.readAt ||= new Date().toISOString()
      }
    },

    async openMessage(msg) {
      this.detailMessage = msg
      this.showInbox = true
      if (!msg.isRead)
        await this.markAsRead(msg.id)
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
            setTimeout(() => {
              if (this.abortCtrl) connect()
            }, 5000)
          }
        }
      }

      this.abortCtrl = abortCtrl
      connect()
    },

    async _handleSSEEvent(type, data) {
      try {
        const parsed = JSON.parse(data)
        if (type === 'init') {
          this.unreadCount = parsed.unreadCount ?? 0
        } else if (type === 'message') {
          this.fetchUnreadCount()
          if (this.inboxMessages.length)
            this.fetchInbox({ pageNo: this.inboxPage, pageSize: this.inboxPageSize })
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
