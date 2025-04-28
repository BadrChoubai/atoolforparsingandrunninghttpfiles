BINS ?= cli
DIST = .dist

all:
	go build


clean:
	echo "nothing to clean"


SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset
.DEFAULT_GOAL = all
.PHONY: clean docker

