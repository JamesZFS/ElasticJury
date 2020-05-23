package common

import (
	"os"
	"strings"
)

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
