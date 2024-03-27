BUILD_DIR = ./build/
PRI_DIR = ./internal/
PUB_DIR = ./pkg/
BIN_DIR = ./bin/
RELEASE_DIR = ./release/
EXEC = karalis

#Get OS and configure based on OS
ifeq ($(OS),Windows_NT)
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CFLAGS += -D AMD64
   		DISTRO = windows64
   	else
	    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
	        CFLAGS += -D AMD64
			DISTRO = windows64
	    endif
	    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
	        CFLAGS += -D IA32 WIN32
			DISTRO = windows32
	    endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CFLAGS += -D LINUX
   		DISTRO = linux
    endif
    ifeq ($(UNAME_S),Darwin)
        CFLAGS += -D OSX
   		DISTRO = mac
    endif
    ifeq ($(UNAME),Solaris)
   		DISTRO = solaris
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CFLAGS += -D ARM
    endif
endif

.PHONY: run
#: Starts the project
run: build .deps
	@$(BIN_DIR)$(EXEC)

.PHONY: build
#: Performs a clean run of the project
build: .dev-deps $(PRI_DIR)** $(PUB_DIR)**
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
	@if [ "linux" = "$(DISTRO)" ]; then sudo apt-get install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libwayland-dev libxkbcommon-dev; fi;
	@if [ "linux" = "$(DISTRO)" ]; then sudo apt-get install -y libgl-dev libx11-dev xorg-dev libxxf86vm-dev; fi;
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
