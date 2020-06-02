package natural

import (
	. "ElasticJury/app/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanyiwu/gojieba"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// Natural language processor

type stringSet map[string]Void
type stringMap map[string]float32
type dictMap map[string]Dict

var (
	stopWords 	= make(stringSet)
	dicts		= make(dictMap)
	jieba     	*gojieba.Jieba
)

// Init stopwords set
func Initialize() {
	// Stopwords
	bytes, err := ioutil.ReadFile(StopWordsPath)
	if err != nil {
		println("[Info] Initializing natural failed.")
		panic(err)
	}
	for _, word := range strings.Split(string(bytes), "\n") {
		stopWords[word] = Voidance
	}

	// Jieba
	jieba = gojieba.NewJieba()

	// Dicts
	dicts["tag"] = BuildDict(TagDictPath)
	dicts["judge"] = BuildDict(JudgeDictPath)
	dicts["law"] = BuildDict(LawsDictPath)
	dicts["word"] = JoinDicts(Dicts{dicts["tag"], dicts["judge"], dicts["law"]})

	// Finish
	println("[Info] Natural initialized.")
	return
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

// Parse misc text into words
func ParseFullText(text string) Conditions {
	all := jieba.ExtractWithWeight(text, SearchWordLimit)
	words := make([]string, len(all))
	weights := make([]float32, len(all))
	var mean float32
	for i := range all {
		words[i] = all[i].Word
		weights[i] = float32(all[i].Weight)
		mean += weights[i]
	}
	if len(all) > 0 {
		mean /= float32(len(all))
	} else {
		mean = 1.0
	}
	for i := range all {
		weights[i] += mean
	}
	reduced := Reduce(words, weights)
	sort.Sort(reduced)
	return reduced
}

func MakeAssociateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		field := context.Param("field")
		item := context.Param("item")
		fmt.Printf("[Associate] Got request:\n")
		fmt.Printf("[Associate] > Field: %s\n", field)
		fmt.Printf("[Associate] > Item: %s\n", item)

		dict, ok := dicts[field]
		if len(item) == 0 || !ok {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		items := dict[item]
		context.JSON(http.StatusOK, gin.H{
			"count": len(items),
			"data": items,
		})
		fmt.Printf("[Associate] Reply with %d items\n", len(items))
	}
}
