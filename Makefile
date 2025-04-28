# This is a file that should include a GITHUB_TOKEN variable for ensuring goreleaser
# can successfully publish packages to the repo
include .env.release
BINS ?= cli
DIST = .dist

all:



VERSION ?= $(shell git describe --tags --always --dirty)
# Used internally.  Users should pass GOOS and/or GOARCH.
OS := $(if $(GOOS),$(GOOS),$(shell GOTOOLCHAIN=local go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell GOTOOLCHAIN=local go env GOARCH))

TAG := $(VERSION)__$(OS)_$(ARCH)
release:
	echo "release $(TAG)"
	goreleaser


version: # @HELP outputs the version string
version:
	@echo $(VERSION)


clean:
	@echo "nothing to clean"


SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset
.DEFAULT_GOAL = all
.PHONY: clean version

