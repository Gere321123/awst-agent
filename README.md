# AWST Agent

A reverse proxy agent that connects to the central server for authentication.

## ⚠️ Prototype Notice

> This is a **prototype**. The central website is still under development. Features may change.

## Quick Install

```bash
# Clone and build
git clone https://github.com/Gere321123/awst-agent.git
cd awst-agent
make build
sudo make install

# Run
sudo awst-agent

## Service Commands

sudo systemctl start awst-agent   # Start
sudo systemctl stop awst-agent    # Stop
sudo systemctl restart awst-agent # Restart
sudo systemctl status awst-agent  # Status
sudo journalctl -u awst-agent -f  # View logs

# Configuration
## Edit .env file to change server URLs:

CENTRAL_SERVER=http://localhost:5000   # Change to your server
BACKEND_SERVER=http://localhost:3000
PROXY_PORT=:8080

# Uninstall

sudo make uninstall