#!/bin/bash

# Exit on any error (but continue through cleanup steps)
# set -e is not used here to allow cleanup to continue even if some steps fail

# Set color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Please run this script as root${NC}"
    exit 1
fi

echo -e "${GREEN}Starting WProxy uninstallation...${NC}"

# Stop service (don't fail if service doesn't exist)
echo "Stopping service..."
if systemctl is-active --quiet wproxy; then
    systemctl stop wproxy || echo -e "${YELLOW}Warning: Failed to stop service${NC}"
else
    echo "Service is not running"
fi

# Disable service (don't fail if service doesn't exist)
echo "Disabling service..."
if systemctl is-enabled --quiet wproxy 2>/dev/null; then
    systemctl disable wproxy || echo -e "${YELLOW}Warning: Failed to disable service${NC}"
else
    echo "Service is not enabled"
fi

# Remove service file
echo "Removing service file..."
if [ -f /etc/systemd/system/wproxy.service ]; then
    rm -f /etc/systemd/system/wproxy.service
else
    echo "Service file not found"
fi

# Reload systemd
systemctl daemon-reload 2>/dev/null || true

# Remove binary file
echo "Removing program files..."
if [ -f /usr/local/bin/wproxy ]; then
    rm -f /usr/local/bin/wproxy
else
    echo "Binary file not found"
fi

# Remove config files (with confirmation for safety)
echo "Removing configuration files..."
if [ -d /etc/wproxy ]; then
    rm -rf /etc/wproxy
else
    echo "Config directory not found"
fi

echo -e "${GREEN}Uninstallation completed!${NC}"