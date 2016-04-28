all: goget app

app: parser.go
	go build

goget:
	go get -v ./...

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
