import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import { marked } from 'marked'

const CODE_BLOCK_RE = /```[\s\S]*?```/g
const INLINE_CODE_RE = /`[^`]*`/g
const IMAGE_RE = /!\[[^\]]*\]\([^)]*\)/g
const LINK_RE = /\[[^\]]*\]\([^)]*\)/g
const LINK_TEXT_PREFIX_RE = /\[[^\]]*\]\(/
const LINK_TEXT_SUFFIX_RE = /\)$/
const HEADING_RE = /#{1,6}\s*/g
const BOLD_RE = /(\*\*|__)(.*?)\1/g
const ITALIC_RE = /(\*|_)(.*?)\1/g
const STRIKETHROUGH_RE = /~~(.*?)~~/g
const BLOCKQUOTE_RE = />\s*/g
const UNORDERED_LIST_RE = /[-*+]\s*/g
const ORDERED_LIST_RE = /\d+\.\s*/g
const HR_RE = /---/g
const MULTIPLE_NEWLINE_RE = /\n{2,}/g
const NEWLINE_RE = /\n/g
const WHITESPACE_RE = /\s+/g

// Configure marked with highlight.js
marked.setOptions({
  breaks: true,
  gfm: true,
  highlight(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(code, { language: lang }).value
      }
      catch { /* fall through */ }
    }
    return hljs.highlightAuto(code).value
  },
})

/**
 * Render Markdown string to sanitized HTML.
 */
export function renderMarkdown(md) {
  if (!md)
    return ''
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
  if (!md)
    return ''
  return md
    .replace(CODE_BLOCK_RE, '') // code blocks
    .replace(INLINE_CODE_RE, '') // inline code
    .replace(IMAGE_RE, '') // images
    .replace(LINK_RE, (m) => { // links -> text only
      return m.replace(LINK_TEXT_PREFIX_RE, '').replace(LINK_TEXT_SUFFIX_RE, '')
    })
    .replace(HEADING_RE, '') // headings
    .replace(BOLD_RE, '$2') // bold
    .replace(ITALIC_RE, '$2') // italic
    .replace(STRIKETHROUGH_RE, '$1') // strikethrough
    .replace(BLOCKQUOTE_RE, '') // blockquotes
    .replace(UNORDERED_LIST_RE, '') // unordered lists
    .replace(ORDERED_LIST_RE, '') // ordered lists
    .replace(HR_RE, '') // hr
    .replace(MULTIPLE_NEWLINE_RE, ' ') // multiple newlines
    .replace(NEWLINE_RE, ' ')
    .replace(WHITESPACE_RE, ' ')
    .trim()
}
