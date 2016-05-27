all: exec

exec: deps parser.go
	go build

deps:
	go get -d -v

parser.go: parser.go.y
	go tool yacc -o parser.go parser.go.y

test: parser.go
	go test -v -cover ./...

ci-test: parser.go
	go test -v -covermode=count -coverprofile=coverage.out
	$(HOME)/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $(COVERALLS_TOKEN)

examples := $(wildcard example/*.sc)
destfiles := $(patsubst example/%.sc,example/%.s,$(examples))
example: $(destfiles)

example/%.s: example/%.sc exec
	./small-c $< > $@

report: report/1.pdf report/2.pdf report/final.pdf

report/%.pdf: %.dvi
	dvipdfmx -o $@ $<

%.dvi: report/%.tex
	platex $<

.PHONY: test ci-test
