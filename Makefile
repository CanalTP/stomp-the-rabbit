VERSION := $(shell git describe --tag --always --dirty)

GO ?= go
GOARCH := amd64
CGO_ENABLE := 1
STATIC :=
LDFLAGS := -X github.com/CanalTP/stomptherabbit.Version=$(VERSION)

ifeq ($(STATIC),1)
LDFLAGS += -w -s
CGO_ENABLE := 0
TAGS := netgo
endif

# used for building docker image only
DRY_RUN ?= false
ifneq (${DRY_RUN},true)
ifneq (${DRY_RUN},false)
$(error DRY_RUN must be 'true' or 'false')
endif
endif

# determine if session is interactive or not
INTERACTIVE:=$(shell tty -s && echo 1)

.PHONY: clean
clean: ## Clean build files. Runs `go clean` internally.
	$(info > Cleaning build cache...)
	@$(GO) clean

.PHONY: build
build: ## Build and install the binary in the current directory
	$(info >  Building binary for linux)
	CGO_ENABLED=$(CGO_ENABLE) $(GO) build -a -tags "$(TAGS)" -ldflags "$(LDFLAGS)" ./cmd/stomptherabbit


.PHONY: run
run: ## Run the program
	$(info > Running binary)
	GOOS=linux GOARCH=$(GOARCH) $(GO) run -race ./cmd/stomptherabbit

.PHONY: docker_login
docker_login: ## Login to dockerhub
ifdef INTERACTIVE
	$(info > Login to dockerhub)
	docker login
else
	$(info > Login skipped, use CI credentials instead)
endif

.PHONY: docker_logout
docker_logout: ## Logout from dockerhub
	$(info > Logout from dockerhub)
	docker logout

.PHONY: release
release: docker_login ## Release docker image
	$(info > Build docker image)
	@deploy/scripts/release.sh $(if $(findstring true,${DRY_RUN}), --dry-run)

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
