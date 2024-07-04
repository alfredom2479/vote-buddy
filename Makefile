build:
	go build -o bin/vote-buddy cmd/main/main.go

run: build
		./bin/vote-buddy

test:
	go test -v ./...