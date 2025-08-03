@echo off
REM This builds natively on Windows using the local Go compiler and C toolchain
REM Requires Go and a C compiler (Visual Studio Build Tools, MinGW, or TDM-GCC)
REM For best results, run from a Developer Command Prompt for Visual Studio

setlocal EnableDelayedExpansion

set BUILDDIR=build
set BINARY=retromansion.exe
set TARGETBIN=%BUILDDIR%\%BINARY%

echo Building Windows binary natively...

REM Create build directory
if not exist "%BUILDDIR%" mkdir "%BUILDDIR%"

REM Check if Go is available
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go compiler not found! Please install Go and add it to PATH.
    exit /b 1
)

REM Set build environment
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

REM Try to detect and use the best available C compiler
where cl.exe >nul 2>&1
if %errorlevel% equ 0 (
    echo Using Visual Studio C compiler...
    set CC=cl.exe
    set CXX=cl.exe
) else (
    where gcc.exe >nul 2>&1
    if %errorlevel% equ 0 (
        echo Using GCC compiler...
        set CC=gcc
        set CXX=g++
    ) else (
        echo ⚠️  No C compiler detected. Build may fail.
        echo    Install Visual Studio Build Tools or MinGW-w64
    )
)

REM Build the binary
echo Running: go build -ldflags="-s -w" -o "%TARGETBIN%" game.go assets.go types.go render.go world.go rooms.go
go build -ldflags="-s -w" -o "%TARGETBIN%" game.go assets.go types.go render.go world.go rooms.go

REM Check build result
if %errorlevel% equ 0 (
    echo.
    echo ✅ Windows build successful!
    echo Binary file details:
    dir "%TARGETBIN%" | findstr retromansion
    echo File size: 
    for %%A in ("%TARGETBIN%") do echo    %%~zA bytes
    if exist "%TARGETBIN%" (
        echo File exists: %TARGETBIN%
    )
) else (
    echo.
    echo ❌ BUILD ERROR
    echo.
    echo Troubleshooting:
    echo - Ensure Go is installed and in PATH
    echo - Install Visual Studio Build Tools or MinGW-w64
    echo - Run from Developer Command Prompt if using Visual Studio
    exit /b 1
)

endlocal