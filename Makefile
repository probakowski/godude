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

.PHONY: build
build:
	@echo "host os: $(HOST_OS), $(OS), $(HOST_UNAME)"
