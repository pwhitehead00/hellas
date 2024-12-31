# hellas

## Developing Cycle

### Requirements

* [docker desktop](https://www.docker.com/products/docker-desktop/)
* [kind](https://kind.sigs.k8s.io/)

### Getting Started

Run `make dev` to create a kind cluster, install and configure cert manager and
Hellas. A self signed cert is automatically created. `make clean` will delete
the test cluster.

Running `make debug` will create a debug pod with the latest version of
Terraform. Run `make debug-test` to test `terraform init` with a minimal
Terraform configuration file. `make debug-clean` to remove test artifacts.

## TODO

* set up logging
