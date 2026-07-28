# HOST_OS is the platform make is running on, which is not necessarily
# the same as the platform we are building for (that's OS).
HOST_UNAME := $(shell uname -s 2>/dev/null)
ifeq ($(HOST_UNAME),Darwin)
HOST_OS := darwin
else ifeq ($(HOST_UNAME),Linux)
HOST_OS := linux
else ifneq (,$(filter MINGW% MSYS% CYGWIN%,$(HOST_UNAME)))
HOST_OS := windows
else ifeq ($(OS),Windows_NT)
# Native Windows make has no uname, but Windows always sets OS=Windows_NT.
HOST_OS := windows
else
HOST_OS := unknown
endif

.PHONY: build-ironrdp-wasm
build-ironrdp-wasm:  ensure-llvm
	@echo "host os: $(HOST_OS), $(OS), $(HOST_UNAME)"
	$(CC) --version
	$(AR) --version
	cargo build --target wasm32-unknown-unknown


.PHONY: ensure-llvm
ifeq ($(HOST_OS),darwin)
BREW_DIR = $(shell brew --prefix)
LLVM_PREFIX = $(shell brew list | grep llvm | head -n 1)
LLVM_DIR = $(shell brew --prefix $(LLVM_PREFIX))
# Prevent these from being exported and expanded for every recipe.
unexport BREW_DIR LLVM_PREFIX LLVM_DIR

# The ironrdp WASM build needs clang/llvm-ar.
# These are applied as target-specific variables so
# brew is only invoked when necessary and so that
# CC and AR are only overwritten for WASM compilation.
build-ironrdp-wasm: override CC = $(LLVM_DIR)/bin/clang
build-ironrdp-wasm: override AR = $(LLVM_DIR)/bin/llvm-ar

ensure-llvm:
	@echo "ensure-llvm brew"
	@if [[ "$(BREW_DIR)" = "$(LLVM_DIR)" ]]; then \
		echo "llvm is required, please run 'brew install llvm' and add '/opt/homebrew/opt/llvm/bin' at the start of PATH variable"; \
		exit 1; \
	fi

else ifeq ($(HOST_OS),windows)
LLVM_DIR=$(shell vswhere.exe -latest -requires Microsoft.VisualStudio.Component.VC.Llvm.Clang -property installationPath)
unexport LLVM_DIR
build-ironrdp-wasm: override CC = $(LLVM_DIR)/VC/Tools/Llvm/x64/bin/clang
build-ironrdp-wasm: override AR = $(LLVM_DIR)/VC/Tools/Llvm/x64/bin/llvm-ar

ensure-llvm:
	@echo "ensure-llvm windows"
	@if [[ "x" = "x$(LLVM_DIR)" ]]; then \
		echo "llvm is required, please install Visual Studio with LLVM component"; \
		exit 1; \
	fi

else
ensure-llvm:
	@echo "ensure-llvm noop $(OS), $(HOST_OS)"
endif
