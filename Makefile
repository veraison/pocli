# Copyright 2023 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL := all

all: build
.PHONY: all

build: pocli
.PHONY: build

pocli: *.go cmd/*.go
	go build -o pocli

clean:
	rm -f pocli
.PHONY: clean
