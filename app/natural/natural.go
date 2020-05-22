package natural

import (
	. "ElasticJury/app/common"
	"errors"
	"fmt"
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
		stopWords[word] = Void{}
	}
	jieba = gojieba.NewJieba()
}

func Finalize() {
	jieba.Free()
}

func PreprocessWord(word string) string {
	word = strings.TrimSpace(word)
	// TODO stem word ?
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

const (
	stateNormal = iota // scanning a normal character
	stateLaw           // scanning between '《' and '》'
)

// Parse misc text and output it into the four fields
// * Todo Need discussion
func ParseFullText(text string) (words []string, tags []string, laws []string, judges []string) {
	state := stateNormal
	normalBuffer := make([]int32, 0, len(text))
	lawBuffer := make([]int32, 0, 512)
	laws = make([]string, 0)
	// a simple state-machine parser
	for i, char := range text {
		switch state {
		case stateNormal:
			switch char {
			case '《':
				state = stateLaw
			case '》':
				fmt.Printf("Bad law '》' at %d\n", i)
				// ignore
			default:
				normalBuffer = append(normalBuffer, char)
			}
		case stateLaw:
			switch char {
			case '《':
				fmt.Printf("Bad law '《' at %d\n", i)
				lawBuffer = lawBuffer[:0] // clear with capacity kept, https://yourbasic.org/golang/clear-slice/
				// ignore
			case '》':
				// finishes a law
				if len(lawBuffer) > 0 {
					laws = append(laws, "《"+string(lawBuffer)+"》")
				}
				lawBuffer = lawBuffer[:0] // clear with capacity kept
				state = stateNormal
				normalBuffer = append(normalBuffer, ' ')
			default:
				lawBuffer = append(lawBuffer, char)
			}
		default:
			panic(errors.New(fmt.Sprintf("unknown state %v", state)))
		}
	}
	normalText := string(normalBuffer)
	words = PreprocessWords(jieba.Cut(normalText, useHmm))
	tags = make([]string, 0)
	judges = make([]string, 0)
	fmt.Println("laws: ", laws)
	fmt.Println("words: ", words)
	return words, tags, laws, judges
}
