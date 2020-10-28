PROJECTNAME := $(shell basename "$(PWD)")

GO ?= go
GOARCH := amd64

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
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) $(GO) build ./cmd/$(PROJECTNAME)


.PHONY: run
run: ## Run the program
	$(info > Running binary)
	GOOS=linux GOARCH=$(GOARCH) $(GO) run -race ./cmd/$(PROJECTNAME)

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
	$(info > Logout from the Kisio Docker registry)
	docker logout ${KISIO_DOCKER_REGISTRY}

.PHONY: build_docker_image
build_docker_image: docker_login ## Build docker image
	$(info > Build docker image)
	@deploy/scripts/build.sh $(if $(findstring true,${DRY_RUN}), --dry-run)

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
