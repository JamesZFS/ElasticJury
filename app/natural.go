package app

import "strings"

// Natural language processor

func preprocessWord(word string) string {
	word = strings.TrimSpace(word)
	// TODO stem word, remove stop tags and escape
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

func FilterStrs(strs []string, predicate func(str string) bool) []string {
	var res []string
	for _, str := range strs {
		if predicate(str) {
			res = append(res, str)
		}
	}
	return res
}

func NotWhiteSpace(str string) bool {
	if strings.TrimSpace(str) == "" {
		return false
	} else {
		return true
	}
}
