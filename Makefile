default: test

deps:
	@go list -f '{{ join .Imports "\n"}}' ./... | xargs -n1 go get -d

testdeps: deps
	@go list -f '{{ join .TestImports "\n"}}' ./... | xargs -n1 go get -d
	@go test -i ./...
	@go build -o /tmp/ping cmds/ping/ping.go

test: testdeps
	@go test ${TESTARGS} ./...

docs:
	@go get golang.org/x/tools/cmd/godoc
	godoc -http=":6060"

.PHONY : deps testdeps test
