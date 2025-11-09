#!/bin/bash

# Simple package builder for Ducla Cloud Agent

set -e

VERSION="1.0.0"
PACKAGE_NAME="ducla-agent"

echo "🚀 Building packages for Ducla Cloud Agent v${VERSION}..."

# Check if fpm is available
if ! command -v fpm &> /dev/null; then
    echo "📦 Installing fpm (Effing Package Management)..."
    
    # Install Ruby and fpm
    sudo apt update
    sudo apt install -y ruby-dev build-essential rpm
    sudo gem install fpm
    
    echo "✅ fpm installed successfully"
fi

# Ensure we have the binary
if [ ! -f "dist/ducla-agent" ]; then
    echo "📦 Building binary first..."
    ./build-v1.sh
fi

# Create package directories
echo "📁 Creating package structure..."
mkdir -p pkg/{deb,rpm}

# Create simple DEB package
echo "📦 Building DEB package..."
fpm -s dir -t deb \
    --name "${PACKAGE_NAME}" \
    --version "${VERSION}" \
    --architecture "amd64" \
    --description "Ducla Cloud Agent - High-performance cloud agent" \
    --maintainer "mandỵhades <mandỵhades@hotmail.com.vn>" \
    --url "https://github.com/duclacloud/DUCLA-CLOUD-AGENT" \
    --license "MIT" \
    --package pkg/deb/ \
    dist/ducla-agent-v1.0.0=/usr/bin/ducla-agent

# Create simple RPM package  
echo "📦 Building RPM package..."
fpm -s dir -t rpm \
    --name "${PACKAGE_NAME}" \
    --version "${VERSION}" \
    --architecture "amd64" \
    --description "Ducla Cloud Agent - High-performance cloud agent" \
    --maintainer "mandỵhades <mandỵhades@hotmail.com.vn>" \
    --url "https://github.com/duclacloud/DUCLA-CLOUD-AGENT" \
    --license "MIT" \
    --package pkg/rpm/ \
    dist/ducla-agent-v1.0.0=/usr/bin/ducla-agent

echo ""
echo "🎉 Packages built successfully!"
echo ""
echo "📦 DEB Package:"
ls -la pkg/deb/*.deb
echo ""
echo "📦 RPM Package:"
ls -la pkg/rpm/*.rpm
echo ""
echo "📋 Installation commands:"
echo "  DEB: sudo dpkg -i pkg/deb/${PACKAGE_NAME}_${VERSION}_amd64.deb"
echo "  RPM: sudo rpm -ivh pkg/rpm/${PACKAGE_NAME}-${VERSION}-1.amd64.rpm"
echo ""
echo "🔍 Test installation:"
echo "  sudo dpkg -i pkg/deb/${PACKAGE_NAME}_${VERSION}_amd64.deb"
echo "  ducla-agent -version"
echo "  sudo dpkg -r ${PACKAGE_NAME}"