# test: run all tests
test:
	go test -v -race -buildvcs ./...

# test/cover: run all tests and display coverage
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

tidy:
	go mod tidy -v
	go fmt ./...
