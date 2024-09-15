package utils

import (
	"regexp"
)

func SanitizeFilename(filename string) string {
	// Remove anything that isn't a word character (a-z, A-Z, 0-9, _), a hyphen, a space, or a period.
	re := regexp.MustCompile(`[^\w\-. ]`)
	sanitized := re.ReplaceAllString(filename, "")

	// Replace spaces with underscores.
	sanitized = regexp.MustCompile(` `).ReplaceAllString(sanitized, "_")

	// If the filename was originally just dangerous characters, it could now be empty, so we provide a default.
	if sanitized == "" {
		sanitized = "chailabs-default-filename"
	}

	return sanitized
}
