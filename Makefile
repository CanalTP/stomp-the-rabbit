PROJECTNAME := $(shell basename "$(PWD)")

GO ?= go
GOARCH := amd64

## `make clean`: Clean build files. Runs `go clean` internally.
.PHONY: clean
clean:
	$(info >  Cleaning build cache...)
	@$(GO) clean

## `make build`: Build and install the binary in the current directory
.PHONY: build
build:
	$(info >  Building binary for linux)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) $(GO) build ./cmd/$(PROJECTNAME)

## `make run`: Run the program
.PHONY: run
run:
	$(info >  Running binary)
	GOOS=linux GOARCH=$(GOARCH) $(GO) run -race ./cmd/$(PROJECTNAME)


.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
