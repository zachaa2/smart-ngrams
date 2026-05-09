// --- schema definitions - matches the Go structs ---
/** @typedef {{ ngram: string, count: number }} NgramEntry */
/** @typedef {{ total: number, unique: number }} NgramStats */
/** @typedef {{ stats: Object.<string,NgramStats>, result: Object.<string,NgramEntry[]> }} NgramResult */
/** @typedef {{ n: number, topN: number }} LevelConfig */

const DEFAULT_LEVELS = [
    { n: 1, topN: 128 },
    { n: 2, topN: 64 },
    { n: 3, topN: 32 },
    { n: 4, topN: 16 },
    { n: 5, topN: 8 },
]

/**
 * Adds a new level row to the UI with the given n and topN values.
 * @param {number} n the n in n-grams
 * @param {number} topN the # of top n-grams to show (by count)
*/
function addLevelRow(n, topN) {
    const row = document.createElement("div");
    row.className = "level-row";
    row.innerHTML = `
        <label>n= <input type="number" class="n-input" value="${n}" min="1" /></label>
        <label>Top <input type="number" class="topn-input" value="${topN}" min="1" /></label>
        <button class="remove-btn outline secondary">Remove</button>
    `
    row.querySelector(".remove-btn").addEventListener("click", () => row.remove());
    document.getElementById("levels").appendChild(row);
}

// init rows
DEFAULT_LEVELS.forEach(level => addLevelRow(level.n, level.topN))
document.getElementById("add-level-btn").addEventListener("click", () => addLevelRow(1, 10));

// TODO: add wasm integration