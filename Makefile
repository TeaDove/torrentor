GO ?= GO111MODULE=on CGO_ENABLED=0 go

run:
	$(GO) run ${TARGET}

build:
	$(GO) build -o bootstrap ${TARGET}

tests:
	$(GO) test ./... -count=1 -p=100

docker-run:
	docker build --build-arg TARGET=${TARGET} -t telemetry-go . && docker run -it telemetry-go

run-tests:
	go test ./... -count=1 -p=100

update-all:
	go get -u ./...
	go mod tidy

show-pprof:
	go tool pprof -web cpu.prof

lint:
	gofumpt -w *.go
	golines --base-formatter=gofumpt --max-len=120 --no-reformat-tags -w .
	wsl --fix ./...
	golangci-lint run --fix
