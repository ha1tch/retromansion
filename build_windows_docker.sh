# This script assumes a MacOS or Linux host, and that you have a 
# container management system running, such as Docker or Orbstack.
# You may use it to build Windows binaries using MinGW cross-compilation 
# in a Linux container, avoiding the need to install MinGW locally or 
# deal with Windows licensing issues.

#!/bin/bash
BUILDDIR=build
BINARY=retromansion.exe
TARGETBIN=$BUILDDIR/$BINARY

mkdir -p $BUILDDIR

docker run --rm --platform linux/amd64 -v $(pwd):/app -w /app golang:1.24 \
  sh -c "
    apt-get update -qq && \
    apt-get install -y -qq \
      gcc-mingw-w64-x86-64 \
      g++-mingw-w64-x86-64 \
      build-essential && \
    export CGO_ENABLED=1 && \
    export GOOS=windows && \
    export GOARCH=amd64 && \
    export CC=x86_64-w64-mingw32-gcc && \
    export CXX=x86_64-w64-mingw32-g++ && \
    go build -ldflags=\"-s -w\" -o $TARGETBIN game.go assets.go types.go render.go world.go rooms.go
  "

if [ $? -eq 0 ]; then
    echo "✅ Windows build successful!"
    echo "Binary file details:"
    ls -lh $TARGETBIN
    file   $TARGETBIN
    du -h  $TARGETBIN
else
    echo ""
    echo "❌ BUILD ERROR"
    exit 1
fi