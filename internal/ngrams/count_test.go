package ngrams

import "testing"

// --- CountWords tests ---

func TestCountWords_CountsWordsFromSegments(t *testing.T) {
	segments := [][]string{{"hello", "world"}, {"hello"}}
	stopCounts := map[string]int{}
	wordCounts := map[string]int{}

	CountWords(segments, stopCounts, wordCounts)

	if wordCounts["hello"] != 2 {
		t.Errorf("expected 'hello' count 2, got %d", wordCounts["hello"])
	}
	if wordCounts["world"] != 1 {
		t.Errorf("expected 'world' count 1, got %d", wordCounts["world"])
	}
}

func TestCountWords_IncludesStopWordCounts(t *testing.T) {
	segments := [][]string{{"quick", "brown"}}
	stopCounts := map[string]int{"the": 3, "and": 1}
	wordCounts := map[string]int{}

	CountWords(segments, stopCounts, wordCounts)

	if wordCounts["the"] != 3 {
		t.Errorf("expected 'the' count 3, got %d", wordCounts["the"])
	}
	if wordCounts["and"] != 1 {
		t.Errorf("expected 'and' count 1, got %d", wordCounts["and"])
	}
}

func TestCountWords_SingleLetterStopWordsExcluded(t *testing.T) {
	// single-letter stop words (e.g. "a", "i") should not be included in word counts
	segments := [][]string{}
	stopCounts := map[string]int{"a": 5, "i": 2, "the": 1}
	wordCounts := map[string]int{}

	CountWords(segments, stopCounts, wordCounts)

	if _, ok := wordCounts["a"]; ok {
		t.Errorf("single-letter stop word 'a' should not be in word counts")
	}
	if _, ok := wordCounts["i"]; ok {
		t.Errorf("single-letter stop word 'i' should not be in word counts")
	}
	if wordCounts["the"] != 1 {
		t.Errorf("expected 'the' count 1, got %d", wordCounts["the"])
	}
}

func TestCountWords_AccumulatesIntoExistingMap(t *testing.T) {
	// counts should accumulate if wordCounts already has entries
	segments := [][]string{{"hello"}}
	stopCounts := map[string]int{}
	wordCounts := map[string]int{"hello": 10}

	CountWords(segments, stopCounts, wordCounts)

	if wordCounts["hello"] != 11 {
		t.Errorf("expected 'hello' count 11, got %d", wordCounts["hello"])
	}
}

func TestCountWords_EmptyInput(t *testing.T) {
	wordCounts := map[string]int{}
	CountWords([][]string{}, map[string]int{}, wordCounts)
	if len(wordCounts) != 0 {
		t.Errorf("expected empty word counts, got %v", wordCounts)
	}
}

// --- CountNGrams tests ---

func TestCountNGrams_Bigrams(t *testing.T) {
	segments := [][]string{{"hello", "world", "foo"}}
	nGramCounts := map[string]int{}

	CountNGrams(segments, nGramCounts, 2)

	if nGramCounts["hello world"] != 1 {
		t.Errorf("expected 'hello world' count 1, got %d", nGramCounts["hello world"])
	}
	if nGramCounts["world foo"] != 1 {
		t.Errorf("expected 'world foo' count 1, got %d", nGramCounts["world foo"])
	}
}

func TestCountNGrams_Trigrams(t *testing.T) {
	segments := [][]string{{"quick", "brown", "fox"}}
	nGramCounts := map[string]int{}

	CountNGrams(segments, nGramCounts, 3)

	if nGramCounts["quick brown fox"] != 1 {
		t.Errorf("expected 'quick brown fox' count 1, got %d", nGramCounts["quick brown fox"])
	}
	if len(nGramCounts) != 1 {
		t.Errorf("expected exactly 1 trigram, got %d", len(nGramCounts))
	}
}

func TestCountNGrams_SegmentTooShortIsSkipped(t *testing.T) {
	// segments shorter than n should produce no ngrams
	segments := [][]string{{"hello"}}
	nGramCounts := map[string]int{}

	CountNGrams(segments, nGramCounts, 2)

	if len(nGramCounts) != 0 {
		t.Errorf("expected no ngrams, got %v", nGramCounts)
	}
}

func TestCountNGrams_NGramsDoNotSpanSegments(t *testing.T) {
	// "world" ends segment 1 and "foo" starts segment 2 — "world foo" should not appear
	segments := [][]string{{"hello", "world"}, {"foo", "bar"}}
	nGramCounts := map[string]int{}

	CountNGrams(segments, nGramCounts, 2)

	if nGramCounts["world foo"] != 0 {
		t.Errorf("ngram 'world foo' should not span segments")
	}
}

func TestCountNGrams_RepeatedNGramsAccumulate(t *testing.T) {
	segments := [][]string{{"hello", "world"}, {"hello", "world"}}
	nGramCounts := map[string]int{}

	CountNGrams(segments, nGramCounts, 2)

	if nGramCounts["hello world"] != 2 {
		t.Errorf("expected 'hello world' count 2, got %d", nGramCounts["hello world"])
	}
}

func TestCountNGrams_EmptySegments(t *testing.T) {
	nGramCounts := map[string]int{}
	CountNGrams([][]string{}, nGramCounts, 2)
	if len(nGramCounts) != 0 {
		t.Errorf("expected no ngrams, got %v", nGramCounts)
	}
}

// --- GetNgramTotalCount tests ---

func TestGetNgramTotalCount_SumsAllCounts(t *testing.T) {
	counts := map[string]int{"hello world": 3, "brown fox": 2, "quick brown": 1}
	got := GetNgramTotalCount(counts)
	if got != 6 {
		t.Errorf("expected 6, got %d", got)
	}
}

func TestGetNgramTotalCount_EmptyMap(t *testing.T) {
	got := GetNgramTotalCount(map[string]int{})
	if got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}
