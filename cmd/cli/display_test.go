package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout runs f and returns everything written to stdout during its execution.
func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// --- printHeader tests ---

func TestPrintHeader_Words(t *testing.T) {
	out := captureStdout(func() { printHeader(1, 10, 50) })
	if !strings.Contains(out, "words") {
		t.Errorf("expected 'words' in output, got %q", out)
	}
}

func TestPrintHeader_Bigrams(t *testing.T) {
	out := captureStdout(func() { printHeader(2, 10, 50) })
	if !strings.Contains(out, "bigrams") {
		t.Errorf("expected 'bigrams' in output, got %q", out)
	}
}

func TestPrintHeader_Trigrams(t *testing.T) {
	out := captureStdout(func() { printHeader(3, 10, 50) })
	if !strings.Contains(out, "trigrams") {
		t.Errorf("expected 'trigrams' in output, got %q", out)
	}
}

func TestPrintHeader_FourGrams(t *testing.T) {
	out := captureStdout(func() { printHeader(4, 10, 50) })
	if !strings.Contains(out, "4-grams") {
		t.Errorf("expected '4-grams' in output, got %q", out)
	}
}

func TestPrintHeader_FiveGrams(t *testing.T) {
	out := captureStdout(func() { printHeader(5, 10, 50) })
	if !strings.Contains(out, "5-grams") {
		t.Errorf("expected '5-grams' in output, got %q", out)
	}
}

func TestPrintHeader_InvalidNgram(t *testing.T) {
	out := captureStdout(func() { printHeader(9, 10, 50) })
	if !strings.Contains(out, "Invalid") {
		t.Errorf("expected invalid usage message, got %q", out)
	}
}

func TestPrintHeader_TopNCappedAtNgramCount(t *testing.T) {
	// ngramCount=3 < topN=10, so output should say "Top 3"
	out := captureStdout(func() { printHeader(2, 10, 3) })
	if !strings.Contains(out, "3") {
		t.Errorf("expected top count to be capped at 3, got %q", out)
	}
}

// --- displayCounts tests ---

func TestDisplayCounts_SortedByFrequencyDesc(t *testing.T) {
	counts := map[string]int{"apple": 1, "banana": 5, "cherry": 3}
	out := captureStdout(func() { displayCounts(counts, 10, 1) })
	bananaPos := strings.Index(out, "banana")
	cherryPos := strings.Index(out, "cherry")
	applePos := strings.Index(out, "apple")
	if !(bananaPos < cherryPos && cherryPos < applePos) {
		t.Errorf("expected banana > cherry > apple order, got:\n%s", out)
	}
}

func TestDisplayCounts_AlphabeticTiebreaker(t *testing.T) {
	counts := map[string]int{"zebra": 5, "apple": 5, "mango": 5}
	out := captureStdout(func() { displayCounts(counts, 10, 1) })
	applePos := strings.Index(out, "apple")
	mangoPos := strings.Index(out, "mango")
	zebraPos := strings.Index(out, "zebra")
	if !(applePos < mangoPos && mangoPos < zebraPos) {
		t.Errorf("expected alphabetical order for equal counts, got:\n%s", out)
	}
}

func TestDisplayCounts_TopNLimitsOutput(t *testing.T) {
	counts := map[string]int{"apple": 3, "banana": 2, "cherry": 1}
	out := captureStdout(func() { displayCounts(counts, 2, 1) })
	if strings.Contains(out, "cherry") {
		t.Errorf("expected 'cherry' to be excluded when topN=2, got:\n%s", out)
	}
}

func TestDisplayCounts_EmptyMap(t *testing.T) {
	out := captureStdout(func() { displayCounts(map[string]int{}, 10, 1) })
	// should output header with "Top 0" and no entries
	if !strings.Contains(out, "0") {
		t.Errorf("expected 'Top 0' for empty map, got %q", out)
	}
}

// --- printStats tests ---

func TestPrintStats_ShowsValidDocCount(t *testing.T) {
out := captureStdout(func() {
printStats(3, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
})
if !strings.Contains(out, "3") {
t.Errorf("expected valid doc count 3 in output, got %q", out)
}
}

func TestPrintStats_ShowsWordCounts(t *testing.T) {
wordCounts := map[string]int{"lion": 5, "mouse": 3}
out := captureStdout(func() {
printStats(1, wordCounts, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
})
if !strings.Contains(out, "Number of words: 8") {
t.Errorf("expected total word count 8 in output, got %q", out)
}
if !strings.Contains(out, "Number of unique words: 2") {
t.Errorf("expected unique word count 2 in output, got %q", out)
}
}

func TestPrintStats_ShowsNgramCounts(t *testing.T) {
bigramCounts := map[string]int{"quick brown": 2, "brown fox": 1}
out := captureStdout(func() {
printStats(1, map[string]int{}, bigramCounts, map[string]int{}, map[string]int{}, map[string]int{})
})
if !strings.Contains(out, "bigrams: 3") {
t.Errorf("expected total bigram count 3 in output, got %q", out)
}
if !strings.Contains(out, "bigrams: 2") {
t.Errorf("expected unique bigram count 2 in output, got %q", out)
}
}

func TestPrintStats_ZeroCountsForEmptyMaps(t *testing.T) {
out := captureStdout(func() {
printStats(0, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{})
})
if !strings.Contains(out, "Number of words: 0") {
t.Errorf("expected zero word count, got %q", out)
}
}
