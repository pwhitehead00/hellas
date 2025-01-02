.PHONY: dev
dev: create-cluster image-load cert-manager install-hellas

.PHONY: redeploy
redeploy: image-load
	kubectl -n hellas rollout restart deployment hellas

.PHONY: create-cluster
create-cluster:
	kind create cluster

.PHONY: cert-manager
cert-manager:
	helm repo add jetstack https://charts.jetstack.io --force-update
	helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --set crds.enabled=true --wait
	kubectl apply -f hack/self-signed.yaml --wait
	kubectl create ns hellas
	kubectl apply -f hack/cert.yaml

.PHONY: install-hellas
install-hellas:
	helm upgrade --install --namespace hellas --create-namespace hellas ./charts/hellas

.PHONY: clean
clean:
	kind delete cluster

.PHONY: docker-build-dev
docker-build-dev:
	# CGO_ENABLED=0 docker build -t hellas:latest .
	docker build -t hellas:latest .

.PHONY: image-load
image-load: docker-build-dev
	kind load docker-image hellas:latest

.PHONY: debug
debug:
	kubectl -n hellas create cm init-script --from-file=./hack/init.sh
	kubectl -n hellas apply -f ./hack/debug.yaml

.PHONY: debug-test
debug-test:
	kubectl -n hellas exec debug -- env TF_LOG=DEBUG terraform -chdir=/terraform init -upgrade

.PHONY: debug-test-errors
debug-test-errors:
	kubectl -n hellas exec debug -- env TF_LOG=DEBUG terraform -chdir=/errors init -upgrade

.PHONY: debug-clean
debug-clean:
	kubectl -n hellas delete cm/init-script po/debug
