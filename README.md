# png2zpl

`png2zpl` is a Go tool to convert PNG images into ZPL (`^GFA`) commands compatible with Zebra thermal printers.

## Features

- Convert PNG images (1-bit, grayscale, or paletted) to ZPL.
- Zebra-compatible run-length RLE compression to reduce ZPL size.
- Cross-platform static binaries (macOS, Linux, Windows).
- Option to invert black/white colors.

## Installation

### Manual Build

Make sure you have Go installed (>= 1.16). Then, in the project folder:

```bash
make release
```

Generated binaries:
- `png2zpl-macos-amd64`
- `png2zpl-macos-arm64`
- `png2zpl-linux-amd64`
- `png2zpl-linux-arm64`
- `png2zpl-windows-amd64.exe`

### Usage

```bash
# Generate ZPL from a PNG image
./png2zpl -input logo.png -output logo.zpl

# Invert colors (useful if the image is negative)
./png2zpl -input logo.png -output logo.zpl -invert
```

The generated `.zpl` file can be sent directly to a Zebra printer.

## Cross-Platform Build Examples

```bash
# macOS x86_64
GOOS=darwin GOARCH=amd64 go build -o png2zpl-macos-amd64 png2zpl.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o png2zpl-linux-arm64 png2zpl.go

# Windows x86_64
GOOS=windows GOARCH=amd64 go build -o png2zpl-windows-amd64.exe png2zpl.go
```

## Cleaning

```bash
make clean
```
Removes all generated binaries.

## License

This project is open-source and free to use.

