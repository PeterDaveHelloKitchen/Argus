# Needs to be defined before including Makefile.common to auto-generate targets
DOCKER_ARCHS ?= amd64 armv7 arm64 ppc64le s390x

UI_PATH = web/ui

GOLANGCI_LINT_OPTS ?= --timeout 4m

include Makefile.common

DOCKER_IMAGE_NAME ?= argus

.PHONY: go-build
go-build: common-build

.PHONY: go-test
go-test:
	go test --tags testing,unit,integration ./...

.PHONY: go-test-coverage
go-test-coverage:
	go test ./...  -coverpkg=./... -coverprofile ./coverage.out --tags testing,unit,integration
	go tool cover -func ./coverage.out

.PHONY: web-install
web-install:
	cd $(UI_PATH) && npx update-browserslist-db@latest && npm install

.PHONY: web-build
web-build:
	cd $(UI_PATH) && PUBLIC_URL=. npm run build

.PHONY: web-test
web-test:
	cd $(UI_PATH) && npm run test:coverage

.PHONY: web-lint
web-lint:
	cd $(UI_PATH) && npm run lint

.PHONY: web
web: web-install web-build

.PHONY: build
build: web common-build

.PHONY: build-all
build-all: compress-web build-darwin build-freebsd build-linux build-openbsd build-windows