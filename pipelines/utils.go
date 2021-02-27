package main

import (
	"fmt"
	"strings"
)

func GetOrDefault(data []string, ix int) string {
	if len(data) >= ix+1 {
		return data[ix]
	}
	return ""
}

func CleanString(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}

func PrepCsv(s string) ([]string, error) {
	data := SplitAtCommas(s)
	fmt.Println(data)
	for i := range data {
		data[i] = CleanString(data[i])
	}
	// Skip header
	if data[0] == "Title" || data[0] == "" {
		return nil, fmt.Errorf("Invalid data")
	}
	return data, nil
}

func ParseRemoved(data []string, ix int) bool {
	d := GetOrDefault(data, ix)
	if d == "" {
		return false
	} else {
		return true
	}
}

// https://stackoverflow.com/questions/59297737/go-split-string-by-comma-but-ignore-comma-within-double-quotes
func SplitAtCommas(s string) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}
