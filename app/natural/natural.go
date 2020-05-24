package natural

import (
	. "ElasticJury/app/common"
	"encoding/json"
	"github.com/yanyiwu/gojieba"
	"io/ioutil"
	"sort"
	"strings"
)

// Natural language processor

type stringSet map[string]Void
type stringMap map[string]float32

const (
	useHmm = true
)

var (
	stopWords = make(stringSet)
	idfDict   = make(stringMap)
	jieba     *gojieba.Jieba
)

// Init stopwords set
func Initialize() {
	// Stopwords
	bytes, err := ioutil.ReadFile(StopWordsPath)
	if err != nil {
		goto ERROR
	}
	for _, word := range strings.Split(string(bytes), "\n") {
		stopWords[word] = Voidance
	}

	// Jieba
	jieba = gojieba.NewJieba()

	// Idf dictionary
	bytes, err = ioutil.ReadFile(IdfDictPath)
	if err != nil {
		goto ERROR
	}
	err = json.Unmarshal(bytes, &idfDict)
	if err != nil {
		goto ERROR
	}
	println("Natural initialized.")
	return

ERROR:
	println("Initializing natural failed.")
	panic(err)
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
		return ""
	} else {
		return word
	}
}

// Trim white space, escape, and filter out empty words
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

func Reduce(words []string, weights []float32) Conditions {
	m := stringMap{}
	for i := range words {
		m[words[i]] = m[words[i]] + weights[i]
	}
	var reduced Conditions
	for k, v := range m {
		reduced = append(reduced, Condition{Item: k, Weight: v})
	}
	return reduced
}

func GetWordsWeights(words []string) []float32 {
	weights := make([]float32, len(words))
	var mean float32
	inDictCount := 0
	for i, word := range words {
		v, in := idfDict[word]
		if in {
			mean += v
			inDictCount++
			weights[i] = v
		}
	}
	if inDictCount == 0 {
		mean = 1.0
	} else {
		mean /= float32(inDictCount) * 2
	}
	for i, word := range words {
		_, in := idfDict[word]
		if !in {
			weights[i] = mean
		}
	}
	return weights
}

// Parse misc text into words
func ParseFullText(text string) Conditions {
	words := PreprocessWords(jieba.CutForSearch(text, useHmm))
	weights := GetWordsWeights(words)
	reduced := Reduce(words, weights)
	sort.Sort(reduced)
	return reduced[:Min(len(reduced), SearchWordLimit)]
}
