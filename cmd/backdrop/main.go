package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	imgpkg "github.com/josepmedialdea/backdrop/internal/image"

	colorpkg "github.com/josepmedialdea/backdrop/internal/color"
	"github.com/spf13/cobra"
)

func main() {
	var (
		colorStr string
		output   string
		force    bool
	)

	rootCmd := &cobra.Command{
		Use:   "backdrop <image>",
		Short: "Fill transparent image backgrounds with a solid color",
		Long:  "backdrop takes an image with a transparent background and fills it with a solid color, saving the result to a specified output path.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], colorStr, output, force)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.Flags().StringVarP(&colorStr, "color", "c", "#000000", "Background color as hex (#rrggbb) or R,G,B")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: <input>_bg.<ext>)")
	rootCmd.Flags().BoolVar(&force, "force", false, "Overwrite output file if it already exists")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(source, colorStr, output string, force bool) error {
	// Parse color.
	bg, err := colorpkg.Parse(colorStr)
	if err != nil {
		return err
	}

	// Load image.
	img, err := imgpkg.Load(source)
	if err != nil {
		return err
	}

	// Validate transparency.
	if !imgpkg.HasTransparency(img) {
		return fmt.Errorf("image has no transparent background")
	}

	// Resolve output path.
	outPath := resolveOutput(source, output)

	// Check overwrite.
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("output already exists. Use --force to overwrite")
		}
	}

	// Fill background and save.
	result := imgpkg.FillBackground(img, bg)
	if err := imgpkg.Save(result, outPath); err != nil {
		return err
	}

	fmt.Printf("Saved: %s\n", outPath)
	return nil
}

func resolveOutput(source, output string) string {
	if output != "" {
		return output
	}

	// For URLs, use the filename in the current working directory.
	name := source
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		parts := strings.Split(source, "/")
		name = parts[len(parts)-1]
		// Strip query params.
		if idx := strings.Index(name, "?"); idx != -1 {
			name = name[:idx]
		}
	}

	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	dir := filepath.Dir(name)

	// For URLs (no directory component), use cwd.
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		dir = "."
	}

	return filepath.Join(dir, stem+"_bg"+ext)
}
