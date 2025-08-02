# This script assumes a MacOS host and that the mingw32 compiler is available
# You can use this script to build Windows binaries from Mac

BUILDDIR=build
BINARY=retromansion.exe
TARGETBIN=$BUILDDIR/$BINARY

mkdir -p $BUILDDIR

export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++

# go build -o $BUILDDIR/debug_$BINARY.exe game.go assets.go types.go render.go world.go rooms.go

go build -ldflags="-s -w" -o $TARGETBIN game.go assets.go types.go render.go world.go rooms.go

ls -al $TARGETBIN
file $TARGETBIN
du -h $TARGETBIN
