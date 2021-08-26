GOARCH=amd64

.PHONY: build
build: linux windows darwin

.PHONY: linux
linux: linux-api linux-eventstore linux-viewbuilder

.PHONY: windows
windows: windows-api windows-eventstore windows-viewbuilder

.PHONY: darwin
darwin: darwin-api darwin-eventstore darwin-viewbuilder

.PHONY: api
api: linux-api windows-api darwin-api

.PHONY: linux-api
linux-api:
	GOOS=linux go build -ldflags="-s -w" -o bin/api-linux-${GOARCH} ./cmd/api/main.go

.PHONY: windows-api
windows-api:
	GOOS=windows go build -ldflags="-s -w" -o bin/api-windows-${GOARCH}.exe ./cmd/api/main.go

.PHONY: darwin-api
darwin-api:
	GOOS=darwin go build -ldflags="-s -w" -o bin/api-darwin-${GOARCH} ./cmd/api/main.go

.PHONY: eventstore
eventstore: linux-eventstore windows-eventstore darwin-eventstore

.PHONY: linux-eventstore
linux-eventstore:
	GOOS=linux go build -ldflags="-s -w" -o bin/eventstore-linux-${GOARCH} ./cmd/eventstore/main.go

.PHONY: windows-eventstore
windows-eventstore:
	GOOS=windows go build -ldflags="-s -w" -o bin/eventstore-windows-${GOARCH}.exe ./cmd/eventstore/main.go

.PHONY: darwin-eventstore
darwin-eventstore:
	GOOS=darwin go build -ldflags="-s -w" -o bin/eventstore-darwin-${GOARCH} ./cmd/api/main.go

.PHONY: viewbuilder
viewbuilder: linux-viewbuilder windows-viewbuilder darwin-viewbuilder

.PHONY: linux-viewbuilder
linux-viewbuilder:
	GOOS=linux go build -ldflags="-s -w" -o bin/viewbuilder-linux-${GOARCH} ./cmd/viewbuilder/main.go

.PHONY: windows-viewbuilder
windows-viewbuilder:
	GOOS=windows go build -ldflags="-s -w" -o bin/viewbuilder-windows-${GOARCH}.exe ./cmd/viewbuilder/main.go

.PHONY: darwin-viewbuilder
darwin-viewbuilder:
	GOOS=darwin go build -ldflags="-s -w" -o bin/viewbuilder-darwin-${GOARCH} ./cmd/api/main.go

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./internal/...

.PHONY: test
test:
	go test -v -race -count=1 ./...

.PHONY: deps
deps:
	go mod verify && \
	go mod tidy

.PHONY: genapi
genapi:
	go generate -v ./...