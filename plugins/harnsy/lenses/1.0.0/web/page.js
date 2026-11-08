// The page of Lenses: the project's lenses, a recording to analyze, the
// analyses. Every text from Lua or the person goes in by textContent or
// the kit's properties, never into innerHTML. Whether a recording needs
// ffmpeg is harnsy's answer (ok, code of the recordings): drawn, never
// decided here. Whether harnsy transcribes at all now (no connection, the
// month's limit) is transcribe.ready(), asked once when the page opens: one
// notice, its why behind ⓘ.

const STYLE = `
  :host { display: block; }
  .stack { display: grid; gap: 16px; }
  .row { display: flex; align-items: center; gap: 8px; padding: 8px 0; border-bottom: 1px solid var(--border); }
  .row:last-child { border-bottom: 0; }
  .row .grow { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .row.pick { cursor: pointer; }
  .row.pick.on { font-weight: 600; }
  .row.off { color: var(--text-muted); cursor: default; }
  .muted { color: var(--text-muted); font-size: 13px; }
  .help { border: 1px solid var(--border); border-radius: 50%; width: 20px; height: 20px; padding: 0;
    background: transparent; color: var(--text-muted); cursor: pointer; font-size: 12px; line-height: 18px; }
  .note { color: var(--text-muted); font-size: 13px; margin: 4px 0 8px; }
  .bar { display: flex; gap: 8px; align-items: center; margin-top: 12px; flex-wrap: wrap; }
  .error { color: var(--danger, #c33); font-size: 13px; }
  .notice { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
  .notice[hidden] { display: none; }
`

const el = (tag, cls, text) => {
  const e = document.createElement(tag)
  if (cls) e.className = cls
  if (text !== undefined) e.textContent = text
  return e
}

export async function mount(root, harnsy) {
  const t = (k, v) => harnsy.t(k, v)
  root.innerHTML = `<style>${STYLE}</style>
    <div class="stack">
      <div class="notice" id="notice" hidden></div>
      <hk-card id="lenses"><div id="lens-rows"></div><hk-button slot="actions" id="lens-add"></hk-button></hk-card>
      <hk-card id="analyze">
        <div id="recs"></div>
        <div class="bar"><hk-chips id="chips"></hk-chips></div>
        <div class="bar"><hk-button variant="primary" id="go"></hk-button><span class="error" id="go-error"></span></div>
      </hk-card>
      <hk-card id="analyses"><div id="rows"></div></hk-card>
    </div>
    <hk-dialog id="dialog">
      <hk-field id="f-name" type="text" required></hk-field>
      <hk-field id="f-instruction" type="textarea" required></hk-field>
    </hk-dialog>`
  const $ = (id) => root.getElementById(id)
  $('lenses').setAttribute('heading', t('lenses'))
  $('analyze').setAttribute('heading', t('analyze.title'))
  $('analyses').setAttribute('heading', t('analyses'))
  $('lens-add').textContent = t('lens.add')
  $('go').textContent = t('analyze.go')
  $('chips').setAttribute('label', t('analyze.lenses'))
  const name = $('f-name'), instruction = $('f-instruction'), dialog = $('dialog')
  name.setAttribute('label', t('lens.name'))
  name.setAttribute('maxlength', '60')
  instruction.setAttribute('label', t('lens.instruction'))
  instruction.setAttribute('maxlength', '2000')
  instruction.setAttribute('rows', '6')
  dialog.setAttribute('confirm', t('save'))
  dialog.setAttribute('cancel', t('cancel'))

  let lenses = [], picked = null, editing = null, asked = false

  const why = (e) => (e && (e.code || e.message)) || String(e)

  const drawLenses = () => {
    const box = $('lens-rows')
    box.replaceChildren()
    if (lenses.length === 0) box.append(el('div', 'muted', t('lens.empty')))
    for (const l of lenses) {
      const row = el('div', 'row')
      row.append(el('span', 'grow', l.name))
      const edit = el('hk-button')
      edit.textContent = t('lens.edit')
      edit.addEventListener('press', () => openLens(l))
      row.append(edit)
      box.append(row)
    }
    const chips = $('chips')
    const chosen = new Set(chips.value || [])
    chips.items = lenses.map((l) => ({ value: l.id, label: l.name }))
    chips.value = lenses.filter((l) => chosen.size === 0 || chosen.has(l.id)).map((l) => l.id)
  }

  const openLens = (l) => {
    editing = l || null
    asked = false
    dialog.setAttribute('heading', t(l ? 'lens.edit' : 'lens.new'))
    if (l) dialog.setAttribute('danger', t('lens.delete'))
    else dialog.removeAttribute('danger')
    dialog.error = ''
    name.value = l ? l.name : ''
    instruction.value = l ? l.instruction : ''
    name.error = instruction.error = ''
    dialog.open = true
  }

  const loadLenses = async () => {
    lenses = (await harnsy.call('lens_list', [harnsy.lang])) || []
    drawLenses()
  }

  dialog.addEventListener('confirm', async () => {
    dialog.busy = true
    try {
      const r = await harnsy.call('lens_save', [editing ? editing.id : 0, name.value || '', instruction.value || ''])
      if (r && r.error) {
        name.error = r.error === 'name' ? t('err.name') : ''
        instruction.error = r.error === 'instruction' ? t('err.instruction') : ''
        return
      }
      dialog.open = false
      await loadLenses()
    } catch (e) {
      dialog.error = t('failed', { why: why(e) })
    } finally {
      dialog.busy = false
    }
  })
  dialog.addEventListener('cancel', () => { dialog.open = false })
  // Delete asks once in the window itself: the second press deletes.
  dialog.addEventListener('danger', async () => {
    if (!editing) return
    if (!asked) {
      asked = true
      dialog.error = t('lens.delete_ask', { name: editing.name })
      return
    }
    try {
      await harnsy.call('lens_delete', [editing.id])
      dialog.open = false
      await loadLenses()
    } catch (e) {
      dialog.error = t('failed', { why: why(e) })
    }
  })
  $('lens-add').addEventListener('press', () => openLens(null))

  // A recording: its name and item, and harnsy's answer whether it goes
  // through now. «Needs ffmpeg» with its «?»: the server's code alone.
  const drawRecordings = (recs) => {
    const box = $('recs')
    box.replaceChildren()
    if (recs.length === 0) {
      box.append(el('div', 'muted', t('analyze.none')))
      return
    }
    for (const r of recs) {
      const usable = !!r.transcript || r.ok
      const row = el('div', 'row ' + (usable ? 'pick' : 'off'))
      if (picked && picked.attachment === r.attachment) row.classList.add('on')
      row.append(el('span', 'grow', r.name))
      row.append(el('span', 'muted', '#' + r.item))
      if (r.transcript) {
        const b = el('hk-badge', '', t('analyze.ready'))
        b.setAttribute('tone', 'ok')
        row.append(b)
      } else if (!r.ok) {
        const b = el('hk-badge', '', r.code === 'transcribe_no_ffmpeg' ? t('analyze.need_ffmpeg') : t('analyze.cannot'))
        b.setAttribute('tone', 'wait')
        b.title = r.code || ''
        row.append(b)
        if (r.code === 'transcribe_no_ffmpeg') {
          const help = el('button', 'help', '?')
          help.type = 'button'
          help.title = t('analyze.ffmpeg_help')
          help.setAttribute('aria-label', t('analyze.ffmpeg_help'))
          help.addEventListener('click', (e) => {
            e.stopPropagation()
            const next = row.nextSibling
            if (next && next.classList && next.classList.contains('note')) next.remove()
            else row.after(el('div', 'note', t('analyze.ffmpeg_help')))
          })
          row.append(help)
        }
      }
      if (usable) {
        row.addEventListener('click', () => {
          picked = r
          drawRecordings(recs)
        })
      }
      box.append(row)
    }
  }

  const loadRecordings = async () => drawRecordings((await harnsy.call('recordings')) || [])

  $('go').addEventListener('press', async () => {
    const ids = $('chips').value || []
    $('go-error').textContent = ''
    if (!picked || ids.length === 0) {
      $('go-error').textContent = t('err.pick')
      return
    }
    $('go').busy = true
    try {
      const r = await harnsy.call('analyze', [picked.item, picked.attachment, ids, harnsy.lang])
      if (r && r.error) {
        const known = { pick: 'err.pick', long: 'err.long', gone: 'err.gone' }[r.error]
        $('go-error').textContent = known ? t(known) : t('failed', { why: r.error })
        return
      }
      picked = null
      await Promise.all([loadRecordings(), loadAnalyses()])
    } catch (e) {
      $('go-error').textContent = t('failed', { why: why(e) })
    } finally {
      $('go').busy = false
    }
  })

  // An analysis: recording · lenses · state · the item, once there is one.
  const stateOf = (a) => {
    switch (a.state) {
      case 'transcribing':
        return a.total > 0 ? [t('state.transcribing', { done: a.done, total: a.total }), 'info'] : [t('state.queued'), 'muted']
      case 'agent': return [t('state.agent'), 'info']
      case 'done': return [t('state.done'), 'ok']
      case 'failed': return [t('state.failed') + ' ⓘ', 'bad']
      default: return [t('state.queued'), 'muted']
    }
  }

  const loadAnalyses = async () => {
    const rows = (await harnsy.call('analyses')) || []
    const box = $('rows')
    box.replaceChildren()
    if (rows.length === 0) box.append(el('div', 'muted', t('analyses.empty')))
    for (const a of rows) {
      const row = el('div', 'row')
      row.append(el('span', 'grow', a.name))
      row.append(el('span', 'muted', t('col.lenses') + ': ' + a.lenses))
      const [text, tone] = stateOf(a)
      const b = el('hk-badge', '', text)
      b.setAttribute('tone', tone)
      if (a.state === 'failed') b.title = a.code || ''
      row.append(b)
      const open = el('hk-button')
      open.textContent = '#' + (a.task || a.item)
      open.addEventListener('press', () => harnsy.open('item:' + (a.task || a.item)))
      row.append(open)
      box.append(row)
    }
  }

  // One notice when harnsy starts no transcription now: a short word, the
  // why behind ⓘ; the code is harnsy's, never decided here.
  const drawReady = (r) => {
    const box = $('notice')
    box.replaceChildren()
    if (!r || r.ok) { box.hidden = true; return }
    const known = { transcribe_no_connection: 'ready.no_connection', transcribe_budget: 'ready.budget' }[r.code]
    const b = el('hk-badge', '', t(known || 'ready.off'))
    b.setAttribute('tone', 'wait')
    b.title = r.code || ''
    const helpText = known ? t(known + '_help') : (r.code || '')
    const help = el('button', 'help', 'ⓘ')
    help.type = 'button'
    help.title = helpText
    help.setAttribute('aria-label', helpText)
    help.addEventListener('click', () => {
      const next = box.querySelector('.note')
      if (next) next.remove()
      else box.append(el('div', 'note', helpText))
    })
    box.append(b, help)
    box.hidden = false
  }
  const loadReady = async () => drawReady(await harnsy.call('node_ready'))

  const loadAll = async () => {
    try {
      await Promise.all([loadReady(), loadLenses(), loadRecordings(), loadAnalyses()])
    } catch (e) {
      $('go-error').textContent = t('failed', { why: why(e) })
    }
  }
  harnsy.on('refresh', () => { loadRecordings().catch(() => {}); loadAnalyses().catch(() => {}) })
  await loadAll()
}

export function unmount() {}
