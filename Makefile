default: test

deps:
	@go list -f '{{ join .Imports "\n"}}' ./... | xargs -n1 go get -d

testdeps: deps
	@go list -f '{{ join .TestImports "\n"}}' ./... | xargs -n1 go get -d
	@go test -i ./...

test: testdeps
	@GOBIN=/usr/local/bin go install cmds/ping.go
	@go test ${TESTARGS} ./...

docs:
	@go get golang.org/x/tools/cmd/godoc
	godoc -http=":6060"

.PHONY : deps testdeps test
