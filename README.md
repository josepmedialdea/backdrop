# 🖼️ backdrop

A CLI tool that fills transparent image backgrounds with a solid color. Supports local files, URLs, and emoji as input.

I built this to generate consistent project avatars for GitLab, but it works anywhere you need a quick, polished image with a colored background — favicons, social cards, README badges, etc.

## Features

- **Multiple input sources** — local files, HTTP(S) URLs, or emoji (e.g. `🏠`)
- **Flexible color input** — hex (`#FF5733`) or RGB (`255,87,51`)
- **Square mode** — force output to a perfect square
- **Padding** — add breathing room in pixels (`--padding 20`) or as a percentage of the image (`--padding 10%`)

## Installation

### From source

```sh
git clone https://github.com/josepmedialdea/backdrop.git
cd backdrop
make build
# binary is at ./bin/backdrop
```

### Go install

```sh
go install github.com/josepmedialdea/backdrop/cmd/backdrop@latest
```

## Usage

```
backdrop <image|emoji> [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--color` | `-c` | `#000000` | Background color as hex (`#rrggbb`) or `R,G,B` |
| `--output` | `-o` | auto | Output file path (default: `<input>_bg.<ext>`) |
| `--force` | | `false` | Overwrite output file if it already exists |
| `--square` | | `false` | Make output image a perfect square |
| `--padding` | | `0` | Padding on all sides: pixels (e.g. `20`) or percentage (e.g. `10%`) |

## Examples

```sh
# Fill a local PNG with a blue background
backdrop logo.png -c "#3498db"

# Fetch an image from a URL and apply a red background
backdrop https://example.com/icon.png -c 255,0,0

# Render an emoji with a dark background
backdrop 🦞 -c "#1a1a2e"

# Square output with 10% padding
backdrop logo.png -c "#2ecc71" --square --padding 10%

# Explicit output path, overwrite if exists
backdrop logo.png -c "#000000" -o result.png --force
```

## License

[MIT](LICENSE)
