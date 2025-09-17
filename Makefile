# Makefile pour png2zpl - builds multiplateformes statiques

BINARY=png2zpl
VERSION=1.0.0

.PHONY: all clean release

all: release

release: clean
	@echo "Building for macOS x86_64"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BINARY}-macos-amd64 png2zpl.go

	@echo "Building for macOS ARM64"
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${BINARY}-macos-arm64 png2zpl.go

	@echo "Building for Linux x86_64"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}-linux-amd64 png2zpl.go

	@echo "Building for Linux ARM64"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BINARY}-linux-arm64 png2zpl.go

	@echo "Building for Windows x86_64"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}-windows-amd64.exe png2zpl.go

clean:
	@echo "Cleaning old binaries..."
	rm -f ${BINARY}-macos-amd64 ${BINARY}-macos-arm64 ${BINARY}-linux-amd64 ${BINARY}-linux-arm64 ${BINARY}-windows-amd64.exe
