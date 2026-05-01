# agents.md — Project Context for Coding Agents

## What This Project Does

Computes "smart" n-grams from text. N-grams never span stop word boundaries — that's the "smart" part. The Go library handles all NLP; WASM exposes it to a browser frontend.

---

## Repository Layout

```
internal/ngrams/   — core NLP library (package ngrams), no main, no I/O
cmd/cli/           — CLI wrapper for local testing; mirrors what WASM will do
wasm/              — Go WASM entry point; exports functions to JS via syscall/js
frontend/          — Vite vanilla JS SPA (no frameworks)
```

---

## Core Library API (`internal/ngrams`)

These are the only functions the WASM and CLI layers call:

```go
Tokenize(text string, stopWordCounts map[string]int) [][]string
// Returns segments (slices of word slices). Stop words delimit segments.
// stopWordCounts is populated as a side effect.

CountWords(segments [][]string, stopWordCounts map[string]int, wordCounts map[string]int)
// Tallies unigrams. Stop words (len > 1) are included. Single-letter words excluded.

CountNGrams(segments [][]string, nGramCounts map[string]int, n int)
// Counts n-grams within each segment. N-grams never cross segment boundaries.
// Segments shorter than n are skipped — produces empty map if n > all segment lengths.

GetNgramTotalCount(counts map[string]int) int
// Sums all values in a count map.
```

Caller pattern (matches both CLI and WASM):
```go
stopWordCounts := make(map[string]int)
tokens := ngrams.Tokenize(text, stopWordCounts)
ngrams.CountWords(tokens, stopWordCounts, wordCounts)    // n=1
ngrams.CountNGrams(tokens, nGramCounts, n)               // n≥2, repeat per n
```

---

## WASM Layer (`wasm/`)

- Build tag: `//go:build js && wasm`
- Exports Go functions to JS via `js.Global().Set("name", js.FuncOf(fn))`
- Function signature for all exported fns: `func(this js.Value, args []js.Value) any`
- `select {}` in `main()` keeps the Go runtime alive
- **Target function**: `computeNgrams(text string, nsJSON string) string`
  - `nsJSON` is a JSON array of ints (e.g. `"[1,2,3,6]"`) — which n values to compute
  - Returns a JSON string with stats and sorted result entries keyed by n (as string)
  - Sort order: count descending, then alphabetically ascending (tiebreaker)
  - If n exceeds all segment lengths: empty array + zero stats for that n

Build command:
```bash
GOOS=js GOARCH=wasm go build -o frontend/public/main.wasm ./wasm/
```

---

## Frontend (`frontend/`)

- Vite + vanilla JS/HTML/CSS — no frameworks, no TypeScript
- `wasm_exec.js` is the Go JS bridge; do not modify it
- WASM is loaded via `WebAssembly.instantiateStreaming` in `main.js`
- After `go.run(instance)`, exported Go functions are available on `window`
- UI: textarea input → dynamic level config (n value + top-N per row) → results tables
- JS slices pre-sorted Go results to each level's top-N; no sorting in JS

---

## Coding Patterns

- **Sorting**: always count desc, alpha asc — matches the CLI's `displayCounts` logic
- **No I/O in `internal/ngrams`**: the library is pure computation, takes strings/maps
- **Maps are always pre-allocated by caller** and passed in; library mutates them
- **n=1 is words** (uses `CountWords`); **n≥2 uses `CountNGrams`** — these are separate functions
- **Tests live next to source**: `*_test.go` in the same package

---

## Module

```
module github.com/zachaa2/smart-ngrams
```
