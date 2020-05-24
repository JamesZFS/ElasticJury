package common

import (
	"fmt"
	"os"
	"strings"
)

type Param struct {
	TableName  string
	FieldName string
	Conditions []string
	Weight []float32
}

func BuildParam(tableName, fieldName, query string) Param {
	conditions := FilterStrs(strings.Split(query, ","), NotWhiteSpace)
	return Param {
		tableName,
		fieldName,
		conditions,
		make([]float32, len(conditions)),
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

func NotWhiteSpace(str string) bool {
	if strings.TrimSpace(str) == "" {
		return false
	} else {
		return true
	}
}

func GetOrExpr(entry int32, field string, conditions []string) string {
	var array []string
	for _, condition := range conditions {
		 array = append(array, fmt.Sprintf("%c.%s = '%s'", entry, field, condition))
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
