package image

import (
	"image"
	"image/color"
	"testing"
)

// newTestImage creates a small NRGBA image with transparent background and a colored center pixel.
func newTestImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	// Leave all pixels fully transparent (zero value), then set center pixel.
	cx, cy := w/2, h/2
	img.SetNRGBA(cx, cy, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	return img
}

func TestFillBackground_NoOptions(t *testing.T) {
	src := newTestImage(10, 20)
	bg := color.NRGBA{R: 0, G: 0, B: 255, A: 255}

	result := FillBackground(src, bg, Options{})

	bounds := result.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 20 {
		t.Errorf("expected 10x20, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Background pixel should be blue.
	c := result.NRGBAAt(0, 0)
	if c != bg {
		t.Errorf("background pixel = %v, want %v", c, bg)
	}

	// Center pixel should be red (composited over blue).
	center := result.NRGBAAt(5, 10)
	if center.R != 255 || center.A != 255 {
		t.Errorf("center pixel = %v, want red", center)
	}
}

func TestFillBackground_SquareOnly(t *testing.T) {
	src := newTestImage(10, 20)
	bg := color.NRGBA{R: 0, G: 255, B: 0, A: 255}

	result := FillBackground(src, bg, Options{Square: true})

	bounds := result.Bounds()
	if bounds.Dx() != 20 || bounds.Dy() != 20 {
		t.Errorf("expected 20x20, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Top-left corner should be background (padding area from squaring).
	c := result.NRGBAAt(0, 0)
	if c != bg {
		t.Errorf("corner pixel = %v, want %v", c, bg)
	}
}

func TestFillBackground_PaddingOnly(t *testing.T) {
	src := newTestImage(10, 10)
	bg := color.NRGBA{R: 128, G: 128, B: 128, A: 255}

	result := FillBackground(src, bg, Options{Padding: 5})

	bounds := result.Bounds()
	if bounds.Dx() != 20 || bounds.Dy() != 20 {
		t.Errorf("expected 20x20, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Top-left corner should be background (padding area).
	c := result.NRGBAAt(0, 0)
	if c != bg {
		t.Errorf("corner pixel = %v, want %v", c, bg)
	}
}

func TestFillBackground_SquareAndPadding(t *testing.T) {
	src := newTestImage(10, 20)
	bg := color.NRGBA{R: 255, G: 255, B: 0, A: 255}

	result := FillBackground(src, bg, Options{Square: true, Padding: 5})

	// Square of max(10,20)=20, plus 2*5 padding = 30x30.
	bounds := result.Bounds()
	if bounds.Dx() != 30 || bounds.Dy() != 30 {
		t.Errorf("expected 30x30, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// All corners should be background.
	for _, pt := range []image.Point{{0, 0}, {29, 0}, {0, 29}, {29, 29}} {
		c := result.NRGBAAt(pt.X, pt.Y)
		if c != bg {
			t.Errorf("corner pixel at %v = %v, want %v", pt, c, bg)
		}
	}
}
