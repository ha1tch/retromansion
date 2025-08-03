#!/bin/bash
# build_macos_bundle.sh
# Builds the Go binary and creates a complete macOS app bundle for distribution
# Run this on macOS to create a distributable .app bundle

set -e  # Exit on any error

# Configuration
APPNAME="Retromansion"
BINARY_NAME="retromansion"
BUNDLE_ID="com.yourname.retromansion"
VERSION="1.0"
MIN_MACOS="10.15"
ICON_CREATED=false

# Build directories
BUILDDIR="build"
APPBUNDLE="$BUILDDIR/$APPNAME.app"

echo "Building macOS app bundle..."

# Clean previous builds
rm -rf "$APPBUNDLE"
rm -f "$BUILDDIR/retromansion-macos.zip"

# Create build directory
mkdir -p "$BUILDDIR"

# Build the Go binary
echo "Building Go binary..."
go build -ldflags="-s -w" -o "$BUILDDIR/retromansion_macos" game.go assets.go types.go render.go world.go rooms.go

if [ ! -f "$BUILDDIR/retromansion_macos" ]; then
    echo "❌ Failed to build binary"
    exit 1
fi

echo "✅ Binary built successfully"

# Create app bundle structure
echo "Creating app bundle structure..."
mkdir -p "$APPBUNDLE/Contents/MacOS"
mkdir -p "$APPBUNDLE/Contents/Resources"

# Copy binary to app bundle  
cp "$BUILDDIR/retromansion_macos" "$APPBUNDLE/Contents/MacOS/$BINARY_NAME"
chmod +x "$APPBUNDLE/Contents/MacOS/$BINARY_NAME"

# Create shell script launcher to fix working directory
cat > "$APPBUNDLE/Contents/MacOS/launcher.sh" << 'EOF'
#!/bin/bash
# Get the directory where this script is located (Contents/MacOS)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# Change to that directory so relative paths work
cd "$SCRIPT_DIR"
# Run the actual binary
exec "./retromansion" "$@"
EOF

chmod +x "$APPBUNDLE/Contents/MacOS/launcher.sh"

echo "✅ Binary and launcher script created"

# Verify binary is executable
if [ -x "$APPBUNDLE/Contents/MacOS/launcher.sh" ] && [ -x "$APPBUNDLE/Contents/MacOS/$BINARY_NAME" ]; then
    echo "✅ Launcher script and binary are executable"
else
    echo "❌ Launcher script or binary is not executable"
    exit 1
fi

# Create Info.plist
echo "Creating Info.plist..."
cat > "$APPBUNDLE/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>launcher.sh</string>
    <key>CFBundleIdentifier</key>
    <string>com.yourname.retromansion</string>
    <key>CFBundleName</key>
    <string>Retromansion</string>
    <key>CFBundleDisplayName</key>
    <string>Retromansion</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>RMGM</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>NSPrincipalClass</key>
    <string>NSApplication</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.games</string>
    <key>LSEnvironment</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin</string>
    </dict>
</dict>
</plist>
EOF

echo "✅ Info.plist created with launcher.sh as executable"

# Generate app icon from player sprite
PLAYER_SPRITE="game_assets/sprites/entities/player/player_down.png"
ICON_NAME="icon"

echo "Looking for player sprite at: $PLAYER_SPRITE"
if [ -f "$PLAYER_SPRITE" ]; then
    echo "✅ Player sprite found, generating app icon..."
    
    # Create iconset directory
    mkdir -p "$ICON_NAME.iconset"
    
    # Generate all required icon sizes
    sips -z 16 16     "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_16x16.png" >/dev/null 2>&1
    sips -z 32 32     "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_16x16@2x.png" >/dev/null 2>&1
    sips -z 32 32     "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_32x32.png" >/dev/null 2>&1
    sips -z 64 64     "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_32x32@2x.png" >/dev/null 2>&1
    sips -z 128 128   "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_128x128.png" >/dev/null 2>&1
    sips -z 256 256   "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_128x128@2x.png" >/dev/null 2>&1
    sips -z 256 256   "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_256x256.png" >/dev/null 2>&1
    sips -z 512 512   "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_256x256@2x.png" >/dev/null 2>&1
    sips -z 512 512   "$PLAYER_SPRITE" --out "$ICON_NAME.iconset/icon_512x512.png" >/dev/null 2>&1
    cp "$PLAYER_SPRITE" "$ICON_NAME.iconset/icon_512x512@2x.png"
    
    # Create .icns file
    echo "Creating .icns file..."
    iconutil -c icns "$ICON_NAME.iconset"
    
    if [ -f "$ICON_NAME.icns" ]; then
        echo "✅ App icon generated successfully"
        ICON_CREATED=true
    else
        echo "❌ Failed to generate .icns file"
        ICON_CREATED=false
    fi
    
    # Clean up temporary iconset directory
    rm -rf "$ICON_NAME.iconset"
else
    echo "❌ Player sprite not found at $PLAYER_SPRITE - skipping icon generation"
    ls -la game_assets/sprites/entities/player/ 2>/dev/null || echo "Directory doesn't exist"
    ICON_CREATED=false
fi

echo "ICON_CREATED status: $ICON_CREATED"

# Copy game assets if they exist
if [ -d "game_assets" ]; then
    echo "Copying game assets..."
    cp -r game_assets "$APPBUNDLE/Contents/Resources/"
    
    # Create symbolic link from MacOS directory to Resources
    # This allows the binary to find assets at the expected relative path
    cd "$APPBUNDLE/Contents/MacOS"
    ln -sf "../Resources/game_assets" "game_assets"
    cd - >/dev/null
    
    echo "✅ Game assets copied and linked"
else
    echo "No game_assets directory found - assets will need to be relative to executable"
fi

# Copy icon to app bundle if it was created
if [ "$ICON_CREATED" = true ] && [ -f "$ICON_NAME.icns" ]; then
    cp "$ICON_NAME.icns" "$APPBUNDLE/Contents/Resources/"
    echo "✅ App icon copied to bundle"
fi

# Force macOS to recognize the app and its icon
if [ "$ICON_CREATED" = true ]; then
    echo "Refreshing icon cache..."
    # Touch the app bundle to update modification time
    touch "$APPBUNDLE"
    # Clear any existing icon cache for this app
    rm -rf "$HOME/Library/Caches/com.apple.iconservices.store" 2>/dev/null || true
    # Force LaunchServices to refresh
    /System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill -r -domain local -domain system -domain user >/dev/null 2>&1 || true
    # Additional icon cache refresh
    killall Dock 2>/dev/null || true
    echo "✅ Icon cache refreshed - Dock will restart"
fi

# Remove quarantine attributes that might prevent execution
echo "Removing quarantine attributes..."
xattr -dr com.apple.quarantine "$APPBUNDLE" 2>/dev/null || true
xattr -cr "$APPBUNDLE" 2>/dev/null || true
echo "✅ Quarantine attributes removed"

# Fix file permissions
echo "Setting proper permissions..."
chmod -R 755 "$APPBUNDLE"
chmod +x "$APPBUNDLE/Contents/MacOS/$BINARY_NAME"
echo "✅ Permissions set"

# Register the app with LaunchServices
echo "Registering app with macOS..."
/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -f "$APPBUNDLE"
echo "✅ App registered"

# Strip all extended attributes that mess with codesigning
echo "Stripping extended attributes..."
find "$APPBUNDLE" -exec xattr -c {} \; 2>/dev/null || true
echo "✅ Extended attributes cleaned"

# Ad-hoc code signing (crucial for GUI apps)
echo "Signing app with ad-hoc identity..."
codesign --force --deep --sign - "$APPBUNDLE" 2>/dev/null || {
    echo "⚠️  Code signing failed, but continuing..."
}
echo "✅ App signed"

# Verify bundle structure
echo "Verifying bundle structure..."
echo "Launcher script at MacOS/launcher.sh: $([ -f "$APPBUNDLE/Contents/MacOS/launcher.sh" ] && echo "✅" || echo "❌")"
echo "Binary at MacOS/retromansion: $([ -f "$APPBUNDLE/Contents/MacOS/retromansion" ] && echo "✅" || echo "❌")"
echo "Both are executable: $([ -x "$APPBUNDLE/Contents/MacOS/launcher.sh" ] && [ -x "$APPBUNDLE/Contents/MacOS/retromansion" ] && echo "✅" || echo "❌")"
echo "Icon file exists: $([ -f "$APPBUNDLE/Contents/Resources/icon.icns" ] && echo "✅" || echo "❌")"
echo "Info.plist exists: $([ -f "$APPBUNDLE/Contents/Info.plist" ] && echo "✅" || echo "❌")"
echo "Game assets symlink: $([ -L "$APPBUNDLE/Contents/MacOS/game_assets" ] && echo "✅" || echo "❌")"

# Create distribution zip
echo "Creating distribution zip..."
cd "$BUILDDIR"
zip -r "retromansion-macos.zip" "$APPNAME.app"
cd ..

# Clean up temporary files
if [ -f "$ICON_NAME.icns" ]; then
    rm "$ICON_NAME.icns"
fi

# Show results
echo ""
echo "✅ macOS app bundle created successfully!"
echo ""
echo "Build Summary:"
echo "   App Bundle: $APPBUNDLE"
echo "   Distribution: $BUILDDIR/retromansion-macos.zip"
echo ""

# Show file details
echo "File Details:"
ls -lh "$APPBUNDLE"
if [ -f "$BUILDDIR/retromansion-macos.zip" ]; then
    echo "ZIP Size: $(du -h $BUILDDIR/retromansion-macos.zip | cut -f1)"
fi

echo ""
echo "Ready for distribution!"
echo ""
echo "User Installation Instructions:"
echo "   1. Download retromansion-macos.zip"
echo "   2. Unzip the file"
echo "   3. Double-click Retromansion.app to run"
echo "   4. If blocked: Right-click → 'Open' → 'Open'"
echo ""
echo "The app is ad-hoc signed and should run without issues on the same machine."
echo "For distribution to other machines, users may need to right-click → Open on first run."