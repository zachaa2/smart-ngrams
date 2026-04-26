package ngrams

import (
	"reflect"
	"testing"
)

// --- handleQuotes tests ---

func TestHandleQuotes_NoQuotes(t *testing.T) {
	input := [][]string{{"hello", "world"}}
	got := handleQuotes(input)
	want := [][]string{{"hello", "world"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestHandleQuotes_Conjunction(t *testing.T) {
	// a single quote within a word (contraction) should be kept as one word
	input := [][]string{{"don't", "stop"}}
	got := handleQuotes(input)
	want := [][]string{{"don't", "stop"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestHandleQuotes_MultipleQuotesSplitsWord(t *testing.T) {
	// two single quotes in a row should split the word into two
	input := [][]string{{"hello''world"}}
	got := handleQuotes(input)
	want := [][]string{{"hello", "world"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestHandleQuotes_SingleLetterWordEndsSegment(t *testing.T) {
	// a single-letter result from a split should terminate the current segment
	input := [][]string{{"a''world"}}
	got := handleQuotes(input)
	// "a" is single-letter so it ends the segment; "world" starts a new one
	want := [][]string{{"world"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestHandleQuotes_EmptyInput(t *testing.T) {
	got := handleQuotes([][]string{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestHandleQuotes_MultipleSegments(t *testing.T) {
	// each input segment should be handled independently
	input := [][]string{{"don't", "stop"}, {"it''works"}}
	got := handleQuotes(input)
	want := [][]string{{"don't", "stop"}, {"it", "works"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// --- pushSegmentIfNotEmpty tests ---

func TestPushSegmentIfNotEmpty(t *testing.T) {
	tests := []struct {
		name           string     // test name
		allSegments    [][]string //first call arg
		currentSegment []string   // second call arg
		wantSegments   [][]string // first return arg
		wantCurrent    []string   // second return arg
	}{
		{
			name:           "empty current segment is not appended",
			allSegments:    [][]string{},
			currentSegment: []string{},
			wantSegments:   [][]string{},
			wantCurrent:    []string{},
		},
		{
			name:           "non-empty current segment is appended and cleared",
			allSegments:    [][]string{},
			currentSegment: []string{"hello", "world"},
			wantSegments:   [][]string{{"hello", "world"}},
			wantCurrent:    []string{},
		},
		{
			name:           "appends to existing segments",
			allSegments:    [][]string{{"foo", "bar"}},
			currentSegment: []string{"baz"},
			wantSegments:   [][]string{{"foo", "bar"}, {"baz"}},
			wantCurrent:    []string{},
		},
		{
			name:           "appended segment is an independent copy",
			allSegments:    [][]string{},
			currentSegment: []string{"hello", "world"},
			wantSegments:   [][]string{{"hello", "world"}},
			wantCurrent:    []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSegments, gotCurrent := pushSegmentIfNotEmpty(tc.allSegments, tc.currentSegment)

			if !reflect.DeepEqual(gotSegments, tc.wantSegments) {
				t.Errorf("allSegments: got %v, want %v", gotSegments, tc.wantSegments)
			}
			if !reflect.DeepEqual(gotCurrent, tc.wantCurrent) {
				t.Errorf("currentSegment: got %v, want %v", gotCurrent, tc.wantCurrent)
			}
		})
	}
}

func TestPushSegmentIfNotEmpty_AppendedSegmentIsACopy(t *testing.T) {
	current := []string{"hello", "world"}
	allSegments, current := pushSegmentIfNotEmpty([][]string{}, current)

	// mutate current after pushing — stored segment should be unaffected
	current = append(current, "mutated")

	if reflect.DeepEqual(allSegments[0], current) {
		t.Errorf("stored segment was mutated along with currentSegment — not a real copy")
	}
}

// --- processWord tests ---

func TestProcessWord(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCleaned string
		wantIsStop  bool
	}{
		{"plain word", "hello", "hello", false},
		{"stop word", "the", "the", true},
		{"leading quote stripped", "'hello", "hello", false},
		{"trailing quote stripped", "hello'", "hello", false},
		{"both quotes stripped", "'hello'", "hello", false},
		{"stop word with quotes", "'the'", "the", true},
		{"empty string", "", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotCleaned, gotIsStop := processWord(tc.input)
			if gotCleaned != tc.wantCleaned {
				t.Errorf("cleaned: got %q, want %q", gotCleaned, tc.wantCleaned)
			}
			if gotIsStop != tc.wantIsStop {
				t.Errorf("isStop: got %v, want %v", gotIsStop, tc.wantIsStop)
			}
		})
	}
}

// --- SegmentText tests ---

func TestSegmentText_EmptyInput(t *testing.T) {
	stopCounts := map[string]int{}
	got := SegmentText("", stopCounts)
	if len(got) != 0 {
		t.Errorf("expected no segments, got %v", got)
	}
}

func TestSegmentText_SingleSegment(t *testing.T) {
	stopCounts := map[string]int{}
	got := SegmentText("quick brown fox", stopCounts)
	want := [][]string{{"quick", "brown", "fox"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_StopWordSplitsSegment(t *testing.T) {
	// "the" is a stop word and should terminate the current segment
	stopCounts := map[string]int{}
	got := SegmentText("quick brown the lazy fox", stopCounts)
	want := [][]string{{"quick", "brown"}, {"lazy", "fox"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_StopWordIsCounted(t *testing.T) {
	stopCounts := map[string]int{}
	SegmentText("quick the brown the fox", stopCounts)
	if stopCounts["the"] != 2 {
		t.Errorf("expected stop word 'the' count 2, got %d", stopCounts["the"])
	}
}

func TestSegmentText_StopWordNotInSegment(t *testing.T) {
	// stop words should not appear in any segment
	stopCounts := map[string]int{}
	segments := SegmentText("quick the brown fox", stopCounts)
	for _, seg := range segments {
		for _, w := range seg {
			if isStopWord(w) {
				t.Errorf("stop word %q found in segment", w)
			}
		}
	}
}

func TestSegmentText_SingleLetterWordSplitsSegment(t *testing.T) {
	// single-letter words act like stop words — they terminate the segment
	stopCounts := map[string]int{}
	got := SegmentText("quick x brown fox", stopCounts)
	want := [][]string{{"quick"}, {"brown", "fox"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_NonAlphaCharsAreDelimiters(t *testing.T) {
	// punctuation and numbers should act as word delimiters
	stopCounts := map[string]int{}
	got := SegmentText("quick,brown.fox", stopCounts)
	want := [][]string{{"quick", "brown", "fox"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_ContractionKeptIntact(t *testing.T) {
	// a single quote within a word (contraction) should be preserved
	stopCounts := map[string]int{}
	got := SegmentText("don't stop", stopCounts)
	want := [][]string{{"don't", "stop"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_LeadingTrailingQuotesStripped(t *testing.T) {
	stopCounts := map[string]int{}
	got := SegmentText("'hello' world", stopCounts)
	want := [][]string{{"hello", "world"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_TextIsLowercased(t *testing.T) {
	stopCounts := map[string]int{}
	got := SegmentText("Quick Brown Fox", stopCounts)
	want := [][]string{{"quick", "brown", "fox"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSegmentText_OnlyStopWords(t *testing.T) {
	// all stop words should produce no segments
	stopCounts := map[string]int{}
	got := SegmentText("the of to a and", stopCounts)
	if len(got) != 0 {
		t.Errorf("expected no segments, got %v", got)
	}
}

// --- Tokenize tests ---

func TestTokenize_BasicText(t *testing.T) {
	stopCounts := map[string]int{}
	got := Tokenize("the quick brown fox", stopCounts)
	want := [][]string{{"quick", "brown", "fox"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTokenize_ContractionSurvivesFullPipeline(t *testing.T) {
	// contractions should pass through both SegmentText and handleQuotes intact
	stopCounts := map[string]int{}
	got := Tokenize("don't stop believing", stopCounts)
	want := [][]string{{"don't", "stop", "believing"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTokenize_StopWordCountsPopulated(t *testing.T) {
	stopCounts := map[string]int{}
	Tokenize("the quick brown fox jumps over the lazy dog", stopCounts)
	if stopCounts["the"] != 2 {
		t.Errorf("expected 'the' count 2, got %d", stopCounts["the"])
	}
	if stopCounts["over"] != 0 {
		// "over" is not in the AP89 stop word list
		t.Errorf("'over' should not be a stop word, got count %d", stopCounts["over"])
	}
}

func TestTokenize_EmptyInput(t *testing.T) {
	stopCounts := map[string]int{}
	got := Tokenize("", stopCounts)
	if len(got) != 0 {
		t.Errorf("expected no tokens, got %v", got)
	}
}

func TestTokenize_MultipleSegmentsFromStopWords(t *testing.T) {
	stopCounts := map[string]int{}
	got := Tokenize("cats and dogs or birds", stopCounts)
	// "and" and "or" are stop words, producing three segments
	want := [][]string{{"cats"}, {"dogs"}, {"birds"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
