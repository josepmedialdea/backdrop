package main

import "testing"

func TestResolveOutput(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		output  string
		isEmoji bool
		want    string
	}{
		{"explicit output", "logo.png", "custom.png", false, "custom.png"},
		{"local file", "logo.png", "", false, "logo_bg.png"},
		{"local file with dir", "images/logo.png", "", false, "images/logo_bg.png"},
		{"simple URL", "https://example.com/logo.png", "", false, "logo_bg.png"},
		{"URL with query", "https://example.com/logo.png?size=large", "", false, "logo_bg.png"},
		{"URL with fragment", "https://example.com/logo.png#section", "", false, "logo_bg.png"},
		{"URL with query and fragment", "https://example.com/logo.png?a=1#frag", "", false, "logo_bg.png"},
		{"URL trailing slash", "https://example.com/logo.png/", "", false, "logo_bg.png"},
		{"URL no filename", "https://example.com/", "", false, "image_bg.png"},
		{"URL no filename no slash", "https://example.com", "", false, "image_bg.png"},
		{"emoji", "🏠", "", true, "emoji_1f3e0_bg.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveOutput(tt.source, tt.output, tt.isEmoji)
			if got != tt.want {
				t.Errorf("resolveOutput(%q, %q, %v) = %q, want %q", tt.source, tt.output, tt.isEmoji, got, tt.want)
			}
		})
	}
}
