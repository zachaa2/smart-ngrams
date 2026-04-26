package main

import (
	"fmt"
	"os"

	"github.com/zachaa2/smart-ngrams/internal/ngrams"
)

func readFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot access %q\n", filename)
		return ""
	}
	return string(content)
}

func ProcessFile(filename string,
	wordCounts map[string]int,
	bigramCounts map[string]int,
	trigramCounts map[string]int,
	fourgramCounts map[string]int,
	fivegramCounts map[string]int) bool {
	// read file first
	fileContent := readFile(filename)
	if len(fileContent) == 0 {
		return false
	}

	// tokenize
	stopWordCounts := make(map[string]int)
	tokens := ngrams.Tokenize(fileContent, stopWordCounts)

	// aggregate counts
	ngrams.CountWords(tokens, stopWordCounts, wordCounts)
	ngrams.CountNGrams(tokens, bigramCounts, 2)
	ngrams.CountNGrams(tokens, trigramCounts, 3)
	ngrams.CountNGrams(tokens, fourgramCounts, 4)
	ngrams.CountNGrams(tokens, fivegramCounts, 5)

	return true
}
