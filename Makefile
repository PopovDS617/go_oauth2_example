BINARY_NAME=app
HTTP_PORT="8000"
LOCAL_BIN:=$(CURDIR)/bin
include .env

## getdeps: installs dependencies
getdeps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ./bin/${BINARY_NAME} ./cmd
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@env HTTP_PORT=${HTTP_PORT} GOOGLE_OAUTH_CLIENT_ID=${GOOGLE_OAUTH_CLIENT_ID} GOOGLE_OAUTH_CLIENT_SECRET=${GOOGLE_OAUTH_CLIENT_SECRET} ./bin/${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ./bin/${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## race: checks for data race
race:
	cd ./cmd && \
	export HTTP_PORT=${HTTP_PORT} && \
	go run -race .

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./bin/${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...

## testrace: checks tests for data race
testrace:
	go test -race -v ./...

## testcov: collects test coverage
testcov:
	go test -coverprofile=coverage.out -v ./...

## testout: displays test coverage as html in browser
testout:
	go tool cover -html=coverage.out

## du: starts all the docker-compose containers in detached mode
du:
	docker-compose up -d

## dd: stops all docker-compose containers
dd :
	docker-compose down

## lint: runs linter for all files in the project
lint:
	cd ./bin && \
	golangci-lint run ../...