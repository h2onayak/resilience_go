CURRENT_DIR=$(shell pwd)
ALL_PACKAGES=$(shell go list -mod=mod ./... | grep -v /vendor)
SOURCE_DIRS=$(shell go list -mod=mod ./... | grep -v /vendor | grep -v /out | cut -d "/" -f5 | uniq)

SHELL := /bin/bash # Use bash syntax

check: fmt lint vet

fmtcheck:
	@gofmt -l -s $(SOURCE_DIRS) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

fmt:
	@echo $(SOURCE_DIRS)
	gofmt -l -s -w $(SOURCE_DIRS)

lint:
	@if [[ `golint $(ALL_PACKAGES) | \
		{ grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } | \
		{ grep -vwE "type name will be used as [a-zA-Z]+[.][a-zA-Z]+ by other packages, and that stutters; consider calling this Client" || true; } | \
		wc -l | tr -d ' '` -ne 0 ]]; then \
			golint $(ALL_PACKAGES) | \
			{ grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } | \
			{ grep -vwE "type name will be used as [a-zA-Z.]+[.][a-zA-Z]+ by other packages, and that stutters; consider calling this Client" || true; }; \
          exit 2; \
    fi;

vet:
	@go vet ./...
