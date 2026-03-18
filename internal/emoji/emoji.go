package emoji

import (
	"fmt"
	"strings"
)

// IsEmoji returns true if the input looks like a pure emoji string rather than
// a file path, URL, or mixed text. It rejects URLs, file paths, and any input
// containing ASCII letters (a-z, A-Z) to avoid ambiguity with mixed inputs
// like "🏠home".
func IsEmoji(s string) bool {
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return false
	}
	if strings.ContainsAny(s, "/\\.") {
		return false
	}
	hasEmoji := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return false
		}
		if r > 0x00FF {
			hasEmoji = true
		}
	}
	return hasEmoji
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
