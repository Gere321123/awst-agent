#!/bin/bash

echo "AWST Agent Installation"

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root: sudo $0"
    exit 1
fi

# Installation paths
INSTALL_DIR="/opt/awst-agent"
BINARY_NAME="awst-agent"
SERVICE_NAME="awst-agent"

# Create directories
mkdir -p $INSTALL_DIR
mkdir -p /etc/awst

# Copy agent binary (assuming it's already built)
if [ -f "./build/$BINARY_NAME" ]; then
    cp "./build/$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
elif [ -f "./$BINARY_NAME" ]; then
    cp "./$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Error: $BINARY_NAME not found. Build it first: make build"
    exit 1
fi

# Create systemd service file
cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=AWST Agent Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/$BINARY_NAME
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Reload and start service
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

echo ""
echo "✓ AWST Agent successfully installed!"
echo "Service status:"
systemctl status $SERVICE_NAME --no-pager
echo ""
echo "Useful commands:"
echo "  sudo systemctl status $SERVICE_NAME  # Check status"
echo "  sudo systemctl restart $SERVICE_NAME # Restart"
echo "  sudo journalctl -u $SERVICE_NAME -f   # View logs"
echo ""
echo "Configuration file: /etc/awst/config.json"