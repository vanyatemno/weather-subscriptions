package util

import "regexp"

// RemoveAiLinks removes markdown links in parentheses like ([text](url)) from the input.
func RemoveAiLinks(text string) string {
	re := regexp.MustCompile(`(?m)\(\[[^\]]+\]\([^)]+\)\)`)
	return re.ReplaceAllString(text, "")
}
