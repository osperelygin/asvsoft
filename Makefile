PWD = $(CURDIR)
# Имя сервиса
SERVICE_NAME := asvsoft
# Директория с локальными бинарниками
LOCAL_BIN = $(PWD)/bin
# Путь до бинарника сервиса
SERVICE_BIN := $(LOCAL_BIN)/$(SERVICE_NAME)
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
# Версия golang-ci
GOLANGCI_TAG := v1.60.1
# Путь до бинарника golang-ci
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
# Версия go-arch-lint
GOARCH_TAG := v1.11.6
# Путь до бинарника go-arch-lint
GOARCH_BIN := $(LOCAL_BIN)/go-arch-lint

# Дефолтное поведение
default: build

# Запуск тестов и подсчет процента покрытия тестами
.PHONY: test
test:
	go test -parallel=10 -cover -coverprofile coverage.out ./internal/pkg/proto/...
	@echo "test passed"

# Cоздание отчета о покрытии тестами
.PHONY: cover
cover: test
	go tool cover -html=coverage.out
	@echo "cover executed"

# Запуск бенчмарков
.PHONY: bench
bench:
	go test -bench=. -benchmem ./internal/pkg/proto/...
	@echo "bench executed"

# Сборка сервиса
.PHONY: build
build: asvsoft
asvsoft:
	@echo "=================================================="
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(SERVICE_BIN) -ldflags=$(LD_FLAGS) $(PWD)/cmd/$(SERVICE_NAME)
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

# Уставнавливаем golangci-lint в локальную диру с бинарниками
bin/golangci-lint:
	GOPROXY="" GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_TAG)

# Линтер проверяет отличия от мастера
.PHONY: lint
lint: bin/golangci-lint
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOLANGCI_BIN) run --config=.golangci.yml --new-from-rev=origin/master ./...
	@echo "lint successfully"

# Линтер проверяет полностью весь код сервиса
.PHONY: full-lint
full-lint: bin/golangci-lint
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOLANGCI_BIN) run --config=.golangci.yml ./...
	@echo "lint successfully"

# Уставнавливаем go-arch-lint в локальную диру с бинарниками
bin/go-arch-lint:
	GOPROXY="" GOBIN=$(LOCAL_BIN) go install github.com/fe3dback/go-arch-lint@$(GOARCH_TAG)

# Cтроит в docs/arch.svg архитектуру сервиса
.PHONY: arch-graph
arch-graph: bin/go-arch-lint
	$(GOARCH_BIN) graph --out docs/arch.svg

# Линтовка структуры кода
.PHONY: arch-lint
arch-lint: bin/go-arch-lint
	$(GOARCH_BIN) check
