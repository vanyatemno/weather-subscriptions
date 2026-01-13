package util

import "regexp"

func RemoveAiLinks(text string) string {
	re := regexp.MustCompile(`(?m)\(\[[^\]]+\]\([^)]+\)\)`)
	return re.ReplaceAllString(text, "")
}
