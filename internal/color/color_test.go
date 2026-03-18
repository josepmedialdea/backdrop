package color

import (
	"image/color"
	"testing"
)

func TestParseHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    color.NRGBA
		wantErr bool
	}{
		{"black", "#000000", color.NRGBA{R: 0, G: 0, B: 0, A: 255}, false},
		{"white", "#FFFFFF", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, false},
		{"red", "#FF0000", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, false},
		{"green", "#00FF00", color.NRGBA{R: 0, G: 255, B: 0, A: 255}, false},
		{"blue", "#0000FF", color.NRGBA{R: 0, G: 0, B: 255, A: 255}, false},
		{"mixed", "#1A1A2E", color.NRGBA{R: 26, G: 26, B: 46, A: 255}, false},
		{"lowercase", "#ff8800", color.NRGBA{R: 255, G: 136, B: 0, A: 255}, false},
		{"leading whitespace", "  #FF0000", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, false},
		{"trailing whitespace", "#FF0000  ", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, false},
		{"too short", "#FFF", color.NRGBA{}, true},
		{"too long", "#FF00000", color.NRGBA{}, true},
		{"invalid chars", "#GGGGGG", color.NRGBA{}, true},
		{"empty after hash", "#", color.NRGBA{}, true},
		{"invalid mid-channel", "#FF00GG", color.NRGBA{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRGB(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    color.NRGBA
		wantErr bool
	}{
		{"zeros", "0,0,0", color.NRGBA{R: 0, G: 0, B: 0, A: 255}, false},
		{"max values", "255,255,255", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, false},
		{"mixed", "128,64,32", color.NRGBA{R: 128, G: 64, B: 32, A: 255}, false},
		{"with spaces", "255 , 0 , 0", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, false},
		{"leading whitespace", "  255,0,0", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, false},
		{"trailing whitespace", "255,0,0  ", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, false},
		{"out of range", "256,0,0", color.NRGBA{}, true},
		{"negative", "-1,0,0", color.NRGBA{}, true},
		{"too few values", "255,0", color.NRGBA{}, true},
		{"empty part", "255,,0", color.NRGBA{}, true},
		{"non-numeric", "red,0,0", color.NRGBA{}, true},
		{"float value", "255.5,0,0", color.NRGBA{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseInvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"plain word", "red"},
		{"empty string", ""},
		{"only whitespace", "   "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatalf("Parse(%q) expected error, got nil", tt.input)
			}
		})
	}
}
