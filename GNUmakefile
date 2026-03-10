default: build

.PHONY: build
build: bin/terraform-provider-tree

bin/terraform-provider-tree:
	go build -a -ldflags '-s -extldflags "-static"' -o bin/terraform-provider-tree

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=10 ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir . -provider-name tree
