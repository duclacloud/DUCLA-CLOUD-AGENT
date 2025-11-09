#!/bin/bash

# Ducla Cloud Agent Uninstall Script

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SERVICE_NAME="ducla-agent"

echo -e "${YELLOW}🗑️ Uninstalling Ducla Cloud Agent...${NC}"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root${NC}"
   exit 1
fi

# Stop and disable service
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${YELLOW}⏹️ Stopping service...${NC}"
    systemctl stop $SERVICE_NAME
fi

if systemctl is-enabled --quiet $SERVICE_NAME; then
    echo -e "${YELLOW}🔧 Disabling service...${NC}"
    systemctl disable $SERVICE_NAME
fi

# Detect OS and remove package
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    
    case $OS in
        ubuntu|debian)
            if dpkg -l | grep -q ducla-agent; then
                echo -e "${YELLOW}📦 Removing DEB package...${NC}"
                dpkg -r ducla-agent
            fi
            ;;
        centos|rhel|fedora)
            if rpm -q ducla-agent >/dev/null 2>&1; then
                echo -e "${YELLOW}📦 Removing RPM package...${NC}"
                rpm -e ducla-agent
            fi
            ;;
        *)
            echo -e "${YELLOW}🗑️ Removing manual installation...${NC}"
            rm -rf /opt/ducla
            rm -f /etc/systemd/system/$SERVICE_NAME.service
            rm -rf /etc/ducla
            systemctl daemon-reload
            ;;
    esac
fi

# Remove user (optional)
if id "ducla" &>/dev/null; then
    echo -e "${YELLOW}👤 Removing user 'ducla'...${NC}"
    userdel ducla 2>/dev/null || true
fi

# Remove logs (optional)
if [ -d "/var/log/ducla" ]; then
    echo -e "${YELLOW}📝 Removing logs...${NC}"
    rm -rf /var/log/ducla
fi

echo -e "${GREEN}✅ Ducla Cloud Agent uninstalled successfully!${NC}"
