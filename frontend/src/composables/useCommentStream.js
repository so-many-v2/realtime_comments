import { ref, shallowRef, onUnmounted } from 'vue'
import { commentStreamUrl } from '../api.js'

/**
 * SSE-поток комментариев с батчевой отрисовкой.
 *
 * Идея: входящие сообщения при высоком RPS нельзя пушить в реактивный массив по
 * одному — каждый push триггерит реактивность и перерисовку, и на 2k rps это
 * кладёт UI. Поэтому:
 *   1. Сообщения копятся в обычном (не реактивном) буфере `buffer`.
 *   2. Раз в кадр (requestAnimationFrame) буфер сливается в реактивный список
 *      одним изменением → одна перерисовка на кадр, а не на каждое сообщение.
 *   3. И буфер, и отрисованный список ограничены сверху (maxRendered): то, что
 *      всё равно уедет за пределы окна, не доходит до DOM. Это держит и память,
 *      и число DOM-узлов под контролем независимо от RPS.
 *
 * @param {string|number} postId
 * @param {object} opts
 * @param {number} opts.maxRendered  сколько комментариев держим в DOM (окно)
 * @param {number} opts.maxPerFrame  сколько максимум вливаем за один кадр
 */
export function useCommentStream(postId, { maxRendered = 200, maxPerFrame = 120 } = {}) {
  const comments = shallowRef([]) // отрисованное окно (новые снизу)
  const totalReceived = ref(0) // счётчик всех принятых (включая вытесненные)
  const connected = ref(false)

  const seen = new Set() // дедуп по id (REST-префилл + SSE могут пересечься)
  let buffer = [] // не реактивный «накопитель»
  let rafId = null
  let es = null
  let stopped = false

  function pushSeen(id) {
    if (id == null) return true
    if (seen.has(id)) return false
    seen.add(id)
    return true
  }

  function enqueue(comment) {
    if (!pushSeen(comment?.id)) return
    buffer.push(comment)
    totalReceived.value++

    // Буфер не растёт бесконечно: всё, что заведомо вытеснится из окна,
    // даже не доходит до отрисовки.
    const overflow = buffer.length - maxRendered
    if (overflow > 0) buffer.splice(0, overflow)

    scheduleFlush()
  }

  function scheduleFlush() {
    if (rafId != null || stopped) return
    rafId = requestAnimationFrame(flush)
  }

  function flush() {
    rafId = null
    if (buffer.length === 0) return

    const batch = buffer.splice(0, maxPerFrame)
    let next = comments.value.concat(batch)

    // обрезаем окно: оставляем только последние maxRendered
    if (next.length > maxRendered) {
      next = next.slice(next.length - maxRendered)
    }
    comments.value = next // одно реактивное изменение на кадр

    if (buffer.length > 0) scheduleFlush() // остаток дольём в следующем кадре
  }

  // Префилл историей из REST: сразу в окно, без анимаций.
  function seed(list) {
    if (!Array.isArray(list)) return
    const fresh = list.filter((c) => pushSeen(c?.id))
    let next = fresh
    if (next.length > maxRendered) next = next.slice(next.length - maxRendered)
    comments.value = next
  }

  function connect() {
    es = new EventSource(commentStreamUrl(postId))
    es.onopen = () => {
      connected.value = true
    }
    es.onerror = () => {
      // EventSource переподключается сам; просто отражаем статус
      connected.value = false
    }
    es.onmessage = (e) => {
      try {
        enqueue(JSON.parse(e.data))
      } catch (_) {
        // битый кадр игнорируем
      }
    }
  }

  function close() {
    stopped = true
    if (es) es.close()
    if (rafId != null) cancelAnimationFrame(rafId)
    rafId = null
  }

  onUnmounted(close)

  return { comments, totalReceived, connected, seed, connect, close }
}