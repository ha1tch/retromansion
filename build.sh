#!/bin/bash
# This builds natively for whatever host platform you're on
# Requires appropriate development tools for each platform:
# - macOS: Xcode command line tools
# - Linux: build-essential and graphics/audio dev packages  
# - Windows: Visual Studio Build Tools or MinGW

# Auto-detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    darwin)  BINARY=retromansion_macos ;;
    linux)   BINARY=retromansion_linux ;;
    mingw*|cygwin*|msys*) BINARY=retromansion.exe ;;
    *)       BINARY=retromansion_${OS} ;;
esac

BUILDDIR=build
TARGETBIN=$BUILDDIR/$BINARY

mkdir -p $BUILDDIR

echo "Building for $(uname -s) $(uname -m)..."

go build -ldflags="-s -w" -o $TARGETBIN game.go assets.go types.go render.go world.go rooms.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    ls -lh $TARGETBIN
    file $TARGETBIN
    du -h $TARGETBIN
else
    echo "❌ Build failed!"
    exit 1
fi