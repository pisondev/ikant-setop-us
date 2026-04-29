package shared

import "strings"

func NormalizeSlug(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	return value
}
