package common

import (
	"fmt"
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

func (c Conditions) Len() int {
	return len(c)
}

func (c Conditions) Less(i, j int) bool {
	return c[i].Weight > c[j].Weight
}

func (c Conditions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
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

func IndexOfStr(strs []string, target string) int {
	for i, str := range strs {
		if str == target {
			return i
		}
	}
	return -1
}
