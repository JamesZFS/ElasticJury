package natural

import (
	. "ElasticJury/app/common"
	"github.com/yanyiwu/gojieba"
	"io/ioutil"
	"strings"
)

// Natural language processor

type stringSet map[string]Void

const (
	useHmm = true
)

var (
	stopWords = make(stringSet)
	jieba     *gojieba.Jieba
)

// Init stopwords set
func Initialize() {
	bytes, err := ioutil.ReadFile(StopWordsPath)
	if err != nil {
		println("Initializing natural failed.")
		panic(err)
	}
	for _, word := range strings.Split(string(bytes), "\n") {
		stopWords[word] = Voidance
	}
	jieba = gojieba.NewJieba()
}

func Finalize() {
	jieba.Free()
}

func PreprocessWord(word string) string {
	word = strings.TrimSpace(word)
	escaped := make([]int32, 0, len(word))
	for _, c := range word { // escape
		switch c {
		case '\'', '"', '`', '\\':
			// ignore
		default:
			escaped = append(escaped, c)
		}
	}
	word = string(escaped)
	if _, isStopWord := stopWords[word]; isStopWord {
		//fmt.Println("Found stop word: ", word)
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

// Parse misc text and output it into the four fields
// * TODO: Need discussion
func ParseFullText(text string) (words []string) {
	return PreprocessWords(jieba.CutForSearch(text, useHmm))
}
