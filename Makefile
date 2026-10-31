.PHONY: build server run test vet fmt check clean

# build the game server binary
build:
	@go build -o bin/server .

# build and start the game server on :4000
server: build
	@./bin/server

# start without building a binary
run:
	@go run .

test:
	@go test ./...

vet:
	@go vet ./...

fmt:
	@gofmt -l -w .

# what to run before committing
check: fmt vet test

clean:
	@rm -rf bin
