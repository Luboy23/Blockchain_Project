build:
	go build -o ./bin/Blockchain_Project

run: build
	./bin/Blockchain_Project

test:
	go test ./...
