package emoji

import (
	"fmt"
	"strings"
)

// IsEmoji returns true if the input looks like an emoji rather than a file path or URL.
// It checks that the string is not a URL, not a file path, and contains at least one
// rune above U+00FF (i.e., outside basic Latin and Latin-1 Supplement).
func IsEmoji(s string) bool {
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return false
	}
	if strings.ContainsAny(s, "/\\.") {
		return false
	}
	for _, r := range s {
		if r > 0x00FF {
			return true
		}
	}
	return false
}

// Codepoint converts an emoji string to lowercase dash-separated hex codepoints.
func Codepoint(s string) string {
	var parts []string
	for _, r := range s {
		parts = append(parts, fmt.Sprintf("%x", r))
	}
	return strings.Join(parts, "-")
}

// CDNURL returns the jsdelivr CDN URL for an Apple emoji PNG given its codepoint string.
func CDNURL(codepoint string) string {
	return fmt.Sprintf("https://cdn.jsdelivr.net/npm/emoji-datasource-apple@16.0.0/img/apple/64/%s.png", codepoint)
}
