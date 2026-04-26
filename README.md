# smart-ngrams

A Go port of an NLP tool for computing "smart" n-grams from text. N-grams are segmented using stop words, so they never span across semantically weak words. The project will eventually compile to WebAssembly (WASM) to power a simple browser-based SPA.

## Overview

Given one or more text documents, the tool:
1. **Tokenizes** the text into segments delimited by stop words and non-alphabetic characters
2. **Counts** unigrams through 5-grams within each segment
3. **Ranks** results by frequency (with alphabetical tiebreaking)

Stop words are based on the AP89 top 50 most frequent English words. N-grams that would span a stop word boundary are intentionally excluded — hence "smart" ngrams.

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

## Running the CLI

```bash
go run ./cmd/cli/ <file1> <file2> ...
```

Example using the included sample files:
```bash
go run ./cmd/cli/ cmd/cli/data/lion.txt
go run ./cmd/cli/ cmd/cli/data/lion.txt cmd/cli/data/kira.txt
```
