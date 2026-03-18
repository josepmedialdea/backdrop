package emoji

import "testing"

func TestIsEmoji(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"🏠", true},
		{"👩\u200d💻", true},
		{"🇪🇸", true},
		{"logo.png", false},
		{"/path/to/file.png", false},
		{"relative/path.png", false},
		{"https://example.com/image.png", false},
		{"http://example.com/image.png", false},
		{"hello", false},
		{"abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := IsEmoji(tt.input)
			if got != tt.want {
				t.Errorf("IsEmoji(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCodepoint(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"single emoji", "🏠", "1f3e0"},
		{"ZWJ sequence", "👩\u200d💻", "1f469-200d-1f4bb"},
		{"flag", "🇪🇸", "1f1ea-1f1f8"},
		{"with variation selector", "☁️", "2601-fe0f"},
		{"heart with variation selector", "❤️", "2764-fe0f"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Codepoint(tt.input)
			if got != tt.want {
				t.Errorf("Codepoint(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCDNURL(t *testing.T) {
	got := CDNURL("1f3e0")
	want := "https://cdn.jsdelivr.net/npm/emoji-datasource-apple@16.0.0/img/apple/64/1f3e0.png"
	if got != want {
		t.Errorf("CDNURL(\"1f3e0\") = %q, want %q", got, want)
	}
}
