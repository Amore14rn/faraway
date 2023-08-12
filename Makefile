# === install mod ===
install:
	go mod download

# === Tests ===
test:
	go clean --testcache
	go test ./...

# === Start (server/client) ===
start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go

start:
	docker-compose up --abort-on-container-exit --force-recreate --build server --build client

# === mod ===
.PHONY:mod
mod:
	go mod tidy && go mod vendor


# === fmt ===
.PHONY: fmt
fmt:
	find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./internal/gen/*" -exec goimports -l -w {} \;

# === Linter ===
.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(GOLANCI_LINT) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) $(GOLANCI_LINT_VERSION)

.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml
