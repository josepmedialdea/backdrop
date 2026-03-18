package main

import (
	"testing"
)

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

func TestParsePadding(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		w, h    int
		want    int
		wantErr bool
	}{
		{"absolute pixels", "20", 100, 200, 20, false},
		{"zero", "0", 100, 100, 0, false},
		{"percent shorter width", "10%", 100, 200, 10, false},
		{"percent shorter width 2", "10%", 200, 400, 20, false},
		{"percent 50", "50%", 64, 64, 32, false},
		{"percent zero", "0%", 100, 100, 0, false},
		{"negative", "-5", 100, 100, 0, true},
		{"over max", "10001", 100, 100, 0, true},
		{"percent over 100", "101%", 100, 100, 0, true},
		{"non-numeric", "abc", 100, 100, 0, true},
		{"bare percent", "%", 100, 100, 0, true},
		{"empty string", "", 100, 100, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePadding(tt.s, tt.w, tt.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePadding(%q, %d, %d) error = %v, wantErr %v", tt.s, tt.w, tt.h, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parsePadding(%q, %d, %d) = %d, want %d", tt.s, tt.w, tt.h, got, tt.want)
			}
		})
	}
}
