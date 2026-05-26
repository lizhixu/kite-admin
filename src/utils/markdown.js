import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import { marked } from 'marked'

// Configure marked with highlight.js
marked.setOptions({
  breaks: true,
  gfm: true,
  highlight(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(code, { language: lang }).value
      } catch { /* fall through */ }
    }
    return hljs.highlightAuto(code).value
  },
})

/**
 * Render Markdown string to sanitized HTML.
 */
export function renderMarkdown(md) {
  if (!md) return ''
  const html = marked.parse(md)
  return DOMPurify.sanitize(html, {
    ADD_TAGS: ['img'],
    ADD_ATTR: ['target', 'rel'],
  })
}

/**
 * Strip Markdown syntax and return plain text (for previews/summaries).
 */
export function stripMarkdown(md) {
  if (!md) return ''
  return md
    .replace(/```[\s\S]*?```/g, '') // code blocks
    .replace(/`[^`]*`/g, '') // inline code
    .replace(/!\[[^\]]*\]\([^)]*\)/g, '') // images
    .replace(/\[[^\]]*\]\([^)]*\)/g, (m) => { // links -> text only
      return m.replace(/\[[^\]]*\]\(/, '').replace(/\)$/, '')
    })
    .replace(/#{1,6}\s*/g, '') // headings
    .replace(/(\*\*|__)(.*?)\1/g, '$2') // bold
    .replace(/(\*|_)(.*?)\1/g, '$2') // italic
    .replace(/~~(.*?)~~/g, '$1') // strikethrough
    .replace(/>\s*/g, '') // blockquotes
    .replace(/[-*+]\s*/g, '') // unordered lists
    .replace(/\d+\.\s*/g, '') // ordered lists
    .replace(/---/g, '') // hr
    .replace(/\n{2,}/g, ' ') // multiple newlines
    .replace(/\n/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}
