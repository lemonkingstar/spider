package pconv

import "strings"

func StrSplit(s, sep string) []string {
	return strings.Split(s, sep)
}

func StrJoin(s []string, sep string) string {
	return strings.Join(s, sep)
}

func StrHasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func StrHasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func StrTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func StrContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func Str2Lower(s string) string {
	return strings.ToLower(s)
}

func Str2Upper(s string) string {
	return strings.ToUpper(s)
}
