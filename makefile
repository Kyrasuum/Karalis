BUILD_DIR = ./build/
PRI_DIR = ./internal/
PUB_DIR = ./pkg/
BIN_DIR = ./bin/
RELEASE_DIR = ./release/
EXEC = karalis

.PHONY: run
#: Starts the project
run: build .deps
	@$(BIN_DIR)$(EXEC)

.PHONY: build
#: Performs a clean run of the project
build: .dev-deps $(PRI_DIR)** $(PUB_DIR)**
	@go mod tidy -compat=1.21
	@go build -o $(BIN_DIR)$(EXEC) cmd/main.go

.PHONY: release
#: packages release target
release: build .deps

.PHONY: clean
#: Cleans build files from project
clean:
	@rm $(BIN_DIR)$(EXEC) || true;

.PHONY: clean-all
#: Cleans slate for project
clean-all:
	@rm .dev-deps || true;
	@rm .deps || true;

# deps include target
.PHONY: deps
.deps:
	$(MAKE) --no-print-directory deps

#: Install dependencies for running this project
deps:
	@touch .deps

# dev-deps include target
.PHONY: dev-deps
.dev-deps:
	$(MAKE) --no-print-directory dev-deps

#: Install dependencies for compiling targets in this makefile
dev-deps: .deps
	@sudo apt-get install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libwayland-dev libxkbcommon-dev
	@sudo apt-get install -y libgl-dev libx11-dev xorg-dev libxxf86vm-dev
	@touch .dev-deps

.PHONY: help
#: Lists available commands
help:
	@echo "Available Commands for project:"
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
	 | grep -v -- -- \
	 | sed 'N;s/\n/###/' \
	 | sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
	 | column -t  -s '###'
