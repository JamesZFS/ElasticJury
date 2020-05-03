package app

import "strings"

// Natural language processor

func preprocessWord(word string) string {
	word = strings.TrimSpace(word)
	// TODO stem word, remove stop words and escape
	return word
}

func PreprocessWords(words []string) []string {
	var res []string
	for _, word := range words {
		x := preprocessWord(word)
		if x != "" {
			res = append(res, x)
		}
	}
	return res
}
