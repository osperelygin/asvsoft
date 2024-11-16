PWD = $(CURDIR)
# Имя сервиса
SERVICE_NAME := asvsoft
# Путь до бинарника сервиса
SERVICE_BIN := $(PWD)/bin/$(SERVICE_NAME)
# Дефолтная ОС
GOOS ?= linux
# Дефолтная архитектура
GOARCH ?= arm
# Время сборки
BUILD_DATE = $(shell TZ=UTC-3 date +%Y-%m-%dT%H:%M:%S)
# Ветка
BRANCH := $(shell git symbolic-ref -q --short HEAD)
# 8 символов последнего коммита
LAST_COMMIT_HASH = $(shell git rev-parse HEAD | cut -c -8)
# ld флаги
LD_FLAGS := "-X 'main.BuildTime=$(BUILD_DATE)' -X 'main.BuildCommit=$(LAST_COMMIT_HASH)' -X 'main.BuildBranch=$(BRANCH)'"
# Путь до бинарника golang-ci
GOLANGCI_BIN := $(shell which golangci-lint)
# Путь до бинарника go
GO_BIN ?= go

# Дефолтное поведение
default: build

# Линтер проверяет отличия от мастера
.PHONY: lint
lint:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOLANGCI_BIN) run --config=.golangci.yml --new-from-rev=origin/master ./...
	@echo "lint successfully"

# Запуск тестов и подсчет процента покрытия тестами
.PHONY: test
test:
	$(GO_BIN) test -parallel=10 -cover -coverprofile coverage.out ./internal/pkg/proto/...
	@echo "test passed"

# Cоздание отчета о покрытии тестами
.PHONY: cover
cover: test
	go tool cover -html=coverage.out
	@echo "cover executed"

# Запуск бенчмарков
.PHONY: bench
bench:
	$(GO_BIN) test -bench=. -benchmem ./internal/pkg/proto/...
	@echo "bench executed"

# Линтер проверяет полностью весь код сервиса
.PHONY: full-lint
full-lint:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOLANGCI_BIN) run --config=.golangci.yml ./...
	@echo "lint successfully"

# Сборка сервиса
.PHONY: build
build: asvsoft
asvsoft:
	@echo "=================================================="
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BIN) build -o $(SERVICE_BIN) -ldflags=$(LD_FLAGS) $(PWD)/cmd/$(SERVICE_NAME)
	@echo "build successfully"

# Деплой опредленного сервиса на плату
.PHONY: deploy-%
deploy-%: build
	@echo "=================================================="
	./scripts/deploy/ssh_deploy.sh "$(SERVICE_BIN)" "$*"
	@echo "deploy successfully"

# Деплой нескольких сервисов на плату
.PHONY: deploy
deploy: build
	@echo "=================================================="
	./scripts/deploy/ssh_deploy.sh "$(SERVICE_BIN)" "${SSH_HOST_LIST}"
	@echo "deploy successfully"
