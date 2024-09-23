PWD = $(CURDIR)
# Имя сервиса
SERVICE_NAME := asvsoft
# Дефолтная ОС
GOOS ?= linux
# Дефолтная архитектура
GOARCH ?= arm
# Время сборки
BUILD_DATE = $(shell TZ=UTC-3 date +%Y-%m-%dT%H:%M)
# Ветка
BRANCH := $(shell git symbolic-ref -q --short HEAD)
# 8 символов последнего коммита
LAST_COMMIT_HASH = $(shell git rev-parse HEAD | cut -c -8)
# ld флаги
LD_FLAGS := "-X 'main.BuildTime=$(BUILD_DATE)' -X 'main.BuildCommit=$(LAST_COMMIT_HASH)' -X 'main.BuildBranch=$(BRANCH)'"
# Путь до бинарника golang-ci
GOLANGCI_BIN := $(shell which golangci-lint)
# Путь до бинарника go/tinygo
GO_BIN := go

# Дефолтное поведение
default: build

# Линтер проверяет отличия от мастера
.PHONY: lint
lint:
	$(GOLANGCI_BIN) run --config=.golangci.yml --new-from-rev=origin/master ./...
	@echo "lint successfully"

# Линтер проверяет полностью весь код сервиса
.PHONY: full-lint
full-lint:
	$(GOLANGCI_BIN) run --config=.golangci.yml ./...
	@echo "lint successfully"

# Сборка сервиса
.PHONY: build
build: asvsoft
asvsoft:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BIN) build -o bin/$(SERVICE_NAME) -ldflags=$(LD_FLAGS) $(PWD)/cmd/$(SERVICE_NAME)
	@echo "build successfully"

# Деплой сервиса на плату
.PHONY: deploy-%
deploy-%: build
	scp $(PWD)/bin/$(SERVICE_NAME) $*:/usr/local/bin/
	@echo "deploy successfully"
