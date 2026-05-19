.PHONY: build install uninstall clean

BINARY=Agent
BUILD_DIR=build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY) ./cmd/agent
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY)"

install: build
	@echo "Installing..."
	sudo cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/
	sudo mkdir -p /etc/awst
	@echo "✓ Installation complete"

uninstall:
	@echo "Uninstalling..."
	sudo rm -f /usr/local/bin/$(BINARY)
	sudo rm -rf /etc/awst
	sudo systemctl stop awst-agent 2>/dev/null || true
	sudo systemctl disable awst-agent 2>/dev/null || true
	sudo rm -f /etc/systemd/system/awst-agent.service
	sudo systemctl daemon-reload
	@echo "✓ Uninstall complete"

clean:
	rm -rf $(BUILD_DIR)