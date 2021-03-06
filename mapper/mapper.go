// package for collect words
package mapper

import (
	"sort"
	"sync"
)

// sortedMap internal structure for class
type sortedMap struct {
	itemsMutex sync.Mutex
	words      map[string]*WordItem
	topCount   int
}

// WordItem structure of word object
type WordItem struct {
	Word   string
	Count  uint32
	line   int
	column int
}

// New Structure constructor
func New() *sortedMap {
	return &sortedMap{
		words: make(map[string]*WordItem),
	}
}

func (s *sortedMap) SetTopCount(value int) {
	s.topCount = value
}

// Insert a new word in words map
func (s *sortedMap) Insert(words []string, position int) {
	s.itemsMutex.Lock()
	defer s.itemsMutex.Unlock()

	for wordIndex, word := range words {
		if _, ok := s.words[word]; !ok {
			s.words[word] = &WordItem{
				word,
				0,
				position,
				wordIndex,
			}
		}
		currentWord := s.words[word]
		currentWord.Count++
		if currentWord.line > position {
			currentWord.line = position
			currentWord.column = wordIndex
		}

		s.words[word] = currentWord
	}

}

// Remove a word from words map
func (s *sortedMap) Remove(word string) {
	s.itemsMutex.Lock()
	defer s.itemsMutex.Unlock()

	if _, ok := s.words[word]; ok {
		delete(s.words, word)
	}
}

// GetResults Get frequently used words in text
func (s *sortedMap) GetResults() []WordItem {
	var sortedResult []WordItem

	for _, word := range s.words {
		sortedResult = append(sortedResult, *word)
	}

	sort.Slice(sortedResult, func(i, j int) bool {
		if sortedResult[i].Count == sortedResult[j].Count {
			return sortedResult[i].line < sortedResult[j].line
		}
		return sortedResult[i].Count > sortedResult[j].Count
	})

	if len(sortedResult) >= s.topCount {
		sortedResult = sortedResult[:s.topCount]
	}

	sort.Slice(sortedResult, func(i, j int) bool {
		if sortedResult[i].line == sortedResult[j].line {
			return sortedResult[i].column < sortedResult[j].column
		}
		return sortedResult[i].line < sortedResult[j].line
	})

	return sortedResult
}
