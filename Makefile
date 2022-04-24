.PHONY: help

help: Makefile
	@sed -n 's/^##//p' $<


## test	: run unit tests
test:
	go test -v -race \
	-covermode=atomic -coverprofile=coverage.out \
	$$(go list ./... | grep -v cmd)
	go tool cover -func coverage.out | grep total | awk '{print "coverage: " $$3}'

## lint	: run linter›
lint:
	@golangci-lint --version
	@golangci-lint cache clean
	@golangci-lint run -v

## fmt	: check code formatting
fmt:
	@gofmt -w -l $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi


## run	: start test server
run :
	source ./.env.local &&\
	go run ./cmd/gophermart/main.go