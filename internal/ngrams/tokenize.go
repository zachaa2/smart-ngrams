package ngrams

import (
	"strings"
	"unicode"
)

// pushSegmentIfNotEmpty appends a copy of currentSegment to allSegments if it is non-empty,
// then clears currentSegment. Returns the updated slices.
func pushSegmentIfNotEmpty(allSegments [][]string, currentSegment []string) ([][]string, []string) {
	if len(currentSegment) > 0 {
		// need to append a copy of current segment
		allSegments = append(allSegments, append([]string{}, currentSegment...))
		currentSegment = currentSegment[:0]
	}
	return allSegments, currentSegment
}

// handleQuotes further parses segmented text by handling internal single quotes.
// Conjunctions (a single quote within a word, e.g. "don't") are kept intact.
// Multiple consecutive single quotes act as a word delimiter, splitting the word into two.
// It is assumed that leading and trailing single quotes have already been stripped from
// each word by SegmentText before this function is called.
func handleQuotes(segments [][]string) [][]string {
	updatedSegments := [][]string{}

	for _, segment := range segments {
		currentSegment := []string{}
		for _, word := range segment {
			var currentWord string
			var quoteCtr int

			for _, c := range word {
				if c == '\'' {
					quoteCtr++
					// split words on multiple single quote
					if quoteCtr > 1 {
						if len(currentWord) > 0 {
							// strip the trailing quote accumulated when quoteCtr was 1
							if currentWord[len(currentWord)-1] == '\'' {
								currentWord = currentWord[:len(currentWord)-1]
							}
							if len(currentWord) > 1 {
								currentSegment = append(currentSegment, currentWord)
							} else {
								// end current segment if word len is 1
								updatedSegments, currentSegment = pushSegmentIfNotEmpty(updatedSegments, currentSegment)
							}
							currentWord = ""
						}
						// reset and continue
						quoteCtr = 0
						continue
					}
				}
				currentWord += string(c)
			}
			// end of a word
			if len(currentWord) > 0 {
				if len(currentWord) > 1 {
					currentSegment = append(currentSegment, currentWord)
				} else {
					updatedSegments, currentSegment = pushSegmentIfNotEmpty(updatedSegments, currentSegment)
				}
			}
		}
		// add any remaining segments
		updatedSegments, currentSegment = pushSegmentIfNotEmpty(updatedSegments, currentSegment)
	}
	return updatedSegments
}

// processWord strips leading and trailing single quotes from word and reports
// whether the result is a stop word.
func processWord(word string) (cleaned string, isStop bool) {
	// trim leading and trailing quotes
	cleaned = strings.Trim(word, "'")
	_, isStop = stopWords[cleaned]
	return
}

// SegmentText splits text into segments, where each segment is a slice of words
// delimited by stop words or single-letter words. Non-alphabetic characters (except
// single quotes in conjunctions) are treated as word delimiters. All words are
// lowercased. Stop word frequencies are tallied into stopWordCounts.
// N-grams should not span across segment boundaries.
func SegmentText(text string, stopWordCounts map[string]int) [][]string {
	var segments [][]string
	var currentSegment []string
	var word string

	// main loop
	for _, c := range text {
		if unicode.IsLetter(c) || string(c) == "'" {
			// build words from alpha or single quite chars
			word += string(unicode.ToLower(c))
		} else { // case where a non-alpha or single quote char is encountered
			if len(word) > 0 {
				// handle word
				isStop := false
				word, isStop = processWord(word)
				isStopWordOrSingleLetter := isStop || len(word) == 1
				if isStopWordOrSingleLetter {
					if isStop {
						stopWordCounts[word]++
					}
					// terminate current segment
					segments, currentSegment = pushSegmentIfNotEmpty(segments, currentSegment)
				} else {
					currentSegment = append(currentSegment, word)
				}
				word = ""
			}
		}
	}

	// add any remaining word/segment
	if len(word) > 0 {
		// handle word ... again
		isStop := false
		word, isStop = processWord(word)
		isStopWordOrSingleLetter := isStop || len(word) == 1
		if isStopWordOrSingleLetter {
			if isStop {
				stopWordCounts[word]++
			}
			// terminate current segment
			segments, currentSegment = pushSegmentIfNotEmpty(segments, currentSegment)
		} else {
			currentSegment = append(currentSegment, word)
		}
		word = ""
	}

	segments, currentSegment = pushSegmentIfNotEmpty(segments, currentSegment)
	return segments
}

// Tokenize fully tokenizes text into segments ready for n-gram calculation.
// It runs SegmentText to produce initial segments, then handleQuotes to further
// parse internal single quotes within words.
func Tokenize(text string, stopWordCounts map[string]int) [][]string {
	segments := SegmentText(text, stopWordCounts)
	tokens := handleQuotes(segments)
	return tokens
}
