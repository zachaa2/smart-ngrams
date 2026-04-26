package main

import (
	"fmt"
	"os"
)

// main is the CLI entry point for local testing of the smart-ngrams library.
// It accepts one or more text file paths as arguments and computes unigrams
// through 5-grams, printing summary statistics and ranked frequency lists.
// Usage: cli <file1> <file2> ...
func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v <file1> <file2> ...\n", os.Args[0])
		os.Exit(1)
	}

	// map of n-grams and their counts
	wordCounts := make(map[string]int)
	bigramCounts := make(map[string]int)
	trigramCounts := make(map[string]int)
	fourgramCounts := make(map[string]int)
	fivegramCounts := make(map[string]int)

	// process each file
	validDocuments := 0
	for _, filename := range args {
		result := processFile(filename, wordCounts, bigramCounts, trigramCounts, fourgramCounts, fivegramCounts)
		if result {
			validDocuments++
		}
	}

	printStats(validDocuments, wordCounts, bigramCounts, trigramCounts, fourgramCounts, fivegramCounts)

	const TOP_WORDS_TO_DISPLAY int = 128
	const TOP_BIGRAMS_TO_DISPLAY int = 64
	const TOP_TRIGRAMS_TO_DISPLAY int = 32
	const TOP_FOURGRAMS_TO_DISPLAY int = 16
	const TOP_FIVEGRAMS_TO_DISPLAY int = 8

	displayCounts(wordCounts, TOP_WORDS_TO_DISPLAY, 1)
	displayCounts(bigramCounts, TOP_BIGRAMS_TO_DISPLAY, 2)
	displayCounts(trigramCounts, TOP_TRIGRAMS_TO_DISPLAY, 3)
	displayCounts(fourgramCounts, TOP_FOURGRAMS_TO_DISPLAY, 4)
	displayCounts(fivegramCounts, TOP_FIVEGRAMS_TO_DISPLAY, 5)
}
