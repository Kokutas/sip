package util

import "strings"

func TrimPrefixAndSuffix(raw string, substring string) (result string) {
	raw = strings.TrimLeft(raw, substring)
	raw = strings.TrimRight(raw,substring)
	raw = strings.TrimPrefix(raw,substring)
	raw = strings.TrimSuffix(raw,substring)
	result = raw
	return result
}
