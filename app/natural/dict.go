package natural

import (
	. "ElasticJury/app/common"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Dict map[string]Strings
type Dicts []Dict

func TitleMarkFilter(word string) string {
	escaped := make([]int32, 0)
	for _, c := range strings.TrimSpace(word) {
		switch c {
		case '《', '》':
			// ignore
		default:
			escaped = append(escaped, c)
		}
	}
	return string(escaped)
}

func JoinDicts(dicts Dicts) Dict {
	merged := make(Dict)
	for _, d := range dicts {
		for k, v := range d {
			merged[k] = append(merged[k], v...)
		}
	}
	for k := range merged {
		merged[k] = UniqueShuffle(merged[k])
		merged[k] = merged[k][:Min(len(merged[k]), TipsCount)]
		sort.Sort(merged[k])
	}
	return merged
}

func BuildDict(path string) Dict {
	dict := make(Dict)
	conditionMap := make(map[string]Conditions)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("[Info] Failed to load file %s\n", err)
	}

	for _, item := range strings.Split(string(bytes), "\n") {
		kv := strings.Split(item, ",")
		full := kv[0]
		index := TitleMarkFilter(kv[0])
		weight, _ := strconv.ParseFloat(kv[1], 32)
		condition := Condition{Item: full, Weight: float32(weight)}
		for i, l := 0, utf8.RuneCountInString(index); i < l; i++ {
			key := Substring(index, 0, i+1)
			conditionMap[key] = append(conditionMap[key], condition)
		}

		words := jieba.CutForSearch(index, true)
		for _, word := range words {
			for i, l := 0, utf8.RuneCountInString(word); i < l; i++ {
				key := Substring(word, 0, i+1)
				conditionMap[key] = append(conditionMap[key], condition)
			}
		}
	}

	for key := range conditionMap {
		conditionMap[key] = Unique(conditionMap[key])
		sort.Sort(conditionMap[key])
		conditionMap[key] = conditionMap[key][0:Min(len(conditionMap[key]), TipsCount)]
		dict[key] = conditionMap[key].ItemArray()
		sort.Sort(dict[key])
	}
	return dict
}