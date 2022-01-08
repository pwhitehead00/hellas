MODULE = $(shell go list -m)
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

.PHONY: up
up:
	./easyrsa.sh
	minikube start --kubernetes-version 1.22.4
	kubectl apply -f ./kubernetes/terraform.yaml
	kubectl cp server.crt terraform:/etc/ssl/certs/
	skaffold dev

.PHONY: clean
clean:
	minikube delete

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/server
