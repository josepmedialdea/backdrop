package color

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// Parse parses a color string in either hex (#rrggbb) or R,G,B format.
func Parse(s string) (color.NRGBA, error) {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "#") {
		return parseHex(s)
	}

	if strings.Contains(s, ",") {
		return parseRGB(s)
	}

	return color.NRGBA{}, fmt.Errorf("invalid color %q. Use #rrggbb or R,G,B", s)
}

func parseHex(s string) (color.NRGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.NRGBA{}, fmt.Errorf("invalid color \"#%s\". Use #rrggbb or R,G,B", s)
	}

	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("invalid color \"#%s\". Use #rrggbb or R,G,B", s)
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("invalid color \"#%s\". Use #rrggbb or R,G,B", s)
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("invalid color \"#%s\". Use #rrggbb or R,G,B", s)
	}

	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

func parseRGB(s string) (color.NRGBA, error) {
	parts := strings.SplitN(s, ",", 3)
	if len(parts) != 3 {
		return color.NRGBA{}, fmt.Errorf("invalid color %q. Use #rrggbb or R,G,B", s)
	}

	vals := [3]uint8{}
	for i, p := range parts {
		v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 8)
		if err != nil {
			return color.NRGBA{}, fmt.Errorf("invalid color %q. Use #rrggbb or R,G,B", s)
		}
		vals[i] = uint8(v)
	}

	return color.NRGBA{R: vals[0], G: vals[1], B: vals[2], A: 255}, nil
}
