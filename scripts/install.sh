#!/bin/bash

# AWST Agent Installer from GitHub
set -e

REPO="Gere321123/awst-agent"
INSTALL_DIR="/opt/awst-agent"
BINARY_NAME="awst-agent"
SERVICE_NAME="awst-agent"

echo "========================================"
echo "   AWST Agent Installation from GitHub"
echo "========================================"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root: sudo $0"
    exit 1
fi

# Create installation directory
mkdir -p $INSTALL_DIR
mkdir -p /etc/awst

# Clone or pull the repository
if [ -d "$INSTALL_DIR/.git" ]; then
    echo "Updating existing installation..."
    cd $INSTALL_DIR
    git pull
else
    echo "Cloning repository from GitHub..."
    git clone https://github.com/$REPO.git $INSTALL_DIR
    cd $INSTALL_DIR
fi

# Build the agent
echo "Building agent..."
make build

# Copy binary
cp build/$BINARY_NAME /usr/local/bin/
chmod +x /usr/local/bin/$BINARY_NAME

# Create .env file if it doesn't exist
if [ ! -f "$INSTALL_DIR/.env" ]; then
    echo "Creating .env template..."
    cat > $INSTALL_DIR/.env << EOF
# Central Server Configuration
CENTRAL_SERVER=http://localhost:5000

# Proxy Configuration
BACKEND_SERVER=http://localhost:3000
PROXY_PORT=:8080
EOF
    echo "✓ .env file created at $INSTALL_DIR/.env"
    echo "  Edit this file to change server configuration"
fi

# Create systemd service
cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=AWST Agent Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=/usr/local/bin/$BINARY_NAME
Restart=on-failure
RestartSec=10
EnvironmentFile=$INSTALL_DIR/.env

[Install]
WantedBy=multi-user.target
EOF

# Reload and start service
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl restart $SERVICE_NAME

echo ""
echo "✓ AWST Agent successfully installed!"
echo ""
echo "Service status:"
systemctl status $SERVICE_NAME --no-pager
echo ""
echo "Configuration file: /etc/awst/config.json"
echo "Environment file: $INSTALL_DIR/.env"
echo ""
echo "Useful commands:"
echo "  sudo systemctl status $SERVICE_NAME  # Check status"
echo "  sudo systemctl restart $SERVICE_NAME # Restart"
echo "  sudo journalctl -u $SERVICE_NAME -f  # View logs"
echo "  sudo nano $INSTALL_DIR/.env          # Edit configuration"