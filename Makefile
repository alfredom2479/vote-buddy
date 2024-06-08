build:
	go build -o bin/vote-buddy

run: build
		./bin/vote-buddy

test:
	go test -v ./...