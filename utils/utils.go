package utils

import "strings"

func SplicingString(list ...string) string {
	var b strings.Builder
	for _, s := range list {
		b.WriteString(s)
	}
	return b.String()
}
