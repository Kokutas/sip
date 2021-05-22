package main

import (
	"fmt"
	"strings"
)

func main() {
	raws := []string{
		"token=xyz&expires=3600&xxxxxxx",
		"yyyyyy",
		"k=v",
	}
	for _, raw := range raws {
		rs := strings.Split(raw, "&")
		// fmt.Println(len(rs), rs)
		for _, r := range rs {
			kvs := strings.Split(r, "=")
			fmt.Println(len(kvs), kvs)
		}

	}
}
