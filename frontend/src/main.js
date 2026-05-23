// --- schema definitions - matches the Go structs ---
/** @typedef {{ ngram: string, count: number }} NgramEntry */
/** @typedef {{ total: number, unique: number }} NgramStats */
/** @typedef {{ stats: Object.<string,NgramStats>, result: Object.<string,NgramEntry[]> }} NgramResult */
/** @typedef {{ n: number, topN: number }} LevelConfig */

// DEFAULT CONFIG
const DEFAULT_LEVELS = [
    { n: 1, topN: 32 },
    { n: 2, topN: 16 },
    { n: 3, topN: 8 },
]

const ACCEPTED_EXTENSIONS = ['.txt', '.md', '.log']
const DEFAULT_FEEDBACK = `Accepted file types: ${ACCEPTED_EXTENSIONS.join(', ')}`

/**
 * Adds a new level row to the UI with the given n and topN values.
 * @param {number} n the n in n-grams
 * @param {number} topN the # of top n-grams to show (by count)
 */
function addLevelRow(n, topN) {
    const row = document.createElement('div')
    row.className = 'level-row'
    row.innerHTML = `
        <label>n= <input type="number" class="n-input" value="${n}" min="1" /></label>
        <label>Top <input type="number" class="topn-input" value="${topN}" min="1" /></label>
        <button class="remove-btn outline secondary">Remove</button>
    `
    row.querySelector('.remove-btn').addEventListener('click', () =>
        row.remove()
    )
    document.getElementById('levels').appendChild(row)
}

async function initWasm() {
    const go = new Go()
    const result = await WebAssembly.instantiateStreaming(
        fetch(import.meta.env.BASE_URL + 'main.wasm'),
        go.importObject
    )
    go.run(result.instance)
    // enable run btn
    document.getElementById('run-btn').disabled = false
}

/** @returns LevelConfig[] */
function collectLevelRows() {
    return Array.from(document.querySelectorAll('.level-row')).map((row) => ({
        n: parseInt(row.querySelector('.n-input').value),
        topN: parseInt(row.querySelector('.topn-input').value),
    }))
}

/**
 * Render results table for specific ngram
 * @param {number} n
 * @param {number} topN
 * @param {NgramEntry[]} entries
 * @param {NgramStats} stats
 * @returns {HTMLElement}
 */
function renderResultTable(n, topN, entries, stats) {
    const slicedData = entries.slice(0, topN)
    const label = n === 1 ? 'Words' : `${n}-grams`

    const details = document.createElement('details')
    details.open = false
    details.innerHTML = `
        <summary>${label} - total: ${stats.total}, unique: ${stats.unique}</summary>
        <table>
            <thead class="result-table-head">
                <tr><th>#</th><th>n-gram</th><th>count</th></tr>
            </thead>
            <tbody>
                ${slicedData
                    .map((entry, i) => {
                        return `<tr><td>${i + 1}</td><td>${entry.ngram}</td><td>${entry.count}</td></tr>`
                    })
                    .join('')}
            </tbody>
        </table>
    `
    return details
}

/**
 * Render all results tables for all ngram levels
 * @param {NgramResult} result
 * @param {LevelConfig} levels
 */
function renderResults(result, levels) {
    const container = document.getElementById('results')
    container.innerHTML = '' //clear old results

    levels.forEach(({ n, topN }) => {
        const entries = result.result[String(n)] ?? []
        const stats = result.stats[String(n)] ?? { total: 0, unique: 0 }
        container.appendChild(renderResultTable(n, topN, entries, stats))
    })
}

// dropzone logic
const dropzone = document.getElementById('dropzone')
const fileInput = dropzone.querySelector('.dropzone-input')
const dropzoneFeedback = document.getElementById('dropzone-feedback')

// validator for filetype
function isAcceptableTextFile(file) {
    if (file.type.startsWith('text/')) return true
    const name = file.name.toLowerCase()
    return ACCEPTED_EXTENSIONS.some((ext) => name.endsWith(ext))
}

// feedback msg setter
/**
 *
 * @param {string} message dropzone msg
 * @param {string} state 'info' | 'success' | 'error'
 */
function setFeedback(message, state = 'info') {
    dropzoneFeedback.textContent = message
    dropzoneFeedback.classList.remove('error', 'success')
    if (state !== 'info') dropzoneFeedback.classList.add(state)
}
// set default msg
setFeedback(DEFAULT_FEEDBACK)

// file loading logic
async function loadFile(file) {
    if (!file) return
    if (!isAcceptableTextFile(file)) {
        setFeedback(`Unsupported file type: ${file.name}`, 'error')
        return
    }
    const text = await file.text()
    document.getElementById('input-textarea').value = text
    setFeedback(`File ${file.name} read successfully`, 'success')
}

// event listeners for dropzone/fileinput
dropzone.addEventListener('click', () => fileInput.click())

fileInput.addEventListener('change', async (e) => {
    await loadFile(e.target.files[0])
    e.target.value = ''
})

dropzone.addEventListener('dragenter', (e) => {
    e.preventDefault()
    dropzone.classList.add('dropzone-dragging')
})
dropzone.addEventListener('dragover', (e) => {
    e.preventDefault()
})
dropzone.addEventListener('dragleave', () => {
    dropzone.classList.remove('dropzone-dragging')
})
dropzone.addEventListener('drop', async (e) => {
    e.preventDefault()
    dropzone.classList.remove('dropzone-dragging')
    await loadFile(e.dataTransfer.files[0])
})

// input selector logic
const inputSelector = document.getElementById('input-selector')
const textarea = document.getElementById('input-textarea')

inputSelector.addEventListener('change', (e) => {
    if (e.target.value === 'text') {
        dropzone.hidden = true
        textarea.hidden = false
    } else {
        dropzone.hidden = false
        textarea.hidden = true
    }
})

// event listener for the run btn
document.getElementById('run-btn').addEventListener('click', () => {
    const text = document.getElementById('input-textarea').value
    const levels = collectLevelRows()
    const ns = levels.map((l) => l.n)

    /** @type {NgramResult} */
    const result = JSON.parse(computeNgrams(text, JSON.stringify(ns)))
    renderResults(result, levels)
})

// INIT
initWasm()

// rows
DEFAULT_LEVELS.forEach((level) => addLevelRow(level.n, level.topN))
document
    .getElementById('add-level-btn')
    .addEventListener('click', () => addLevelRow(1, 10))
