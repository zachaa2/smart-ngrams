package ngrams

// stopWords is the AP89 top 50 stop words. N-grams cannot span across these words.
var stopWords = map[string]struct{}{
	"the": {}, "of": {}, "to": {}, "a": {}, "and": {}, "in": {}, "said": {}, "for": {},
	"that": {}, "was": {}, "on": {}, "he": {}, "is": {}, "with": {}, "at": {}, "by": {},
	"it": {}, "from": {}, "as": {}, "be": {}, "were": {}, "an": {}, "have": {}, "his": {},
	"but": {}, "has": {}, "are": {}, "not": {}, "who": {}, "they": {}, "its": {}, "had": {},
	"will": {}, "would": {}, "about": {}, "i": {}, "been": {}, "this": {}, "their": {},
	"new": {}, "or": {}, "which": {}, "we": {}, "more": {}, "after": {}, "us": {},
	"percent": {}, "up": {}, "one": {}, "people": {},
}

func isStopWord(w string) bool {
	_, ok := stopWords[w]
	return ok
}
