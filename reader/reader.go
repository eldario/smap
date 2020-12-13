// Read and parse given text
package reader

import (
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

// textStructure internal structure of reader
type textStructure struct {
	itemsMutex    sync.Mutex
	mapper        mapper
	minWordLength int
	rgp           *regexp.Regexp
	excludeWords  map[string]bool
}

// mapper internal interface of mapperClass
type mapper interface {
	Insert([]string, uint32)
	Remove(string)
}

// New reader Constructor
func New(mapper mapper, minWordLength int) *textStructure {
	return &textStructure{
		mapper:        mapper,
		minWordLength: minWordLength,
		rgp:           regexp.MustCompile(`[^a-zA-Z\s.0-9]+`),
		excludeWords:  make(map[string]bool),
	}
}

// Read read and parse each line from the text
func (t *textStructure) Read(content string, paragraphNumber uint32) {
	line := t.rgp.ReplaceAllString(content, "")

	if line = strings.TrimSpace(line); line != "" {
		var resultWords []string

		for _, line := range strings.Split(line, ".") {
			t.parseLine(strings.ToLower(strings.TrimSpace(line)), &resultWords)
		}
		t.mapper.Insert(resultWords, paragraphNumber)
	}

}

// parseLine split a line to word
func (t *textStructure) parseLine(line string, resultWords *[]string) {
	words := strings.Split(line, " ")

	wordsCount := len(words)
	for index, word := range words {
		if index == 0 || index == wordsCount-1 {
			t.updateExcludeList(word)
			t.mapper.Remove(word)
			continue
		}

		if t.isWordValid(word) {
			*resultWords = append(*resultWords, word)
		}
	}
}

// updateExcludeList Update exclude list with words
func (t *textStructure) updateExcludeList(word string) {
	t.itemsMutex.Lock()
	defer t.itemsMutex.Unlock()

	if _, ok := t.excludeWords[word]; !ok {
		t.excludeWords[word] = true
	}
}

// isWordValid returns true if word is valid
func (t *textStructure) isWordValid(word string) bool {
	return utf8.RuneCountInString(word) > t.minWordLength && !t.isWordExcluded(word)
}

// isWordExcluded returns if given word in exclude list
func (t *textStructure) isWordExcluded(word string) bool {
	t.itemsMutex.Lock()
	defer t.itemsMutex.Unlock()

	_, ok := t.excludeWords[word]

	return ok
}
