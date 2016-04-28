all: deps build

build: parser.go
	go build

deps:
	go get -d -v

parser.go: parser.go.y
	go tool yacc -o parser.go parser.go.y

test: parser.go
	go test -cover ./...

report: report/1.pdf

report/%.pdf: %.dvi
	dvipdfmx -o $@ $<

%.dvi: report/%.tex
	platex $<

.PHONY: test
