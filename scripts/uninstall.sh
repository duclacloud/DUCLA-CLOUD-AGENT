#!/bin/bash

# Uninstallation script for Ducla Cloud Agent

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/ducla"
DATA_DIR="/opt/ducla"
LOG_DIR="/var/log/ducla"
SERVICE_USER="ducla"

echo -e "${YELLOW}Ducla Cloud Agent Uninstaller${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root or with sudo${NC}"
    exit 1
fi

# Stop and disable service
if command -v systemctl &> /dev/null; then
    if systemctl is-active --quiet ducla-agent; then
        echo -e "${YELLOW}Stopping service...${NC}"
        systemctl stop ducla-agent
    fi
    
    if systemctl is-enabled --quiet ducla-agent; then
        echo -e "${YELLOW}Disabling service...${NC}"
        systemctl disable ducla-agent
    fi
    
    if [ -f /etc/systemd/system/ducla-agent.service ]; then
        rm /etc/systemd/system/ducla-agent.service
        systemctl daemon-reload
    fi
fi

# Remove binary
if [ -f "$INSTALL_DIR/ducla-agent" ]; then
    echo -e "${YELLOW}Removing binary...${NC}"
    rm "$INSTALL_DIR/ducla-agent"
fi

# Ask about data removal
read -p "Remove configuration and data? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Removing configuration and data...${NC}"
    rm -rf "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"
fi

# Remove user
if id "$SERVICE_USER" &>/dev/null; then
    read -p "Remove service user '$SERVICE_USER'? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        userdel "$SERVICE_USER"
    fi
fi

echo ""
echo -e "${GREEN}Uninstallation complete!${NC}"
