#!/bin/bash
set -e

APP_NAME="memcached-go"
OUTPUT_DIR="dist"

# List of major GOOS/GOARCH targets
TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

mkdir -p "$OUTPUT_DIR"

VERSION=$(cat version.txt)

for target in "${TARGETS[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "$target"
    output_name="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    [ "$GOOS" = "windows" ] && output_name="${output_name}.exe"
    echo "Building for $GOOS/$GOARCH..."
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "-X 'memcached-go/internal/version.Version=${VERSION}'" -o "${OUTPUT_DIR}/${output_name}" ./cmd/
done

echo "Builds complete. Distributables are in the '$OUTPUT_DIR' directory."