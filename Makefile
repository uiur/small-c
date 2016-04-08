main: parser
	go build

parser:
	go tool yacc -o parser.go parser.go.y

test:
	go test -cover ./...

.PHONY: test
