# This script assumes a MacOS host, and that you have a 
# container management system running, such as Docker or Orbstack.
# You may use it to build Linux binaries from Mac.

#!/bin/bash
BUILDDIR=build
BINARY=retromansion_linux
TARGETBIN=$BUILDDIR/$BINARY

mkdir -p $BUILDDIR

docker run --rm --platform linux/amd64 -v $(pwd):/app -w /app golang:1.24 \
  sh -c "
    apt-get update -qq && \
    apt-get install -y -qq build-essential \
      libgl1-mesa-dev \
      libxi-dev \
      libxcursor-dev \
      libxrandr-dev \
      libxinerama-dev \
      libasound2-dev \
      libwayland-dev \
      libxkbcommon-dev \
      wayland-protocols && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=\"-s -w\" -o $TARGETBIN game.go assets.go types.go render.go world.go rooms.go
  "

if [ $? -eq 0 ]; then
    echo "Binary file details:"
    ls -lh $TARGETBIN
    file   $TARGETBIN
    du -h  $TARGETBIN
else
    echo ""
    echo "❌ BUILD ERROR"
    exit 1
fi
