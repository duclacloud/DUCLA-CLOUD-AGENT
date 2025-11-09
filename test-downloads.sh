#!/bin/bash

# Test script for Ducla Cloud Agent v1.0.0 downloads
# This script simulates downloading from GitHub releases

set -e

RELEASE_DIR="releases/v1.0.0"
BASE_URL="https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0"

echo "🧪 Testing Ducla Cloud Agent v1.0.0 Download Links"
echo "=================================================="

# Function to simulate wget download
simulate_download() {
    local filename=$1
    local url="${BASE_URL}/${filename}"
    
    echo "📥 Testing: wget $url"
    
    if [ -f "${RELEASE_DIR}/${filename}" ]; then
        echo "✅ File exists: ${filename}"
        echo "📊 Size: $(du -h ${RELEASE_DIR}/${filename} | cut -f1)"
        echo "🔐 SHA256: $(sha256sum ${RELEASE_DIR}/${filename} | cut -d' ' -f1)"
    else
        echo "❌ File missing: ${filename}"
        return 1
    fi
    echo ""
}

# Test all download links
echo "🐧 Testing Ubuntu/Debian DEB package:"
simulate_download "ducla-agent_1.0.0_amd64.deb"

echo "🎩 Testing RHEL/CentOS RPM package:"
simulate_download "ducla-agent-1.0.0-1.x86_64.rpm"

echo "📦 Testing Binary tar.gz package:"
simulate_download "ducla-agent-linux-amd64.tar.gz"

echo "🔐 Testing checksums file:"
simulate_download "checksums.txt"

echo "📋 Testing release notes:"
simulate_download "RELEASE-NOTES.md"

# Verify checksums
echo "🔍 Verifying checksums:"
if cd "${RELEASE_DIR}" && sha256sum -c checksums.txt --quiet; then
    echo "✅ All checksums verified successfully!"
else
    echo "❌ Checksum verification failed!"
    exit 1
fi

echo ""
echo "🎉 All download links are ready for GitHub release!"
echo ""
echo "📋 Summary of files ready for upload:"
echo "======================================"
ls -la "${RELEASE_DIR}/" | grep -v "^d" | grep -v "binary-linux-amd64"

echo ""
echo "🚀 GitHub Release Commands:"
echo "=========================="
echo "1. Create release on GitHub:"
echo "   https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/new"
echo ""
echo "2. Tag: v1.0.0"
echo "3. Title: Ducla Cloud Agent v1.0.0 - First Stable Release"
echo "4. Upload these files:"
echo "   - ducla-agent_1.0.0_amd64.deb"
echo "   - ducla-agent-1.0.0-1.x86_64.rpm"
echo "   - ducla-agent-linux-amd64.tar.gz"
echo "   - checksums.txt"
echo "   - RELEASE-NOTES.md"
echo ""
echo "5. Use RELEASE-NOTES.md content as release description"
echo ""
echo "✨ After upload, these wget commands will work:"
echo "wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent_1.0.0_amd64.deb"
echo "wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent-1.0.0-1.x86_64.rpm"
echo "wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent-linux-amd64.tar.gz"