//go:build js && wasm

package main

import (
	"encoding/json"
	"sort"
	"syscall/js"
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

}

func hello(this js.Value, args []js.Value) any {
	return js.ValueOf("Hello from Go WASM!")
}

func main() {
	js.Global().Set("goHello", js.FuncOf(hello))
	select {} // keep Go runtime alive
}
