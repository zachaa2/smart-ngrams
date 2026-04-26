package ngrams

// CountWords tallies word (unigram) frequencies into wordCounts from the tokenized
// segments and stopWordCounts. Stop words are included in the tally unless they are
// a single letter. Single-letter words are excluded entirely.
func CountWords(segments [][]string, stopWordCounts map[string]int, wordCounts map[string]int) {
	// add stop word tallies
	for word, counts := range stopWordCounts {
		if len(word) > 1 {
			wordCounts[word] += counts
		}
	}

	// tally add words from tokenized text
	for _, segment := range segments {
		for _, word := range segment {
			wordCounts[word]++
		}
	}
}

// CountNGrams counts all n-grams of length n within each segment of the tokenized
// text, accumulating frequencies into nGramCounts. N-grams do not span segment
// boundaries. Segments shorter than n are skipped.
func CountNGrams(segments [][]string, nGramCounts map[string]int, n int) {
	for _, segment := range segments {
		if len(segment) < n {
			continue
		}

		// accumulate all ngrams for tokenized text segments
		for i := 0; i <= len(segment)-n; i++ {
			var nGram string = segment[i]
			for j := 1; j < n; j++ {
				nGram += " " + segment[i+j]
			}
			nGramCounts[nGram]++
		}
	}
}

// GetNgramTotalCount returns the total number of n-gram occurrences in counts,
// summing all frequency values in the map.
func GetNgramTotalCount(counts map[string]int) int {
	totCount := 0
	for _, count := range counts {
		totCount += count
	}
	return totCount
}
