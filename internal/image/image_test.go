package image

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

// newOpaqueImage creates an image with no transparent pixels.
func newOpaqueImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	return img
}

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

func TestHasTransparency_WithTransparent(t *testing.T) {
	img := newTestImage(10, 10) // mostly transparent
	if !HasTransparency(img) {
		t.Error("expected true for image with transparent pixels")
	}
}

func TestHasTransparency_AllOpaque(t *testing.T) {
	img := newOpaqueImage(10, 10)
	if HasTransparency(img) {
		t.Error("expected false for fully opaque image")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/image.png")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidFormat(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.png")
	if err := os.WriteFile(tmp, []byte("not an image"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(tmp)
	if err == nil {
		t.Fatal("expected error for invalid image data")
	}
}

func TestLoadFile_ValidPNG(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "valid.png")
	img := newTestImage(4, 4)
	if err := Save(img, tmp); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Bounds().Dx() != 4 || loaded.Bounds().Dy() != 4 {
		t.Errorf("expected 4x4, got %dx%d", loaded.Bounds().Dx(), loaded.Bounds().Dy())
	}
}

func TestSave_AndReload(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.png")
	src := newTestImage(8, 8)

	if err := Save(src, tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("output file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}

	loaded, err := Load(tmp)
	if err != nil {
		t.Fatalf("failed to reload saved image: %v", err)
	}
	if loaded.Bounds().Dx() != 8 || loaded.Bounds().Dy() != 8 {
		t.Errorf("reloaded image size = %dx%d, want 8x8", loaded.Bounds().Dx(), loaded.Bounds().Dy())
	}
}

func TestSave_BadPath(t *testing.T) {
	img := newTestImage(4, 4)
	err := Save(img, "/nonexistent/dir/out.png")
	if err == nil {
		t.Fatal("expected error for bad output path")
	}
}
