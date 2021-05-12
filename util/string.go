package util

import (
	"strings"
)

func TrimPrefixAndSuffix(raw string, substring string) (result string) {
	for strings.HasPrefix(raw, substring) && strings.HasSuffix(raw, substring) {
		raw = strings.TrimPrefix(raw, substring)
		raw = strings.TrimSuffix(raw, substring)
	}
	result = raw
	//fmt.Println("raw",raw,"----")
	//fmt.Println("result",result,"---")
	return result
}
