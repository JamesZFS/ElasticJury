package common

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Condition struct {
	Item   string
	Weight float32
}

type Conditions []Condition

type Param struct {
	TableName  string
	FieldName  string
	Conditions Conditions
}

type ResultList []int32

type Strings []string

func (l ResultList) ToByteArray() []byte {
	ids := make([]byte, len(l) * 3)
	for i, item := range l {
		ids[i * 3 + 0] = byte((item >>  0) & 0xff)
		ids[i * 3 + 1] = byte((item >>  8) & 0xff)
		ids[i * 3 + 2] = byte((item >> 16) & 0xff)
	}
	return ids
}

func (c Conditions) ItemArray() []string {
	var array []string
	for _, condition := range c {
		array = append(array, condition.Item)
	}
	return array
}

func (c Conditions) Len() int {
	return len(c)
}

func (c Conditions) Less(i, j int) bool {
	return c[i].Weight > c[j].Weight
}

func (c Conditions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (s Strings) Len() int {
	return len(s)
}

func (s Strings) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func (s Strings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func MakeDefaultConditions(items []string) Conditions {
	items = FilterStrs(items, NotWhiteSpace)
	conditions := make(Conditions, len(items))
	for i, item := range items {
		conditions[i] = Condition{item, 1.0}
	}
	return conditions
}

func BuildParam(tableName, fieldName string, conditions Conditions) Param {
	return Param{
		tableName,
		fieldName,
		conditions,
	}
}

func GetEnvVar(key string, dft string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return dft
	}
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

func Unique(items []Condition) []Condition {
	set := make(map[string]bool)
	var reduced []Condition
	for _, item := range items {
		if !set[item.Item] {
			set[item.Item] = true
			reduced = append(reduced, item)
		}
	}
	return reduced
}

func UniqueShuffle(items []string) []string {
	set := make(map[string]bool)
	var shuffled []string
	for _, item := range items {
		if !set[item] {
			set[item] = true
			shuffled = append(shuffled, item)
		}
	}
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func NotWhiteSpace(str string) bool {
	if strings.TrimSpace(str) == "" {
		return false
	} else {
		return true
	}
}

func GetOrExpr(entry int32, field string, conditions Conditions) string {
	var array []string
	for _, condition := range conditions {
		array = append(array, fmt.Sprintf("%c.%s = '%s'", entry, field, condition.Item))
	}
	return strings.Join(array, " OR ")
}
