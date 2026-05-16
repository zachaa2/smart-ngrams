# smart-ngrams

A Go port of an NLP tool for computing "smart" n-grams from text. N-grams are segmented using stop words, so they never span across semantically weak words. The Go library compiles to WebAssembly (WASM) to power a simple browser-based SPA.

**Live demo:** [zachaa2.github.io/smart-ngrams](https://zachaa2.github.io/smart-ngrams/)

## Overview

Given a block of text, the tool:
1. **Tokenizes** the text into segments delimited by stop words and non-alphabetic characters
2. **Counts** n-grams of any size within each segment
3. **Ranks** results by frequency (with alphabetical tiebreaking)

Stop words are based on the AP89 top 50 most frequent English words. N-grams that would span a stop word boundary are intentionally excluded — hence "smart" ngrams.

The browser UI lets you paste text, configure which n-gram levels to compute (any n ≥ 1) and how many top results to display per level, then view ranked frequency tables.

## Repository Structure

```
smart-ngrams/
├── go.mod                          # Go module definition
├── ngrams.cpp                      # Original C++ implementation (reference)
│
├── cmd/
│   └── cli/                        # CLI entry point for local testing (package main)
│       ├── main.go                 # Arg parsing and orchestration
│       ├── file.go                 # File reading and per-file processing
│       ├── display.go              # Stdout formatting (stats, ranked counts)
│       ├── file_test.go            # Unit tests for file.go
│       ├── display_test.go         # Unit tests for display.go
│       └── data/                   # Sample text files for local testing
│
├── internal/
│   └── ngrams/                     # Core NLP library (package ngrams)
│       ├── stopwords.go            # AP89 stop word list and lookup helper
│       ├── tokenize.go             # Segmentation and tokenization logic
│       ├── tokenize_test.go        # Unit tests for tokenize.go
│       ├── count.go                # N-gram counting logic
│       └── count_test.go           # Unit tests for count.go
│
├── wasm/
│   └── main.go                     # Go WASM entry point — exports computeNgrams to JS
│
├── frontend/                       # Vite vanilla JS SPA
│   ├── index.html                  # Analyzer page
│   ├── algorithm.html              # Algorithm explainer (in progress)
│   ├── src/
│   │   ├── main.js                 # WASM loading, UI logic, result rendering
│   │   └── style.css
│   └── public/
│       ├── main.wasm               # Compiled Go WASM binary (see build instructions)
│       └── wasm_exec.js            # Go JS bridge (from Go stdlib)
│
└── .github/
    └── workflows/
        └── deploy.yml              # Builds WASM + frontend, deploys to GitHub Pages
```

## Running Unit Tests

Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./internal/ngrams/
```

Run with verbose output (shows each test name):
```bash
go test ./internal/ngrams/ -v
```

Run a specific test or group of tests by name pattern:
```bash
go test ./internal/ngrams/ -v -run TestSegmentText
go test ./internal/ngrams/ -v -run TestTokenize
go test ./internal/ngrams/ -v -run TestHandleQuotes
go test ./internal/ngrams/ -v -run TestCountWords
go test ./internal/ngrams/ -v -run TestCountNGrams
go test ./cmd/cli/ -v -run TestProcessFile
go test ./cmd/cli/ -v -run TestDisplayCounts
go test ./cmd/cli/ -v -run TestPrintStats
```

Show test coverage:
```bash
go test ./... -cover
```

## Building the WASM Binary

The `main.wasm` binary is not committed to the repo. Build it from the project root before running the frontend:

```bash
GOOS=js GOARCH=wasm go build -o frontend/public/main.wasm ./wasm/
```

If you need to refresh `wasm_exec.js` (e.g. after a Go version upgrade):
```bash
cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/
```

## Running the Frontend

```bash
cd frontend
npm install   # first time only
npm run dev   # starts Vite dev server
```

## Running the CLI

```bash
go run ./cmd/cli/ <file1> <file2> ...
```

Example using the included sample files:
```bash
go run ./cmd/cli/ cmd/cli/data/lion.txt
go run ./cmd/cli/ cmd/cli/data/lion.txt cmd/cli/data/kira.txt
```

## LLM Disclosure

This project is mainly meant to be a learning exercise in Go and WebAssembly, where I put my tutorial knowledge to the test by porting some cpp code I wrote in college and making a small project out of it. The boundary between self-written and LLM-assisted code roughly follows the boundary between the algorithm and the infrastructure around it.

**Written independently (no LLM assistance):**
- The entire core NLP library (`internal/ngrams/`) — tokenizer, stop word handling, n-gram counting, and all unit tests
- The CLI wrapper (`cmd/cli/`)
- The original C++ implementation this was ported from

**Written with LLM guidance (GitHub Copilot):**
- `wasm/main.go` — the WASM layer was implemented by me with Copilot providing explanations of the `syscall/js` API, Go/JS interop patterns, and reviewing the implementation. Key esign decisions and actual coding is done by me.
- `frontend/` — HTML structure, CSS, and JS were written by me with Copilot advising on Vite multi-page setup, Pico CSS conventions, and vanilla JS DOM patterns. The core rendering and WASM wiring logic was written by me.
- `.github/workflows/deploy.yml` — generated by Copilot and lightly adjusted
- Like 90% of this `readme`

The intent was to use the LLM as a knowledgeable collaborator for unfamiliar tooling (Go, WASM, Vite) while keeping the algorithmic work fully self-directed.

