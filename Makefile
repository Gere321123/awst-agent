.PHONY: build install uninstall clean run test

BINARY=awst-agent
BUILD_DIR=build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go mod download
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/agent
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY)"

run: build
	@echo "Running agent..."
	./$(BUILD_DIR)/$(BINARY)

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

test:
	go test -v ./...

docker-build:
	docker build -t awst-agent .

docker-run:
	docker run -d --name awst-agent --network host awst-agent