package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"
)

// Load reads and decodes an image from a local path or HTTP(S) URL.
func Load(source string) (image.Image, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return loadURL(source)
	}
	return loadFile(source)
}

func loadFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("unsupported image format")
	}
	return img, nil
}

func loadURL(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	req.Header.Set("User-Agent", "backdrop/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image (HTTP %d)", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unsupported image format")
	}
	return img, nil
}

// HasTransparency checks whether the image contains at least one fully transparent pixel.
func HasTransparency(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a == 0 {
				return true
			}
		}
	}
	return false
}

// Options controls canvas adjustments applied during background fill.
type Options struct {
	Square  bool
	Padding int
}

// FillBackground creates a new image with the given background color,
// optionally squaring and/or padding the canvas, then composites the source on top.
func FillBackground(src image.Image, bg color.NRGBA, opts Options) *image.NRGBA {
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()

	canvasW, canvasH := w, h
	if opts.Square {
		side := max(w, h)
		canvasW, canvasH = side, side
	}
	if opts.Padding > 0 {
		canvasW += 2 * opts.Padding
		canvasH += 2 * opts.Padding
	}

	offsetX := (canvasW - w) / 2
	offsetY := (canvasH - h) / 2

	dst := image.NewNRGBA(image.Rect(0, 0, canvasW, canvasH))

	// Fill with background color.
	draw.Draw(dst, dst.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Composite source on top at the centered offset.
	dstRect := image.Rect(offsetX, offsetY, offsetX+w, offsetY+h)
	draw.Draw(dst, dstRect, src, src.Bounds().Min, draw.Over)

	return dst
}

// Save encodes the image as PNG and writes it to the given path.
func Save(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
