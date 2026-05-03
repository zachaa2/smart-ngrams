//go:build js && wasm

package main

import (
	"encoding/json"
	"sort"
	"strconv"
	"syscall/js"

	"github.com/zachaa2/smart-ngrams/internal/ngrams"
)

type NgramEntry struct {
	Ngram string `json:"ngram"`
	Count int    `json:"count"`
}

type NgramStats struct {
	Total  int `json:"total"`
	Unique int `json:"unique"`
}

type NgramResult struct {
	Stats  map[string]NgramStats   `json:"stats"`
	Result map[string][]NgramEntry `json:"result"`
}

func sortNgramsResult(ngrams []NgramEntry) []NgramEntry {
	// Sort by count descending, then alphabetically
	sort.Slice(ngrams, func(i, j int) bool {
		if ngrams[i].Count != ngrams[j].Count {
			return ngrams[i].Count > ngrams[j].Count
		}
		return ngrams[i].Ngram < ngrams[j].Ngram
	})
	return ngrams
}

// builds the result for a particular n-gram size
func buildResult(result NgramResult, counts map[string]int, n int) NgramResult {
	entries := make([]NgramEntry, 0, len(counts))
	for ngram, count := range counts {
		entries = append(entries, NgramEntry{Ngram: ngram, Count: count})
	}
	// sort
	entries = sortNgramsResult(entries)

	// add to result struct
	result.Result[strconv.Itoa(n)] = entries
	result.Stats[strconv.Itoa(n)] = NgramStats{
		Total:  ngrams.GetNgramTotalCount(counts),
		Unique: len(counts),
	}
	return result
}

func computeNgrams(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return js.ValueOf("ERROR: Expected two args: text and ns")
	}
	text := args[0].String()
	nsRaw := args[1].String()
	var ns []int
	if err := json.Unmarshal([]byte(nsRaw), &ns); err != nil {
		return js.ValueOf("ERROR: Failed to parse ns")
	}
	// one tokenization pass
	stopWordCounts := make(map[string]int)
	segments := ngrams.Tokenize(text, stopWordCounts)

	// init result struct
	result := NgramResult{
		Stats:  make(map[string]NgramStats),
		Result: make(map[string][]NgramEntry),
	}

	// compute ngrams on tokenized text, for each n
	for _, n := range ns {
		if n < 1 { // we'll just no-op for n < 1
			continue
		} else if n == 1 {
			wordCounts := make(map[string]int)
			ngrams.CountWords(segments, stopWordCounts, wordCounts)
			result = buildResult(result, wordCounts, 1)

		} else if n >= 2 {
			nGramCounts := make(map[string]int)
			ngrams.CountNGrams(segments, nGramCounts, n)
			result = buildResult(result, nGramCounts, n)
		}
	}

	// return
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return js.ValueOf("ERROR: failed to marshal results")
	}
	return js.ValueOf(string(jsonBytes))
}

func main() {
	js.Global().Set("computeNgrams", js.FuncOf(computeNgrams))
	select {} // keep Go runtime alive
}
