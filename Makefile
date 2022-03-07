.PHONY: all run coordinator docker-compose-up docker-compose-down

# Go Flags
GOFLAGS ?= $(GOFLAGS:)
# We need to export GOBIN to allow it to be set
# for processes spawned from the Makefile
export GOBIN ?= $(PWD)/bin
GO=go
LDFLAGS="-s -w"

# ENV
COORDINATOR_CMD="./cmd/coordinator"
COORDINATOR_OUT="coordinator"


all: run # Alias

run: coordinator

coordinator:
	@echo Starting Coordinator
	
	# $(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(COORDINATOR_CMD)
	$(GO) run $(GOFLAGS) $(COORDINATOR_CMD)

build_coordinator:
	@echo Building Coordinator

	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(COORDINATOR_OUT) $(COORDINATOR_CMD)

dev-compose: docker-compose-down docker-compose-up coordinator

docker-compose-up:
	@docker-compose up -d

docker-compose-clean:
	@docker-compose down
