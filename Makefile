GO ?= GO111MODULE=on CGO_ENABLED=1 go

run-backend:
	$(GO) run main.go

build-backend:
	$(GO) build -o bootstrap main.go
