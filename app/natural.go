package app

import "strings"

// Natural language processor

func preprocessWord(word string) string {
	word = strings.TrimSpace(word)
	// TODO stem word, remove stop words and escape
	return word
}
