default: build

build:
	GOTOOLCHAIN=1.23.7 go build -a -ldflags '-s -extldflags "-static"'

fmt:
	gofmt -s -w -e .

test:
	GOTOOLCHAIN=go1.23.7 go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
