package main

import "testing"

// --- readFile tests ---

func TestReadFile_ValidFile(t *testing.T) {
	content := readFile("data/lion.txt")
	if len(content) == 0 {
		t.Error("expected non-empty content for lion.txt")
	}
}

func TestReadFile_EmptyFile(t *testing.T) {
	content := readFile("data/empty.txt")
	if len(content) != 0 {
		t.Errorf("expected empty content for empty.txt, got %d bytes", len(content))
	}
}

func TestReadFile_MissingFile(t *testing.T) {
	content := readFile("data/does_not_exist.txt")
	if content != "" {
		t.Errorf("expected empty string for missing file, got %q", content)
	}
}

// --- ProcessFile tests ---

func TestProcessFile_ReturnsFalseForMissingFile(t *testing.T) {
	result := ProcessFile("data/does_not_exist.txt",
		map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
	if result {
		t.Error("expected false for missing file")
	}
}

func TestProcessFile_ReturnsFalseForEmptyFile(t *testing.T) {
	result := ProcessFile("data/empty.txt",
		map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
	if result {
		t.Error("expected false for empty file")
	}
}

func TestProcessFile_ReturnsTrueForValidFile(t *testing.T) {
	result := ProcessFile("data/lion.txt",
		map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
	if !result {
		t.Error("expected true for valid file lion.txt")
	}
}

func TestProcessFile_PopulatesWordCounts(t *testing.T) {
	wordCounts := map[string]int{}
	ProcessFile("data/lion.txt", wordCounts,
		map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
	if len(wordCounts) == 0 {
		t.Error("expected word counts to be populated")
	}
}

func TestProcessFile_PopulatesBigramCounts(t *testing.T) {
	bigramCounts := map[string]int{}
	ProcessFile("data/lion.txt", map[string]int{}, bigramCounts,
		map[string]int{}, map[string]int{}, map[string]int{})
	if len(bigramCounts) == 0 {
		t.Error("expected bigram counts to be populated")
	}
}

func TestProcessFile_NgramsDoNotSpanStopWords(t *testing.T) {
	// "the" is a stop word — any bigram spanning it should not exist
	bigramCounts := map[string]int{}
	ProcessFile("data/lion.txt", map[string]int{}, bigramCounts,
		map[string]int{}, map[string]int{}, map[string]int{})
	for bigram := range bigramCounts {
		for _, stop := range []string{" the ", " of ", " and ", " in "} {
			if len(bigram) > len(stop) {
				// check if the bigram contains a stop word in between its two words
				// a valid bigram "word1 word2" won't have interior stop words
				_ = stop
			}
		}
	}
	// specifically: "lion the" or "the mouse" should not appear as bigrams
	if bigramCounts["lion the"] > 0 {
		t.Errorf("bigram 'lion the' should not exist — spans a stop word")
	}
	if bigramCounts["the mouse"] > 0 {
		t.Errorf("bigram 'the mouse' should not exist — spans a stop word")
	}
}

func TestProcessFile_AccumulatesAcrossMultipleCalls(t *testing.T) {
	wordCounts := map[string]int{}
	bigramCounts := map[string]int{}
	trigramCounts := map[string]int{}
	fourgramCounts := map[string]int{}
	fivegramCounts := map[string]int{}

	ProcessFile("data/lion.txt", wordCounts, bigramCounts, trigramCounts, fourgramCounts, fivegramCounts)
	firstTotal := len(wordCounts)

	ProcessFile("data/kira.txt", wordCounts, bigramCounts, trigramCounts, fourgramCounts, fivegramCounts)
	secondTotal := len(wordCounts)

	if secondTotal <= firstTotal {
		t.Errorf("expected word count to grow after processing second file, got %d -> %d", firstTotal, secondTotal)
	}
}
