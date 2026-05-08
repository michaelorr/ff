BINARY := ff

.PHONY: build clean lint test run watch

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

lint:
	golantci-lint run ./...

test:
	go test ./...

watch:
	find . -name '*.go' | entr -cs 'go build -o $(BINARY) .'
