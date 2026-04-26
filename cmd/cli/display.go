package main

import (
	"fmt"
	"sort"

	"github.com/zachaa2/smart-ngrams/internal/ngrams"
)

type entry struct {
	key   string
	count int
}

// printHeader prints the section header for an ngram frequency block.
// ngram selects the label (1=words, 2=bigrams, ..., 5=5-grams), topN and
// ngramCount together determine the displayed count (capped at ngramCount).
func printHeader(ngram, topN, ngramCount int) {
	var outstr string
	switch ngram {
	case 1:
		outstr = " words:\n"
	case 2:
		outstr = " interesting bigrams:\n"
	case 3:
		outstr = " interesting trigrams:\n"
	case 4:
		outstr = " interesting 4-grams:\n"
	case 5:
		outstr = " interesting 5-grams:\n"
	default:
		fmt.Println("Invalid usage of ngram arg in print_header()")
		return
	}
	fmt.Printf("\nTop %v %v", min(topN, ngramCount), outstr)
}

// displayCounts prints the top-N entries from counts, sorted by frequency
// descending with alphabetical ordering as a tiebreaker. ngram controls the
// section header label printed above the results.
func displayCounts(counts map[string]int, topN int, ngrams int) {
	// need to convert to a slice of key, count pairs to enfore sorting
	entries := make([]entry, 0, len(counts))
	for k, v := range counts {
		entries = append(entries, entry{k, v})
	}

	// sort by freq desc. then by alpha as tiebraker
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].key < entries[j].key
	})

	// output top n items or less if there are less than n items
	printHeader(ngrams, topN, len(entries))
	for i := 0; i < min(topN, len(entries)); i++ {
		fmt.Printf("%v %v\n", entries[i].count, entries[i].key)
	}
}

// printStats prints a summary of n-gram counts across all processed documents,
// including total and unique counts for words and each n-gram level (2–5).
func printStats(validDocs int,
	wordCounts map[string]int,
	bigramCounts map[string]int,
	trigramCounts map[string]int,
	fourgramCounts map[string]int,
	fivegramCounts map[string]int) {

	totWords := ngrams.GetNgramTotalCount(wordCounts)
	uniqueWords := len(wordCounts)

	totBigrams := ngrams.GetNgramTotalCount(bigramCounts)
	uniqueBigrams := len(bigramCounts)

	totTrigrams := ngrams.GetNgramTotalCount(trigramCounts)
	uniqueTrigrams := len(trigramCounts)

	totFourgrams := ngrams.GetNgramTotalCount(fourgramCounts)
	uniqueFourgrams := len(fourgramCounts)

	totFivegrams := ngrams.GetNgramTotalCount(fivegramCounts)
	uniqueFivegrams := len(fivegramCounts)

	fmt.Printf("Number of valid documents: %v\n", validDocs)
	fmt.Printf("Number of words: %v\n", totWords)
	fmt.Printf("Number of unique words: %v\n", uniqueWords)
	fmt.Printf("Number of \"interesting\" bigrams: %v\n", totBigrams)
	fmt.Printf("Number of unique \"interesting\" bigrams: %v\n", uniqueBigrams)
	fmt.Printf("Number of \"interesting\" trigrams: %v\n", totTrigrams)
	fmt.Printf("Number of unique \"interesting\" trigrams: %v\n", uniqueTrigrams)
	fmt.Printf("Number of \"interesting\" 4-grams: %v\n", totFourgrams)
	fmt.Printf("Number of unique \"interesting\" 4-grams: %v\n", uniqueFourgrams)
	fmt.Printf("Number of \"interesting\" 5-grams: %v\n", totFivegrams)
	fmt.Printf("Number of unique \"interesting\" 5-grams: %v\n", uniqueFivegrams)

}
