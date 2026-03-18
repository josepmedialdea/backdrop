package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	colorpkg "github.com/josepmedialdea/backdrop/internal/color"
	"github.com/josepmedialdea/backdrop/internal/emoji"
	imgpkg "github.com/josepmedialdea/backdrop/internal/image"
	"github.com/spf13/cobra"
)

func main() {
	var (
		colorStr string
		output   string
		force    bool
		square   bool
		padding  string
	)

	rootCmd := &cobra.Command{
		Use:   "backdrop <image|emoji>",
		Short: "Fill transparent image backgrounds with a solid color",
		Long:  "backdrop takes an image (file, URL, or emoji) with a transparent background and fills it with a solid color. Supports --square to force square output and --padding to add breathing room.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], colorStr, output, force, square, padding)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.Flags().StringVarP(&colorStr, "color", "c", "#000000", "Background color as hex (#rrggbb) or R,G,B")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: <input>_bg.<ext>)")
	rootCmd.Flags().BoolVar(&force, "force", false, "Overwrite output file if it already exists")
	rootCmd.Flags().BoolVar(&square, "square", false, "Make output image a perfect square")
	rootCmd.Flags().StringVar(&padding, "padding", "0", "Padding on all sides: pixels (e.g. 20) or percentage (e.g. 10%)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(source, colorStr, output string, force, square bool, paddingStr string) error {
	// Parse color.
	bg, err := colorpkg.Parse(colorStr)
	if err != nil {
		return err
	}

	// Detect emoji input and convert to CDN URL.
	isEmoji := emoji.IsEmoji(source)
	loadSource := source
	if isEmoji {
		cp := emoji.Codepoint(source)
		loadSource = emoji.CDNURL(cp)
	}

	// Load image.
	img, err := imgpkg.Load(loadSource)
	if err != nil {
		return err
	}

	// Resolve padding (needs image dimensions for percentage mode).
	bounds := img.Bounds()
	padding, err := parsePadding(paddingStr, bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	// Validate transparency.
	if !imgpkg.HasTransparency(img) {
		return fmt.Errorf("image has no transparent background")
	}

	// Resolve output path.
	outPath := resolveOutput(source, output, isEmoji)

	// Check overwrite.
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("output already exists. Use --force to overwrite")
		}
	}

	// Fill background and save.
	opts := imgpkg.Options{Square: square, Padding: padding}
	result := imgpkg.FillBackground(img, bg, opts)
	if err := imgpkg.Save(result, outPath); err != nil {
		return err
	}

	fmt.Printf("Saved: %s\n", outPath)
	return nil
}

func parsePadding(s string, w, h int) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("padding value must not be empty")
	}

	if strings.HasSuffix(s, "%") {
		numStr := strings.TrimSuffix(s, "%")
		if numStr == "" {
			return 0, fmt.Errorf("invalid padding percentage: %q", s)
		}
		pct, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid padding percentage: %q", s)
		}
		if pct < 0 || pct > 100 {
			return 0, fmt.Errorf("padding percentage must be between 0%% and 100%%")
		}
		shorter := w
		if h < w {
			shorter = h
		}
		px := int(math.Round(pct / 100 * float64(shorter)))
		if px > 10000 {
			px = 10000
		}
		return px, nil
	}

	px, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid padding value: %q", s)
	}
	if px < 0 || px > 10000 {
		return 0, fmt.Errorf("padding must be between 0 and 10000")
	}
	return px, nil
}

func resolveOutput(source, output string, isEmoji bool) string {
	if output != "" {
		return output
	}

	if isEmoji {
		cp := emoji.Codepoint(source)
		return "emoji_" + cp + "_bg.png"
	}

	// For URLs, extract the filename and place output in cwd.
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		u := source
		// Strip fragment and query params.
		if idx := strings.Index(u, "#"); idx != -1 {
			u = u[:idx]
		}
		if idx := strings.Index(u, "?"); idx != -1 {
			u = u[:idx]
		}
		u = strings.TrimRight(u, "/")
		parts := strings.Split(u, "/")
		name := parts[len(parts)-1]
		ext := strings.ToLower(filepath.Ext(name))
		// Fallback if no usable image filename.
		if name == "" || (ext != ".png" && ext != ".jpg" && ext != ".jpeg") {
			return "image_bg.png"
		}
		stem := strings.TrimSuffix(name, ext)
		return stem + "_bg" + ext
	}

	// Local file path.
	dir := filepath.Dir(source)
	ext := filepath.Ext(source)
	stem := strings.TrimSuffix(filepath.Base(source), ext)

	return filepath.Join(dir, stem+"_bg"+ext)
}
