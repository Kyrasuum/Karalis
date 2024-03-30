BUILD_DIR = ./build/
PRI_DIR = ./internal/
PUB_DIR = ./pkg/
BIN_DIR = ./bin/
RELEASE_DIR = ./release/
EXEC = karalis
ZIG ?= false

#Get OS and configure based on OS
ifeq ($(OS),Windows_NT)
    DISTRO ?= windows
    LDFLAGS ?= -ldflags='-H=windowsgui'
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        ARCH ?= amd64
		CFLAGS ?= x86_64-windows-gnu -I/usr/local/include -L/usr/local/lib
    else
		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		    ARCH ?= amd64
		    CFLAGS ?= x86_64-windows-gnu -I/usr/local/include -L/usr/local/lib
		endif
		ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		    ARCH ?= ia32
		endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
   		DISTRO ?= linux
    endif
    ifeq ($(UNAME_S),Darwin)
   		DISTRO ?= mac
    endif
    ifeq ($(UNAME),Solaris)
	   	DISTRO ?= solaris
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        ARCH ?= amd64
		CFLAGS ?= x86_64-linux-gnu -isystem /usr/include -I/usr/include/GL -I/usr/local/include -L/usr/lib/x86_64-linux-gnu -L/usr/local/lib
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        ARCH ?= ia32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        ARCH ?= arm64
		CFLAGS ?= aarch64-linux-gnu -isystem /usr/include -I/usr/include/GL -I/usr/local/include -L/usr/lib/aarch64-linux-gnu -L/usr/local/lib
    endif
endif

.PHONY: run
#: Starts the project
run: build .deps
	@$(BIN_DIR)$(EXEC)

.PHONY: build
# check if zig enabled and linux or windows os
ifneq (,$(filter $(ZIG)$(DISTRO),truelinux truewindows))
# zig build
build: .dev-deps $(PRI_DIR)** $(PUB_DIR)**
	@CC="zig cc -target $(CFLAGS)" \
	CXX="zig c++ -target $(CFLAGS)" \
	CGO_ENABLED=1 \
	GOOS=$(DISTRO) \
	GOARCH=$(ARCH) \
	go build -o $(BIN_DIR)$(EXEC) $(LDFLAGS) cmd/main.go
else
#: Performs a clean run of the project
build: .dev-deps $(PRI_DIR)** $(PUB_DIR)**
	@go build -o $(BIN_DIR)$(EXEC) cmd/main.go
endif

build-wasm:
	@ZIG=$(ZIG) \
	DISTRO=js \
	ARCH=wasm \
	$(MAKE) --no-print-directory build
build-mac:
	@ZIG=$(ZIG) \
	DISTRO=mac \
	ARCH=amd64 \
	$(MAKE) --no-print-directory build
build-linux-amd64:
	@ZIG=$(ZIG) \
	DISTRO=linux \
	ARCH=amd64 \
	CFLAGS="x86_64-linux-gnu -isystem /usr/include -L/usr/lib/x86_64-linux-gnu" \
	$(MAKE) --no-print-directory build
build-linux-arm64:
	@ZIG=$(ZIG) \
	DISTRO=linux \
	ARCH=arm64 \
	CFLAGS="aarch64-linux-gnu -isystem /usr/include -L/usr/lib/aarch64-linux-gnu" \
	$(MAKE) --no-print-directory build
build-windows-amd64:
	@ZIG=$(ZIG) \
	DISTRO=windows \
	ARCH=amd64 \
	CFLAGS="x86_64-windows-gnu" \
	LDFLAGS="-ldflags='-H=windowsgui'" \
	$(MAKE) --no-print-directory build


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
	@$(MAKE) --no-print-directory deps

#: Install dependencies for running this project
deps:
	@touch .deps

# dev-deps include target
.PHONY: dev-deps
.dev-deps:
	@$(MAKE) --no-print-directory dev-deps

# dev-deps for linux
ifeq ($(DISTRO),linux)
ifeq ($(ARCH),amd64)
dev-deps: .dev-deps-linux-amd64
.PHONY: .dev-deps-linux-amd64
.dev-deps-linux-amd64:
	@sudo apt-get install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libwayland-dev libxkbcommon-dev
	@sudo apt-get install -y libgl-dev libx11-dev xorg-dev libxxf86vm-dev
endif
ifeq ($(ARCH),arm64)
dev-deps: .dev-deps-linux-arm64
.PHONY: .dev-deps-linux-arm64
.dev-deps-linux-arm64:
	@sudo dpkg --add-architecture arm64
	@sudo apt-get install -y libgl1-mesa-dev:arm64 libxi-dev:arm64 libxcursor-dev:arm64 libxrandr-dev:arm64 libxinerama-dev:arm64 libwayland-dev:arm64 libxkbcommon-dev:arm64
	@sudo apt-get install -y libgl-dev:arm64 libx11-dev:arm64 xorg-dev:arm64 libxxf86vm-dev:arm64
endif
endif

#: Install dependencies for compiling targets in this makefile
dev-deps: .deps
	@go mod tidy
	@go get -v -u github.com/gen2brain/raylib-go/raylib
	@touch .dev-deps

.PHONY: test
#: Perfrom unit tests for application
test:
	@gofmt -e -w .
	@go test ./... -cover || true;

.PHONY: help
#: Lists available commands
help:
	@echo "Available Commands for project:"
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
	 | grep -v -- -- \
	 | sed 'N;s/\n/###/' \
	 | sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
	 | column -t  -s '###'
