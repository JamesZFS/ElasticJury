package natural

import (
	. "ElasticJury/app/common"
	"fmt"
	"io/ioutil"
	"strings"
)

// Natural language processor

type stringSet map[string]Void

var (
	stopWords = make(stringSet)
)

// Init stopwords set
func Initialize() {
	bytes, err := ioutil.ReadFile(StopWordsPath)
	if err != nil {
		println("Initializing natural failed.")
		panic(err)
	}
	for _, word := range strings.Split(string(bytes), "\n") {
		stopWords[word] = Void{}
	}
}

func PreprocessWord(word string) string {
	word = strings.TrimSpace(word)
	// TODO stem word ?
	if _, isStopWord := stopWords[word]; isStopWord {
		fmt.Println("Found stop word: ", word)
		return ""
	} else {
		return word
	}
}

func PreprocessWords(words []string) []string {
	var res []string
	for _, word := range words {
		x := PreprocessWord(word)
		if x != "" {
			res = append(res, x)
		}
	}
	return res
}

func ParseFullText(text string) (words []string, tags []string, laws []string, judges []string) {
	panic("Unimplemented!")
}
