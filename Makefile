MODULE = $(shell go list -m)
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

up:
	minikube start --kubernetes-version 1.22.4
	skaffold dev

clean:
	skaffold delete
	minikube delete

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/server
